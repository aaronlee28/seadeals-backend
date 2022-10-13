package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"seadeals-backend/apperror"
	"seadeals-backend/dto"
	"seadeals-backend/model"
)

func (h *Handler) Ping(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, dto.StatusOKResponse("Pong"))
}

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
		FullName: result.FullName,
	}
	accessToken, refreshToken, err := h.authService.AuthAfterRegister(userJWT, &result.Wallet, tx)
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

	ctx.JSON(http.StatusCreated, dto.StatusCreatedResponse(gin.H{"data": gin.H{"user": result, "id_token": accessToken}}))
}

func (h *Handler) SignIn(ctx *gin.Context) {
	value, _ := ctx.Get("payload")
	json, _ := value.(*dto.SignInReq)

	accessToken, refreshToken, err := h.authService.SignIn(json)
	ctx.SetSameSite(http.SameSiteNoneMode)
	if os.Getenv("ENV") == "dev" {
		ctx.SetCookie("refresh_token", refreshToken, 60*60*24, "/", ctx.Request.Header.Get("Origin"), false, false)
	} else {
		ctx.SetCookie("refresh_token", refreshToken, 60*60*24, "/", ctx.Request.Header.Get("Origin"), true, true)
	}

	if err != nil {
		_ = ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, dto.StatusOKResponse(gin.H{"id_token": accessToken}))
}

func (h *Handler) SignOut(ctx *gin.Context) {
	user, exists := ctx.Get("user")
	if !exists {
		_ = ctx.Error(apperror.BadRequestError("User is invalid"))
		return
	}

	value, _ := ctx.Get("payload")
	json, _ := value.(*dto.SignOutReq)
	jwtUser := user.(dto.UserJWT)
	if jwtUser.UserID != json.UserID {
		_ = ctx.Error(apperror.UnauthorizedError("Cannot log out another user"))
		return
	}

	err := h.authService.SignOut(jwtUser.UserID)
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	ctx.SetSameSite(http.SameSiteNoneMode)
	if os.Getenv("ENV") == "dev" {
		ctx.SetCookie("refresh_token", "", 60*60*24, "/", ctx.Request.Header.Get("Origin"), false, false)
	} else {
		ctx.SetCookie("refresh_token", "", 60*60*24, "/", ctx.Request.Header.Get("Origin"), true, true)
	}
	ctx.JSON(http.StatusOK, dto.StatusOKResponse(gin.H{"logout_user": user}))
}

func (h *Handler) StepUpPassword(ctx *gin.Context) {
	value, _ := ctx.Get("payload")
	json, _ := value.(*dto.StepUpPasswordRes)

	userPayload, _ := ctx.Get("user")
	user, _ := userPayload.(dto.UserJWT)
	userID := user.UserID

	err := h.authService.StepUpPassword(userID, json)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}
	successResponse := dto.StatusOKResponse(err)
	ctx.JSON(http.StatusOK, successResponse)
}

func (h *Handler) RegisterAsSeller(ctx *gin.Context) {
	value, _ := ctx.Get("payload")
	json, _ := value.(*dto.RegisterAsSellerReq)

	userPayload, _ := ctx.Get("user")
	user, _ := userPayload.(dto.UserJWT)
	userID := user.UserID

	seller, accessToken, err := h.userService.RegisterAsSeller(json, userID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}
	successResponse := dto.StatusOKResponse(gin.H{"seller": seller, "id_token": accessToken})
	ctx.JSON(http.StatusOK, successResponse)
}

func (h *Handler) UserDetails(ctx *gin.Context) {
	userPayload, _ := ctx.Get("user")
	user, _ := userPayload.(dto.UserJWT)
	userID := user.UserID

	res, err := h.userService.UserDetails(userID)
	fmt.Println("userid", userID)
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, dto.StatusOKResponse(res))
}

func (h *Handler) ChangeUserDetails(ctx *gin.Context) {

	userPayload, _ := ctx.Get("user")
	user, _ := userPayload.(dto.UserJWT)
	userID := user.UserID

	value, _ := ctx.Get("payload")
	json, _ := value.(*dto.ChangeUserDetails)

	res, err := h.userService.ChangeUserDetailsLessPassword(userID, json)
	fmt.Println("userid", userID)

	if err != nil {
		_ = ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, dto.StatusOKResponse(res))
}

func (h *Handler) ChangeUserPassword(ctx *gin.Context) {

	userPayload, _ := ctx.Get("user")
	user, _ := userPayload.(dto.UserJWT)
	userID := user.UserID

	value, _ := ctx.Get("payload")
	json, _ := value.(*dto.ChangePasswordReq)

	err := h.userService.ChangeUserPassword(userID, json)

	if err != nil {
		_ = ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, dto.StatusOKResponse("Ok"))
}
