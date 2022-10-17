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
	UpdateCart(userID uint, req *dto.UpdateCartItemReq) (*model.CartItem, error)
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

func (c *cartItemService) DeleteCartItem(orderItemID uint, userID uint) (*model.CartItem, error) {
	tx := c.db.Begin()
	var err error
	defer helper.CommitOrRollback(tx, &err)

	deleteOrder, err := c.cartItemRepository.DeleteCartItem(tx, orderItemID, userID)
	if err != nil {
		return nil, err
	}
	return deleteOrder, nil
}

func (c *cartItemService) AddToCart(userID uint, req *dto.AddToCartReq) (*model.CartItem, error) {
	tx := c.db.Begin()
	var err error
	defer helper.CommitOrRollback(tx, &err)

	productVarDet, err := c.productVarDetRepo.GetProductVariantDetailByID(tx, req.ProductVariantDetailID)
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
	addedItem, err := c.cartItemRepository.AddToCart(tx, cartItem)
	if err != nil {
		return nil, err
	}

	if productVarDet.Stock < addedItem.Quantity || (productVarDet.Product.MaxQuantity < addedItem.Quantity && productVarDet.Product.MaxQuantity != 0) {
		err = apperror.BadRequestError("Kuantitas pembelian melebih stock atau maximum pembelian")
		return nil, err
	}

	if productVarDet.Product.MinQuantity > addedItem.Quantity && productVarDet.Product.MinQuantity != 0 {
		err = apperror.BadRequestError("Kuantitas pembelian Kurang dari minimum pembelian")
		return nil, err
	}

	return addedItem, nil
}

func (c *cartItemService) UpdateCart(userID uint, req *dto.UpdateCartItemReq) (*model.CartItem, error) {
	tx := c.db.Begin()
	var err error
	defer helper.CommitOrRollback(tx, &err)

	var cartItem *model.CartItem
	cartItem, err = c.cartItemRepository.UpdateCart(tx, req, userID)
	if err != nil {
		return nil, err
	}

	productVarDet, err := c.productVarDetRepo.GetProductVariantDetailByID(tx, cartItem.ProductVariantDetailID)
	if err != nil {
		return nil, err
	}

	if productVarDet.Stock < cartItem.Quantity || (productVarDet.Product.MaxQuantity < cartItem.Quantity && productVarDet.Product.MaxQuantity != 0) {
		err = apperror.BadRequestError("Kuantitas pembelian melebih stock atau maximum pembelian")
		return nil, err
	}

	if productVarDet.Product.MinQuantity > cartItem.Quantity && productVarDet.Product.MinQuantity != 0 {
		err = apperror.BadRequestError("Kuantitas pembelian Kurang dari minimum pembelian")
		return nil, err
	}

	return cartItem, err
}

func (c *cartItemService) GetCartItems(query *repository.Query, userID uint) ([]*dto.CartItemRes, int64, int64, error) {
	tx := c.db.Begin()
	var err error
	defer helper.CommitOrRollback(tx, &err)

	orderItems, totalPage, totalData, err := c.cartItemRepository.GetCartItem(tx, query, userID)
	if err != nil {
		return nil, 0, 0, err
	}

	var cartItems = make([]*dto.CartItemRes, 0)
	for _, item := range orderItems {
		subtotal := float64(item.Quantity) * item.ProductVariantDetail.Price
		now := time.Now()
		fullPrice := item.ProductVariantDetail.Price
		promotion := item.ProductVariantDetail.Product.Promotion
		currentPricePerItem := item.ProductVariantDetail.Price
		var discountNominal *float64
		var discountPercent *int
		if promotion != nil && now.After(promotion.StartDate) && now.Before(promotion.EndDate) && promotion.Quota >= item.Quantity {
			if promotion.AmountType == "percent" {
				subtotal = (100 - promotion.Amount) / 100 * subtotal
			} else {
				discountNominal = &promotion.Amount
				currentPricePerItem -= promotion.Amount
				percent := 100 - int((currentPricePerItem/fullPrice)*100)
				discountPercent = &percent
				subtotal = float64(item.Quantity) * (currentPricePerItem)
			}
		}

		var imageURL string
		if len(item.ProductVariantDetail.Product.ProductPhotos) > 0 {
			imageURL = item.ProductVariantDetail.Product.ProductPhotos[0].PhotoURL
		}

		var variantDetail string
		if item.ProductVariantDetail.ProductVariant1 != nil {
			variantDetail += *item.ProductVariantDetail.Variant1Value
		}
		if item.ProductVariantDetail.ProductVariant2 != nil {
			variantDetail += ", " + *item.ProductVariantDetail.Variant2Value
		}
		cartItem := &dto.CartItemRes{
			ID:                  item.ID,
			Quantity:            item.Quantity,
			ProductVariant:      variantDetail,
			MinQuantity:         item.ProductVariantDetail.Product.MinQuantity,
			MaxQuantity:         item.ProductVariantDetail.Product.MaxQuantity,
			Stock:               item.ProductVariantDetail.Stock,
			ProductSlug:         item.ProductVariantDetail.Product.Slug,
			DiscountPercent:     discountPercent,
			DiscountNominal:     discountNominal,
			PriceBeforeDiscount: fullPrice,
			PricePerItem:        currentPricePerItem,
			SellerID:            item.ProductVariantDetail.Product.SellerID,
			SellerName:          item.ProductVariantDetail.Product.Seller.Name,
			ImageURL:            imageURL,
			Subtotal:            subtotal,
			ProductName:         item.ProductVariantDetail.Product.Name,
		}
		cartItems = append(cartItems, cartItem)
	}

	return cartItems, totalPage, totalData, nil
}
