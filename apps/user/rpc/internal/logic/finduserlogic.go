package logic

import (
	"context"
	"github.com/jinzhu/copier"
	"penguin/apps/user/user_models"

	"penguin/apps/user/rpc/internal/svc"
	"penguin/apps/user/rpc/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type FindUserLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFindUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FindUserLogic {
	return &FindUserLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *FindUserLogic) FindUser(in *user.FindUserReq) (*user.FindUserResp, error) {
	// todo: add your logic here and delete this line

	var (
		userEntitys []*user_models.Users
		err         error
		resp        []*user.UserEntity
	)

	if in.Phone != "" {
		userEntity, err := l.svcCtx.FindByPhone(l.ctx, in.Phone)
		if err != nil {
			userEntitys = append(userEntitys, userEntity)
		}
	} else if in.Name != "" {
		userEntitys, err = l.svcCtx.UsersModel.ListByNickname(l.ctx, in.Name)
	} else if len(in.Ids) > 0 {
		userEntitys, err = l.svcCtx.UsersModel.ListByIds(l.ctx, in.Ids)
	}

	if err != nil {
		return nil, err
	}

	err = copier.Copy(&resp, userEntitys)
	if err != nil {
		return nil, err
	}

	return &user.FindUserResp{
		User: resp,
	}, nil
}
