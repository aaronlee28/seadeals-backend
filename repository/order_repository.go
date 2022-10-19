package repository

import (
	"fmt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"seadeals-backend/apperror"
	"seadeals-backend/db"
	"seadeals-backend/dto"
	"seadeals-backend/model"
	"time"
)

type OrderQuery struct {
	Filter string `json:"filter"`
	Page   int    `json:"page"`
	Limit  int    `json:"limit"`
}

type OrderRepository interface {
	GetOrderBySellerID(tx *gorm.DB, sellerID uint, query *OrderQuery) ([]*model.Order, int64, int64, error)
	GetOrderByUserID(tx *gorm.DB, userID uint, query *OrderQuery) ([]*model.Order, int64, int64, error)
	GetOrderDetailByID(tx *gorm.DB, orderID uint) (*model.Order, error)

	UpdateOrderStatus(tx *gorm.DB, orderID uint, status string) (*model.Order, error)
	CheckAndUpdateOnDelivery() []*model.Order
	CheckAndUpdateWaitingForSeller() []*model.Order
	RefundToWalletByUserID(userID uint, refundedAmount float64) *model.Wallet
	AddToWalletTransaction(walletID uint, refundAmount float64)
	GetOrderItemsByOrderID(orderID uint) []*model.OrderItem
	UpdateStockByProductVariantDetailID(pvdID uint, quantity uint)
	UpdateOrderStatusByTransID(tx *gorm.DB, transactionID uint, status string) ([]*model.Order, error)
	CheckAndUpdateOnOrderDelivered() []*model.Order
}

type orderRepository struct {
}

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
	result = result.Preload("Delivery.Courier")
	result = result.Preload("Seller")
	result = result.Preload("Complaint")
	result = result.Preload("Voucher")
	result = result.Preload("OrderItems.ProductVariantDetail.Product.Category")
	result = result.Preload("OrderItems.ProductVariantDetail.Product.Promotion")
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

func (o *orderRepository) GetOrderByUserID(tx *gorm.DB, userID uint, query *OrderQuery) ([]*model.Order, int64, int64, error) {
	var orders []*model.Order
	result := tx.Model(&orders).Where("user_id = ?", userID)
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
	result = result.Preload("Delivery.Courier")
	result = result.Preload("Seller")
	result = result.Preload("Complaint")
	result = result.Preload("Voucher")
	result = result.Preload("OrderItems.ProductVariantDetail.Product.Category")
	result = result.Preload("OrderItems.ProductVariantDetail.Product.Promotion")
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
	result := tx.Model(&order).Preload("OrderItems").Preload("Complaint.ComplaintPhotos").Preload("Transaction").First(&order)
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

func (o *orderRepository) CheckAndUpdateOnDelivery() []*model.Order {
	var order []*model.Order
	tx := db.Get().Begin()
	_ = tx.Clauses(clause.Returning{}).Where("status = ?", dto.OrderOnDelivery).Where("? >= updated_at at time zone 'UTC' + interval '2 day'", time.Now()).Find(&order).Update("status", dto.OrderDelivered)

	tx.Commit()
	return order

}

func (o *orderRepository) CheckAndUpdateOnOrderDelivered() []*model.Order {
	var order []*model.Order
	tx := db.Get().Begin()
	_ = tx.Clauses(clause.Returning{}).Where("status = ?", dto.OrderDelivered).Where("? >= updated_at at time zone 'UTC' + interval '2 day'", time.Now()).Find(&order).Update("status", dto.OrderDone)

	tx.Commit()
	return order

}

func (o *orderRepository) CheckAndUpdateWaitingForSeller() []*model.Order {
	tx := db.Get().Begin()
	var orders []*model.Order
	result := tx.Clauses(clause.Returning{}).Where("status = ?", dto.OrderWaitingSeller).Where("? >= updated_at at time zone 'UTC' + interval '3 day'", time.Now()).Find(&orders).Update("status", dto.OrderRefunded)
	if result.Error != nil {
		tx.Rollback()
		fmt.Println("error:", result.Error)
		return nil
	}
	tx.Commit()
	return orders
}

func (o *orderRepository) RefundToWalletByUserID(userID uint, refundedAmount float64) *model.Wallet {
	tx := db.Get().Begin()
	var wallet *model.Wallet
	result := tx.Clauses(clause.Returning{}).Where("id = ?", userID).First(&wallet).Update("balance", wallet.Balance+refundedAmount)
	if result.Error != nil {
		tx.Rollback()
		fmt.Println("error:", result.Error)
		return nil
	}
	tx.Commit()
	return wallet
}

func (o *orderRepository) AddToWalletTransaction(walletID uint, refundAmount float64) {
	tx := db.Get().Begin()
	walletTransaction := model.WalletTransaction{
		WalletID:      walletID,
		TransactionID: nil,
		Total:         refundAmount,
		PaymentMethod: dto.Wallet,
		PaymentType:   "credit",
		Description:   "refund",
	}
	result := tx.Create(&walletTransaction)
	if result.Error != nil {
		tx.Rollback()
		fmt.Println("error:", result.Error)
		return
	}
	tx.Commit()
	return
}

func (o *orderRepository) GetOrderItemsByOrderID(orderID uint) []*model.OrderItem {
	tx := db.Get().Begin()
	var orderItems []*model.OrderItem
	result := tx.Clauses(clause.Returning{}).Where("order_id = ?", orderID).Find(&orderItems)
	if result.Error != nil {
		tx.Rollback()
		fmt.Println("error:", result.Error)
		return nil
	}
	tx.Commit()
	return orderItems
}

func (o *orderRepository) UpdateStockByProductVariantDetailID(pvdID uint, quantity uint) {
	tx := db.Get().Begin()
	var pvd *model.ProductVariantDetail
	result := tx.Clauses(clause.Returning{}).Where("id = ?", pvdID).Find(&pvd).Update("stock", pvd.Stock+quantity)
	if result.Error != nil {
		tx.Rollback()
		fmt.Println("error:", result.Error)
		return
	}
	tx.Commit()
	return
}

func (o *orderRepository) UpdateOrderStatusByTransID(tx *gorm.DB, transactionID uint, status string) ([]*model.Order, error) {
	var orders []*model.Order
	result := tx.Model(&orders).Clauses(clause.Returning{}).Where("transaction_id = ?", transactionID).Update("status", status)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, apperror.BadRequestError("order doesn't exists")
		}
		return nil, apperror.InternalServerError("Cannot find order")
	}
	result = result.Model(&orders).Where("transaction_id = ?", transactionID).Preload("Delivery").Find(&orders)
	if result.Error != nil {
		return nil, apperror.InternalServerError("Cannot find order")
	}
	return orders, nil
}
