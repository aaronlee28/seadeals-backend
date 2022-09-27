package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"net/http"
	"seadeals-backend/dto"
)

func RequestValidator(creator func() any) gin.HandlerFunc {
	return func(context *gin.Context) {
		model := creator()
		if err := context.ShouldBindJSON(&model); err != nil {
			var e []map[string]any
			for _, valErr := range err.(validator.ValidationErrors) {
				tmp := map[string]any{"key": valErr.Namespace(), "tag": valErr.Tag()}
				e = append(e, tmp)
			}
			badRequest := &dto.AppResponse{
				StatusCode: http.StatusBadRequest,
				Status:     "BAD_REQUEST_ERROR",
				Data:       e,
			}
			context.AbortWithStatusJSON(badRequest.StatusCode, badRequest)
			return
		}

		context.Set("payload", model)
		context.Next()
	}
}
