package msgTransfer

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/zeromicro/go-zero/core/logx"
	"penguin/apps/im/im_models"
	"penguin/apps/im/ws/websocket"
	"penguin/apps/task/mq/internal/svc"
	"penguin/apps/task/mq/mq"
	"penguin/pkg/constants"
	"time"
)

type (
	MsgChatTransfer struct {
		logx.Logger
		svc *svc.ServiceContext
	}
)

func NewMsgChatTransfer(svcContext *svc.ServiceContext) *MsgChatTransfer {
	return &MsgChatTransfer{
		Logger: logx.WithContext(context.Background()),
		svc:    svcContext,
	}
}

func (m *MsgChatTransfer) Consume(ctx context.Context, key, value string) error {
	fmt.Printf("key: %s, value: %s\n", key, value)
	var (
		data mq.MsgChatTransfer
		//ctx  = context.Background()
	)

	if err := json.Unmarshal([]byte(value), &data); err != nil {
		return err
	}

	// 记录数据
	if err := m.addChatLog(ctx, &data); err != nil {
		return err
	}

	// 推送消息
	return m.svc.WsClient.Send(websocket.Message{
		FrameType: websocket.FrameData,
		Method:    "push",
		FormId:    constants.SYSTEM_ROOT_UID,
		Data:      data,
	})
}

func (m *MsgChatTransfer) addChatLog(ctx context.Context, data *mq.MsgChatTransfer) error {
	// 记录消息
	chatlog := im_models.ChatLog{
		ConversationID: data.ConversationId,
		SendID:         data.SendId,
		RecvID:         data.RecvId,
		MsgFrom:        0,
		ChatType:       data.ChatType,
		MsgType:        data.MType,
		MsgContent:     data.Content,
		SendTime:       data.SendTime,
		UpdateAt:       time.Time{},
		CreateAt:       time.Time{},
	}
	return m.svc.ChatLogModel.Insert(ctx, &chatlog)
}
