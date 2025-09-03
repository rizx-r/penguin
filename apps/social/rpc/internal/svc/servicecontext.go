package svc

import (
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"penguin/apps/social/rpc/internal/config"
	"penguin/apps/social/social_models"
)

type ServiceContext struct {
	Config config.Config

	social_models.FriendsModel
	social_models.FriendRequestsModel
	social_models.GroupsModel
	social_models.GroupRequestsModel
	social_models.GroupMembersModel
}

func NewServiceContext(c config.Config) *ServiceContext {

	sqlConn := sqlx.NewMysql(c.Mysql.Datasource)

	return &ServiceContext{
		Config: c,

		FriendsModel:        social_models.NewFriendsModel(sqlConn, c.Cache),
		FriendRequestsModel: social_models.NewFriendRequestsModel(sqlConn, c.Cache),
		GroupsModel:         social_models.NewGroupsModel(sqlConn, c.Cache),
		GroupRequestsModel:  social_models.NewGroupRequestsModel(sqlConn, c.Cache),
		GroupMembersModel:   social_models.NewGroupMembersModel(sqlConn, c.Cache),
	}
}
