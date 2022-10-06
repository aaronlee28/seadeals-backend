package repository

import (
	"gorm.io/gorm"
	"seadeals-backend/apperror"
	"seadeals-backend/model"
	"strconv"
)

type CartItemRepository interface {
	AddToCart(tx *gorm.DB, orderItem *model.CartItem) (*model.CartItem, error)
	DeleteCartItem(tx *gorm.DB, orderItemID uint, userID uint) (*model.CartItem, error)
	GetCartItem(tx *gorm.DB, query *Query, userID uint) ([]*model.CartItem, int64, int64, error)
}

type cartItemRepository struct{}

func NewCartItemRepository() CartItemRepository {
	return &cartItemRepository{}
}

func (o *cartItemRepository) AddToCart(tx *gorm.DB, cartItem *model.CartItem) (*model.CartItem, error) {
	var existingCartItem = &model.CartItem{}
	result := tx.Where("user_id = ?", cartItem.UserID).Where("product_variant_detail_id = ?", cartItem.ProductVariantDetailID).First(&existingCartItem)
	if result.Error == nil {
		existingCartItem.Quantity += cartItem.Quantity
		result = tx.Updates(&existingCartItem)
		if result.Error != nil {
			return nil, apperror.InternalServerError("Cannot update order item")
		}
		return existingCartItem, nil
	}

	result = tx.Create(&cartItem)
	if result.Error != nil {
		return nil, apperror.InternalServerError("Cannot create cart item")
	}

	return cartItem, nil
}

func (o *cartItemRepository) DeleteCartItem(tx *gorm.DB, cartItemID uint, userID uint) (*model.CartItem, error) {
	var existingCartItem = &model.CartItem{ID: cartItemID}
	result := tx.First(&existingCartItem)
	if result.Error != nil {
		return nil, apperror.NotFoundError("Cannot find cart item")
	}

	if existingCartItem.UserID != userID {
		return nil, apperror.UnauthorizedError("Cannot delete other user cart item")
	}

	result = tx.Model(&existingCartItem).Update("quantity", 0)
	if result.Error != nil {
		return nil, apperror.InternalServerError("Cannot delete cart item")
	}
	return existingCartItem, nil
}

func (o *cartItemRepository) GetCartItem(tx *gorm.DB, query *Query, userID uint) ([]*model.CartItem, int64, int64, error) {
	var cartItems []*model.CartItem
	var count int64
	result := tx.Model(&model.CartItem{})
	result = result.Where("user_id = ?", userID).Where("quantity != ?", 0).Count(&count)
	if result.Error != nil {
		return nil, 0, 0, apperror.InternalServerError("Cannot count cart item")
	}

	limit, _ := strconv.Atoi(query.Limit)
	if limit != 0 {
		result = result.Limit(limit)
	}

	result = result.Preload("ProductVariantDetail").Preload("ProductVariantDetail.Product.ProductPhotos").Preload("ProductVariantDetail.Product.Seller").Preload("ProductVariantDetail.Product.Promotion").Find(&cartItems)
	if result.Error != nil {
		return nil, 0, 0, apperror.NotFoundError("Cannot get cart item")
	}

	totalPage := int64(1)
	if limit != 0 {
		totalPage = count / int64(limit)
		if count%int64(limit) != 0 {
			totalPage += 1
		}
	}

	return cartItems, totalPage, count, nil
}
