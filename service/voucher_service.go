package service

import (
	"gorm.io/gorm"
	"seadeals-backend/apperror"
	"seadeals-backend/dto"
	"seadeals-backend/model"
	"seadeals-backend/repository"
	"strings"
)

type VoucherService interface {
	CreateVoucher(req *dto.PostVoucherReq, userID uint) (*model.Voucher, error)
	UpdateVoucher(req *dto.PatchVoucherReq, id, userID uint) (*model.Voucher, error)
}

type voucherService struct {
	db          *gorm.DB
	voucherRepo repository.VoucherRepository
	sellerRepo  repository.SellerRepository
}

type VoucherServiceConfig struct {
	DB          *gorm.DB
	VoucherRepo repository.VoucherRepository
	SellerRepo  repository.SellerRepository
}

func NewVoucherService(c *VoucherServiceConfig) VoucherService {
	return &voucherService{
		db:          c.DB,
		voucherRepo: c.VoucherRepo,
		sellerRepo:  c.SellerRepo,
	}
}

func validateModel(v *model.Voucher, seller *model.Seller) error {
	if v.Code != "" {
		username := seller.User.Username[:4]
		v.Code = strings.ToUpper(username + v.Code)
	}

	v.AmountType = strings.ToLower(v.AmountType)
	if v.AmountType != model.PercentageType && v.AmountType != model.NominalType {
		v.AmountType = model.NominalType
	}

	if v.AmountType == model.PercentageType {
		if v.Amount > 100 {
			return apperror.BadRequestError("percentage amount must be in range 1-100")
		}
	}
	return nil
}

func (s *voucherService) CreateVoucher(req *dto.PostVoucherReq, userID uint) (*model.Voucher, error) {
	tx := s.db.Begin()

	seller, err := s.sellerRepo.FindSellerByID(tx, req.SellerID)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	if seller.UserID != userID {
		tx.Rollback()
		return nil, apperror.UnauthorizedError("cannot add other shop voucher")
	}

	voucher := &model.Voucher{
		SellerID:    req.SellerID,
		Name:        req.Name,
		Code:        req.Code,
		StartDate:   req.StartDate,
		EndDate:     req.EndDate,
		Quota:       req.Quota,
		AmountType:  req.AmountType,
		Amount:      req.Amount,
		MinSpending: req.MinSpending,
	}

	err = validateModel(voucher, seller)
	if err != nil {
		return nil, apperror.BadRequestError(err.Error())
	}

	voucher, err = s.voucherRepo.CreateVoucher(tx, voucher)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	tx.Commit()
	return voucher, nil
}

func (s *voucherService) UpdateVoucher(req *dto.PatchVoucherReq, id, userID uint) (*model.Voucher, error) {
	tx := s.db.Begin()

	seller, err := s.sellerRepo.FindSellerByUserID(tx, userID)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	if seller.UserID != userID {
		tx.Rollback()
		return nil, apperror.UnauthorizedError("cannot update other shop voucher")
	}

	voucher := &model.Voucher{
		Name:        req.Name,
		StartDate:   req.StartDate,
		EndDate:     req.EndDate,
		Quota:       req.Quota,
		AmountType:  req.AmountType,
		Amount:      req.Amount,
		MinSpending: req.MinSpending,
	}

	err = validateModel(voucher, seller)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	v, err := s.voucherRepo.UpdateVoucher(tx, voucher, id)
	if err != nil {
		tx.Callback()
		return nil, err
	}

	tx.Commit()
	return v, nil
}
