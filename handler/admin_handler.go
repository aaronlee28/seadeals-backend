package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"seadeals-backend/dto"
)

func (h *Handler) CreateGlobalVoucher(ctx *gin.Context) {

	value, _ := ctx.Get("payload")
	json, _ := value.(*dto.CreateGlobalVoucher)
	res, err := h.adminService.CreateGlobalVoucher(json)

	if err != nil {
		_ = ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, dto.StatusOKResponse(res))
}

func (h *Handler) CreateCategory(ctx *gin.Context) {
	value, _ := ctx.Get("payload")
	json, _ := value.(*dto.CreateCategory)
	res, err := h.adminService.CreateCategory(json)

	if err != nil {
		_ = ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, dto.StatusOKResponse(res))
}
