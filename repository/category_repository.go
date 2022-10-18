package repository

import (
	"seadeals-backend/apperror"
	"seadeals-backend/model"
	"strconv"

	"gorm.io/gorm"
)

type ProductCategoryRepository interface {
	FindCategories(tx *gorm.DB, query *model.CategoryQuery) ([]*model.ProductCategory, int64, int64, error)
}

type productCategoryRepository struct{}

func NewProductCategoryRepository() ProductCategoryRepository {
	return &productCategoryRepository{}
}

func (r *productCategoryRepository) FindCategories(tx *gorm.DB, query *model.CategoryQuery) ([]*model.ProductCategory, int64, int64, error) {
	var categories []*model.ProductCategory

	result := tx.Model(&model.ProductCategory{})
	result = result.Distinct().Select("product_categories.id, product_categories.name, product_categories.slug, product_categories.icon_url, product_categories.parent_id")
	result = result.Joins("LEFT JOIN product_categories as c2 ON product_categories.id = c2.parent_id")
	result = result.Joins("LEFT JOIN product_categories as c3 ON c2.id = c3.parent_id")
	stringJoin := "INNER JOIN products as p ON p.category_id = product_categories.id OR p.category_id = c2.id OR p.category_id = c3.id"
	if query.FindAll {
		stringJoin = "LEFT JOIN products as p ON p.category_id = product_categories.id OR p.category_id = c2.id OR p.category_id = c3.id"
	}
	result = result.Joins(stringJoin)

	if query.ParentID == 0 {
		result = result.Where("product_categories.parent_id IS NULL")
	} else {
		result = result.Where("product_categories.parent_id = ?", query.ParentID)
	}

	if query.SellerID != 0 {
		result = result.Where("p.seller_id = ?", query.SellerID)
	}

	var totalData int64
	result = result.Order("product_categories.id")
	table := tx.Table("(?) as sub", result).Count(&totalData)
	if table.Error != nil {
		return nil, 0, 0, apperror.InternalServerError("cannot fetch categories count")
	}

	limit, _ := strconv.Atoi(query.Limit)
	totalPage := int64(1)
	if limit != 0 {
		totalPage = totalData / int64(limit)
		if totalData%int64(limit) != 0 {
			totalPage += 1
		}
		table = table.Limit(limit)
	}

	page, _ := strconv.Atoi(query.Page)
	if page != 0 {
		table = table.Offset((page - 1) * limit)
	}

	table = table.Unscoped()
	if query.FindAll {
		table = table.Preload("Children.Children.Children")
	}
	table = table.Find(&categories)
	if table.Error != nil {
		return nil, 0, 0, apperror.InternalServerError("Cannot fetch categories")
	}

	return categories, totalPage, totalData, nil
}
