package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"seadeals-backend/apperror"
	"seadeals-backend/dto"
	"seadeals-backend/helper"
	"seadeals-backend/model"
	"strconv"
)

const (
	SortReviewDefault   = "desc"
	SortByReviewDefault = ""
	LimitReviewDefault  = "6"
	PageReviewDefault   = "1"
)

func (h *Handler) FindReviewByProductID(ctx *gin.Context) {
	idParam, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		_ = ctx.Error(apperror.BadRequestError("Invalid id format"))
		return
	}

	limit, _ := strconv.Atoi(helper.GetQuery(ctx, "limit", LimitReviewDefault))
	page, _ := strconv.Atoi(helper.GetQuery(ctx, "page", PageReviewDefault))
	queryParam := &model.ReviewQueryParam{
		SortBy: helper.GetQuery(ctx, "sortBy", SortByReviewDefault),
		Sort:   helper.GetQuery(ctx, "sort", SortReviewDefault),
		Limit:  limit,
		Page:   page,
	}

	res, err := h.reviewService.FindReviewByProductID(uint(idParam), queryParam)
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, dto.StatusOKResponse(res))
}
