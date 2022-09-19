package middleware

//
//import (
//	"encoding/json"
//	"git.garena.com/sea-labs-id/batch-01/reivaldo-julianto/final-project-togetherbnb-backend/dto"
//	"git.garena.com/sea-labs-id/batch-01/reivaldo-julianto/final-project-togetherbnb-backend/httperror"
//	"git.garena.com/sea-labs-id/batch-01/reivaldo-julianto/final-project-togetherbnb-backend/utils"
//	"github.com/gin-gonic/gin"
//	"os"
//	"strings"
//)
//
//func AuthorizeJWTFor(role string) gin.HandlerFunc {
//	return func(c *gin.Context) {
//		if os.Getenv("ENV") == "testing" {
//			user := dto.UserJWT{
//				ID:       1,
//				Email:    "test",
//				CityName: "test",
//				WalletID: 1,
//			}
//			c.Set("user", user)
//			c.Set("token", "test")
//			return
//		} else if os.Getenv("ENV") != "testingNoUser" {
//			authHeader := c.GetHeader("Authorization")
//
//			splitAuthHeader := strings.Split(authHeader, "Bearer ")
//			unauthorizedError := httperror.UnauthorizeError()
//			if len(splitAuthHeader) < 2 {
//				c.AbortWithStatusJSON(unauthorizedError.StatusCode, unauthorizedError)
//				return
//			}
//
//			encodedToken := splitAuthHeader[1]
//			token, err := utils.ValidateToken(encodedToken)
//			if err != nil || !token.Valid {
//				c.AbortWithStatusJSON(unauthorizedError.StatusCode, unauthorizedError)
//				return
//			}
//
//			claims, ok := token.Claims.(jwt.MapClaims)
//			if !ok {
//				c.AbortWithStatusJSON(unauthorizedError.StatusCode, unauthorizedError)
//				return
//			}
//
//			scopeJson, err := json.Marshal(claims["scope"])
//			var scope string
//
//			err = json.Unmarshal(scopeJson, &scope)
//			if err != nil {
//				c.AbortWithStatusJSON(unauthorizedError.StatusCode, unauthorizedError)
//				return
//			}
//			splitScope := strings.Split(scope, " ")
//			isAuthorize := false
//			for _, s := range splitScope {
//				if s == role {
//					isAuthorize = true
//					break
//				}
//			}
//			if !isAuthorize {
//				c.AbortWithStatusJSON(unauthorizedError.StatusCode, unauthorizedError)
//				return
//			}
//
//			userJson, err := json.Marshal(claims["user"])
//			var user dto.UserJWT
//
//			err = json.Unmarshal(userJson, &user)
//			if err != nil {
//				c.AbortWithStatusJSON(unauthorizedError.StatusCode, unauthorizedError)
//				return
//			}
//
//			c.Set("user", user)
//			c.Set("token", encodedToken)
//		}
//	}
//}
