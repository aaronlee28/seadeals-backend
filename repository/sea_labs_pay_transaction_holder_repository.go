package repository

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"seadeals-backend/apperror"
	"seadeals-backend/model"
)

type SeaLabsPayTransactionHolderRepository interface {
	GetTransHolderFromTransactionID(tx *gorm.DB, transactionID uint) (*model.SeaLabsPayTransactionHolder, error)

	CreateTransactionHolder(tx *gorm.DB, model *model.SeaLabsPayTransactionHolder) (*model.SeaLabsPayTransactionHolder, error)
	UpdateTransactionHolder(tx *gorm.DB, txnID uint, status string) (*model.SeaLabsPayTransactionHolder, error)
}

type seaLabsPayTransactionHolderRepository struct{}

func NewSeaLabsPayTransactionHolderRepository() SeaLabsPayTransactionHolderRepository {
	return &seaLabsPayTransactionHolderRepository{}
}

func (s *seaLabsPayTransactionHolderRepository) GetTransHolderFromTransactionID(tx *gorm.DB, transactionID uint) (*model.SeaLabsPayTransactionHolder, error) {
	var transHolder = &model.SeaLabsPayTransactionHolder{}
	result := tx.Model(&transHolder).Where("transaction_id = ?", transactionID).First(&transHolder)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, apperror.BadRequestError("That Transactions Holder doesn't exists")
		}
		return nil, apperror.InternalServerError("Cannot find Transaction Holder")
	}
	return transHolder, nil
}

func (s *seaLabsPayTransactionHolderRepository) CreateTransactionHolder(tx *gorm.DB, model *model.SeaLabsPayTransactionHolder) (*model.SeaLabsPayTransactionHolder, error) {
	result := tx.Create(&model)
	if result.Error != nil {
		return nil, apperror.InternalServerError("Cannot create sea labs pay transaction holder")
	}

	return model, nil
}

func (s *seaLabsPayTransactionHolderRepository) UpdateTransactionHolder(tx *gorm.DB, txnID uint, status string) (*model.SeaLabsPayTransactionHolder, error) {
	var existingData = &model.SeaLabsPayTransactionHolder{}
	result := tx.Model(existingData).Where("txn_id = ?", txnID).First(existingData)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, apperror.BadRequestError("No such transaction exist")
		}
		return nil, apperror.InternalServerError("Cannot find top up transaction holder")
	}

	existingData.TransactionStatus = &status
	result = tx.Model(&existingData).Clauses(clause.Returning{}).Updates(&existingData)
	if result.Error != nil {
		return nil, apperror.InternalServerError("Cannot create sea labs pay top up holder")
	}

	return existingData, nil
}
