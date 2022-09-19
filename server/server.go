package server

import (
	"log"
	"seadeals-backend/config"
	"seadeals-backend/db"
	"seadeals-backend/repository"
	"seadeals-backend/service"
)

func Init() {
	userRepository := repository.NewUserRepository(&repository.UserRepositoryConfig{})
	userRoleRepository := repository.NewUserRoleRepository(&repository.UserRoleRepositoryConfig{})
	walletRepository := repository.NewWalletRepository(&repository.WalletRepositoryConfig{})
	refreshTokenRepository := repository.NewRefreshTokenRepo(&repository.RefreshTokenRepositoryConfig{})

	userService := service.NewUserService(&service.UserServiceConfig{
		DB:               db.Get(),
		UserRepository:   userRepository,
		UserRoleRepo:     userRoleRepository,
		WalletRepository: walletRepository,
	})

	authService := service.NewAuthService(&service.AuthSConfig{
		DB:               db.Get(),
		RefreshTokenRepo: refreshTokenRepository,
		AppConfig:        config.Config,
	})

	router := NewRouter(&RouterConfig{
		UserService: userService,
		AuthService: authService,
	})
	log.Fatalln(router.Run(":" + config.Config.Port))
}
