package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"seadeals-backend/apperror"
	"seadeals-backend/dto"
)

func (h *Handler) RegisterSeaLabsPayAccount(ctx *gin.Context) {
	user, exists := ctx.Get("user")
	if !exists {
		_ = ctx.Error(apperror.BadRequestError("User is invalid"))
		return
	}

	value, _ := ctx.Get("payload")
	json, _ := value.(*dto.RegisterSeaLabsPayReq)
	if user.(dto.UserJWT).UserID != json.UserID {
		_ = ctx.Error(apperror.ForbiddenError("Cannot register another user Sea Labs Pay Account"))
		return
	}

	result, err := h.seaLabsPayAccServ.RegisterSeaLabsPayAccount(json)
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusCreated, dto.StatusCreatedResponse(result))
}

func (h *Handler) CheckSeaLabsPayAccount(ctx *gin.Context) {
	user, exists := ctx.Get("user")
	if !exists {
		_ = ctx.Error(apperror.BadRequestError("User is invalid"))
		return
	}

	value, _ := ctx.Get("payload")
	json, _ := value.(*dto.CheckSeaLabsPayReq)
	if user.(dto.UserJWT).UserID != json.UserID {
		_ = ctx.Error(apperror.ForbiddenError("Cannot check another user Sea Labs Pay Account"))
		return
	}

	result, err := h.seaLabsPayAccServ.CheckSeaLabsAccountExists(json)
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, dto.StatusOKResponse(result))
}

func (h *Handler) UpdateSeaLabsPayToMain(ctx *gin.Context) {
	user, exists := ctx.Get("user")
	if !exists {
		_ = ctx.Error(apperror.BadRequestError("User is invalid"))
		return
	}

	value, _ := ctx.Get("payload")
	json, _ := value.(*dto.UpdateSeaLabsPayToMainReq)
	if user.(dto.UserJWT).UserID != json.UserID {
		_ = ctx.Error(apperror.ForbiddenError("Cannot check another user Sea Labs Pay Account"))
		return
	}

	result, err := h.seaLabsPayAccServ.UpdateSeaLabsAccountToMain(json)
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, dto.StatusOKResponse(result))
}

func (h *Handler) GetSeaLabsPayAccount(ctx *gin.Context) {
	user, exists := ctx.Get("user")
	if !exists {
		_ = ctx.Error(apperror.BadRequestError("User is invalid"))
		return
	}

	result, err := h.seaLabsPayAccServ.GetSeaLabsAccountByUserID(user.(dto.UserJWT).UserID)
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, dto.StatusOKResponse(result))
}
