package logic

import (
	"context"
	"database/sql"
	"fmt"
	"penguin/pkg/xerr"

	//"errors"
	"penguin/apps/user/models"
	"penguin/pkg/ctxdata"
	"penguin/pkg/encrypt"
	"penguin/pkg/wuid"
	"time"

	"penguin/apps/user/rpc/internal/svc"
	"penguin/apps/user/rpc/user"

	"github.com/zeromicro/go-zero/core/logx"
)

//var (
//	ErrPhoneRegistered = errors.New("phone is already registered")
//)

type RegisterLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *RegisterLogic) Register(in *user.RegisterReq) (*user.RegisterResp, error) {
	// todo: add your logic here and delete this line

	userEntity, err := l.svcCtx.UsersModel.FindByPhone(l.ctx, in.Phone)
	if err != nil && err != models.ErrNotFound {
		return nil, err
	}

	if userEntity != nil {
		fmt.Println(userEntity)
		return nil, xerr.ErrPhoneRegistered
	}

	userEntity = &models.Users{
		Id:       wuid.GenUid(l.svcCtx.Config.Mysql.Datasource),
		Avatar:   in.Avatar,
		Nickname: in.Nickname,
		Phone:    in.Phone,
		Sex: sql.NullInt64{
			Int64: int64(in.Sex),
			Valid: true,
		},
	}

	if len(in.Password) > 0 {
		encodedPassword, err := encrypt.GenPasswordHash([]byte(in.Password))
		if err != nil {
			return nil, err
		}
		userEntity.Password = sql.NullString{
			String: string(encodedPassword),
			Valid:  true,
		}
	}

	_, err = l.svcCtx.UsersModel.Insert(l.ctx, userEntity)
	if err != nil {
		return nil, err
	}

	now := time.Now().Unix()
	token, err := ctxdata.GetJwtFromToken(l.svcCtx.Config.Jwt.AccessSecret, now, l.svcCtx.Config.Jwt.AccessExpire, userEntity.Id)
	if err != nil {
		return nil, err
	}

	return &user.RegisterResp{
		Token:  token,
		Expire: now + l.svcCtx.Config.Jwt.AccessExpire,
	}, nil
}
