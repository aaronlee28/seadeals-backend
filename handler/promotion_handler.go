package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"seadeals-backend/dto"
)

func (h *Handler) GetPromotion(ctx *gin.Context) {
	payload, _ := ctx.Get("user")
	user, _ := payload.(dto.UserJWT)
	userID := user.UserID

	result, err := h.promotionService.GetPromotionByUserID(userID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}
	successResponse := dto.StatusOKResponse(result)
	ctx.JSON(http.StatusOK, successResponse)
}

func (h *Handler) CreatePromotion(ctx *gin.Context) {
	payload, _ := ctx.Get("user")
	user, _ := payload.(dto.UserJWT)
	userID := user.UserID

	value, _ := ctx.Get("payload")
	json, _ := value.(*dto.CreatePromotionReq)

	result, err := h.promotionService.CreatePromotion(userID, json)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}
	successResponse := dto.StatusOKResponse(result)
	ctx.JSON(http.StatusOK, successResponse)
}
