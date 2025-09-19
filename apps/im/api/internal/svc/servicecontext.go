package svc

import (
	"github.com/zeromicro/go-zero/zrpc"
	"penguin/apps/im/api/internal/config"
	"penguin/apps/im/rpc/imclient"
	"penguin/apps/social/rpc/socialclient"
	"penguin/apps/user/rpc/userclient"
)

type ServiceContext struct {
	Config config.Config

	imclient.Im
	userclient.User
	socialclient.Social
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config: c,
		// RPC客户端
		Im:     imclient.NewIm(zrpc.MustNewClient(c.ImRpc)),
		User:   userclient.NewUser(zrpc.MustNewClient(c.UserRpc)),
		Social: socialclient.NewSocial(zrpc.MustNewClient(c.SocialRpc)),
	}
}
