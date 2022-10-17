package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"seadeals-backend/dto"
)

func (h *Handler) FavoriteToProduct(ctx *gin.Context) {
	payload, _ := ctx.Get("user")
	user, _ := payload.(dto.UserJWT)
	userID := user.UserID

	value, _ := ctx.Get("payload")
	json, _ := value.(*dto.FavoriteProductReq)

	favorite, favoriteCount, err := h.favoriteService.FavoriteToProduct(userID, json.ProductID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}
	successResponse := dto.StatusOKResponse(gin.H{"favorites": favorite, "new_favorite_count": favoriteCount})
	ctx.JSON(http.StatusOK, successResponse)
}

func (h *Handler) GetUserFavoriteCount(ctx *gin.Context) {
	payload, _ := ctx.Get("user")
	user, _ := payload.(dto.UserJWT)
	userID := user.UserID

	userFavCount, err := h.favoriteService.GetUserFavoriteCount(userID)
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	successResponse := dto.StatusOKResponse(gin.H{"favorite_count": userFavCount})
	ctx.JSON(http.StatusOK, successResponse)
}
