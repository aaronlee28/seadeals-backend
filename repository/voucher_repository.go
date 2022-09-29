package repository

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"seadeals-backend/apperror"
	"seadeals-backend/model"
	"time"
)

type VoucherRepository interface {
	CreateVoucher(tx *gorm.DB, v *model.Voucher) (*model.Voucher, error)
	UpdateVoucher(tx *gorm.DB, v *model.Voucher, id uint) (*model.Voucher, error)
	FindVoucherByID(tx *gorm.DB, id uint) (*model.Voucher, error)
	FindVoucherDetailByID(tx *gorm.DB, id uint) (*model.Voucher, error)
	FindVoucherBySellerID(tx *gorm.DB, sellerID uint, qp *model.VoucherQueryParam) ([]*model.Voucher, error)
	FindVoucherByCode(tx *gorm.DB, code string) (*model.Voucher, error)
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
		return nil, apperror.NotFoundError(new(apperror.VoucherNotFoundError).Error())
	}
	return updatedVoucher, result.Error
}

func (r *voucherRepository) FindVoucherByID(tx *gorm.DB, id uint) (*model.Voucher, error) {
	var v *model.Voucher
	result := tx.First(&v, id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, apperror.NotFoundError(new(apperror.VoucherNotFoundError).Error())
	}
	return v, result.Error
}

func (r *voucherRepository) FindVoucherDetailByID(tx *gorm.DB, id uint) (*model.Voucher, error) {
	var v *model.Voucher
	result := tx.Preload("Seller.User").First(&v, id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, apperror.NotFoundError(new(apperror.VoucherNotFoundError).Error())
	}
	return v, result.Error
}

func (r *voucherRepository) FindVoucherBySellerID(tx *gorm.DB, sellerID uint, qp *model.VoucherQueryParam) ([]*model.Voucher, error) {
	offset := (qp.Page - 1) * qp.Limit
	orderStmt := fmt.Sprintf("%s %s", qp.SortBy, qp.Sort)

	var queryDB = tx
	switch qp.Status {
	case model.StatusUpcoming:
		queryDB = tx.Where("start_date > ?", time.Now())
	case model.StatusOnGoing:
		queryDB = tx.Where("start_date <= ? AND end_date >= ?", time.Now(), time.Now())
	case model.StatusEnded:
		queryDB = tx.Where("end_date < ?", time.Now())
	}

	var vouchers []*model.Voucher
	result := queryDB.Limit(int(qp.Limit)).Offset(int(offset)).Where("seller_id = ?", sellerID).Preload("Seller.User").Order(orderStmt).Find(&vouchers)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, apperror.NotFoundError(new(apperror.VoucherNotFoundError).Error())
	}
	return vouchers, result.Error
}

func (r *voucherRepository) FindVoucherByCode(tx *gorm.DB, code string) (*model.Voucher, error) {
	var v *model.Voucher
	result := tx.Where("code = ?", code).First(&v)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, apperror.NotFoundError(new(apperror.VoucherNotFoundError).Error())
	}
	return v, result.Error
}

func (r *voucherRepository) DeleteVoucherByID(tx *gorm.DB, id uint) (bool, error) {
	var deletedVoucher *model.Voucher
	result := tx.Delete(&deletedVoucher, id)
	if result.RowsAffected == 0 {
		return false, apperror.NotFoundError(new(apperror.VoucherNotFoundError).Error())
	}
	return true, result.Error
}
