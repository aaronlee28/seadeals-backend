package middleware

import (
	"github.com/gin-gonic/gin"
	"seadeals-backend/apperror"
)

func RequestValidator(creator func() any) gin.HandlerFunc {
	return func(context *gin.Context) {
		model := creator()
		if err := context.ShouldBindJSON(&model); err != nil {
			badRequest := apperror.BadRequestError("")
			context.AbortWithStatusJSON(badRequest.StatusCode, badRequest)
			return
		}

		context.Set("payload", model)
		context.Next()
	}
}
