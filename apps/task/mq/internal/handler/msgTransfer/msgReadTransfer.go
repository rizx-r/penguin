package msgTransfer

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"github.com/rogpeppe/go-internal/cache"
	"github.com/zeromicro/go-queue/kq"
	"penguin/apps/im/ws/ws"
	"penguin/apps/task/mq/internal/svc"
	"penguin/apps/task/mq/mq"
	"penguin/pkg/bitmap"
	"penguin/pkg/constants"
	_default "penguin/pkg/default"
	"sync"
	"time"
)

type (
	MsgReadTransfer struct {
		*BaseMsgTransfer

		cache.Cache
		mu        sync.Mutex
		groupMsgs map[string]*GroupMsgRead

		push chan *ws.Push
	}
)

var (
	GroupMsgReadRecordDelayTime  = _default.GroupMsgReadRecordDelayTime
	GroupMsgReadRecordDelayCount = _default.GroupMsgReadRecordDelayCount
)

func NewMsgReadTransfer(svcCtx *svc.ServiceContext) kq.ConsumeHandler {
	m := &MsgReadTransfer{
		BaseMsgTransfer: NewBaseMsgTransfer(svcCtx),
		groupMsgs:       make(map[string]*GroupMsgRead),
		push:            make(chan *ws.Push),
	}
	if svcCtx.Config.MsgReadHandler.GroupMsgReadHandler != constants.GroupMsgReadHandlerAtTransfer {
		GroupMsgReadRecordDelayCount = svcCtx.Config.MsgReadHandler.GroupMsgReadRecordDelayTime
	}
	if svcCtx.Config.MsgReadHandler.GroupMsgReadRecordDelayTime > 0 {
		GroupMsgReadRecordDelayTime = time.Duration(svcCtx.Config.MsgReadHandler.GroupMsgReadRecordDelayTime) * time.Second
	}

	go m.transfer()
	return m
}

func (receiver *MsgReadTransfer) Consume(ctx context.Context, key, value string) error {
	receiver.Infof("=>[MsgReadTransfer] key: %s, value: %s\n", key, value)
	var (
		data mq.MsgMarkRead
	)
	if err := json.Unmarshal([]byte(value), &data); err != nil {
		return err
	}

	// 更新
	readRecords, err := receiver.UpdateChatLogRead(ctx, &data)
	if err != nil {
		return err
	}

	push := &ws.Push{
		ConversationId: data.ConversationId,
		ChatType:       data.ChatType,
		SendId:         data.SendId,
		RecvId:         data.RecvId,
		ContentType:    constants.ContentMarkRead,
		ReadRecords:    readRecords,
	}

	switch data.ChatType {
	case constants.SingleChatType:
		// 直接推送
		receiver.push <- push
	case constants.GroupChatType:
		// 判断是否开启消息合并
		if receiver.svcCtx.MsgReadHandler.GroupMsgReadHandler == constants.GroupMsgReadHandlerAtTransfer {
			receiver.push <- push
		}
		receiver.mu.Lock()
		defer receiver.mu.Unlock()

		push.SendId = ""
		// 判断是否有记录
		if _, ok := receiver.groupMsgs[push.ConversationId]; ok {
			receiver.Infof("merge push %v", push.ConversationId)
			// 合并请求
			receiver.groupMsgs[push.ConversationId].mergePush(push)
		} else {
			// 创建记录
			receiver.Infof("new GroupMsgRead push %v", push.ConversationId)
			receiver.groupMsgs[push.ConversationId] = NewGroupMsgRead(push, receiver.push)
		}
	}

	return nil
}

func (receiver *MsgReadTransfer) UpdateChatLogRead(ctx context.Context, data *mq.MsgMarkRead) (map[string]string, error) {
	res := make(map[string]string)
	chatLogs, err := receiver.svcCtx.ChatLogModel.ListByMsgIds(ctx, data.MsgIds)
	if err != nil {
		return nil, err
	}
	// 处理已读
	for _, chatLog := range chatLogs {
		switch chatLog.ChatType {
		case constants.SingleChatType:
			chatLog.ReadRecords = []byte{1}
		case constants.GroupChatType:
			readRecords := bitmap.LoadBitmap(chatLog.ReadRecords)
			readRecords.Set(data.SendId)
			chatLog.ReadRecords = readRecords.Export()
		default:
			panic("func [UpdateChatLogRead]: unknow chatLog.ChatType")
		}
		// 记录已读
		res[chatLog.ID.Hex()] = base64.StdEncoding.EncodeToString(chatLog.ReadRecords)
		err = receiver.svcCtx.ChatLogModel.UpdateMarkRead(ctx, chatLog.ID, chatLog.ReadRecords)
		if err != nil {
			return nil, err
		}
	}
	return res, nil
}

func (receiver *MsgReadTransfer) transfer() {
	for push := range receiver.push {
		if push.RecvId != "" || len(push.RecvIds) > 0 {
			if err := receiver.Transfer(context.Background(), push); err != nil {
				receiver.Errorf("MsgReadTransfer err: %v, push %v", err, push)
			}
		}
		if push.ChatType == constants.SingleChatType {
			// 不处理私聊
			continue
		}
		if receiver.svcCtx.Config.MsgReadHandler.GroupMsgReadHandler != constants.GroupMsgReadHandlerAtTransfer {
			// 如果设置了不处理
			continue
		}
		// 清空数据
		receiver.mu.Lock()
		if _, ok := receiver.groupMsgs[push.ConversationId]; ok && receiver.groupMsgs[push.ConversationId].IsIdle() {
			receiver.groupMsgs[push.ConversationId].clear()
			delete(receiver.groupMsgs, push.ConversationId)
		}
	}
}
