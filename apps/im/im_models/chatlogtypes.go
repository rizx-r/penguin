package im_models

import (
	"penguin/pkg/constants"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	DefaultChatLogLimit int64 = 100
)

type (
	ChatLog struct {
		ID             primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"` // !
		ConversationID string             `bson:"conversationId"`                     // bson 是 MongoDB 使用的一种二进制 JSON（Binary JSON）格式。
		SendID         string             `bson:"sendId"`
		RecvID         string             `bson:"recvId"`
		MsgFrom        int                `bson:"msgFrom"`
		ChatType       constants.ChatType `bson:"chatType"`
		MsgType        constants.MType    `bson:"msgType"`
		MsgContent     string             `bson:"msgContent"`
		SendTime       int64              `bson:"sendTime"`
		Status         int                `bson:"status"`
		ReadRecords    []byte             `bson:"readRecords"` // 记录已读

		UpdateAt time.Time `bson:"updateAt,omitempty" json:"update_at,omitempty"` // omitempty: 当字段的值是该类型的零值时，序列化（转成 bson/json）时会自动忽略这个字段，不输出到结果里。
		CreateAt time.Time `bson:"createAt,omitempty" json:"create_at,omitempty"`
	}
)
