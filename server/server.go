package server

import (
	"log"
	"seadeals-backend/config"
	"seadeals-backend/db"
	"seadeals-backend/repository"
	"seadeals-backend/service"
)

func Init() {
	userRepository := repository.NewUserRepository()
	userRoleRepository := repository.NewUserRoleRepository()
	walletRepository := repository.NewWalletRepository()
	refreshTokenRepository := repository.NewRefreshTokenRepo()
	addressRepository := repository.NewAddressRepository()
	cityRepository := repository.NewCityRepository()
	districtRepository := repository.NewDistrictRepository()
	provinceRepository := repository.NewProvinceRepository()
	subDistrictRepository := repository.NewSubDistrictRepository()
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

	cityService := service.NewCityService(&service.CityServiceConfig{
		DB:             db.Get(),
		CityRepository: cityRepository,
	})

	districtService := service.NewDistrictService(&service.DistrictServiceConfig{
		DB:                 db.Get(),
		DistrictRepository: districtRepository,
	})

	provinceService := service.NewProvinceService(&service.ProvinceServiceConfig{
		DB:                 db.Get(),
		ProvinceRepository: provinceRepository,
	})

	subDistrictService := service.NewSubDistrictService(&service.SubDistrictServiceConfig{
		DB:                    db.Get(),
		SubDistrictRepository: subDistrictRepository,
	})

	addressService := service.NewAddressService(&service.AddressServiceConfig{
		DB:                db.Get(),
		AddressRepository: addressRepository,
	})

	productService := service.NewProductService(&service.ProductConfig{
		DB:          db.Get(),
		ProductRepo: productRepository,
	})

	walletService := service.NewWalletService(&service.WalletServiceConfig{
		DB:               db.Get(),
		WalletRepository: walletRepository,
	})
	router := NewRouter(&RouterConfig{
		UserService:        userService,
		AuthService:        authService,
		ProvinceService:    provinceService,
		CityService:        cityService,
		DistrictService:    districtService,
		SubDistrictService: subDistrictService,
		AddressService:     addressService,
		WalletService: walletService,
		ProductService:     productService,
	})
	log.Fatalln(router.Run(":" + config.Config.Port))
}
