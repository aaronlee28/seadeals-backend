package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"net/http"
	"reflect"
	"seadeals-backend/apperror"
)

func RequestValidator(creator func() any) gin.HandlerFunc {
	return func(context *gin.Context) {
		model := creator()
		if err := context.ShouldBindJSON(&model); err != nil {
			if reflect.TypeOf(err).String() == "validator.ValidationErrors" {
				e := err.(validator.ValidationErrors)[0]
				msg := fmt.Sprintf("%s failed on '%s' validation", e.Field(), e.Tag())
				badRequest := apperror.BadRequestError(msg)
				context.AbortWithStatusJSON(badRequest.StatusCode, badRequest)
				return
			}
			context.AbortWithStatusJSON(http.StatusBadRequest, apperror.BadRequestError(err.Error()))
			return
		}

		context.Set("payload", model)
		context.Next()
	}
}
