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

type GroupPutinListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGroupPutinListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupPutinListLogic {
	return &GroupPutinListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// GroupPutinList 获取加群申请列表
func (l *GroupPutinListLogic) GroupPutinList(in *social.GroupPutinListReq) (*social.GroupPutinListResp, error) {
	groupReqs, err := l.svcCtx.GroupRequestsModel.ListNoHandler(l.ctx, in.GroupId)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewDBErr(), "list group req err: %v, req: %v", err, in)
	}

	var respList []*social.GroupRequests
	if err := copier.Copy(&respList, &groupReqs); err != nil {
		return nil, err
	}
	return &social.GroupPutinListResp{
		List: respList,
	}, nil
}
