package handler

import (
	"github.com/gin-gonic/gin"
	"math"
	"net/http"
	"seadeals-backend/apperror"
	"seadeals-backend/dto"
	"seadeals-backend/helper"
	"seadeals-backend/model"
	"seadeals-backend/repository"
	"strconv"
)

func (h *Handler) FindProductDetailByID(ctx *gin.Context) {
	idParam := ctx.Param("id")
	productID, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		_ = ctx.Error(apperror.BadRequestError("Product id is invalid"))
		return
	}

	payload, _ := ctx.Get("user")
	user, isValid := payload.(dto.UserJWT)
	userID := uint(0)
	if isValid {
		userID = user.UserID
	}

	res, err := h.productService.FindProductDetailByID(uint(productID), userID)
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, dto.StatusOKResponse(res))
}

func (h *Handler) FindSimilarProduct(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		_ = ctx.Error(apperror.BadRequestError("Invalid id format"))
		return
	}

	products, err := h.productService.FindSimilarProducts(uint(id))
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, dto.StatusOKResponse(products))
}

func (h *Handler) GetProductsBySellerID(ctx *gin.Context) {
	query := map[string]string{
		"page":      ctx.Query("page"),
		"s":         ctx.Query("s"),
		"sortBy":    ctx.Query("sortBy"),
		"sort":      ctx.Query("sort"),
		"limit":     ctx.Query("limit"),
		"minAmount": ctx.Query("minAmount"),
		"maxAmount": ctx.Query("maxAmount"),
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

	ctx.JSON(http.StatusOK, dto.StatusOKResponse(gin.H{"products": res, "total_data": totalData, "total_page": totalPage, "current_page": productQuery.Page, "limit": productQuery.Limit}))
}

func (h *Handler) GetProductsByCategoryID(ctx *gin.Context) {
	query := map[string]string{
		"page":      ctx.Query("page"),
		"s":         ctx.Query("s"),
		"sortBy":    ctx.Query("sortBy"),
		"sort":      ctx.Query("sort"),
		"limit":     ctx.Query("limit"),
		"minAmount": ctx.Query("minAmount"),
		"maxAmount": ctx.Query("maxAmount"),
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

	ctx.JSON(http.StatusOK, dto.StatusOKResponse(gin.H{"products": res, "total_data": totalData, "total_page": totalPage, "current_page": productQuery.Page, "limit": productQuery.Limit}))
}

func (h *Handler) SearchProducts(ctx *gin.Context) {
	query := &repository.SearchQuery{
		Search:     helper.GetQuery(ctx, "s", ""),
		SortBy:     helper.GetQuery(ctx, "sortBy", ""),
		Sort:       helper.GetQuery(ctx, "sort", model.SortByReviewDefault),
		Limit:      helper.GetQuery(ctx, "limit", "20"),
		Page:       helper.GetQuery(ctx, "page", "1"),
		MinAmount:  helper.GetQueryToFloat64(ctx, "minAmount", 0),
		MaxAmount:  helper.GetQueryToFloat64(ctx, "maxAmount", math.MaxFloat64),
		City:       helper.GetQuery(ctx, "city", ""),
		Rating:     helper.GetQuery(ctx, "rating", "0"),
		Category:   helper.GetQuery(ctx, "category", ""),
		CategoryID: helper.GetQueryToUint(ctx, "categoryID", 0),
		SellerID:   helper.GetQueryToUint(ctx, "sellerID", 0),
	}
	limit, _ := strconv.ParseUint(query.Limit, 10, 64)
	if limit == 0 {
		limit = 20
	}
	page, _ := strconv.ParseUint(query.Page, 10, 64)
	if page == 0 {
		page = 1
	}

	result, totalPage, totalData, err := h.productService.GetProducts(query)
	if err != nil {
		_ = ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusOK, dto.StatusOKResponse(gin.H{"products": result, "total_data": totalData, "total_page": totalPage, "current_data": page, "limit": limit}))
}

func (h *Handler) SearchRecommendProduct(ctx *gin.Context) {
	query := &repository.SearchQuery{
		Search:     helper.GetQuery(ctx, "s", ""),
		SortBy:     helper.GetQuery(ctx, "sortBy", "total_sold"),
		Sort:       helper.GetQuery(ctx, "sort", model.SortByReviewDefault),
		Limit:      helper.GetQuery(ctx, "limit", "30"),
		Page:       helper.GetQuery(ctx, "page", "1"),
		MinAmount:  helper.GetQueryToFloat64(ctx, "minAmount", 0),
		MaxAmount:  helper.GetQueryToFloat64(ctx, "maxAmount", math.MaxFloat64),
		City:       helper.GetQuery(ctx, "city", ""),
		Rating:     helper.GetQuery(ctx, "rating", "0"),
		Category:   helper.GetQuery(ctx, "category", ""),
		CategoryID: helper.GetQueryToUint(ctx, "categoryID", 0),
		SellerID:   helper.GetQueryToUint(ctx, "sellerID", 0),
	}

	result, err := h.productService.SearchRecommendProduct(query)
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	successResponse := dto.StatusOKResponse(result)
	ctx.JSON(http.StatusOK, successResponse)

}
