package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"seadeals-backend/apperror"
	"seadeals-backend/dto"
	"strconv"
)

func (h *Handler) GetPromotion(ctx *gin.Context) {
	payload, _ := ctx.Get("user")
	user, _ := payload.(dto.UserJWT)
	userID := user.UserID

	result, err := h.promotionService.GetPromotionByUserID(userID)
	if err != nil {
		_ = ctx.Error(err)
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

func (h *Handler) ViewDetailPromotionByID(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		_ = ctx.Error(apperror.BadRequestError("Invalid id format"))
		return
	}
	uintID := uint(id)

	result, err := h.promotionService.ViewDetailPromotionByID(uintID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}
	successResponse := dto.StatusOKResponse(result)
	ctx.JSON(http.StatusOK, successResponse)
}
