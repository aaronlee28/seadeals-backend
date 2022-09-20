package handler

import "seadeals-backend/service"

type Handler struct {
	userService    service.UserService
	authService    service.AuthService
	productService service.ProductService
}

type Config struct {
	UserService    service.UserService
	AuthService    service.AuthService
	ProductService service.ProductService
}

func New(config *Config) *Handler {
	return &Handler{
		userService:    config.UserService,
		authService:    config.AuthService,
		productService: config.ProductService,
	}
}
