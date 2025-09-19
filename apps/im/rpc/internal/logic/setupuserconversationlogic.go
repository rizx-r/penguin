package logic

import (
	"context"
	"errors"
	perr "github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"penguin/apps/im/im_models"
	"penguin/apps/im/rpc/im"
	"penguin/apps/im/rpc/internal/svc"
	"penguin/pkg/constants"
	"penguin/pkg/wuid"
	"penguin/pkg/xerr"
)

type SetUpUserConversationLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSetUpUserConversationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SetUpUserConversationLogic {
	return &SetUpUserConversationLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 建立会话: 群聊, 私聊
func (l *SetUpUserConversationLogic) SetUpUserConversation(in *im.SetUpUserConversationReq) (*im.SetUpUserConversationResp, error) {
	// todo: add your logic here and delete this line

	switch constants.ChatType(in.ChatType) {
	case constants.SingleChatType:
		// 生成会话id
		conversationId := wuid.CombineId(in.SendId, in.RecvId) // "in.SendId_in.RecvId"

		// 验证是否建立过会话
		conversationRes, err := l.svcCtx.ConversationModel.FindOne(l.ctx, conversationId)
		if err != nil {
			if errors.Is(err, im_models.ErrNotFound) {
				err = l.svcCtx.ConversationModel.Insert(l.ctx, &im_models.Conversation{
					ConversationId: conversationId,
					ChatType:       constants.SingleChatType,
				})
				if err != nil {
					return nil, perr.Wrapf(xerr.NewDBErr(), "ConversationModel.Insert error at SetUpUserConversation: %v", err)
				}
			} else {
				return nil, perr.Wrapf(err, "ConversationModel.FindOne error at SetUpUserConversation: %v", err)
			}
		} else if conversationRes != nil {
			// 会话已经建立了
			return nil, nil
		}
		// 建立两者的会话
		err = l.setUpUserConversation(conversationId, in.SendId, in.RecvId, constants.SingleChatType, true)
		if err != nil {
			return nil, err
		}
		err = l.setUpUserConversation(conversationId, in.RecvId, in.SendId, constants.SingleChatType, false)
		if err != nil {
			return nil, err
		}
	case constants.GroupChatType:
		err := l.setUpUserConversation(in.RecvId, in.SendId, in.RecvId, constants.GroupChatType, true)
		if err != nil {
			return nil, err
		}
	default:
		panic("unhandled default case at SetUpUserConversation")
	}

	return &im.SetUpUserConversationResp{}, nil
}

/*
	func (l *SetUpUserConversationLogic) setUpUserConversation(conversationId, userId, recvId string, chatType constants.ChatType, isShow bool) error {
		// 用户的会话列表
		conversations, err := l.svcCtx.ConversationsModel.FindByUserId(l.ctx, conversationId)
		if err != nil {
			if errors.Is(err, im_models.ErrNotFound) {
				conversations = &im_models.Conversations{
					ID:               primitive.ObjectID{},
					UserId:           userId,
					ConversationList: make(map[string]*im_models.Conversation),
				}
			} else {
				return perr.Wrapf(xerr.NewDBErr(), "ConversationsModel.FindByUserId error at setUpUserConversation: %v, req: %v", err, userId)
			}
		}

		// 更新会话记录
		if _, ok := conversations.ConversationList[conversationId]; ok {
			return nil
		}

		// 添加会话记录
		conversations.ConversationList[conversationId] = &im_models.Conversation{
			ConversationId: conversationId,
			ChatType:       constants.SingleChatType,
			IsShow:         isShow,
		}

		// 更新
		err = l.svcCtx.ConversationsModel.Update(l.ctx, conversations)
		if err != nil {
			return perr.Wrapf(xerr.NewDBErr(), "ConversationsModel.Update error at SetUpUserConversation: %v", err)
		}
		return nil
	}
*/
func (l *SetUpUserConversationLogic) setUpUserConversation(conversationId, userId, recvId string,
	chatType constants.ChatType, isShow bool) error {
	// 用户的会话列表
	conversations, err := l.svcCtx.ConversationsModel.FindByUserId(l.ctx, userId)
	if err != nil {
		if err == im_models.ErrNotFound {
			conversations = &im_models.Conversations{
				ID:               primitive.NewObjectID(),
				UserId:           userId,
				ConversationList: make(map[string]*im_models.Conversation),
			}
		} else {
			return perr.Wrapf(xerr.NewDBErr(), "ConversationsModel.FindOne err %v, req %v", err, userId)
		}
	}

	// 更新会话记录
	if _, ok := conversations.ConversationList[conversationId]; ok {
		return nil
	}

	// 添加会话记录
	conversations.ConversationList[conversationId] = &im_models.Conversation{
		ConversationId: conversationId,
		ChatType:       constants.SingleChatType,
		IsShow:         isShow,
	}

	// 更新
	err = l.svcCtx.ConversationsModel.Update(l.ctx, conversations)
	if err != nil {
		return perr.Wrapf(xerr.NewDBErr(), "ConversationsModel.Insert err %v, req %v", err, conversations)
	}
	return nil
}

/*// 建立会话: 群聊, 私聊
func (l *SetUpUserConversationLogic) SetUpUserConversation(in *im.SetUpUserConversationReq) (*im.SetUpUserConversationResp, error) {
	// todo: add your logic here and delete this line

	var res im.SetUpUserConversationResp
	switch constants.ChatType(in.ChatType) {
	case constants.SingleChatType:
		// 生成会话的id
		conversationId := wuid.CombineId(in.SendId, in.RecvId)
		// 验证是否建立过会话
		conversationRes, err := l.svcCtx.ConversationModel.FindOne(l.ctx, conversationId)
		if err != nil {
			// 建立会话
			if err == im_models.ErrNotFound {
				err = l.svcCtx.ConversationModel.Insert(l.ctx, &im_models.Conversation{
					ConversationId: conversationId,
					ChatType:       constants.SingleChatType,
				})

				if err != nil {
					return nil, errors.Wrapf(xerr.NewDBErr(), "ConversationsModel.Insert err %v", err)
				}
			} else {
				return nil, errors.Wrapf(xerr.NewDBErr(), "ConversationsModel.FindOne err %v, req %v", err, conversationId)
			}
		} else if conversationRes != nil {
			return &res, nil
		}
		// 建立两者的会话
		err = l.setUpUserConversation(conversationId, in.SendId, in.RecvId, constants.SingleChatType, true)
		if err != nil {
			return nil, err
		}
		err = l.setUpUserConversation(conversationId, in.RecvId, in.SendId, constants.SingleChatType, false)
		if err != nil {
			return nil, err
		}
	}

	return &res, nil
}

func (l *SetUpUserConversationLogic) setUpUserConversation(conversationId, userId, recvId string,
	chatType constants.ChatType, isShow bool) error {
	// 用户的会话列表
	conversations, err := l.svcCtx.ConversationsModel.FindByUserId(l.ctx, userId)
	if err != nil {
		if err == im_models.ErrNotFound {
			conversations = &im_models.Conversations{
				ID:               primitive.NewObjectID(),
				UserId:           userId,
				ConversationList: make(map[string]*im_models.Conversation),
			}
		} else {
			return errors.Wrapf(xerr.NewDBErr(), "ConversationsModel.FindOne err %v, req %v", err, userId)
		}
	}

	// 更新会话记录
	if _, ok := conversations.ConversationList[conversationId]; ok {
		return nil
	}

	// 添加会话记录
	conversations.ConversationList[conversationId] = &im_models.Conversation{
		ConversationId: conversationId,
		ChatType:       constants.SingleChatType,
		IsShow:         isShow,
	}

	// 更新
	err = l.svcCtx.ConversationsModel.Update(l.ctx, conversations)
	if err != nil {
		return errors.Wrapf(xerr.NewDBErr(), "ConversationsModel.Insert err %v, req %v", err, conversations)
	}
	return nil
}*/
