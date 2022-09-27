package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"seadeals-backend/apperror"
	"seadeals-backend/dto"
	"seadeals-backend/repository"
)

func (h *Handler) WalletDataTransactions(ctx *gin.Context) {
	payload, _ := ctx.Get("user")
	user, _ := payload.(dto.UserJWT)
	userID := user.UserID

	result, err := h.walletService.UserWalletData(userID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}
	successResponse := dto.StatusOKResponse(result)
	ctx.JSON(http.StatusOK, successResponse)
}

func (h *Handler) TransactionDetails(ctx *gin.Context) {
	value, _ := ctx.Get("payload")
	json, _ := value.(*dto.TransactionDetailsReq)
	transactionID := json.TransactionID

	result, err := h.walletService.TransactionDetails(transactionID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err)
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
		ctx.JSON(http.StatusBadRequest, err)
		return
	}
	successResponse := dto.StatusOKResponse(result)
	ctx.JSON(http.StatusOK, successResponse)
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
		ctx.JSON(http.StatusBadRequest, err)
		return
	}
	successResponse := dto.StatusCreatedResponse(err)
	ctx.JSON(http.StatusCreated, successResponse)
}

func (h *Handler) RequestWalletChangeByEmail(ctx *gin.Context) {
	payload, _ := ctx.Get("user")
	user, _ := payload.(dto.UserJWT)
	userID := user.UserID

	res, key, err := h.walletService.RequestPinChangeWithEmail(userID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}
	successResponse := dto.StatusOKResponse(gin.H{"mailjet_response": res, "key": key})
	ctx.JSON(http.StatusOK, successResponse)
}

func (h *Handler) ValidateIfRequestByEmailIsValid(ctx *gin.Context) {
	payload, _ := ctx.Get("user")
	user, _ := payload.(dto.UserJWT)
	userID := user.UserID

	value, _ := ctx.Get("payload")
	json, _ := value.(*dto.KeyRequestByEmailReq)
	key := json.Key

	res, err := h.walletService.ValidateRequestIsValid(userID, key)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}
	successResponse := dto.StatusOKResponse(gin.H{"message": res})
	ctx.JSON(http.StatusOK, successResponse)
}

func (h *Handler) ValidateIfRequestChangeByEmailCodeIsValid(ctx *gin.Context) {
	payload, _ := ctx.Get("user")
	user, _ := payload.(dto.UserJWT)
	userID := user.UserID

	value, _ := ctx.Get("payload")
	json, _ := value.(*dto.CodeKeyRequestByEmailReq)

	res, err := h.walletService.ValidateCodeToRequestByEmail(userID, json)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}
	successResponse := dto.StatusOKResponse(gin.H{"message": res})
	ctx.JSON(http.StatusOK, successResponse)
}

func (h *Handler) ChangeWalletPinByEmail(ctx *gin.Context) {
	payload, _ := ctx.Get("user")
	user, _ := payload.(dto.UserJWT)
	userID := user.UserID

	value, _ := ctx.Get("payload")
	json, _ := value.(*dto.ChangePinByEmailReq)

	res, err := h.walletService.ChangeWalletPinByEmail(userID, json)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}
	successResponse := dto.StatusOKResponse(res)
	ctx.JSON(http.StatusOK, successResponse)
}

func (h *Handler) ValidateWalletPin(ctx *gin.Context) {
	user, exists := ctx.Get("user")
	if !exists {
		_ = ctx.Error(apperror.BadRequestError("User is invalid"))
		return
	}
	userID := user.(dto.UserJWT).UserID

	value, _ := ctx.Get("payload")
	json, _ := value.(*dto.PinReq)
	pin := json.Pin

	result, err := h.walletService.ValidateWalletPin(userID, pin)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}
	status := "success"
	if !result {
		status = "failed"
	}

	successResponse := dto.StatusOKResponse(gin.H{"status": status})
	ctx.JSON(http.StatusOK, successResponse)
}

func (h *Handler) GetWalletStatus(ctx *gin.Context) {
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

	result, err := h.walletService.CheckoutCart(userID, json)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}

	successResponse := dto.StatusOKResponse(gin.H{"status": result})
	ctx.JSON(http.StatusOK, successResponse)
}
