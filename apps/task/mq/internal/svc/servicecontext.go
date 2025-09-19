package svc

import (
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/zrpc"
	"net/http"
	"penguin/apps/im/im_models"
	"penguin/apps/im/ws/websocket"
	"penguin/apps/social/rpc/socialclient"
	"penguin/apps/task/mq/internal/config"
	"penguin/pkg/constants"
)

type (
	ServiceContext struct {
		config.Config
		WsClient websocket.Client
		*redis.Redis
		socialclient.Social
		im_models.ChatLogModel
		im_models.ConversationModel
	}
)

func NewServiceContext(c config.Config) *ServiceContext {
	svc := &ServiceContext{
		Config:            c,
		Redis:             redis.MustNewRedis(c.Redisx),
		ChatLogModel:      im_models.MustChatLogModel(c.Mongo.Url, c.Mongo.Db),
		ConversationModel: im_models.MustConversationModel(c.Mongo.Url, c.Mongo.Db),
		Social:            socialclient.NewSocial(zrpc.MustNewClient(c.SocialRpc)),
	}
	token, err := svc.GetSystemToken()
	if err != nil {
		panic(err)
	}

	header := http.Header{}
	header.Set("Authorization", token)
	svc.WsClient = websocket.NewClient(c.Ws.Host, websocket.WithClientHeader(header))
	return svc
}

func (s *ServiceContext) GetSystemToken() (string, error) {
	return s.Redis.Get(constants.REDIS_SYSTEM_ROOT_TOKEN)
}
