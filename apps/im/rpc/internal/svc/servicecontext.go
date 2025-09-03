package svc

import (
	"penguin/apps/im/im_models"
	"penguin/apps/im/rpc/internal/config"
)

type ServiceContext struct {
	Config config.Config

	im_models.ChatLogModel
	im_models.ConversationsModel
	im_models.ConversationModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:             c,
		ChatLogModel:       im_models.MustChatLogModel(c.Mongo.Url, c.Mongo.Db),
		ConversationsModel: im_models.MustConversationsModel(c.Mongo.Url, c.Mongo.Db),
		ConversationModel:  im_models.MustConversationModel(c.Mongo.Url, c.Mongo.Db),
	}
}
