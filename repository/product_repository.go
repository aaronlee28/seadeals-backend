package repository

import (
	"gorm.io/gorm"
	"seadeals-backend/apperror"
	"seadeals-backend/dto"
	"seadeals-backend/model"
	"strconv"
)

type ProductRepository interface {
	FindProductDetailByID(tx *gorm.DB, id uint) (*model.Product, error)
	FindProductBySlug(tx *gorm.DB, slug string) (*model.Product, error)

	SearchProduct(tx *gorm.DB, q *SearchQuery) (*[]model.Product, error)
	SearchProduct2(tx *gorm.DB, q *SearchQuery) ([]*dto.SearchedProductRes, error)
	SearchImageURL(tx *gorm.DB, productID uint) (string, error)
	SearchMinMaxPrice(tx *gorm.DB, productID uint) (uint, uint, error)
	SearchPromoPrice(tx *gorm.DB, productID uint) (float64, error)
	SearchRating(tx *gorm.DB, productID uint) ([]int, error)
	SearchCity(tx *gorm.DB, productID uint) (string, error)
	SearchCategory(tx *gorm.DB, productID uint) (string, error)
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

func (r *productRepository) SearchProduct2(tx *gorm.DB, q *SearchQuery) ([]*dto.SearchedProductRes, error) {
	search := "%" + q.Search + "%"
	limit, _ := strconv.Atoi(q.Limit)
	page, _ := strconv.Atoi(q.Page)
	offset := (limit * page) - limit

	var res []*dto.SearchedProductRes
	result := tx.Raw("SELECT d.product_id, slug, media_url, min_price, max as max_price, bought, updated_at FROM " +
		"(SELECT b.product_id, slug, media_url, min as min_price, bought, updated_at FROM (SELECT a.product_id as product_id, seller_id, name, slug, category_id, bought, updated_at, media_url  FROM (SELECT id as product_id, seller_id, name, slug, category_id, sold_count as bought, updated_at FROM Products WHERE UPPER(name) like UPPER('" + search + "') Limit " + strconv.Itoa(limit) + " Offset " + strconv.Itoa(offset) + ") a join (select product_id, min(photo_url) as media_url from product_photos group by product_id) as one_photo_url on a.product_id = one_photo_url.product_id) b join (select min(price), product_id from product_variant_details group by product_id) c on b.product_id = c.product_id) d " +
		"join (select max(price), product_id from product_variant_details group by product_id) e on d.product_id = e.product_id").Scan(&res)

	if result.Error != nil {
		return nil, apperror.InternalServerError("cannot find product")
	}
	return res, nil
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
		return 0, apperror.InternalServerError("cannot find promo price")
	}
	return promo, nil
}

func (r *productRepository) SearchRating(tx *gorm.DB, productID uint) ([]int, error) {
	var rating []int
	ratingQuery := tx.Select("rating").Table("reviews").Where("product_id = ?", productID).Scan(&rating)
	if ratingQuery.Error != nil {
		return nil, apperror.InternalServerError("cannot find rating")
	}

	return rating, nil
}

func (r *productRepository) SearchCity(tx *gorm.DB, productID uint) (string, error) {
	var city string
	result := tx.Raw("SELECT cities.name FROM (SELECT districts.city_id FROM (SELECT sub_districts.district_id FROM (SELECT addresses.sub_district_id FROM (SELECT products.id as product_id, sellers.address_id FROM products JOIN sellers ON products.seller_id = sellers.id WHERE products.id = ?) a JOIN addresses on a.address_id = addresses.id) b JOIN sub_districts on b.sub_district_id = sub_districts.id) c join districts on c.district_id = districts.id) d join cities on d.city_id = cities.id", productID).Scan(&city)
	if result.Error != nil {
		return "", apperror.InternalServerError("cannot find city")
	}
	return city, nil
}

func (r *productRepository) SearchCategory(tx *gorm.DB, productID uint) (string, error) {
	var category string
	categoryQuery := tx.Table("product_categories").Select("product_categories.name").Joins("join products on products.category_id = product_categories.id").Where("products.id = ?", productID).Scan(&category)
	if categoryQuery.Error != nil {
		return "", apperror.InternalServerError("cannot find category")
	}

	return category, nil
}
