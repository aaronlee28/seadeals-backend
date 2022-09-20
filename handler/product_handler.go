package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"seadeals-backend/dto"
)

func (h *Handler) FindProductDetailBySlug(ctx *gin.Context) {
	slug := ctx.Param("slug")

	res, err := h.productService.FindProductDetailBySlug(slug)
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, dto.StatusOKResponse(res))
}
