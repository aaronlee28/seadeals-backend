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
	UserService             service.UserService
	AuthService             service.AuthService
	AddressService          service.AddressService
	WalletService           service.WalletService
	ProductCategoryService  service.ProductCategoryService
	ProductService          service.ProductService
	ProductVariantService   service.ProductVariantService
	ReviewService           service.ReviewService
	SellerService           service.SellerService
	UserSeaLabsPayAccServ   service.UserSeaPayAccountServ
	OrderItemService        service.CartItemService
	RefreshTokenService     service.RefreshTokenService
	SealabsPayService       service.SealabsPayService
	FavoriteService         service.FavoriteService
	SocialGraphService      service.SocialGraphService
	VoucherService          service.VoucherService
	PromotionService        service.PromotionService
	CourierService          service.CourierService
	OrderService            service.OrderService
	DeliveryService         service.DeliveryService
	SellerAvailableCourServ service.SellerAvailableCourService
	AdminService            service.AdminService
}

func NewRouter(c *RouterConfig) *gin.Engine {
	r := gin.Default()

	h := handler.New(&handler.Config{
		UserService:             c.UserService,
		AuthService:             c.AuthService,
		AddressService:          c.AddressService,
		ProductCategoryService:  c.ProductCategoryService,
		ProductService:          c.ProductService,
		ProductVariantService:   c.ProductVariantService,
		ReviewService:           c.ReviewService,
		SellerService:           c.SellerService,
		WalletService:           c.WalletService,
		SeaLabsPayAccServ:       c.UserSeaLabsPayAccServ,
		OrderItemService:        c.OrderItemService,
		RefreshTokenService:     c.RefreshTokenService,
		SealabsPayService:       c.SealabsPayService,
		FavoriteService:         c.FavoriteService,
		SocialGraphService:      c.SocialGraphService,
		VoucherService:          c.VoucherService,
		PromotionService:        c.PromotionService,
		CourierService:          c.CourierService,
		OrderService:            c.OrderService,
		DeliveryService:         c.DeliveryService,
		SellerAvailableCourServ: c.SellerAvailableCourServ,
		AdminService:            c.AdminService,
	})
	r.Use(middleware.ErrorHandler)
	r.Use(middleware.AllowCrossOrigin)
	r.NoRoute()

	// PING
	r.GET("/", h.Ping)

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

	// USER
	r.GET("/user/profiles", middleware.AuthorizeJWTFor(model.UserRoleName), h.UserDetails)
	r.PATCH("/user/change-profiles", middleware.AuthorizeJWTFor(model.UserRoleName), middleware.RequestValidator(func() any {
		return &dto.ChangeUserDetails{}
	}), h.ChangeUserDetails)
	r.PATCH("/user/change-password", middleware.AuthorizeJWTFor(model.UserRoleName), middleware.RequestValidator(func() any {
		return &dto.ChangePasswordReq{}
	}), h.ChangeUserPassword)
	// ADDRESS
	r.POST("/user/profiles/addresses", middleware.AuthorizeJWTFor(model.UserRoleName), middleware.RequestValidator(func() any {
		return &dto.CreateAddressReq{}
	}), h.CreateNewAddress)
	r.PATCH("/user/profiles/addresses", middleware.AuthorizeJWTFor(model.UserRoleName), middleware.RequestValidator(func() any {
		return &dto.UpdateAddressReq{}
	}), h.UpdateAddress)
	r.PATCH("/user/profiles/addresses/:id", middleware.AuthorizeJWTFor(model.UserRoleName), h.ChangeMainAddress)
	r.GET("/user/profiles/addresses/main", middleware.AuthorizeJWTFor(model.UserRoleName), h.GetUserMainAddress)
	r.GET("/user/profiles/addresses", middleware.AuthorizeJWTFor(model.UserRoleName), h.GetAddressesByUserID)

	// COURIER
	r.GET("/couriers", middleware.AuthorizeJWTFor(model.SellerRoleName), h.GetAllCouriers)
	r.GET("/sellers/:id/couriers", middleware.AuthorizeJWTFor(model.UserRoleName), h.GetAvailableCourierForBuyer)

	// CATEGORIES
	r.GET("/categories", h.FindCategories)

	// PRODUCTS
	r.GET("/products/:id/variant", h.FindAllProductVariantByProductID)
	r.GET("/products/:id/promotion-price", h.GetVariantPriceAfterPromotionByProductID)
	r.GET("/products/:id/similar-products", h.FindSimilarProduct)
	r.GET("/search-recommend-product", h.SearchRecommendProduct)
	r.GET("/products/detail/:id", middleware.OptionalAuthorizeJWTFor(model.UserRoleName), h.FindProductDetailByID)
	r.GET("/sellers/:id/products", h.GetProductsBySellerID)
	r.GET("/categories/:id/products", h.GetProductsByCategoryID)
	r.GET("/products", h.SearchProducts)
	r.POST("/sellers/create-product", middleware.AuthorizeJWTFor(model.UserRoleName), middleware.RequestValidator(func() any {
		return &dto.PostCreateProductReq{}
	}), h.CreateSellerProduct)
	r.PATCH("/sellers/:id/update-product-and-details", middleware.AuthorizeJWTFor(model.UserRoleName), middleware.RequestValidator(func() any {
		return &dto.PatchProductAndDetailsReq{}
	}), h.UpdateProductAndDetails)
	r.PATCH("/sellers/:id/update-variant-and-details", middleware.AuthorizeJWTFor(model.UserRoleName), middleware.RequestValidator(func() any {
		return &dto.PatchVariantAndDetails{}
	}), h.UpdateVariantAndDetails)
	r.DELETE("/sellers/:id/delete-variant-and-details", middleware.AuthorizeJWTFor(model.UserRoleName), middleware.RequestValidator(func() any {
		return &dto.DefaultPrice{}
	}), h.DeleteVariantAndDetails)
	r.POST("/sellers/:id/add-variant-and-details", middleware.AuthorizeJWTFor(model.UserRoleName), middleware.RequestValidator(func() any {
		return &dto.AddVariantAndDetails{}
	}), h.AddVariantDetails)
	r.POST("/sellers/:id/add-product-photo", middleware.AuthorizeJWTFor(model.UserRoleName), middleware.RequestValidator(func() any {
		return &dto.ProductPhotoReq{}
	}), h.AddProductPhoto)
	r.DELETE("/sellers/:id/delete-product-photo", middleware.AuthorizeJWTFor(model.UserRoleName), middleware.RequestValidator(func() any {
		return &dto.DeleteProductPhoto{}
	}), h.DeleteProductPhoto)
	r.DELETE("/sellers/:id/delete-product", middleware.AuthorizeJWTFor(model.UserRoleName), h.DeleteProduct)

	// NOTIFICATION
	r.POST("/products/favorites", middleware.AuthorizeJWTFor(model.UserRoleName), middleware.RequestValidator(func() any {
		return &dto.FavoriteProductReq{}
	}), h.FavoriteToProduct)
	r.POST("/sellers/follow", middleware.AuthorizeJWTFor(model.UserRoleName), middleware.RequestValidator(func() any {
		return &dto.FollowSellerReq{}
	}), h.FollowToSeller)
	r.GET("/sellers/products", middleware.AuthorizeJWTFor(model.SellerRoleName), h.GetProductsByUserIDUnscoped)

	// SELLER
	r.GET("/sellers/:id", middleware.OptionalAuthorizeJWTFor(model.UserRoleName), h.FindSellerByID)
	r.POST("/sellers", middleware.AuthorizeJWTFor(model.UserRoleName), middleware.RequestValidator(func() any {
		return &dto.RegisterAsSellerReq{}
	}), h.RegisterAsSeller)
	r.POST("/sellers/couriers", middleware.AuthorizeJWTFor(model.SellerRoleName), middleware.RequestValidator(func() any {
		return &dto.AddDeliveryReq{}
	}), h.CreateOrUpdateSellerAvailableCour)
	r.GET("/sellers/couriers", middleware.AuthorizeJWTFor(model.SellerRoleName), h.GetAvailableCourierForSeller)
	r.GET("/sellers/:id/vouchers", middleware.AuthorizeJWTFor(model.UserRoleName), h.GetVouchersBySellerID)

	// ORDER
	r.GET("/sellers/orders", middleware.AuthorizeJWTFor(model.SellerRoleName), h.GetSellerOrders)
	r.GET("/user/order/:id", middleware.AuthorizeJWTFor(model.UserRoleName), h.GetOrderByID)
	r.GET("/user/orders", middleware.AuthorizeJWTFor(model.UserRoleName), h.GetBuyerOrders)
	r.POST("/user/finish/orders", middleware.AuthorizeJWTFor(model.UserRoleName), middleware.RequestValidator(func() any {
		return &dto.FinishOrderReq{}
	}), h.FinishOrder)

	// RECEIPT
	r.GET("/user/orders/receipt/:id", middleware.AuthorizeJWTFor(model.UserRoleName), h.GetDetailOrderForReceipt)

	// DELIVERY
	r.POST("/seller/deliver/order", middleware.AuthorizeJWTFor(model.SellerRoleName), middleware.RequestValidator(func() any {
		return &dto.DeliverOrderReq{}
	}), h.DeliverOrder)
	r.GET("/seller/settings/print", middleware.AuthorizeJWTFor(model.SellerRoleName), h.GetSellerPrintSettings)
	r.PATCH("/seller/settings/print", middleware.AuthorizeJWTFor(model.SellerRoleName), middleware.RequestValidator(func() any {
		return &dto.DeliverSettingsPrint{}
	}), h.UpdatePrintSettings)

	// THERMAL
	r.GET("/seller/orders/thermal/:id", middleware.AuthorizeJWTFor(model.SellerRoleName), h.GetDetailOrderForThermal)

	// VOUCHER
	r.POST("/vouchers", middleware.AuthorizeJWTFor(model.SellerRoleName), middleware.RequestValidator(func() any {
		return &dto.PostVoucherReq{}
	}), h.CreateVoucher)
	r.POST("/validate-voucher", middleware.AuthorizeJWTFor(model.UserRoleName), middleware.RequestValidator(func() any {
		return &dto.PostValidateVoucherReq{}
	}), h.ValidateVoucher)
	r.GET("/vouchers", middleware.AuthorizeJWTFor(model.SellerRoleName), h.FindVoucherByUserID)
	r.GET("/vouchers/:id/detail", middleware.AuthorizeJWTFor(model.SellerRoleName), h.FindVoucherDetailByID)
	r.GET("/vouchers/:id", h.FindVoucherByID)
	r.PATCH("/vouchers/:id", middleware.AuthorizeJWTFor(model.SellerRoleName), middleware.RequestValidator(func() any {
		return &dto.PatchVoucherReq{}
	}), h.UpdateVoucher)
	r.DELETE("/vouchers/:id", middleware.AuthorizeJWTFor(model.SellerRoleName), h.DeleteVoucherByID)
	r.GET("/global-vouchers", middleware.AuthorizeJWTFor(model.UserRoleName), h.GetAvailableGlobalVouchers)

	// REFUND
	r.POST("/cancel/orders", middleware.AuthorizeJWTFor(model.SellerRoleName), middleware.RequestValidator(func() any {
		return &dto.SellerCancelOrderReq{}
	}), h.CancelOrderBySeller)
	r.POST("/request-refund/orders", middleware.AuthorizeJWTFor(model.UserRoleName), middleware.RequestValidator(func() any {
		return &dto.CreateComplaintReq{}
	}), h.RequestRefundByBuyer)
	r.POST("/reject-refund/orders", middleware.AuthorizeJWTFor(model.SellerRoleName), middleware.RequestValidator(func() any {
		return &dto.RejectAcceptRefundReq{}
	}), h.RejectRefundRequest)
	r.POST("/accept-refund/orders", middleware.AuthorizeJWTFor(model.SellerRoleName), middleware.RequestValidator(func() any {
		return &dto.RejectAcceptRefundReq{}
	}), h.AcceptRefundRequest)

	// PROMOTION
	r.GET("/promotions", middleware.AuthorizeJWTFor(model.UserRoleName), h.GetPromotion)
	r.POST("/promotions", middleware.RequestValidator(func() any { return &dto.CreatePromotionArrayReq{} }), middleware.AuthorizeJWTFor(model.UserRoleName), h.CreatePromotion)
	r.GET("/view-detail-promotion/:id", middleware.AuthorizeJWTFor(model.UserRoleName), h.ViewDetailPromotionByID)
	r.PATCH("/patch-promotions", middleware.AuthorizeJWTFor(model.UserRoleName), middleware.RequestValidator(func() any {
		return &dto.PatchPromotionArrayReq{}
	}), h.UpdatePromotion)

	// WALLET
	r.GET("/user-wallet", middleware.AuthorizeJWTFor(model.UserRoleName), h.WalletDataTransactions)
	r.GET("/transactions/:id", middleware.AuthorizeJWTFor(model.UserRoleName), h.TransactionDetails)
	r.GET("/paginated-transaction", middleware.AuthorizeJWTFor(model.UserRoleName), h.PaginatedTransactions)
	r.GET("/user/wallet/transactions", middleware.AuthorizeJWTFor(model.UserRoleName), h.GetWalletTransactions)
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

	// TOP UP WITH SEA LABS
	r.POST("/user/wallet/top-up/sea-labs-pay", middleware.AuthorizeJWTFor(model.UserRoleName), middleware.RequestValidator(func() any {
		return &dto.TopUpWalletWithSeaLabsPayReq{}
	}), h.TopUpWithSeaLabsPay)
	r.POST("/user/wallet/top-up/sea-labs-pay/callback", middleware.RequestValidator(func() any {
		return &dto.SeaLabsPayReq{}
	}), h.TopUpWithSeaLabsPayCallback)

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

	// PAY WITH SEA LABS
	r.POST("/order/pay/sea-labs-pay", middleware.AuthorizeJWTFor(model.UserRoleName), middleware.RequestValidator(func() any {
		return &dto.CheckoutCartReq{}
	}), h.PayWithSeaLabsPay)
	// PAY WITH SEA LABS
	r.POST("/order/pay/sea-labs-pay/callback", middleware.RequestValidator(func() any {
		return &dto.SeaLabsPayReq{}
	}), h.PayWithSeaLabsPayCallback)

	// CART ITEM
	r.GET("/user/cart", middleware.AuthorizeJWTFor(model.UserRoleName), h.GetCartItem)
	r.POST("/user/cart", middleware.AuthorizeJWTFor(model.UserRoleName), middleware.RequestValidator(func() any {
		return &dto.AddToCartReq{}
	}), h.AddToCart)
	r.PATCH("/user/cart", middleware.AuthorizeJWTFor(model.UserRoleName), middleware.RequestValidator(func() any {
		return &dto.UpdateCartItemReq{}
	}), h.UpdateCart)
	r.DELETE("/user/cart", middleware.AuthorizeJWTFor(model.UserRoleName), middleware.RequestValidator(func() any {
		return &dto.DeleteFromCartReq{}
	}), h.DeleteCartItem)

	// FAVORITES
	r.GET("/user/favorite-counts", middleware.AuthorizeJWTFor(model.UserRoleName), h.GetUserFavoriteCount)

	// PAYMENT
	r.POST("/checkout-cart", middleware.RequestValidator(func() any { return &dto.CheckoutCartReq{} }), middleware.AuthorizeJWTFor(model.Level1RoleName), h.CheckoutCart)
	r.POST("/predicted-price", middleware.AuthorizeJWTFor(model.UserRoleName), middleware.RequestValidator(func() any {
		return &dto.PredictedPriceReq{}
	}), h.GetTotalPredictedPrice)

	// ADMIN
	r.POST("/create-global-voucher", middleware.RequestValidator(func() any { return &dto.CreateGlobalVoucher{} }), middleware.AuthorizeJWTFor(model.AdminRoleName), h.CreateGlobalVoucher)
	r.POST("/create-category", middleware.RequestValidator(func() any { return &dto.CreateCategory{} }), middleware.AuthorizeJWTFor(model.AdminRoleName), h.CreateCategory)

	// REVIEWS
	r.GET("/products/:id/reviews", h.FindReviewByProductID)
	r.POST("/product/review", middleware.RequestValidator(func() any { return &dto.CreateUpdateReview{} }), middleware.AuthorizeJWTFor(model.UserRoleName), h.CreateUpdateReview)
	r.GET("/user/review-history", middleware.AuthorizeJWTFor(model.UserRoleName), h.UserReviewHistory)
	r.POST("/user/existing-review", middleware.AuthorizeJWTFor(model.UserRoleName), middleware.RequestValidator(func() any { return &dto.GetExistingReview{} }), h.FindReviewByProductIDAndSellerID)

	return r

}
