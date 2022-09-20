package service

import (
	"gorm.io/gorm"
	"seadeals-backend/dto"
	"seadeals-backend/repository"
)

type WalletService interface {
	UserWalletData(id uint) (*dto.WalletDataRes, error)
}

type walletService struct {
	db               *gorm.DB
	walletRepository repository.WalletRepository
}

type WalletServiceConfig struct {
	DB               *gorm.DB
	WalletRepository repository.WalletRepository
}

func NewWalletService(c *WalletServiceConfig) WalletService {
	return &walletService{
		db:               c.DB,
		walletRepository: c.WalletRepository,
	}
}

func (w *walletService) UserWalletData(id uint) (*dto.WalletDataRes, error) {
	tx := w.db.Begin()
	wallet, err := w.walletRepository.GetWalletByUserID(tx, id)
	if err != nil {
		return nil, err
	}
	walletData := &dto.WalletDataRes{
		UserID:  wallet.UserID,
		Balance: wallet.Balance,
	}
	return walletData, nil
}
