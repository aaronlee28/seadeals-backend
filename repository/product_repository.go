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
	SearchImageURL(tx *gorm.DB, productID uint) (string, error)
	SearchMinMaxPrice(tx *gorm.DB, productID uint) (uint, uint, error)
	SearchPromoPrice(tx *gorm.DB, productID uint) (float64, error)
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
	page, _ := strconv.Atoi(q.Page)
	offset := (limit * page) - limit

	result := tx.Where("UPPER(name) like UPPER(?)", search).Limit(limit).Offset(offset).Find(&p)
	if result.Error != nil {
		return nil, apperror.InternalServerError("cannot find product")
	}
	return p, nil
}
func (r *productRepository) SearchImageURL(tx *gorm.DB, productID uint) (string, error) {
	var url string
	result := tx.Raw("SELECT photo_url FROM (select product_id, min(id) as First from product_photos group by product_id) foo join product_photos p on foo.product_id = p.product_id and foo.First = p.id where p.product_id=?", productID).Scan(&url)
	if result.Error != nil {
		return "", apperror.InternalServerError("cannot find image")
	}
	return url, nil
}

func (r *productRepository) SearchMinMaxPrice(tx *gorm.DB, productID uint) (uint, uint, error) {
	var min, max uint

	minQuery := tx.Select("price").Table("product_variant_details").Where("product_id = ?", productID).Order("price asc").Limit(1).Scan(&min)

	if minQuery.Error != nil {
		return 0, 0, apperror.InternalServerError("cannot find price")
	}

	maxQuery := tx.Select("price").Table("product_variant_details").Where("product_id = ?", productID).Order("price desc").Limit(1).Scan(&max)

	if maxQuery.Error != nil {
		return 0, 0, apperror.InternalServerError("cannot find price")
	}
	return min, max, nil
}

func (r *productRepository) SearchPromoPrice(tx *gorm.DB, productID uint) (float64, error) {
	var promo float64

	promoQuery := tx.Select("amount").Table("promotions").Where("product_id = ?", productID).Order("amount asc").Limit(1).Scan(&promo)
	if promoQuery.Error != nil {
		return 0, apperror.InternalServerError("cannot find price")
	}
	return promo, nil
}
