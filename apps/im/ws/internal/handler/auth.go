package handler

import (
	"context"
	"github.com/golang-jwt/jwt/v4"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/token"
	"net/http"
	"penguin/apps/im/ws/internal/svc"
	"penguin/pkg/ctxdata"
)

type (
	JwtAuth struct {
		svc    *svc.ServiceContext
		parser *token.TokenParser
		logx.Logger
	}
)

func NewJwtAuth(svc *svc.ServiceContext) *JwtAuth {
	return &JwtAuth{
		svc:    svc,
		parser: token.NewTokenParser(),
		Logger: logx.WithContext(context.Background()),
	}
}

func (j JwtAuth) Authenticate(w http.ResponseWriter, r *http.Request) (ok bool) {
	tok, err := j.parser.ParseToken(r, j.svc.Config.JwtAuth.AccessSecret, "")
	if err != nil {
		j.Errorf("pare token err at im/ws/internal/handler/auth.go:Authenticate: %v", err)
		return false
	}

	if !tok.Valid {
		return false
	}

	claims, ok := tok.Claims.(jwt.MapClaims)
	if !ok {
		return false
	}

	*r = *r.WithContext(context.WithValue(r.Context(), ctxdata.Identify, claims[ctxdata.Identify]))
	return true
}

func (j JwtAuth) UserId(r *http.Request) (uid string) {
	//TODO implement me
	return ctxdata.GetUid(r.Context())
}
