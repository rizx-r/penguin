package ctxdata

import (
	"github.com/golang-jwt/jwt/v4"
)

const Identify = "fum-im.com"

func GetJwtFromToken(secretKey string, iat, duration int64, uid string) (string, error) {
	claims := make(jwt.MapClaims)
	claims["exp"] = iat + duration
	claims["iat"] = iat
	claims[Identify] = uid

	token := jwt.New(jwt.SigningMethodHS256)
	token.Claims = claims

	return token.SignedString([]byte(secretKey))
}
