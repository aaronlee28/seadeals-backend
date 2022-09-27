package repository

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"seadeals-backend/apperror"
	"seadeals-backend/model"
)

type SeaLabsPayTopUpHolderRepository interface {
	CreateTopUpHolder(*gorm.DB, *model.SeaLabsPayTopUpHolder) (*model.SeaLabsPayTopUpHolder, error)
	UpdateTopUpHolder(tx *gorm.DB, txnID uint, status string) (*model.SeaLabsPayTopUpHolder, error)
}

type seaLabsPayTopUpHolderRepository struct{}

func NewSeaLabsPayTopUpHolderRepository() SeaLabsPayTopUpHolderRepository {
	return &seaLabsPayTopUpHolderRepository{}
}

func (s *seaLabsPayTopUpHolderRepository) CreateTopUpHolder(tx *gorm.DB, model *model.SeaLabsPayTopUpHolder) (*model.SeaLabsPayTopUpHolder, error) {
	result := tx.Create(&model)
	if result.Error != nil {
		return nil, apperror.InternalServerError("Cannot create sea labs pay top up holder")
	}

	return model, nil
}

func (s *seaLabsPayTopUpHolderRepository) UpdateTopUpHolder(tx *gorm.DB, txnID uint, status string) (*model.SeaLabsPayTopUpHolder, error) {
	var existingData = &model.SeaLabsPayTopUpHolder{}
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
