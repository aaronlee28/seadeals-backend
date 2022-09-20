package server

import (
	"github.com/gin-gonic/gin"
	"seadeals-backend/dto"
	"seadeals-backend/handler"
	"seadeals-backend/middleware"
	"seadeals-backend/model"
	"seadeals-backend/service"
)

type RouterConfig struct {
	UserService           service.UserService
	AuthService           service.AuthService
	ProvinceService       service.ProvinceService
	CityService           service.CityService
	DistrictService       service.DistrictService
	SubDistrictService    service.SubDistrictService
	AddressService        service.AddressService
	WalletService         service.WalletService
	ProductService        service.ProductService
	UserSeaLabsPayAccServ service.UserSeaPayAccountServ
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
		ProductService:     c.ProductService,
		WalletService:      c.WalletService,
		SeaLabsPayAccServ:  c.UserSeaLabsPayAccServ,
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

	// ADDRESS
	r.GET("/provinces", h.GetProvinces)
	r.GET("/provinces/:id/cities", h.GetCitiesByProvinceID)
	r.GET("/cities/:id/districts", h.GetDistrictsByCityID)
	r.GET("/districts/:id/sub-districts", h.GetSubDistrictsByCityID)
	r.POST("/user/profiles/addresses", middleware.AuthorizeJWTFor(model.UserRoleName), middleware.RequestValidator(func() any {
		return &dto.CreateAddressReq{}
	}), h.CreateNewAddress)
	r.PATCH("/user/profiles/addresses", middleware.AuthorizeJWTFor(model.UserRoleName), middleware.RequestValidator(func() any {
		return &dto.UpdateAddressReq{}
	}), h.UpdateAddress)
	r.GET("/user/profiles/addresses", middleware.AuthorizeJWTFor(model.UserRoleName), h.GetAddressesByUserID)

	// PRODUCTS
	r.GET("/products/:slug", h.FindProductDetailBySlug)

	// WALLET
	r.GET("/user-wallet", middleware.AuthorizeJWTFor("user"), h.WalletDataTransactions)

	// SEA LABS ACCOUNT
	r.POST("/user/sea-labs-pay/register", middleware.AuthorizeJWTFor("user"), middleware.RequestValidator(func() any {
		return &dto.RegisterSeaLabsPayReq{}
	}), h.RegisterSeaLabsPayAccount)
	r.POST("/user/sea-labs-pay/validator", middleware.AuthorizeJWTFor("user"), middleware.RequestValidator(func() any {
		return &dto.CheckSeaLabsPayReq{}
	}), h.CheckSeaLabsPayAccount)
	r.PATCH("/user/sea-labs-pay", middleware.AuthorizeJWTFor("user"), middleware.RequestValidator(func() any {
		return &dto.UpdateSeaLabsPayToMainReq{}
	}), h.UpdateSeaLabsPayToMain)

	return r
}
