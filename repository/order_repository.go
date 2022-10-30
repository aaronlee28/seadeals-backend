package repository

import (
	"fmt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"seadeals-backend/apperror"
	"seadeals-backend/config"
	"seadeals-backend/db"
	"seadeals-backend/dto"
	"seadeals-backend/model"
	"strings"
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
	// GetOrderDetailForReceipt function below is just to prevent heavy loading when fetching data
	GetOrderDetailForReceipt(tx *gorm.DB, orderID uint) (*model.Order, error)
	// GetOrderDetailForThermal function below is just to prevent heavy loading when fetching data
	GetOrderDetailForThermal(tx *gorm.DB, orderID uint) (*model.Order, error)

	UpdateOrderStatus(tx *gorm.DB, orderID uint, status string) (*model.Order, error)
	FindAndUpdateWaitingForSellerToRefunded() []*model.Order
	RefundToWalletByUserID(userID uint, refundedAmount float64) *model.Wallet
	AddToWalletTransaction(walletID uint, refundAmount float64)
	GetOrderItemsByOrderID(orderID uint) []*model.OrderItem
	UpdateStockByProductVariantDetailID(pvdID uint, quantity uint)
	UpdateOrderStatusByTransID(tx *gorm.DB, transactionID uint, status string) ([]*model.Order, error)
	FindAndUpdateDeliveredOrderToDone() []*model.Order
	GetOrderByID(tx *gorm.DB, userID uint, orderID uint) (*model.Order, error)
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
		var filters []string
		filters = strings.Split(query.Filter, ",")
		result = result.Where("status IN (?)", filters)
	}
	result.Where("status NOT LIKE ?", dto.OrderWaitingPayment)

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
	result = result.Preload("User")
	result = result.Preload("Delivery.Courier")
	result = result.Preload("Seller")
	result = result.Preload("Complaint.ComplaintPhotos")
	result = result.Preload("Voucher")
	result = result.Preload("OrderItems.ProductVariantDetail.ProductVariant1")
	result = result.Preload("OrderItems.ProductVariantDetail.ProductVariant2")
	result = result.Preload("OrderItems.ProductVariantDetail.Product.ProductPhotos")
	result = result.Preload("OrderItems.ProductVariantDetail.Product.Category")
	result = result.Preload("OrderItems.ProductVariantDetail.Product.Promotion")
	result = result.Preload("Transaction")
	result = result.Order("updated_at desc").Order("id").Find(&orders)
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
		var filters []string
		filters = strings.Split(query.Filter, ",")
		result = result.Where("status IN (?)", filters)
	}
	result.Where("status NOT LIKE ?", dto.OrderWaitingPayment)

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
	result = result.Preload("User")
	result = result.Preload("Delivery.Courier")
	result = result.Preload("Seller")
	result = result.Preload("Complaint.ComplaintPhotos")
	result = result.Preload("Voucher")
	result = result.Preload("OrderItems.ProductVariantDetail.ProductVariant1")
	result = result.Preload("OrderItems.ProductVariantDetail.ProductVariant2")
	result = result.Preload("OrderItems.ProductVariantDetail.Product.ProductPhotos")
	result = result.Preload("OrderItems.ProductVariantDetail.Product.Category")
	result = result.Preload("OrderItems.ProductVariantDetail.Product.Promotion")
	result = result.Preload("OrderItems.ProductVariantDetail.Product.Review", "user_id = ?", userID)
	result = result.Preload("Transaction")
	result = result.Order("updated_at desc").Order("id").Find(&orders)
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
	result := tx.Model(&order).Preload("OrderItems.ProductVariantDetail").Preload("Complaint.ComplaintPhotos").Preload("Transaction").First(&order)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, apperror.BadRequestError("order doesn't exists")
		}
		return nil, apperror.InternalServerError("Cannot find order")
	}
	return order, nil
}

func (o *orderRepository) GetOrderDetailForReceipt(tx *gorm.DB, orderID uint) (*model.Order, error) {
	var order = &model.Order{}
	order.ID = orderID
	result := tx.Model(&order).Preload("OrderItems.ProductVariantDetail.Product.ProductDetail")
	result = result.Preload("Transaction.Orders.OrderItems.ProductVariantDetail.Product")
	result = result.Preload("Transaction.Orders.Seller.Address")
	result = result.Preload("Transaction.Orders.Delivery")
	result = result.Preload("Transaction.Voucher")
	result = result.Preload("Transaction.Orders.OrderItems.ProductVariantDetail.ProductVariant1")
	result = result.Preload("Transaction.Orders.OrderItems.ProductVariantDetail.ProductVariant2")
	result = result.Preload("Voucher")
	result = result.Preload("Delivery.Courier")
	result = result.Preload("Transaction.Voucher")
	result = result.Preload("Seller.Address")
	result = result.Preload("User")
	result = result.Preload("Complaint.ComplaintPhotos").Preload("Transaction")
	result = result.First(&order)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, apperror.BadRequestError("order doesn't exists")
		}
		return nil, apperror.InternalServerError("Cannot find order")
	}
	return order, nil
}

func (o *orderRepository) GetOrderDetailForThermal(tx *gorm.DB, orderID uint) (*model.Order, error) {
	var order = &model.Order{}
	order.ID = orderID
	result := tx.Model(&order).Preload("OrderItems.ProductVariantDetail.Product.ProductDetail")
	result = result.Preload("Transaction.Orders.OrderItems.ProductVariantDetail.ProductVariant1")
	result = result.Preload("Transaction.Orders.OrderItems.ProductVariantDetail.ProductVariant2")
	result = result.Preload("Delivery.Courier")
	result = result.Preload("Seller.Address")
	result = result.Preload("User")
	result = result.First(&order)
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

func (o *orderRepository) FindAndUpdateDeliveredOrderToDone() []*model.Order {
	var order []*model.Order
	tx := db.Get().Begin()
	_ = tx.Model(&order).Clauses(clause.Returning{}).Where("status = ?", dto.OrderDelivered).Where("? >= updated_at at time zone '"+config.Config.TZ+"' + interval '"+config.Config.Interval.DeliveredOrderToDone+"'", time.Now()).Update("status", dto.OrderDone)

	tx.Commit()
	return order

}

func (o *orderRepository) FindAndUpdateWaitingForSellerToRefunded() []*model.Order {
	tx := db.Get().Begin()
	var orders []*model.Order
	result := tx.Clauses(clause.Returning{}).Where("status = ?", dto.OrderWaitingSeller).Where("? >= updated_at at time zone '"+config.Config.TZ+"' + interval '"+config.Config.Interval.WaitingForSellerToRefund+"'", time.Now()).Find(&orders).Update("status", dto.OrderRefunded)
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
	result := tx.Clauses(clause.Returning{}).Where("user_id = ?", userID).First(&wallet).Update("balance", wallet.Balance+refundedAmount)
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
		PaymentType:   "CREDIT",
		Description:   "refund",
	}
	result := tx.Create(&walletTransaction)
	if result.Error != nil {
		tx.Rollback()
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
	result = result.Model(&orders).Where("transaction_id = ?", transactionID).Preload("OrderItems.ProductVariantDetail").Preload("Delivery").Find(&orders)
	if result.Error != nil {
		return nil, apperror.InternalServerError("Cannot find order")
	}
	return orders, nil
}

func (o *orderRepository) GetOrderByID(tx *gorm.DB, userID uint, orderID uint) (*model.Order, error) {
	var order *model.Order
	result := tx.Model(&order).Where("id = ? AND user_id = ?", orderID, userID)
	result = result.Preload("Delivery.DeliveryActivity")
	result = result.Preload("User")
	result = result.Preload("Delivery.Courier")
	result = result.Preload("Seller")
	result = result.Preload("Complaint")
	result = result.Preload("Voucher")
	result = result.Preload("OrderItems.ProductVariantDetail.ProductVariant1")
	result = result.Preload("OrderItems.ProductVariantDetail.ProductVariant2")
	result = result.Preload("OrderItems.ProductVariantDetail.Product.ProductPhotos")
	result = result.Preload("OrderItems.ProductVariantDetail.Product.Category")
	result = result.Preload("OrderItems.ProductVariantDetail.Product.Promotion")
	result = result.Preload("OrderItems.ProductVariantDetail.Product.Review", "user_id = ?", userID)
	result = result.Preload("Transaction")
	result = result.Order("updated_at desc").Order("id").First(&order)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, apperror.BadRequestError("order doesn't exists")
		}
		return nil, apperror.InternalServerError("Cannot find order")
	}
	return order, nil
}
