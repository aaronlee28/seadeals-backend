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
	productRepository := repository.NewProductRepository()

	userService := service.NewUserService(&service.UserServiceConfig{
		DB:               db.Get(),
		UserRepository:   userRepository,
		UserRoleRepo:     userRoleRepository,
		WalletRepository: walletRepository,
	})

	authService := service.NewAuthService(&service.AuthSConfig{
		DB:               db.Get(),
		RefreshTokenRepo: refreshTokenRepository,
		UserRepository:   userRepository,
		UserRoleRepo:     userRoleRepository,
		WalletRepository: walletRepository,
		AppConfig:        config.Config,
	})

	productService := service.NewProductService(&service.ProductConfig{
		DB:          db.Get(),
		ProductRepo: productRepository,
	})

	router := NewRouter(&RouterConfig{
		UserService:    userService,
		AuthService:    authService,
		ProductService: productService,
	})
	log.Fatalln(router.Run(":" + config.Config.Port))
}
