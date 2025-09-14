package group

import (
	"context"
	"penguin/apps/im/rpc/im"
	"penguin/apps/social/rpc/socialclient"
	"penguin/pkg/constants"
	"penguin/pkg/ctxdata"

	"penguin/apps/social/api/internal/svc"
	"penguin/apps/social/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GroupPutInLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewGroupPutInLogic 申请进群
func NewGroupPutInLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupPutInLogic {
	return &GroupPutInLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GroupPutInLogic) GroupPutIn(req *types.GroupPutInRep) (resp *types.GroupPutInResp, err error) {
	uid := ctxdata.GetUid(l.ctx)
	res, err := l.svcCtx.Social.GroupPutin(l.ctx, &socialclient.GroupPutinReq{
		GroupId:    req.GroupId,
		ReqId:      uid,
		ReqMsg:     req.ReqMsg,
		ReqTime:    req.ReqTime,
		JoinSource: int32(req.JoinSource),
	})
	if err != nil || res.GroupId == "" {
		return nil, err
	}
	_, err = l.svcCtx.Im.SetUpUserConversation(l.ctx, &im.SetUpUserConversationReq{
		SendId:   uid,
		RecvId:   res.GroupId,
		ChatType: int32(constants.GroupChatType),
	})
	return nil, err
}
