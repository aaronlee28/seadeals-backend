package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"seadeals-backend/apperror"
	"seadeals-backend/dto"
	"strconv"
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

	result, err := h.addressService.UpdateAddress(json, user.(dto.UserJWT).UserID)
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

func (h *Handler) GetUserMainAddress(ctx *gin.Context) {
	payload, _ := ctx.Get("user")
	user, _ := payload.(dto.UserJWT)
	userID := user.UserID

	address, err := h.addressService.GetUserMainAddress(userID)
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, dto.StatusOKResponse(address))
}

func (h *Handler) ChangeMainAddress(ctx *gin.Context) {
	payload, _ := ctx.Get("user")
	user, _ := payload.(dto.UserJWT)
	userID := user.UserID

	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		_ = ctx.Error(apperror.BadRequestError("Invalid id format"))
		return
	}

	address, err := h.addressService.ChangeMainAddress(uint(id), userID)
	if err != nil {
		_ = ctx.Error(apperror.BadRequestError("Failed to change main address"))
		return
	}

	ctx.JSON(http.StatusOK, dto.StatusOKResponse(address))
}
