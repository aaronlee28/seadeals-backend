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

//func (a *Handler) TransactionHistory(c *gin.Context) {
//
//	query := &repository.SearchQuery{
//		SortBy:     c.Query("sortBy"),
//		Sort:       c.Query("sort"),
//		Limit:      c.Query("limit"),
//		Page:       c.Query("page"),
//		Search:     c.Query("search"),
//		FilterTime: c.Query("filterTime"),
//		MinAmount:  c.Query("minAmount"),
//		MaxAmount:  c.Query("maxAmount"),
//		Type:       c.Query("type"),
//	}
//
//	result, err := a.WalletService.TransactionHistory(query, userid)
//	if err != nil {
//		e := c.Error(err)
//		c.JSON(http.StatusBadRequest, e)
//		return
//	}
//	successResponse := httpsuccess.OkSuccess("Ok", result)
//	c.JSON(http.StatusOK, successResponse)
//
//}
