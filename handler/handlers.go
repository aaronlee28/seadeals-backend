package handler

import "seadeals-backend/service"

type Handler struct {
	userService service.UserService
	authService service.AuthService
}

type Config struct {
	UserService service.UserService
	AuthService service.AuthService
}

func New(config *Config) *Handler {
	return &Handler{
		userService: config.UserService,
		authService: config.AuthService,
	}
}
