package user

import (
	"context"
	"fmt"
	"github.com/jinzhu/copier"
	"penguin/apps/user/rpc/user"

	"penguin/apps/user/api/internal/svc"
	"penguin/apps/user/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type LoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 用户登入
func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LoginLogic) Login(req *types.LoginReq) (resp *types.LoginResp, err error) {
	// todo: add your logic here and delete this line
	loginResp, err := l.svcCtx.User.Login(l.ctx, &user.LoginReq{
		Phone:    req.Phone,
		Password: req.Password,
	})

	if err != nil {
		fmt.Println("login error: ", err)
		return nil, err
	}

	var res types.LoginResp
	err = copier.Copy(&res, loginResp)
	if err != nil {
		fmt.Println("copier.Copy failed:", err)
		return nil, err
	}

	return &res, nil
}
