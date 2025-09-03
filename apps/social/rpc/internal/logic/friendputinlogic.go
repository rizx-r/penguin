package logic

import (
	"context"
	"database/sql"
	"github.com/pkg/errors"
	"penguin/apps/social/social_models"
	"penguin/pkg/constants"
	"penguin/pkg/xerr"
	"time"

	"penguin/apps/social/rpc/internal/svc"
	"penguin/apps/social/rpc/social"

	"github.com/zeromicro/go-zero/core/logx"
)

type FriendPutInLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFriendPutInLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FriendPutInLogic {
	return &FriendPutInLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *FriendPutInLogic) FriendPutIn(in *social.FriendPutInReq) (*social.FriendPutInResp, error) {
	// todo: add your logic here and delete this line

	// 申请者与被申请者是否已经是好友关系
	friend, err := l.svcCtx.FriendsModel.FindByUidAndFid(l.ctx, in.UserId, in.ReqUid)
	if err != nil && err != social_models.ErrNotFound {
		return nil, errors.Wrapf(xerr.NewDBErr(), "find friends by uid and fid error: %v req: %v", err, in)
	}
	if friend != nil {
		return &social.FriendPutInResp{}, err
	}

	// 是否已申请过
	friendRequested, err := l.svcCtx.FriendRequestsModel.FindByReqUidAndUserId(l.ctx, in.ReqUid, in.UserId)
	if err != nil && err != social_models.ErrNotFound {
		return nil, errors.Wrapf(xerr.NewDBErr(), "[FindByReqUidAndUserId] error: %v req: %v", err, in)
	}
	if friendRequested != nil {
		return &social.FriendPutInResp{}, err
	}

	// 创建好友申请
	_, err = l.svcCtx.FriendRequestsModel.Insert(l.ctx, &social_models.FriendRequests{
		UserId: in.UserId,
		ReqUid: in.ReqUid,
		ReqMsg: sql.NullString{
			Valid:  true,
			String: in.ReqMsg,
		},
		ReqTime: time.Unix(in.ReqTime, 0),
		HandleResult: sql.NullInt64{
			Int64: int64(constants.NoHandleResult),
			Valid: true,
		},
	})

	if err != nil {
		return nil, errors.Wrapf(xerr.NewDBErr(), "insert friend error: %v req: %v", err, in)

	}

	return &social.FriendPutInResp{}, nil
}
