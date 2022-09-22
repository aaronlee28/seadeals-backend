package repository

import (
	"gorm.io/gorm"
	"seadeals-backend/apperror"
	"seadeals-backend/dto"
)

type ProductVariantDetailRepository interface {
	GetProductsBySellerID(tx *gorm.DB, query *dto.SellerProductSearchQuery, sellerID uint) ([]*dto.SellerProductsCustomTable, int64, int64, error)
}

type productVariantDetailRepository struct{}

func NewProductVariantDetailRepository() ProductVariantDetailRepository {
	return &productVariantDetailRepository{}
}

func (p *productVariantDetailRepository) GetProductsBySellerID(tx *gorm.DB, query *dto.SellerProductSearchQuery, sellerID uint) ([]*dto.SellerProductsCustomTable, int64, int64, error) {
	var products []*dto.SellerProductsCustomTable

	result := tx.Model(&dto.SellerProductsCustomTable{})
	result = result.Select("min(price) as min, max(price) as max, product_id, products.seller_id")
	result = result.Where("products.seller_id = ?", sellerID).Group("product_id").Group("products.seller_id").Joins("FULL JOIN products on products.id = product_id")
	if query.Search != "" {
		result = result.Where("name ILIKE ?", query.Search)
	}

	orderByString := query.SortBy
	if query.SortBy == "price" {
		orderByString = "min"
	} else {
		if query.SortBy == "" {
			orderByString = "products." + "sold_count"
			result.Select("min(price), max(price), product_id, seller_id, products.sold_count")
			result.Group("products.sold_count")
		} else {
			orderByString = "products." + query.SortBy
			result.Select("min(price), max(price), product_id, seller_id, products." + query.SortBy)
			result.Group("products." + query.SortBy)
		}
	}

	if query.SortBy == "" {
		if query.Sort != "asc" {
			orderByString += " desc"
		}
	} else {
		if query.Sort == "desc" {
			orderByString += " desc"
		}
	}

	var totalData int64
	result = result.Order(orderByString).Order("product_id")
	table := tx.Table("(?) as s1", result).Where("min >= ?", query.MinPrice).Where("min <= ?", query.MaxPrice)
	tx.Table("(?) as s2", table).Count(&totalData)
	if result.Error != nil {
		return nil, 0, 0, apperror.InternalServerError("cannot fetch products count")
	}

	limit := 20
	if query.Limit != 0 {
		limit = query.Limit
		table = table.Limit(limit)
	}
	if query.Page != 0 {
		table = table.Offset((query.Page - 1) * limit)
	}

	table = table.Unscoped()
	table = table.Preload("Product.ProductPhotos").Preload("Product.Seller.Address.SubDistrict.District.City")
	table = table.Find(&products)
	if table.Error != nil {
		return nil, 0, 0, apperror.InternalServerError("cannot fetch products")
	}

	totalPage := totalData / int64(limit)
	if totalData%int64(limit) != 0 {
		totalPage += 1
	}

	return products, totalPage, totalData, nil
}
