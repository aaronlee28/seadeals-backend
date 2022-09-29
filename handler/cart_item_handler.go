package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"seadeals-backend/apperror"
	"seadeals-backend/dto"
	"seadeals-backend/repository"
)

func (h *Handler) DeleteCartItem(ctx *gin.Context) {
	payload, _ := ctx.Get("user")
	user, _ := payload.(dto.UserJWT)
	userID := user.UserID

	value, _ := ctx.Get("payload")
	json, _ := value.(*dto.DeleteFromCartReq)
	if json.UserID != userID {
		ctx.JSON(http.StatusBadRequest, apperror.UnauthorizedError("Cannot delete other user order item"))
		return
	}

	result, err := h.orderItemService.DeleteCartItem(json.CartItemID, json.UserID)
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

	result, err := h.orderItemService.AddToCart(userID, json)
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	successResponse := dto.StatusOKResponse(result)
	ctx.JSON(http.StatusOK, successResponse)
}

func (h *Handler) GetCartItem(ctx *gin.Context) {
	payload, _ := ctx.Get("user")
	user, _ := payload.(dto.UserJWT)
	userID := user.UserID

	query := &repository.Query{
		Limit: ctx.Query("limit"),
	}

	result, totalPage, totalData, err := h.orderItemService.GetCartItems(query, userID)
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	successResponse := dto.StatusOKResponse(gin.H{"total_page": totalPage, "total_data": totalData, "current_page": 1, "limit": 5, "cart_items": result})
	ctx.JSON(http.StatusOK, successResponse)
}
