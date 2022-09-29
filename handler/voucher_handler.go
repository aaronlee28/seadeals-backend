package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"seadeals-backend/apperror"
	"seadeals-backend/dto"
	"seadeals-backend/helper"
	"seadeals-backend/model"
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

func (h *Handler) FindVoucherDetailByID(ctx *gin.Context) {
	userJWT, _ := ctx.Get("user")
	user := userJWT.(dto.UserJWT)

	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		_ = ctx.Error(apperror.BadRequestError("Invalid id format"))
		return
	}

	voucher, err := h.voucherService.FindVoucherDetailByID(uint(id), user.UserID)
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, dto.StatusOKResponse(voucher))
}

func (h *Handler) FindVoucherByID(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		_ = ctx.Error(apperror.BadRequestError("Invalid id format"))
		return
	}

	voucher, err := h.voucherService.FindVoucherByID(uint(id))
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, dto.StatusOKResponse(voucher))
}

func (h *Handler) FindVoucherBySellerID(ctx *gin.Context) {
	userJWT, _ := ctx.Get("user")
	user := userJWT.(dto.UserJWT)

	sellerID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		_ = ctx.Error(apperror.BadRequestError("Invalid id format"))
		return
	}

	qp := &model.VoucherQueryParam{
		SortBy: helper.GetQuery(ctx, "sortBy", model.SortByVoucherDefault),
		Sort:   helper.GetQuery(ctx, "sort", model.SortVoucherDefault),
		Limit:  helper.GetQueryToUint(ctx, "limit", model.LimitVoucherDefault),
		Page:   helper.GetQueryToUint(ctx, "page", model.PageVoucherDefault),
	}
	vouchers, err := h.voucherService.FindVoucherBySellerID(uint(sellerID), user.UserID, qp)
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, dto.StatusOKResponse(vouchers))
}

func (h *Handler) ValidateVoucher(ctx *gin.Context) {
	payload, _ := ctx.Get("payload")
	req := payload.(*dto.PostValidateVoucherReq)

	voucher, err := h.voucherService.ValidateVoucher(req)
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
