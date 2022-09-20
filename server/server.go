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
	addressRepository := repository.NewAddressRepository(&repository.AddressRepositoryConfig{})
	cityRepository := repository.NewCityRepository(&repository.CityRepositoryConfig{})
	districtRepository := repository.NewDistrictRepository(&repository.DistrictRepositoryConfig{})
	provinceRepository := repository.NewProvinceRepository(&repository.ProvinceRepositoryConfig{})
	subDistrictRepository := repository.NewSubDistrictRepository(&repository.SubDistrictRepositoryConfig{})

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
	})
	log.Fatalln(router.Run(":" + config.Config.Port))
}
