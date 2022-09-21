package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"seadeals-backend/apperror"
	"seadeals-backend/dto"
	"strconv"
)

func (h *Handler) FindSellerByID(ctx *gin.Context) {
	idParam, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, apperror.BadRequestError("Invalid id format"))
		return
	}

	seller, err := h.sellerService.FindSellerByID(uint(idParam))
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, dto.StatusOKResponse(seller))
}
