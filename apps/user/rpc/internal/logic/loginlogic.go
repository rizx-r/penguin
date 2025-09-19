package logic

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"penguin/apps/user/user_models"
	"penguin/pkg/ctxdata"
	"penguin/pkg/encrypt"
	"penguin/pkg/xerr"
	"time"

	"penguin/apps/user/rpc/internal/svc"
	"penguin/apps/user/rpc/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type LoginLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *LoginLogic) Login(in *user.LoginReq) (*user.LoginResp, error) {
	// todo: add your logic here and delete this line

	userEntity, err := l.svcCtx.UsersModel.FindByPhone(l.ctx, in.Phone)
	userEntity1, err := l.svcCtx.UsersModel.FindByPhone(l.ctx, "13700001111")
	userEntity2, err := l.svcCtx.UsersModel.FindByPhone(l.ctx, "13700001112")
	fmt.Println("userEntity: ", userEntity)
	fmt.Println("userEntity1: ", userEntity1)
	fmt.Println("userEntity2: ", userEntity2)

	if err != nil {
		if err == user_models.ErrNotFound {
			return nil, errors.WithStack(xerr.ErrPhoneNotRegistered)
		}
		return nil, errors.Wrapf(xerr.NewDBErr(), "error: find user by phone. err: %v, request: %v", err, in.Phone)
	}

	if !encrypt.ValidatePasswordHash(in.Password, userEntity.Password.String) {
		//zap.S().Infof("[ErrUserPwdNotMatched], %v", in)
		fmt.Printf("[ErrUserPwdNotMatched], %v, (%s) \n", in, userEntity.Password.String)
		return nil, errors.WithStack(xerr.ErrUserPwdNotMatched)
	}

	now := time.Now().Unix()
	token, err := ctxdata.GetJwtFromToken(l.svcCtx.Config.Jwt.AccessSecret, now, l.svcCtx.Config.Jwt.AccessExpire, userEntity.Id)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewInternalErr(), "error: get token. err: %v, request: %v", err, in.Phone)
	}

	return &user.LoginResp{
		Token:  token,
		Expire: now + l.svcCtx.Config.Jwt.AccessExpire,
	}, nil
}
