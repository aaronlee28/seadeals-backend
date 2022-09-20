package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"seadeals-backend/apperror"
	"seadeals-backend/dto"
	"seadeals-backend/repository"
)

func (h *Handler) DeleteOrderItem(ctx *gin.Context) {
	payload, _ := ctx.Get("user")
	user, _ := payload.(dto.UserJWT)
	userID := user.UserID

	value, _ := ctx.Get("payload")
	json, _ := value.(*dto.DeleteFromCartReq)
	if json.UserID != userID {
		ctx.JSON(http.StatusBadRequest, apperror.UnauthorizedError("Cannot delete other user order item"))
		return
	}

	result, err := h.orderItemService.DeleteOrderItem(json.OrderItemID, json.UserID)
	if err != nil {
		_ = ctx.Error(err)
		return
	}
	successResponse := dto.StatusOKResponse(result)
	ctx.JSON(http.StatusOK, successResponse)

}

func (h *Handler) AddToCart(ctx *gin.Context) {
	payload, _ := ctx.Get("user")
	user, _ := payload.(dto.UserJWT)
	userID := user.UserID

	value, _ := ctx.Get("payload")
	json, _ := value.(*dto.AddToCartReq)
	if json.UserID != userID {
		ctx.JSON(http.StatusBadRequest, apperror.UnauthorizedError("Cannot add other user order item"))
		return
	}

	result, err := h.orderItemService.AddToCart(json)
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	successResponse := dto.StatusOKResponse(result)
	ctx.JSON(http.StatusOK, successResponse)
}

func (h *Handler) GetOrderItem(ctx *gin.Context) {
	payload, _ := ctx.Get("user")
	user, _ := payload.(dto.UserJWT)
	userID := user.UserID

	query := &repository.Query{
		Limit: ctx.Query("limit"),
	}

	result, totalPage, totalData, err := h.orderItemService.GetOrderItem(query, userID)
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	successResponse := dto.StatusOKResponse(gin.H{"total_page": totalPage, "total_data": totalData, "current_page": 1, "limit": 5, "order_items": result})
	ctx.JSON(http.StatusOK, successResponse)
}
