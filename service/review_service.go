package service

import (
	"fmt"
	"gorm.io/gorm"
	"math"
	"seadeals-backend/dto"
	"seadeals-backend/helper"
	"seadeals-backend/model"
	"seadeals-backend/repository"
)

type ReviewService interface {
	FindReviewByProductID(productID uint, qp *model.ReviewQueryParam) (*dto.GetReviewsRes, error)
	CreateUpdateReview(userID uint, req *dto.CreateUpdateReview) (*model.Review, error)
	UserReviewHistory(userID uint) ([]*model.Review, error)
}

type reviewService struct {
	db          *gorm.DB
	reviewRepo  repository.ReviewRepository
	sellerRepo  repository.SellerRepository
	productRepo repository.ProductRepository
}

type ReviewServiceConfig struct {
	DB          *gorm.DB
	ReviewRepo  repository.ReviewRepository
	SellerRepo  repository.SellerRepository
	ProductRepo repository.ProductRepository
}

func NewReviewService(config *ReviewServiceConfig) ReviewService {
	return &reviewService{
		db:          config.DB,
		reviewRepo:  config.ReviewRepo,
		sellerRepo:  config.SellerRepo,
		productRepo: config.ProductRepo,
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
	//if qp.Limit == 0 {
	//	qp.Limit = ""
	//}
	if qp.Rating < 0 && qp.Rating > 5 {
		qp.Rating = 0
	}
}

func (s *reviewService) FindReviewByProductID(productID uint, qp *model.ReviewQueryParam) (*dto.GetReviewsRes, error) {
	validateReviewQueryParam(qp)

	tx := s.db.Begin()
	var err error
	defer helper.CommitOrRollback(tx, &err)
	reviewsRaw, err := s.reviewRepo.FindReviewByProductIDNoQuery(tx, productID)
	fmt.Println("totalReview", len(reviewsRaw))
	reviews, err := s.reviewRepo.FindReviewByProductID(tx, productID, qp)
	if err != nil {
		return nil, err
	}

	totalReviews := uint(len(reviews))
	totalRawReviews := uint(len(reviewsRaw))

	totalPages := (totalReviews + qp.Limit - 1) / qp.Limit

	var reviewsRes = make([]*dto.GetReviewRes, 0)
	var avgRating float64

	for _, review := range reviews {
		reviewsRes = append(reviewsRes, new(dto.GetReviewRes).From(review))
	}
	for _, rev := range reviewsRaw {
		avgRating += float64(rev.Rating)
	}
	if totalRawReviews > 0 {
		avgRating = avgRating / float64(totalRawReviews)
	}

	ratio := math.Pow(10, float64(1))
	RoundedAvgRating := math.Round(avgRating*ratio) / ratio

	res := &dto.GetReviewsRes{
		Limit:         qp.Limit,
		Page:          qp.Page,
		TotalPages:    totalPages,
		TotalReviews:  totalReviews,
		AverageRating: RoundedAvgRating,
		Reviews:       reviewsRes,
	}

	return res, nil
}

func (s *reviewService) CreateUpdateReview(userID uint, req *dto.CreateUpdateReview) (*model.Review, error) {

	tx := s.db.Begin()
	var err error
	defer helper.CommitOrRollback(tx, &err)

	_, err = s.reviewRepo.ValidateUserOrderItem(tx, userID, req.ProductID)
	if err != nil {
		return nil, err
	}

	var existingReview *model.Review

	existingReview, err = s.reviewRepo.FindReviewByProductIDAndSellerID(tx, userID, req.ProductID)

	newReview := model.Review{
		UserID:      userID,
		ProductID:   req.ProductID,
		Rating:      int(req.Rating),
		ImageURL:    req.ImageURL,
		ImageName:   req.ImageName,
		Description: req.Description,
	}
	var createdReview *model.Review

	if existingReview.ID == 0 {
		createdReview, err = s.reviewRepo.CreateReview(tx, &newReview)
		if err != nil {
			return nil, err
		}
	} else {
		createdReview, err = s.reviewRepo.UpdateReview(tx, existingReview.ID, &newReview)
		if err != nil {
			return nil, err
		}
	}
	return createdReview, nil
}

func (s *reviewService) UserReviewHistory(userID uint) ([]*model.Review, error) {

	tx := s.db.Begin()
	var err error
	defer helper.CommitOrRollback(tx, &err)

	var reviewHistory []*model.Review

	reviewHistory, err = s.reviewRepo.UserReviewHistory(tx, userID)
	if err != nil {
		return nil, err
	}
	return reviewHistory, nil
}
