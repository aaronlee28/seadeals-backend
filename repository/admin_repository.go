package repository

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"seadeals-backend/apperror"
	"seadeals-backend/model"
)

type AdminRepository interface {
	CreateGlobalVoucher(tx *gorm.DB, req *model.Voucher) (*model.Voucher, error)
}

type adminRepository struct{}

func NewAdminRepository() AdminRepository {
	return &adminRepository{}
}

func (a *adminRepository) CreateGlobalVoucher(tx *gorm.DB, req *model.Voucher) (*model.Voucher, error) {
	result := tx.Clauses(clause.Returning{}).Create(&req)
	if result.Error != nil {
		return nil, apperror.InternalServerError("Cannot create global voucher")
	}
	return req, nil
}
