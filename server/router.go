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
	UserService        service.UserService
	AuthService        service.AuthService
	ProvinceService    service.ProvinceService
	CityService        service.CityService
	DistrictService    service.DistrictService
	SubDistrictService service.SubDistrictService
	AddressService     service.AddressService
	WalletService      service.WalletService
	ProductService     service.ProductService
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

	r.GET("/transaction-details", middleware.RequestValidator(func() any { return &dto.TransactionDetailsReq{} }), middleware.AuthorizeJWTFor("user"), h.TransactionDetails)

	r.GET("/paginated-transaction", middleware.AuthorizeJWTFor("user"), h.TransactionDetails)
	return r
}
