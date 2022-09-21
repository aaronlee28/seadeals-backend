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
	products, err := p.productRepo.SearchProduct(tx, q)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	var searchedProducts []*dto.SearchedProductRes
	for _, product := range *products {
		pr := new(dto.SearchedProductRes).FromProduct(&product)
		searchedProducts = append(searchedProducts, pr)
	}
	length := len(*products)
	searchedSortFilterProducts := dto.SearchedSortFilterProduct{
		TotalLength:     length,
		SearchedProduct: searchedProducts,
	}
	for _, product := range searchedSortFilterProducts.SearchedProduct {
		url, err2 := p.productRepo.SearchImageURL(tx, product.ProductID)

		min, max, err3 := p.productRepo.SearchMinMaxPrice(tx, product.ProductID)
		product.MediaURL = url
		product.MinPrice = min
		product.MaxPrice = max

		if err2 != nil {
			tx.Rollback()
			return nil, err2
		}
		if err3 != nil {
			tx.Rollback()
			return nil, err3
		}
	}
	return &searchedSortFilterProducts, nil
}
