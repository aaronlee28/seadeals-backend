package middleware

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"os"
	"seadeals-backend/dto"
	"seadeals-backend/helper"
	"strings"
)

func OptionalAuthorizeJWTFor(role string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if os.Getenv("ENV") == "testing" {
			user := dto.UserJWT{
				UserID:   1,
				Email:    "test",
				Username: "test",
				WalletID: 1,
			}
			c.Set("user", user)
			c.Set("token", "test")
			return
		} else if os.Getenv("ENV") != "testingNoUser" {
			authHeader := c.GetHeader("Authorization")

			splitAuthHeader := strings.Split(authHeader, "Bearer ")
			if len(splitAuthHeader) < 2 {
				c.Next()
				return
			}

			encodedToken := splitAuthHeader[1]
			token, err := helper.ValidateToken(encodedToken)
			if err != nil || !token.Valid {
				c.Next()
				return
			}

			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				c.Next()
				return
			}

			typeJson, err := json.Marshal(claims["type"])
			var typeString string
			_ = json.Unmarshal(typeJson, &typeString)
			if typeString != dto.JWTAccessToken {
				fmt.Println("1")
				c.Next()
				return
			}

			scopeJson, err := json.Marshal(claims["scope"])
			var scope string

			err = json.Unmarshal(scopeJson, &scope)
			if err != nil {
				c.Next()
				return
			}
			splitScope := strings.Split(scope, " ")
			isAuthorize := false
			for _, s := range splitScope {
				if s == role {
					isAuthorize = true
					break
				}
			}
			if !isAuthorize {
				c.Next()
				return
			}

			userJson, err := json.Marshal(claims["user"])
			var user dto.UserJWT

			err = json.Unmarshal(userJson, &user)
			if err != nil {
				c.Next()
				return
			}

			c.Set("user", user)
			c.Set("token", encodedToken)
		}
	}
}
