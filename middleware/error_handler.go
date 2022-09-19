package middleware

import (
	"github.com/gin-gonic/gin"
	"seadeals-backend/apperror"
)

func ErrorHandler(context *gin.Context) {
	context.Next()

	if len(context.Errors) <= 0 {
		return
	}

	firstError := context.Errors[0].Err
	appError, isAppError := firstError.(apperror.AppError)
	if isAppError {
		context.JSON(appError.StatusCode, appError)
		return
	}
	serverError := apperror.InternalServerError(firstError.Error())
	context.JSON(serverError.StatusCode, serverError)
}
