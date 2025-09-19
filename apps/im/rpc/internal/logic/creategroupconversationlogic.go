package logic

import (
	"context"
	"github.com/pkg/errors"
	"penguin/apps/im/im_models"
	"penguin/apps/im/rpc/im"
	"penguin/apps/im/rpc/internal/svc"
	"penguin/pkg/constants"
	"penguin/pkg/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateGroupConversationLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateGroupConversationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateGroupConversationLogic {
	return &CreateGroupConversationLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 创建群聊
func (l *CreateGroupConversationLogic) CreateGroupConversation(in *im.CreateGroupConversationReq) (*im.CreateGroupConversationResp, error) {
	// todo: add your logic here and delete this line

	res := &im.CreateGroupConversationResp{}

	// 是否已经存在群会话了
	_, err := l.svcCtx.ConversationModel.FindOne(l.ctx, in.GroupId)
	if err == nil {
		return res, nil
	}
	if err != im_models.ErrNotFound {
		return nil, errors.Wrapf(xerr.NewDBErr(), "ConversationModel.FindOne error at CreateGroupConversation: %v, req: %v", err, in)
	}

	// 添加群会话
	err = l.svcCtx.ConversationModel.Insert(l.ctx, &im_models.Conversation{
		ConversationId: in.GroupId,
		ChatType:       constants.GroupChatType,
	})

	if err != nil {
		return nil, errors.Wrapf(err, "ConversationModel.Insert error at CreateGroupConversation: %v, req: %v", err, in)
	}

	_, err = NewSetUpUserConversationLogic(l.ctx, l.svcCtx).SetUpUserConversation(&im.SetUpUserConversationReq{
		SendId:   in.CreateId,
		RecvId:   in.GroupId,
		ChatType: int32(constants.GroupChatType),
	})
	return res, err

	return &im.CreateGroupConversationResp{}, nil
}
