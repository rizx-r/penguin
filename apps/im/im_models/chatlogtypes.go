package im_models

import (
	"penguin/pkg/constants"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type (
	ChatLog struct {
		ID             primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"` // !
		ConversationID string             `bson:"conversation_id"`                    // bson 是 MongoDB 使用的一种二进制 JSON（Binary JSON）格式。
		SendID         string             `bson:"send_id"`
		RecvID         string             `bson:"recv_id"`
		MsgFrom        int                `bson:"msg_from"`
		ChatType       constants.ChatType `bson:"chat_type"`
		MsgType        constants.MType    `bson:"msg_type"`
		MsgContent     string             `bson:"msg_content"`
		SendTime       int64              `bson:"send_time"`
		Status         int                `bson:"status"`
		UpdateAt       time.Time          `bson:"update_at,omitempty" json:"update_at,omitempty"` // omitempty: 当字段的值是该类型的零值时，序列化（转成 bson/json）时会自动忽略这个字段，不输出到结果里。
		CreateAt       time.Time          `bson:"create_at,omitempty" json:"create_at,omitempty"`
	}
)
