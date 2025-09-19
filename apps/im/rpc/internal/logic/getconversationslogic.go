package logic

import (
	"context"
	"fmt"
	"github.com/jinzhu/copier"
	"github.com/pkg/errors"
	"penguin/apps/im/im_models"
	"penguin/pkg/xerr"

	"penguin/apps/im/rpc/im"
	"penguin/apps/im/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetConversationsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetConversationsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetConversationsLogic {
	return &GetConversationsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取会话
func (l *GetConversationsLogic) GetConversations(in *im.GetConversationsReq) (*im.GetConversationsResp, error) {
	// todo: add your logic here and delete this line
	fmt.Println("GetConversations: ", in.UserId)
	// 根据用户查询会话列表
	data, err := l.svcCtx.ConversationsModel.FindByUserId(l.ctx, in.UserId)
	if err != nil {
		if err == im_models.ErrNotFound {
			// 如果还未和任何人创建会话
			fmt.Println("user %v have no conversation list", in.UserId)
			return &im.GetConversationsResp{}, nil
		}
		return nil, errors.Wrapf(xerr.NewDBErr(), "ConversationsModel.FindByUserId error at GetConversations: %v, req: %v", err, in)
	}

	var res im.GetConversationsResp

	if err = copier.Copy(&res, &data); err != nil {
		fmt.Println("copier.Copy error data:", data)
		return nil, err
	}

	// 根据会话列表，查询出具体的会话
	ids := make([]string, 0, len(data.ConversationList))

	for _, conversation := range data.ConversationList {
		ids = append(ids, conversation.ConversationId)
	}

	conversations, err := l.svcCtx.ConversationModel.ListByConversationIds(l.ctx, ids)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewDBErr(), "ConversationModel.ListByConversationIds error at GetConversations: %v, req: %v", err, in)
	}

	// 计算是否存在未读消息
	for _, conversation := range conversations {
		if _, ok := res.ConversationList[conversation.ConversationId]; !ok {
			continue
		}
		// 用户读取的消息量
		total := res.ConversationList[conversation.ConversationId].Total
		if total < int32(conversation.Total) {
			// 有新的消息未被用户读到
			res.ConversationList[conversation.ConversationId].Total = int32(conversation.Total)
			// 未读消息量
			res.ConversationList[conversation.ConversationId].ToRead = int32(conversation.Total) - total
			// 更改当前会话为显示状态
			res.ConversationList[conversation.ConversationId].IsShow = true
		}
	}
	/*	for _, conversation := range conversations {
		if conv, ok := res.ConversationList[conversation.ConversationId]; ok {
			// 用户读取的消息量
			total := conv.Total
			if total < int32(conversation.Total) {
				// 有新的消息未被用户读到
				conv.Total = int32(conversation.Total)
				// 未读消息量
				conv.ToRead = int32(conversation.Total) - total
				// 更改当前会话为显示状态
				conv.IsShow = true
			}
		}
	}*/

	return &res, nil
}
