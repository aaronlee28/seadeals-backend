package service

import (
	"gorm.io/gorm"
	"seadeals-backend/dto"
	"seadeals-backend/helper"
	"seadeals-backend/model"
	"seadeals-backend/repository"
)

type ProductService interface {
	FindProductDetailByID(productID uint, userID uint) (*model.Product, error)
	FindSimilarProducts(productID uint) ([]*dto.ProductRes, error)
	SearchRecommendProduct(q *repository.SearchQuery) ([]*dto.ProductRes, int64, int64, error)
	GetProductsBySellerID(query *dto.SellerProductSearchQuery, sellerID uint) ([]*dto.ProductRes, int64, int64, error)
	GetProductsByCategoryID(query *dto.SellerProductSearchQuery, categoryID uint) ([]*dto.ProductRes, int64, int64, error)
	GetProducts(q *repository.SearchQuery) ([]*dto.ProductRes, int64, int64, error)
	CreateSellerProduct(userID uint, req *dto.PostCreateProductReq) (*dto.PostCreateProductRes, error)
}

type productService struct {
	db                *gorm.DB
	productRepo       repository.ProductRepository
	reviewRepo        repository.ReviewRepository
	productVarDetRepo repository.ProductVariantDetailRepository
	sellerRepo        repository.SellerRepository
}

type ProductConfig struct {
	DB                *gorm.DB
	ProductRepo       repository.ProductRepository
	ReviewRepo        repository.ReviewRepository
	ProductVarDetRepo repository.ProductVariantDetailRepository
	SellerRepo        repository.SellerRepository
}

func NewProductService(config *ProductConfig) ProductService {
	return &productService{
		db:                config.DB,
		productRepo:       config.ProductRepo,
		reviewRepo:        config.ReviewRepo,
		productVarDetRepo: config.ProductVarDetRepo,
		sellerRepo:        config.SellerRepo,
	}
}

func (p *productService) FindProductDetailByID(productID uint, userID uint) (*model.Product, error) {
	tx := p.db.Begin()
	var err error
	defer helper.CommitOrRollback(tx, &err)

	product, err := p.productRepo.FindProductDetailByID(tx, productID, userID)
	if err != nil {
		return nil, err
	}

	return product, nil
}

func (p *productService) GetProductsBySellerID(query *dto.SellerProductSearchQuery, sellerID uint) ([]*dto.ProductRes, int64, int64, error) {
	tx := p.db.Begin()
	var err error
	defer helper.CommitOrRollback(tx, &err)

	variantDetails, totalPage, totalData, err := p.productVarDetRepo.GetProductsBySellerID(tx, query, sellerID)
	if err != nil {
		return nil, 0, 0, err
	}

	var productsRes = make([]*dto.ProductRes, 0)
	for _, variantDetail := range variantDetails {
		var photoURL string
		if len(variantDetail.Product.ProductPhotos) > 0 {
			photoURL = variantDetail.Product.ProductPhotos[0].PhotoURL
		}

		dtoProduct := &dto.ProductRes{
			MinPrice: variantDetail.Min,
			MaxPrice: variantDetail.Max,
			Product: &dto.GetProductRes{
				ID:            variantDetail.ProductID,
				Price:         variantDetail.Min,
				Name:          variantDetail.Product.Name,
				Slug:          variantDetail.Product.Slug,
				MediaURL:      photoURL,
				City:          variantDetail.Product.Seller.Address.City,
				Rating:        variantDetail.Avg,
				TotalReviewer: variantDetail.Count,
				TotalSold:     uint(variantDetail.Product.SoldCount),
			},
		}
		productsRes = append(productsRes, dtoProduct)
	}

	return productsRes, totalPage, totalData, nil
}

func (p *productService) GetProductsByCategoryID(query *dto.SellerProductSearchQuery, categoryID uint) ([]*dto.ProductRes, int64, int64, error) {
	tx := p.db.Begin()
	var err error
	defer helper.CommitOrRollback(tx, &err)

	variantDetails, totalPage, totalData, err := p.productVarDetRepo.GetProductsByCategoryID(tx, query, categoryID)
	if err != nil {
		return nil, 0, 0, err
	}

	var productsRes = make([]*dto.ProductRes, 0)
	for _, variantDetail := range variantDetails {
		var photoURL string
		if len(variantDetail.Product.ProductPhotos) > 0 {
			photoURL = variantDetail.Product.ProductPhotos[0].PhotoURL
		}

		dtoProduct := &dto.ProductRes{
			MinPrice: variantDetail.Min,
			MaxPrice: variantDetail.Max,
			Product: &dto.GetProductRes{
				ID:            variantDetail.ID,
				Price:         variantDetail.Min,
				Name:          variantDetail.Product.Name,
				Slug:          variantDetail.Product.Slug,
				MediaURL:      photoURL,
				City:          variantDetail.Product.Seller.Address.City,
				Rating:        variantDetail.Avg,
				TotalReviewer: variantDetail.Count,
				TotalSold:     uint(variantDetail.Product.SoldCount),
			},
		}
		productsRes = append(productsRes, dtoProduct)
	}

	return productsRes, totalPage, totalData, nil
}

func (p *productService) FindSimilarProducts(productID uint) ([]*dto.ProductRes, error) {
	tx := p.db.Begin()
	var err error
	defer helper.CommitOrRollback(tx, &err)

	product, err := p.productRepo.FindProductByID(tx, productID)
	if err != nil {
		return nil, err
	}

	products, err := p.productRepo.FindSimilarProduct(tx, product.CategoryID)
	if err != nil {
		return nil, err
	}

	var productsRes = make([]*dto.ProductRes, 0)
	for _, pdt := range products {
		if pdt.ID == productID {
			continue
		}

		var photoURL string
		if len(pdt.ProductPhotos) > 0 {
			photoURL = pdt.ProductPhotos[0].PhotoURL
		}

		var minPrice, maxPrice uint
		minPrice, maxPrice, err = p.productRepo.SearchMinMaxPrice(tx, pdt.ID)
		if err != nil {
			return nil, err
		}

		var ratings []int
		ratings, err = p.productRepo.SearchRating(tx, pdt.ID)
		if err != nil {
			return nil, err
		}
		reviewCount := len(ratings)
		avgRating := float64(helper.SumInt(ratings)) / float64(reviewCount)

		dtoProduct := &dto.ProductRes{
			MinPrice: float64(minPrice),
			MaxPrice: float64(maxPrice),
			Product: &dto.GetProductRes{
				ID:            pdt.ID,
				Price:         float64(minPrice),
				Name:          pdt.Name,
				Slug:          pdt.Slug,
				MediaURL:      photoURL,
				Rating:        avgRating,
				TotalReviewer: int64(reviewCount),
				TotalSold:     uint(pdt.SoldCount),
			},
		}
		productsRes = append(productsRes, dtoProduct)
	}

	return productsRes, nil
}

func (p *productService) GetProducts(query *repository.SearchQuery) ([]*dto.ProductRes, int64, int64, error) {
	tx := p.db.Begin()
	var err error
	defer helper.CommitOrRollback(tx, &err)

	variantDetails, totalPage, totalData, err := p.productVarDetRepo.SearchProducts(tx, query)
	if err != nil {
		return nil, 0, 0, err
	}

	var productsRes = make([]*dto.ProductRes, 0)
	for _, variantDetail := range variantDetails {
		var photoURL string
		if len(variantDetail.ProductPhotos) > 0 {
			photoURL = variantDetail.ProductPhotos[0].PhotoURL
		}

		dtoProduct := &dto.ProductRes{
			MinPrice: variantDetail.Min,
			MaxPrice: variantDetail.Max,
			Product: &dto.GetProductRes{
				ID:            variantDetail.ID,
				Price:         variantDetail.Min,
				Name:          variantDetail.Name,
				Slug:          variantDetail.Slug,
				MediaURL:      photoURL,
				City:          variantDetail.Seller.Address.City,
				Rating:        variantDetail.Avg,
				TotalReviewer: variantDetail.Count,
				TotalSold:     uint(variantDetail.Product.SoldCount),
			},
		}
		productsRes = append(productsRes, dtoProduct)
	}

	return productsRes, totalPage, totalData, nil
}

func (p *productService) SearchRecommendProduct(q *repository.SearchQuery) ([]*dto.ProductRes, int64, int64, error) {
	tx := p.db.Begin()
	var err error
	defer helper.CommitOrRollback(tx, &err)

	products, totalPage, totalData, err := p.productRepo.SearchRecommendProduct(tx, q)
	if err != nil {
		return nil, 0, 0, err
	}

	var productsRes = make([]*dto.ProductRes, 0)
	for _, product := range products {
		var photoURL string
		if len(product.ProductPhotos) > 0 {
			photoURL = product.ProductPhotos[0].PhotoURL
		}

		dtoProduct := &dto.ProductRes{
			MinPrice: product.Min,
			MaxPrice: product.Max,
			Product: &dto.GetProductRes{
				ID:            product.ID,
				Price:         product.Min,
				Name:          product.Name,
				Slug:          product.Slug,
				MediaURL:      photoURL,
				City:          product.Seller.Address.City,
				Rating:        product.Avg,
				TotalReviewer: product.Count,
				TotalSold:     uint(product.Product.SoldCount),
			},
		}
		productsRes = append(productsRes, dtoProduct)
	}
	return productsRes, totalPage, totalData, nil
}

func (p *productService) CreateSellerProduct(userID uint, req *dto.PostCreateProductReq) (*dto.PostCreateProductRes, error) {
	tx := p.db.Begin()
	var err error
	var err2 error
	var err3 error
	var err4 error
	var err5 error
	var err6 error
	defer helper.CommitOrRollback(tx, &err)
	defer helper.CommitOrRollback(tx, &err2)
	defer helper.CommitOrRollback(tx, &err3)
	defer helper.CommitOrRollback(tx, &err4)
	defer helper.CommitOrRollback(tx, &err5)
	defer helper.CommitOrRollback(tx, &err6)
	//get seller id
	seller, _ := p.sellerRepo.FindSellerByUserID(tx, userID)
	//create product
	product, err := p.productRepo.CreateProduct(tx, req.Name, req.CategoryID, seller.ID, req.IsBulkEnabled, req.MinQuantity, req.MaxQuantity)
	if err != nil {
		return nil, err
	}
	//create product details
	productDetail, err2 := p.productRepo.CreateProductDetail(tx, product.ID, req.ProductDetail)
	if err2 != nil {
		return nil, err2
	}
	//create product photos table
	var productPhotos []*model.ProductPhoto
	for _, ph := range req.ProductPhotos {
		productPhoto, err3 := p.productRepo.CreateProductPhoto(tx, product.ID, ph)
		if err3 != nil {
			return nil, err3
		}
		productPhotos = append(productPhotos, productPhoto)
	}
	var productVariant1 *model.ProductVariant
	var productVariant2 *model.ProductVariant
	if req.HasVariant {
		//create product variants
		productVariant1, err4 = p.productRepo.CreateProductVariant(tx, req.Variant1Name)
		if err4 != nil {
			return nil, err4
		}

		if req.Variant2Name != nil {
			productVariant2, err5 = p.productRepo.CreateProductVariant(tx, *req.Variant2Name)
			if err5 != nil {
				return nil, err5
			}
		} else {
			productVariant2 = nil
		}
	}
	//create product variant details
	productVariantDetail, err6 := p.productRepo.CreateProductVariantDetail(tx, product.ID, productVariant1.ID, productVariant2, req.ProductVariantDetails)
	if err6 != nil {
		return nil, err6
	}
	ret := dto.PostCreateProductRes{
		Product:              product,
		ProductDetail:        productDetail,
		ProductPhoto:         productPhotos,
		ProductVariant1:      productVariant1,
		ProductVariant2:      productVariant2,
		ProductVariantDetail: productVariantDetail,
	}
	return &ret, nil
}
