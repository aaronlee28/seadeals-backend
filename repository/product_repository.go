package repository

import (
	"gorm.io/gorm"
	"seadeals-backend/apperror"
	"seadeals-backend/model"
	"strconv"
)

type ProductRepository interface {
	FindProductDetailByID(tx *gorm.DB, id uint) (*model.Product, error)
	FindProductBySlug(tx *gorm.DB, slug string) (*model.Product, error)
	SearchProduct(tx *gorm.DB, q *SearchQuery) (*[]model.Product, error)
}

type productRepository struct{}

func NewProductRepository() ProductRepository {
	return &productRepository{}
}

type SearchQuery struct {
	Search    string
	SortBy    string
	Sort      string
	Limit     string
	Page      string
	MinAmount string
	MaxAmount string
}

func (r *productRepository) FindProductDetailByID(tx *gorm.DB, id uint) (*model.Product, error) {
	var product *model.Product
	result := tx.Preload("ProductPhotos", "product_id = ?", id).Preload("ProductDetail", "product_id = ?", id).First(&product, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return product, nil
}

func (r *productRepository) FindProductBySlug(tx *gorm.DB, slug string) (*model.Product, error) {
	var product *model.Product
	result := tx.First(&product, "slug = ?", slug)
	if result.Error != nil {
		return nil, result.Error
	}
	return product, nil
}

func (r *productRepository) SearchProduct(tx *gorm.DB, q *SearchQuery) (*[]model.Product, error) {
	//need to join product with location and rating
	//for location, need to join product -> seller -> addresses -> sub_district-> district->cities
	//for rating, need to join product -> reviews
	//do this on a different function
	var p *[]model.Product
	search := "%" + q.Search + "%"
	limit, _ := strconv.Atoi(q.Limit)
	//page, _ := strconv.Atoi(q.Page)

	result := tx.Where("UPPER(name) like UPPER(?)", search).Limit(limit).Find(&p)
	if result.Error != nil {
		return nil, apperror.InternalServerError("cannot find product")
	}
	return p, nil
}
