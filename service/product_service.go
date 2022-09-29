package service

import (
	"errors"
	"gorm.io/gorm"
	"seadeals-backend/apperror"
	"seadeals-backend/dto"
	"seadeals-backend/helper"
	"seadeals-backend/model"
	"seadeals-backend/repository"
)

type ProductService interface {
	FindProductDetailBySlug(slug string, userID uint) (*model.Product, error)
	FindSimilarProducts(productID uint) ([]*dto.ProductRes, error)
	SearchRecommendProduct(q *repository.SearchQuery) (*dto.SearchedSortFilterProduct, error)
	GetProductsBySellerID(query *dto.SellerProductSearchQuery, sellerID uint) ([]*dto.ProductRes, int64, int64, error)
	GetProductsByCategoryID(query *dto.SellerProductSearchQuery, categoryID uint) ([]*dto.ProductRes, int64, int64, error)
	GetProducts(q *repository.SearchQuery) ([]*dto.ProductRes, int64, int64, error)
}

type productService struct {
	db                *gorm.DB
	productRepo       repository.ProductRepository
	reviewRepo        repository.ReviewRepository
	productVarDetRepo repository.ProductVariantDetailRepository
}

type ProductConfig struct {
	DB                *gorm.DB
	ProductRepo       repository.ProductRepository
	ReviewRepo        repository.ReviewRepository
	ProductVarDetRepo repository.ProductVariantDetailRepository
}

func NewProductService(config *ProductConfig) ProductService {
	return &productService{
		db:                config.DB,
		productRepo:       config.ProductRepo,
		reviewRepo:        config.ReviewRepo,
		productVarDetRepo: config.ProductVarDetRepo,
	}
}

func (p *productService) FindProductDetailBySlug(slug string, userID uint) (*model.Product, error) {
	tx := p.db.Begin()

	product, err := p.productRepo.FindProductBySlug(tx, slug)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	product, err = p.productRepo.FindProductDetailByID(tx, product.ID, userID)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	tx.Commit()
	return product, nil
}

func (p *productService) GetProductsBySellerID(query *dto.SellerProductSearchQuery, sellerID uint) ([]*dto.ProductRes, int64, int64, error) {
	tx := p.db.Begin()

	variantDetails, totalPage, totalData, err := p.productVarDetRepo.GetProductsBySellerID(tx, query, sellerID)
	if err != nil {
		tx.Rollback()
		return nil, 0, 0, err
	}
	if totalData == 0 {
		tx.Rollback()
		return nil, 0, 0, apperror.NotFoundError("No Products were found")
	}

	var productsRes []*dto.ProductRes
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

	tx.Commit()
	return productsRes, totalPage, totalData, nil
}

func (p *productService) GetProductsByCategoryID(query *dto.SellerProductSearchQuery, categoryID uint) ([]*dto.ProductRes, int64, int64, error) {
	tx := p.db.Begin()

	variantDetails, totalPage, totalData, err := p.productVarDetRepo.GetProductsByCategoryID(tx, query, categoryID)
	if err != nil {
		tx.Rollback()
		return nil, 0, 0, err
	}
	if totalData == 0 {
		tx.Rollback()
		return nil, 0, 0, apperror.NotFoundError("No Products were found")
	}

	var productsRes []*dto.ProductRes
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

	tx.Commit()
	return productsRes, totalPage, totalData, nil
}

func (p *productService) FindSimilarProducts(productID uint) ([]*dto.ProductRes, error) {
	tx := p.db.Begin()

	products, err := p.productRepo.FindSimilarProduct(tx, productID)
	if err != nil {
		tx.Rollback()
		if errors.Is(err, &apperror.ProductNotFoundError{}) {
			return nil, apperror.NotFoundError(err.Error())
		}
		return nil, err
	}

	var productsRes []*dto.ProductRes
	for _, product := range products {
		if product.ID == productID {
			continue
		}

		var photoURL string
		if len(product.ProductPhotos) > 0 {
			photoURL = product.ProductPhotos[0].PhotoURL
		}

		minPrice, maxPrice, err := p.productRepo.SearchMinMaxPrice(tx, product.ID)
		if err != nil {
			tx.Rollback()
			return nil, err
		}

		ratings, err := p.productRepo.SearchRating(tx, product.ID)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
		reviewCount := len(ratings)
		avgRating := float64(helper.SumInt(ratings)) / float64(reviewCount)

		dtoProduct := &dto.ProductRes{
			MinPrice: float64(minPrice),
			MaxPrice: float64(maxPrice),
			Product: &dto.GetProductRes{
				ID:            product.ID,
				Price:         float64(minPrice),
				Name:          product.Name,
				Slug:          product.Slug,
				MediaURL:      photoURL,
				Rating:        avgRating,
				TotalReviewer: int64(reviewCount),
				TotalSold:     uint(product.SoldCount),
			},
		}
		productsRes = append(productsRes, dtoProduct)
	}

	tx.Commit()
	return productsRes, nil
}

func (p *productService) GetProducts(query *repository.SearchQuery) ([]*dto.ProductRes, int64, int64, error) {
	tx := p.db.Begin()

	variantDetails, totalPage, totalData, err := p.productVarDetRepo.SearchProducts(tx, query)
	if err != nil {
		tx.Rollback()
		return nil, 0, 0, err
	}
	if totalData == 0 {
		tx.Rollback()
		return nil, 0, 0, apperror.NotFoundError("No Products were found")
	}

	var productsRes []*dto.ProductRes
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

	tx.Commit()
	return productsRes, totalPage, totalData, nil
}

func (p *productService) SearchRecommendProduct(q *repository.SearchQuery) (*dto.SearchedSortFilterProduct, error) {
	tx := p.db.Begin()

	products, err := p.productRepo.SearchRecommendProduct(tx, q)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	searchedSortFilterProducts := dto.SearchedSortFilterProduct{
		TotalLength:     len(products),
		SearchedProduct: products,
	}

	tx.Commit()
	return &searchedSortFilterProducts, nil
}
