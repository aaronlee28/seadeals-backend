package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"seadeals-backend/apperror"
	"seadeals-backend/dto"
	"strconv"
)

func (h *Handler) GetSubDistrictsByCityID(ctx *gin.Context) {
	id := ctx.Param("id")

	idUint, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		_ = ctx.Error(apperror.BadRequestError("Invalid id format"))
		return
	}

	result, err := h.subDistrictService.GetSubDistrictsByDistrictID(uint(idUint))
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, dto.StatusOKResponse(result))
}
