package group

import (
	"context"
	"github.com/jinzhu/copier"
	"penguin/apps/social/rpc/socialclient"
	"penguin/pkg/ctxdata"

	"penguin/apps/social/api/internal/svc"
	"penguin/apps/social/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GroupListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 用户申群列表
func NewGroupListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupListLogic {
	return &GroupListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// GroupList 获取群列表
func (l *GroupListLogic) GroupList(req *types.GroupListRep) (resp *types.GroupListResp, err error) {
	uid := ctxdata.GetUid(l.ctx)
	list, err := l.svcCtx.Social.GroupList(l.ctx, &socialclient.GroupListReq{
		UserId: uid,
	})
	if err != nil {
		return nil, err
	}

	var respList []*types.Groups
	err = copier.Copy(&respList, &list.List)
	if err != nil {
		return nil, err
	}

	return &types.GroupListResp{List: respList}, nil
}
