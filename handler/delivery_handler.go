package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"seadeals-backend/apperror"
	"seadeals-backend/dto"
)

func (h *Handler) DeliverOrder(ctx *gin.Context) {
	payload, _ := ctx.Get("user")
	user, isValid := payload.(dto.UserJWT)
	if !isValid {
		_ = ctx.Error(apperror.BadRequestError("Invalid user"))
		return
	}

	value, _ := ctx.Get("payload")
	json, _ := value.(*dto.DeliverOrderReq)

	couriers, err := h.deliveryService.DeliverOrder(json, user.UserID)
	if err != nil {
		_ = ctx.Error(err)
		return
	}
	successResponse := dto.StatusOKResponse(couriers)
	ctx.JSON(http.StatusOK, successResponse)
}