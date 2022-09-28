package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"seadeals-backend/dto"
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
