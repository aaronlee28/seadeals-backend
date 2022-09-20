package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"seadeals-backend/dto"
	"seadeals-backend/model"
)

func (h *Handler) WalletDataTransactions(ctx *gin.Context) {
	payload, _ := ctx.Get("user")
	user, _ := payload.(model.User)
	userID := user.ID

	result, err := h.walletService.UserWalletData(userID)
	if err != nil {
		e := ctx.Error(err)
		ctx.JSON(http.StatusBadRequest, e)
		return
	}

	successResponse := dto.StatusOKResponse(result)
	ctx.JSON(http.StatusOK, successResponse)

}
