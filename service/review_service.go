package service

import (
	"errors"
	"gorm.io/gorm"
	"seadeals-backend/apperror"
	"seadeals-backend/dto"
	"seadeals-backend/model"
	"seadeals-backend/repository"
)

type ReviewService interface {
	FindReviewByProductID(productID uint, qp *model.ReviewQueryParam) (*dto.GetReviewsRes, error)
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

	if qp.Page == 0 {
		qp.Page = 1
	}
	if qp.Limit == 0 {
		qp.Limit = 6
	}
}

func (s *reviewService) FindReviewByProductID(productID uint, qp *model.ReviewQueryParam) (*dto.GetReviewsRes, error) {
	validateReviewQueryParam(qp)

	tx := s.db.Begin()
	reviews, err := s.reviewRepo.FindReviewByProductID(tx, productID, qp)
	if err != nil {
		if errors.Is(err, &apperror.ReviewNotFoundError{}) {
			return nil, apperror.NotFoundError(err.Error())
		}
		return nil, err
	}

	totalReviews := len(reviews)
	totalPages := (totalReviews + qp.Limit - 1) / qp.Limit

	var reviewsRes []*dto.GetReviewRes
	var avgRating float64
	for _, review := range reviews {
		reviewsRes = append(reviewsRes, new(dto.GetReviewRes).From(review))
		avgRating += float64(review.Rating)
	}
	avgRating = avgRating / float64(totalReviews)

	res := &dto.GetReviewsRes{
		Limit:         uint(qp.Limit),
		Page:          uint(qp.Page),
		TotalPages:    uint(totalPages),
		TotalReviews:  uint(totalReviews),
		AverageRating: avgRating,
		Reviews:       reviewsRes,
	}

	return res, nil
}
