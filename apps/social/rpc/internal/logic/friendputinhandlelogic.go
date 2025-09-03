package logic

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"penguin/apps/social/social_models"
	"penguin/pkg/constants"
	"penguin/pkg/xerr"

	"penguin/apps/social/rpc/internal/svc"
	"penguin/apps/social/rpc/social"

	"github.com/zeromicro/go-zero/core/logx"
)

type FriendPutInHandleLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFriendPutInHandleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FriendPutInHandleLogic {
	return &FriendPutInHandleLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *FriendPutInHandleLogic) FriendPutInHandle(in *social.FriendPutInHandleReq) (*social.FriendPutInHandleResp, error) {
	// todo: add your logic here and delete this line

	friendReq, err := l.svcCtx.FriendRequestsModel.FindOne(l.ctx, uint64(in.FriendReqId))

	if err != nil {
		return nil, errors.Wrapf(xerr.NewDBErr(), "[FriendPutInHandle] error: %v, request: %v", err, in.FriendReqId)
	}

	// 好友申请是否有处理过？
	switch constants.HandlerResult(friendReq.HandleResult.Int64) {
	case constants.PassHandleResult:
		return nil, errors.WithStack(xerr.ErrFriendRequestAlreadyPassed)
	case constants.RefuseHandleResult:
		return nil, errors.WithStack(xerr.ErrFriendRequestAlreadyRefused)
	}

	friendReq.HandleResult.Int64 = int64(in.FriendReqId)

	// 修改申请结果
	err = l.svcCtx.FriendRequestsModel.Trans(l.ctx, func(ctx context.Context, session sqlx.Session) error {
		fmt.Println("0xf124515451a")
		if err := l.svcCtx.FriendRequestsModel.Update(l.ctx, session, friendReq); err != nil {
			return errors.Wrapf(xerr.NewDBErr(), "update friend request error: %v, request: %v", err, friendReq)
		}

		if constants.HandlerResult(in.HandleResult) != constants.PassHandleResult {
			fmt.Println("0xf124515451a-0-1")
			return nil
		}

		friends := []*social_models.Friends{
			{
				UserId:    friendReq.UserId,
				FriendUid: friendReq.ReqUid,
			},
			{
				UserId:    friendReq.ReqUid,
				FriendUid: friendReq.UserId,
			},
		}
		fmt.Println("0xf124515451a-1")

		_, err = l.svcCtx.FriendsModel.InsertBatch(l.ctx, session, friends...)
		fmt.Println("0xf124515451a-2")

		if err != nil {
			fmt.Println(err.Error())
			fmt.Printf("insert friends error: %v, request: %v", err, in.FriendReqId)

			return errors.Wrapf(xerr.NewDBErr(), "insert friends error: %v, request: %v", err, in.FriendReqId)
		}

		fmt.Println("0xf124515451a-3")

		return nil
	})

	return &social.FriendPutInHandleResp{}, nil
}
