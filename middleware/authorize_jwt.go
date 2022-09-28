package middleware

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"os"
	"seadeals-backend/apperror"
	"seadeals-backend/dto"
	"seadeals-backend/helper"
	"strings"
)

func AuthorizeJWTFor(role string) gin.HandlerFunc {
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
			unauthorizedError := apperror.UnauthorizedError("Unauthorized")
			if len(splitAuthHeader) < 2 {
				c.AbortWithStatusJSON(unauthorizedError.StatusCode, unauthorizedError)
				return
			}

			encodedToken := splitAuthHeader[1]
			token, err := helper.ValidateToken(encodedToken)
			if err != nil || !token.Valid {
				c.AbortWithStatusJSON(unauthorizedError.StatusCode, unauthorizedError)
				return
			}

			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				c.AbortWithStatusJSON(unauthorizedError.StatusCode, unauthorizedError)
				return
			}

			typeJson, err := json.Marshal(claims["type"])
			var typeString string
			_ = json.Unmarshal(typeJson, &typeString)
			if typeString != dto.JWTAccessToken {
				c.AbortWithStatusJSON(unauthorizedError.StatusCode, unauthorizedError)
				return
			}

			scopeJson, err := json.Marshal(claims["scope"])
			var scope string

			err = json.Unmarshal(scopeJson, &scope)
			if err != nil {
				c.AbortWithStatusJSON(unauthorizedError.StatusCode, unauthorizedError)
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
				c.AbortWithStatusJSON(unauthorizedError.StatusCode, unauthorizedError)
				return
			}

			userJson, err := json.Marshal(claims["user"])
			var user dto.UserJWT

			err = json.Unmarshal(userJson, &user)
			if err != nil {
				c.AbortWithStatusJSON(unauthorizedError.StatusCode, unauthorizedError)
				return
			}

			c.Set("user", user)
			c.Set("token", encodedToken)
		}
	}
}
