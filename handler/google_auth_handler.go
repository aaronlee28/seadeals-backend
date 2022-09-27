package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"net/http"
	"os"
	"seadeals-backend/apperror"
	"seadeals-backend/dto"
	"time"
)

func (h *Handler) GoogleSignIn(ctx *gin.Context) {
	value, _ := ctx.Get("payload")
	googleLogin, _ := value.(*dto.GoogleLogin)
	claims := jwt.MapClaims{}
	_, _ = jwt.ParseWithClaims(googleLogin.TokenID, &claims, nil)

	timeNow := time.Now().Unix()
	timeExp := int64(claims["exp"].(float64))
	aud := claims["aud"].(string)
	clientID := os.Getenv("GOOGLE_CLIENT_ID")
	if claims["iss"].(string) != "https://accounts.google.com" || claims["email_verified"].(bool) != true || timeExp < timeNow || aud != clientID {
		_ = ctx.Error(apperror.UnauthorizedError("Unauthorized token"))
		return
	}

	user, err := h.userService.CheckGoogleAccount(claims["email"].(string))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.AppResponse{
			StatusCode: http.StatusBadRequest,
			Status:     "BAD_REQUEST_ERROR",
			Data:       gin.H{"error": err.Error(), "user": gin.H{"email": claims["email"], "name": claims["name"]}},
		})
		return
	}

	accessToken, refreshToken, err := h.authService.SignInWithGoogle(user)
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	ctx.SetSameSite(http.SameSiteNoneMode)
	if os.Getenv("ENV") == "dev" {
		ctx.SetCookie("refresh_token", refreshToken, 60*60*24, "/", ctx.Request.Header.Get("Origin"), false, false)
	} else {
		ctx.SetCookie("refresh_token", refreshToken, 60*60*24, "/", ctx.Request.Header.Get("Origin"), true, true)
	}

	ctx.JSON(http.StatusOK, dto.StatusOKResponse(gin.H{"user": user, "has_login": true, "id_token": accessToken}))
}
