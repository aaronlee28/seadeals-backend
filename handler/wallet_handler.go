package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"seadeals-backend/apperror"
	"seadeals-backend/dto"
	"seadeals-backend/helper"
	"seadeals-backend/repository"
	"strconv"
)

func (h *Handler) WalletDataTransactions(ctx *gin.Context) {
	payload, _ := ctx.Get("user")
	user, _ := payload.(dto.UserJWT)
	userID := user.UserID

	result, err := h.walletService.UserWalletData(userID)
	if err != nil {
		_ = ctx.Error(err)
		return
	}
	successResponse := dto.StatusOKResponse(result)
	ctx.JSON(http.StatusOK, successResponse)
}

func (h *Handler) TransactionDetails(ctx *gin.Context) {
	payload, _ := ctx.Get("user")
	user, _ := payload.(dto.UserJWT)
	userID := user.UserID

	idParam, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		_ = ctx.Error(apperror.BadRequestError("Invalid id format"))
		return
	}

	result, err := h.walletService.TransactionDetails(userID, uint(idParam))
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	successResponse := dto.StatusOKResponse(result)
	ctx.JSON(http.StatusOK, successResponse)
}

func (h *Handler) PaginatedTransactions(ctx *gin.Context) {
	payload, _ := ctx.Get("user")
	user, _ := payload.(dto.UserJWT)
	userID := user.UserID

	query := &repository.Query{
		Limit: ctx.Query("limit"),
		Page:  ctx.Query("page"),
	}

	result, err := h.walletService.PaginatedTransactions(query, userID)
	if err != nil {
		_ = ctx.Error(err)
		return
	}
	successResponse := dto.StatusOKResponse(result)
	ctx.JSON(http.StatusOK, successResponse)
}

func (h *Handler) GetWalletTransactions(ctx *gin.Context) {
	payload, _ := ctx.Get("user")
	user, _ := payload.(dto.UserJWT)
	userID := user.UserID

	query := &dto.WalletTransactionsQuery{
		SortBy: helper.GetQuery(ctx, "sortBy", ""),
		Sort:   helper.GetQuery(ctx, "sort", ""),
		Limit:  helper.GetQuery(ctx, "limit", "10"),
		Page:   helper.GetQuery(ctx, "page", "1"),
	}
	limit, _ := strconv.ParseUint(query.Limit, 10, 64)
	if limit == 0 {
		limit = 20
	}
	page, _ := strconv.ParseUint(query.Page, 10, 64)
	if page == 0 {
		page = 1
	}

	result, totalPage, totalData, err := h.walletService.GetWalletTransactionsByUserID(query, userID)
	if err != nil {
		_ = ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusOK, dto.StatusOKResponse(gin.H{"transactions": result, "total_data": totalData, "total_page": totalPage, "current_data": page, "limit": limit}))
}

func (h *Handler) WalletPin(ctx *gin.Context) {
	payload, _ := ctx.Get("user")
	user, _ := payload.(dto.UserJWT)
	userID := user.UserID

	value, _ := ctx.Get("payload")
	json, _ := value.(*dto.PinReq)
	pin := json.Pin

	err := h.walletService.WalletPin(userID, pin)
	if err != nil {
		_ = ctx.Error(err)
		return
	}
	successResponse := dto.StatusCreatedResponse(err)
	ctx.JSON(http.StatusCreated, successResponse)
}

func (h *Handler) RequestWalletChangeByEmail(ctx *gin.Context) {

	if os.Getenv("ENV") == "testing" {
		fmt.Println("disable cron")
		return
	}

	payload, _ := ctx.Get("user")
	user, _ := payload.(dto.UserJWT)
	userID := user.UserID

	recipientEmail, key, err := h.walletService.RequestPinChangeWithEmail(userID)
	if err != nil {
		_ = ctx.Error(err)
		return
	}
	successResponse := dto.StatusOKResponse(gin.H{"message": "Verification code is send to " + recipientEmail + ", please check your email", "key": key})
	ctx.JSON(http.StatusOK, successResponse)
}

func (h *Handler) ValidateIfRequestByEmailIsValid(ctx *gin.Context) {

	if os.Getenv("ENV") == "testing" {
		fmt.Println("disable cron")
		return
	}

	payload, _ := ctx.Get("user")
	user, _ := payload.(dto.UserJWT)
	userID := user.UserID

	value, _ := ctx.Get("payload")
	json, _ := value.(*dto.KeyRequestByEmailReq)
	key := json.Key

	res, err := h.walletService.ValidateRequestIsValid(userID, key)
	if err != nil {
		_ = ctx.Error(err)
		return
	}
	successResponse := dto.StatusOKResponse(gin.H{"message": res})
	ctx.JSON(http.StatusOK, successResponse)
}

func (h *Handler) ValidateIfRequestChangeByEmailCodeIsValid(ctx *gin.Context) {

	if os.Getenv("ENV") == "testing" {
		fmt.Println("disable cron")
		return
	}

	payload, _ := ctx.Get("user")
	user, _ := payload.(dto.UserJWT)
	userID := user.UserID

	value, _ := ctx.Get("payload")
	json, _ := value.(*dto.CodeKeyRequestByEmailReq)

	res, err := h.walletService.ValidateCodeToRequestByEmail(userID, json)
	if err != nil {
		_ = ctx.Error(err)
		return
	}
	successResponse := dto.StatusOKResponse(gin.H{"message": res})
	ctx.JSON(http.StatusOK, successResponse)
}

func (h *Handler) ChangeWalletPinByEmail(ctx *gin.Context) {

	if os.Getenv("ENV") == "testing" {
		fmt.Println("disable cron")
		return
	}

	payload, _ := ctx.Get("user")
	user, _ := payload.(dto.UserJWT)
	userID := user.UserID

	value, _ := ctx.Get("payload")
	json, _ := value.(*dto.ChangePinByEmailReq)

	res, err := h.walletService.ChangeWalletPinByEmail(userID, json)
	if err != nil {
		_ = ctx.Error(err)
		return
	}
	successResponse := dto.StatusOKResponse(res)
	ctx.JSON(http.StatusOK, successResponse)
}

func (h *Handler) ValidateWalletPin(ctx *gin.Context) {

	if os.Getenv("ENV") == "testing" {
		fmt.Println("disable cron")
		return
	}

	user, exists := ctx.Get("user")
	if !exists {
		_ = ctx.Error(apperror.BadRequestError("User is invalid"))
		return
	}
	userJwt := user.(dto.UserJWT)

	value, _ := ctx.Get("payload")
	json, _ := value.(*dto.PinReq)
	pin := json.Pin

	idToken, result, err := h.walletService.ValidateWalletPin(&userJwt, pin)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}
	status := "success"
	if !result {
		status = "failed"
	}

	successResponse := dto.StatusOKResponse(gin.H{"status": status, "id_token": idToken})
	ctx.JSON(http.StatusOK, successResponse)
}

func (h *Handler) GetWalletStatus(ctx *gin.Context) {

	if os.Getenv("ENV") == "testing" {
		fmt.Println("disable cron")
		return
	}
	user, exists := ctx.Get("user")
	if !exists {
		_ = ctx.Error(apperror.BadRequestError("User is invalid"))
		return
	}
	userID := user.(dto.UserJWT).UserID

	result, err := h.walletService.GetWalletStatus(userID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}

	successResponse := dto.StatusOKResponse(gin.H{"status": result})
	ctx.JSON(http.StatusOK, successResponse)
}

func (h *Handler) CheckoutCart(ctx *gin.Context) {
	user, exists := ctx.Get("user")
	value, _ := ctx.Get("payload")
	json, _ := value.(*dto.CheckoutCartReq)

	if !exists {
		_ = ctx.Error(apperror.BadRequestError("User is invalid"))
		return
	}
	userID := user.(dto.UserJWT).UserID

	result, err := h.walletService.PayOrderWithWallet(userID, json)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}

	successResponse := dto.StatusOKResponse(gin.H{"transaction": result})
	ctx.JSON(http.StatusOK, successResponse)
}
