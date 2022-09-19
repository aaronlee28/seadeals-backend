package server

import (
	"github.com/gin-gonic/gin"
	"seadeals-backend/dto"
	"seadeals-backend/handler"
	"seadeals-backend/middleware"
	"seadeals-backend/service"
)

type RouterConfig struct {
	UserService service.UserService
	AuthService service.AuthService
}

func NewRouter(c *RouterConfig) *gin.Engine {
	r := gin.Default()
	r.NoRoute()

	h := handler.New(&handler.Config{
		UserService: c.UserService,
		AuthService: c.AuthService,
	})

	r.Use(middleware.ErrorHandler)
	r.Use(middleware.AllowCrossOrigin)

	r.POST("/register", middleware.RequestValidator(func() any {
		return &dto.RegisterRequest{}
	}), h.Register)

	return r
}
