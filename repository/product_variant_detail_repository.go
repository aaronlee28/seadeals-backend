package repository

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"seadeals-backend/apperror"
	"seadeals-backend/dto"
	"seadeals-backend/model"
	"strconv"
	"strings"
)

type ProductVariantDetailRepository interface {
	GetProductsBySellerID(tx *gorm.DB, query *dto.SellerProductSearchQuery, sellerID uint) ([]*dto.SellerProductsCustomTable, int64, int64, error)
	GetProductsByCategoryID(tx *gorm.DB, query *dto.SellerProductSearchQuery, sellerID uint) ([]*dto.SellerProductsCustomTable, int64, int64, error)
	SearchProducts(tx *gorm.DB, query *SearchQuery) ([]*dto.SellerProductsCustomTable, int64, int64, error)

	AddProductVariantStock(tx *gorm.DB, productVarDetID uint, amount uint) (*model.ProductVariantDetail, error)
}

type productVariantDetailRepository struct{}

func NewProductVariantDetailRepository() ProductVariantDetailRepository {
	return &productVariantDetailRepository{}
}

func (p *productVariantDetailRepository) GetProductsBySellerID(tx *gorm.DB, query *dto.SellerProductSearchQuery, sellerID uint) ([]*dto.SellerProductsCustomTable, int64, int64, error) {
	var products []*dto.SellerProductsCustomTable

	s1 := tx.Model(&model.ProductVariantDetail{})
	s1 = s1.Select("min(price), max(price), product_id")
	s1 = s1.Group("product_id")

	s2 := tx.Model(&model.Review{})
	s2 = s2.Select("count(*), AVG(rating), product_id")
	s2 = s2.Group("product_id")

	result := tx.Model(&dto.SellerProductsCustomTable{})
	result = result.Select("*")
	result = result.Joins("JOIN product_categories as c ON products.category_id = c.id")
	result = result.Joins("JOIN (?) as s1 ON products.id = s1.product_id", s1)
	result = result.Joins("LEFT JOIN (?) as s2 ON products.id = s2.product_id", s2)

	// CHANGE THIS CODE BELLOW TO CHANGE LIST OF PRODUCT BY...
	result = result.Where("seller_id = ?", sellerID)

	orderByString := query.SortBy
	if query.SortBy == "price" {
		orderByString = "min"
	} else {
		if query.SortBy == "" {
			orderByString = "sold_count"
		} else {
			orderByString = "sold_count"
			if query.SortBy == "date" {
				orderByString = "products.created_at"
			}
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
	result = result.Order(orderByString).Order("products.id")
	result = result.Where("min >= ?", query.MinAmount).Where("min <= ?", query.MaxAmount).Where("products.name ILIKE ?", "%"+query.Search+"%")
	table := tx.Table("(?) as s3", result).Count(&totalData)
	if table.Error != nil {
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

	table = table.Preload("ProductPhotos").Preload("Seller.Address")
	table = table.Unscoped().Find(&products)
	if table.Error != nil {
		return nil, 0, 0, apperror.InternalServerError("cannot fetch products")
	}

	totalPage := totalData / int64(limit)
	if totalData%int64(limit) != 0 {
		totalPage += 1
	}
	return products, totalPage, totalData, nil
}

func (p *productVariantDetailRepository) GetProductsByCategoryID(tx *gorm.DB, query *dto.SellerProductSearchQuery, categoryID uint) ([]*dto.SellerProductsCustomTable, int64, int64, error) {
	var products []*dto.SellerProductsCustomTable

	s1 := tx.Model(&model.ProductVariantDetail{})
	s1 = s1.Select("min(price), max(price), product_id")
	s1 = s1.Group("product_id")

	s2 := tx.Model(&model.Review{})
	s2 = s2.Select("count(*), AVG(rating), product_id")
	s2 = s2.Group("product_id")

	result := tx.Model(&dto.SellerProductsCustomTable{})
	result = result.Select("*")
	result = result.Joins("JOIN product_categories as c ON products.category_id = c.id")
	result = result.Joins("JOIN (?) as s1 ON products.id = s1.product_id", s1)
	result = result.Joins("LEFT JOIN (?) as s2 ON products.id = s2.product_id", s2)

	// CHANGE THIS CODE BELLOW TO CHANGE LIST OF PRODUCT BY...
	result = result.Where("(category_id = ? OR parent_id = ?)", categoryID, categoryID)

	orderByString := query.SortBy
	if query.SortBy == "price" {
		orderByString = "min"
	} else {
		if query.SortBy == "" {
			orderByString = "sold_count"
		} else {
			orderByString = "sold_count"
			if query.SortBy == "date" {
				orderByString = "products.created_at"
			}
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
	result = result.Order(orderByString).Order("products.id")
	result = result.Where("min >= ?", query.MinAmount).Where("min <= ?", query.MaxAmount).Where("products.name ILIKE ?", "%"+query.Search+"%")
	table := tx.Table("(?) as s3", result).Count(&totalData)
	if table.Error != nil {
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

	table = table.Preload("ProductPhotos").Preload("Seller.Address")
	table = table.Unscoped().Find(&products)
	if table.Error != nil {
		return nil, 0, 0, apperror.InternalServerError("cannot fetch products")
	}

	totalPage := totalData / int64(limit)
	if totalData%int64(limit) != 0 {
		totalPage += 1
	}
	return products, totalPage, totalData, nil
}

func (p *productVariantDetailRepository) SearchProducts(tx *gorm.DB, query *SearchQuery) ([]*dto.SellerProductsCustomTable, int64, int64, error) {
	var products []*dto.SellerProductsCustomTable

	s1 := tx.Model(&model.ProductVariantDetail{})
	s1 = s1.Select("min(price), max(price), product_id")
	s1 = s1.Group("product_id")

	s2 := tx.Model(&model.Review{})
	s2 = s2.Select("count(*), AVG(rating), product_id")
	s2 = s2.Group("product_id")

	seller := tx.Model(&model.Seller{})
	seller = seller.Joins("Address")
	seller = seller.Select("city, city_id, sellers.id, name")

	result := tx.Model(&dto.SellerProductsCustomTable{})
	result = result.Select("products.name, min, max, city, city_id, products.id, products.slug, products.category_id, products.favorite_count, products.seller_id, products.sold_count, avg, count, parent_id, products.created_at")
	result = result.Joins("JOIN product_categories as c ON products.category_id = c.id")
	result = result.Joins("JOIN (?) as seller ON products.seller_id = seller.id", seller)
	result = result.Joins("JOIN (?) as s1 ON products.id = s1.product_id", s1)
	result = result.Joins("LEFT JOIN (?) as s2 ON products.id = s2.product_id", s2)

	// CHANGE THIS CODE BELLOW TO CHANGE LIST OF PRODUCT BY...
	if query.CategoryID != 0 {
		result = result.Where("(category_id = ? OR parent_id = ?)", query.CategoryID, query.CategoryID)
	}
	if query.SellerID != 0 {
		result = result.Where("seller_id = ?", query.SellerID)
	}
	if query.City != "" {
		citiesSplit := strings.Split(query.City, ",")
		result = result.Where("city IN ?", citiesSplit)
	}

	orderByString := query.SortBy
	if query.SortBy == "" {
		orderByString = "sold_count"
	} else {
		orderByString = "sold_count"
		if query.SortBy == "price" {
			orderByString = "min"
		}
		if query.SortBy == "date" {
			orderByString = "products.created_at"
		}
		if query.SortBy == "favorite" {
			orderByString = "products.favorite_count"
			if query.Sort == "" {
				orderByString += " desc"
			}
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
	result = result.Order(orderByString).Order("products.id")
	result = result.Where("min >= ?", query.MinAmount).Where("min <= ?", query.MaxAmount).Where("products.name ILIKE ?", "%"+query.Search+"%")

	rating, _ := strconv.Atoi(query.Rating)
	if rating != 0 {
		result = result.Where("avg >= ? AND avg IS NOT NULL", rating)
	}

	table := tx.Table("(?) as s3", result).Count(&totalData)
	if table.Error != nil {
		return nil, 0, 0, apperror.InternalServerError("cannot fetch products count")
	}

	limit, _ := strconv.Atoi(query.Limit)
	if limit == 0 {
		limit = 20
	}
	table = table.Limit(limit)

	page, _ := strconv.Atoi(query.Page)
	if page == 0 {
		page = 1
	}
	table = table.Offset((page - 1) * limit)

	table = table.Preload("ProductPhotos").Preload("Seller.Address")
	table = table.Unscoped().Find(&products)
	if table.Error != nil {
		return nil, 0, 0, apperror.InternalServerError("cannot fetch products")
	}

	totalPage := totalData / int64(limit)
	if totalData%int64(limit) != 0 {
		totalPage += 1
	}
	return products, totalPage, totalData, nil
}

func (p *productVariantDetailRepository) AddProductVariantStock(tx *gorm.DB, productVarDetID uint, amount uint) (*model.ProductVariantDetail, error) {
	var productVarDet = &model.ProductVariantDetail{}
	productVarDet.ID = productVarDetID
	result := tx.Model(&productVarDet).Clauses(clause.Returning{}).Update("stock", gorm.Expr("stock + ?", amount))
	if result.Error != nil {
		return nil, apperror.InternalServerError("Cannot update product variant stock")
	}
	return productVarDet, nil
}
