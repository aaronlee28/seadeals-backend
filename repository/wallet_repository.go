package repository

import (
	"context"
	"fmt"
	"seadeals-backend/apperror"
	"seadeals-backend/dto"
	"seadeals-backend/model"
	"seadeals-backend/redisutils"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type WalletRepository interface {
	CreateWallet(*gorm.DB, *model.Wallet) (*model.Wallet, error)
	UpdateWallet(tx *gorm.DB, userID uint, newBalance float64) error
	GetWalletByUserID(*gorm.DB, uint) (*model.Wallet, error)
	GetTransactionsByUserID(tx *gorm.DB, userID uint) ([]*model.Transaction, error)
	TransactionDetails(tx *gorm.DB, transactionID uint) (*model.Transaction, error)
	PaginatedTransactions(tx *gorm.DB, q *Query, userID uint) (int, []*model.Transaction, error)
	TopUp(tx *gorm.DB, wallet *model.Wallet, amount float64) (*model.Wallet, error)

	WalletPin(tx *gorm.DB, userID uint, pin string) error
	RequestChangePinByEmail(userID uint, key string, code string) error
	ValidateRequestIsValid(userID uint, key string) error
	ValidateRequestByEmailCodeIsValid(userID uint, req *dto.CodeKeyRequestByEmailReq) error
	ChangeWalletPinByEmail(tx *gorm.DB, userID uint, sellerID uint, req *dto.ChangePinByEmailReq) (*model.Wallet, error)

	ValidateWalletPin(tx *gorm.DB, userID uint, pin string) error
	GetWalletStatus(tx *gorm.DB, userID uint) (string, error)
	StepUpPassword(tx *gorm.DB, userID uint, password string) error
	GetCartItem(tx *gorm.DB, cartID uint) (*model.CartItem, error)
	GetVoucher(tx *gorm.DB, voucherCode string) (*model.Voucher, error)
	CreateTransaction(tx *gorm.DB, transaction *model.Transaction) (*model.Transaction, error)
	CreateOrder(tx *gorm.DB, sellerID uint, voucherID *uint, transactionID uint, userID uint) (*model.Order, error)
	CreateOrderItemAndRemoveFromCart(tx *gorm.DB, productVariantDetailID uint, product *model.Product, orderID uint, userID uint, quantity uint, subtotal float64, cartItem *model.CartItem) error
	UpdateOrder(tx *gorm.DB, order *model.Order) error
	UpdateTransaction(tx *gorm.DB, transaction *model.Transaction) error
	UpdateStock(tx *gorm.DB, productVariantDetail *model.ProductVariantDetail, newStock uint) error
	CreateWalletTransaction(tx *gorm.DB, walletID uint, transaction *model.Transaction) error
	UpdateWalletBalance(tx *gorm.DB, wallet *model.Wallet, totalTransaction float64) error
}

type walletRepository struct{}

func NewWalletRepository() WalletRepository {
	return &walletRepository{}
}

type Query struct {
	Limit string
	Page  string
}

const (
	WalletBlocked string = "blocked"
	WalletActive  string = "active"
)

func (w *walletRepository) CreateWallet(tx *gorm.DB, wallet *model.Wallet) (*model.Wallet, error) {
	result := tx.Create(&wallet)
	if result.Error != nil {
		return nil, apperror.InternalServerError("cannot create new wallet")
	}

	return wallet, result.Error
}

func (w *walletRepository) GetWalletByUserID(tx *gorm.DB, userID uint) (*model.Wallet, error) {
	var wallet = &model.Wallet{}
	result := tx.Model(&wallet).Where("user_id = ?", userID).First(&wallet)
	if result.Error != nil {
		return nil, apperror.InternalServerError("cannot find wallet")
	}

	return wallet, nil
}

func (w *walletRepository) GetTransactionsByUserID(tx *gorm.DB, userID uint) ([]*model.Transaction, error) {
	var transactions []*model.Transaction
	result := tx.Where("user_id = ?", userID).Find(&transactions)
	if result.Error != nil {
		return nil, apperror.InternalServerError("cannot find transactions")
	}
	return transactions, nil
}

func (w *walletRepository) TransactionDetails(tx *gorm.DB, transactionID uint) (*model.Transaction, error) {
	var transaction = &model.Transaction{}
	var orders []*model.Order
	result := tx.Where("id = ?", transactionID).Preload("Voucher").First(&transaction)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, apperror.NotFoundError("No such transaction exists")
		}
		return nil, apperror.InternalServerError("cannot find transactions")
	}
	result2 := tx.Preload("Seller").Preload("OrderItems.ProductVariantDetail.ProductVariant1").Preload("OrderItems.ProductVariantDetail.ProductVariant2").Preload("OrderItems.ProductVariantDetail.Product").Preload("OrderItems.ProductVariantDetail.Product.Promotion").Where("transaction_id = ?", transactionID).Find(&orders)

	if result2.Error != nil {
		return nil, apperror.InternalServerError("cannot find order in the transaction")
	}
	for _, order := range orders {
		fmt.Println(order.Seller.Name)
	}
	transaction.Orders = orders
	return transaction, nil
}

func (w *walletRepository) PaginatedTransactions(tx *gorm.DB, q *Query, userID uint) (int, []*model.Transaction, error) {
	var trans []*model.Transaction
	limit, _ := strconv.Atoi(q.Limit)
	page, _ := strconv.Atoi(q.Page)
	offset := (limit * page) - limit

	result1 := tx.Where("user_id = ?", userID).Find(&trans)
	if result1.Error != nil {
		return 0, nil, apperror.InternalServerError("cannot find transactions")
	}

	result2 := tx.Limit(limit).Offset(offset).Order("created_at desc").Find(&trans)
	if result2.Error != nil {
		return 0, nil, apperror.InternalServerError("cannot find transactions")
	}
	totalLength := len(trans)
	return totalLength, trans, nil
}

func (w *walletRepository) TopUp(tx *gorm.DB, wallet *model.Wallet, amount float64) (*model.Wallet, error) {
	wallet.Balance += amount
	result := tx.Model(wallet).Clauses(clause.Returning{}).Updates(wallet)
	if result.Error != nil {
		return nil, apperror.InternalServerError("Cannot top up wallet")
	}
	return wallet, nil
}

func (w *walletRepository) WalletPin(tx *gorm.DB, userID uint, pin string) error {
	var wallet *model.Wallet
	result1 := tx.Model(&wallet).Where("user_id = ?", userID).First(&wallet)
	if result1.Error != nil {
		return apperror.InternalServerError("cannot find wallet")
	}

	result2 := tx.Model(&wallet).Update("pin", pin)
	if result2.Error != nil {
		return apperror.InternalServerError("failed to update pin")
	}
	return nil
}

func (w *walletRepository) RequestChangePinByEmail(userID uint, key string, code string) error {
	rds := redisutils.Use()
	ctx := context.Background()
	keyTries := "user:" + strconv.Itoa(int(userID)) + ":wallet:tries"

	tries, err := rds.Get(ctx, keyTries).Int()
	if err != nil && err != redis.Nil {
		return apperror.InternalServerError("Cannot get data in redis")
	}
	if tries >= 3 {
		return apperror.BadRequestError("Wallet is blocked because too many wrong attempts")
	}

	keyWallet := "user:" + strconv.FormatUint(uint64(userID), 10) + ":wallet:pin:request:key"
	codeWallet := "user:" + strconv.FormatUint(uint64(userID), 10) + ":wallet:pin:request:code"

	rds.Set(ctx, keyWallet, key, 5*time.Minute)
	rds.Set(ctx, codeWallet, code, 5*time.Minute)
	return nil
}

func (w *walletRepository) ValidateRequestIsValid(userID uint, key string) error {
	rds := redisutils.Use()
	ctx := context.Background()
	keyTries := "user:" + strconv.Itoa(int(userID)) + ":wallet:tries"

	tries, err := rds.Get(ctx, keyTries).Int()
	if err != nil && err != redis.Nil {
		return apperror.InternalServerError("Cannot get data in redis")
	}
	if tries >= 3 {
		return apperror.BadRequestError("Wallet is blocked because too many wrong attempts")
	}

	keyWallet := "user:" + strconv.FormatUint(uint64(userID), 10) + ":wallet:pin:request:key"

	keyRedis, err := rds.Get(ctx, keyWallet).Result()
	if err != nil && err != redis.Nil {
		return apperror.InternalServerError("Cannot get data in redis")
	}

	if key != keyRedis {
		return apperror.BadRequestError("Request is invalid or expired")
	}

	return nil
}

func (w *walletRepository) ValidateRequestByEmailCodeIsValid(userID uint, req *dto.CodeKeyRequestByEmailReq) error {
	rds := redisutils.Use()
	ctx := context.Background()
	keyTries := "user:" + strconv.Itoa(int(userID)) + ":wallet:tries"

	tries, err := rds.Get(ctx, keyTries).Int()
	if err != nil && err != redis.Nil {
		return apperror.InternalServerError("Cannot get data in redis")
	}
	if tries >= 3 {
		return apperror.BadRequestError("Wallet is blocked because too many wrong attempts")
	}

	keyWallet := "user:" + strconv.FormatUint(uint64(userID), 10) + ":wallet:pin:request:key"
	codeWallet := "user:" + strconv.FormatUint(uint64(userID), 10) + ":wallet:pin:request:code"

	keyRedis, err := rds.Get(ctx, keyWallet).Result()
	if err != nil && err != redis.Nil {
		return apperror.InternalServerError("Cannot get data in redis")
	}
	if req.Key != keyRedis {
		return apperror.BadRequestError("Request is invalid or expired")
	}

	codeRedis, err := rds.Get(ctx, codeWallet).Result()
	if err != nil && err != redis.Nil {
		return apperror.InternalServerError("Cannot get data in redis")
	}
	if req.Code != codeRedis {
		return apperror.BadRequestError("Code is invalid")
	}

	return nil
}

func (w *walletRepository) ChangeWalletPinByEmail(tx *gorm.DB, userID uint, walletID uint, req *dto.ChangePinByEmailReq) (*model.Wallet, error) {
	rds := redisutils.Use()
	ctx := context.Background()
	keyTries := "user:" + strconv.Itoa(int(userID)) + ":wallet:tries"

	tries, err := rds.Get(ctx, keyTries).Int()
	if err != nil && err != redis.Nil {
		return nil, apperror.InternalServerError("Cannot get data in redis")
	}
	if tries >= 3 {
		return nil, apperror.BadRequestError("Wallet is blocked because too many wrong attempts")
	}

	keyWallet := "user:" + strconv.FormatUint(uint64(userID), 10) + ":wallet:pin:request:key"
	codeWallet := "user:" + strconv.FormatUint(uint64(userID), 10) + ":wallet:pin:request:code"

	keyRedis, err := rds.Get(ctx, keyWallet).Result()
	if err != nil && err != redis.Nil {
		return nil, apperror.InternalServerError("Cannot get data in redis")
	}
	if req.Key != keyRedis {
		return nil, apperror.BadRequestError("Request is invalid or expired")
	}

	codeRedis, err := rds.Get(ctx, codeWallet).Result()
	if err != nil && err != redis.Nil {
		return nil, apperror.InternalServerError("Cannot get data in redis")
	}
	if req.Code != codeRedis {
		return nil, apperror.BadRequestError("Code is invalid")
	}

	wallet := &model.Wallet{ID: walletID}
	result := tx.Model(&wallet).Update("pin", req.Pin)
	if result.Error != nil {
		return nil, apperror.InternalServerError("failed to update pin")
	}

	rds.Del(ctx, keyWallet)
	rds.Del(ctx, codeWallet)
	return wallet, nil
}

func (w *walletRepository) ValidateWalletPin(tx *gorm.DB, userID uint, pin string) error {
	rds := redisutils.Use()
	ctx := context.Background()
	keyTries := "user:" + strconv.Itoa(int(userID)) + ":wallet:tries"

	tries, err := rds.Get(ctx, keyTries).Int()
	if err != nil && err != redis.Nil {
		return apperror.InternalServerError("Cannot get data in redis")
	}
	if tries >= 3 {
		return apperror.BadRequestError("Wallet is blocked because too many wrong attempts")
	}

	var wallet *model.Wallet
	result1 := tx.Model(&wallet).Where("user_id = ?", userID).First(&wallet)
	if result1.Error != nil {
		return apperror.InternalServerError("cannot find wallet")
	}
	if wallet.Pin == nil {
		return apperror.BadRequestError("Wallet does not have pin")
	}

	if *wallet.Pin != pin {
		tries += 1
		rds.Set(ctx, keyTries, tries, 15*time.Minute)
		if tries >= 3 {
			return apperror.BadRequestError("Too many wrong attempts, wallet is blocked for 15 minutes")
		}
		return apperror.BadRequestError("Pin is incorrect")
	}

	rds.Del(ctx, keyTries)
	return nil
}

func (w *walletRepository) GetWalletStatus(tx *gorm.DB, userID uint) (string, error) {
	rds := redisutils.Use()
	ctx := context.Background()
	keyTries := "user:" + strconv.Itoa(int(userID)) + ":wallet:tries"

	tries, err := rds.Get(ctx, keyTries).Int()
	if err != nil && err != redis.Nil {
		return "", apperror.InternalServerError("Cannot get data in redis")
	}
	if tries >= 3 {
		return WalletBlocked, nil
	}

	if err == redis.Nil {
		var wallet = &model.Wallet{}
		result1 := tx.Model(&wallet).Where("user_id = ?", userID).First(&wallet)
		if result1.Error != nil {
			return "", apperror.InternalServerError("cannot find wallet")
		}

		return wallet.Status, nil
	}

	return WalletActive, nil
}

func (w *walletRepository) StepUpPassword(tx *gorm.DB, userID uint, password string) error {
	var user = &model.User{}
	result1 := tx.Model(&user).Where("id = ?", userID).First(&user)
	if result1.Error != nil {
		return apperror.InternalServerError("cannot find wallet")
	}

	match := checkPasswordHash(password, user.Password)
	if !match {
		return apperror.BadRequestError("Invalid email or password")
	}

	return nil
}

func (w *walletRepository) UpdateWallet(tx *gorm.DB, userID uint, newBalance float64) error {
	var wallet = &model.Wallet{}

	result := tx.Model(&wallet).Where("user_id = ?", userID).First(&wallet)
	if result.Error != nil {
		return apperror.InternalServerError("cannot find wallet")
	}
	result2 := tx.Model(&wallet).Update("Balance", newBalance)
	if result2.Error != nil {
		return apperror.InternalServerError("Failed to update wallet")
	}
	return nil
}

func (w *walletRepository) GetCartItem(tx *gorm.DB, cartID uint) (*model.CartItem, error) {
	var cartItem = &model.CartItem{}

	result := tx.Preload("ProductVariantDetail.Product.Promotion", "start_date <= ? AND end_date >= ?", time.Now(), time.Now()).Preload("ProductVariantDetail.Product.ProductDetail").Where("id = ?", cartID).First(&cartItem)
	if result.Error != nil {
		return nil, apperror.InternalServerError("cannot find cart item")
	}

	return cartItem, nil
}

func (w *walletRepository) GetVoucher(tx *gorm.DB, voucherCode string) (*model.Voucher, error) {
	var voucher = &model.Voucher{}
	if voucherCode == "" {
		return nil, nil
	}
	result := tx.Where("code = ?", voucherCode).First(&voucher)
	if result.Error != nil {
		return nil, apperror.InternalServerError("cannot find voucher")
	}
	return voucher, nil
}

func (w *walletRepository) CreateOrderItemAndRemoveFromCart(tx *gorm.DB, productVariantDetailID uint, product *model.Product, orderID uint, userID uint, quantity uint, subtotal float64, cartItem *model.CartItem) error {
	createOrderItem := &model.OrderItem{
		ProductVariantDetailID: productVariantDetailID,
		OrderID:                &orderID,
		UserID:                 userID,
		Quantity:               quantity,
		Subtotal:               subtotal,
	}
	if product.Promotion != nil {
		createOrderItem.PromotionID = &product.Promotion.ID
	}
	result := tx.Create(&createOrderItem)
	if result.Error != nil {
		return apperror.InternalServerError("Failed to create order")
	}

	if cartItem != nil {
		result2 := tx.Model(&cartItem).Update("deleted_at", time.Now())
		if result2.Error != nil {
			return apperror.InternalServerError("Failed to delete cart item")
		}
	}

	return nil

}
func (w *walletRepository) CreateTransaction(tx *gorm.DB, transaction *model.Transaction) (*model.Transaction, error) {
	result := tx.Create(&transaction)
	if result.Error != nil {
		return nil, apperror.InternalServerError("Failed to create transaction")
	}
	return transaction, nil
}
func (w *walletRepository) CreateOrder(tx *gorm.DB, sellerID uint, voucherID *uint, transactionID uint, userID uint) (*model.Order, error) {
	order := &model.Order{
		SellerID:      sellerID,
		VoucherID:     nil,
		TransactionID: transactionID,
		UserID:        userID,
		Total:         0,
	}
	if voucherID != nil {
		order.VoucherID = voucherID
	}

	result := tx.Create(&order)
	if result.Error != nil {
		return nil, apperror.InternalServerError("Failed to create order")
	}

	return order, nil
}
func (w *walletRepository) UpdateOrder(tx *gorm.DB, order *model.Order) error {
	result := tx.Model(&order).Updates(&order)
	if result.Error != nil {
		return apperror.InternalServerError("failed to update order")
	}

	// GORM cannot update null value
	result = tx.Model(&order).Updates(map[string]interface{}{"voucher_id": order.VoucherID})
	if result.Error != nil {
		return apperror.InternalServerError("failed to update order")
	}
	return nil
}
func (w *walletRepository) UpdateTransaction(tx *gorm.DB, transaction *model.Transaction) error {
	result := tx.Model(&transaction).Updates(&transaction)
	if result.Error != nil {
		return apperror.InternalServerError("failed to update transaction")
	}

	// GORM cannot update null value
	result = tx.Model(&transaction).Updates(map[string]interface{}{"voucher_id": transaction.VoucherID})
	if result.Error != nil {
		return apperror.InternalServerError("failed to update transaction")
	}
	return nil
}

func (w *walletRepository) UpdateStock(tx *gorm.DB, productVariantDetail *model.ProductVariantDetail, newStock uint) error {
	result := tx.Model(&productVariantDetail).Update("stock", newStock)

	if result.Error != nil {
		return apperror.InternalServerError("failed to update stock")
	}
	return nil
}

func (w *walletRepository) CreateWalletTransaction(tx *gorm.DB, walletID uint, transaction *model.Transaction) error {

	walletTransaction := &model.WalletTransaction{
		WalletID:      walletID,
		TransactionID: &transaction.ID,
		Total:         transaction.Total,
		PaymentMethod: transaction.PaymentMethod,
		PaymentType:   "DEBIT",
		Description:   "Payment from wallet",
		CreatedAt:     transaction.CreatedAt,
	}

	result := tx.Create(&walletTransaction)
	if result.Error != nil {
		return apperror.InternalServerError("Failed to create wallet transaction")
	}
	return nil
}

func (w *walletRepository) UpdateWalletBalance(tx *gorm.DB, wallet *model.Wallet, totalTransaction float64) error {
	newBalance := wallet.Balance - totalTransaction
	result := tx.Model(&wallet).Update("balance", newBalance)

	if result.Error != nil {
		return apperror.InternalServerError("failed to update wallet's balance")
	}
	return nil
}
