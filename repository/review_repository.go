package repository

import (
	"fmt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"seadeals-backend/apperror"
	"seadeals-backend/model"
)

type ReviewRepository interface {
	GetReviewsAvgAndCountBySellerID(tx *gorm.DB, sellerID uint) (float64, int64, error)
	GetReviewsAvgAndCountByProductID(tx *gorm.DB, productID uint) (float64, int64, error)
	FindReviewByProductID(tx *gorm.DB, productID uint, qp *model.ReviewQueryParam) ([]*model.Review, error)
	FindReviewByProductIDAndSellerID(tx *gorm.DB, userID uint, productID uint) (*model.Review, error)
	ValidateUserOrderItem(tx *gorm.DB, userID uint, productID uint) (*model.OrderItem, error)
	CreateReview(tx *gorm.DB, req *model.Review) (*model.Review, error)
	UpdateReview(tx *gorm.DB, reviewID uint, req *model.Review) (*model.Review, error)
	UserReviewHistory(tx *gorm.DB, userID uint) ([]*model.Review, error)
	FindReviewByProductIDNoQuery(tx *gorm.DB, productID uint) ([]*model.Review, error)
	FindReviewByProductIDNoLimit(tx *gorm.DB, productID uint, qp *model.ReviewQueryParam) ([]*model.Review, error)
}

type reviewRepository struct{}

func NewReviewRepository() ReviewRepository {
	return &reviewRepository{}
}

func (r *reviewRepository) GetReviewsAvgAndCountBySellerID(tx *gorm.DB, sellerID uint) (float64, int64, error) {
	var average *float64
	var zero float64
	var totalReview int64
	result := tx.Model(&model.Review{}).Joins("JOIN products ON products.id = reviews.product_id").Where("products.seller_id = ? AND products.deleted_at IS NULL", sellerID).Count(&totalReview)
	if result.Error != nil {
		return 0, 0, apperror.InternalServerError("Cannot count total review")
	}

	result = result.Select("avg(rating) as average").Scan(&average)
	if result.Error != nil {
		return 0, 0, apperror.InternalServerError("Cannot count average review")
	}
	if average == nil {
		average = &zero
	}

	return *average, totalReview, nil
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
	offset := (qp.Page - 1) * qp.Limit
	orderStmt := fmt.Sprintf("%s %s", qp.SortBy, qp.Sort)

	queryDB := tx
	if qp.Rating != 0 {
		queryDB = queryDB.Where("rating = ?", qp.Rating)
	}
	if qp.WithImageOnly == true {
		queryDB = queryDB.Where("image_url IS NOT NULL")
	}
	if qp.WithDescriptionOnly == true {
		queryDB = queryDB.Where("description IS NOT NULL")
	}

	result := queryDB.Limit(int(qp.Limit)).Offset(int(offset)).Where("product_id = ?", productID).Order(orderStmt).Preload("User").Find(&reviews)

	return reviews, result.Error
}

func (r *reviewRepository) FindReviewByProductIDAndSellerID(tx *gorm.DB, userID uint, productID uint) (*model.Review, error) {
	var review *model.Review
	result := tx.Clauses(clause.Returning{}).Where("user_id = ?", userID).Where("product_id = ?", productID).First(&review)
	if result.Error != gorm.ErrRecordNotFound && result.Error != nil {

		return nil, result.Error
	}

	return review, nil
}

func (r *reviewRepository) ValidateUserOrderItem(tx *gorm.DB, userID uint, productID uint) (*model.OrderItem, error) {
	var order *model.OrderItem
	//var productVariantDetail *model.ProductVariantDetail
	result := tx.Clauses(clause.Returning{}).Preload("ProductVariantDetail", "product_id = ?", productID).Where("user_id = ?", userID).First(&order)
	//result2 := tx.Clauses(clause.Returning{}).Where("product_id = ?", productID).

	return order, result.Error
}

func (r *reviewRepository) CreateReview(tx *gorm.DB, req *model.Review) (*model.Review, error) {
	result := tx.Clauses(clause.Returning{}).Create(&req)
	if result.Error != nil {
		return nil, apperror.InternalServerError("Cannot create review")
	}
	return req, result.Error
}

func (r *reviewRepository) UpdateReview(tx *gorm.DB, reviewID uint, req *model.Review) (*model.Review, error) {
	result := tx.Clauses(clause.Returning{}).Where("id = ?", reviewID).Updates(&req)
	if result.Error != nil {
		return nil, apperror.InternalServerError("Cannot update review")
	}
	return req, result.Error
}

func (r *reviewRepository) UserReviewHistory(tx *gorm.DB, userID uint) ([]*model.Review, error) {
	var reviewHistory []*model.Review
	result := tx.Clauses(clause.Returning{}).Where("user_id = ?", userID).Find(&reviewHistory)
	if result.Error != nil {
		return nil, apperror.InternalServerError("Cannot find transaction")
	}
	return reviewHistory, result.Error
}

func (r *reviewRepository) FindReviewByProductIDNoQuery(tx *gorm.DB, productID uint) ([]*model.Review, error) {
	var reviews []*model.Review
	result := tx.Clauses(clause.Returning{}).Where("product_id = ?", productID).Find(&reviews)
	if result.Error != nil {
		return nil, apperror.InternalServerError("Cannot find reviews")
	}
	return reviews, nil
}

func (r *reviewRepository) FindReviewByProductIDNoLimit(tx *gorm.DB, productID uint, qp *model.ReviewQueryParam) ([]*model.Review, error) {
	var reviews []*model.Review

	orderStmt := fmt.Sprintf("%s %s", qp.SortBy, qp.Sort)

	queryDB := tx
	if qp.Rating != 0 {
		queryDB = queryDB.Where("rating = ?", qp.Rating)
	}
	if qp.WithImageOnly == true {
		queryDB = queryDB.Where("image_url IS NOT NULL")
	}
	if qp.WithDescriptionOnly == true {
		queryDB = queryDB.Where("description IS NOT NULL")
	}

	result := queryDB.Where("product_id = ?", productID).Order(orderStmt).Preload("User").Find(&reviews)

	return reviews, result.Error
}
