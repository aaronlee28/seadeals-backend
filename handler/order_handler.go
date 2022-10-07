package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"seadeals-backend/apperror"
	"seadeals-backend/dto"
	"seadeals-backend/helper"
	"seadeals-backend/repository"
)

func (h *Handler) CancelOrderBySeller(ctx *gin.Context) {
	payload, _ := ctx.Get("user")
	user, isValid := payload.(dto.UserJWT)
	if !isValid {
		_ = ctx.Error(apperror.BadRequestError("Invalid user"))
		return
	}

	value, _ := ctx.Get("payload")
	json, _ := value.(*dto.SellerCancelOrderReq)

	message, err := h.orderService.CancelOrderBySeller(json.OrderID, user.UserID)
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, dto.StatusOKResponse(gin.H{"message": message}))
}

func (h *Handler) GetSellerOrders(ctx *gin.Context) {
	query := &repository.OrderQuery{
		Filter: helper.GetQuery(ctx, "filter", ""),
		Limit:  helper.GetQueryToInt(ctx, "limit", 10),
		Page:   helper.GetQueryToInt(ctx, "page", 1),
	}

	payload, _ := ctx.Get("user")
	user, isValid := payload.(dto.UserJWT)
	if !isValid {
		_ = ctx.Error(apperror.BadRequestError("Invalid user"))
		return
	}

	result, totalPage, totalData, err := h.orderService.GetOrderBySellerID(user.UserID, query)
	if err != nil {
		_ = ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusOK, dto.StatusOKResponse(gin.H{"orders": result, "total_data": totalData, "total_page": totalPage, "current_page": query.Page, "limit": query.Limit}))
}
