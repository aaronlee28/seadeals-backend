package server

import (
	"github.com/gin-gonic/gin"
	"seadeals-backend/middleware"
)

type RouterConfig struct {
}

func NewRouter(c *RouterConfig) *gin.Engine {
	r := gin.Default()
	r.NoRoute()
	r.Use(middleware.ErrorHandler)
	r.Use(middleware.AllowCrossOrigin)
	return r
}
