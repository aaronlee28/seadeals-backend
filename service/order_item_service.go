package service

import (
	"gorm.io/gorm"
	"seadeals-backend/dto"
	"seadeals-backend/model"
	"seadeals-backend/repository"
	"time"
)

type OrderItemService interface {
	DeleteOrderItem(orderItemID uint, userID uint) (*model.OrderItem, error)
	AddToCart(req *dto.AddToCartReq) (*model.OrderItem, error)
	GetOrderItem(query *repository.Query, userID uint) ([]*dto.CartItemRes, int64, int64, error)
}

type orderItemService struct {
	db                  *gorm.DB
	orderItemRepository repository.OrderItemRepository
}

type OrderItemServiceConfig struct {
	DB                  *gorm.DB
	OrderItemRepository repository.OrderItemRepository
}

func NewOrderItemService(config *OrderItemServiceConfig) OrderItemService {
	return &orderItemService{
		db:                  config.DB,
		orderItemRepository: config.OrderItemRepository,
	}
}

func (o *orderItemService) DeleteOrderItem(orderItemID uint, userID uint) (*model.OrderItem, error) {
	tx := o.db.Begin()
	deleteOrder, err := o.orderItemRepository.DeleteCartItem(tx, orderItemID, userID)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	tx.Commit()
	return deleteOrder, nil
}

func (o *orderItemService) AddToCart(req *dto.AddToCartReq) (*model.OrderItem, error) {
	tx := o.db.Begin()
	orderItem := &model.OrderItem{
		ProductVariantDetailID: req.ProductVariantDetailID,
		UserID:                 req.UserID,
		Quantity:               req.Quantity,
	}
	addedItem, err := o.orderItemRepository.AddToCart(tx, orderItem)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	tx.Commit()
	return addedItem, nil
}

func (o *orderItemService) GetOrderItem(query *repository.Query, userID uint) ([]*dto.CartItemRes, int64, int64, error) {
	tx := o.db.Begin()
	orderItems, totalPage, totalData, err := o.orderItemRepository.GetOrderItem(tx, query, userID)
	if err != nil {
		tx.Rollback()
		return nil, 0, 0, err
	}

	var cartItems []*dto.CartItemRes
	for _, item := range orderItems {
		subtotal := float64(item.Quantity) * item.ProductVariantDetail.Price
		now := time.Now()
		promotion := item.ProductVariantDetail.Product.Promotion
		pricePerItem := item.ProductVariantDetail.Price
		if promotion != nil && now.After(promotion.StartDate) && now.Before(promotion.EndDate) {
			if promotion.AmountType == "percent" {
				subtotal = (100 - promotion.Amount) / 100 * subtotal
			} else {
				pricePerItem -= promotion.Amount
				subtotal = float64(item.Quantity) * (pricePerItem)
			}
		}
		cartItem := &dto.CartItemRes{
			Quantity:     item.Quantity,
			Subtotal:     subtotal,
			PricePerItem: pricePerItem,
			SellerID:     item.ProductVariantDetail.Product.SellerID,
			ProductName:  item.ProductVariantDetail.Product.Name,
		}
		cartItems = append(cartItems, cartItem)
	}

	tx.Commit()
	return cartItems, totalPage, totalData, nil
}
