package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"seadeals-backend/apperror"
	"seadeals-backend/dto"
	"strconv"
)

func (h *Handler) CreateVoucher(ctx *gin.Context) {
	userJWT, _ := ctx.Get("user")
	user := userJWT.(dto.UserJWT)

	payload, _ := ctx.Get("payload")
	req := payload.(*dto.PostVoucherReq)

	voucher, err := h.voucherService.CreateVoucher(req, user.UserID)
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, dto.StatusOKResponse(voucher))
}

func (h *Handler) UpdateVoucher(ctx *gin.Context) {
	userJWT, _ := ctx.Get("user")
	user := userJWT.(dto.UserJWT)

	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		_ = ctx.Error(apperror.BadRequestError("Invalid id format"))
		return
	}

	payload, _ := ctx.Get("payload")
	req := payload.(*dto.PatchVoucherReq)

	voucher, err := h.voucherService.UpdateVoucher(req, uint(id), user.UserID)
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, dto.StatusOKResponse(voucher))
}

func (h *Handler) DeleteVoucherByID(ctx *gin.Context) {
	userJWT, _ := ctx.Get("user")
	user := userJWT.(dto.UserJWT)

	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		_ = ctx.Error(apperror.BadRequestError("Invalid id format"))
		return
	}

	isDeleted, err := h.voucherService.DeleteVoucherByID(uint(id), user.UserID)
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, dto.StatusOKResponse(gin.H{"is_deleted": isDeleted}))
}
