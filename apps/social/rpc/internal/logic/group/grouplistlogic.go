package group

import (
	"context"
	"github.com/jinzhu/copier"
	"github.com/pkg/errors"
	"penguin/apps/social/rpc/internal/svc"
	"penguin/apps/social/rpc/social"
	"penguin/pkg/xerr"
	"strconv"

	"github.com/zeromicro/go-zero/core/logx"
)

type GroupListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGroupListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupListLogic {
	return &GroupListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// GroupList 获取用户的群列表
func (l *GroupListLogic) GroupList(in *social.GroupListReq) (*social.GroupListResp, error) {
	userGroup, err := l.svcCtx.GroupMembersModel.ListByUserId(l.ctx, in.UserId)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewDBErr(), "list group member err %v req %v", err, in.UserId)
	}
	if len(userGroup) == 0 {
		return &social.GroupListResp{}, nil
	}

	gids := make([]string, 0, len(userGroup))
	for _, g := range userGroup {
		gids = append(gids, strconv.FormatUint(g.Id, 10)) // string(g.Id)
	}

	groups, err := l.svcCtx.GroupsModel.ListByGroupIds(l.ctx, gids)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewDBErr(), " list group err: %v,, req: %v", err, in)
	}

	var respList []*social.Groups
	if err := copier.Copy(&respList, &groups); err != nil {
		return nil, err
	}

	return &social.GroupListResp{
		List: respList,
	}, nil
}
