package service

import (
	"gorm.io/gorm"
	"seadeals-backend/apperror"
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
	productVarDetRepo  repository.ProductVariantDetailRepository
}

type CartItemServiceConfig struct {
	DB                 *gorm.DB
	CartItemRepository repository.CartItemRepository
	ProductVarDetRepo  repository.ProductVariantDetailRepository
}

func NewCartItemService(config *CartItemServiceConfig) CartItemService {
	return &cartItemService{
		db:                 config.DB,
		cartItemRepository: config.CartItemRepository,
		productVarDetRepo:  config.ProductVarDetRepo,
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

	productVarDet, err := o.productVarDetRepo.GetProductVariantDetailByID(tx, req.ProductVariantDetailID)
	if err != nil {
		return nil, err
	}

	if productVarDet.Product.Seller.UserID == userID {
		err = apperror.BadRequestError("Cannot buy your own product")
		return nil, err
	}

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

		var imageURL string
		if len(item.ProductVariantDetail.Product.ProductPhotos) > 0 {
			imageURL = item.ProductVariantDetail.Product.ProductPhotos[0].PhotoURL
		}
		cartItem := &dto.CartItemRes{
			ID:           item.ID,
			Quantity:     item.Quantity,
			Subtotal:     subtotal,
			PricePerItem: pricePerItem,
			SellerID:     item.ProductVariantDetail.Product.SellerID,
			ImageURL:     imageURL,
			SellerName:   item.ProductVariantDetail.Product.Seller.Name,
			ProductName:  item.ProductVariantDetail.Product.Name,
		}
		cartItems = append(cartItems, cartItem)
	}

	return cartItems, totalPage, totalData, nil
}
