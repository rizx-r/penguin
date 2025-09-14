package group

import (
	"context"
	"penguin/apps/im/rpc/imclient"
	"penguin/apps/social/rpc/socialclient"
	"penguin/pkg/ctxdata"

	"penguin/apps/social/api/internal/svc"
	"penguin/apps/social/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateGroupLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewCreateGroupLogic 创群
func NewCreateGroupLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateGroupLogic {
	return &CreateGroupLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// CreateGroup 创建群聊的业务逻辑
func (l *CreateGroupLogic) CreateGroup(req *types.GroupCreateReq) (resp *types.GroupCreateResp, err error) {
	uid := ctxdata.GetUid(l.ctx)

	// 创建群
	res, err := l.svcCtx.GroupCreate(l.ctx, &socialclient.GroupCreateReq{
		Name:       req.Name,
		Icon:       req.Icon,
		CreatorUid: uid,
	})
	if err != nil || res.Id == "" {
		return nil, err
	}

	// 建立会话
	_, err = l.svcCtx.Im.CreateGroupConversation(l.ctx, &imclient.CreateGroupConversationReq{
		GroupId:  res.Id,
		CreateId: uid,
	})
	return nil, err
}
