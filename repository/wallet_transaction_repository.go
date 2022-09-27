package repository

import (
	"gorm.io/gorm"
	"seadeals-backend/apperror"
	"seadeals-backend/model"
)

type WalletTransactionRepository interface {
	CreateTransaction(tx *gorm.DB, model *model.WalletTransaction) (*model.WalletTransaction, error)
}

type walletTransactionRepository struct{}

func NewWalletTransactionRepository() WalletTransactionRepository {
	return &walletTransactionRepository{}
}

func (w *walletTransactionRepository) CreateTransaction(tx *gorm.DB, model *model.WalletTransaction) (*model.WalletTransaction, error) {
	result := tx.Create(&model)
	if result.Error != nil {
		return nil, apperror.InternalServerError("Cannot create wallet transaction")
	}
	return model, nil
}
