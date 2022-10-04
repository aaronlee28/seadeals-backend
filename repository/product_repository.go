package repository

import (
	"gorm.io/gorm"
	"seadeals-backend/apperror"
	"seadeals-backend/dto"
	"seadeals-backend/model"
	"strconv"
)

type ProductRepository interface {
	FindProductByID(tx *gorm.DB, productID uint) (*model.Product, error)
	FindProductDetailByID(tx *gorm.DB, productID uint, userID uint) (*model.Product, error)
	FindProductBySlug(tx *gorm.DB, slug string) (*model.Product, error)
	FindSimilarProduct(tx *gorm.DB, categoryID uint) ([]*model.Product, error)

	SearchProduct(tx *gorm.DB, q *SearchQuery) (*[]model.Product, error)
	SearchRecommendProduct(tx *gorm.DB, q *SearchQuery) ([]*dto.SellerProductsCustomTable, int64, int64, error)
	SearchImageURL(tx *gorm.DB, productID uint) (string, error)
	SearchMinMaxPrice(tx *gorm.DB, productID uint) (uint, uint, error)
	SearchPromoPrice(tx *gorm.DB, productID uint) (float64, error)
	SearchRating(tx *gorm.DB, productID uint) ([]int, error)
	SearchCity(tx *gorm.DB, productID uint) (string, error)
	SearchCategory(tx *gorm.DB, productID uint) (string, error)
	GetProductDetail(tx *gorm.DB, id uint) (*model.Product, error)
	GetProductPhotoURL(tx *gorm.DB, productID uint) (string, error)
	CreateProduct(tx *gorm.DB, name string, categoryID uint, sellerID uint, bulk bool, minQuantity uint, maxQuantity uint) (*model.Product, error)
	CreateProductDetail(tx *gorm.DB, productID uint, req *dto.ProductDetailReq) (*model.ProductDetail, error)
	CreateProductPhoto(tx *gorm.DB, productID uint, req *dto.ProductPhoto) (*model.ProductPhoto, error)
	CreateProductVariant(tx *gorm.DB, name string) (*model.ProductVariant, error)
	CreateProductVariantDetail(tx *gorm.DB, productID uint, variant1ID *uint, variant2 *model.ProductVariant, req *dto.ProductVariantDetail) (*model.ProductVariantDetail, error)
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

func (r *productRepository) FindProductByID(tx *gorm.DB, productID uint) (*model.Product, error) {
	var product *model.Product
	result := tx.First(&product, productID)
	if result.Error == gorm.ErrRecordNotFound {
		return nil, apperror.NotFoundError(new(apperror.ProductNotFoundError).Error())
	}
	return product, result.Error
}

func (r *productRepository) FindProductDetailByID(tx *gorm.DB, productID uint, userID uint) (*model.Product, error) {
	var product *model.Product
	result := tx.Preload("ProductPhotos", "product_id = ?", productID)
	result = result.Preload("ProductDetail", "product_id = ?", productID)
	result = result.Preload("ProductVariantDetail", "product_id = ?", productID, func(db *gorm.DB) *gorm.DB {
		return db.Order("product_variant_details.price")
	})
	result = result.Preload("ProductVariantDetail.ProductVariant1")
	result = result.Preload("ProductVariantDetail.ProductVariant2")
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
	return products, result.Error
}

func (r *productRepository) GetProductDetail(tx *gorm.DB, id uint) (*model.Product, error) {
	var product *model.Product
	result := tx.Preload("ProductVariantDetail", "product_id = ?", id).Preload("Promotion", "product_id = ?", id).First(&product, id)
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

func (r *productRepository) SearchRecommendProduct(tx *gorm.DB, query *SearchQuery) ([]*dto.SellerProductsCustomTable, int64, int64, error) {
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
	orderByString := "favorite_count desc"

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
		limit = 18
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

func (r *productRepository) CreateProduct(tx *gorm.DB, name string, categoryID uint, sellerID uint, bulk bool, minQuantity uint, maxQuantity uint) (*model.Product, error) {
	product := &model.Product{
		Name:          name,
		CategoryID:    categoryID,
		SellerID:      sellerID,
		IsBulkEnabled: bulk,
		MinQuantity:   minQuantity,
		MaxQuantity:   maxQuantity,
	}
	result := tx.Create(&product)
	if result.Error != nil {
		return nil, apperror.InternalServerError("Cannot create product")
	}
	return product, nil
}
func (r *productRepository) CreateProductDetail(tx *gorm.DB, productID uint, req *dto.ProductDetailReq) (*model.ProductDetail, error) {

	productDetail := &model.ProductDetail{
		ProductID:       productID,
		Description:     req.Description,
		VideoURL:        req.VideoURL,
		IsHazardous:     *req.IsHazardous,
		ConditionStatus: req.ConditionStatus,
		Length:          req.Length,
		Width:           req.Width,
		Height:          req.Height,
		Weight:          req.Weight,
	}
	result := tx.Create(&productDetail)
	if result.Error != nil {
		return nil, apperror.InternalServerError("Cannot create product details")
	}
	return productDetail, nil
}

func (r *productRepository) CreateProductPhoto(tx *gorm.DB, productID uint, req *dto.ProductPhoto) (*model.ProductPhoto, error) {

	productPhoto := &model.ProductPhoto{
		ProductID: productID,
		PhotoURL:  req.PhotoURL,
		Name:      req.Name,
	}
	result := tx.Create(&productPhoto)
	if result.Error != nil {
		return nil, apperror.InternalServerError("Cannot create product photo")
	}
	return productPhoto, nil
}

func (r *productRepository) CreateProductVariant(tx *gorm.DB, name string) (*model.ProductVariant, error) {

	productVariant := &model.ProductVariant{
		Name: name,
	}
	result := tx.Create(&productVariant)
	if result.Error != nil {
		return nil, apperror.InternalServerError("Cannot create product variant")
	}
	return productVariant, nil
}

func (r *productRepository) CreateProductVariantDetail(tx *gorm.DB, productID uint, variant1ID *uint, variant2 *model.ProductVariant, req *dto.ProductVariantDetail) (*model.ProductVariantDetail, error) {

	productVariantDetail := &model.ProductVariantDetail{
		ProductID:     productID,
		Price:         req.Price,
		Variant1Value: req.Variant1Value,
		Variant1ID:    variant1ID,
		VariantCode:   req.VariantCode,
		PictureURL:    req.PictureURL,
		Stock:         req.Stock,
	}
	if req.Variant2Value != nil {
		productVariantDetail.Variant2Value = req.Variant2Value
	}
	if variant2 != nil {
		productVariantDetail.Variant2ID = variant2.ID
	}
	result := tx.Create(&productVariantDetail)
	if result.Error != nil {
		return nil, apperror.InternalServerError("Cannot create product variant detail")
	}
	return productVariantDetail, nil
}
