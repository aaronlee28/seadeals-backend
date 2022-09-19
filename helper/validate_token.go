package helper

import (
	"github.com/golang-jwt/jwt/v4"
	"os"
	"seadeals-backend/apperror"
	"seadeals-backend/config"
)

type fakeClaim map[string]interface{}

func (f fakeClaim) Valid() error {
	return nil
}

var fake = fakeClaim{}

func ValidateToken(encodedToken string) (*jwt.Token, error) {
	if os.Getenv("ENV") == "testing" {
		var claim = jwt.MapClaims{}
		if os.Getenv("TEST") == "invalid user" {
			claim["user"] = map[string]interface{}{
				"email": 123,
			}
		} else {
			claim["user"] = map[string]interface{}{
				"id":       1,
				"email":    "test",
				"city":     "test",
				"walletID": 1,
			}
		}

		if os.Getenv("TEST") == "invalid scope" {
			claim["scope"] = 123
		} else {
			claim["scope"] = "test"
		}

		claim["token"] = "test"
		claim["exp"] = "test"
		claim["iat"] = "test"
		claim["iss"] = "test"

		if os.Getenv("TEST") == "invalid token" {
			return &jwt.Token{
				Raw:       "test",
				Method:    nil,
				Header:    nil,
				Claims:    fake,
				Signature: string(config.Config.JWTSecret),
				Valid:     true,
			}, nil
		}
		return &jwt.Token{
			Raw:       "test",
			Method:    nil,
			Header:    nil,
			Claims:    claim,
			Signature: string(config.Config.JWTSecret),
			Valid:     true,
		}, nil
	}

	return jwt.Parse(encodedToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, apperror.UnauthorizedError("")
		}

		return config.Config.JWTSecret, nil
	})
}
