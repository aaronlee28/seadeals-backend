package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"seadeals-backend/apperror"
	"seadeals-backend/dto"
)

func (h *Handler) CreateNewAddress(ctx *gin.Context) {
	user, exists := ctx.Get("user")
	if !exists {
		_ = ctx.Error(apperror.BadRequestError("User is invalid"))
		return
	}

	value, _ := ctx.Get("payload")
	json, _ := value.(*dto.CreateAddressReq)

	result, err := h.addressService.CreateAddress(json, user.(dto.UserJWT).UserID)
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusCreated, dto.StatusCreatedResponse(result))
}

func (h *Handler) UpdateAddress(ctx *gin.Context) {
	user, exists := ctx.Get("user")
	if !exists {
		_ = ctx.Error(apperror.BadRequestError("User is invalid"))
		return
	}

	value, _ := ctx.Get("payload")
	json, _ := value.(*dto.UpdateAddressReq)
	if user.(dto.UserJWT).UserID != json.UserID {
		_ = ctx.Error(apperror.ForbiddenError("Cannot update another user address"))
		return
	}

	result, err := h.addressService.UpdateAddress(json)
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, dto.StatusOKResponse(result))
}

func (h *Handler) GetAddressesByUserID(ctx *gin.Context) {
	user, exists := ctx.Get("user")
	if !exists {
		_ = ctx.Error(apperror.BadRequestError("User is invalid"))
		return
	}

	result, err := h.addressService.GetAddressesByUserID(user.(dto.UserJWT).UserID)
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, dto.StatusOKResponse(result))
}
