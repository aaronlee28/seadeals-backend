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
	sellerService          service.SellerService
	seaLabsPayAccServ      service.UserSeaPayAccountServ
	orderItemService       service.OrderItemService
	refreshTokenService    service.RefreshTokenService
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
	SellerService          service.SellerService
	WalletService          service.WalletService
	SeaLabsPayAccServ      service.UserSeaPayAccountServ
	OrderItemService       service.OrderItemService
	RefreshTokenService    service.RefreshTokenService
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
		sellerService:          config.SellerService,
		seaLabsPayAccServ:      config.SeaLabsPayAccServ,
		orderItemService:       config.OrderItemService,
		refreshTokenService:    config.RefreshTokenService,
	}
}
