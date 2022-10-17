package repository

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"seadeals-backend/apperror"
	"seadeals-backend/model"
)

type AccountHolderRepository interface {
	SendToAccountHolder(tx *gorm.DB, model *model.AccountHolder) (*model.AccountHolder, error)
	TakeMoneyFromAccountHolderByOrderID(tx *gorm.DB, orderID uint) (*model.AccountHolder, error)
}

type accountHolderRepository struct{}

func NewAccountHolderRepository() AccountHolderRepository {
	return &accountHolderRepository{}
}

func (a *accountHolderRepository) SendToAccountHolder(tx *gorm.DB, model *model.AccountHolder) (*model.AccountHolder, error) {
	result := tx.Create(&model)
	if result.Error != nil {
		return nil, apperror.InternalServerError("Tidak bisa membuat account holding")
	}
	return model, nil
}

func (a *accountHolderRepository) TakeMoneyFromAccountHolderByOrderID(tx *gorm.DB, orderID uint) (*model.AccountHolder, error) {
	var accountHolder = &model.AccountHolder{HasTaken: true}
	result := tx.Model(&accountHolder).Clauses(clause.Returning{}).Where("order_id = ?", orderID).Updates(&accountHolder)
	if result.Error != nil {
		return nil, apperror.InternalServerError("Tidak bisa update account holder")
	}
	return accountHolder, nil
}
