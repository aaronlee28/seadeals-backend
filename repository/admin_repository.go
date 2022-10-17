package repository

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"seadeals-backend/apperror"
	"seadeals-backend/model"
)

type AdminRepository interface {
	CreateGlobalVoucher(tx *gorm.DB, req *model.Voucher) (*model.Voucher, error)
	GetCategoryByID(tx *gorm.DB, categoryID uint) (*model.ProductCategory, error)
	CreateCategory(tx *gorm.DB, req *model.ProductCategory) (*model.ProductCategory, error)
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

func (a *adminRepository) GetCategoryByID(tx *gorm.DB, categoryID uint) (*model.ProductCategory, error) {
	var find *model.ProductCategory
	result := tx.Clauses(clause.Returning{}).Where("id = ?", categoryID).First(&find)
	if result.Error != nil {
		return nil, apperror.InternalServerError("category id not found")
	}
	return find, nil
}

func (a *adminRepository) CreateCategory(tx *gorm.DB, req *model.ProductCategory) (*model.ProductCategory, error) {
	result := tx.Clauses(clause.Returning{}).Create(&req)
	if result.Error != nil {
		return nil, apperror.InternalServerError("Cannot create category")
	}
	return req, nil
}
