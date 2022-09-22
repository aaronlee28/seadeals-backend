package service

import (
	"gorm.io/gorm"
	"seadeals-backend/apperror"
	"seadeals-backend/dto"
	"seadeals-backend/model"
	"seadeals-backend/repository"
)

type ProductService interface {
	FindProductDetailBySlug(slug string) (*model.Product, error)
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

func (p *productService) FindProductDetailBySlug(slug string) (*model.Product, error) {
	tx := p.db.Begin()

	product, err := p.productRepo.FindProductBySlug(tx, slug)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	product, err = p.productRepo.FindProductDetailByID(tx, product.ID)
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
				PictureURL:    photoURL,
				City:          variantDetail.Product.Seller.Address.SubDistrict.District.City.Name,
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

func (s *productService) GetProductsByCategoryID(query *dto.SellerProductSearchQuery, categoryID uint) ([]*dto.ProductRes, int64, int64, error) {
	tx := s.db.Begin()

	variantDetails, totalPage, totalData, err := s.productVarDetRepo.GetProductsByCategoryID(tx, query, categoryID)
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
				PictureURL:    photoURL,
				City:          variantDetail.Product.Seller.Address.SubDistrict.District.City.Name,
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

func (s *productService) GetProducts(query *repository.SearchQuery) ([]*dto.ProductRes, int64, int64, error) {
	tx := s.db.Begin()

	variantDetails, totalPage, totalData, err := s.productVarDetRepo.SearchProducts(tx, query)
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
				PictureURL:    photoURL,
				City:          variantDetail.Product.Seller.Address.SubDistrict.District.City.Name,
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

	return &searchedSortFilterProducts, nil
}
