package msgTransfer

import (
	"context"
	"fmt"
	"github.com/zeromicro/go-zero/core/logx"
	"penguin/apps/im/ws/websocket"
	"penguin/apps/im/ws/ws"
	"penguin/apps/social/rpc/socialclient"
	"penguin/apps/task/mq/internal/svc"
	"penguin/pkg/constants"
)

type (
	BaseMsgTransfer struct {
		svcCtx *svc.ServiceContext
		logx.Logger
	}
)

func NewBaseMsgTransfer(svcCtx *svc.ServiceContext) *BaseMsgTransfer {
	return &BaseMsgTransfer{
		svcCtx: svcCtx,
		Logger: logx.WithContext(context.Background()),
	}
}

// Transfer 转发消息
func (m *BaseMsgTransfer) Transfer(ctx context.Context, data *ws.Push) error {
	var err error

	// TODO:
	switch data.ChatType {
	case constants.GroupChatType:
		err = m.group(ctx, data)
	case constants.SingleChatType:
		err = m.single(ctx, data)
	default:
		err = fmt.Errorf("unknown chat type: %s", data.ChatType)
	}

	return err
}

// single 私聊
func (m *BaseMsgTransfer) single(ctx context.Context, data *ws.Push) error {
	return m.sendWebsocketMessage(data)
}

// group 群聊
func (m *BaseMsgTransfer) group(ctx context.Context, data *ws.Push) error {
	// 查询群用户
	users, err := m.svcCtx.Social.GroupUsers(ctx, &socialclient.GroupUsersReq{GroupId: data.RecvId})
	if err != nil {
		return err
	}
	// 构建要发送的用户列表：群用户 除了自己
	data.RecvIds = make([]string, 0, len(users.List))
	for _, user := range users.List {
		if user.UserId == data.SendId {
			continue
		}
		data.RecvIds = append(data.RecvIds, data.SendId)
	}
	// 发送消息
	return m.sendWebsocketMessage(data)
}

func (m *BaseMsgTransfer) sendWebsocketMessage(data *ws.Push) error {
	return m.svcCtx.WsClient.Send(websocket.Message{
		FrameType: websocket.FrameData,
		Method:    "push",
		FormId:    constants.SYSTEM_ROOT_UID,
		Data:      data,
	})
}
