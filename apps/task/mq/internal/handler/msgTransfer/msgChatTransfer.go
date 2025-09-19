package msgTransfer

import (
	"context"
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"penguin/apps/im/im_models"
	"penguin/apps/im/ws/ws"
	"penguin/apps/task/mq/internal/svc"
	"penguin/apps/task/mq/mq"
	"penguin/pkg/bitmap"
	"time"
)

type (
	MsgChatTransfer struct {
		*BaseMsgTransfer
	}
)

func NewMsgChatTransfer(svcContext *svc.ServiceContext) *MsgChatTransfer {
	return &MsgChatTransfer{
		BaseMsgTransfer: NewBaseMsgTransfer(svcContext),
	}
}

func (m *MsgChatTransfer) Consume(ctx context.Context, key, value string) error {
	fmt.Printf("[Consume] key: %s, value: %s\n", key, value)
	var (
		data mq.MsgChatTransfer
		//ctx  = context.Background()
		msgId = primitive.NewObjectID()
	)

	if err := json.Unmarshal([]byte(value), &data); err != nil {
		return err
	}

	// 记录数据
	if err := m.addChatLog(ctx, msgId, &data); err != nil {
		return err
	}

	return m.Transfer(ctx, &ws.Push{
		ConversationId: data.ConversationId,
		ChatType:       data.ChatType,
		SendId:         data.SendId,
		RecvId:         data.RecvId,
		RecvIds:        data.RecvIds,
		SendTime:       data.SendTime,
		MType:          data.MType,
		MsgId:          msgId.Hex(),
		Content:        data.Content,
	})
}

func (m *MsgChatTransfer) addChatLog(ctx context.Context, msgId primitive.ObjectID, data *mq.MsgChatTransfer) error {
	// 记录消息
	chatlog := im_models.ChatLog{
		ConversationID: data.ConversationId,
		ID:             msgId,
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

	readRecords := bitmap.NewBitmap(0)
	readRecords.Set(chatlog.SendID)
	chatlog.ReadRecords = readRecords.Export()

	err := m.svcCtx.ChatLogModel.Insert(ctx, &chatlog)
	if err != nil {
		return err
	}
	return m.svcCtx.ConversationModel.UpdateMsg(ctx, &chatlog)
}
