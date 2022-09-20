package server

import (
	"github.com/gin-gonic/gin"
	"seadeals-backend/dto"
	"seadeals-backend/handler"
	"seadeals-backend/middleware"
	"seadeals-backend/service"
)

type RouterConfig struct {
	UserService        service.UserService
	AuthService        service.AuthService
	ProvinceService    service.ProvinceService
	CityService        service.CityService
	DistrictService    service.DistrictService
	SubDistrictService service.SubDistrictService
	AddressService     service.AddressService
	WalletService service.WalletService
}

func NewRouter(c *RouterConfig) *gin.Engine {
	r := gin.Default()
	r.NoRoute()

	h := handler.New(&handler.Config{
		UserService:        c.UserService,
		AuthService:        c.AuthService,
		ProvinceService:    c.ProvinceService,
		CityService:        c.CityService,
		DistrictService:    c.DistrictService,
		SubDistrictService: c.SubDistrictService,
		AddressService:     c.AddressService,
		WalletService: c.WalletService,
	})

	r.Use(middleware.ErrorHandler)
	r.Use(middleware.AllowCrossOrigin)

	// AUTH
	r.POST("/register", middleware.RequestValidator(func() any {
		return &dto.RegisterRequest{}
	}), h.Register)

	r.POST("/google/sign-in", middleware.RequestValidator(func() any {
		return &dto.GoogleLogin{}
	}), h.SignInWithGoogleEmail)

	// GUEST ROUTE
	r.GET("/provinces", h.GetProvinces)
	r.GET("/provinces/:id/cities", h.GetCitiesByProvinceID)
	r.GET("/cities/:id/districts", h.GetDistrictsByCityID)
	r.GET("/districts/:id/sub-districts", h.GetSubDistrictsByCityID)

	// USER ROUTE
	r.POST("/user/profiles/addresses", middleware.AuthorizeJWTFor("user"), middleware.RequestValidator(func() any {
		return &dto.CreateAddressReq{}
	}), h.CreateNewAddress)
	r.PATCH("/user/profiles/addresses", middleware.AuthorizeJWTFor("user"), middleware.RequestValidator(func() any {
		return &dto.UpdateAddressReq{}
	}), h.UpdateAddress)
	r.GET("/user/profiles/addresses", middleware.AuthorizeJWTFor("user"), h.GetAddressesByUserID)

	r.GET("/userwallet", middleware.AuthorizeJWTFor("user"), h.WalletDataTransactions)

	return r
}
