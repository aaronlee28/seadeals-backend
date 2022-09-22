package repository

import (
	"fmt"
	"gorm.io/gorm"
	"seadeals-backend/apperror"
	"seadeals-backend/model"
)

type ReviewRepository interface {
	GetReviewsAvgAndCountBySellerID(tx *gorm.DB, sellerID uint) (float64, int64, error)
	GetReviewsAvgAndCountByProductID(tx *gorm.DB, productID uint) (float64, int64, error)
	FindReviewByProductID(tx *gorm.DB, productID uint, qp *model.ReviewQueryParam) ([]*model.Review, error)
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

func (r *reviewRepository) GetReviewsAvgAndCountByProductID(tx *gorm.DB, productID uint) (float64, int64, error) {
	var average float64
	var totalReview int64
	result := tx.Model(&model.Review{}).Where("product_id = ?", productID).Count(&totalReview)
	if result.Error != nil {
		return 0, 0, apperror.InternalServerError("Cannot count total review")
	}

	if totalReview == 0 {
		return 0, 0, nil
	}

	result = result.Select("avg(rating) as total").Find(&average)
	if result.Error != nil {
		return 0, 0, apperror.InternalServerError("Cannot count average review")
	}

	return average, totalReview, nil
}

func (r *reviewRepository) FindReviewByProductID(tx *gorm.DB, productID uint, qp *model.ReviewQueryParam) ([]*model.Review, error) {
	var reviews []*model.Review
	orderStmt := fmt.Sprintf("%s %s", qp.SortBy, qp.Sort)

	result := tx.Where("product_id = ?", productID).Order(orderStmt).Preload("User").Find(&reviews)
	if result.Error != nil {
		return nil, result.Error
	}

	return reviews, nil
}
