package friend

import (
	"context"
	"penguin/apps/social/rpc/socialclient"
	"penguin/apps/user/rpc/userclient"
	"penguin/pkg/ctxdata"

	"penguin/apps/social/api/internal/svc"
	"penguin/apps/social/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type FriendListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 好友列表
func NewFriendListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FriendListLogic {
	return &FriendListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FriendListLogic) FriendList(req *types.FriendListReq) (resp *types.FriendListResp, err error) {
	// todo: add your logic here and delete this line

	uid := ctxdata.GetUid(l.ctx)

	friends, err := l.svcCtx.Social.FriendList(l.ctx, &socialclient.FriendListReq{
		UserId: uid,
	})

	//if err != nil {
	//	return nil, err
	//}

	if err != nil {
		return &types.FriendListResp{}, nil
	}

	// 根据好友id获取好友列表
	uids := make([]string, 0, len(friends.List))
	for _, friend := range friends.List {
		uids = append(uids, friend.FriendUid)
	}

	// 根据uid查询用户信息
	users, err := l.svcCtx.User.FindUser(l.ctx, &userclient.FindUserReq{
		Ids: uids,
	})
	if err != nil {
		return nil, err
	}

	userRecords := make(map[string]*userclient.UserEntity, len(users.User))
	for _, user := range users.User {
		userRecords[user.Id] = user
	}

	respList := make([]*types.Friends, 0, len(friends.List))
	for _, friend := range friends.List {
		friend_ := &types.Friends{
			Id:        friend.Id,
			FriendUid: friend.FriendUid,
		}

		if u, ok := userRecords[friend.FriendUid]; ok {
			friend_.Nickname = u.Nickname
			friend_.Avatar = u.Avatar
		}

		respList = append(respList, friend_)
	}

	return &types.FriendListResp{
		List: respList,
	}, nil
}
