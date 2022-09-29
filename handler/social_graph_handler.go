package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"seadeals-backend/dto"
)

func (h *Handler) FollowToSeller(ctx *gin.Context) {
	payload, _ := ctx.Get("user")
	user, _ := payload.(dto.UserJWT)
	userID := user.UserID

	value, _ := ctx.Get("payload")
	json, _ := value.(*dto.FollowSellerReq)

	favorite, err := h.socialGraphService.FollowToSeller(userID, json.SellerID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}
	successResponse := dto.StatusOKResponse(favorite)
	ctx.JSON(http.StatusOK, successResponse)
}
