package repository

import (
	"fmt"
	"gorm.io/gorm"
	"seadeals-backend/apperror"
	"seadeals-backend/model"
)

type ReviewRepository interface {
	GetReviewsAvgAndCountBySellerID(tx *gorm.DB, sellerID uint) (float64, int64, error)
}

type reviewRepository struct{}

func NewReviewRepository() ReviewRepository {
	return &reviewRepository{}
}

func (r *reviewRepository) GetReviewsAvgAndCountBySellerID(tx *gorm.DB, sellerID uint) (float64, int64, error) {
	var average float64
	var totalReview int64
	result := tx.Model(&model.Review{}).Joins("Product", tx.Where(&model.Product{SellerID: sellerID})).Count(&totalReview)
	if result.Error != nil {
		fmt.Println(result.Error)
		return 0, 0, apperror.InternalServerError("Cannot count total review")
	}

	result = result.Select("avg(rating) as total").Find(&average)
	if result.Error != nil {
		return 0, 0, apperror.InternalServerError("Cannot count average review")
	}

	return average, totalReview, nil
}
