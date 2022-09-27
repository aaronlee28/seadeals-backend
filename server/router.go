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
	AddressService         service.AddressService
	WalletService          service.WalletService
	ProductCategoryService service.ProductCategoryService
	ProductService         service.ProductService
	ProductVariantService  service.ProductVariantService
	ReviewService          service.ReviewService
	SellerService          service.SellerService
	UserSeaLabsPayAccServ  service.UserSeaPayAccountServ
	OrderItemService       service.CartItemService
	RefreshTokenService    service.RefreshTokenService
	SealabsPayService      service.SealabsPayService
}

func NewRouter(c *RouterConfig) *gin.Engine {
	r := gin.Default()

	h := handler.New(&handler.Config{
		UserService:            c.UserService,
		AuthService:            c.AuthService,
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
		SealabsPayService:      c.SealabsPayService,
	})

	r.Use(middleware.ErrorHandler)
	r.Use(middleware.AllowCrossOrigin)
	r.NoRoute()

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
	r.POST("/step-up-password", middleware.AuthorizeJWTFor(model.UserRoleName), middleware.RequestValidator(func() any {
		return &dto.StepUpPasswordRes{}
	}), h.StepUpPassword)
	// GOOGLE AUTH
	r.POST("/google/sign-in", middleware.RequestValidator(func() any {
		return &dto.GoogleLogin{}
	}), h.GoogleSignIn)

	// ADDRESS
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
	r.GET("/products/:id/promotion-price", h.GetVariantPriceAfterPromotionByProductID)
	r.GET("/products/:id/similar-products", h.FindSimilarProduct)
	r.GET("/search-recommend-product", h.SearchRecommendProduct)
	r.GET("/products/detail/:slug", h.FindProductDetailBySlug)
	r.GET("/sellers/:id/products", h.GetProductsBySellerID)
	r.GET("/categories/:id/products", h.GetProductsByCategoryID)
	r.GET("/products", h.SearchProducts)

	// REVIEWS
	r.GET("/products/:id/reviews", h.FindReviewByProductID)

	// SELLER
	r.GET("/sellers/:id", h.FindSellerByID)
	r.POST("/sellers", middleware.AuthorizeJWTFor(model.UserRoleName), middleware.RequestValidator(func() any {
		return &dto.RegisterAsSellerReq{}
	}), h.RegisterAsSeller)

	// WALLET
	r.GET("/user-wallet", middleware.AuthorizeJWTFor(model.UserRoleName), h.WalletDataTransactions)
	r.GET("/transaction-details", middleware.RequestValidator(func() any { return &dto.TransactionDetailsReq{} }), middleware.AuthorizeJWTFor("user"), h.TransactionDetails)
	r.GET("/paginated-transaction", middleware.AuthorizeJWTFor(model.UserRoleName), h.PaginatedTransactions)
	r.PATCH("/wallet-pin", middleware.AuthorizeJWTFor(model.UserRoleName), middleware.RequestValidator(func() any { return &dto.PinReq{} }), h.WalletPin)
	r.POST("/wallet/pin-by-email/", middleware.AuthorizeJWTFor(model.UserRoleName), h.RequestWalletChangeByEmail)
	r.POST("/wallet/validator/pin-by-email", middleware.AuthorizeJWTFor(model.UserRoleName), middleware.RequestValidator(func() any {
		return &dto.KeyRequestByEmailReq{}
	}), h.ValidateIfRequestByEmailIsValid)
	r.POST("/wallet/validator/pin-by-email/code", middleware.AuthorizeJWTFor(model.UserRoleName), middleware.RequestValidator(func() any {
		return &dto.CodeKeyRequestByEmailReq{}
	}), h.ValidateIfRequestChangeByEmailCodeIsValid)
	r.PATCH("/wallet/pin-by-email", middleware.AuthorizeJWTFor(model.UserRoleName), middleware.RequestValidator(func() any {
		return &dto.ChangePinByEmailReq{}
	}), h.ChangeWalletPinByEmail)
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
	r.POST("create-signature", middleware.RequestValidator(func() any { return &dto.SeaDealspayReq{} }), h.CreateSignature)
	r.GET("/user/sea-labs-pay", middleware.AuthorizeJWTFor(model.UserRoleName), h.GetSeaLabsPayAccount)

	// CART ITEM
	r.GET("/user/cart", middleware.AuthorizeJWTFor(model.UserRoleName), h.GetCartItem)
	r.POST("/user/cart", middleware.AuthorizeJWTFor(model.UserRoleName), middleware.RequestValidator(func() any {
		return &dto.AddToCartReq{}
	}), h.AddToCart)
	r.DELETE("/user/cart", middleware.AuthorizeJWTFor(model.UserRoleName), middleware.RequestValidator(func() any {
		return &dto.DeleteFromCartReq{}
	}), h.DeleteCartItem)

	return r
}
