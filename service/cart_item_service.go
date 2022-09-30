package service

import (
	"gorm.io/gorm"
	"seadeals-backend/dto"
	"seadeals-backend/helper"
	"seadeals-backend/model"
	"seadeals-backend/repository"
	"time"
)

type CartItemService interface {
	DeleteCartItem(orderItemID uint, userID uint) (*model.CartItem, error)
	AddToCart(userID uint, req *dto.AddToCartReq) (*model.CartItem, error)
	GetCartItems(query *repository.Query, userID uint) ([]*dto.CartItemRes, int64, int64, error)
}

type cartItemService struct {
	db                 *gorm.DB
	cartItemRepository repository.CartItemRepository
}

type CartItemServiceConfig struct {
	DB                 *gorm.DB
	CartItemRepository repository.CartItemRepository
}

func NewCartItemService(config *CartItemServiceConfig) CartItemService {
	return &cartItemService{
		db:                 config.DB,
		cartItemRepository: config.CartItemRepository,
	}
}

func (o *cartItemService) DeleteCartItem(orderItemID uint, userID uint) (*model.CartItem, error) {
	tx := o.db.Begin()
	var err error
	defer helper.CommitOrRollback(tx, &err)

	deleteOrder, err := o.cartItemRepository.DeleteCartItem(tx, orderItemID, userID)
	if err != nil {
		return nil, err
	}
	return deleteOrder, nil
}

func (o *cartItemService) AddToCart(userID uint, req *dto.AddToCartReq) (*model.CartItem, error) {
	tx := o.db.Begin()
	var err error
	defer helper.CommitOrRollback(tx, &err)

	cartItem := &model.CartItem{
		ProductVariantDetailID: req.ProductVariantDetailID,
		UserID:                 userID,
		Quantity:               req.Quantity,
	}
	addedItem, err := o.cartItemRepository.AddToCart(tx, cartItem)
	if err != nil {
		return nil, err
	}
	return addedItem, nil
}

func (o *cartItemService) GetCartItems(query *repository.Query, userID uint) ([]*dto.CartItemRes, int64, int64, error) {
	tx := o.db.Begin()
	var err error
	defer helper.CommitOrRollback(tx, &err)

	orderItems, totalPage, totalData, err := o.cartItemRepository.GetCartItem(tx, query, userID)
	if err != nil {
		return nil, 0, 0, err
	}

	var cartItems = make([]*dto.CartItemRes, 0)
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
			ID:           item.ID,
			Quantity:     item.Quantity,
			Subtotal:     subtotal,
			PricePerItem: pricePerItem,
			SellerID:     item.ProductVariantDetail.Product.SellerID,
			ProductName:  item.ProductVariantDetail.Product.Name,
		}
		cartItems = append(cartItems, cartItem)
	}

	return cartItems, totalPage, totalData, nil
}
