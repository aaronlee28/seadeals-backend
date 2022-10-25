package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"seadeals-backend/apperror"
	"seadeals-backend/dto"
	"seadeals-backend/helper"
	"seadeals-backend/repository"
	"strconv"
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

func (h *Handler) RequestRefundByBuyer(ctx *gin.Context) {
	payload, _ := ctx.Get("user")
	user, isValid := payload.(dto.UserJWT)
	if !isValid {
		_ = ctx.Error(apperror.BadRequestError("Invalid user"))
		return
	}

	value, _ := ctx.Get("payload")
	json, _ := value.(*dto.CreateComplaintReq)

	response, err := h.orderService.RequestRefundByBuyer(json, user.UserID)
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, dto.StatusOKResponse(response))
}

func (h *Handler) AcceptRefundRequest(ctx *gin.Context) {
	payload, _ := ctx.Get("user")
	user, isValid := payload.(dto.UserJWT)
	if !isValid {
		_ = ctx.Error(apperror.BadRequestError("Invalid user"))
		return
	}

	value, _ := ctx.Get("payload")
	json, _ := value.(*dto.RejectAcceptRefundReq)

	response, err := h.orderService.AcceptRefundRequest(json, user.UserID)
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, dto.StatusOKResponse(response))
}

func (h *Handler) RejectRefundRequest(ctx *gin.Context) {
	payload, _ := ctx.Get("user")
	user, isValid := payload.(dto.UserJWT)
	if !isValid {
		_ = ctx.Error(apperror.BadRequestError("Invalid user"))
		return
	}

	value, _ := ctx.Get("payload")
	json, _ := value.(*dto.RejectAcceptRefundReq)

	response, err := h.orderService.RejectRefundRequest(json, user.UserID)
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, dto.StatusOKResponse(response))
}

func (h *Handler) FinishOrder(ctx *gin.Context) {
	payload, _ := ctx.Get("user")
	user, isValid := payload.(dto.UserJWT)
	if !isValid {
		_ = ctx.Error(apperror.BadRequestError("Invalid user"))
		return
	}

	value, _ := ctx.Get("payload")
	json, _ := value.(*dto.FinishOrderReq)

	response, err := h.orderService.FinishOrder(json, user.UserID)
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, dto.StatusOKResponse(response))
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

func (h *Handler) GetBuyerOrders(ctx *gin.Context) {
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

	result, totalPage, totalData, err := h.orderService.GetOrderByUserID(user.UserID, query)
	if err != nil {
		_ = ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusOK, dto.StatusOKResponse(gin.H{"orders": result, "total_data": totalData, "total_page": totalPage, "current_page": query.Page, "limit": query.Limit}))
}

func (h *Handler) GetDetailOrderForThermal(ctx *gin.Context) {
	idParam := ctx.Param("id")
	orderID, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		_ = ctx.Error(apperror.BadRequestError("Order id tidak valid"))
		return
	}

	payload, _ := ctx.Get("user")
	user, isValid := payload.(dto.UserJWT)
	if !isValid {
		_ = ctx.Error(apperror.BadRequestError("Invalid user"))
		return
	}

	result, err := h.orderService.GetDetailOrderForThermal(uint(orderID), user.UserID)
	if err != nil {
		_ = ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusOK, dto.StatusOKResponse(result))
}

func (h *Handler) GetDetailOrderForReceipt(ctx *gin.Context) {
	idParam := ctx.Param("id")
	orderID, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		_ = ctx.Error(apperror.BadRequestError("Order id tidak valid"))
		return
	}

	payload, _ := ctx.Get("user")
	user, isValid := payload.(dto.UserJWT)
	if !isValid {
		_ = ctx.Error(apperror.BadRequestError("Invalid user"))
		return
	}

	result, err := h.orderService.GetDetailOrderForReceipt(uint(orderID), user.UserID)
	if err != nil {
		_ = ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusOK, dto.StatusOKResponse(result))
}

func (h *Handler) GetTotalPredictedPrice(ctx *gin.Context) {
	payload, _ := ctx.Get("user")
	user, isValid := payload.(dto.UserJWT)
	if !isValid {
		_ = ctx.Error(apperror.BadRequestError("Invalid user"))
		return
	}

	value, _ := ctx.Get("payload")
	json, _ := value.(*dto.PredictedPriceReq)

	response, err := h.orderService.GetTotalPredictedPrice(json, user.UserID)
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, dto.StatusOKResponse(response))
}

func (h *Handler) GetOrderByID(ctx *gin.Context) {
	payload, _ := ctx.Get("user")
	user, _ := payload.(dto.UserJWT)
	userID := user.UserID

	idParam, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		_ = ctx.Error(apperror.BadRequestError("Invalid id format"))
		return
	}
	result, err := h.orderService.GetOrderByID(userID, uint(idParam))
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	successResponse := dto.StatusOKResponse(result)
	ctx.JSON(http.StatusOK, successResponse)
}
