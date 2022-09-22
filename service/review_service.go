package service

import (
	"gorm.io/gorm"
	"seadeals-backend/dto"
	"seadeals-backend/model"
	"seadeals-backend/repository"
)

type ReviewService interface {
	FindReviewByProductID(productID uint, qp *model.ReviewQueryParam) ([]*dto.GetReviewRes, error)
}

type reviewService struct {
	db         *gorm.DB
	reviewRepo repository.ReviewRepository
}

type ReviewServiceConfig struct {
	DB         *gorm.DB
	ReviewRepo repository.ReviewRepository
}

func NewReviewService(config *ReviewServiceConfig) ReviewService {
	return &reviewService{
		db:         config.DB,
		reviewRepo: config.ReviewRepo,
	}
}

func validateReviewQueryParam(qp *model.ReviewQueryParam) {
	if !(qp.Sort == "asc" || qp.Sort == "desc") {
		qp.Sort = "desc"
	}
	qp.SortBy = "created_at"
}

func (s *reviewService) FindReviewByProductID(productID uint, qp *model.ReviewQueryParam) ([]*dto.GetReviewRes, error) {
	validateReviewQueryParam(qp)

	tx := s.db.Begin()
	reviews, err := s.reviewRepo.FindReviewByProductID(tx, productID, qp)
	if err != nil {
		return nil, err
	}

	var res []*dto.GetReviewRes
	for _, review := range reviews {
		res = append(res, new(dto.GetReviewRes).From(review))
	}

	return res, nil
}
