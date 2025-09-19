package logic

import (
	"context"
	"errors"
	"penguin/apps/user/user_models"
	"penguin/pkg/xerr"

	"penguin/apps/user/rpc/internal/svc"
	"penguin/apps/user/rpc/user"

	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserInfoLogic {
	return &GetUserInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetUserInfoLogic) GetUserInfo(in *user.GetUserInfoReq) (*user.GetUserInfoResp, error) {
	// todo: add your logic here and delete this line

	userEntity, err := l.svcCtx.UsersModel.FindOne(l.ctx, in.Id)
	if err != nil {
		if err == user_models.ErrNotFound {
			return nil, xerr.ErrUserPwdNotMatched
		}
		return nil, err
	}

	var resp user.UserEntity

	err = copier.Copy(&resp, userEntity)
	if err != nil {
		return nil, errors.New("copy userEntity fail")
	}

	return &user.GetUserInfoResp{User: &resp}, nil
}
