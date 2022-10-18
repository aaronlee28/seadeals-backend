package service

import (
	"gorm.io/gorm"
	"seadeals-backend/dto"
	"seadeals-backend/helper"
	"seadeals-backend/model"
	"seadeals-backend/repository"
)

type SellerAvailableCourService interface {
	CreateOrUpdateCourier(req *dto.AddDeliveryReq, userID uint) (*model.SellerAvailableCourier, error)
	GetAvailableCourierForSeller(userID uint) ([]*model.SellerAvailableCourier, error)
	GetAvailableCourierForBuyer(sellerID uint) ([]*model.SellerAvailableCourier, error)
}

type sellerAvailableCourService struct {
	db                  *gorm.DB
	sellerAvailCourRepo repository.SellerAvailableCourierRepository
	sellerRepository    repository.SellerRepository
}

type SellerAvailableCourServiceConfig struct {
	DB                  *gorm.DB
	SellerAvailCourRepo repository.SellerAvailableCourierRepository
	SellerRepository    repository.SellerRepository
}

func NewSellerAvailableCourService(c *SellerAvailableCourServiceConfig) SellerAvailableCourService {
	return &sellerAvailableCourService{
		db:                  c.DB,
		sellerAvailCourRepo: c.SellerAvailCourRepo,
		sellerRepository:    c.SellerRepository,
	}
}

func (s *sellerAvailableCourService) CreateOrUpdateCourier(req *dto.AddDeliveryReq, userID uint) (*model.SellerAvailableCourier, error) {
	tx := s.db.Begin()
	var err error
	defer helper.CommitOrRollback(tx, &err)

	seller, err := s.sellerRepository.FindSellerByUserID(tx, userID)
	if err != nil {
		return nil, err
	}

	courier, err := s.sellerAvailCourRepo.AddSellerAvailableDeliveryMethod(tx, req, seller.ID)
	if err != nil {
		return nil, err
	}

	return courier, nil
}

func (s *sellerAvailableCourService) GetAvailableCourierForSeller(userID uint) ([]*model.SellerAvailableCourier, error) {
	tx := s.db.Begin()
	var err error
	defer helper.CommitOrRollback(tx, &err)

	seller, err := s.sellerRepository.FindSellerByUserID(tx, userID)
	if err != nil {
		return nil, err
	}

	couriers, err := s.sellerAvailCourRepo.GetAllSellerAvailableCourier(tx, seller.ID)
	if err != nil {
		return nil, err
	}

	return couriers, nil
}

func (s *sellerAvailableCourService) GetAvailableCourierForBuyer(sellerID uint) ([]*model.SellerAvailableCourier, error) {
	tx := s.db.Begin()
	var err error
	defer helper.CommitOrRollback(tx, &err)

	couriers, err := s.sellerAvailCourRepo.GetAllSellerAvailableCourier(tx, sellerID)
	if err != nil {
		return nil, err
	}

	return couriers, nil
}
