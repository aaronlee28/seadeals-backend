package server

import (
	"log"
	"seadeals-backend/config"
	"seadeals-backend/cronjob"
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
	productCategoryRepository := repository.NewProductCategoryRepository()
	productRepository := repository.NewProductRepository()
	productVariantRepository := repository.NewProductVariantRepository()
	reviewRepository := repository.NewReviewRepository()
	sellerRepository := repository.NewSellerRepository()
	userSeaLabsPayAccountRepo := repository.NewSeaPayAccountRepo()
	orderItemRepository := repository.NewCartItemRepository()
	socialGraphRepo := repository.NewSocialGraphRepository()
	productVarDetRepo := repository.NewProductVariantDetailRepository()
	seaLabsPayTopUpHolderRepo := repository.NewSeaLabsPayTopUpHolderRepository()
	seaLabsPayTransactionHolderRepo := repository.NewSeaLabsPayTransactionHolderRepository()
	favoriteRepository := repository.NewFavoriteRepository()
	voucherRepo := repository.NewVoucherRepository()
	promotionRepository := repository.NewPromotionRepository()
	courierRepository := repository.NewCourierRepository()
	orderRepository := repository.NewOrderRepo()
	sellerAvailableCourRepo := repository.NewSellerAvailableCourierRepository()
	transactionRepo := repository.NewTransactionRepository()
	adminRepository := repository.NewAdminRepository()
	complaintRepo := repository.NewComplaintRepository()
	complaintPhotoRepo := repository.NewComplaintPhotoRepository()
	notificationRepo := repository.NewNotificationRepository()
	deliveryRepository := repository.NewDeliveryRepository()
	deliveryActivityRepo := repository.NewDeliveryActivityRepository()
	accountHolderRepo := repository.NewAccountHolderRepository()

	userService := service.NewUserService(&service.UserServiceConfig{
		DB:               db.Get(),
		UserRepository:   userRepository,
		UserRoleRepo:     userRoleRepository,
		AddressRepo:      addressRepository,
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

	addressService := service.NewAddressService(&service.AddressServiceConfig{
		DB:                db.Get(),
		AddressRepository: addressRepository,
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
		SellerRepo:        sellerRepository,
		SocialGraphRepo:   socialGraphRepo,
		NotificationRepo:  notificationRepo,
	})

	productVariantService := service.NewProductVariantService(&service.ProductVariantServiceConfig{
		DB:                 db.Get(),
		ProductVariantRepo: productVariantRepository,
		ProductRepo:        productRepository,
		ProductVarDetRepo:  productVarDetRepo,
	})

	reviewService := service.NewReviewService(&service.ReviewServiceConfig{
		DB:          db.Get(),
		ReviewRepo:  reviewRepository,
		SellerRepo:  sellerRepository,
		ProductRepo: productRepository,
	})

	sellerService := service.NewSellerService(&service.SellerServiceConfig{
		DB:              db.Get(),
		SellerRepo:      sellerRepository,
		ReviewRepo:      reviewRepository,
		SocialGraphRepo: socialGraphRepo,
		ProductRepo:     productRepository,
	})

	walletService := service.NewWalletService(&service.WalletServiceConfig{
		DB:                db.Get(),
		AddressRepository: addressRepository,
		WalletRepository:  walletRepository,
		CourierRepository: courierRepository,
		DeliveryRepo:      deliveryRepository,
		DeliveryActRepo:   deliveryActivityRepo,
		UserRepository:    userRepository,
		WalletTransRepo:   walletTransactionRepo,
		UserRoleRepo:      userRoleRepository,
		SellerRepository:  sellerRepository,
		AccountHolderRepo: accountHolderRepo,
	})

	userSeaLabsPayAccountServ := service.NewUserSeaPayAccountServ(&service.UserSeaPayAccountServConfig{
		DB:                          db.Get(),
		AddressRepository:           addressRepository,
		UserSeaPayAccountRepo:       userSeaLabsPayAccountRepo,
		DeliveryRepo:                deliveryRepository,
		DeliveryActRepo:             deliveryActivityRepo,
		CourierRepository:           courierRepository,
		OrderRepo:                   orderRepository,
		SeaLabsPayTopUpHolderRepo:   seaLabsPayTopUpHolderRepo,
		SeaLabsPayTransactionHolder: seaLabsPayTransactionHolderRepo,
		SellerRepository:            sellerRepository,
		WalletRepository:            walletRepository,
		WalletTransactionRepo:       walletTransactionRepo,
		AccountHolderRepo:           accountHolderRepo,
	})

	orderItemService := service.NewCartItemService(&service.CartItemServiceConfig{
		DB:                 db.Get(),
		CartItemRepository: orderItemRepository,
		ProductVarDetRepo:  productVarDetRepo,
	})

	refreshTokenService := service.NewRefreshTokenService(&service.RefreshTokenServiceConfig{
		DB:               db.Get(),
		RefreshTokenRepo: refreshTokenRepository,
	})

	sealabsPayService := service.NewSealabsPayService(&service.SealabsServiceConfig{
		DB: db.Get(),
	})

	favoriteService := service.NewFavoriteService(&service.FavoriteServiceConfig{
		DB:                 db.Get(),
		FavoriteRepository: favoriteRepository,
		ProductRepository:  productRepository,
	})

	socialGraphService := service.NewSocialGraphService(&service.SocialGraphServiceConfig{
		DB:                    db.Get(),
		SocialGraphRepository: socialGraphRepo,
	})

	voucherService := service.NewVoucherService(&service.VoucherServiceConfig{
		DB:          db.Get(),
		VoucherRepo: voucherRepo,
		SellerRepo:  sellerRepository,
	})

	promotionService := service.NewPromotionService(&service.PromotionServiceConfig{
		DB:                  db.Get(),
		PromotionRepository: promotionRepository,
		SellerRepo:          sellerRepository,
		ProductRepo:         productRepository,
		SocialGraphRepo:     socialGraphRepo,
		NotificationRepo:    notificationRepo,
	})

	courierService := service.NewCourierService(&service.CourierServiceConfig{
		DB:                db.Get(),
		CourierRepository: courierRepository,
	})

	orderService := service.NewOrderService(&service.OrderServiceConfig{
		DB:                        db.Get(),
		OrderRepository:           orderRepository,
		AccountHolderRepo:         accountHolderRepo,
		AddressRepository:         addressRepository,
		CourierRepository:         courierRepository,
		SellerRepository:          sellerRepository,
		VoucherRepo:               voucherRepo,
		DeliveryRepo:              deliveryRepository,
		TransactionRepo:           transactionRepo,
		WalletRepository:          walletRepository,
		WalletTransRepo:           walletTransactionRepo,
		ProductVarDetRepo:         productVarDetRepo,
		ProductRepo:               productRepository,
		SeaLabsPayTransHolderRepo: seaLabsPayTransactionHolderRepo,
		ComplainRepo:              complaintRepo,
		ComplaintPhotoRepo:        complaintPhotoRepo,
		NotificationRepo:          notificationRepo,
	})

	deliveryService := service.NewDeliveryService(&service.DeliveryServiceConfig{
		DB:                     db.Get(),
		DeliveryRepository:     deliveryRepository,
		DeliverActivityRepo:    deliveryActivityRepo,
		AddressRepository:      addressRepository,
		OrderRepository:        orderRepository,
		SellerRepository:       sellerRepository,
		NotificationRepository: notificationRepo,
	})

	sellerAvailableCourServ := service.NewSellerAvailableCourService(&service.SellerAvailableCourServiceConfig{
		DB:                  db.Get(),
		SellerAvailCourRepo: sellerAvailableCourRepo,
		SellerRepository:    sellerRepository,
	})

	adminService := service.NewAdminRService(&service.AdminConfig{
		DB:        db.Get(),
		AdminRepo: adminRepository,
	})

	runCronJobHelper := cronjob.NewCronJob(&cronjob.RunCronJobConfig{
		DB:           db.Get(),
		OrderService: orderService,
	})

	runCronJobHelper.RunCronJobs()

	router := NewRouter(&RouterConfig{
		UserService:             userService,
		AuthService:             authService,
		AddressService:          addressService,
		WalletService:           walletService,
		ProductCategoryService:  productCategoryService,
		ProductService:          productService,
		ProductVariantService:   productVariantService,
		ReviewService:           reviewService,
		SellerService:           sellerService,
		UserSeaLabsPayAccServ:   userSeaLabsPayAccountServ,
		OrderItemService:        orderItemService,
		RefreshTokenService:     refreshTokenService,
		SealabsPayService:       sealabsPayService,
		FavoriteService:         favoriteService,
		SocialGraphService:      socialGraphService,
		VoucherService:          voucherService,
		PromotionService:        promotionService,
		CourierService:          courierService,
		OrderService:            orderService,
		DeliveryService:         deliveryService,
		SellerAvailableCourServ: sellerAvailableCourServ,
		AdminService:            adminService,
	})

	log.Fatalln(router.Run(":" + config.Config.Port))
}
