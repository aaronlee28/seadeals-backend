package helper

import (
	"github.com/golang-jwt/jwt/v4"
	"seadeals-backend/config"
	"seadeals-backend/dto"
	"time"
)

type idTokenClaims struct {
	jwt.RegisteredClaims
	User  *dto.UserJWT `json:"user"`
	Scope string       `json:"scope"`
	Type  string       `json:"type"`
}

func GenerateJWTToken(user *dto.UserJWT, role string, idExp int64, jwtType string) (string, error) {
	unixTime := time.Now().Unix()
	tokenExp := unixTime + idExp

	timeExpire := jwt.NumericDate{Time: time.Unix(tokenExp, 0)}
	timeNow := jwt.NumericDate{Time: time.Now()}
	accessClaims := &idTokenClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: &timeExpire,
			IssuedAt:  &timeNow,
			Issuer:    config.Config.AppName,
		},
		User:  user,
		Scope: role,
		Type:  jwtType,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	tokenString, _ := token.SignedString(config.Config.JWTSecret)

	return tokenString, nil
}
