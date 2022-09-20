package repository

import (
	"gorm.io/gorm"
	"seadeals-backend/apperror"
	"seadeals-backend/model"
)

type WalletRepository interface {
	CreateWallet(*gorm.DB, *model.Wallet) (*model.Wallet, error)
	GetWalletByUserID(*gorm.DB, uint) (*model.Wallet, error)
	GetTransactionsByUserID(tx *gorm.DB, userID uint) (*[]model.Transaction, error)
}

type walletRepository struct{}

func NewWalletRepository() WalletRepository {
	return &walletRepository{}
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
	result := tx.Model(&wallet).First(&wallet)
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
