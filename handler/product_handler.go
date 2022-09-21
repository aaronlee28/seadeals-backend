package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"seadeals-backend/dto"
	"seadeals-backend/repository"
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

func (a *Handler) SearchProduct(c *gin.Context) {

	query := &repository.SearchQuery{
		Search:    c.Query("search"),
		SortBy:    c.Query("sortBy"),
		Sort:      c.Query("sort"),
		Limit:     c.Query("limit"),
		Page:      c.Query("page"),
		MinAmount: c.Query("minAmount"),
		MaxAmount: c.Query("maxAmount"),
	}

	result, err := a.productService.SearchProduct(query)
	if err != nil {
		e := c.Error(err)
		c.JSON(http.StatusBadRequest, e)
		return
	}
	successResponse := dto.StatusOKResponse(result)
	c.JSON(http.StatusOK, successResponse)

}
