package repository

import (
	"errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"seadeals-backend/apperror"
	"seadeals-backend/model"
)

type VoucherRepository interface {
	CreateVoucher(tx *gorm.DB, v *model.Voucher) (*model.Voucher, error)
	UpdateVoucher(tx *gorm.DB, v *model.Voucher, id uint) (*model.Voucher, error)
	FindVoucherByID(tx *gorm.DB, id uint) (*model.Voucher, error)
	FindVoucherDetailByID(tx *gorm.DB, id uint) (*model.Voucher, error)
	DeleteVoucherByID(tx *gorm.DB, id uint) (bool, error)
}

type voucherRepository struct{}

func NewVoucherRepository() VoucherRepository {
	return &voucherRepository{}
}

func (r *voucherRepository) CreateVoucher(tx *gorm.DB, v *model.Voucher) (*model.Voucher, error) {
	result := tx.Clauses(clause.OnConflict{DoNothing: true}).Create(&v)
	if int(result.RowsAffected) == 0 {
		return nil, apperror.BadRequestError("code already exist")
	}
	return v, nil
}

func (r *voucherRepository) UpdateVoucher(tx *gorm.DB, v *model.Voucher, id uint) (*model.Voucher, error) {
	var updatedVoucher *model.Voucher
	result := tx.First(&updatedVoucher, id).Updates(&v)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, apperror.NotFoundError("voucher not found error")
	}
	return updatedVoucher, result.Error
}

func (r *voucherRepository) FindVoucherByID(tx *gorm.DB, id uint) (*model.Voucher, error) {
	var v *model.Voucher
	result := tx.First(&v, id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, apperror.NotFoundError("voucher not found error")
	}
	return v, result.Error
}

func (r *voucherRepository) FindVoucherDetailByID(tx *gorm.DB, id uint) (*model.Voucher, error) {
	var v *model.Voucher
	result := tx.Preload("Seller.User").First(&v, id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, apperror.NotFoundError("voucher not found error")
	}
	return v, result.Error
}

func (r *voucherRepository) DeleteVoucherByID(tx *gorm.DB, id uint) (bool, error) {
	var deletedVoucher *model.Voucher
	result := tx.Delete(&deletedVoucher, id)
	if result.RowsAffected == 0 {
		return false, apperror.NotFoundError("voucher not found error")
	}
	return true, result.Error
}
