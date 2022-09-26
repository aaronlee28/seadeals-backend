package handler

import "seadeals-backend/service"

type Handler struct {
	userService            service.UserService
	authService            service.AuthService
	provinceService        service.ProvinceService
	cityService            service.CityService
	districtService        service.DistrictService
	subDistrictService     service.SubDistrictService
	addressService         service.AddressService
	walletService          service.WalletService
	productCategoryService service.ProductCategoryService
	productService         service.ProductService
	productVariantService  service.ProductVariantService
	reviewService          service.ReviewService
	sellerService          service.SellerService
	seaLabsPayAccServ      service.UserSeaPayAccountServ
	orderItemService       service.CartItemService
	refreshTokenService    service.RefreshTokenService
	sealabsPayService      service.SealabsPayService
}

type Config struct {
	UserService            service.UserService
	AuthService            service.AuthService
	ProvinceService        service.ProvinceService
	CityService            service.CityService
	DistrictService        service.DistrictService
	SubDistrictService     service.SubDistrictService
	AddressService         service.AddressService
	ProductCategoryService service.ProductCategoryService
	ProductService         service.ProductService
	ProductVariantService  service.ProductVariantService
	ReviewService          service.ReviewService
	SellerService          service.SellerService
	WalletService          service.WalletService
	SeaLabsPayAccServ      service.UserSeaPayAccountServ
	OrderItemService       service.CartItemService
	RefreshTokenService    service.RefreshTokenService
	SealabsPayService      service.SealabsPayService
}

func New(config *Config) *Handler {
	return &Handler{
		userService:            config.UserService,
		authService:            config.AuthService,
		walletService:          config.WalletService,
		cityService:            config.CityService,
		provinceService:        config.ProvinceService,
		districtService:        config.DistrictService,
		subDistrictService:     config.SubDistrictService,
		addressService:         config.AddressService,
		productCategoryService: config.ProductCategoryService,
		productService:         config.ProductService,
		productVariantService:  config.ProductVariantService,
		reviewService:          config.ReviewService,
		sellerService:          config.SellerService,
		seaLabsPayAccServ:      config.SeaLabsPayAccServ,
		orderItemService:       config.OrderItemService,
		refreshTokenService:    config.RefreshTokenService,
		sealabsPayService:      config.SealabsPayService,
	}
}
