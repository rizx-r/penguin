package group

import (
	"context"
	"github.com/jinzhu/copier"
	"penguin/apps/social/rpc/socialclient"

	"penguin/apps/social/api/internal/svc"
	"penguin/apps/social/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GroupPutInListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 申请进群列表
func NewGroupPutInListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupPutInListLogic {
	return &GroupPutInListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GroupPutInListLogic) GroupPutInList(req *types.GroupPutInListRep) (resp *types.GroupPutInListResp, err error) {
	list, err := l.svcCtx.Social.GroupPutinList(l.ctx, &socialclient.GroupPutinListReq{
		GroupId: req.GroupId,
	})
	if err != nil {
		return nil, err
	}
	var respList []*types.GroupRequests
	if err = copier.Copy(&respList, &list.List); err != nil {
		return nil, err
	}
	return &types.GroupPutInListResp{List: respList}, nil
}
