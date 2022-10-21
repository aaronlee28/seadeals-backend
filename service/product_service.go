package service

import (
	"gorm.io/gorm"
	"math"
	"seadeals-backend/apperror"
	"seadeals-backend/dto"
	"seadeals-backend/helper"
	"seadeals-backend/model"
	"seadeals-backend/repository"
)

type ProductService interface {
	FindProductDetailByID(productID uint, userID uint) (*dto.ProductDetailRes, *dto.GetSellerRes, error)
	FindSimilarProducts(productID uint, query *repository.SearchQuery) ([]*dto.ProductRes, int64, int64, error)
	SearchRecommendProduct(q *repository.SearchQuery) ([]*dto.ProductRes, int64, int64, error)
	GetProductsBySellerID(query *dto.SellerProductSearchQuery, sellerID uint) ([]*dto.ProductRes, int64, int64, error)
	GetProductsByUserIDUnscoped(query *dto.SellerProductSearchQuery, userID uint) ([]*dto.GetSellerSummaryProductRes, int64, int64, error)
	GetProductsByCategoryID(query *dto.SellerProductSearchQuery, categoryID uint) ([]*dto.ProductRes, int64, int64, error)
	GetProducts(q *repository.SearchQuery) ([]*dto.ProductRes, int64, int64, error)

	CreateSellerProduct(userID uint, req *dto.PostCreateProductReq) (*dto.PostCreateProductRes, error)
	UpdateProductAndDetails(userID uint, productID uint, req *dto.PatchProductAndDetailsReq) (*dto.PatchProductAndDetailsRes, error)
	UpdateVariantAndDetails(userID uint, variantDetailsID uint, req *dto.PatchVariantAndDetails) (*dto.VariantAndDetailsUpdateRes, error)
	DeleteProductVariantDetails(userID uint, variantDetailsID uint, defaultPrice *float64) error
	AddVariantDetails(userID uint, productID uint, req *dto.AddVariantAndDetails) ([]*model.ProductVariantDetail, error)
	AddProductPhoto(userID uint, productID uint, req *dto.ProductPhotoReq) ([]*model.ProductPhoto, error)
	DeleteProductPhoto(userID uint, productID uint, req *dto.DeleteProductPhoto) ([]*model.ProductPhoto, error)
	DeleteProduct(userID uint, productID uint) (*model.Product, error)
}

type productService struct {
	db                *gorm.DB
	productRepo       repository.ProductRepository
	reviewRepo        repository.ReviewRepository
	productVarDetRepo repository.ProductVariantDetailRepository
	sellerRepo        repository.SellerRepository
	socialGraphRepo   repository.SocialGraphRepository
	notificationRepo  repository.NotificationRepository
}

type ProductConfig struct {
	DB                *gorm.DB
	ProductRepo       repository.ProductRepository
	ReviewRepo        repository.ReviewRepository
	ProductVarDetRepo repository.ProductVariantDetailRepository
	SellerRepo        repository.SellerRepository
	SocialGraphRepo   repository.SocialGraphRepository
	NotificationRepo  repository.NotificationRepository
}

func NewProductService(config *ProductConfig) ProductService {
	return &productService{
		db:                config.DB,
		productRepo:       config.ProductRepo,
		reviewRepo:        config.ReviewRepo,
		productVarDetRepo: config.ProductVarDetRepo,
		sellerRepo:        config.SellerRepo,
		socialGraphRepo:   config.SocialGraphRepo,
		notificationRepo:  config.NotificationRepo,
	}
}

func (p *productService) FindProductDetailByID(productID uint, userID uint) (*dto.ProductDetailRes, *dto.GetSellerRes, error) {
	tx := p.db.Begin()
	var err error
	defer helper.CommitOrRollback(tx, &err)

	product, err := p.productRepo.FindProductDetailByID(tx, productID, userID)
	if err != nil {
		return nil, nil, err
	}

	if len(product.ProductVariantDetail) > 0 {
		product.MinPrice = product.ProductVariantDetail[0].Price
		product.MaxPrice = product.ProductVariantDetail[len(product.ProductVariantDetail)-1].Price
	}

	seller, err := p.sellerRepo.FindSellerDetailByID(tx, product.SellerID, userID)
	if err != nil {
		return nil, nil, err
	}

	sellerRes := new(dto.GetSellerRes).From(seller)
	averageReview, totalReview, err := p.reviewRepo.GetReviewsAvgAndCountBySellerID(tx, seller.ID)
	if err != nil {
		return nil, nil, err
	}

	ratio := math.Pow(10, float64(1))
	RoundedAvgRating := math.Round(averageReview*ratio) / ratio

	sellerRes.TotalReviewer = uint(totalReview)
	sellerRes.Rating = RoundedAvgRating

	followers, err := p.socialGraphRepo.GetFollowerCountBySellerID(tx, seller.ID)
	if err != nil {
		return nil, nil, err
	}
	sellerRes.Followers = uint(followers)

	following, err := p.socialGraphRepo.GetFollowingCountByUserID(tx, seller.UserID)
	if err != nil {
		return nil, nil, err
	}
	sellerRes.Following = uint(following)

	totalProduct, err := p.productRepo.GetProductCountBySellerID(tx, seller.ID)
	if err != nil {
		return nil, nil, err
	}
	sellerRes.TotalProduct = totalProduct
	product.ID = productID

	return product, sellerRes, nil
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
			MinPriceBeforeDisc: variantDetail.MinBeforeDisc,
			MaxPriceBeforeDisc: variantDetail.MaxBeforeDisc,
			MinPrice:           variantDetail.Min,
			MaxPrice:           variantDetail.Max,
			Product: &dto.GetProductRes{
				ID:              variantDetail.ID,
				Price:           variantDetail.Min,
				Name:            variantDetail.Product.Name,
				Slug:            variantDetail.Product.Slug,
				MediaURL:        photoURL,
				City:            variantDetail.Product.Seller.Address.City,
				Rating:          variantDetail.Avg,
				TotalReviewer:   variantDetail.Count,
				PromotionAmount: variantDetail.PromotionAmount,
				TotalSold:       uint(variantDetail.Product.SoldCount),
			},
		}
		productsRes = append(productsRes, dtoProduct)
	}

	return productsRes, totalPage, totalData, nil
}

func (p *productService) GetProductsByUserIDUnscoped(query *dto.SellerProductSearchQuery, userID uint) ([]*dto.GetSellerSummaryProductRes, int64, int64, error) {
	tx := p.db.Begin()
	var err error
	defer helper.CommitOrRollback(tx, &err)

	seller, err := p.sellerRepo.FindSellerByUserID(tx, userID)
	if err != nil {
		return nil, 0, 0, err
	}

	products, totalData, err := p.productVarDetRepo.GetProductsBySellerIDUnscoped(tx, query, seller.ID)
	if err != nil {
		return nil, 0, 0, err
	}

	totalPage := (totalData + int64(query.Limit) - 1) / int64(query.Limit)

	var productsRes = make([]*dto.GetSellerSummaryProductRes, 0)
	for _, product := range products {
		dtoProduct := new(dto.GetSellerSummaryProductRes).From(product)
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
			MinPriceBeforeDisc: variantDetail.MinBeforeDisc,
			MaxPriceBeforeDisc: variantDetail.MaxBeforeDisc,
			MinPrice:           variantDetail.Min,
			MaxPrice:           variantDetail.Max,
			Product: &dto.GetProductRes{
				ID:              variantDetail.ID,
				Price:           variantDetail.Min,
				Name:            variantDetail.Product.Name,
				Slug:            variantDetail.Product.Slug,
				MediaURL:        photoURL,
				City:            variantDetail.Product.Seller.Address.City,
				Rating:          variantDetail.Avg,
				TotalReviewer:   variantDetail.Count,
				PromotionAmount: variantDetail.PromotionAmount,
				TotalSold:       uint(variantDetail.Product.SoldCount),
			},
		}
		productsRes = append(productsRes, dtoProduct)
	}

	return productsRes, totalPage, totalData, nil
}

func (p *productService) FindSimilarProducts(productID uint, query *repository.SearchQuery) ([]*dto.ProductRes, int64, int64, error) {
	tx := p.db.Begin()
	var err error
	defer helper.CommitOrRollback(tx, &err)

	product, err := p.productRepo.FindProductByID(tx, productID)
	if err != nil {
		return nil, 0, 0, err
	}

	products, totalPage, _, err := p.productRepo.FindSimilarProduct(tx, product.CategoryID, query)
	if err != nil {
		return nil, 0, 0, err
	}

	var productsRes = make([]*dto.ProductRes, 0)
	for _, pdt := range products {
		if pdt.ID == productID {
			continue
		}

		imageURL := ""
		if len(pdt.ProductPhotos) > 0 {
			imageURL = pdt.ProductPhotos[0].PhotoURL
		}
		dtoProduct := &dto.ProductRes{
			MinPriceBeforeDisc: pdt.MinBeforeDisc,
			MaxPriceBeforeDisc: pdt.MaxBeforeDisc,
			MinPrice:           pdt.Min,
			MaxPrice:           pdt.Max,
			Product: &dto.GetProductRes{
				ID:              pdt.ID,
				Price:           pdt.Min,
				Name:            pdt.Name,
				Slug:            pdt.Slug,
				MediaURL:        imageURL,
				Rating:          pdt.Avg,
				TotalReviewer:   pdt.Count,
				TotalSold:       uint(pdt.Product.SoldCount),
				PromotionAmount: pdt.PromotionAmount,
				City:            pdt.Seller.Address.City,
			},
		}
		productsRes = append(productsRes, dtoProduct)
	}

	return productsRes, totalPage, int64(len(productsRes)), nil
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
			MinPriceBeforeDisc: variantDetail.MinBeforeDisc,
			MaxPriceBeforeDisc: variantDetail.MaxBeforeDisc,
			MinPrice:           variantDetail.Min,
			MaxPrice:           variantDetail.Max,
			Product: &dto.GetProductRes{
				ID:              variantDetail.ID,
				Price:           variantDetail.Min,
				Name:            variantDetail.Name,
				Slug:            variantDetail.Slug,
				MediaURL:        photoURL,
				City:            variantDetail.Seller.Address.City,
				Rating:          variantDetail.Avg,
				TotalReviewer:   variantDetail.Count,
				PromotionAmount: variantDetail.PromotionAmount,
				TotalSold:       uint(variantDetail.Product.SoldCount),
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
			MinPriceBeforeDisc: product.MinBeforeDisc,
			MaxPriceBeforeDisc: product.MaxBeforeDisc,
			MinPrice:           product.Min,
			MaxPrice:           product.Max,
			Product: &dto.GetProductRes{
				ID:              product.ID,
				Price:           product.Min,
				Name:            product.Name,
				Slug:            product.Slug,
				MediaURL:        photoURL,
				City:            product.Seller.Address.City,
				Rating:          product.Avg,
				TotalReviewer:   product.Count,
				PromotionAmount: product.PromotionAmount,
				TotalSold:       uint(product.Product.SoldCount),
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
	var seller *model.Seller
	seller, err = p.sellerRepo.FindSellerByUserID(tx, userID)
	if err != nil {
		return nil, err
	}
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
		var productVariant1 *model.ProductVariant
		var productVariant2 *model.ProductVariant
		productVariant1, err = p.productRepo.CreateProductVariant(tx, *req.Variant1Name)
		if err != nil {
			return nil, err
		}
		if req.Variant2Name != nil {
			productVariant2, err = p.productRepo.CreateProductVariant(tx, *req.Variant2Name)
			if err != nil {
				return nil, err
			}
		} else {
			productVariant2 = nil
		}
		for _, v := range req.VariantArray {
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
	var userArray []*model.SocialGraph
	userArray, err = p.socialGraphRepo.GetFollowerUserID(tx, seller.ID)
	for _, user := range userArray {
		newNotification := &model.Notification{
			UserID:   user.UserID,
			SellerID: seller.ID,
			Title:    dto.NotificationFollowProduk,
			Detail:   "Seller adds new product",
		}
		p.notificationRepo.AddToNotificationFromModel(tx, newNotification)
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

func (p *productService) UpdateVariantAndDetails(userID uint, variantDetailsID uint, req *dto.PatchVariantAndDetails) (*dto.VariantAndDetailsUpdateRes, error) {
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
	ret := &dto.VariantAndDetailsUpdateRes{
		Variant1Name:          &updatedProductVariant1.Name,
		Variant2Name:          &updatedProductVariant2.Name,
		ProductVariantDetails: pvdRet,
	}
	var userArray []*model.Favorite
	var title string
	var title2 string
	var detail string
	var detail2 string
	if req.ProductVariantDetails.Stock != 0 {
		title = dto.NotificationFavoriteStok
		detail = "Favorite item stock change"
	}
	if productVariantDetails.Price == 0 && req.ProductVariantDetails.Price != 0 {
		title = dto.NotificationFavoriteHarga
		detail = "Favorite item price change"
	}
	userArray, err = p.socialGraphRepo.GetFavoriteUserID(tx, updatedProductVariantDetails.ProductID)

	if productVariantDetails.Price == 0 && req.ProductVariantDetails.Price != 0 && req.ProductVariantDetails.Stock != 0 {
		title = dto.NotificationFavoriteStok
		detail = "Favorite item stock change"
		title2 = dto.NotificationFavoriteHarga
		detail2 = "Favorite item price change"
		for _, user := range userArray {
			newNotificationStock := &model.Notification{
				UserID:   user.UserID,
				SellerID: seller.ID,
				Title:    title,
				Detail:   detail,
			}

			p.notificationRepo.AddToNotificationFromModel(tx, newNotificationStock)
			newNotificationHarga := &model.Notification{
				UserID:   user.UserID,
				SellerID: seller.ID,
				Title:    title2,
				Detail:   detail2,
			}

			p.notificationRepo.AddToNotificationFromModel(tx, newNotificationHarga)
		}
	} else {
		for _, user := range userArray {
			newNotification := &model.Notification{
				UserID:   user.UserID,
				SellerID: seller.ID,
				Title:    title,
				Detail:   detail,
			}
			p.notificationRepo.AddToNotificationFromModel(tx, newNotification)
		}
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
	if len(pvds) == 1 && defaultPrice != nil {
		createPVD := &model.ProductVariantDetail{
			ProductID:     checkPID.ID,
			Price:         *defaultPrice,
			Variant1Value: nil,
			Variant2Value: nil,
			Variant1ID:    nil,
			Variant2ID:    nil,
			VariantCode:   nil,
			PictureURL:    nil,
			Stock:         0,
		}
		_, err = p.productRepo.CreateProductVariantDetailWithModel(tx, createPVD)
		if err != nil {
			return err
		}
	}
	err = p.productRepo.DeleteProductVariantDetailsByID(tx, variantDetailsID)
	if err != nil {
		return err
	}
	return nil
}

//add product variant detail

func (p *productService) AddVariantDetails(userID uint, productID uint, req *dto.AddVariantAndDetails) ([]*model.ProductVariantDetail, error) {
	tx := p.db.Begin()
	var err error
	var err2 error

	defer helper.CommitOrRollback(tx, &err)
	defer helper.CommitOrRollback(tx, &err2)

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

	err = p.productRepo.DeleteNullProductVariantDetailsByID(tx, productID)

	var productVariant1 *model.ProductVariant
	var productVariant2 *model.ProductVariant
	productVariant1, err = p.productRepo.GetVariantByName(tx, *req.Variant1Name)
	if err != nil {

		productVariant1, err2 = p.productRepo.CreateVariantWithName(tx, *req.Variant1Name)
		if err2 != nil {
			return nil, err
		}
	}
	if req.Variant2Name != nil {
		productVariant2, err = p.productRepo.GetVariantByName(tx, *req.Variant2Name)

		if err != nil {
			productVariant2, err2 = p.productRepo.CreateVariantWithName(tx, *req.Variant2Name)
			if err2 != nil {
				return nil, err
			}
		}
	}
	var createdVariantDetail []*model.ProductVariantDetail
	for _, pvd := range req.ProductVariantDetails {

		addPVD := model.ProductVariantDetail{
			ProductID:     productID,
			Price:         pvd.Price,
			Variant1Value: &productVariant1.Name,
			Variant2Value: &productVariant2.Name,
			Variant1ID:    productVariant1.ID,
			Variant2ID:    productVariant2.ID,
			VariantCode:   pvd.VariantCode,
			PictureURL:    pvd.PictureURL,
			Stock:         pvd.Stock,
		}
		var newPVD *model.ProductVariantDetail
		newPVD, err = p.productRepo.CreateProductVariantDetailWithModel(tx, &addPVD)
		createdVariantDetail = append(createdVariantDetail, newPVD)
	}

	return createdVariantDetail, nil
}

//add product photo
func (p *productService) AddProductPhoto(userID uint, productID uint, req *dto.ProductPhotoReq) ([]*model.ProductPhoto, error) {
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
		err = apperror.BadRequestError("Product does not belong to seller")
		return nil, err
	}

	productPhotos, err := p.productRepo.CreateProductPhotos(tx, productID, req)
	if err != nil {
		return nil, err
	}
	return productPhotos, nil
}

//delete product photo
func (p *productService) DeleteProductPhoto(userID uint, productID uint, req *dto.DeleteProductPhoto) ([]*model.ProductPhoto, error) {
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
		err = apperror.BadRequestError("Product does not belong to seller")
		return nil, err
	}

	productPhotos, err := p.productRepo.DeleteProductPhotos(tx, req)
	if err != nil {
		return nil, err
	}
	return productPhotos, nil
}

//delete product
func (p *productService) DeleteProduct(userID uint, productID uint) (*model.Product, error) {
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
		err = apperror.BadRequestError("Product does not belong to seller")
		return nil, err
	}

	product, err := p.productRepo.DeleteProduct(tx, productID)
	if err != nil {
		return nil, err
	}
	return product, nil
}
