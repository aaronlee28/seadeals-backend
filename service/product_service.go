package service

import (
	"gorm.io/gorm"
	"seadeals-backend/apperror"
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
	UpdateProductAndDetails(userID uint, productID uint, req *dto.PatchProductAndDetailsReq) (*dto.PatchProductAndDetailsRes, error)
	UpdateVariantAndDetails(userID uint, variantDetailsID uint, req *dto.PatchVariantAndDetails) (*dto.VariantAndDetails, error)
	DeleteProductVariantDetails(userID uint, variantDetailsID uint, defaultPrice *float64) error
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

	defer helper.CommitOrRollback(tx, &err)

	if req.DefaultPrice == nil && len(req.VariantArray) == 0 && req.DefaultPrice == nil {
		err = apperror.BadRequestError("default price is required if there is no variant")
	}

	//get seller id
	seller, _ := p.sellerRepo.FindSellerByUserID(tx, userID)
	//create product
	var product *model.Product
	product, err = p.productRepo.CreateProduct(tx, req.Name, req.CategoryID, seller.ID, req.IsBulkEnabled, req.MinQuantity, req.MaxQuantity)
	if err != nil {
		return nil, err
	}
	//create product details
	var productDetail *model.ProductDetail
	productDetail, err = p.productRepo.CreateProductDetail(tx, product.ID, req.ProductDetail)
	if err != nil {
		return nil, err
	}
	//create product photos table
	var productPhotos []*model.ProductPhoto
	for _, ph := range req.ProductPhotos {
		var productPhoto *model.ProductPhoto
		productPhoto, err = p.productRepo.CreateProductPhoto(tx, product.ID, ph)
		if err != nil {
			return nil, err
		}
		productPhotos = append(productPhotos, productPhoto)
	}
	var productVariantDetail *model.ProductVariantDetail
	var productVariantDetails []*model.ProductVariantDetail
	if len(req.VariantArray) == 0 {
		defaultProductVariantDetail := dto.ProductVariantDetail{
			Price:         *req.DefaultPrice,
			Variant1Value: nil,
			Variant2Value: nil,
			VariantCode:   nil,
			PictureURL:    nil,
			Stock:         *req.DefaultStock,
		}
		productVariantDetail, err = p.productRepo.CreateProductVariantDetail(tx, product.ID, nil, nil, &defaultProductVariantDetail)
		productVariantDetails = append(productVariantDetails, productVariantDetail)
		if err != nil {
			return nil, err
		}
	}
	//create product variant details
	if len(req.VariantArray) > 0 {
		for _, v := range req.VariantArray {
			var productVariant1 *model.ProductVariant
			var productVariant2 *model.ProductVariant
			productVariant1, err = p.productRepo.CreateProductVariant(tx, *v.Variant1Name)
			if err != nil {
				return nil, err
			}
			if v.Variant2Name != nil {
				productVariant2, err = p.productRepo.CreateProductVariant(tx, *v.Variant1Name)
				if err != nil {
					return nil, err
				}
			} else {
				productVariant2 = nil
			}
			productVariantDetail, err = p.productRepo.CreateProductVariantDetail(tx, product.ID, productVariant1.ID, productVariant2, v.ProductVariantDetails)
			if err != nil {
				return nil, err
			}
			productVariantDetails = append(productVariantDetails, productVariantDetail)
		}
	}

	ret := dto.PostCreateProductRes{
		Product:              product,
		ProductDetail:        productDetail,
		ProductPhoto:         productPhotos,
		ProductVariantDetail: productVariantDetails,
	}
	return &ret, nil
}
func (p *productService) UpdateProductAndDetails(userID uint, productID uint, req *dto.PatchProductAndDetailsReq) (*dto.PatchProductAndDetailsRes, error) {
	tx := p.db.Begin()
	var err error

	defer helper.CommitOrRollback(tx, &err)

	var seller *model.Seller
	seller, err = p.sellerRepo.FindSellerByUserID(tx, userID)
	if err != nil {
		return nil, err
	}
	var checkPID *model.Product
	checkPID, err = p.productRepo.FindProductByID(tx, productID)

	if seller.ID != checkPID.SellerID {
		err = apperror.BadRequestError("Product does not belong to seller ")
		return nil, err
	}
	product := model.Product{
		Name:          req.Product.Name,
		IsBulkEnabled: req.Product.IsBulkEnabled,
		MinQuantity:   req.Product.MinQuantity,
		MaxQuantity:   req.Product.MaxQuantity,
	}

	var updatedProduct *model.Product
	updatedProduct, err = p.productRepo.UpdateProduct(tx, productID, &product)

	productDetail := model.ProductDetail{
		Description:     req.ProductDetail.Description,
		VideoURL:        req.ProductDetail.VideoURL,
		IsHazardous:     req.ProductDetail.IsHazardous,
		ConditionStatus: req.ProductDetail.ConditionStatus,
		Length:          req.ProductDetail.Length,
		Width:           req.ProductDetail.Width,
		Height:          req.ProductDetail.Height,
		Weight:          req.ProductDetail.Weight,
	}
	var updatedProductDetail *model.ProductDetail
	updatedProductDetail, err = p.productRepo.UpdateProductDetail(tx, productID, &productDetail)

	res := dto.PatchProductAndDetailsRes{
		Product:       updatedProduct,
		ProductDetail: updatedProductDetail,
	}

	return &res, nil
}

//when update, check kalo ada variant yang nil, kalau ada, delete variant null
func (p *productService) UpdateVariantAndDetails(userID uint, variantDetailsID uint, req *dto.PatchVariantAndDetails) (*dto.VariantAndDetails, error) {
	tx := p.db.Begin()
	var err error

	defer helper.CommitOrRollback(tx, &err)

	var seller *model.Seller
	seller, err = p.sellerRepo.FindSellerByUserID(tx, userID)
	if err != nil {
		return nil, err
	}
	var productVariantDetails *model.ProductVariantDetail
	productVariantDetails, err = p.productRepo.FindProductVariantDetailsByID(tx, variantDetailsID)
	var checkPID *model.Product
	checkPID, err = p.productRepo.FindProductByID(tx, productVariantDetails.ProductID)
	if seller.ID != checkPID.SellerID {
		err = apperror.BadRequestError("Product does not belong to seller ")
		return nil, err
	}
	var updatedProductVariant1 *model.ProductVariant
	var updatedProductVariant2 *model.ProductVariant

	if req.Variant1Name != nil {
		updateProductVariant := &model.ProductVariant{
			Name: *req.Variant1Name,
		}
		updatedProductVariant1, err = p.productRepo.UpdateProductVariantByID(tx, *productVariantDetails.Variant1ID, updateProductVariant)
		if err != nil {
			return nil, err
		}
	}
	if req.Variant2Name != nil {
		updateProductVariant := &model.ProductVariant{
			Name: *req.Variant2Name,
		}
		updatedProductVariant2, err = p.productRepo.UpdateProductVariantByID(tx, *productVariantDetails.Variant1ID, updateProductVariant)
		if err != nil {
			return nil, err
		}
	}
	updateProductVariantDetail := &model.ProductVariantDetail{
		Price:         req.ProductVariantDetails.Price,
		Variant1Value: req.ProductVariantDetails.Variant1Value,
		Variant2Value: req.ProductVariantDetails.Variant2Value,
		VariantCode:   req.ProductVariantDetails.VariantCode,
		PictureURL:    req.ProductVariantDetails.PictureURL,
		Stock:         req.ProductVariantDetails.Stock,
	}
	var updatedProductVariantDetails *model.ProductVariantDetail

	updatedProductVariantDetails, err = p.productRepo.UpdateProductVariantDetailByID(tx, variantDetailsID, updateProductVariantDetail)
	if err != nil {
		return nil, err
	}
	pvdRet := &dto.ProductVariantDetail{
		Price:         updatedProductVariantDetails.Price,
		Variant1Value: updatedProductVariantDetails.Variant1Value,
		Variant2Value: updatedProductVariantDetails.Variant2Value,
		VariantCode:   updatedProductVariantDetails.VariantCode,
		PictureURL:    updatedProductVariantDetails.PictureURL,
		Stock:         updatedProductVariantDetails.Stock,
	}
	ret := &dto.VariantAndDetails{
		Variant1Name:          &updatedProductVariant1.Name,
		Variant2Name:          &updatedProductVariant2.Name,
		ProductVariantDetails: pvdRet,
	}
	return ret, nil
}

//when delete, check kalo product variant = 1, kalo 1 then add default price variant

func (p *productService) DeleteProductVariantDetails(userID uint, variantDetailsID uint, defaultPrice *float64) error {
	tx := p.db.Begin()
	var err error

	defer helper.CommitOrRollback(tx, &err)

	var seller *model.Seller
	seller, err = p.sellerRepo.FindSellerByUserID(tx, userID)
	if err != nil {
		return err
	}
	var productVariantDetails *model.ProductVariantDetail
	productVariantDetails, err = p.productRepo.FindProductVariantDetailsByID(tx, variantDetailsID)
	var checkPID *model.Product
	checkPID, err = p.productRepo.FindProductByID(tx, productVariantDetails.ProductID)
	if seller.ID != checkPID.SellerID {
		err = apperror.BadRequestError("Product does not belong to seller ")
		return err
	}
	var pvds []*model.ProductVariantDetail

	pvds, err = p.productRepo.FindProductVariantDetailsByProductID(tx, productVariantDetails.ProductID)
	if len(pvds) == 1 && defaultPrice == nil {
		err = apperror.BadRequestError("default price is required")
		return err
	}
	err = p.productRepo.DeleteProductVariantDetailsByID(tx, variantDetailsID)
	if err != nil {
		return err
	}
	return nil
}

//add product variant detail
//func (p *productService) CreateVariantAndDetails(userID uint, variantDetailsID, req *dto.VariantAndDetails) (*dto.VariantAndDetails, error) {
//	tx := p.db.Begin()
//	var err error
//
//	defer helper.CommitOrRollback(tx, &err)
//
//	var seller *model.Seller
//	seller, err = p.sellerRepo.FindSellerByUserID(tx, userID)
//	if err != nil {
//		return nil, err
//	}
//	var productVariantDetails *model.ProductVariantDetail
//	productVariantDetails, err = p.productRepo.FindProductVariantDetailsByID(tx, variantDetailsID)
//	var checkPID *model.Product
//	checkPID, err = p.productRepo.FindProductByID(tx, productVariantDetails.ProductID)
//	if seller.ID != checkPID.SellerID {
//		err = apperror.BadRequestError("Product does not belong to seller ")
//		return nil, err
//	}
//	var updatedProductVariant1 *model.ProductVariant
//	var updatedProductVariant2 *model.ProductVariant
//
//	if req.Variant1Name != nil {
//		updateProductVariant := &model.ProductVariant{
//			Name: *req.Variant1Name,
//		}
//		updatedProductVariant1, err = p.productRepo.UpdateProductVariantByID(tx, *productVariantDetails.Variant1ID, updateProductVariant)
//		if err != nil {
//			return nil, err
//		}
//	}
//	if req.Variant2Name != nil {
//		updateProductVariant := &model.ProductVariant{
//			Name: *req.Variant2Name,
//		}
//		updatedProductVariant2, err = p.productRepo.UpdateProductVariantByID(tx, *productVariantDetails.Variant1ID, updateProductVariant)
//		if err != nil {
//			return nil, err
//		}
//	}
//	updateProductVariantDetail := &model.ProductVariantDetail{
//		Price:         req.ProductVariantDetails.Price,
//		Variant1Value: req.ProductVariantDetails.Variant1Value,
//		Variant2Value: req.ProductVariantDetails.Variant2Value,
//		VariantCode:   req.ProductVariantDetails.VariantCode,
//		PictureURL:    req.ProductVariantDetails.PictureURL,
//		Stock:         req.ProductVariantDetails.Stock,
//	}
//	var updatedProductVariantDetails *model.ProductVariantDetail
//
//	updatedProductVariantDetails, err = p.productRepo.UpdateProductVariantDetailByID(tx, variantDetailsID, updateProductVariantDetail)
//	if err != nil {
//		return nil, err
//	}
//	pvdRet := &dto.ProductVariantDetail{
//		Price:         updatedProductVariantDetails.Price,
//		Variant1Value: updatedProductVariantDetails.Variant1Value,
//		Variant2Value: updatedProductVariantDetails.Variant2Value,
//		VariantCode:   updatedProductVariantDetails.VariantCode,
//		PictureURL:    updatedProductVariantDetails.PictureURL,
//		Stock:         updatedProductVariantDetails.Stock,
//	}
//	ret := &dto.VariantAndDetails{
//		Variant1Name:          &updatedProductVariant1.Name,
//		Variant2Name:          &updatedProductVariant2.Name,
//		ProductVariantDetails: pvdRet,
//	}
//	return ret, nil
//}
//update product photo

//delete product
