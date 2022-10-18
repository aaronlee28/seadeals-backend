package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"seadeals-backend/apperror"
	"seadeals-backend/dto"
	"strconv"
)

func (h *Handler) CreateOrUpdateSellerAvailableCour(ctx *gin.Context) {
	payload, _ := ctx.Get("user")
	user, isValid := payload.(dto.UserJWT)
	if !isValid {
		_ = ctx.Error(apperror.BadRequestError("User is invalid"))
		return
	}

	value, _ := ctx.Get("payload")
	json, _ := value.(*dto.AddDeliveryReq)

	result, err := h.sellerAvailableCourServ.CreateOrUpdateCourier(json, user.UserID)
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, dto.StatusOKResponse(result))
}

func (h *Handler) GetAvailableCourierForSeller(ctx *gin.Context) {
	payload, _ := ctx.Get("user")
	user, isValid := payload.(dto.UserJWT)
	if !isValid {
		_ = ctx.Error(apperror.BadRequestError("User is invalid"))
		return
	}

	result, err := h.sellerAvailableCourServ.GetAvailableCourierForSeller(user.UserID)
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, dto.StatusOKResponse(result))
}

func (h *Handler) GetAvailableCourierForBuyer(ctx *gin.Context) {
	idParam := ctx.Param("id")
	sellerID, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		_ = ctx.Error(apperror.BadRequestError("Seller id tidak valid"))
		return
	}

	payload, _ := ctx.Get("user")
	_, isValid := payload.(dto.UserJWT)
	if !isValid {
		_ = ctx.Error(apperror.BadRequestError("User tidak valid"))
		return
	}

	result, err := h.sellerAvailableCourServ.GetAvailableCourierForBuyer(uint(sellerID))
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, dto.StatusOKResponse(result))
}
