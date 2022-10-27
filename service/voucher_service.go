package service

import (
	"fmt"
	"seadeals-backend/apperror"
	"seadeals-backend/dto"
	"seadeals-backend/helper"
	"seadeals-backend/model"
	"seadeals-backend/repository"
	"strings"
	"time"

	"gorm.io/gorm"
)

type VoucherService interface {
	CreateVoucher(req *dto.PostVoucherReq, userID uint) (*dto.GetVoucherRes, error)
	FindVoucherDetailByID(id, userID uint) (*dto.GetVoucherRes, error)
	FindVoucherByID(id uint) (*dto.GetVoucherRes, error)
	FindVoucherByUserID(userID uint, qp *model.VoucherQueryParam) (*dto.GetVouchersRes, error)
	ValidateVoucher(req *dto.PostValidateVoucherReq) (*dto.GetVoucherRes, error)
	UpdateVoucher(req *dto.PatchVoucherReq, id, userID uint) (*dto.GetVoucherRes, error)
	//
	DeleteVoucherByID(id, userID uint) (bool, error)
	GetVouchersBySellerID(sellerID uint) ([]*dto.GetVoucherRes, error)
	GetAvailableGlobalVouchers() ([]*dto.GetVoucherRes, error)
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
	if qp.Month == 0 {
		qp.Month = model.MonthVoucherDefault
	}
}

func validateVoucher(voucher *model.Voucher, req *dto.PostValidateVoucherReq) error {
	if !(voucher.StartDate.Before(time.Now()) && voucher.EndDate.After(time.Now())) {
		return apperror.BadRequestError(new(apperror.VoucherNotFoundError).Error())
	}
	if *voucher.SellerID != req.SellerID {
		return apperror.BadRequestError("voucher cannot be used in this shop")
	}
	if voucher.MinSpending > req.Spend {
		return apperror.BadRequestError(fmt.Sprintf("need Rp%.2f more spending", voucher.MinSpending-req.Spend))
	}
	return nil
}

func (s *voucherService) CreateVoucher(req *dto.PostVoucherReq, userID uint) (*dto.GetVoucherRes, error) {
	tx := s.db.Begin()
	var err error
	defer helper.CommitOrRollback(tx, &err)

	seller, err := s.sellerRepo.FindSellerByUserID(tx, userID)
	if err != nil {
		return nil, err
	}

	voucher := &model.Voucher{
		SellerID:    &seller.ID,
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
		return nil, err
	}

	res := new(dto.GetVoucherRes).From(voucher)
	return res, nil
}

func (s *voucherService) FindVoucherDetailByID(id, userID uint) (*dto.GetVoucherRes, error) {
	tx := s.db.Begin()
	var err error
	defer helper.CommitOrRollback(tx, &err)

	voucher, err := s.voucherRepo.FindVoucherDetailByID(tx, id)
	if err != nil {
		return nil, err
	}

	if voucher.Seller.UserID != userID {
		err = apperror.UnauthorizedError("cannot fetch other shop detail voucher")
		return nil, err
	}

	res := new(dto.GetVoucherRes).From(voucher)
	return res, nil
}

func (s *voucherService) FindVoucherByID(id uint) (*dto.GetVoucherRes, error) {
	tx := s.db.Begin()
	var err error
	defer helper.CommitOrRollback(tx, &err)

	voucher, err := s.voucherRepo.FindVoucherByID(tx, id)
	if err != nil {
		return nil, err
	}

	res := new(dto.GetVoucherRes).From(voucher)
	return res, nil
}

func (s *voucherService) FindVoucherByUserID(userID uint, qp *model.VoucherQueryParam) (*dto.GetVouchersRes, error) {
	tx := s.db.Begin()
	var err error
	defer helper.CommitOrRollback(tx, &err)

	seller, err := s.sellerRepo.FindSellerByUserID(tx, userID)
	if err != nil {
		return nil, err
	}

	validateVoucherQueryParam(qp)
	vouchers, totalVouchers, err := s.voucherRepo.FindVoucherBySellerID(tx, seller.ID, qp)
	if err != nil {
		return nil, err
	}

	totalPages := (uint(totalVouchers) + qp.Limit - 1) / qp.Limit

	var voucherRes = make([]*dto.GetVoucherRes, 0)
	for _, voucher := range vouchers {
		voucherRes = append(voucherRes, new(dto.GetVoucherRes).From(voucher))
	}

	res := &dto.GetVouchersRes{
		Limit:         qp.Limit,
		Page:          qp.Page,
		TotalPages:    totalPages,
		TotalVouchers: uint(totalVouchers),
		Vouchers:      voucherRes,
	}

	return res, nil
}

func (s *voucherService) ValidateVoucher(req *dto.PostValidateVoucherReq) (*dto.GetVoucherRes, error) {
	tx := s.db.Begin()
	var err error
	defer helper.CommitOrRollback(tx, &err)

	voucher, err := s.voucherRepo.FindVoucherByCode(tx, req.Code)
	if err != nil {
		return nil, err
	}

	err = validateVoucher(voucher, req)
	if err != nil {
		return nil, err
	}

	res := new(dto.GetVoucherRes).From(voucher)
	return res, nil
}

func (s *voucherService) UpdateVoucher(req *dto.PatchVoucherReq, id, userID uint) (*dto.GetVoucherRes, error) {
	tx := s.db.Begin()
	var err error
	defer helper.CommitOrRollback(tx, &err)

	v, err := s.voucherRepo.FindVoucherDetailByID(tx, id)
	if err != nil {
		return nil, err
	}

	if v.Seller.UserID != userID {
		err = apperror.UnauthorizedError("cannot update other shop voucher")
		return nil, err
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
		return nil, err
	}

	v, err = s.voucherRepo.UpdateVoucher(tx, voucher, id)
	if err != nil {
		return nil, err
	}

	res := new(dto.GetVoucherRes).From(v)
	return res, nil
}

func (s *voucherService) DeleteVoucherByID(id, userID uint) (bool, error) {
	tx := s.db.Begin()
	var err error
	defer helper.CommitOrRollback(tx, &err)

	v, err := s.voucherRepo.FindVoucherDetailByID(tx, id)
	if err != nil {
		return false, err
	}

	if v.Seller.UserID != userID {
		err = apperror.UnauthorizedError("cannot delete other shop voucher")
		return false, err
	}
	if v.StartDate.Before(time.Now()) {
		err = apperror.BadRequestError("cannot delete voucher that has been started")
		return false, err
	}
	isDeleted, err := s.voucherRepo.DeleteVoucherByID(tx, id)
	if err != nil {
		return false, err
	}
	return isDeleted, nil
}

func (s *voucherService) GetVouchersBySellerID(sellerID uint) ([]*dto.GetVoucherRes, error) {
	tx := s.db.Begin()
	var err error
	defer helper.CommitOrRollback(tx, &err)

	vouchers, err := s.voucherRepo.GetVouchersBySellerID(tx, sellerID)
	if err != nil {
		return nil, err
	}

	var voucherRes = make([]*dto.GetVoucherRes, 0)
	for _, voucher := range vouchers {
		voucherRes = append(voucherRes, new(dto.GetVoucherRes).From(voucher))
	}

	return voucherRes, nil
}

func (s *voucherService) GetAvailableGlobalVouchers() ([]*dto.GetVoucherRes, error) {
	tx := s.db.Begin()
	var err error
	defer helper.CommitOrRollback(tx, &err)

	vouchers, err := s.voucherRepo.GetAvailableGlobalVouchers(tx)
	if err != nil {
		return nil, err
	}

	var voucherRes = make([]*dto.GetVoucherRes, 0)
	for _, voucher := range vouchers {
		var sellerID uint
		voucher.SellerID = &sellerID
		voucherRes = append(voucherRes, new(dto.GetVoucherRes).From(voucher))
	}

	return voucherRes, nil
}
