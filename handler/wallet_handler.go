package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"seadeals-backend/dto"
	"seadeals-backend/repository"
)

func (h *Handler) WalletDataTransactions(ctx *gin.Context) {
	payload, _ := ctx.Get("user")
	user, _ := payload.(dto.UserJWT)
	userID := user.UserID

	result, err := h.walletService.UserWalletData(userID)
	if err != nil {
		e := ctx.Error(err)
		ctx.JSON(http.StatusBadRequest, e)
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
		e := ctx.Error(err)
		ctx.JSON(http.StatusBadRequest, e)
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
		e := ctx.Error(err)
		ctx.JSON(http.StatusBadRequest, e)
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
	fmt.Println("id", userID)
	err := h.walletService.WalletPin(userID, pin)
	if err != nil {
		e := ctx.Error(err)
		ctx.JSON(http.StatusBadRequest, e)
		return
	}
	successResponse := dto.StatusOKResponse(err)
	ctx.JSON(http.StatusOK, successResponse)

}
