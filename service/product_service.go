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
	SearchProduct(q *repository.SearchQuery) (*dto.SearchedSortFilterProduct, error)
}

type productService struct {
	db          *gorm.DB
	productRepo repository.ProductRepository
}

type ProductConfig struct {
	DB          *gorm.DB
	ProductRepo repository.ProductRepository
}

func NewProductService(config *ProductConfig) ProductService {
	return &productService{db: config.DB, productRepo: config.ProductRepo}
}

func (s *productService) FindProductDetailBySlug(slug string) (*model.Product, error) {
	tx := s.db.Begin()

	product, err := s.productRepo.FindProductBySlug(tx, slug)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	product, err = s.productRepo.FindProductDetailByID(tx, product.ID)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	tx.Commit()
	return product, nil
}

func (p *productService) SearchProduct(q *repository.SearchQuery) (*dto.SearchedSortFilterProduct, error) {
	tx := p.db.Begin()

	if q.Search == "" {
		return nil, apperror.BadRequestError("Search required")
	}
	if q.SortBy == "" {
		q.SortBy = "bought"
	}
	if q.Sort == "" {
		q.Sort = "desc"
	}
	length, products, err := p.productRepo.SearchProduct(tx, q)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	var searchedProducts []dto.SearchedProductRes
	for _, product := range *products {
		pr := dto.SearchedProductRes{
			ProductID:  product.ID,
			Slug:       product.Slug,
			MediaURL:   "",
			MinPrice:   0,
			MaxPrice:   0,
			PromoPrice: product.Promotion.Amount,
			Rating:     0,
			Bought:     product.SoldCount,
			City:       "",
		}
		searchedProducts = append(searchedProducts, pr)
	}

	searchedSortFilterProducts := dto.SearchedSortFilterProduct{
		TotalLength:     length,
		SearchedProduct: searchedProducts,
	}
	return &searchedSortFilterProducts, nil
}
