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
