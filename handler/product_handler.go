package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"seadeals-backend/apperror"
	"seadeals-backend/dto"
	"seadeals-backend/repository"
	"strconv"
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

func (h *Handler) GetProductsBySellerID(ctx *gin.Context) {
	query := map[string]string{
		"page":     ctx.Query("page"),
		"s":        ctx.Query("s"),
		"sortBy":   ctx.Query("sortBy"),
		"sort":     ctx.Query("sort"),
		"limit":    ctx.Query("limit"),
		"minPrice": ctx.Query("minPrice"),
		"maxPrice": ctx.Query("maxPrice"),
	}
	productQuery, err := new(dto.SellerProductSearchQuery).FromQuery(query)
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	idString := ctx.Param("id")
	sellerID, err := strconv.ParseUint(idString, 10, 32)
	if err != nil {
		_ = ctx.Error(apperror.BadRequestError("Seller id is in invalid form"))
		return
	}

	res, totalPage, totalData, err := h.productService.GetProductsBySellerID(productQuery, uint(sellerID))
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	if len(res) == 0 {
		_ = ctx.Error(apperror.NotFoundError("No products were found"))
		return
	}

	ctx.JSON(http.StatusOK, dto.StatusOKResponse(gin.H{"products": res, "total_data": totalData, "total_page": totalPage, "current_page": productQuery.Page, "limit": productQuery.Limit}))
}

func (h *Handler) GetProductsByCategoryID(ctx *gin.Context) {
	query := map[string]string{
		"page":     ctx.Query("page"),
		"s":        ctx.Query("s"),
		"sortBy":   ctx.Query("sortBy"),
		"sort":     ctx.Query("sort"),
		"limit":    ctx.Query("limit"),
		"minPrice": ctx.Query("minPrice"),
		"maxPrice": ctx.Query("maxPrice"),
	}
	productQuery, err := new(dto.SellerProductSearchQuery).FromQuery(query)
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	idString := ctx.Param("id")
	categoryID, err := strconv.ParseUint(idString, 10, 32)
	if err != nil {
		_ = ctx.Error(apperror.BadRequestError("Category id is in invalid form"))
		return
	}

	res, totalPage, totalData, err := h.productService.GetProductsByCategoryID(productQuery, uint(categoryID))
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	if len(res) == 0 {
		_ = ctx.Error(apperror.NotFoundError("No products were found"))
		return
	}

	ctx.JSON(http.StatusOK, dto.StatusOKResponse(gin.H{"products": res, "total_data": totalData, "total_page": totalPage, "current_page": productQuery.Page, "limit": productQuery.Limit}))
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
		City:      c.Query("city"),
		Rating:    c.Query("rating"),
		Category:  c.Query("category"),
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
