package service

import (
	"fmt"
	"gorm.io/gorm"
	"seadeals-backend/apperror"
	"seadeals-backend/dto"
	"seadeals-backend/model"
	"seadeals-backend/repository"
	"strings"
	"time"
)

type VoucherService interface {
	CreateVoucher(req *dto.PostVoucherReq, userID uint) (*dto.GetVoucherRes, error)
	FindVoucherDetailByID(id, userID uint) (*dto.GetVoucherRes, error)
	FindVoucherByID(id uint) (*dto.GetVoucherRes, error)
	FindVoucherBySellerID(sellerID, userID uint, qp *model.VoucherQueryParam) (*dto.GetVouchersRes, error)
	ValidateVoucher(req *dto.PostValidateVoucherReq) (*dto.GetVoucherRes, error)
	UpdateVoucher(req *dto.PatchVoucherReq, id, userID uint) (*dto.GetVoucherRes, error)
	DeleteVoucherByID(id, userID uint) (bool, error)
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

func validateVoucherQueryParam(qp *model.VoucherQueryParam) {
	if !(qp.Sort == "asc" || qp.Sort == "desc") {
		qp.Sort = "desc"
	}
	qp.SortBy = "created_at"

	if qp.Page == 0 {
		qp.Page = model.PageVoucherDefault
	}
	if qp.Limit == 0 {
		qp.Limit = model.LimitVoucherDefault
	}
	if !(qp.Status == model.StatusUpcoming || qp.Status == model.StatusOnGoing || qp.Status == model.StatusEnded) {
		qp.Status = ""
	}
}

func validateVoucher(voucher *model.Voucher, req *dto.PostValidateVoucherReq) error {
	if !(voucher.StartDate.Before(time.Now()) && voucher.EndDate.After(time.Now())) {
		return apperror.BadRequestError(new(apperror.VoucherNotFoundError).Error())
	}
	if voucher.SellerID != req.SellerID {
		return apperror.BadRequestError("voucher cannot be used in this shop")
	}
	if (voucher.Quota - int(req.Quantity)) <= 0 {
		return apperror.BadRequestError(fmt.Sprintf("the number of purchases exceeds the voucher quota (quota: %d)", voucher.Quota))
	}
	if voucher.MinSpending > req.Price {
		return apperror.BadRequestError(fmt.Sprintf("need %.2f more spending", voucher.MinSpending-req.Price))
	}
	return nil
}

func (s *voucherService) CreateVoucher(req *dto.PostVoucherReq, userID uint) (*dto.GetVoucherRes, error) {
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

	res := new(dto.GetVoucherRes).From(voucher)

	tx.Commit()
	return res, nil
}

func (s *voucherService) FindVoucherDetailByID(id, userID uint) (*dto.GetVoucherRes, error) {
	tx := s.db.Begin()
	voucher, err := s.voucherRepo.FindVoucherDetailByID(tx, id)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	if voucher.Seller.UserID != userID {
		tx.Rollback()
		return nil, apperror.UnauthorizedError("cannot fetch other shop detail voucher")
	}

	res := new(dto.GetVoucherRes).From(voucher)

	tx.Commit()
	return res, nil
}

func (s *voucherService) FindVoucherByID(id uint) (*dto.GetVoucherRes, error) {
	tx := s.db.Begin()
	voucher, err := s.voucherRepo.FindVoucherByID(tx, id)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	res := new(dto.GetVoucherRes).From(voucher)

	tx.Commit()
	return res, nil
}

func (s *voucherService) FindVoucherBySellerID(sellerID, userID uint, qp *model.VoucherQueryParam) (*dto.GetVouchersRes, error) {
	tx := s.db.Begin()

	seller, err := s.sellerRepo.FindSellerByID(tx, sellerID)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	if seller.UserID != userID {
		tx.Rollback()
		return nil, apperror.UnauthorizedError("cannot fetch other shop voucher")
	}

	validateVoucherQueryParam(qp)
	vouchers, err := s.voucherRepo.FindVoucherBySellerID(tx, sellerID, qp)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	totalVouchers := uint(len(vouchers))
	totalPages := (totalVouchers + qp.Limit - 1) / qp.Limit

	var voucherRes []*dto.GetVoucherRes
	for _, voucher := range vouchers {
		voucherRes = append(voucherRes, new(dto.GetVoucherRes).From(voucher))
	}

	res := &dto.GetVouchersRes{
		Limit:         qp.Limit,
		Page:          qp.Page,
		TotalPages:    totalPages,
		TotalVouchers: totalVouchers,
		Vouchers:      voucherRes,
	}

	tx.Commit()
	return res, nil
}

func (s *voucherService) ValidateVoucher(req *dto.PostValidateVoucherReq) (*dto.GetVoucherRes, error) {
	tx := s.db.Begin()
	voucher, err := s.voucherRepo.FindVoucherByCode(tx, req.Code)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	err = validateVoucher(voucher, req)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	res := new(dto.GetVoucherRes).From(voucher)

	tx.Commit()
	return res, nil
}

func (s *voucherService) UpdateVoucher(req *dto.PatchVoucherReq, id, userID uint) (*dto.GetVoucherRes, error) {
	tx := s.db.Begin()

	v, err := s.voucherRepo.FindVoucherDetailByID(tx, id)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	if v.Seller.UserID != userID {
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

	err = validateModel(voucher, v.Seller)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	v, err = s.voucherRepo.UpdateVoucher(tx, voucher, id)
	if err != nil {
		tx.Callback()
		return nil, err
	}

	res := new(dto.GetVoucherRes).From(v)

	tx.Commit()
	return res, nil
}

func (s *voucherService) DeleteVoucherByID(id, userID uint) (bool, error) {
	tx := s.db.Begin()

	v, err := s.voucherRepo.FindVoucherDetailByID(tx, id)
	if err != nil {
		tx.Rollback()
		return false, err
	}

	if v.Seller.User.ID != userID {
		tx.Rollback()
		return false, apperror.UnauthorizedError("cannot delete other shop voucher")
	}

	voucher, err := s.voucherRepo.FindVoucherDetailByID(tx, id)
	if err != nil {
		tx.Rollback()
		return false, err
	}
	if voucher.StartDate.Before(time.Now()) {
		tx.Rollback()
		return false, apperror.BadRequestError("cannot delete voucher that has been started")
	}

	isDeleted, err := s.voucherRepo.DeleteVoucherByID(tx, id)
	if err != nil {
		tx.Rollback()
		return false, err
	}

	tx.Commit()
	return isDeleted, nil
}
