package handler

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"net/http"
	"os"
	"seadeals-backend/apperror"
	"seadeals-backend/config"
	"seadeals-backend/dto"
	"seadeals-backend/helper"
	"time"
)

type idTokenClaims struct {
	jwt.RegisteredClaims
	User  *dto.UserJWT `json:"user"`
	Scope string       `json:"scope"`
	Type  string       `json:"type"`
}

func (h *Handler) RefreshAccessToken(c *gin.Context) {
	unauthorizedError := apperror.UnauthorizedError("Invalid Refresh Token")

	refreshJWT, err := c.Cookie("refresh_token")
	if err != nil {
		c.AbortWithStatusJSON(unauthorizedError.StatusCode, unauthorizedError)
		return
	}

	isInDatabase, userID, err := h.refreshTokenService.CheckIfTokenExist(refreshJWT)
	if err != nil || !isInDatabase {
		c.AbortWithStatusJSON(unauthorizedError.StatusCode, unauthorizedError)
		return
	}

	token, err := helper.ValidateToken(refreshJWT)
	if err != nil || !token.Valid {
		c.AbortWithStatusJSON(unauthorizedError.StatusCode, unauthorizedError)
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		c.AbortWithStatusJSON(unauthorizedError.StatusCode, unauthorizedError)
		return
	}

	userJson, _ := json.Marshal(claims["user"])
	var user dto.UserJWT

	err = json.Unmarshal(userJson, &user)
	if err != nil {
		c.AbortWithStatusJSON(unauthorizedError.StatusCode, unauthorizedError)
		return
	}

	if userID != user.UserID {
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

	var idExp = config.Config.JWTExpiredInMinuteTime * 60
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
		User:  &user,
		Scope: scope,
		Type:  dto.JWTAccessToken,
	}
	newToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	tokenString, err := newToken.SignedString(config.Config.JWTSecret)

	if os.Getenv("ENV") == "testing" {
		tokenString = "test"
	}

	c.JSON(http.StatusOK, dto.StatusOKResponse(gin.H{"user_id": userID, "id_token": tokenString}))
}
