package group

import (
	"context"
	"github.com/jinzhu/copier"
	"github.com/pkg/errors"
	"penguin/pkg/xerr"

	"penguin/apps/social/rpc/internal/svc"
	"penguin/apps/social/rpc/social"

	"github.com/zeromicro/go-zero/core/logx"
)

type GroupUsersLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGroupUsersLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupUsersLogic {
	return &GroupUsersLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// GroupUsers 获取群成员列表
func (l *GroupUsersLogic) GroupUsers(in *social.GroupUsersReq) (*social.GroupUsersResp, error) {
	groupMembers, err := l.svcCtx.GroupMembersModel.ListByGroupId(l.ctx, in.GroupId)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewDBErr(), "list group member err: %v, req: %v", err, in)
	}
	var respList []*social.GroupMembers
	if err := copier.Copy(&respList, &groupMembers); err != nil {
		return nil, errors.Wrapf(err, "copy group members err: %v, req: %v", err, in)
	}
	return &social.GroupUsersResp{List: respList}, nil
}
