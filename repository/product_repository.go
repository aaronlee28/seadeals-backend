package repository

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"seadeals-backend/apperror"
	"seadeals-backend/dto"
	"seadeals-backend/model"
	"strconv"
	"time"
)

type ProductRepository interface {
	FindProductByID(tx *gorm.DB, productID uint) (*model.Product, error)
	FindProductDetailByID(tx *gorm.DB, productID uint, userID uint) (*dto.ProductDetailRes, error)
	FindProductBySlug(tx *gorm.DB, slug string) (*model.Product, error)
	FindSimilarProduct(tx *gorm.DB, categoryID uint, query *SearchQuery) ([]*dto.SellerProductsCustomTable, int64, int64, error)

	GetProductCountBySellerID(tx *gorm.DB, sellerID uint) (int64, error)

	UpdateProductFavoriteCount(tx *gorm.DB, productID uint, isFavorite bool) (*model.Product, error)
	AddProductSoldCount(tx *gorm.DB, productID uint, total int) (*model.Product, error)

	SearchProduct(tx *gorm.DB, q *SearchQuery) (*[]model.Product, error)
	SearchRecommendProduct(tx *gorm.DB, q *SearchQuery) ([]*dto.SellerProductsCustomTable, int64, int64, error)
	SearchImageURL(tx *gorm.DB, productID uint) (string, error)
	SearchMinMaxPrice(tx *gorm.DB, productID uint) (float64, float64, error)
	SearchPromoPrice(tx *gorm.DB, productID uint) (float64, error)
	SearchRating(tx *gorm.DB, productID uint) ([]int, error)
	SearchCity(tx *gorm.DB, productID uint) (string, error)
	SearchCategory(tx *gorm.DB, productID uint) (string, error)
	GetProductDetail(tx *gorm.DB, id uint) (*model.Product, error)
	GetProductPhotoURL(tx *gorm.DB, productID uint) (string, error)

	CreateProduct(tx *gorm.DB, name string, categoryID uint, sellerID uint, bulk bool, minQuantity uint, maxQuantity uint) (*model.Product, error)
	CreateProductDetail(tx *gorm.DB, productID uint, req *dto.ProductDetailsReq) (*model.ProductDetail, error)
	CreateProductPhoto(tx *gorm.DB, productID uint, req *dto.ProductPhoto) (*model.ProductPhoto, error)
	CreateProductVariant(tx *gorm.DB, name string) (*model.ProductVariant, error)
	CreateProductVariantDetail(tx *gorm.DB, productID uint, variant1ID *uint, variant2 *model.ProductVariant, req *dto.ProductVariantDetail) (*model.ProductVariantDetail, error)
	UpdateProduct(tx *gorm.DB, productID uint, p *model.Product) (*model.Product, error)
	UpdateProductDetail(tx *gorm.DB, productID uint, pd *model.ProductDetail) (*model.ProductDetail, error)
	FindProductVariantDetailsByID(tx *gorm.DB, id uint) (*model.ProductVariantDetail, error)
	DeleteProductVariantDetailsByID(tx *gorm.DB, id uint) error
	UpdateProductVariantDetailByID(tx *gorm.DB, id uint, pvd *model.ProductVariantDetail) (*model.ProductVariantDetail, error)
	UpdateProductVariantByID(tx *gorm.DB, id uint, pvd *model.ProductVariant) (*model.ProductVariant, error)
	FindProductVariantDetailsByProductID(tx *gorm.DB, ProductID uint) ([]*model.ProductVariantDetail, error)
	GetVariantByName(tx *gorm.DB, name string) (*model.ProductVariant, error)
	CreateVariantWithName(tx *gorm.DB, name string) (*model.ProductVariant, error)
	CreateProductVariantDetailWithModel(tx *gorm.DB, pvd *model.ProductVariantDetail) (*model.ProductVariantDetail, error)
	DeleteNullProductVariantDetailsByID(tx *gorm.DB, ProductID uint) error
	CreateProductPhotos(tx *gorm.DB, productID uint, req *dto.ProductPhotoReq) ([]*model.ProductPhoto, error)
	DeleteProductPhotos(tx *gorm.DB, req *dto.DeleteProductPhoto) ([]*model.ProductPhoto, error)
	DeleteProduct(tx *gorm.DB, productID uint) (*model.Product, error)
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
	ExcludedID uint
}

func (r *productRepository) FindProductByID(tx *gorm.DB, productID uint) (*model.Product, error) {
	var product *model.Product
	result := tx.First(&product, productID)
	if result.Error == gorm.ErrRecordNotFound {
		return nil, apperror.NotFoundError(new(apperror.ProductNotFoundError).Error())
	}
	return product, result.Error
}

func (r *productRepository) FindProductDetailByID(tx *gorm.DB, productID uint, userID uint) (*dto.ProductDetailRes, error) {
	var productVariantDetail *model.ProductVariantDetail
	variant := tx.Model(&productVariantDetail)
	variant = variant.Select("SUM(product_variant_details.stock) as total_stock, product_variant_details.product_id")
	variant = variant.Group("product_variant_details.product_id")

	var productReview *model.Review
	review := tx.Model(&productReview)
	review = review.Select("AVG(rating) as average_rating, COUNT(*) as total_review, reviews.product_id")
	review = review.Group("reviews.product_id")

	var product *dto.ProductDetailRes
	result := tx.Model(&product)
	result = result.Select("*")
	result = result.Preload("ProductPhotos", "product_id = ?", productID)
	result = result.Preload("ProductDetail", "product_id = ?", productID)
	result = result.Preload("ProductVariantDetail", "product_id = ?", productID, func(db *gorm.DB) *gorm.DB {
		return db.Order("product_variant_details.price")
	})
	result = result.Preload("Category")
	result = result.Preload("Promotion")
	result = result.Preload("ProductVariantDetail.ProductVariant1")
	result = result.Preload("ProductVariantDetail.ProductVariant2")
	result = result.Preload("Favorite", "product_id = ? AND user_id = ? AND is_favorite IS TRUE", productID, userID)
	result = result.Joins("JOIN (?) AS s1 ON s1.product_id = products.id", variant)
	result = result.Joins("LEFT JOIN (?) AS s2 ON s2.product_id = products.id", review)
	result = result.First(&product, productID)
	if result.Error != nil {
		return nil, result.Error
	}
	return product, nil
}

func (r *productRepository) FindSimilarProduct(tx *gorm.DB, categoryID uint, query *SearchQuery) ([]*dto.SellerProductsCustomTable, int64, int64, error) {
	var products []*dto.SellerProductsCustomTable

	promotions := tx.Model(&model.Promotion{})
	promotions = promotions.Where("start_date <= ? AND end_date >= ?", time.Now(), time.Now())
	promotions = promotions.Order("product_id, id")
	promotions = promotions.Select("DISTINCT ON (product_id) *")

	s1 := tx.Model(&model.ProductVariantDetail{})
	s1 = s1.Select("min(price - COALESCE(promotions.amount, 0)) as min, max(price - COALESCE(promotions.amount, 0)) as max, min(price) as min_before_disc, max(price) as max_before_disc, product_variant_details.product_id")
	s1 = s1.Joins("LEFT JOIN (?) as promotions ON promotions.product_id = product_variant_details.product_id", promotions)
	s1 = s1.Group("product_variant_details.product_id")

	s2 := tx.Model(&model.Review{})
	s2 = s2.Select("count(*), AVG(rating), product_id")
	s2 = s2.Group("product_id")

	seller := tx.Model(&model.Seller{})
	seller = seller.Joins("Address")
	seller = seller.Select("city, city_id, sellers.id, name")

	result := tx.Model(&dto.SellerProductsCustomTable{})
	result = result.Select("products.name, min, max, min_before_disc, max_before_disc, city, city_id, products.id, products.slug, p.amount as promotion_amount, p.id as promotion_id, products.category_id, products.favorite_count, products.seller_id, products.sold_count, avg, count, c.parent_id, c2.parent_id, products.created_at")
	result = result.Where("products.id != ?", query.ExcludedID)
	result = result.Joins("JOIN product_categories as c ON products.category_id = c.id")
	result = result.Joins("LEFT JOIN product_categories as c2 ON c.parent_id = c2.id")
	result = result.Joins("LEFT JOIN (?) as p ON p.product_id = products.id", promotions)
	result = result.Joins("JOIN (?) as seller ON products.seller_id = seller.id", seller)
	result = result.Joins("JOIN (?) as s1 ON products.id = s1.product_id", s1)
	result = result.Joins("LEFT JOIN (?) as s2 ON products.id = s2.product_id", s2)

	// CHANGE THIS CODE BELLOW TO CHANGE LIST OF PRODUCT BY...
	result = result.Where("(category_id = ? OR c.parent_id = ? OR c2.parent_id = ?)", categoryID, categoryID, categoryID)

	var totalData int64

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

func (r *productRepository) GetProductCountBySellerID(tx *gorm.DB, sellerID uint) (int64, error) {
	var totalProduct int64
	result := tx.Model(&model.Product{}).Where("seller_id = ?", sellerID).Where("is_archived IS FALSE").Count(&totalProduct)
	if result.Error != nil {
		return 0, apperror.InternalServerError("Cannot count total product")
	}
	return totalProduct, nil
}

func (r *productRepository) UpdateProductFavoriteCount(tx *gorm.DB, productID uint, isFavorite bool) (*model.Product, error) {
	var product = &model.Product{}
	product.ID = productID
	var result *gorm.DB
	if isFavorite {
		result = tx.Model(&product).Clauses(clause.Returning{}).Update("favorite_count", gorm.Expr("favorite_count + 1"))
	} else {
		result = tx.Model(&product).Clauses(clause.Returning{}).Update("favorite_count", gorm.Expr("favorite_count - 1"))
	}
	if result.Error != nil {
		return nil, apperror.InternalServerError("Cannot update product")
	}
	return product, nil
}

func (r *productRepository) AddProductSoldCount(tx *gorm.DB, productID uint, total int) (*model.Product, error) {
	var product = &model.Product{ID: productID}
	result := tx.Model(&product).Update("sold_count", gorm.Expr("sold_count + ?", total))
	if result.Error != nil {
		return nil, apperror.InternalServerError("Tidak bisa menambah sold_count product")
	}
	return product, nil
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

	promotions := tx.Model(&model.Promotion{})
	promotions = promotions.Where("start_date <= ? AND end_date >= ?", time.Now(), time.Now())
	promotions = promotions.Order("product_id, id")
	promotions = promotions.Select("DISTINCT ON (product_id) *")

	s1 := tx.Model(&model.ProductVariantDetail{})
	s1 = s1.Select("min(price - COALESCE(promotions.amount, 0)) as min, max(price - COALESCE(promotions.amount, 0)) as max, min(price) as min_before_disc, max(price) as max_before_disc, product_variant_details.product_id")
	s1 = s1.Joins("LEFT JOIN (?) as promotions ON promotions.product_id = product_variant_details.product_id", promotions)
	s1 = s1.Group("product_variant_details.product_id")

	s2 := tx.Model(&model.Review{})
	s2 = s2.Select("count(*), AVG(rating), product_id")
	s2 = s2.Group("product_id")

	seller := tx.Model(&model.Seller{})
	seller = seller.Joins("Address")
	seller = seller.Select("city, city_id, sellers.id, name")

	result := tx.Model(&dto.SellerProductsCustomTable{})
	result = result.Select("products.name, min, max, min_before_disc, max_before_disc, city, city_id, products.id, products.slug, p.amount as promotion_amount, p.id as promotion_id, products.category_id, products.favorite_count, products.seller_id, products.sold_count, avg, count, parent_id, products.created_at")
	result = result.Joins("JOIN product_categories as c ON products.category_id = c.id")
	result = result.Joins("LEFT JOIN (?) as p ON p.product_id = products.id", promotions)
	result = result.Joins("JOIN (?) as seller ON products.seller_id = seller.id", seller)
	result = result.Joins("JOIN (?) as s1 ON products.id = s1.product_id", s1)
	result = result.Joins("LEFT JOIN (?) as s2 ON products.id = s2.product_id", s2)
	orderByString := "favorite_count desc"

	if query.ExcludedID != 0 {
		result = result.Where("products.id != ?", query.ExcludedID)
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

func (r *productRepository) SearchMinMaxPrice(tx *gorm.DB, productID uint) (float64, float64, error) {
	var min, max float64

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
func (r *productRepository) CreateProductDetail(tx *gorm.DB, productID uint, req *dto.ProductDetailsReq) (*model.ProductDetail, error) {

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
func (r *productRepository) UpdateProduct(tx *gorm.DB, productID uint, p *model.Product) (*model.Product, error) {
	var updatedProduct *model.Product
	result := tx.First(&updatedProduct, productID).Updates(&p)
	return updatedProduct, result.Error
}

func (r *productRepository) UpdateProductDetail(tx *gorm.DB, productID uint, pd *model.ProductDetail) (*model.ProductDetail, error) {

	var updatedProductDetail *model.ProductDetail
	result := tx.First(&updatedProductDetail, "product_id = ?", productID).Updates(&pd)
	return updatedProductDetail, result.Error
}

func (r *productRepository) FindProductVariantDetailsByID(tx *gorm.DB, id uint) (*model.ProductVariantDetail, error) {
	var productVariantDetails *model.ProductVariantDetail
	result := tx.First(&productVariantDetails, id)
	if result.Error == gorm.ErrRecordNotFound {
		return nil, apperror.NotFoundError(new(apperror.ProductNotFoundError).Error())
	}
	return productVariantDetails, result.Error
}
func (r *productRepository) FindProductVariantDetailsByProductID(tx *gorm.DB, ProductID uint) ([]*model.ProductVariantDetail, error) {
	var productVariantDetails []*model.ProductVariantDetail
	result := tx.Find(&productVariantDetails, "product_id = ?", ProductID)
	if result.Error == gorm.ErrRecordNotFound {
		return nil, apperror.NotFoundError(new(apperror.ProductNotFoundError).Error())
	}
	return productVariantDetails, result.Error
}
func (r *productRepository) DeleteProductVariantDetailsByID(tx *gorm.DB, id uint) error {
	var deletedVariantDetail *model.ProductVariantDetail
	result := tx.Delete(&deletedVariantDetail, id)
	return result.Error
}
func (r *productRepository) DeleteNullProductVariantDetailsByID(tx *gorm.DB, ProductID uint) error {
	var productVariantDetails *model.ProductVariantDetail
	result := tx.First(&productVariantDetails, "product_id = ?", ProductID).Where("variant1_id = null")
	if result.Error == gorm.ErrRecordNotFound {
		return apperror.NotFoundError(new(apperror.ProductNotFoundError).Error())
	}
	result2 := tx.Delete(&productVariantDetails)

	return result2.Error
}
func (r *productRepository) UpdateProductVariantDetailByID(tx *gorm.DB, id uint, pvd *model.ProductVariantDetail) (*model.ProductVariantDetail, error) {

	var updatedProductVariantDetail *model.ProductVariantDetail
	result := tx.First(&updatedProductVariantDetail, id).Updates(&pvd)
	return updatedProductVariantDetail, result.Error
}

func (r *productRepository) UpdateProductVariantByID(tx *gorm.DB, id uint, pvd *model.ProductVariant) (*model.ProductVariant, error) {
	var updatedProductVariant *model.ProductVariant
	result := tx.First(&updatedProductVariant, id).Updates(&pvd)
	return updatedProductVariant, result.Error
}
func (r *productRepository) GetVariantByName(tx *gorm.DB, name string) (*model.ProductVariant, error) {
	var updatedProductVariant *model.ProductVariant
	result := tx.First(&updatedProductVariant, "name = ?", name)
	return updatedProductVariant, result.Error
}

func (r *productRepository) CreateVariantWithName(tx *gorm.DB, name string) (*model.ProductVariant, error) {
	createPV := model.ProductVariant{
		Name: name,
	}
	result := tx.Create(&createPV)

	if result.Error != nil {
		return nil, apperror.InternalServerError("Cannot create variant")
	}
	return &createPV, result.Error
}

func (r *productRepository) CreateProductVariantDetailWithModel(tx *gorm.DB, pvd *model.ProductVariantDetail) (*model.ProductVariantDetail, error) {
	result := tx.Create(&pvd)
	if result.Error != nil {
		return nil, apperror.InternalServerError("Cannot create variant detail")
	}
	return pvd, result.Error
}
func (r *productRepository) CreateProductPhotos(tx *gorm.DB, productID uint, req *dto.ProductPhotoReq) ([]*model.ProductPhoto, error) {
	var ret []*model.ProductPhoto
	for _, ph := range req.ProductPhoto {

		productPhoto := model.ProductPhoto{
			ProductID: productID,
			PhotoURL:  ph.PhotoURL,
			Name:      ph.Name,
		}
		result := tx.Clauses(clause.Returning{}).Create(&productPhoto)
		if result.Error != nil {
			return nil, apperror.InternalServerError("Cannot create product photo")
		}

		ret = append(ret, &productPhoto)
	}

	return ret, nil
}
func (r *productRepository) DeleteProductPhotos(tx *gorm.DB, req *dto.DeleteProductPhoto) ([]*model.ProductPhoto, error) {
	var ret []*model.ProductPhoto
	for _, ph := range req.PhotoId {
		var p *model.ProductPhoto
		find := tx.Where("id = ?", ph).First(&p)
		if find.Error != nil {
			return nil, apperror.InternalServerError("Cannot find id")
		}
		ret = append(ret, p)
	}
	result := tx.Clauses(clause.Returning{}).Delete(&ret)
	if result.Error != nil {
		return nil, apperror.InternalServerError("Cannot delete product photo")
	}
	return ret, nil
}

func (r *productRepository) DeleteProduct(tx *gorm.DB, productID uint) (*model.Product, error) {
	var ret *model.Product

	result := tx.Clauses(clause.Returning{}).Where("id = ?", productID).First(&ret).Delete(&ret)
	if result.Error != nil {
		return nil, apperror.InternalServerError("Cannot delete product")
	}
	return ret, nil
}
