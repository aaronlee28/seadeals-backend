package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"seadeals-backend/apperror"
	"seadeals-backend/dto"
	"seadeals-backend/helper"
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

func (a *Handler) SearchRecommendProduct(c *gin.Context) {

	query := &repository.SearchQuery{
		Search:    helper.GetQuery(c, "sortBy", ""),
		SortBy:    helper.GetQuery(c, "sortBy", "bought"),
		Sort:      helper.GetQuery(c, "sort", SortByReviewDefault),
		Limit:     helper.GetQuery(c, "limit", "30"),
		Page:      helper.GetQuery(c, "page", "1"),
		MinAmount: helper.GetQuery(c, "minAmount", "0"),
		MaxAmount: helper.GetQuery(c, "maxAmount", "99999999999"),
		City:      helper.GetQuery(c, "city", ""),
		Rating:    helper.GetQuery(c, "rating", ""),
		Category:  helper.GetQuery(c, "category", ""),
	}

	result, err := a.productService.SearchRecommendProduct(query)
	if err != nil {
		e := c.Error(err)
		c.JSON(http.StatusBadRequest, e)
		return
	}
	successResponse := dto.StatusOKResponse(result)
	c.JSON(http.StatusOK, successResponse)

}
