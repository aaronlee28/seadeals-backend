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
	sortReviewDefault   = "desc"
	sortByReviewDefault = ""
)

func (h *Handler) FindReviewByProductID(ctx *gin.Context) {
	idParam, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		_ = ctx.Error(apperror.BadRequestError("Invalid id format"))
		return
	}

	queryParam := &model.ReviewQueryParam{
		SortBy: helper.GetQuery(ctx, "sortBy", sortByReviewDefault),
		Sort:   helper.GetQuery(ctx, "sort", sortReviewDefault),
	}

	res, err := h.reviewService.FindReviewByProductID(uint(idParam), queryParam)
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, dto.StatusOKResponse(res))
}
