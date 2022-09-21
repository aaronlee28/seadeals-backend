package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"seadeals-backend/apperror"
	"seadeals-backend/dto"
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
		"page":   ctx.Query("page"),
		"s":      ctx.Query("s"),
		"sortBy": ctx.Query("sortBy"),
		"sort":   ctx.Query("sort"),
		"limit":  ctx.Query("limit"),
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
