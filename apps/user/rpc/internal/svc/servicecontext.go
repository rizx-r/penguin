package svc

import (
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"penguin/apps/user/models"
	"penguin/apps/user/rpc/internal/config"
	"penguin/pkg/constants"
	"penguin/pkg/ctxdata"
	"time"
)

type ServiceContext struct {
	Config config.Config
	*redis.Redis
	models.UsersModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	sqlConn := sqlx.NewMysql(c.Mysql.Datasource)

	return &ServiceContext{
		Config:     c,
		Redis:      redis.MustNewRedis(c.Redisx),
		UsersModel: models.NewUsersModel(sqlConn, c.Cache),
	}
}

func (s *ServiceContext) SetRootToken() error {
	// 生成jwt
	systemToken, err := ctxdata.GetJwtFromToken(s.Config.Jwt.AccessSecret, time.Now().Unix(), 999999999, constants.SYSTEM_ROOT_UID)
	if err != nil {
		return err
	}
	// 写入到redis
	return s.Redis.Set(constants.REDIS_SYSTEM_ROOT_TOKEN, systemToken)
}
