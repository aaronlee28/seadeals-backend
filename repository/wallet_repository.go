package repository

import (
	"gorm.io/gorm"
	"seadeals-backend/apperror"
	"seadeals-backend/model"
	"strconv"
)

type WalletRepository interface {
	CreateWallet(*gorm.DB, *model.Wallet) (*model.Wallet, error)
	GetWalletByUserID(*gorm.DB, uint) (*model.Wallet, error)
	GetTransactionsByUserID(tx *gorm.DB, userID uint) (*[]model.Transaction, error)
	TransactionDetails(tx *gorm.DB, transactionID uint) (*model.Transaction, error)
	PaginatedTransactions(tx *gorm.DB, q *Query, userID uint) (int, *[]model.Transaction, error)
}

type walletRepository struct{}

func NewWalletRepository() WalletRepository {
	return &walletRepository{}
}

type Query struct {
	Limit string
	Page  string
}

func (w *walletRepository) CreateWallet(tx *gorm.DB, wallet *model.Wallet) (*model.Wallet, error) {
	result := tx.Create(&wallet)
	if result.Error != nil {
		return nil, apperror.InternalServerError("cannot create new wallet")
	}

	return wallet, result.Error
}

func (w *walletRepository) GetWalletByUserID(tx *gorm.DB, userID uint) (*model.Wallet, error) {
	var wallet = &model.Wallet{UserID: userID}
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
	totalLength := len(*trans)

	result2 := tx.Limit(limit).Offset(offset).Order("created_at desc").Find(&trans)
	if result2.Error != nil {
		return 0, nil, apperror.InternalServerError("cannot find transactions")
	}
	return totalLength, trans, nil
}
