package repository

import (
	"gorm.io/gorm"
	"seadeals-backend/apperror"
	"seadeals-backend/model"
)

type TransactionRepository interface {
	GetPriceBeforeGlobalDisc(tx *gorm.DB, transactionID uint) (float64, error)
}

type transactionRepository struct{}

func NewTransactionRepository() TransactionRepository {
	return &transactionRepository{}
}

func (t *transactionRepository) GetPriceBeforeGlobalDisc(tx *gorm.DB, transactionID uint) (float64, error) {
	var price float64
	result := tx.Model(&model.Order{}).Select("SUM(total) as price").Where("transaction_id = ?", transactionID).Scan(&price)
	if result.Error != nil {
		return 0, apperror.InternalServerError("Cannot sum total of orders")
	}
	return price, nil
}
