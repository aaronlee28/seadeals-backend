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

func (h *Handler) FindReviewByProductID(ctx *gin.Context) {
	idParam, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		_ = ctx.Error(apperror.BadRequestError("Invalid id format"))
		return
	}

	withImgOnly, _ := strconv.ParseBool(helper.GetQuery(ctx, "withImgOnly", "false"))
	withDescOnly, _ := strconv.ParseBool(helper.GetQuery(ctx, "withDescOnly", "false"))
	queryParam := &model.ReviewQueryParam{
		SortBy:              helper.GetQuery(ctx, "sortBy", model.SortByReviewDefault),
		Sort:                helper.GetQuery(ctx, "sort", model.SortReviewDefault),
		Limit:               helper.GetQueryToUint(ctx, "limit", model.LimitReviewDefault),
		Page:                helper.GetQueryToUint(ctx, "page", model.PageReviewDefault),
		Rating:              helper.GetQueryToUint(ctx, "rating", 0),
		WithImageOnly:       withImgOnly,
		WithDescriptionOnly: withDescOnly,
	}

	res, err := h.reviewService.FindReviewByProductID(uint(idParam), queryParam)
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, dto.StatusOKResponse(res))
}
