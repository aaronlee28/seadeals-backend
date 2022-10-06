package repository

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"seadeals-backend/apperror"
	"seadeals-backend/model"
)

type OrderQuery struct {
	Filter string `json:"filter"`
	Page   int    `json:"page"`
	Limit  int    `json:"limit"`
}

type OrderRepository interface {
	GetOrderBySellerID(tx *gorm.DB, sellerID uint, query *OrderQuery) ([]*model.Order, int64, int64, error)
	GetOrderDetailByID(tx *gorm.DB, orderID uint) (*model.Order, error)

	UpdateOrderStatus(tx *gorm.DB, orderID uint, status string) (*model.Order, error)
}

type orderRepository struct{}

func NewOrderRepo() OrderRepository {
	return &orderRepository{}
}

func (o *orderRepository) GetOrderBySellerID(tx *gorm.DB, sellerID uint, query *OrderQuery) ([]*model.Order, int64, int64, error) {
	var orders []*model.Order
	result := tx.Model(&orders).Where("seller_id = ?", sellerID)
	if query.Filter != "" {
		result = result.Where("status ILIKE ?", query.Filter)
	}

	var totalData int64
	table := result.Count(&totalData)
	if table.Error != nil {
		return nil, 0, 0, apperror.InternalServerError("Cannot count order")
	}

	limit := 0
	if query.Limit != 0 {
		limit = query.Limit
	}
	result = result.Limit(limit)
	if query.Page != 0 {
		result = result.Offset((query.Page - 1) * limit)
	}

	result = result.Preload("Delivery.DeliveryActivity")
	result = result.Preload("OrderItems.ProductVariantDetail.Product")
	result = result.Preload("Transaction")
	result = result.Order("created_at desc").Order("id").Find(&orders)
	if result.Error != nil {
		return nil, 0, 0, apperror.InternalServerError("Cannot find order")
	}

	totalPage := totalData / int64(limit)
	if totalData%int64(limit) != 0 {
		totalPage += 1
	}

	return orders, totalPage, totalData, nil
}

func (o *orderRepository) GetOrderDetailByID(tx *gorm.DB, orderID uint) (*model.Order, error) {
	var order = &model.Order{}
	order.ID = orderID
	result := tx.Model(&order).Preload("OrderItems").Preload("Transaction").First(&order)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, apperror.BadRequestError("order doesn't exists")
		}
		return nil, apperror.InternalServerError("Cannot find order")
	}
	return order, nil
}

func (o *orderRepository) UpdateOrderStatus(tx *gorm.DB, orderID uint, status string) (*model.Order, error) {
	var order = &model.Order{}
	order.ID = orderID
	result := tx.Model(&order).Clauses(clause.Returning{}).Update("status", status)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, apperror.BadRequestError("order doesn't exists")
		}
		return nil, apperror.InternalServerError("Cannot find order")
	}
	return order, nil
}
