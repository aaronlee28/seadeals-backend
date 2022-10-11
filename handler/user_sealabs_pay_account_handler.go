package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"seadeals-backend/apperror"
	"seadeals-backend/dto"
	"strconv"
)

func (h *Handler) RegisterSeaLabsPayAccount(ctx *gin.Context) {
	user, exists := ctx.Get("user")
	if !exists {
		_ = ctx.Error(apperror.BadRequestError("User is invalid"))
		return
	}

	value, _ := ctx.Get("payload")
	json, _ := value.(*dto.RegisterSeaLabsPayReq)

	result, err := h.seaLabsPayAccServ.RegisterSeaLabsPayAccount(json, user.(dto.UserJWT).UserID)
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusCreated, dto.StatusCreatedResponse(result))
}

func (h *Handler) CheckSeaLabsPayAccount(ctx *gin.Context) {
	user, exists := ctx.Get("user")
	if !exists {
		_ = ctx.Error(apperror.BadRequestError("User is invalid"))
		return
	}

	value, _ := ctx.Get("payload")
	json, _ := value.(*dto.CheckSeaLabsPayReq)

	result, err := h.seaLabsPayAccServ.CheckSeaLabsAccountExists(json, user.(dto.UserJWT).UserID)
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, dto.StatusOKResponse(result))
}

func (h *Handler) UpdateSeaLabsPayToMain(ctx *gin.Context) {
	user, exists := ctx.Get("user")
	if !exists {
		_ = ctx.Error(apperror.BadRequestError("User is invalid"))
		return
	}

	value, _ := ctx.Get("payload")
	json, _ := value.(*dto.UpdateSeaLabsPayToMainReq)

	result, err := h.seaLabsPayAccServ.UpdateSeaLabsAccountToMain(json, user.(dto.UserJWT).UserID)
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, dto.StatusOKResponse(result))
}

func (h *Handler) GetSeaLabsPayAccount(ctx *gin.Context) {
	user, exists := ctx.Get("user")
	if !exists {
		_ = ctx.Error(apperror.BadRequestError("User is invalid"))
		return
	}

	result, err := h.seaLabsPayAccServ.GetSeaLabsAccountByUserID(user.(dto.UserJWT).UserID)
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, dto.StatusOKResponse(result))
}

func (h *Handler) PayWithSeaLabsPay(ctx *gin.Context) {
	payload, _ := ctx.Get("user")
	user, _ := payload.(dto.UserJWT)
	userID := user.UserID

	value, _ := ctx.Get("payload")
	json, _ := value.(*dto.CheckoutCartReq)

	if json.AccountNumber == "" {
		_ = ctx.Error(apperror.BadRequestError("Invalid account number"))
		return
	}

	redirectURL, transaction, err := h.seaLabsPayAccServ.PayWithSeaLabsPay(userID, json)
	if err != nil {
		_ = ctx.Error(err)
		return
	}
	if redirectURL == "" {
		_ = ctx.Error(apperror.InternalServerError("Cannot send url"))
		return
	}

	ctx.JSON(http.StatusOK, dto.StatusOKResponse(gin.H{"redirect_url": redirectURL, "transaction": transaction}))
}

func (h *Handler) PayWithSeaLabsPayCallback(ctx *gin.Context) {
	value, _ := ctx.Get("payload")
	json, _ := value.(*dto.SeaLabsPayReq)
	txnID, err := strconv.ParseUint(json.TxnID, 10, 64)
	if err != nil {
		_ = ctx.Error(apperror.BadRequestError("Cannot convert txnID to uint"))
		return
	}
	response, err := h.seaLabsPayAccServ.PayWithSeaLabsPayCallback(uint(txnID), json.Status)
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, dto.StatusOKResponse(response))
}

func (h *Handler) TopUpWithSeaLabsPay(ctx *gin.Context) {
	user, exists := ctx.Get("user")
	if !exists {
		_ = ctx.Error(apperror.BadRequestError("User is invalid"))
		return
	}

	value, _ := ctx.Get("payload")
	json, _ := value.(*dto.TopUpWalletWithSeaLabsPayReq)
	response, redirectURL, err := h.seaLabsPayAccServ.TopUpWithSeaLabsPay(json.Amount, user.(dto.UserJWT).UserID, json.AccountNumber)
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, dto.StatusOKResponse(gin.H{"status": response, "redirect_url": redirectURL}))
}

func (h *Handler) TopUpWithSeaLabsPayCallback(ctx *gin.Context) {
	value, _ := ctx.Get("payload")
	json, _ := value.(*dto.SeaLabsPayReq)
	txnID, err := strconv.ParseUint(json.TxnID, 10, 64)
	if err != nil {
		_ = ctx.Error(apperror.BadRequestError("Cannot convert txnID to uint"))
		return
	}
	response, err := h.seaLabsPayAccServ.TopUpWithSeaLabsPayCallback(uint(txnID), json.Status)
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, dto.StatusOKResponse(response))
}
