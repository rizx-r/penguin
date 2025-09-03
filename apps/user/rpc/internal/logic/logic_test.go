package logic

import (
	"github.com/zeromicro/go-zero/core/conf"
	"path/filepath"
	"penguin/apps/user/rpc/internal/config"
	"penguin/apps/user/rpc/internal/svc"
)

var (
	svcCtx *svc.ServiceContext
)

func init() {
	var c config.Config
	conf.MustLoad(filepath.Join("../../etc/dev/user.yaml"), &c)
	svcCtx = svc.NewServiceContext(c)
}
