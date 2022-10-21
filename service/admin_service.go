package service

import (
	"gorm.io/gorm"
	"seadeals-backend/apperror"
	"seadeals-backend/dto"
	"seadeals-backend/helper"
	"seadeals-backend/model"
	"seadeals-backend/repository"
)

type AdminService interface {
	CreateGlobalVoucher(req *dto.CreateGlobalVoucher) (*model.Voucher, error)
	CreateCategory(req *dto.CreateCategory) (*model.ProductCategory, error)
}

type adminService struct {
	db        *gorm.DB
	adminRepo repository.AdminRepository
}

type AdminConfig struct {
	DB        *gorm.DB
	AdminRepo repository.AdminRepository
}

func NewAdminRService(config *AdminConfig) AdminService {
	return &adminService{
		db:        config.DB,
		adminRepo: config.AdminRepo,
	}
}

func (a *adminService) CreateGlobalVoucher(req *dto.CreateGlobalVoucher) (*model.Voucher, error) {
	tx := a.db.Begin()
	var err error
	defer helper.CommitOrRollback(tx, &err)

	if req.AmountType != "percentage" && req.AmountType != "nominal" {
		err = apperror.BadRequestError("invalid amount type")
		return nil, err
	}

	if req.AmountType == "percentage" && req.Amount > 100 {
		err = apperror.BadRequestError("percent amount exceed")
		return nil, err
	}

	createGlobalVoucher := model.Voucher{
		SellerID:    nil,
		Name:        req.Name,
		Code:        req.Code,
		StartDate:   req.StartDate,
		EndDate:     req.EndDate,
		Quota:       int(req.Quota),
		AmountType:  req.AmountType,
		Amount:      req.Amount,
		MinSpending: req.MinSpending,
	}
	var createdGlobalVoucher *model.Voucher
	createdGlobalVoucher, err = a.adminRepo.CreateGlobalVoucher(tx, &createGlobalVoucher)
	if err != nil {
		return nil, err
	}

	return createdGlobalVoucher, nil
}
func (a *adminService) CreateCategory(req *dto.CreateCategory) (*model.ProductCategory, error) {
	tx := a.db.Begin()
	var err error
	defer helper.CommitOrRollback(tx, &err)

	createCategory := model.ProductCategory{
		Name:    req.Name,
		Slug:    req.Slug,
		IconURL: req.IconURL,
	}

	if req.ParentID != nil {
		_, err = a.adminRepo.GetCategoryByID(tx, *req.ParentID)
		if err != nil {
			return nil, err
		}

		createCategory.ParentID = req.ParentID
	}

	var createdCategory *model.ProductCategory
	createdCategory, err = a.adminRepo.CreateCategory(tx, &createCategory)
	if err != nil {
		return nil, err
	}

	return createdCategory, nil
}
