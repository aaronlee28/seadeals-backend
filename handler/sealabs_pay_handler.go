package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"seadeals-backend/dto"
)

func (h *Handler) CreateSignature(ctx *gin.Context) {
	value, _ := ctx.Get("payload")
	json, _ := value.(*dto.SeaDealspayReq)
	result := h.sealabsPayService.CreateSignature(json)

	successResponse := dto.StatusOKResponse(result)
	ctx.JSON(http.StatusOK, successResponse)
}
