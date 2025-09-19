package svc

import (
	"github.com/zeromicro/go-zero/zrpc"
	"penguin/apps/im/rpc/imclient"
	"penguin/apps/social/api/internal/config"
	"penguin/apps/social/rpc/socialclient"
	"penguin/apps/user/rpc/userclient"
)

type ServiceContext struct {
	Config config.Config

	socialclient.Social
	userclient.User
	imclient.Im
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config: c,
		Social: socialclient.NewSocial(zrpc.MustNewClient(c.SocialRpc)),
		User:   userclient.NewUser(zrpc.MustNewClient(c.UserRpc)),
		Im:     imclient.NewIm(zrpc.MustNewClient(c.ImRpc)),
	}
}
