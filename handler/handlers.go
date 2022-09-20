package handler

import "seadeals-backend/service"

type Handler struct {
	userService        service.UserService
	authService        service.AuthService
	provinceService    service.ProvinceService
	cityService        service.CityService
	districtService    service.DistrictService
	subDistrictService service.SubDistrictService
	addressService     service.AddressService
	productService service.ProductService
}

type Config struct {
	UserService        service.UserService
	AuthService        service.AuthService
	ProvinceService    service.ProvinceService
	CityService        service.CityService
	DistrictService    service.DistrictService
	SubDistrictService service.SubDistrictService
	AddressService     service.AddressService
	ProductService service.ProductService
}

func New(config *Config) *Handler {
	return &Handler{
		userService:        config.UserService,
		authService:        config.AuthService,
		cityService:        config.CityService,
		provinceService:    config.ProvinceService,
		districtService:    config.DistrictService,
		subDistrictService: config.SubDistrictService,
		addressService:     config.AddressService,
		productService: config.ProductService,
	}
}
