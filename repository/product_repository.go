package repository

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"seadeals-backend/apperror"
	"seadeals-backend/dto"
	"seadeals-backend/model"
	"strconv"
)

type ProductRepository interface {
	FindProductDetailByID(tx *gorm.DB, productID uint, userID uint) (*model.Product, error)
	FindProductBySlug(tx *gorm.DB, slug string) (*model.Product, error)
	FindSimilarProduct(tx *gorm.DB, productID uint) ([]*model.Product, error)

	SearchProduct(tx *gorm.DB, q *SearchQuery) (*[]model.Product, error)
	SearchRecommendProduct(tx *gorm.DB, q *SearchQuery) ([]*dto.SearchedProductRes, error)
	SearchImageURL(tx *gorm.DB, productID uint) (string, error)
	SearchMinMaxPrice(tx *gorm.DB, productID uint) (uint, uint, error)
	SearchPromoPrice(tx *gorm.DB, productID uint) (float64, error)
	SearchRating(tx *gorm.DB, productID uint) ([]int, error)
	SearchCity(tx *gorm.DB, productID uint) (string, error)
	SearchCategory(tx *gorm.DB, productID uint) (string, error)
	GetProductDetail(tx *gorm.DB, id uint) (*model.Product, error)
	GetProductPhotoURL(tx *gorm.DB, productID uint) (string, error)
}

type productRepository struct{}

func NewProductRepository() ProductRepository {
	return &productRepository{}
}

type SearchQuery struct {
	Search     string
	SortBy     string
	Sort       string
	Limit      string
	Page       string
	MinAmount  float64
	MaxAmount  float64
	City       string
	Rating     string
	Category   string
	SellerID   uint
	CategoryID uint
}

func (r *productRepository) FindProductDetailByID(tx *gorm.DB, productID uint, userID uint) (*model.Product, error) {
	var product *model.Product
	result := tx.Preload("ProductPhotos", "product_id = ?", productID)
	result = result.Preload("ProductDetail", "product_id = ?", productID)
	result = result.Preload("ProductVariantDetail", "product_id = ?", productID)
	result = result.Preload("Favorite", "product_id = ? AND user_id = ? AND is_favorite IS TRUE", productID, userID)
	result = result.First(&product, productID)
	if result.Error != nil {
		return nil, result.Error
	}
	return product, nil
}

func (r *productRepository) FindSimilarProduct(tx *gorm.DB, categoryID uint) ([]*model.Product, error) {
	var products []*model.Product
	result := tx.Limit(24).Where("category_id = ?", categoryID).Preload("ProductVariantDetail").Preload("ProductPhotos").Find(&products)
	if result.Error != nil {
		return nil, result.Error
	}
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, &apperror.ProductNotFoundError{}
	}
	return products, nil
}

func (r *productRepository) GetProductDetail(tx *gorm.DB, id uint) (*model.Product, error) {
	var product *model.Product
	result := tx.Preload("ProductVariantDetail", "product_id = ?", id).Preload("Promotion", "product_id = ?", id).Where("id = ?", id).First(&product, id)
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

func (r *productRepository) SearchRecommendProduct(tx *gorm.DB, q *SearchQuery) ([]*dto.SearchedProductRes, error) {
	search := "%" + q.Search + "%"
	city := "%" + q.City + "%"
	category := "%" + q.Category + "%"
	limit, _ := strconv.Atoi(q.Limit)
	page, _ := strconv.Atoi(q.Page)
	offset := (limit * page) - limit

	var res []*dto.SearchedProductRes
	result := tx.Raw("SELECT product_id as id, product_name as name, slug, media_url, min_price as price, min_price, max_price, total_sold, views_count as views, promo_price, rating, count as total_reviewer, city, category, updated_at FROM " +
		"(SELECT j.product_id as product_id, product_name, slug, media_url, min_price, max_price, total_sold, promo_price, rating, count, name as city, category_id, views_count, updated_at FROM " +
		"(SELECT h.product_id, product_name, slug, media_url, min_price, max_price, total_sold, promo_price, rating, count, updated_at FROM" +
		"(SELECT f.product_id, product_name, slug, media_url, min_price, max_price, total_sold, updated_at, min as promo_price FROM " +
		"(SELECT d.product_id, product_name, slug, media_url, min_price, max as max_price, total_sold, updated_at FROM " +
		"(SELECT b.product_id, name as product_name, slug, media_url, min as min_price, total_sold, updated_at FROM (SELECT a.product_id as product_id, seller_id, name, slug, category_id, views_count, total_sold, updated_at, media_url  FROM (SELECT id as product_id, seller_id, name, slug, category_id, views_count, sold_count as total_sold, updated_at FROM Products WHERE UPPER(name) like UPPER('" + search + "') Limit " + strconv.Itoa(limit) + " Offset " + strconv.Itoa(offset) + ") a left join (select ab.product_id, ab.photo_url as media_url FROM (SELECT product_id, min(id) AS First FROM product_photos GROUP BY product_id) foo join product_photos ab on foo.product_id = ab.product_id and foo.First = ab.id) as one_photo_url on a.product_id = one_photo_url.product_id) b left join (select min(price), product_id from product_variant_details group by product_id) c on b.product_id = c.product_id) d " +
		"left join (select max(price), product_id from product_variant_details group by product_id) e on d.product_id = e.product_id) f " +
		"left join (select product_id, min(amount) from promotions group by product_id) g on f.product_id = g.product_id) h " +
		"left join (select avg(rating) as rating, count(rating) as count, product_id from reviews group by product_id) i on h.product_id = i.product_id) j " +
		"left join (SELECT product_id, cities.name, category_id, views_count FROM (SELECT product_id, districts.city_id, category_id, views_count FROM (SELECT product_id, sub_districts.district_id, category_id, views_count FROM (SELECT product_id, addresses.sub_district_id, category_id, views_count FROM (SELECT products.id as product_id, sellers.address_id, products.category_id as category_id, products.views_count FROM products JOIN sellers ON products.seller_id = sellers.id) aa JOIN addresses on aa.address_id = addresses.id) bb JOIN sub_districts on bb.sub_district_id = sub_districts.id) cc join districts on cc.district_id = districts.id) dd join cities on dd.city_id = cities.id) k " +
		"on j.product_id = k.product_id) l " +
		"left join (SELECT id as category_id, name as category from product_categories) m " +
		"on l.category_id = m.category_id" +
		" where min_price >= " +
		fmt.Sprintf("%f", q.MinAmount) +
		" and " +
		"max_price <= " +
		fmt.Sprintf("%f", q.MaxAmount) +
		" and " +
		"UPPER(city) like UPPER('" +
		city +
		"')" +
		" and " +
		"rating >= " +
		q.Rating +
		" or rating is null " +
		"and " +
		"UPPER(category) like UPPER('" +
		category +
		"')" +
		" order by " +
		q.SortBy +
		" " +
		q.Sort).Scan(&res)

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
	result := tx.Raw("SELECT cities.name FROM (SELECT districts.city_id FROM (SELECT sub_districts.district_id FROM (SELECT addresses.sub_district_id FROM (SELECT products.id as product_id, sellers.address_id FROM products JOIN sellers ON products.seller_id = sellers.id WHERE products.id = ?) aa JOIN addresses on aa.address_id = addresses.id) bb JOIN sub_districts on bb.sub_district_id = sub_districts.id) cc join districts on cc.district_id = districts.id) dd join cities on dd.city_id = cities.id", productID).Scan(&city)
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

func (r *productRepository) GetProductPhotoURL(tx *gorm.DB, productID uint) (string, error) {
	var photoURL string
	photoQuery := tx.Table("product_photos").Select("photo_url").Where("product_id = ?", productID).Limit(1).Find(&photoURL)
	if photoQuery.Error != nil {
		return "", apperror.InternalServerError("cannot find photo")
	}

	return photoURL, nil
}
