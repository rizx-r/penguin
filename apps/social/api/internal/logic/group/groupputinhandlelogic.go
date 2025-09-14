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

type GroupPutInHandleLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewGroupPutInHandleLogic 申请进群处理
func NewGroupPutInHandleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupPutInHandleLogic {
	return &GroupPutInHandleLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GroupPutInHandleLogic) GroupPutInHandle(req *types.GroupPutInHandleRep) (resp *types.GroupPutInHandleResp, err error) {
	uid := ctxdata.GetUid(l.ctx)
	res, err := l.svcCtx.Social.GroupPutInHandle(l.ctx, &socialclient.GroupPutInHandleReq{
		GroupReqId:   req.GroupReqId,
		GroupId:      req.GroupId,
		HandleUid:    uid,
		HandleResult: req.HandleResult,
	})

	if constants.HandlerResult(req.HandleResult) != constants.PassHandleResult {
		return
	}

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
