package handler

import (
	"encoding/base64"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"seadeals-backend/apperror"
	"seadeals-backend/config"
	"seadeals-backend/dto"
	"time"
)

func generateStateOauthCookie(ctx *gin.Context) string {
	var expiration = time.Now().Add(365 * 24 * time.Hour)

	b := make([]byte, 16)
	rand.Read(b)
	state := base64.URLEncoding.EncodeToString(b)
	cookie := http.Cookie{Name: "oauth_state", Value: state, Expires: expiration}
	http.SetCookie(ctx.Writer, &cookie)

	return state
}

func (h *Handler) GoogleSignIn(ctx *gin.Context) {
	googleConfig := config.SetupGoogleAuthConfig()
	randState := generateStateOauthCookie(ctx)
	url := googleConfig.AuthCodeURL(randState)
	ctx.Redirect(http.StatusTemporaryRedirect, url)
}

func (h *Handler) GoogleCallback(ctx *gin.Context) {
	queryParams := ctx.Request.URL.Query()

	state := queryParams["state"][0]
	oauthState, _ := ctx.Request.Cookie("oauth_state")
	if state != oauthState.Value {
		_ = ctx.Error(apperror.BadRequestError("state doesn't match"))
		return
	}

	code := queryParams["code"][0]
	googleConfig := config.SetupGoogleAuthConfig()
	token, err := googleConfig.Exchange(ctx, code)
	if err != nil {
		_ = ctx.Error(apperror.UnauthorizedError("wrong token authentication"))
		return
	}

	res, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	if err != nil {
		_ = ctx.Error(apperror.InternalServerError("failed to fetch user data"))
		return
	}

	var googleCallbackRes *dto.GoogleLogin
	userData, err := ioutil.ReadAll(res.Body)
	err = json.Unmarshal(userData, &googleCallbackRes)
	if err != nil {
		_ = ctx.Error(apperror.InternalServerError("failed to parse JSON"))
		return
	}

	user, err := h.userService.CheckGoogleAccount(googleCallbackRes)
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	accessToken, refreshToken, err := h.authService.SignInWithGoogle(user)
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

	ctx.JSON(http.StatusOK, dto.StatusOKResponse(gin.H{"data": gin.H{"user": user, "has_login": true, "id_token": accessToken}}))
}
