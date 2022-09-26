package repository

import (
	"context"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
	"seadeals-backend/apperror"
	"seadeals-backend/dto"
	"seadeals-backend/model"
	"seadeals-backend/redisutils"
	"strconv"
	"time"
)

type WalletRepository interface {
	CreateWallet(*gorm.DB, *model.Wallet) (*model.Wallet, error)
	GetWalletByUserID(*gorm.DB, uint) (*model.Wallet, error)
	GetTransactionsByUserID(tx *gorm.DB, userID uint) (*[]model.Transaction, error)
	TransactionDetails(tx *gorm.DB, transactionID uint) (*model.Transaction, error)
	PaginatedTransactions(tx *gorm.DB, q *Query, userID uint) (int, *[]model.Transaction, error)
	WalletPin(tx *gorm.DB, userID uint, pin string) error

	RequestChangePinByEmail(userID uint, key string, code string) error
	ValidateRequestIsValid(userID uint, key string) error
	ValidateRequestByEmailCodeIsValid(userID uint, req *dto.CodeKeyRequestByEmailReq) error
	ChangeWalletPinByEmail(tx *gorm.DB, userID uint, sellerID uint, req *dto.ChangePinByEmailReq) (*model.Wallet, error)

	ValidateWalletPin(tx *gorm.DB, userID uint, pin string) error
	GetWalletStatus(tx *gorm.DB, userID uint) (string, error)
	StepUpPassword(tx *gorm.DB, userID uint, password string) error
	PayWithWallet(tx *gorm.DB, userID uint) (*model.Transaction, error)
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
	var wallet *model.Wallet
	result := tx.Model(&wallet).Where("user_id = ?", userID).First(&wallet)
	if result.Error != nil {
		return nil, apperror.InternalServerError("cannot find wallet")
	}

	return wallet, nil
}

func (w *walletRepository) GetTransactionsByUserID(tx *gorm.DB, userID uint) (*[]model.Transaction, error) {
	var transactions *[]model.Transaction
	result := tx.Where("user_id = ?", userID).Find(&transactions)
	if result.Error != nil {
		return nil, apperror.InternalServerError("cannot find transactions")
	}
	return transactions, nil
}

func (w *walletRepository) TransactionDetails(tx *gorm.DB, transactionID uint) (*model.Transaction, error) {
	var transaction *model.Transaction
	result := tx.Where("id = ?", transactionID).First(&transaction)
	if result.Error != nil {
		return nil, apperror.InternalServerError("cannot find transactions")
	}
	return transaction, nil
}

func (w *walletRepository) PaginatedTransactions(tx *gorm.DB, q *Query, userID uint) (int, *[]model.Transaction, error) {
	var trans *[]model.Transaction
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
	totalLength := len(*trans)
	return totalLength, trans, nil
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
		var wallet *model.Wallet
		result1 := tx.Model(&wallet).Where("user_id = ?", userID).First(&wallet)
		if result1.Error != nil {
			return "", apperror.InternalServerError("cannot find wallet")
		}

		return wallet.Status, nil
	}

	return WalletActive, nil
}

func (w *walletRepository) StepUpPassword(tx *gorm.DB, userID uint, password string) error {
	var user *model.User
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

func (w *walletRepository) PayWithWallet(tx *gorm.DB, userID uint) (*model.Transaction, error) {
	var wallet *model.Wallet
	var orderItems *model.OrderItem
	var transaction *model.Transaction

	result1 := tx.Model(&wallet).Where("user_id = ?", userID).First(&wallet)
	if result1.Error != nil {
		return nil, apperror.InternalServerError("cannot find wallet")
	}

	result2 := tx.Model(&orderItems).Where("order_id = ?", userID).Where("order_id is null").First(&orderItems)
	if result2.Error != nil {
		return nil, apperror.InternalServerError("cannot find wallet")
	}

	return transaction, nil
}
