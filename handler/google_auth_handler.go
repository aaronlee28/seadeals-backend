package handler

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"seadeals-backend/apperror"
	"seadeals-backend/config"
	"seadeals-backend/dto"
)

func (h *Handler) GoogleSignIn(ctx *gin.Context) {
	googleConfig := config.SetupGoogleAuthConfig()
	url := googleConfig.AuthCodeURL(config.RandomState)
	ctx.Redirect(http.StatusTemporaryRedirect, url)
}

func (h *Handler) GoogleCallback(ctx *gin.Context) {
	queryParams := ctx.Request.URL.Query()

	state := queryParams["state"][0]
	if state != config.RandomState {
		ctx.JSON(http.StatusBadRequest, apperror.BadRequestError("state doesn't match	"))
		return
	}

	code := queryParams["code"][0]
	googleConfig := config.SetupGoogleAuthConfig()
	token, err := googleConfig.Exchange(ctx, code)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, apperror.UnauthorizedError("wrong token authentication"))
		return
	}

	res, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, apperror.InternalServerError("failed to fetch user data"))
		return
	}

	var userRes map[string]any
	userData, err := ioutil.ReadAll(res.Body)
	err = json.Unmarshal(userData, &userRes)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, apperror.InternalServerError("failed to parse JSON"))
		return
	}

	ctx.JSON(http.StatusOK, dto.StatusOKResponse(userRes))
}
