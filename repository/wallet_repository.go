package repository

import (
	"gorm.io/gorm"
	"seadeals-backend/apperror"
	"seadeals-backend/model"
)

type WalletRepository interface {
	CreateWallet(tx *gorm.DB, wallet *model.Wallet) (*model.Wallet, error)
}

type walletRepository struct {
}

type WalletRepositoryConfig struct {
}

func NewWalletRepository(c *WalletRepositoryConfig) WalletRepository {
	return &walletRepository{}
}

func (w *walletRepository) CreateWallet(tx *gorm.DB, wallet *model.Wallet) (*model.Wallet, error) {
	result := tx.Create(&wallet)
	if result.Error != nil {
		return nil, apperror.InternalServerError("cannot create new wallet")
	}

	return wallet, result.Error
}
