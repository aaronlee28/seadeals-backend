package repository

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"seadeals-backend/apperror"
	"seadeals-backend/model"
)

type VoucherRepository interface {
	CreateVoucher(tx *gorm.DB, v *model.Voucher) (*model.Voucher, error)
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
