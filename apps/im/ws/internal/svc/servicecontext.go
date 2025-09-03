package svc

import (
	"penguin/apps/im/im_models"
	"penguin/apps/im/ws/internal/config"
	"penguin/apps/task/mq/mqclient"
)

type (
	ServiceContext struct {
		Config config.Config
		mqclient.MsgChatTransferClient
		im_models.ChatLogModel
	}
)

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:                c,
		MsgChatTransferClient: mqclient.NewMsgChatTransferClient(c.MsgChatTransfer.Addrs, c.MsgChatTransfer.Topic),
		ChatLogModel:          im_models.MustChatLogModel(c.Mongo.Url, c.Mongo.Db),
	}
}
