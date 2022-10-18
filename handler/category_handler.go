package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"seadeals-backend/dto"
	"seadeals-backend/helper"
	"seadeals-backend/model"
	"strconv"
)

func (h *Handler) FindCategories(ctx *gin.Context) {
	query := &model.CategoryQuery{
		Search:   helper.GetQuery(ctx, "s", ""),
		Limit:    helper.GetQuery(ctx, "limit", "0"),
		Page:     helper.GetQuery(ctx, "page", "1"),
		SellerID: helper.GetQueryToUint(ctx, "sellerID", 0),
		ParentID: helper.GetQueryToUint(ctx, "parentID", 0),
		FindAll:  helper.GetQueryToBool(ctx, "findAll", false),
	}
	categories, totalPage, totalData, err := h.productCategoryService.FindCategories(query)
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	limit, _ := strconv.ParseUint(query.Limit, 10, 64)
	page, _ := strconv.ParseUint(query.Page, 10, 64)
	if page == 0 {
		page = 1
	}

	ctx.JSON(http.StatusOK, dto.StatusOKResponse(gin.H{"categories": categories, "current_page": page, "limit": limit, "total_page": totalPage, "total_data": totalData}))
}
