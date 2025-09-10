package friend

import (
	"context"
	"fmt"
	"github.com/jinzhu/copier"
	"github.com/pkg/errors"
	"penguin/pkg/xerr"

	"penguin/apps/social/rpc/internal/svc"
	"penguin/apps/social/rpc/social"

	"github.com/zeromicro/go-zero/core/logx"
)

type FriendListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFriendListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FriendListLogic {
	return &FriendListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *FriendListLogic) FriendList(in *social.FriendListReq) (*social.FriendListResp, error) {
	// todo: add your logic here and delete this line

	friendslist, err := l.svcCtx.FriendsModel.ListByUserId(l.ctx, in.UserId)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewDBErr(), "list friend by uid err %v req %v", err, in.UserId)
	}

	var respList []*social.Friends
	err = copier.Copy(&respList, &friendslist)
	if err != nil {
		fmt.Println("[friendlistlogic.go] copier.Copy friendslist err:", err.Error())
		return nil, err
	}

	return &social.FriendListResp{
		List: respList,
	}, nil
}
