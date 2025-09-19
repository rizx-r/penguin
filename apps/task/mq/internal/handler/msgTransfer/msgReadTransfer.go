package msgTransfer

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"github.com/zeromicro/go-queue/kq"
	"penguin/apps/im/ws/ws"
	"penguin/apps/task/mq/internal/svc"
	"penguin/apps/task/mq/mq"
	"penguin/pkg/bitmap"
	"penguin/pkg/constants"
)

type (
	MsgReadTransfer struct {
		*BaseMsgTransfer
	}
)

func NewMsgReadTransfer(svcCtx *svc.ServiceContext) kq.ConsumeHandler {
	return &MsgReadTransfer{
		NewBaseMsgTransfer(svcCtx),
	}
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

	return receiver.Transfer(ctx, &ws.Push{
		ConversationId: data.ConversationId,
		ChatType:       data.ChatType,
		SendId:         data.SendId,
		RecvId:         data.RecvId,
		ContentType:    constants.ContentMarkRead,
		ReadRecords:    readRecords,
	})
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
