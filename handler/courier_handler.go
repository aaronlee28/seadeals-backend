package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"seadeals-backend/dto"
)

func (h *Handler) GetAllCouriers(ctx *gin.Context) {
	couriers, err := h.courierService.GetAllCouriers()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}
	successResponse := dto.StatusOKResponse(couriers)
	ctx.JSON(http.StatusOK, successResponse)
}
