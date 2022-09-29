package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"seadeals-backend/apperror"
	"seadeals-backend/dto"
	"strconv"
)

func (h *Handler) FindSellerByID(ctx *gin.Context) {
	sellerID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		_ = ctx.Error(apperror.BadRequestError("Invalid id format"))
		return
	}

	payload, _ := ctx.Get("user")
	user, isValid := payload.(dto.UserJWT)
	userID := uint(0)
	if isValid {
		userID = user.UserID
	}

	seller, err := h.sellerService.FindSellerByID(uint(sellerID), userID)
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, dto.StatusOKResponse(seller))
}
