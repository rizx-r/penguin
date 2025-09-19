package mqclient

import (
	"context"
	"encoding/json"
	"github.com/zeromicro/go-queue/kq"
	"penguin/apps/task/mq/mq"
)

type (
	MsgChatTransferClient interface {
		Push(msg *mq.MsgChatTransfer) error
	}
	msgChatTransferClient struct {
		pusher *kq.Pusher
	}
	MsgReadTransferClient interface {
		Push(msg *mq.MsgMarkRead) error
	}
	msgReadTransferClient struct {
		pusher *kq.Pusher
	}
)

func NewMsgChatTransferClient(addr []string, topic string, opts ...kq.PushOption) MsgChatTransferClient {
	return &msgChatTransferClient{
		pusher: kq.NewPusher(addr, topic),
	}
}

func (c *msgChatTransferClient) Push(msg *mq.MsgChatTransfer) error {
	body, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	return c.pusher.Push(context.Background(), string(body)) // ?
}

func NewMsgReadTransferClient(addr []string, topic string, opts ...kq.PushOption) MsgReadTransferClient {
	return &msgReadTransferClient{
		pusher: kq.NewPusher(addr, topic),
	}
}

func (c *msgReadTransferClient) Push(msg *mq.MsgMarkRead) error {
	body, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	return c.pusher.Push(context.Background(), string(body)) // ?
}
