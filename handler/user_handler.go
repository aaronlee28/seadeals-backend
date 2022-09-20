package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"seadeals-backend/dto"
	"seadeals-backend/model"
)

func (h *Handler) Register(ctx *gin.Context) {
	value, _ := ctx.Get("payload")
	json, _ := value.(*dto.RegisterRequest)

	result, tx, err := h.userService.Register(json)
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	userJWT := &model.User{
		ID:       result.ID,
		Email:    result.Email,
		Username: result.Username,
	}
	accessToken, refreshToken, err := h.authService.AuthAfterRegister(userJWT, &result.Wallet, tx)
	if err != nil {
		_ = ctx.Error(err)
		return
	}
	ctx.SetSameSite(http.SameSiteNoneMode)
	if os.Getenv("ENV") == "dev" {
		ctx.SetCookie("access_token", refreshToken, 60*60*24, "/", ctx.Request.Header.Get("Origin"), false, false)
	} else {
		ctx.SetCookie("access_token", refreshToken, 60*60*24, "/", ctx.Request.Header.Get("Origin"), true, true)
	}

	ctx.JSON(http.StatusCreated, dto.StatusCreatedResponse(gin.H{"data": gin.H{"user": result, "id_token": accessToken}}))
}
func (h *Handler) SignInWithGoogleEmail(ctx *gin.Context) {
	value, _ := ctx.Get("payload")
	json, _ := value.(*dto.GoogleLogin)

	result, err := h.userService.CheckGoogleAccount(json)
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	if result == nil {
		ctx.JSON(http.StatusCreated, dto.StatusCreatedResponse(gin.H{"data": gin.H{"user": result, "has_login": false, "id_token": ""}}))
		return
	}

	accessToken, refreshToken, err := h.authService.SignInWithGoogle(result)
	if err != nil {
		_ = ctx.Error(err)
		return
	}
	ctx.SetSameSite(http.SameSiteNoneMode)
	if os.Getenv("ENV") == "dev" {
		ctx.SetCookie("access_token", refreshToken, 60*60*24, "/", ctx.Request.Header.Get("Origin"), false, false)
	} else {
		ctx.SetCookie("access_token", refreshToken, 60*60*24, "/", ctx.Request.Header.Get("Origin"), true, true)
	}

	ctx.JSON(http.StatusCreated, dto.StatusCreatedResponse(gin.H{"data": gin.H{"user": result, "has_login": true, "id_token": accessToken}}))
}
