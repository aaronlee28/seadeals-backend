package service

import (
	"fmt"
	"gorm.io/gorm"
	"seadeals-backend/dto"
	"seadeals-backend/helper"
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
		qp.Page = model.PageReviewDefault
	}
	if qp.Limit == 0 {
		qp.Limit = model.LimitReviewDefault
	}
	if qp.Rating < 0 && qp.Rating > 5 {
		qp.Rating = 0
	}
}

func (s *reviewService) FindReviewByProductID(productID uint, qp *model.ReviewQueryParam) (*dto.GetReviewsRes, error) {
	validateReviewQueryParam(qp)

	tx := s.db.Begin()
	var err error
	defer helper.CommitOrRollback(tx, &err)

	reviews, err := s.reviewRepo.FindReviewByProductID(tx, productID, qp)
	if err != nil {
		return nil, err
	}

	totalReviews := uint(len(reviews))
	totalPages := (totalReviews + qp.Limit - 1) / qp.Limit

	var reviewsRes = make([]*dto.GetReviewRes, 0)
	var avgRating float64
	for _, review := range reviews {
		reviewsRes = append(reviewsRes, new(dto.GetReviewRes).From(review))
		avgRating += float64(review.Rating)
	}
	if totalReviews > 0 {
		avgRating = avgRating / float64(totalReviews)
	}

	res := &dto.GetReviewsRes{
		Limit:         qp.Limit,
		Page:          qp.Page,
		TotalPages:    totalPages,
		TotalReviews:  totalReviews,
		AverageRating: avgRating,
		Reviews:       reviewsRes,
	}
	fmt.Println(res)

	return res, nil
}
