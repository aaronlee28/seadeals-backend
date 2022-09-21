package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"seadeals-backend/dto"
)

func (h *Handler) FindAllProductCategories(ctx *gin.Context) {
	categories, err := h.productCategoryService.FindAllProductCategories()
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, dto.StatusOKResponse(categories))
}
