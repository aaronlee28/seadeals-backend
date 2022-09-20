package service

import (
	"gorm.io/gorm"
	"seadeals-backend/dto"
	"seadeals-backend/repository"
)

type WalletService interface {
	UserWalletData(id uint) (*dto.WalletDataRes, error)
	TransactionDetails(id uint) (*dto.TransactionDetailsRes, error)
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
	transactions, err := w.walletRepository.GetTransactionsByUserID(tx, id)
	walletData := &dto.WalletDataRes{
		UserID:       wallet.UserID,
		Balance:      wallet.Balance,
		Transactions: transactions,
	}
	return walletData, nil
}

func (w *walletService) TransactionDetails(id uint) (*dto.TransactionDetailsRes, error) {
	tx := w.db.Begin()
	t, err := w.walletRepository.TransactionDetails(tx, id)
	if err != nil {
		return nil, err
	}
	transaction := &dto.TransactionDetailsRes{
		Id:            t.Id,
		VoucherID:     t.VoucherID,
		Total:         t.Total,
		PaymentType:   t.PaymentType,
		PaymentMethod: t.PaymentMethod,
		CreatedAt:     t.CreatedAt,
		UpdatedAt:     t.UpdatedAt,
	}
	return transaction, nil
}

//
//func (w *walletService) PaginatedTransactions(id uint) (*dto.PaginatedTransactionRes, error) {
//	tx := w.db.Begin()
//	transactions, err := w.walletRepository.GetTransactionsByUserID(tx, id)
//
//}
