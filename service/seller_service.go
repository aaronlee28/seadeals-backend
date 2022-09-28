package service

import (
	"gorm.io/gorm"
	"seadeals-backend/dto"
	"seadeals-backend/repository"
)

type SellerService interface {
	FindSellerByID(id uint) (*dto.GetSellerRes, error)
}

type sellerService struct {
	db              *gorm.DB
	sellerRepo      repository.SellerRepository
	reviewRepo      repository.ReviewRepository
	socialGraphRepo repository.SocialGraphRepository
}

type SellerServiceConfig struct {
	DB              *gorm.DB
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
	}
}

func (s *sellerService) FindSellerByID(id uint) (*dto.GetSellerRes, error) {
	tx := s.db.Begin()
	seller, err := s.sellerRepo.FindSellerByID(tx, id)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	res := new(dto.GetSellerRes).From(seller)

	averageReview, totalReview, err := s.reviewRepo.GetReviewsAvgAndCountBySellerID(tx, id)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	res.TotalReviewer = uint(totalReview)
	res.Rating = averageReview

	followers, err := s.socialGraphRepo.GetFollowerCountBySellerID(tx, seller.ID)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	res.Followers = uint(followers)

	following, err := s.socialGraphRepo.GetFollowingCountByUserID(tx, seller.UserID)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	res.Following = uint(following)

	tx.Commit()
	return res, nil
}
