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
	UserService            service.UserService
	AuthService            service.AuthService
	ProvinceService        service.ProvinceService
	CityService            service.CityService
	DistrictService        service.DistrictService
	SubDistrictService     service.SubDistrictService
	AddressService         service.AddressService
	WalletService          service.WalletService
	ProductCategoryService service.ProductCategoryService
	ProductService         service.ProductService
	ProductVariantService  service.ProductVariantService
	ReviewService          service.ReviewService
	SellerService          service.SellerService
	UserSeaLabsPayAccServ  service.UserSeaPayAccountServ
	OrderItemService       service.OrderItemService
	RefreshTokenService    service.RefreshTokenService
}

func NewRouter(c *RouterConfig) *gin.Engine {
	r := gin.Default()
	r.NoRoute()

	h := handler.New(&handler.Config{
		UserService:            c.UserService,
		AuthService:            c.AuthService,
		ProvinceService:        c.ProvinceService,
		CityService:            c.CityService,
		DistrictService:        c.DistrictService,
		SubDistrictService:     c.SubDistrictService,
		AddressService:         c.AddressService,
		ProductCategoryService: c.ProductCategoryService,
		ProductService:         c.ProductService,
		ProductVariantService:  c.ProductVariantService,
		ReviewService:          c.ReviewService,
		SellerService:          c.SellerService,
		WalletService:          c.WalletService,
		SeaLabsPayAccServ:      c.UserSeaLabsPayAccServ,
		OrderItemService:       c.OrderItemService,
		RefreshTokenService:    c.RefreshTokenService,
	})

	r.Use(middleware.ErrorHandler)
	r.Use(middleware.AllowCrossOrigin)

	// AUTH
	r.POST("/register", middleware.RequestValidator(func() any {
		return &dto.RegisterRequest{}
	}), h.Register)
	r.GET("/refresh/access-token", h.RefreshAccessToken)
	r.POST("/sign-in", middleware.RequestValidator(func() any {
		return &dto.SignInReq{}
	}), h.SignIn)
	r.POST("/sign-out", middleware.AuthorizeJWTFor(model.UserRoleName), middleware.RequestValidator(func() any {
		return &dto.SignOutReq{}
	}), h.SignOut)

	// GOOGLE AUTH
	r.GET("/google/sign-in", h.GoogleSignIn)
	r.GET("/google/callback", h.GoogleCallback)

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

	// CATEGORIES
	r.GET("/categories", h.FindAllProductCategories)

	// PRODUCTS
	r.GET("/products/:id/variant", h.FindAllProductVariantByProductID)
	r.GET("/search-recommend-product/", h.SearchRecommendProduct)
	r.GET("/products/detail/:slug", h.FindProductDetailBySlug)
	r.GET("/sellers/:id/products", h.GetProductsBySellerID)

	// REVIEWS
	r.GET("/products/:id/reviews", h.FindReviewByProductID)

	// SELLER
	r.GET("/sellers/:id", h.FindSellerByID)

	// WALLET
	r.GET("/user-wallet", middleware.AuthorizeJWTFor(model.UserRoleName), h.WalletDataTransactions)
	r.GET("/transaction-details", middleware.RequestValidator(func() any { return &dto.TransactionDetailsReq{} }), middleware.AuthorizeJWTFor("user"), h.TransactionDetails)
	r.GET("/paginated-transaction", middleware.AuthorizeJWTFor(model.UserRoleName), h.PaginatedTransactions)
	r.PATCH("/wallet-pin", middleware.AuthorizeJWTFor(model.UserRoleName), middleware.RequestValidator(func() any { return &dto.PinReq{} }), h.WalletPin)
	r.POST("/user/validator/wallet-pin", middleware.AuthorizeJWTFor(model.UserRoleName), middleware.RequestValidator(func() any {
		return &dto.PinReq{}
	}), h.ValidateWalletPin)
	r.GET("/user/wallet/status", middleware.AuthorizeJWTFor(model.UserRoleName), h.GetWalletStatus)

	// SEA LABS ACCOUNT
	r.POST("/user/sea-labs-pay/register", middleware.AuthorizeJWTFor(model.UserRoleName), middleware.RequestValidator(func() any {
		return &dto.RegisterSeaLabsPayReq{}
	}), h.RegisterSeaLabsPayAccount)
	r.POST("/user/sea-labs-pay/validator", middleware.AuthorizeJWTFor(model.UserRoleName), middleware.RequestValidator(func() any {
		return &dto.CheckSeaLabsPayReq{}
	}), h.CheckSeaLabsPayAccount)
	r.PATCH("/user/sea-labs-pay", middleware.AuthorizeJWTFor(model.UserRoleName), middleware.RequestValidator(func() any {
		return &dto.UpdateSeaLabsPayToMainReq{}
	}), h.UpdateSeaLabsPayToMain)

	// ORDER ITEM
	r.GET("/user/cart", middleware.AuthorizeJWTFor(model.UserRoleName), h.GetOrderItem)
	r.POST("/user/cart", middleware.AuthorizeJWTFor(model.UserRoleName), middleware.RequestValidator(func() any {
		return &dto.AddToCartReq{}
	}), h.AddToCart)
	r.DELETE("/user/cart", middleware.AuthorizeJWTFor(model.UserRoleName), middleware.RequestValidator(func() any {
		return &dto.DeleteFromCartReq{}
	}), h.DeleteOrderItem)

	return r
}
