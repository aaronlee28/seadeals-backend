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
	walletTransactionRepo := repository.NewWalletTransactionRepository()
	refreshTokenRepository := repository.NewRefreshTokenRepo()
	addressRepository := repository.NewAddressRepository()
	cityRepository := repository.NewCityRepository()
	districtRepository := repository.NewDistrictRepository()
	provinceRepository := repository.NewProvinceRepository()
	subDistrictRepository := repository.NewSubDistrictRepository()
	productCategoryRepository := repository.NewProductCategoryRepository()
	productRepository := repository.NewProductRepository()
	productVariantRepository := repository.NewProductVariantRepository()
	reviewRepository := repository.NewReviewRepository()
	sellerRepository := repository.NewSellerRepository()
	userSeaLabsPayAccountRepo := repository.NewSeaPayAccountRepo()
	orderItemRepository := repository.NewCartItemRepository()
	socialGraphRepo := repository.NewSocialGraphRepository()
	productVarDetRepo := repository.NewProductVariantDetailRepository()
	userAddressRepo := repository.NewUserAddressRepository()
	seaLabsPayTopUpHolderRepo := repository.NewSeaLabsPayTopUpHolderRepository()

	userService := service.NewUserService(&service.UserServiceConfig{
		DB:               db.Get(),
		UserRepository:   userRepository,
		UserRoleRepo:     userRoleRepository,
		WalletRepository: walletRepository,
		AppConfig:        config.Config,
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
		UserAddressRepo:   userAddressRepo,
	})

	productCategoryService := service.NewProductCategoryService(&service.ProductCategoryServiceConfig{
		DB:                        db.Get(),
		ProductCategoryRepository: productCategoryRepository,
	})

	productService := service.NewProductService(&service.ProductConfig{
		DB:                db.Get(),
		ProductRepo:       productRepository,
		ReviewRepo:        reviewRepository,
		ProductVarDetRepo: productVarDetRepo,
	})

	productVariantService := service.NewProductVariantService(&service.ProductVariantServiceConfig{
		DB:                 db.Get(),
		ProductVariantRepo: productVariantRepository,
		ProductRepo:        productRepository,
		ProductVarDetRepo:  productVarDetRepo,
	})

	reviewService := service.NewReviewService(&service.ReviewServiceConfig{
		DB:         db.Get(),
		ReviewRepo: reviewRepository,
	})

	sellerService := service.NewSellerService(&service.SellerServiceConfig{
		DB:              db.Get(),
		SellerRepo:      sellerRepository,
		ReviewRepo:      reviewRepository,
		SocialGraphRepo: socialGraphRepo,
	})

	walletService := service.NewWalletService(&service.WalletServiceConfig{
		DB:               db.Get(),
		WalletRepository: walletRepository,
		UserRepository:   userRepository,
	})

	userSeaLabsPayAccountServ := service.NewUserSeaPayAccountServ(&service.UserSeaPayAccountServConfig{
		DB:                        db.Get(),
		UserSeaPayAccountRepo:     userSeaLabsPayAccountRepo,
		SeaLabsPayTopUpHolderRepo: seaLabsPayTopUpHolderRepo,
		WalletRepository:          walletRepository,
		WalletTransactionRepo:     walletTransactionRepo,
	})

	orderItemService := service.NewCartItemService(&service.CartItemServiceConfig{
		DB:                 db.Get(),
		CartItemRepository: orderItemRepository,
	})

	refreshTokenService := service.NewRefreshTokenService(&service.RefreshTokenServiceConfig{
		DB:               db.Get(),
		RefreshTokenRepo: refreshTokenRepository,
	})
	sealabsPayService := service.NewSealabsPayService(&service.SealabsServiceConfig{
		DB: db.Get(),
	})
	router := NewRouter(&RouterConfig{
		UserService:            userService,
		AuthService:            authService,
		ProvinceService:        provinceService,
		CityService:            cityService,
		DistrictService:        districtService,
		SubDistrictService:     subDistrictService,
		AddressService:         addressService,
		WalletService:          walletService,
		ProductCategoryService: productCategoryService,
		ProductService:         productService,
		ProductVariantService:  productVariantService,
		ReviewService:          reviewService,
		SellerService:          sellerService,
		UserSeaLabsPayAccServ:  userSeaLabsPayAccountServ,
		OrderItemService:       orderItemService,
		RefreshTokenService:    refreshTokenService,
		SealabsPayService:      sealabsPayService,
	})
	log.Fatalln(router.Run(":" + config.Config.Port))
}
