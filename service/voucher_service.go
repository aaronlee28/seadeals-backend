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

func validateRequest(req *dto.PostVoucherReq, seller *model.Seller) error {
	username := seller.User.Username[:4]
	req.Code = strings.ToUpper(username + req.Code)

	req.AmountType = strings.ToLower(req.AmountType)
	if req.AmountType != model.PercentageType && req.AmountType != model.NominalType {
		req.AmountType = model.PercentageType
	}

	if req.AmountType == model.PercentageType {
		if req.Amount > 100 {
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

	err = validateRequest(req, seller)
	if err != nil {
		tx.Rollback()
		return nil, err
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

	voucher, err = s.voucherRepo.CreateVoucher(tx, voucher)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	tx.Commit()
	return voucher, nil
}
