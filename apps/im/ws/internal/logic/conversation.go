package logic

import (
	"context"
	"penguin/apps/im/im_models"
	"penguin/apps/im/ws/internal/svc"
	"penguin/apps/im/ws/websocket"
	"penguin/apps/im/ws/ws"
	"penguin/pkg/wuid"
	"time"
)

type (
	Conversation struct {
		ctx context.Context
		srv *websocket.Server
		svc *svc.ServiceContext
	}
)

func NewConversation(ctx context.Context, srv *websocket.Server, svc *svc.ServiceContext) *Conversation {
	return &Conversation{
		ctx: ctx,
		srv: srv,
		svc: svc,
	}
}

func (c *Conversation) SingleChat(data *ws.Chat, userId string) error {
	if data.ConversationId == "" {
		data.ConversationId = wuid.CombineId(userId, data.RecvId)
	}

	// 记录消息
	chatLog := im_models.ChatLog{
		ConversationID: data.ConversationId,
		SendID:         userId,
		RecvID:         data.RecvId,
		ChatType:       data.ChatType,
		MsgFrom:        0,
		MsgType:        data.MType,
		MsgContent:     data.Content,
		SendTime:       time.Now().UnixNano(),
	}

	err := c.svc.ChatLogModel.Insert(c.ctx, &chatLog)
	return err
}
