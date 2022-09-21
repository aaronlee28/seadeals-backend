package repository

import (
	"gorm.io/gorm"
	"seadeals-backend/apperror"
	"seadeals-backend/model"
	"strconv"
)

type OrderItemRepository interface {
	AddToCart(tx *gorm.DB, orderItem *model.OrderItem) (*model.OrderItem, error)
	DeleteCartItem(tx *gorm.DB, orderItemID uint, userID uint) (*model.OrderItem, error)
	GetOrderItem(tx *gorm.DB, query *Query, userID uint) ([]*model.OrderItem, int64, int64, error)
}

type orderItemRepository struct{}

func NewOrderItemRepository() OrderItemRepository {
	return &orderItemRepository{}
}

func (o *orderItemRepository) AddToCart(tx *gorm.DB, orderItem *model.OrderItem) (*model.OrderItem, error) {
	var existingOrderItem = &model.OrderItem{}
	result := tx.Where("user_id = ?", orderItem.UserID).Where("product_variant_detail_id = ?", orderItem.ProductVariantDetailID).Where("order_id IS NULL").First(&existingOrderItem)
	if result.Error == nil {
		existingOrderItem.Quantity += orderItem.Quantity
		result = tx.Updates(&existingOrderItem)
		if result.Error != nil {
			return nil, apperror.InternalServerError("Cannot update order item")
		}
		return existingOrderItem, nil
	}

	result = tx.Create(&orderItem)
	if result.Error != nil {
		return nil, apperror.InternalServerError("Cannot create order item")
	}

	return orderItem, nil
}

func (o *orderItemRepository) DeleteCartItem(tx *gorm.DB, orderItemID uint, userID uint) (*model.OrderItem, error) {
	var existingOrderItem = &model.OrderItem{ID: orderItemID}
	result := tx.First(&existingOrderItem)
	if result.Error != nil {
		return nil, apperror.NotFoundError("Cannot find order item")
	}

	if existingOrderItem.UserID != userID {
		return nil, apperror.UnauthorizedError("Cannot delete other user order item")
	}

	result = tx.Model(&existingOrderItem).Update("quantity", 0)
	if result.Error != nil {
		return nil, apperror.InternalServerError("Cannot delete order item")
	}
	return existingOrderItem, nil
}

func (o *orderItemRepository) GetOrderItem(tx *gorm.DB, query *Query, userID uint) ([]*model.OrderItem, int64, int64, error) {
	var orderItems []*model.OrderItem
	var count int64
	result := tx.Model(&model.OrderItem{})
	result = result.Where("user_id = ?", userID).Where("order_id IS NULL").Where("quantity != ?", 0).Count(&count)
	if result.Error != nil {
		return nil, 0, 0, apperror.InternalServerError("Cannot count order item")
	}

	limit, _ := strconv.Atoi(query.Limit)
	if limit != 0 {
		result = result.Limit(limit)
	}

	result = result.Preload("ProductVariantDetail").Preload("ProductVariantDetail.Product").Preload("ProductVariantDetail.Product.Promotion").Find(&orderItems)
	if result.Error != nil {
		return nil, 0, 0, apperror.NotFoundError("Cannot get order item")
	}

	totalPage := count / int64(limit)
	if count%int64(limit) != 0 {
		totalPage += 1
	}

	return orderItems, totalPage, count, nil
}
