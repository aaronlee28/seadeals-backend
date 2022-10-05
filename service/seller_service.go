package service

import (
	"gorm.io/gorm"
	"seadeals-backend/dto"
	"seadeals-backend/helper"
	"seadeals-backend/repository"
)

type SellerService interface {
	FindSellerByID(sellerID uint, userID uint) (*dto.GetSellerRes, error)
}

type sellerService struct {
	db              *gorm.DB
	productRepo     repository.ProductRepository
	sellerRepo      repository.SellerRepository
	reviewRepo      repository.ReviewRepository
	socialGraphRepo repository.SocialGraphRepository
}

type SellerServiceConfig struct {
	DB              *gorm.DB
	ProductRepo     repository.ProductRepository
	SellerRepo      repository.SellerRepository
	ReviewRepo      repository.ReviewRepository
	SocialGraphRepo repository.SocialGraphRepository
}

func NewSellerService(c *SellerServiceConfig) SellerService {
	return &sellerService{
		db:              c.DB,
		sellerRepo:      c.SellerRepo,
		reviewRepo:      c.ReviewRepo,
		socialGraphRepo: c.SocialGraphRepo,
		productRepo:     c.ProductRepo,
	}
}

func (s *sellerService) FindSellerByID(sellerID uint, userID uint) (*dto.GetSellerRes, error) {
	tx := s.db.Begin()
	var err error
	defer helper.CommitOrRollback(tx, &err)

	seller, err := s.sellerRepo.FindSellerDetailByID(tx, sellerID, userID)
	if err != nil {
		return nil, err
	}

	res := new(dto.GetSellerRes).From(seller)

	averageReview, totalReview, err := s.reviewRepo.GetReviewsAvgAndCountBySellerID(tx, sellerID)
	if err != nil {
		return nil, err
	}
	res.TotalReviewer = uint(totalReview)
	res.Rating = averageReview

	followers, err := s.socialGraphRepo.GetFollowerCountBySellerID(tx, seller.ID)
	if err != nil {
		return nil, err
	}
	res.Followers = uint(followers)

	following, err := s.socialGraphRepo.GetFollowingCountByUserID(tx, seller.UserID)
	if err != nil {
		return nil, err
	}
	res.Following = uint(following)

	totalProduct, err := s.productRepo.GetProductCountBySellerID(tx, seller.ID)
	if err != nil {
		return nil, err
	}
	res.TotalProduct = totalProduct

	return res, nil
}
