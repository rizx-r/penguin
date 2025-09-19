package svc

import (
	"penguin/apps/im/im_models"
	"penguin/apps/im/ws/internal/config"
	"penguin/apps/task/mq/mqclient"
)

type (
	ServiceContext struct {
		Config config.Config
		im_models.ChatLogModel
		mqclient.MsgChatTransferClient
		mqclient.MsgReadTransferClient
	}
)

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:                c,
		ChatLogModel:          im_models.MustChatLogModel(c.Mongo.Url, c.Mongo.Db),
		MsgChatTransferClient: mqclient.NewMsgChatTransferClient(c.MsgChatTransfer.Addrs, c.MsgChatTransfer.Topic),
		MsgReadTransferClient: mqclient.NewMsgReadTransferClient(c.MsgReadTransfer.Addrs, c.MsgReadTransfer.Topic),
	}
}
