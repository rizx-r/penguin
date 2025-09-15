package mq

import "penguin/pkg/constants"

type (
	MsgChatTransfer struct {
		ConversationId     string `json:"conversationId"`
		constants.ChatType `json:"chatType"`
		SendId             string   `json:"sendId"`
		RecvId             string   `json:"recvId"`
		RecvIds            []string `json:"recvIds"`
		SendTime           int64    `json:"sendTime"`
		constants.MType    `json:"mType"`
		Content            string `json:"content"`
	}

	// MsgMarkRead 处理已读的消费者
	MsgMarkRead struct {
		constants.ChatType `json:",chatType"`
		ConversationId     string   `json:"conversationId"` // 会话id
		SendId             string   `json:"sendId"`         // 发送者id
		RecvId             string   `json:"recvId"`         // 接受者id
		MsgIds             []string `json:"msgIds"`         // 消息id
	}
)
