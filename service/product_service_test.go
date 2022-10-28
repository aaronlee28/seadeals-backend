package service_test

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"seadeals-backend/dto"
	"seadeals-backend/mocks"
	"seadeals-backend/model"
	"seadeals-backend/repository"
	"seadeals-backend/service"
	"seadeals-backend/testutil"
	"testing"
)

func TestProductService_FindProductDetailByID(t *testing.T) {
	t.Run("Should return response body", func(t *testing.T) {

		gormDB := testutil.MockDB()
		mockRepo1 := new(mocks.ProductRepository)
		mockRepo2 := new(mocks.ReviewRepository)
		mockRepo3 := new(mocks.ProductVariantDetailRepository)
		mockRepo4 := new(mocks.SellerRepository)
		mockRepo5 := new(mocks.SocialGraphRepository)
		mockRepo6 := new(mocks.NotificationRepository)
		cfg := &service.ProductConfig{
			DB:                gormDB,
			ProductRepo:       mockRepo1,
			ReviewRepo:        mockRepo2,
			ProductVarDetRepo: mockRepo3,
			SellerRepo:        mockRepo4,
			SocialGraphRepo:   mockRepo5,
			NotificationRepo:  mockRepo6,
		}
		s := service.NewProductService(cfg)

		mockAddress := &model.Address{City: "test"}
		mockSeller := &model.Seller{Address: mockAddress}
		mockProduct := &model.Product{Seller: mockSeller}

		mockArray := []*dto.SellerProductsCustomTable{}
		MockSellerCustomTable := &dto.SellerProductsCustomTable{Product: *mockProduct}
		mockArray = append(mockArray, MockSellerCustomTable)
		mockRepo3.On("GetProductsBySellerID", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("*dto.SellerProductSearchQuery"), mock.AnythingOfType("uint")).Return(mockArray, int64(0), int64(0), nil)

		res, _, _, err := s.GetProductsBySellerID(&dto.SellerProductSearchQuery{}, uint(1))

		assert.Nil(t, err)
		assert.NotNil(t, res)
	})

	t.Run("Should return error", func(t *testing.T) {

		gormDB := testutil.MockDB()
		mockRepo1 := new(mocks.ProductRepository)
		mockRepo2 := new(mocks.ReviewRepository)
		mockRepo3 := new(mocks.ProductVariantDetailRepository)
		mockRepo4 := new(mocks.SellerRepository)
		mockRepo5 := new(mocks.SocialGraphRepository)
		mockRepo6 := new(mocks.NotificationRepository)
		cfg := &service.ProductConfig{
			DB:                gormDB,
			ProductRepo:       mockRepo1,
			ReviewRepo:        mockRepo2,
			ProductVarDetRepo: mockRepo3,
			SellerRepo:        mockRepo4,
			SocialGraphRepo:   mockRepo5,
			NotificationRepo:  mockRepo6,
		}
		s := service.NewProductService(cfg)

		mockAddress := &model.Address{City: "test"}
		mockSeller := &model.Seller{Address: mockAddress}
		mockProduct := &model.Product{Seller: mockSeller}

		mockArray := []*dto.SellerProductsCustomTable{}
		MockSellerCustomTable := &dto.SellerProductsCustomTable{Product: *mockProduct}
		mockArray = append(mockArray, MockSellerCustomTable)
		mockRepo3.On("GetProductsBySellerID", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("*dto.SellerProductSearchQuery"), mock.AnythingOfType("uint")).Return(nil, int64(0), int64(0), errors.New(""))

		res, _, _, err := s.GetProductsBySellerID(&dto.SellerProductSearchQuery{}, uint(1))

		assert.Nil(t, res)
		assert.NotNil(t, err)
	})

	t.Run("Should return response body", func(t *testing.T) {

		gormDB := testutil.MockDB()
		mockRepo1 := new(mocks.ProductRepository)
		mockRepo2 := new(mocks.ReviewRepository)
		mockRepo3 := new(mocks.ProductVariantDetailRepository)
		mockRepo4 := new(mocks.SellerRepository)
		mockRepo5 := new(mocks.SocialGraphRepository)
		mockRepo6 := new(mocks.NotificationRepository)
		cfg := &service.ProductConfig{
			DB:                gormDB,
			ProductRepo:       mockRepo1,
			ReviewRepo:        mockRepo2,
			ProductVarDetRepo: mockRepo3,
			SellerRepo:        mockRepo4,
			SocialGraphRepo:   mockRepo5,
			NotificationRepo:  mockRepo6,
		}
		s := service.NewProductService(cfg)

		mockAddress := &model.Address{City: "test"}
		mockSeller := &model.Seller{Address: mockAddress}
		mockProductPhotoArr := []*model.ProductPhoto{}
		mockProductPhoto := &model.ProductPhoto{PhotoURL: "test"}
		mockProductPhotoArr = append(mockProductPhotoArr, mockProductPhoto)
		mockProduct := &model.Product{Seller: mockSeller, ProductPhotos: mockProductPhotoArr}

		mockArray := []*dto.SellerProductsCustomTable{}
		MockSellerCustomTable := &dto.SellerProductsCustomTable{Product: *mockProduct}
		mockArray = append(mockArray, MockSellerCustomTable)
		mockRepo3.On("GetProductsBySellerID", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("*dto.SellerProductSearchQuery"), mock.AnythingOfType("uint")).Return(mockArray, int64(0), int64(0), nil)

		res, _, _, err := s.GetProductsBySellerID(&dto.SellerProductSearchQuery{}, uint(1))

		assert.Nil(t, err)
		assert.NotNil(t, res)
	})
}

func TestProductService_GetProductsByUserIDUnscoped(t *testing.T) {
	t.Run("Should return response body", func(t *testing.T) {

		gormDB := testutil.MockDB()
		mockRepo1 := new(mocks.ProductRepository)
		mockRepo2 := new(mocks.ReviewRepository)
		mockRepo3 := new(mocks.ProductVariantDetailRepository)
		mockRepo4 := new(mocks.SellerRepository)
		mockRepo5 := new(mocks.SocialGraphRepository)
		mockRepo6 := new(mocks.NotificationRepository)
		cfg := &service.ProductConfig{
			DB:                gormDB,
			ProductRepo:       mockRepo1,
			ReviewRepo:        mockRepo2,
			ProductVarDetRepo: mockRepo3,
			SellerRepo:        mockRepo4,
			SocialGraphRepo:   mockRepo5,
			NotificationRepo:  mockRepo6,
		}
		s := service.NewProductService(cfg)

		mockRepo4.On("FindSellerByUserID", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("uint")).Return(&model.Seller{}, nil)

		mockRepo3.On("GetProductsBySellerIDUnscoped", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("*dto.SellerProductSearchQuery"), mock.AnythingOfType("uint")).Return([]*model.Product{}, int64(1), nil)

		res, _, _, err := s.GetProductsByUserIDUnscoped(&dto.SellerProductSearchQuery{Limit: 1}, uint(1))

		assert.Nil(t, err)
		assert.NotNil(t, res)
	})
	t.Run("Should return err", func(t *testing.T) {

		gormDB := testutil.MockDB()
		mockRepo1 := new(mocks.ProductRepository)
		mockRepo2 := new(mocks.ReviewRepository)
		mockRepo3 := new(mocks.ProductVariantDetailRepository)
		mockRepo4 := new(mocks.SellerRepository)
		mockRepo5 := new(mocks.SocialGraphRepository)
		mockRepo6 := new(mocks.NotificationRepository)
		cfg := &service.ProductConfig{
			DB:                gormDB,
			ProductRepo:       mockRepo1,
			ReviewRepo:        mockRepo2,
			ProductVarDetRepo: mockRepo3,
			SellerRepo:        mockRepo4,
			SocialGraphRepo:   mockRepo5,
			NotificationRepo:  mockRepo6,
		}
		s := service.NewProductService(cfg)

		mockRepo4.On("FindSellerByUserID", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("uint")).Return(nil, errors.New(""))
		mockProductArray := []*model.Product{}
		mockProduct := &model.Product{}
		mockProductArray = append(mockProductArray, mockProduct)
		mockRepo3.On("GetProductsBySellerIDUnscoped", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("*dto.SellerProductSearchQuery"), mock.AnythingOfType("uint")).Return(mockProductArray, int64(1), nil)

		res, _, _, err := s.GetProductsByUserIDUnscoped(&dto.SellerProductSearchQuery{Limit: 1}, uint(1))

		assert.Nil(t, res)
		assert.NotNil(t, err)
	})
	t.Run("Should return error", func(t *testing.T) {

		gormDB := testutil.MockDB()
		mockRepo1 := new(mocks.ProductRepository)
		mockRepo2 := new(mocks.ReviewRepository)
		mockRepo3 := new(mocks.ProductVariantDetailRepository)
		mockRepo4 := new(mocks.SellerRepository)
		mockRepo5 := new(mocks.SocialGraphRepository)
		mockRepo6 := new(mocks.NotificationRepository)
		cfg := &service.ProductConfig{
			DB:                gormDB,
			ProductRepo:       mockRepo1,
			ReviewRepo:        mockRepo2,
			ProductVarDetRepo: mockRepo3,
			SellerRepo:        mockRepo4,
			SocialGraphRepo:   mockRepo5,
			NotificationRepo:  mockRepo6,
		}
		s := service.NewProductService(cfg)

		mockRepo4.On("FindSellerByUserID", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("uint")).Return(&model.Seller{}, nil)

		mockRepo3.On("GetProductsBySellerIDUnscoped", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("*dto.SellerProductSearchQuery"), mock.AnythingOfType("uint")).Return(nil, int64(1), errors.New(""))

		res, _, _, err := s.GetProductsByUserIDUnscoped(&dto.SellerProductSearchQuery{Limit: 1}, uint(1))

		assert.Nil(t, res)
		assert.NotNil(t, err)
	})
}

func TestProductService_GetProductsByCategoryID(t *testing.T) {
	t.Run("Should return response body", func(t *testing.T) {

		gormDB := testutil.MockDB()
		mockRepo1 := new(mocks.ProductRepository)
		mockRepo2 := new(mocks.ReviewRepository)
		mockRepo3 := new(mocks.ProductVariantDetailRepository)
		mockRepo4 := new(mocks.SellerRepository)
		mockRepo5 := new(mocks.SocialGraphRepository)
		mockRepo6 := new(mocks.NotificationRepository)
		cfg := &service.ProductConfig{
			DB:                gormDB,
			ProductRepo:       mockRepo1,
			ReviewRepo:        mockRepo2,
			ProductVarDetRepo: mockRepo3,
			SellerRepo:        mockRepo4,
			SocialGraphRepo:   mockRepo5,
			NotificationRepo:  mockRepo6,
		}
		s := service.NewProductService(cfg)

		mockAddress := &model.Address{City: "test"}
		mockSeller := &model.Seller{Address: mockAddress}
		mockProductPhotoArr := []*model.ProductPhoto{}
		mockProductPhoto := &model.ProductPhoto{PhotoURL: "test"}
		mockProductPhotoArr = append(mockProductPhotoArr, mockProductPhoto)
		mockProduct := &model.Product{Seller: mockSeller, ProductPhotos: mockProductPhotoArr}

		mockArray := []*dto.SellerProductsCustomTable{}
		MockSellerCustomTable := &dto.SellerProductsCustomTable{Product: *mockProduct}
		mockArray = append(mockArray, MockSellerCustomTable)
		mockRepo3.On("GetProductsByCategoryID", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("*dto.SellerProductSearchQuery"), mock.AnythingOfType("uint")).Return(mockArray, int64(1), int64(1), nil)

		res, _, _, err := s.GetProductsByCategoryID(&dto.SellerProductSearchQuery{Limit: 1}, uint(1))

		assert.Nil(t, err)
		assert.NotNil(t, res)
	})

	t.Run("Should return error", func(t *testing.T) {

		gormDB := testutil.MockDB()
		mockRepo1 := new(mocks.ProductRepository)
		mockRepo2 := new(mocks.ReviewRepository)
		mockRepo3 := new(mocks.ProductVariantDetailRepository)
		mockRepo4 := new(mocks.SellerRepository)
		mockRepo5 := new(mocks.SocialGraphRepository)
		mockRepo6 := new(mocks.NotificationRepository)
		cfg := &service.ProductConfig{
			DB:                gormDB,
			ProductRepo:       mockRepo1,
			ReviewRepo:        mockRepo2,
			ProductVarDetRepo: mockRepo3,
			SellerRepo:        mockRepo4,
			SocialGraphRepo:   mockRepo5,
			NotificationRepo:  mockRepo6,
		}
		s := service.NewProductService(cfg)

		mockRepo3.On("GetProductsByCategoryID", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("*dto.SellerProductSearchQuery"), mock.AnythingOfType("uint")).Return(nil, int64(1), int64(1), errors.New(""))

		res, _, _, err := s.GetProductsByCategoryID(&dto.SellerProductSearchQuery{Limit: 1}, uint(1))

		assert.Nil(t, res)
		assert.NotNil(t, err)
	})
}

func TestProductService_FindSimilarProducts(t *testing.T) {
	t.Run("Should return error", func(t *testing.T) {

		gormDB := testutil.MockDB()
		mockRepo1 := new(mocks.ProductRepository)
		mockRepo2 := new(mocks.ReviewRepository)
		mockRepo3 := new(mocks.ProductVariantDetailRepository)
		mockRepo4 := new(mocks.SellerRepository)
		mockRepo5 := new(mocks.SocialGraphRepository)
		mockRepo6 := new(mocks.NotificationRepository)
		cfg := &service.ProductConfig{
			DB:                gormDB,
			ProductRepo:       mockRepo1,
			ReviewRepo:        mockRepo2,
			ProductVarDetRepo: mockRepo3,
			SellerRepo:        mockRepo4,
			SocialGraphRepo:   mockRepo5,
			NotificationRepo:  mockRepo6,
		}
		s := service.NewProductService(cfg)

		mockRepo1.On("FindProductByID", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("uint")).Return(&model.Product{}, nil)

		mockAddress := &model.Address{City: "test"}
		mockSeller := &model.Seller{Address: mockAddress}
		mockProduct := &model.Product{ID: 1, Seller: mockSeller}

		mockArray := []*dto.SellerProductsCustomTable{}
		MockSellerCustomTable := &dto.SellerProductsCustomTable{Product: *mockProduct}
		mockArray = append(mockArray, MockSellerCustomTable)
		mockRepo1.On("FindSimilarProduct", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("uint"), mock.AnythingOfType("*repository.SearchQuery")).Return(mockArray, int64(1), int64(1), nil)

		res, _, _, err := s.FindSimilarProducts(uint(1), &repository.SearchQuery{})

		assert.Nil(t, err)
		assert.NotNil(t, res)
	})

	t.Run("Should return error", func(t *testing.T) {

		gormDB := testutil.MockDB()
		mockRepo1 := new(mocks.ProductRepository)
		mockRepo2 := new(mocks.ReviewRepository)
		mockRepo3 := new(mocks.ProductVariantDetailRepository)
		mockRepo4 := new(mocks.SellerRepository)
		mockRepo5 := new(mocks.SocialGraphRepository)
		mockRepo6 := new(mocks.NotificationRepository)
		cfg := &service.ProductConfig{
			DB:                gormDB,
			ProductRepo:       mockRepo1,
			ReviewRepo:        mockRepo2,
			ProductVarDetRepo: mockRepo3,
			SellerRepo:        mockRepo4,
			SocialGraphRepo:   mockRepo5,
			NotificationRepo:  mockRepo6,
		}
		s := service.NewProductService(cfg)

		mockRepo1.On("FindProductByID", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("uint")).Return(nil, errors.New(""))

		mockAddress := &model.Address{City: "test"}
		mockSeller := &model.Seller{Address: mockAddress}
		mockProduct := &model.Product{Seller: mockSeller}

		mockArray := []*dto.SellerProductsCustomTable{}
		MockSellerCustomTable := &dto.SellerProductsCustomTable{Product: *mockProduct}
		mockArray = append(mockArray, MockSellerCustomTable)
		mockRepo1.On("FindSimilarProduct", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("uint"), mock.AnythingOfType("*repository.SearchQuery")).Return(mockArray, int64(1), int64(1), nil)

		res, _, _, err := s.FindSimilarProducts(uint(1), &repository.SearchQuery{})

		assert.Nil(t, res)
		assert.NotNil(t, err)
	})

	t.Run("Should return error", func(t *testing.T) {

		gormDB := testutil.MockDB()
		mockRepo1 := new(mocks.ProductRepository)
		mockRepo2 := new(mocks.ReviewRepository)
		mockRepo3 := new(mocks.ProductVariantDetailRepository)
		mockRepo4 := new(mocks.SellerRepository)
		mockRepo5 := new(mocks.SocialGraphRepository)
		mockRepo6 := new(mocks.NotificationRepository)
		cfg := &service.ProductConfig{
			DB:                gormDB,
			ProductRepo:       mockRepo1,
			ReviewRepo:        mockRepo2,
			ProductVarDetRepo: mockRepo3,
			SellerRepo:        mockRepo4,
			SocialGraphRepo:   mockRepo5,
			NotificationRepo:  mockRepo6,
		}
		s := service.NewProductService(cfg)

		mockRepo1.On("FindProductByID", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("uint")).Return(&model.Product{}, nil)

		mockAddress := &model.Address{City: "test"}
		mockSeller := &model.Seller{Address: mockAddress}
		mockProduct := &model.Product{Seller: mockSeller}

		mockArray := []*dto.SellerProductsCustomTable{}
		MockSellerCustomTable := &dto.SellerProductsCustomTable{Product: *mockProduct}
		mockArray = append(mockArray, MockSellerCustomTable)
		mockRepo1.On("FindSimilarProduct", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("uint"), mock.AnythingOfType("*repository.SearchQuery")).Return(nil, int64(1), int64(1), errors.New(""))

		res, _, _, err := s.FindSimilarProducts(uint(1), &repository.SearchQuery{})

		assert.Nil(t, res)
		assert.NotNil(t, err)
	})
}

func TestProductService_GetProducts(t *testing.T) {
	t.Run("Should return response body", func(t *testing.T) {

		gormDB := testutil.MockDB()
		mockRepo1 := new(mocks.ProductRepository)
		mockRepo2 := new(mocks.ReviewRepository)
		mockRepo3 := new(mocks.ProductVariantDetailRepository)
		mockRepo4 := new(mocks.SellerRepository)
		mockRepo5 := new(mocks.SocialGraphRepository)
		mockRepo6 := new(mocks.NotificationRepository)
		cfg := &service.ProductConfig{
			DB:                gormDB,
			ProductRepo:       mockRepo1,
			ReviewRepo:        mockRepo2,
			ProductVarDetRepo: mockRepo3,
			SellerRepo:        mockRepo4,
			SocialGraphRepo:   mockRepo5,
			NotificationRepo:  mockRepo6,
		}
		s := service.NewProductService(cfg)
		mockProdPhArr := []*model.ProductPhoto{}
		mockProdPhoto := &model.ProductPhoto{PhotoURL: ""}
		mockProdPhArr = append(mockProdPhArr, mockProdPhoto)
		mockAddress := &model.Address{City: "test"}
		mockSeller := &model.Seller{Address: mockAddress}
		mockProduct := &model.Product{Seller: mockSeller, ProductPhotos: mockProdPhArr}

		mockArray := []*dto.SellerProductsCustomTable{}
		MockSellerCustomTable := &dto.SellerProductsCustomTable{Product: *mockProduct}
		mockArray = append(mockArray, MockSellerCustomTable)

		mockRepo3.On("SearchProducts", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("*repository.SearchQuery")).Return(mockArray, int64(1), int64(1), nil)

		res, _, _, err := s.GetProducts(&repository.SearchQuery{})

		assert.Nil(t, err)
		assert.NotNil(t, res)
	})

	t.Run("Should return error", func(t *testing.T) {

		gormDB := testutil.MockDB()
		mockRepo1 := new(mocks.ProductRepository)
		mockRepo2 := new(mocks.ReviewRepository)
		mockRepo3 := new(mocks.ProductVariantDetailRepository)
		mockRepo4 := new(mocks.SellerRepository)
		mockRepo5 := new(mocks.SocialGraphRepository)
		mockRepo6 := new(mocks.NotificationRepository)
		cfg := &service.ProductConfig{
			DB:                gormDB,
			ProductRepo:       mockRepo1,
			ReviewRepo:        mockRepo2,
			ProductVarDetRepo: mockRepo3,
			SellerRepo:        mockRepo4,
			SocialGraphRepo:   mockRepo5,
			NotificationRepo:  mockRepo6,
		}
		s := service.NewProductService(cfg)

		mockAddress := &model.Address{City: "test"}
		mockSeller := &model.Seller{Address: mockAddress}
		mockProduct := &model.Product{Seller: mockSeller}

		mockArray := []*dto.SellerProductsCustomTable{}
		MockSellerCustomTable := &dto.SellerProductsCustomTable{Product: *mockProduct}
		mockArray = append(mockArray, MockSellerCustomTable)

		mockRepo3.On("SearchProducts", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("*repository.SearchQuery")).Return(nil, int64(1), int64(1), errors.New(""))

		res, _, _, err := s.GetProducts(&repository.SearchQuery{})

		assert.Nil(t, res)
		assert.NotNil(t, err)
	})
}

func TestProductService_SearchRecommendProduct(t *testing.T) {
	t.Run("Should return response body", func(t *testing.T) {

		gormDB := testutil.MockDB()
		mockRepo1 := new(mocks.ProductRepository)
		mockRepo2 := new(mocks.ReviewRepository)
		mockRepo3 := new(mocks.ProductVariantDetailRepository)
		mockRepo4 := new(mocks.SellerRepository)
		mockRepo5 := new(mocks.SocialGraphRepository)
		mockRepo6 := new(mocks.NotificationRepository)
		cfg := &service.ProductConfig{
			DB:                gormDB,
			ProductRepo:       mockRepo1,
			ReviewRepo:        mockRepo2,
			ProductVarDetRepo: mockRepo3,
			SellerRepo:        mockRepo4,
			SocialGraphRepo:   mockRepo5,
			NotificationRepo:  mockRepo6,
		}
		s := service.NewProductService(cfg)

		mockProdPhArr := []*model.ProductPhoto{}
		mockProdPhoto := &model.ProductPhoto{PhotoURL: ""}
		mockProdPhArr = append(mockProdPhArr, mockProdPhoto)
		mockAddress := &model.Address{City: "test"}
		mockSeller := &model.Seller{Address: mockAddress}
		mockProduct := &model.Product{Seller: mockSeller, ProductPhotos: mockProdPhArr}
		mockArray := []*dto.SellerProductsCustomTable{}
		MockSellerCustomTable := &dto.SellerProductsCustomTable{Product: *mockProduct}
		mockArray = append(mockArray, MockSellerCustomTable)

		mockRepo1.On("SearchRecommendProduct", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("*repository.SearchQuery")).Return(mockArray, int64(1), int64(1), nil)

		res, _, _, err := s.SearchRecommendProduct(&repository.SearchQuery{})

		assert.Nil(t, err)
		assert.NotNil(t, res)
	})

	t.Run("Should return error", func(t *testing.T) {

		gormDB := testutil.MockDB()
		mockRepo1 := new(mocks.ProductRepository)
		mockRepo2 := new(mocks.ReviewRepository)
		mockRepo3 := new(mocks.ProductVariantDetailRepository)
		mockRepo4 := new(mocks.SellerRepository)
		mockRepo5 := new(mocks.SocialGraphRepository)
		mockRepo6 := new(mocks.NotificationRepository)
		cfg := &service.ProductConfig{
			DB:                gormDB,
			ProductRepo:       mockRepo1,
			ReviewRepo:        mockRepo2,
			ProductVarDetRepo: mockRepo3,
			SellerRepo:        mockRepo4,
			SocialGraphRepo:   mockRepo5,
			NotificationRepo:  mockRepo6,
		}
		s := service.NewProductService(cfg)

		mockProdPhArr := []*model.ProductPhoto{}
		mockProdPhoto := &model.ProductPhoto{PhotoURL: ""}
		mockProdPhArr = append(mockProdPhArr, mockProdPhoto)
		mockAddress := &model.Address{City: "test"}
		mockSeller := &model.Seller{Address: mockAddress}
		mockProduct := &model.Product{Seller: mockSeller, ProductPhotos: mockProdPhArr}
		mockArray := []*dto.SellerProductsCustomTable{}
		MockSellerCustomTable := &dto.SellerProductsCustomTable{Product: *mockProduct}
		mockArray = append(mockArray, MockSellerCustomTable)

		mockRepo1.On("SearchRecommendProduct", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("*repository.SearchQuery")).Return(nil, int64(1), int64(1), errors.New(""))

		res, _, _, err := s.SearchRecommendProduct(&repository.SearchQuery{})

		assert.Nil(t, res)
		assert.NotNil(t, err)
	})
}

func TestProductService_CreateSellerProduct(t *testing.T) {
	t.Run("Should return response body", func(t *testing.T) {

		gormDB := testutil.MockDB()
		mockRepo1 := new(mocks.ProductRepository)
		mockRepo2 := new(mocks.ReviewRepository)
		mockRepo3 := new(mocks.ProductVariantDetailRepository)
		mockRepo4 := new(mocks.SellerRepository)
		mockRepo5 := new(mocks.SocialGraphRepository)
		mockRepo6 := new(mocks.NotificationRepository)
		cfg := &service.ProductConfig{
			DB:                gormDB,
			ProductRepo:       mockRepo1,
			ReviewRepo:        mockRepo2,
			ProductVarDetRepo: mockRepo3,
			SellerRepo:        mockRepo4,
			SocialGraphRepo:   mockRepo5,
			NotificationRepo:  mockRepo6,
		}
		s := service.NewProductService(cfg)

		price := float64(10000)
		defaultStock := uint(1)

		mockPPArray := []*dto.ProductPhoto{}
		mockPP := &dto.ProductPhoto{}
		mockPPArray = append(mockPPArray, mockPP)

		mockVarArr := []*dto.VariantAndDetails{}
		mockVar := &dto.VariantAndDetails{}
		mockVarArr = append(mockVarArr, mockVar)

		mockVN1 := ""
		mockVN2 := ""
		req := &dto.PostCreateProductReq{DefaultPrice: &price, DefaultStock: &defaultStock, ProductPhotos: mockPPArray, VariantArray: mockVarArr, Variant1Name: &mockVN1, Variant2Name: &mockVN2}

		mockRepo4.On("FindSellerByUserID", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("uint")).Return(&model.Seller{}, nil)

		mockRepo1.On("CreateProduct", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("string"), mock.AnythingOfType("uint"), mock.AnythingOfType("uint"), mock.AnythingOfType("bool"), mock.AnythingOfType("uint"), mock.AnythingOfType("uint")).Return(&model.Product{}, nil)

		mockRepo1.On("CreateProductDetail", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("uint"), mock.AnythingOfType("*dto.ProductDetailsReq")).Return(&model.ProductDetail{}, nil)

		mockRepo1.On("CreateProductPhoto", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("uint"), mock.AnythingOfType("*dto.ProductPhoto")).Return(&model.ProductPhoto{}, nil)

		mockRepo1.On("CreateProductVariantDetail", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("uint"), mock.AnythingOfType("*uint"), mock.AnythingOfType("*model.ProductVariant"), mock.AnythingOfType("*dto.ProductVariantDetail")).Return(&model.ProductVariantDetail{}, nil)

		mockRepo1.On("CreateProductVariant", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("string")).Return(&model.ProductVariant{}, nil)

		mockSocialArr := []*model.SocialGraph{}
		mockSocial := &model.SocialGraph{}
		mockSocialArr = append(mockSocialArr, mockSocial)
		mockRepo5.On("GetFollowerUserID", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("uint")).Return(mockSocialArr, nil)

		mockRepo6.On("AddToNotificationFromModel", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("*model.Notification"))

		res, err := s.CreateSellerProduct(uint(1), req)

		assert.Nil(t, err)
		assert.NotNil(t, res)
	})

	t.Run("Should return error", func(t *testing.T) {

		gormDB := testutil.MockDB()
		mockRepo1 := new(mocks.ProductRepository)
		mockRepo2 := new(mocks.ReviewRepository)
		mockRepo3 := new(mocks.ProductVariantDetailRepository)
		mockRepo4 := new(mocks.SellerRepository)
		mockRepo5 := new(mocks.SocialGraphRepository)
		mockRepo6 := new(mocks.NotificationRepository)
		cfg := &service.ProductConfig{
			DB:                gormDB,
			ProductRepo:       mockRepo1,
			ReviewRepo:        mockRepo2,
			ProductVarDetRepo: mockRepo3,
			SellerRepo:        mockRepo4,
			SocialGraphRepo:   mockRepo5,
			NotificationRepo:  mockRepo6,
		}
		s := service.NewProductService(cfg)

		price := float64(10000)
		defaultStock := uint(1)

		mockPPArray := []*dto.ProductPhoto{}
		mockPP := &dto.ProductPhoto{}
		mockPPArray = append(mockPPArray, mockPP)

		mockVarArr := []*dto.VariantAndDetails{}
		mockVar := &dto.VariantAndDetails{}
		mockVarArr = append(mockVarArr, mockVar)

		mockVN1 := ""
		mockVN2 := ""
		req := &dto.PostCreateProductReq{DefaultPrice: &price, DefaultStock: &defaultStock, ProductPhotos: mockPPArray, Variant1Name: &mockVN1, Variant2Name: &mockVN2}

		mockRepo4.On("FindSellerByUserID", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("uint")).Return(&model.Seller{}, nil)

		mockRepo1.On("CreateProduct", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("string"), mock.AnythingOfType("uint"), mock.AnythingOfType("uint"), mock.AnythingOfType("bool"), mock.AnythingOfType("uint"), mock.AnythingOfType("uint")).Return(&model.Product{}, nil)

		mockRepo1.On("CreateProductDetail", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("uint"), mock.AnythingOfType("*dto.ProductDetailsReq")).Return(&model.ProductDetail{}, nil)

		mockRepo1.On("CreateProductPhoto", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("uint"), mock.AnythingOfType("*dto.ProductPhoto")).Return(&model.ProductPhoto{}, nil)

		mockRepo1.On("CreateProductVariantDetail", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("uint"), mock.AnythingOfType("*uint"), mock.AnythingOfType("*model.ProductVariant"), mock.AnythingOfType("*dto.ProductVariantDetail")).Return(&model.ProductVariantDetail{}, nil)

		mockRepo1.On("CreateProductVariant", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("string")).Return(&model.ProductVariant{}, nil)

		mockSocialArr := []*model.SocialGraph{}
		mockSocial := &model.SocialGraph{}
		mockSocialArr = append(mockSocialArr, mockSocial)
		mockRepo5.On("GetFollowerUserID", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("uint")).Return(mockSocialArr, nil)

		mockRepo6.On("AddToNotificationFromModel", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("*model.Notification"))

		res, err := s.CreateSellerProduct(uint(1), req)

		assert.Nil(t, err)
		assert.NotNil(t, res)
	})

	t.Run("Should return error", func(t *testing.T) {

		gormDB := testutil.MockDB()
		mockRepo1 := new(mocks.ProductRepository)
		mockRepo2 := new(mocks.ReviewRepository)
		mockRepo3 := new(mocks.ProductVariantDetailRepository)
		mockRepo4 := new(mocks.SellerRepository)
		mockRepo5 := new(mocks.SocialGraphRepository)
		mockRepo6 := new(mocks.NotificationRepository)
		cfg := &service.ProductConfig{
			DB:                gormDB,
			ProductRepo:       mockRepo1,
			ReviewRepo:        mockRepo2,
			ProductVarDetRepo: mockRepo3,
			SellerRepo:        mockRepo4,
			SocialGraphRepo:   mockRepo5,
			NotificationRepo:  mockRepo6,
		}
		s := service.NewProductService(cfg)

		price := float64(10000)
		defaultStock := uint(1)

		mockPPArray := []*dto.ProductPhoto{}
		mockPP := &dto.ProductPhoto{}
		mockPPArray = append(mockPPArray, mockPP)

		mockVarArr := []*dto.VariantAndDetails{}
		mockVar := &dto.VariantAndDetails{}
		mockVarArr = append(mockVarArr, mockVar)

		mockVN1 := ""
		mockVN2 := ""
		req := &dto.PostCreateProductReq{DefaultPrice: &price, DefaultStock: &defaultStock, ProductPhotos: mockPPArray, VariantArray: mockVarArr, Variant1Name: &mockVN1, Variant2Name: &mockVN2}

		mockRepo4.On("FindSellerByUserID", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("uint")).Return(nil, errors.New(""))

		mockRepo1.On("CreateProduct", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("string"), mock.AnythingOfType("uint"), mock.AnythingOfType("uint"), mock.AnythingOfType("bool"), mock.AnythingOfType("uint"), mock.AnythingOfType("uint")).Return(&model.Product{}, nil)

		mockRepo1.On("CreateProductDetail", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("uint"), mock.AnythingOfType("*dto.ProductDetailsReq")).Return(&model.ProductDetail{}, nil)

		mockRepo1.On("CreateProductPhoto", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("uint"), mock.AnythingOfType("*dto.ProductPhoto")).Return(&model.ProductPhoto{}, nil)

		mockRepo1.On("CreateProductVariantDetail", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("uint"), mock.AnythingOfType("*uint"), mock.AnythingOfType("*model.ProductVariant"), mock.AnythingOfType("*dto.ProductVariantDetail")).Return(&model.ProductVariantDetail{}, nil)

		mockRepo1.On("CreateProductVariant", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("string")).Return(&model.ProductVariant{}, nil)

		mockSocialArr := []*model.SocialGraph{}
		mockSocial := &model.SocialGraph{}
		mockSocialArr = append(mockSocialArr, mockSocial)
		mockRepo5.On("GetFollowerUserID", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("uint")).Return(mockSocialArr, nil)

		mockRepo6.On("AddToNotificationFromModel", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("*model.Notification"))

		res, err := s.CreateSellerProduct(uint(1), req)

		assert.Nil(t, res)
		assert.NotNil(t, err)
	})

	t.Run("Should return error", func(t *testing.T) {

		gormDB := testutil.MockDB()
		mockRepo1 := new(mocks.ProductRepository)
		mockRepo2 := new(mocks.ReviewRepository)
		mockRepo3 := new(mocks.ProductVariantDetailRepository)
		mockRepo4 := new(mocks.SellerRepository)
		mockRepo5 := new(mocks.SocialGraphRepository)
		mockRepo6 := new(mocks.NotificationRepository)
		cfg := &service.ProductConfig{
			DB:                gormDB,
			ProductRepo:       mockRepo1,
			ReviewRepo:        mockRepo2,
			ProductVarDetRepo: mockRepo3,
			SellerRepo:        mockRepo4,
			SocialGraphRepo:   mockRepo5,
			NotificationRepo:  mockRepo6,
		}
		s := service.NewProductService(cfg)

		price := float64(10000)
		defaultStock := uint(1)

		mockPPArray := []*dto.ProductPhoto{}
		mockPP := &dto.ProductPhoto{}
		mockPPArray = append(mockPPArray, mockPP)

		mockVarArr := []*dto.VariantAndDetails{}
		mockVar := &dto.VariantAndDetails{}
		mockVarArr = append(mockVarArr, mockVar)

		mockVN1 := ""
		mockVN2 := ""
		req := &dto.PostCreateProductReq{DefaultPrice: &price, DefaultStock: &defaultStock, ProductPhotos: mockPPArray, VariantArray: mockVarArr, Variant1Name: &mockVN1, Variant2Name: &mockVN2}

		mockRepo4.On("FindSellerByUserID", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("uint")).Return(&model.Seller{}, nil)

		mockRepo1.On("CreateProduct", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("string"), mock.AnythingOfType("uint"), mock.AnythingOfType("uint"), mock.AnythingOfType("bool"), mock.AnythingOfType("uint"), mock.AnythingOfType("uint")).Return(&model.Product{}, nil)

		mockRepo1.On("CreateProductDetail", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("uint"), mock.AnythingOfType("*dto.ProductDetailsReq")).Return(nil, errors.New(""))

		mockRepo1.On("CreateProductPhoto", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("uint"), mock.AnythingOfType("*dto.ProductPhoto")).Return(&model.ProductPhoto{}, nil)

		mockRepo1.On("CreateProductVariantDetail", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("uint"), mock.AnythingOfType("*uint"), mock.AnythingOfType("*model.ProductVariant"), mock.AnythingOfType("*dto.ProductVariantDetail")).Return(&model.ProductVariantDetail{}, nil)

		mockRepo1.On("CreateProductVariant", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("string")).Return(&model.ProductVariant{}, nil)

		mockSocialArr := []*model.SocialGraph{}
		mockSocial := &model.SocialGraph{}
		mockSocialArr = append(mockSocialArr, mockSocial)
		mockRepo5.On("GetFollowerUserID", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("uint")).Return(mockSocialArr, nil)

		mockRepo6.On("AddToNotificationFromModel", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("*model.Notification"))

		res, err := s.CreateSellerProduct(uint(1), req)

		assert.Nil(t, res)
		assert.NotNil(t, err)
	})

	t.Run("Should return error", func(t *testing.T) {

		gormDB := testutil.MockDB()
		mockRepo1 := new(mocks.ProductRepository)
		mockRepo2 := new(mocks.ReviewRepository)
		mockRepo3 := new(mocks.ProductVariantDetailRepository)
		mockRepo4 := new(mocks.SellerRepository)
		mockRepo5 := new(mocks.SocialGraphRepository)
		mockRepo6 := new(mocks.NotificationRepository)
		cfg := &service.ProductConfig{
			DB:                gormDB,
			ProductRepo:       mockRepo1,
			ReviewRepo:        mockRepo2,
			ProductVarDetRepo: mockRepo3,
			SellerRepo:        mockRepo4,
			SocialGraphRepo:   mockRepo5,
			NotificationRepo:  mockRepo6,
		}
		s := service.NewProductService(cfg)

		price := float64(10000)
		defaultStock := uint(1)

		mockPPArray := []*dto.ProductPhoto{}
		mockPP := &dto.ProductPhoto{}
		mockPPArray = append(mockPPArray, mockPP)

		mockVarArr := []*dto.VariantAndDetails{}
		mockVar := &dto.VariantAndDetails{}
		mockVarArr = append(mockVarArr, mockVar)

		mockVN1 := ""
		mockVN2 := ""
		req := &dto.PostCreateProductReq{DefaultPrice: &price, DefaultStock: &defaultStock, ProductPhotos: mockPPArray, VariantArray: mockVarArr, Variant1Name: &mockVN1, Variant2Name: &mockVN2}

		mockRepo4.On("FindSellerByUserID", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("uint")).Return(&model.Seller{}, nil)

		mockRepo1.On("CreateProduct", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("string"), mock.AnythingOfType("uint"), mock.AnythingOfType("uint"), mock.AnythingOfType("bool"), mock.AnythingOfType("uint"), mock.AnythingOfType("uint")).Return(nil, errors.New(""))

		mockRepo1.On("CreateProductDetail", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("uint"), mock.AnythingOfType("*dto.ProductDetailsReq")).Return(&model.ProductDetail{}, nil)

		mockRepo1.On("CreateProductPhoto", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("uint"), mock.AnythingOfType("*dto.ProductPhoto")).Return(&model.ProductPhoto{}, nil)

		mockRepo1.On("CreateProductVariantDetail", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("uint"), mock.AnythingOfType("*uint"), mock.AnythingOfType("*model.ProductVariant"), mock.AnythingOfType("*dto.ProductVariantDetail")).Return(&model.ProductVariantDetail{}, nil)

		mockRepo1.On("CreateProductVariant", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("string")).Return(&model.ProductVariant{}, nil)

		mockSocialArr := []*model.SocialGraph{}
		mockSocial := &model.SocialGraph{}
		mockSocialArr = append(mockSocialArr, mockSocial)
		mockRepo5.On("GetFollowerUserID", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("uint")).Return(mockSocialArr, nil)

		mockRepo6.On("AddToNotificationFromModel", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("*model.Notification"))

		res, err := s.CreateSellerProduct(uint(1), req)

		assert.Nil(t, res)
		assert.NotNil(t, err)
	})

	t.Run("Should return error", func(t *testing.T) {

		gormDB := testutil.MockDB()
		mockRepo1 := new(mocks.ProductRepository)
		mockRepo2 := new(mocks.ReviewRepository)
		mockRepo3 := new(mocks.ProductVariantDetailRepository)
		mockRepo4 := new(mocks.SellerRepository)
		mockRepo5 := new(mocks.SocialGraphRepository)
		mockRepo6 := new(mocks.NotificationRepository)
		cfg := &service.ProductConfig{
			DB:                gormDB,
			ProductRepo:       mockRepo1,
			ReviewRepo:        mockRepo2,
			ProductVarDetRepo: mockRepo3,
			SellerRepo:        mockRepo4,
			SocialGraphRepo:   mockRepo5,
			NotificationRepo:  mockRepo6,
		}
		s := service.NewProductService(cfg)

		price := float64(10000)
		defaultStock := uint(1)

		mockPPArray := []*dto.ProductPhoto{}
		mockPP := &dto.ProductPhoto{}
		mockPPArray = append(mockPPArray, mockPP)

		mockVarArr := []*dto.VariantAndDetails{}
		mockVar := &dto.VariantAndDetails{}
		mockVarArr = append(mockVarArr, mockVar)

		mockVN1 := ""
		mockVN2 := ""
		req := &dto.PostCreateProductReq{DefaultPrice: &price, DefaultStock: &defaultStock, ProductPhotos: mockPPArray, VariantArray: mockVarArr, Variant1Name: &mockVN1, Variant2Name: &mockVN2}

		mockRepo4.On("FindSellerByUserID", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("uint")).Return(&model.Seller{}, nil)

		mockRepo1.On("CreateProduct", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("string"), mock.AnythingOfType("uint"), mock.AnythingOfType("uint"), mock.AnythingOfType("bool"), mock.AnythingOfType("uint"), mock.AnythingOfType("uint")).Return(&model.Product{}, nil)

		mockRepo1.On("CreateProductDetail", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("uint"), mock.AnythingOfType("*dto.ProductDetailsReq")).Return(&model.ProductDetail{}, nil)

		mockRepo1.On("CreateProductPhoto", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("uint"), mock.AnythingOfType("*dto.ProductPhoto")).Return(nil, errors.New(""))

		mockRepo1.On("CreateProductVariantDetail", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("uint"), mock.AnythingOfType("*uint"), mock.AnythingOfType("*model.ProductVariant"), mock.AnythingOfType("*dto.ProductVariantDetail")).Return(&model.ProductVariantDetail{}, nil)

		mockRepo1.On("CreateProductVariant", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("string")).Return(&model.ProductVariant{}, nil)

		mockSocialArr := []*model.SocialGraph{}
		mockSocial := &model.SocialGraph{}
		mockSocialArr = append(mockSocialArr, mockSocial)
		mockRepo5.On("GetFollowerUserID", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("uint")).Return(mockSocialArr, nil)

		mockRepo6.On("AddToNotificationFromModel", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("*model.Notification"))

		res, err := s.CreateSellerProduct(uint(1), req)

		assert.Nil(t, res)
		assert.NotNil(t, err)
	})

	t.Run("Should return error", func(t *testing.T) {

		gormDB := testutil.MockDB()
		mockRepo1 := new(mocks.ProductRepository)
		mockRepo2 := new(mocks.ReviewRepository)
		mockRepo3 := new(mocks.ProductVariantDetailRepository)
		mockRepo4 := new(mocks.SellerRepository)
		mockRepo5 := new(mocks.SocialGraphRepository)
		mockRepo6 := new(mocks.NotificationRepository)
		cfg := &service.ProductConfig{
			DB:                gormDB,
			ProductRepo:       mockRepo1,
			ReviewRepo:        mockRepo2,
			ProductVarDetRepo: mockRepo3,
			SellerRepo:        mockRepo4,
			SocialGraphRepo:   mockRepo5,
			NotificationRepo:  mockRepo6,
		}
		s := service.NewProductService(cfg)

		price := float64(10000)
		defaultStock := uint(1)

		mockPPArray := []*dto.ProductPhoto{}
		mockPP := &dto.ProductPhoto{}
		mockPPArray = append(mockPPArray, mockPP)

		mockVarArr := []*dto.VariantAndDetails{}
		mockVar := &dto.VariantAndDetails{}
		mockVarArr = append(mockVarArr, mockVar)

		mockVN1 := ""
		req := &dto.PostCreateProductReq{DefaultPrice: &price, DefaultStock: &defaultStock, ProductPhotos: mockPPArray, VariantArray: mockVarArr, Variant1Name: &mockVN1}

		mockRepo4.On("FindSellerByUserID", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("uint")).Return(&model.Seller{}, nil)

		mockRepo1.On("CreateProduct", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("string"), mock.AnythingOfType("uint"), mock.AnythingOfType("uint"), mock.AnythingOfType("bool"), mock.AnythingOfType("uint"), mock.AnythingOfType("uint")).Return(&model.Product{}, nil)

		mockRepo1.On("CreateProductDetail", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("uint"), mock.AnythingOfType("*dto.ProductDetailsReq")).Return(&model.ProductDetail{}, nil)

		mockRepo1.On("CreateProductPhoto", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("uint"), mock.AnythingOfType("*dto.ProductPhoto")).Return(&model.ProductPhoto{}, nil)

		mockRepo1.On("CreateProductVariantDetail", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("uint"), mock.AnythingOfType("*uint"), mock.AnythingOfType("*model.ProductVariant"), mock.AnythingOfType("*dto.ProductVariantDetail")).Return(nil, errors.New(""))

		mockRepo1.On("CreateProductVariant", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("string")).Return(nil, errors.New(""))

		mockSocialArr := []*model.SocialGraph{}
		mockSocial := &model.SocialGraph{}
		mockSocialArr = append(mockSocialArr, mockSocial)
		mockRepo5.On("GetFollowerUserID", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("uint")).Return(mockSocialArr, nil)

		mockRepo6.On("AddToNotificationFromModel", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("*model.Notification"))

		res, err := s.CreateSellerProduct(uint(1), req)

		assert.Nil(t, res)
		assert.NotNil(t, err)
	})
	t.Run("Should return error", func(t *testing.T) {

		gormDB := testutil.MockDB()
		mockRepo1 := new(mocks.ProductRepository)
		mockRepo2 := new(mocks.ReviewRepository)
		mockRepo3 := new(mocks.ProductVariantDetailRepository)
		mockRepo4 := new(mocks.SellerRepository)
		mockRepo5 := new(mocks.SocialGraphRepository)
		mockRepo6 := new(mocks.NotificationRepository)
		cfg := &service.ProductConfig{
			DB:                gormDB,
			ProductRepo:       mockRepo1,
			ReviewRepo:        mockRepo2,
			ProductVarDetRepo: mockRepo3,
			SellerRepo:        mockRepo4,
			SocialGraphRepo:   mockRepo5,
			NotificationRepo:  mockRepo6,
		}
		s := service.NewProductService(cfg)

		price := float64(10000)
		defaultStock := uint(1)

		mockPPArray := []*dto.ProductPhoto{}
		mockPP := &dto.ProductPhoto{}
		mockPPArray = append(mockPPArray, mockPP)

		mockVarArr := []*dto.VariantAndDetails{}
		mockVar := &dto.VariantAndDetails{}
		mockVarArr = append(mockVarArr, mockVar)

		mockVN1 := ""
		req := &dto.PostCreateProductReq{DefaultPrice: &price, DefaultStock: &defaultStock, ProductPhotos: mockPPArray, VariantArray: mockVarArr, Variant1Name: &mockVN1}

		mockRepo4.On("FindSellerByUserID", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("uint")).Return(&model.Seller{}, nil)

		mockRepo1.On("CreateProduct", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("string"), mock.AnythingOfType("uint"), mock.AnythingOfType("uint"), mock.AnythingOfType("bool"), mock.AnythingOfType("uint"), mock.AnythingOfType("uint")).Return(&model.Product{}, nil)

		mockRepo1.On("CreateProductDetail", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("uint"), mock.AnythingOfType("*dto.ProductDetailsReq")).Return(&model.ProductDetail{}, nil)

		mockRepo1.On("CreateProductPhoto", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("uint"), mock.AnythingOfType("*dto.ProductPhoto")).Return(&model.ProductPhoto{}, nil)

		mockRepo1.On("CreateProductVariantDetail", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("uint"), mock.AnythingOfType("*uint"), mock.AnythingOfType("*model.ProductVariant"), mock.AnythingOfType("*dto.ProductVariantDetail")).Return(nil, errors.New(""))

		mockRepo1.On("CreateProductVariant", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("string")).Return(nil, errors.New(""))

		mockSocialArr := []*model.SocialGraph{}
		mockSocial := &model.SocialGraph{}
		mockSocialArr = append(mockSocialArr, mockSocial)
		mockRepo5.On("GetFollowerUserID", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("uint")).Return(mockSocialArr, nil)

		mockRepo6.On("AddToNotificationFromModel", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("*model.Notification"))

		res, err := s.CreateSellerProduct(uint(1), req)

		assert.Nil(t, res)
		assert.NotNil(t, err)
	})

	t.Run("Should return error", func(t *testing.T) {

		gormDB := testutil.MockDB()
		mockRepo1 := new(mocks.ProductRepository)
		mockRepo2 := new(mocks.ReviewRepository)
		mockRepo3 := new(mocks.ProductVariantDetailRepository)
		mockRepo4 := new(mocks.SellerRepository)
		mockRepo5 := new(mocks.SocialGraphRepository)
		mockRepo6 := new(mocks.NotificationRepository)
		cfg := &service.ProductConfig{
			DB:                gormDB,
			ProductRepo:       mockRepo1,
			ReviewRepo:        mockRepo2,
			ProductVarDetRepo: mockRepo3,
			SellerRepo:        mockRepo4,
			SocialGraphRepo:   mockRepo5,
			NotificationRepo:  mockRepo6,
		}
		s := service.NewProductService(cfg)

		price := float64(10000)
		defaultStock := uint(1)

		mockPPArray := []*dto.ProductPhoto{}
		mockPP := &dto.ProductPhoto{}
		mockPPArray = append(mockPPArray, mockPP)

		mockVarArr := []*dto.VariantAndDetails{}
		mockVar := &dto.VariantAndDetails{}
		mockVarArr = append(mockVarArr, mockVar)

		mockVN1 := ""
		mockVN2 := ""
		req := &dto.PostCreateProductReq{DefaultPrice: &price, DefaultStock: &defaultStock, ProductPhotos: mockPPArray, VariantArray: mockVarArr, Variant1Name: &mockVN1, Variant2Name: &mockVN2}

		mockRepo4.On("FindSellerByUserID", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("uint")).Return(&model.Seller{}, nil)

		mockRepo1.On("CreateProduct", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("string"), mock.AnythingOfType("uint"), mock.AnythingOfType("uint"), mock.AnythingOfType("bool"), mock.AnythingOfType("uint"), mock.AnythingOfType("uint")).Return(&model.Product{}, nil)

		mockRepo1.On("CreateProductDetail", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("uint"), mock.AnythingOfType("*dto.ProductDetailsReq")).Return(&model.ProductDetail{}, nil)

		mockRepo1.On("CreateProductPhoto", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("uint"), mock.AnythingOfType("*dto.ProductPhoto")).Return(&model.ProductPhoto{}, nil)

		mockRepo1.On("CreateProductVariantDetail", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("uint"), mock.AnythingOfType("*uint"), mock.AnythingOfType("*model.ProductVariant"), mock.AnythingOfType("*dto.ProductVariantDetail")).Return(&model.ProductVariantDetail{}, nil)

		mockRepo1.On("CreateProductVariant", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("string")).Return(nil, errors.New(""))

		mockSocialArr := []*model.SocialGraph{}
		mockSocial := &model.SocialGraph{}
		mockSocialArr = append(mockSocialArr, mockSocial)
		mockRepo5.On("GetFollowerUserID", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("uint")).Return(mockSocialArr, nil)

		mockRepo6.On("AddToNotificationFromModel", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("*model.Notification"))

		res, err := s.CreateSellerProduct(uint(1), req)

		assert.Nil(t, res)
		assert.NotNil(t, err)
	})
}

func TestProductService_UpdateProductAndDetails(t *testing.T) {
	t.Run("Should return response body", func(t *testing.T) {

		gormDB := testutil.MockDB()
		mockRepo1 := new(mocks.ProductRepository)
		mockRepo2 := new(mocks.ReviewRepository)
		mockRepo3 := new(mocks.ProductVariantDetailRepository)
		mockRepo4 := new(mocks.SellerRepository)
		mockRepo5 := new(mocks.SocialGraphRepository)
		mockRepo6 := new(mocks.NotificationRepository)
		cfg := &service.ProductConfig{
			DB:                gormDB,
			ProductRepo:       mockRepo1,
			ReviewRepo:        mockRepo2,
			ProductVarDetRepo: mockRepo3,
			SellerRepo:        mockRepo4,
			SocialGraphRepo:   mockRepo5,
			NotificationRepo:  mockRepo6,
		}
		s := service.NewProductService(cfg)

		mockProduct := &dto.PatchProduct{Name: "", IsBulkEnabled: false, MinQuantity: 1, MaxQuantity: 100}
		mockProductDetail := &dto.PatchProductDetailReq{}

		req := &dto.PatchProductAndDetailsReq{Product: mockProduct, ProductDetail: mockProductDetail}

		mockRepo4.On("FindSellerByUserID", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("uint")).Return(&model.Seller{}, nil)

		mockRepo1.On("FindProductByID", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("uint")).Return(&model.Product{}, nil)

		mockRepo1.On("UpdateProduct", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("uint"), mock.AnythingOfType("*model.Product")).Return(&model.Product{}, nil)

		mockRepo1.On("UpdateProductDetail", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("uint"), mock.AnythingOfType("*model.ProductDetail")).Return(&model.ProductDetail{}, nil)
		res, err := s.UpdateProductAndDetails(uint(1), uint(1), req)

		assert.Nil(t, err)
		assert.NotNil(t, res)
	})

	t.Run("Should return error", func(t *testing.T) {

		gormDB := testutil.MockDB()
		mockRepo1 := new(mocks.ProductRepository)
		mockRepo2 := new(mocks.ReviewRepository)
		mockRepo3 := new(mocks.ProductVariantDetailRepository)
		mockRepo4 := new(mocks.SellerRepository)
		mockRepo5 := new(mocks.SocialGraphRepository)
		mockRepo6 := new(mocks.NotificationRepository)
		cfg := &service.ProductConfig{
			DB:                gormDB,
			ProductRepo:       mockRepo1,
			ReviewRepo:        mockRepo2,
			ProductVarDetRepo: mockRepo3,
			SellerRepo:        mockRepo4,
			SocialGraphRepo:   mockRepo5,
			NotificationRepo:  mockRepo6,
		}
		s := service.NewProductService(cfg)

		mockProduct := &dto.PatchProduct{Name: "", IsBulkEnabled: false, MinQuantity: 1, MaxQuantity: 100}
		mockProductDetail := &dto.PatchProductDetailReq{}

		req := &dto.PatchProductAndDetailsReq{Product: mockProduct, ProductDetail: mockProductDetail}

		mockRepo4.On("FindSellerByUserID", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("uint")).Return(nil, errors.New(""))

		mockRepo1.On("FindProductByID", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("uint")).Return(&model.Product{}, nil)

		mockRepo1.On("UpdateProduct", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("uint"), mock.AnythingOfType("*model.Product")).Return(&model.Product{}, nil)

		mockRepo1.On("UpdateProductDetail", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("uint"), mock.AnythingOfType("*model.ProductDetail")).Return(&model.ProductDetail{}, nil)
		res, err := s.UpdateProductAndDetails(uint(1), uint(1), req)

		assert.Nil(t, res)
		assert.NotNil(t, err)
	})

	t.Run("Should return error", func(t *testing.T) {

		gormDB := testutil.MockDB()
		mockRepo1 := new(mocks.ProductRepository)
		mockRepo2 := new(mocks.ReviewRepository)
		mockRepo3 := new(mocks.ProductVariantDetailRepository)
		mockRepo4 := new(mocks.SellerRepository)
		mockRepo5 := new(mocks.SocialGraphRepository)
		mockRepo6 := new(mocks.NotificationRepository)
		cfg := &service.ProductConfig{
			DB:                gormDB,
			ProductRepo:       mockRepo1,
			ReviewRepo:        mockRepo2,
			ProductVarDetRepo: mockRepo3,
			SellerRepo:        mockRepo4,
			SocialGraphRepo:   mockRepo5,
			NotificationRepo:  mockRepo6,
		}
		s := service.NewProductService(cfg)

		mockProduct := &dto.PatchProduct{Name: "", IsBulkEnabled: false, MinQuantity: 1, MaxQuantity: 100}
		mockProductDetail := &dto.PatchProductDetailReq{}

		req := &dto.PatchProductAndDetailsReq{Product: mockProduct, ProductDetail: mockProductDetail}

		mockRepo4.On("FindSellerByUserID", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("uint")).Return(&model.Seller{ID: 1}, nil)

		mockRepo1.On("FindProductByID", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("uint")).Return(&model.Product{SellerID: 2}, nil)

		mockRepo1.On("UpdateProduct", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("uint"), mock.AnythingOfType("*model.Product")).Return(&model.Product{}, nil)

		mockRepo1.On("UpdateProductDetail", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("uint"), mock.AnythingOfType("*model.ProductDetail")).Return(&model.ProductDetail{}, nil)
		res, err := s.UpdateProductAndDetails(uint(1), uint(1), req)

		assert.Nil(t, res)
		assert.NotNil(t, err)
	})
}

func TestProductService_UpdateVariantAndDetails(t *testing.T) {
	t.Run("Should return response body", func(t *testing.T) {

		gormDB := testutil.MockDB()
		mockRepo1 := new(mocks.ProductRepository)
		mockRepo2 := new(mocks.ReviewRepository)
		mockRepo3 := new(mocks.ProductVariantDetailRepository)
		mockRepo4 := new(mocks.SellerRepository)
		mockRepo5 := new(mocks.SocialGraphRepository)
		mockRepo6 := new(mocks.NotificationRepository)
		cfg := &service.ProductConfig{
			DB:                gormDB,
			ProductRepo:       mockRepo1,
			ReviewRepo:        mockRepo2,
			ProductVarDetRepo: mockRepo3,
			SellerRepo:        mockRepo4,
			SocialGraphRepo:   mockRepo5,
			NotificationRepo:  mockRepo6,
		}
		s := service.NewProductService(cfg)
		mockVN1 := ""
		mockVN2 := ""

		mockPVD := &dto.PatchProductVariantDetail{}
		req := &dto.PatchVariantAndDetails{Variant1Name: &mockVN1, Variant2Name: &mockVN2, ProductVariantDetails: mockPVD}

		mockRepo4.On("FindSellerByUserID", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("uint")).Return(&model.Seller{}, nil)

		vID := uint(1)
		mockRepo1.On("FindProductVariantDetailsByID", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("uint")).Return(&model.ProductVariantDetail{ProductID: 1, Variant1ID: &vID}, nil)

		mockRepo1.On("FindProductByID", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("uint")).Return(&model.Product{}, nil)

		mockRepo1.On("UpdateProductVariantDetailByID", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("uint"), mock.AnythingOfType("*model.ProductVariantDetail")).Return(&model.ProductVariantDetail{}, nil)

		mockRepo1.On("UpdateProductVariantByID", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("uint"), mock.AnythingOfType("*model.ProductVariant")).Return(&model.ProductVariant{}, nil)

		mockRepo5.On("GetFavoriteUserID", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("uint")).Return([]*model.Favorite{}, nil)

		mockRepo6.On("AddToNotificationFromModel", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("*model.Notification"))

		res, err := s.UpdateVariantAndDetails(uint(1), uint(1), req)

		assert.Nil(t, err)
		assert.NotNil(t, res)
	})

	t.Run("Should return response body", func(t *testing.T) {

		gormDB := testutil.MockDB()
		mockRepo1 := new(mocks.ProductRepository)
		mockRepo2 := new(mocks.ReviewRepository)
		mockRepo3 := new(mocks.ProductVariantDetailRepository)
		mockRepo4 := new(mocks.SellerRepository)
		mockRepo5 := new(mocks.SocialGraphRepository)
		mockRepo6 := new(mocks.NotificationRepository)
		cfg := &service.ProductConfig{
			DB:                gormDB,
			ProductRepo:       mockRepo1,
			ReviewRepo:        mockRepo2,
			ProductVarDetRepo: mockRepo3,
			SellerRepo:        mockRepo4,
			SocialGraphRepo:   mockRepo5,
			NotificationRepo:  mockRepo6,
		}
		s := service.NewProductService(cfg)
		mockVN1 := ""
		mockVN2 := ""

		mockPVD := &dto.PatchProductVariantDetail{Stock: 1}
		req := &dto.PatchVariantAndDetails{Variant1Name: &mockVN1, Variant2Name: &mockVN2, ProductVariantDetails: mockPVD}

		mockRepo4.On("FindSellerByUserID", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("uint")).Return(&model.Seller{}, nil)

		vID := uint(1)
		mockRepo1.On("FindProductVariantDetailsByID", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("uint")).Return(&model.ProductVariantDetail{ProductID: 1, Variant1ID: &vID}, nil)

		mockRepo1.On("FindProductByID", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("uint")).Return(&model.Product{}, nil)

		mockRepo1.On("UpdateProductVariantDetailByID", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("uint"), mock.AnythingOfType("*model.ProductVariantDetail")).Return(&model.ProductVariantDetail{}, nil)

		mockRepo1.On("UpdateProductVariantByID", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("uint"), mock.AnythingOfType("*model.ProductVariant")).Return(&model.ProductVariant{}, nil)

		mockRepo5.On("GetFavoriteUserID", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("uint")).Return([]*model.Favorite{}, nil)

		mockRepo6.On("AddToNotificationFromModel", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("*model.Notification"))

		res, err := s.UpdateVariantAndDetails(uint(1), uint(1), req)

		assert.Nil(t, err)
		assert.NotNil(t, res)
	})

	t.Run("Should return response body", func(t *testing.T) {

		gormDB := testutil.MockDB()
		mockRepo1 := new(mocks.ProductRepository)
		mockRepo2 := new(mocks.ReviewRepository)
		mockRepo3 := new(mocks.ProductVariantDetailRepository)
		mockRepo4 := new(mocks.SellerRepository)
		mockRepo5 := new(mocks.SocialGraphRepository)
		mockRepo6 := new(mocks.NotificationRepository)
		cfg := &service.ProductConfig{
			DB:                gormDB,
			ProductRepo:       mockRepo1,
			ReviewRepo:        mockRepo2,
			ProductVarDetRepo: mockRepo3,
			SellerRepo:        mockRepo4,
			SocialGraphRepo:   mockRepo5,
			NotificationRepo:  mockRepo6,
		}
		s := service.NewProductService(cfg)
		mockVN1 := ""
		mockVN2 := ""

		mockPVD := &dto.PatchProductVariantDetail{Price: 1}
		req := &dto.PatchVariantAndDetails{Variant1Name: &mockVN1, Variant2Name: &mockVN2, ProductVariantDetails: mockPVD}

		mockRepo4.On("FindSellerByUserID", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("uint")).Return(&model.Seller{}, nil)

		vID := uint(1)
		mockRepo1.On("FindProductVariantDetailsByID", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("uint")).Return(&model.ProductVariantDetail{ProductID: 1, Variant1ID: &vID, Price: 0}, nil)

		mockRepo1.On("FindProductByID", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("uint")).Return(&model.Product{}, nil)

		mockRepo1.On("UpdateProductVariantDetailByID", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("uint"), mock.AnythingOfType("*model.ProductVariantDetail")).Return(&model.ProductVariantDetail{}, nil)

		mockRepo1.On("UpdateProductVariantByID", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("uint"), mock.AnythingOfType("*model.ProductVariant")).Return(&model.ProductVariant{}, nil)

		mockRepo5.On("GetFavoriteUserID", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("uint")).Return([]*model.Favorite{}, nil)

		mockRepo6.On("AddToNotificationFromModel", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("*model.Notification"))

		res, err := s.UpdateVariantAndDetails(uint(1), uint(1), req)

		assert.Nil(t, err)
		assert.NotNil(t, res)
	})

	t.Run("Should return response body", func(t *testing.T) {

		gormDB := testutil.MockDB()
		mockRepo1 := new(mocks.ProductRepository)
		mockRepo2 := new(mocks.ReviewRepository)
		mockRepo3 := new(mocks.ProductVariantDetailRepository)
		mockRepo4 := new(mocks.SellerRepository)
		mockRepo5 := new(mocks.SocialGraphRepository)
		mockRepo6 := new(mocks.NotificationRepository)
		cfg := &service.ProductConfig{
			DB:                gormDB,
			ProductRepo:       mockRepo1,
			ReviewRepo:        mockRepo2,
			ProductVarDetRepo: mockRepo3,
			SellerRepo:        mockRepo4,
			SocialGraphRepo:   mockRepo5,
			NotificationRepo:  mockRepo6,
		}
		s := service.NewProductService(cfg)
		mockVN1 := ""
		mockVN2 := ""

		mockPVD := &dto.PatchProductVariantDetail{Price: 1, Stock: 1}
		req := &dto.PatchVariantAndDetails{Variant1Name: &mockVN1, Variant2Name: &mockVN2, ProductVariantDetails: mockPVD}

		mockRepo4.On("FindSellerByUserID", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("uint")).Return(&model.Seller{}, nil)

		vID := uint(1)
		mockRepo1.On("FindProductVariantDetailsByID", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("uint")).Return(&model.ProductVariantDetail{ProductID: 1, Variant1ID: &vID, Price: 0}, nil)

		mockRepo1.On("FindProductByID", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("uint")).Return(&model.Product{}, nil)

		mockRepo1.On("UpdateProductVariantDetailByID", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("uint"), mock.AnythingOfType("*model.ProductVariantDetail")).Return(&model.ProductVariantDetail{}, nil)

		mockRepo1.On("UpdateProductVariantByID", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("uint"), mock.AnythingOfType("*model.ProductVariant")).Return(&model.ProductVariant{}, nil)
		mockFav := &model.Favorite{}
		mockFavArr := []*model.Favorite{mockFav}

		mockRepo5.On("GetFavoriteUserID", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("uint")).Return(mockFavArr, nil)

		mockRepo6.On("AddToNotificationFromModel", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("*model.Notification"))

		res, err := s.UpdateVariantAndDetails(uint(1), uint(1), req)

		assert.Nil(t, err)
		assert.NotNil(t, res)
	})

	t.Run("Should return response body", func(t *testing.T) {

		gormDB := testutil.MockDB()
		mockRepo1 := new(mocks.ProductRepository)
		mockRepo2 := new(mocks.ReviewRepository)
		mockRepo3 := new(mocks.ProductVariantDetailRepository)
		mockRepo4 := new(mocks.SellerRepository)
		mockRepo5 := new(mocks.SocialGraphRepository)
		mockRepo6 := new(mocks.NotificationRepository)
		cfg := &service.ProductConfig{
			DB:                gormDB,
			ProductRepo:       mockRepo1,
			ReviewRepo:        mockRepo2,
			ProductVarDetRepo: mockRepo3,
			SellerRepo:        mockRepo4,
			SocialGraphRepo:   mockRepo5,
			NotificationRepo:  mockRepo6,
		}
		s := service.NewProductService(cfg)
		mockVN1 := ""
		mockVN2 := ""

		mockPVD := &dto.PatchProductVariantDetail{Price: 1, Stock: 1}
		req := &dto.PatchVariantAndDetails{Variant1Name: &mockVN1, Variant2Name: &mockVN2, ProductVariantDetails: mockPVD}

		mockRepo4.On("FindSellerByUserID", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("uint")).Return(&model.Seller{}, nil)

		vID := uint(1)
		mockRepo1.On("FindProductVariantDetailsByID", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("uint")).Return(&model.ProductVariantDetail{ProductID: 1, Variant1ID: &vID, Price: 1}, nil)

		mockRepo1.On("FindProductByID", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("uint")).Return(&model.Product{}, nil)

		mockRepo1.On("UpdateProductVariantDetailByID", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("uint"), mock.AnythingOfType("*model.ProductVariantDetail")).Return(&model.ProductVariantDetail{}, nil)

		mockRepo1.On("UpdateProductVariantByID", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("uint"), mock.AnythingOfType("*model.ProductVariant")).Return(&model.ProductVariant{}, nil)
		mockFav := &model.Favorite{}
		mockFavArr := []*model.Favorite{mockFav}

		mockRepo5.On("GetFavoriteUserID", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("uint")).Return(mockFavArr, nil)

		mockRepo6.On("AddToNotificationFromModel", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("*model.Notification"))

		res, err := s.UpdateVariantAndDetails(uint(1), uint(1), req)

		assert.Nil(t, err)
		assert.NotNil(t, res)
	})

	t.Run("Should return error", func(t *testing.T) {

		gormDB := testutil.MockDB()
		mockRepo1 := new(mocks.ProductRepository)
		mockRepo2 := new(mocks.ReviewRepository)
		mockRepo3 := new(mocks.ProductVariantDetailRepository)
		mockRepo4 := new(mocks.SellerRepository)
		mockRepo5 := new(mocks.SocialGraphRepository)
		mockRepo6 := new(mocks.NotificationRepository)
		cfg := &service.ProductConfig{
			DB:                gormDB,
			ProductRepo:       mockRepo1,
			ReviewRepo:        mockRepo2,
			ProductVarDetRepo: mockRepo3,
			SellerRepo:        mockRepo4,
			SocialGraphRepo:   mockRepo5,
			NotificationRepo:  mockRepo6,
		}
		s := service.NewProductService(cfg)
		mockVN1 := ""
		mockVN2 := ""

		mockPVD := &dto.PatchProductVariantDetail{Price: 1, Stock: 1}
		req := &dto.PatchVariantAndDetails{Variant1Name: &mockVN1, Variant2Name: &mockVN2, ProductVariantDetails: mockPVD}

		mockRepo4.On("FindSellerByUserID", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("uint")).Return(nil, errors.New(""))

		vID := uint(1)
		mockRepo1.On("FindProductVariantDetailsByID", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("uint")).Return(&model.ProductVariantDetail{ProductID: 1, Variant1ID: &vID, Price: 1}, nil)

		mockRepo1.On("FindProductByID", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("uint")).Return(&model.Product{}, nil)

		mockRepo1.On("UpdateProductVariantDetailByID", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("uint"), mock.AnythingOfType("*model.ProductVariantDetail")).Return(&model.ProductVariantDetail{}, nil)

		mockRepo1.On("UpdateProductVariantByID", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("uint"), mock.AnythingOfType("*model.ProductVariant")).Return(&model.ProductVariant{}, nil)
		mockFav := &model.Favorite{}
		mockFavArr := []*model.Favorite{mockFav}

		mockRepo5.On("GetFavoriteUserID", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("uint")).Return(mockFavArr, nil)

		mockRepo6.On("AddToNotificationFromModel", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("*model.Notification"))

		res, err := s.UpdateVariantAndDetails(uint(1), uint(1), req)

		assert.Nil(t, res)
		assert.NotNil(t, err)
	})

	t.Run("Should return error", func(t *testing.T) {

		gormDB := testutil.MockDB()
		mockRepo1 := new(mocks.ProductRepository)
		mockRepo2 := new(mocks.ReviewRepository)
		mockRepo3 := new(mocks.ProductVariantDetailRepository)
		mockRepo4 := new(mocks.SellerRepository)
		mockRepo5 := new(mocks.SocialGraphRepository)
		mockRepo6 := new(mocks.NotificationRepository)
		cfg := &service.ProductConfig{
			DB:                gormDB,
			ProductRepo:       mockRepo1,
			ReviewRepo:        mockRepo2,
			ProductVarDetRepo: mockRepo3,
			SellerRepo:        mockRepo4,
			SocialGraphRepo:   mockRepo5,
			NotificationRepo:  mockRepo6,
		}
		s := service.NewProductService(cfg)
		mockVN1 := ""
		mockVN2 := ""

		mockPVD := &dto.PatchProductVariantDetail{Price: 1, Stock: 1}
		req := &dto.PatchVariantAndDetails{Variant1Name: &mockVN1, Variant2Name: &mockVN2, ProductVariantDetails: mockPVD}

		mockRepo4.On("FindSellerByUserID", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("uint")).Return(&model.Seller{ID: 1}, nil)

		vID := uint(1)
		mockRepo1.On("FindProductVariantDetailsByID", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("uint")).Return(&model.ProductVariantDetail{ProductID: 1, Variant1ID: &vID, Price: 1}, nil)

		mockRepo1.On("FindProductByID", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("uint")).Return(&model.Product{SellerID: 2}, nil)

		mockRepo1.On("UpdateProductVariantDetailByID", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("uint"), mock.AnythingOfType("*model.ProductVariantDetail")).Return(&model.ProductVariantDetail{}, nil)

		mockRepo1.On("UpdateProductVariantByID", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("uint"), mock.AnythingOfType("*model.ProductVariant")).Return(&model.ProductVariant{}, nil)
		mockFav := &model.Favorite{}
		mockFavArr := []*model.Favorite{mockFav}

		mockRepo5.On("GetFavoriteUserID", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("uint")).Return(mockFavArr, nil)

		mockRepo6.On("AddToNotificationFromModel", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("*model.Notification"))

		res, err := s.UpdateVariantAndDetails(uint(1), uint(1), req)

		assert.Nil(t, res)
		assert.NotNil(t, err)
	})

	t.Run("Should return error", func(t *testing.T) {

		gormDB := testutil.MockDB()
		mockRepo1 := new(mocks.ProductRepository)
		mockRepo2 := new(mocks.ReviewRepository)
		mockRepo3 := new(mocks.ProductVariantDetailRepository)
		mockRepo4 := new(mocks.SellerRepository)
		mockRepo5 := new(mocks.SocialGraphRepository)
		mockRepo6 := new(mocks.NotificationRepository)
		cfg := &service.ProductConfig{
			DB:                gormDB,
			ProductRepo:       mockRepo1,
			ReviewRepo:        mockRepo2,
			ProductVarDetRepo: mockRepo3,
			SellerRepo:        mockRepo4,
			SocialGraphRepo:   mockRepo5,
			NotificationRepo:  mockRepo6,
		}
		s := service.NewProductService(cfg)
		mockVN1 := ""
		mockVN2 := ""

		mockPVD := &dto.PatchProductVariantDetail{Price: 1, Stock: 1}
		req := &dto.PatchVariantAndDetails{Variant1Name: &mockVN1, Variant2Name: &mockVN2, ProductVariantDetails: mockPVD}

		mockRepo4.On("FindSellerByUserID", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("uint")).Return(&model.Seller{}, nil)

		vID := uint(1)
		mockRepo1.On("FindProductVariantDetailsByID", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("uint")).Return(&model.ProductVariantDetail{ProductID: 1, Variant1ID: &vID, Price: 1}, nil)

		mockRepo1.On("FindProductByID", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("uint")).Return(&model.Product{}, nil)

		mockRepo1.On("UpdateProductVariantDetailByID", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("uint"), mock.AnythingOfType("*model.ProductVariantDetail")).Return(&model.ProductVariantDetail{}, nil)

		mockRepo1.On("UpdateProductVariantByID", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("uint"), mock.AnythingOfType("*model.ProductVariant")).Return(nil, errors.New(""))

		mockFav := &model.Favorite{}
		mockFavArr := []*model.Favorite{mockFav}

		mockRepo5.On("GetFavoriteUserID", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("uint")).Return(mockFavArr, nil)

		mockRepo6.On("AddToNotificationFromModel", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("*model.Notification"))

		res, err := s.UpdateVariantAndDetails(uint(1), uint(1), req)

		assert.Nil(t, res)
		assert.NotNil(t, err)
	})

	t.Run("Should return error", func(t *testing.T) {

		gormDB := testutil.MockDB()
		mockRepo1 := new(mocks.ProductRepository)
		mockRepo2 := new(mocks.ReviewRepository)
		mockRepo3 := new(mocks.ProductVariantDetailRepository)
		mockRepo4 := new(mocks.SellerRepository)
		mockRepo5 := new(mocks.SocialGraphRepository)
		mockRepo6 := new(mocks.NotificationRepository)
		cfg := &service.ProductConfig{
			DB:                gormDB,
			ProductRepo:       mockRepo1,
			ReviewRepo:        mockRepo2,
			ProductVarDetRepo: mockRepo3,
			SellerRepo:        mockRepo4,
			SocialGraphRepo:   mockRepo5,
			NotificationRepo:  mockRepo6,
		}
		s := service.NewProductService(cfg)
		mockVN1 := ""
		mockVN2 := ""

		mockPVD := &dto.PatchProductVariantDetail{Price: 1, Stock: 1}
		req := &dto.PatchVariantAndDetails{Variant1Name: &mockVN1, Variant2Name: &mockVN2, ProductVariantDetails: mockPVD}

		mockRepo4.On("FindSellerByUserID", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("uint")).Return(&model.Seller{}, nil)

		vID := uint(1)
		mockRepo1.On("FindProductVariantDetailsByID", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("uint")).Return(&model.ProductVariantDetail{ProductID: 1, Variant1ID: &vID, Price: 1}, nil)

		mockRepo1.On("FindProductByID", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("uint")).Return(&model.Product{}, nil)

		mockRepo1.On("UpdateProductVariantDetailByID", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("uint"), mock.AnythingOfType("*model.ProductVariantDetail")).Return(&model.ProductVariantDetail{}, nil)

		mockRepo1.On("UpdateProductVariantByID", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("uint"), mock.AnythingOfType("*model.ProductVariant")).Return(nil, errors.New(""))
		mockFav := &model.Favorite{}
		mockFavArr := []*model.Favorite{mockFav}

		mockRepo5.On("GetFavoriteUserID", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("uint")).Return(mockFavArr, nil)

		mockRepo6.On("AddToNotificationFromModel", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("*model.Notification"))

		res, err := s.UpdateVariantAndDetails(uint(1), uint(1), req)

		assert.Nil(t, res)
		assert.NotNil(t, err)
	})

	t.Run("Should return error", func(t *testing.T) {

		gormDB := testutil.MockDB()
		mockRepo1 := new(mocks.ProductRepository)
		mockRepo2 := new(mocks.ReviewRepository)
		mockRepo3 := new(mocks.ProductVariantDetailRepository)
		mockRepo4 := new(mocks.SellerRepository)
		mockRepo5 := new(mocks.SocialGraphRepository)
		mockRepo6 := new(mocks.NotificationRepository)
		cfg := &service.ProductConfig{
			DB:                gormDB,
			ProductRepo:       mockRepo1,
			ReviewRepo:        mockRepo2,
			ProductVarDetRepo: mockRepo3,
			SellerRepo:        mockRepo4,
			SocialGraphRepo:   mockRepo5,
			NotificationRepo:  mockRepo6,
		}
		s := service.NewProductService(cfg)
		mockVN1 := ""
		mockVN2 := ""

		mockPVD := &dto.PatchProductVariantDetail{Price: 1, Stock: 1}
		req := &dto.PatchVariantAndDetails{Variant1Name: &mockVN1, Variant2Name: &mockVN2, ProductVariantDetails: mockPVD}

		mockRepo4.On("FindSellerByUserID", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("uint")).Return(&model.Seller{}, nil)

		vID := uint(1)
		mockRepo1.On("FindProductVariantDetailsByID", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("uint")).Return(&model.ProductVariantDetail{ProductID: 1, Variant1ID: &vID, Price: 1}, nil)

		mockRepo1.On("FindProductByID", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("uint")).Return(&model.Product{}, nil)

		mockRepo1.On("UpdateProductVariantDetailByID", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("uint"), mock.AnythingOfType("*model.ProductVariantDetail")).Return(nil, errors.New(""))

		mockRepo1.On("UpdateProductVariantByID", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("uint"), mock.AnythingOfType("*model.ProductVariant")).Return(&model.ProductVariant{}, nil)
		mockFav := &model.Favorite{}
		mockFavArr := []*model.Favorite{mockFav}

		mockRepo5.On("GetFavoriteUserID", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("uint")).Return(mockFavArr, nil)

		mockRepo6.On("AddToNotificationFromModel", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("*model.Notification"))

		res, err := s.UpdateVariantAndDetails(uint(1), uint(1), req)

		assert.Nil(t, res)
		assert.NotNil(t, err)
	})
}

func TestProductService_DeleteProductVariantDetails(t *testing.T) {
	t.Run("Should return response body", func(t *testing.T) {

		gormDB := testutil.MockDB()
		mockRepo1 := new(mocks.ProductRepository)
		mockRepo2 := new(mocks.ReviewRepository)
		mockRepo3 := new(mocks.ProductVariantDetailRepository)
		mockRepo4 := new(mocks.SellerRepository)
		mockRepo5 := new(mocks.SocialGraphRepository)
		mockRepo6 := new(mocks.NotificationRepository)
		cfg := &service.ProductConfig{
			DB:                gormDB,
			ProductRepo:       mockRepo1,
			ReviewRepo:        mockRepo2,
			ProductVarDetRepo: mockRepo3,
			SellerRepo:        mockRepo4,
			SocialGraphRepo:   mockRepo5,
			NotificationRepo:  mockRepo6,
		}
		s := service.NewProductService(cfg)

		mockRepo4.On("FindSellerByUserID", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("uint")).Return(&model.Seller{}, nil)

		vID := uint(1)
		mockRepo1.On("FindProductVariantDetailsByID", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("uint")).Return(&model.ProductVariantDetail{ProductID: 1, Variant1ID: &vID, Price: 0}, nil)

		mockRepo1.On("FindProductByID", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("uint")).Return(&model.Product{}, nil)

		mockPVD := &model.ProductVariantDetail{}
		mockPVDArr := []*model.ProductVariantDetail{mockPVD}
		mockRepo1.On("FindProductVariantDetailsByProductID", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("uint")).Return(mockPVDArr, nil)

		mockRepo1.On("CreateProductVariantDetailWithModel", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("*model.ProductVariantDetail")).Return(&model.ProductVariantDetail{}, nil)

		mockRepo1.On("DeleteProductVariantDetailsByID", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("uint")).Return(nil)

		resFloat := float64(1)
		err := s.DeleteProductVariantDetails(uint(1), uint(1), &resFloat)

		assert.Nil(t, err)
	})

	t.Run("Should return err", func(t *testing.T) {

		gormDB := testutil.MockDB()
		mockRepo1 := new(mocks.ProductRepository)
		mockRepo2 := new(mocks.ReviewRepository)
		mockRepo3 := new(mocks.ProductVariantDetailRepository)
		mockRepo4 := new(mocks.SellerRepository)
		mockRepo5 := new(mocks.SocialGraphRepository)
		mockRepo6 := new(mocks.NotificationRepository)
		cfg := &service.ProductConfig{
			DB:                gormDB,
			ProductRepo:       mockRepo1,
			ReviewRepo:        mockRepo2,
			ProductVarDetRepo: mockRepo3,
			SellerRepo:        mockRepo4,
			SocialGraphRepo:   mockRepo5,
			NotificationRepo:  mockRepo6,
		}
		s := service.NewProductService(cfg)

		mockRepo4.On("FindSellerByUserID", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("uint")).Return(&model.Seller{}, nil)

		vID := uint(1)
		mockRepo1.On("FindProductVariantDetailsByID", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("uint")).Return(&model.ProductVariantDetail{ProductID: 1, Variant1ID: &vID, Price: 0}, nil)

		mockRepo1.On("FindProductByID", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("uint")).Return(&model.Product{}, nil)

		mockPVD := &model.ProductVariantDetail{}
		mockPVDArr := []*model.ProductVariantDetail{mockPVD}
		mockRepo1.On("FindProductVariantDetailsByProductID", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("uint")).Return(mockPVDArr, nil)

		mockRepo1.On("CreateProductVariantDetailWithModel", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("*model.ProductVariantDetail")).Return(&model.ProductVariantDetail{}, nil)

		mockRepo1.On("DeleteProductVariantDetailsByID", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("uint")).Return(nil)

		//resFloat := float64(1)
		err := s.DeleteProductVariantDetails(uint(1), uint(1), nil)

		assert.NotNil(t, err)
	})
	t.Run("Should return err", func(t *testing.T) {

		gormDB := testutil.MockDB()
		mockRepo1 := new(mocks.ProductRepository)
		mockRepo2 := new(mocks.ReviewRepository)
		mockRepo3 := new(mocks.ProductVariantDetailRepository)
		mockRepo4 := new(mocks.SellerRepository)
		mockRepo5 := new(mocks.SocialGraphRepository)
		mockRepo6 := new(mocks.NotificationRepository)
		cfg := &service.ProductConfig{
			DB:                gormDB,
			ProductRepo:       mockRepo1,
			ReviewRepo:        mockRepo2,
			ProductVarDetRepo: mockRepo3,
			SellerRepo:        mockRepo4,
			SocialGraphRepo:   mockRepo5,
			NotificationRepo:  mockRepo6,
		}
		s := service.NewProductService(cfg)

		mockRepo4.On("FindSellerByUserID", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("uint")).Return(nil, errors.New(""))

		vID := uint(1)
		mockRepo1.On("FindProductVariantDetailsByID", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("uint")).Return(&model.ProductVariantDetail{ProductID: 1, Variant1ID: &vID, Price: 0}, nil)

		mockRepo1.On("FindProductByID", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("uint")).Return(&model.Product{}, nil)

		mockPVD := &model.ProductVariantDetail{}
		mockPVDArr := []*model.ProductVariantDetail{mockPVD}
		mockRepo1.On("FindProductVariantDetailsByProductID", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("uint")).Return(mockPVDArr, nil)

		mockRepo1.On("CreateProductVariantDetailWithModel", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("*model.ProductVariantDetail")).Return(&model.ProductVariantDetail{}, nil)

		mockRepo1.On("DeleteProductVariantDetailsByID", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("uint")).Return(nil)

		resFloat := float64(1)
		err := s.DeleteProductVariantDetails(uint(1), uint(1), &resFloat)

		assert.NotNil(t, err)
	})
	t.Run("Should return err", func(t *testing.T) {

		gormDB := testutil.MockDB()
		mockRepo1 := new(mocks.ProductRepository)
		mockRepo2 := new(mocks.ReviewRepository)
		mockRepo3 := new(mocks.ProductVariantDetailRepository)
		mockRepo4 := new(mocks.SellerRepository)
		mockRepo5 := new(mocks.SocialGraphRepository)
		mockRepo6 := new(mocks.NotificationRepository)
		cfg := &service.ProductConfig{
			DB:                gormDB,
			ProductRepo:       mockRepo1,
			ReviewRepo:        mockRepo2,
			ProductVarDetRepo: mockRepo3,
			SellerRepo:        mockRepo4,
			SocialGraphRepo:   mockRepo5,
			NotificationRepo:  mockRepo6,
		}
		s := service.NewProductService(cfg)

		mockRepo4.On("FindSellerByUserID", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("uint")).Return(&model.Seller{ID: 1}, nil)

		vID := uint(1)
		mockRepo1.On("FindProductVariantDetailsByID", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("uint")).Return(&model.ProductVariantDetail{ProductID: 1, Variant1ID: &vID, Price: 0}, nil)

		mockRepo1.On("FindProductByID", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("uint")).Return(&model.Product{SellerID: 2}, nil)

		mockPVD := &model.ProductVariantDetail{}
		mockPVDArr := []*model.ProductVariantDetail{mockPVD}
		mockRepo1.On("FindProductVariantDetailsByProductID", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("uint")).Return(mockPVDArr, nil)

		mockRepo1.On("CreateProductVariantDetailWithModel", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("*model.ProductVariantDetail")).Return(&model.ProductVariantDetail{}, nil)

		mockRepo1.On("DeleteProductVariantDetailsByID", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("uint")).Return(nil)

		resFloat := float64(1)
		err := s.DeleteProductVariantDetails(uint(1), uint(1), &resFloat)

		assert.NotNil(t, err)
	})

	t.Run("Should return err", func(t *testing.T) {

		gormDB := testutil.MockDB()
		mockRepo1 := new(mocks.ProductRepository)
		mockRepo2 := new(mocks.ReviewRepository)
		mockRepo3 := new(mocks.ProductVariantDetailRepository)
		mockRepo4 := new(mocks.SellerRepository)
		mockRepo5 := new(mocks.SocialGraphRepository)
		mockRepo6 := new(mocks.NotificationRepository)
		cfg := &service.ProductConfig{
			DB:                gormDB,
			ProductRepo:       mockRepo1,
			ReviewRepo:        mockRepo2,
			ProductVarDetRepo: mockRepo3,
			SellerRepo:        mockRepo4,
			SocialGraphRepo:   mockRepo5,
			NotificationRepo:  mockRepo6,
		}
		s := service.NewProductService(cfg)

		mockRepo4.On("FindSellerByUserID", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("uint")).Return(&model.Seller{}, nil)

		vID := uint(1)
		mockRepo1.On("FindProductVariantDetailsByID", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("uint")).Return(&model.ProductVariantDetail{ProductID: 1, Variant1ID: &vID, Price: 0}, nil)

		mockRepo1.On("FindProductByID", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("uint")).Return(&model.Product{}, nil)

		mockPVD := &model.ProductVariantDetail{}
		mockPVDArr := []*model.ProductVariantDetail{mockPVD}
		mockRepo1.On("FindProductVariantDetailsByProductID", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("uint")).Return(mockPVDArr, nil)

		mockRepo1.On("CreateProductVariantDetailWithModel", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("*model.ProductVariantDetail")).Return(nil, errors.New(""))

		mockRepo1.On("DeleteProductVariantDetailsByID", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("uint")).Return(nil)

		resFloat := float64(1)
		err := s.DeleteProductVariantDetails(uint(1), uint(1), &resFloat)

		assert.NotNil(t, err)
	})
	t.Run("Should return err", func(t *testing.T) {

		gormDB := testutil.MockDB()
		mockRepo1 := new(mocks.ProductRepository)
		mockRepo2 := new(mocks.ReviewRepository)
		mockRepo3 := new(mocks.ProductVariantDetailRepository)
		mockRepo4 := new(mocks.SellerRepository)
		mockRepo5 := new(mocks.SocialGraphRepository)
		mockRepo6 := new(mocks.NotificationRepository)
		cfg := &service.ProductConfig{
			DB:                gormDB,
			ProductRepo:       mockRepo1,
			ReviewRepo:        mockRepo2,
			ProductVarDetRepo: mockRepo3,
			SellerRepo:        mockRepo4,
			SocialGraphRepo:   mockRepo5,
			NotificationRepo:  mockRepo6,
		}
		s := service.NewProductService(cfg)

		mockRepo4.On("FindSellerByUserID", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("uint")).Return(&model.Seller{}, nil)

		vID := uint(1)
		mockRepo1.On("FindProductVariantDetailsByID", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("uint")).Return(&model.ProductVariantDetail{ProductID: 1, Variant1ID: &vID, Price: 0}, nil)

		mockRepo1.On("FindProductByID", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("uint")).Return(&model.Product{}, nil)

		mockPVD := &model.ProductVariantDetail{}
		mockPVDArr := []*model.ProductVariantDetail{mockPVD}
		mockRepo1.On("FindProductVariantDetailsByProductID", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("uint")).Return(mockPVDArr, nil)

		mockRepo1.On("CreateProductVariantDetailWithModel", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("*model.ProductVariantDetail")).Return(&model.ProductVariantDetail{}, nil)

		mockRepo1.On("DeleteProductVariantDetailsByID", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("uint")).Return(errors.New(""))

		resFloat := float64(1)
		err := s.DeleteProductVariantDetails(uint(1), uint(1), &resFloat)

		assert.NotNil(t, err)
	})
}

func TestProductService_AddVariantDetails(t *testing.T) {
	t.Run("Should return response body", func(t *testing.T) {

		gormDB := testutil.MockDB()
		mockRepo1 := new(mocks.ProductRepository)
		mockRepo2 := new(mocks.ReviewRepository)
		mockRepo3 := new(mocks.ProductVariantDetailRepository)
		mockRepo4 := new(mocks.SellerRepository)
		mockRepo5 := new(mocks.SocialGraphRepository)
		mockRepo6 := new(mocks.NotificationRepository)
		cfg := &service.ProductConfig{
			DB:                gormDB,
			ProductRepo:       mockRepo1,
			ReviewRepo:        mockRepo2,
			ProductVarDetRepo: mockRepo3,
			SellerRepo:        mockRepo4,
			SocialGraphRepo:   mockRepo5,
			NotificationRepo:  mockRepo6,
		}
		s := service.NewProductService(cfg)

		mockVN1 := ""
		mockVN2 := ""
		mockPVD2 := &dto.ProductVariantDetail{}
		mockPVDArr := []*dto.ProductVariantDetail{mockPVD2}

		req := &dto.AddVariantAndDetails{Variant1Name: &mockVN1, Variant2Name: &mockVN2, ProductVariantDetails: mockPVDArr}

		mockRepo4.On("FindSellerByUserID", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("uint")).Return(&model.Seller{}, nil)

		mockRepo1.On("FindProductByID", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("uint")).Return(&model.Product{}, nil)

		mockRepo1.On("DeleteNullProductVariantDetailsByID", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("uint")).Return(nil)

		mockRepo1.On("GetVariantByName", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("string")).Return(&model.ProductVariant{}, nil)

		mockRepo1.On("CreateVariantWithName", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("string")).Return(&model.ProductVariant{}, nil)

		mockRepo1.On("CreateProductVariantDetailWithModel", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("*model.ProductVariantDetail")).Return(&model.ProductVariantDetail{}, nil)

		res, err := s.AddVariantDetails(uint(1), uint(1), req)

		assert.Nil(t, err)
		assert.NotNil(t, res)
	})

	t.Run("Should return error", func(t *testing.T) {

		gormDB := testutil.MockDB()
		mockRepo1 := new(mocks.ProductRepository)
		mockRepo2 := new(mocks.ReviewRepository)
		mockRepo3 := new(mocks.ProductVariantDetailRepository)
		mockRepo4 := new(mocks.SellerRepository)
		mockRepo5 := new(mocks.SocialGraphRepository)
		mockRepo6 := new(mocks.NotificationRepository)
		cfg := &service.ProductConfig{
			DB:                gormDB,
			ProductRepo:       mockRepo1,
			ReviewRepo:        mockRepo2,
			ProductVarDetRepo: mockRepo3,
			SellerRepo:        mockRepo4,
			SocialGraphRepo:   mockRepo5,
			NotificationRepo:  mockRepo6,
		}
		s := service.NewProductService(cfg)

		mockVN1 := ""
		mockVN2 := ""
		mockPVD2 := &dto.ProductVariantDetail{}
		mockPVDArr := []*dto.ProductVariantDetail{mockPVD2}

		req := &dto.AddVariantAndDetails{Variant1Name: &mockVN1, Variant2Name: &mockVN2, ProductVariantDetails: mockPVDArr}

		mockRepo4.On("FindSellerByUserID", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("uint")).Return(nil, errors.New(""))

		mockRepo1.On("FindProductByID", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("uint")).Return(&model.Product{}, nil)

		mockRepo1.On("DeleteNullProductVariantDetailsByID", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("uint")).Return(nil)

		mockRepo1.On("GetVariantByName", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("string")).Return(&model.ProductVariant{}, nil)

		mockRepo1.On("CreateVariantWithName", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("string")).Return(&model.ProductVariant{}, nil)

		mockRepo1.On("CreateProductVariantDetailWithModel", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("*model.ProductVariantDetail")).Return(&model.ProductVariantDetail{}, nil)

		res, err := s.AddVariantDetails(uint(1), uint(1), req)

		assert.Nil(t, res)
		assert.NotNil(t, err)
	})

	t.Run("Should return error", func(t *testing.T) {

		gormDB := testutil.MockDB()
		mockRepo1 := new(mocks.ProductRepository)
		mockRepo2 := new(mocks.ReviewRepository)
		mockRepo3 := new(mocks.ProductVariantDetailRepository)
		mockRepo4 := new(mocks.SellerRepository)
		mockRepo5 := new(mocks.SocialGraphRepository)
		mockRepo6 := new(mocks.NotificationRepository)
		cfg := &service.ProductConfig{
			DB:                gormDB,
			ProductRepo:       mockRepo1,
			ReviewRepo:        mockRepo2,
			ProductVarDetRepo: mockRepo3,
			SellerRepo:        mockRepo4,
			SocialGraphRepo:   mockRepo5,
			NotificationRepo:  mockRepo6,
		}
		s := service.NewProductService(cfg)

		mockVN1 := ""
		mockVN2 := ""
		mockPVD2 := &dto.ProductVariantDetail{}
		mockPVDArr := []*dto.ProductVariantDetail{mockPVD2}

		req := &dto.AddVariantAndDetails{Variant1Name: &mockVN1, Variant2Name: &mockVN2, ProductVariantDetails: mockPVDArr}

		mockRepo4.On("FindSellerByUserID", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("uint")).Return(&model.Seller{ID: 1}, nil)

		mockRepo1.On("FindProductByID", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("uint")).Return(&model.Product{SellerID: 2}, nil)

		mockRepo1.On("DeleteNullProductVariantDetailsByID", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("uint")).Return(nil)

		mockRepo1.On("GetVariantByName", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("string")).Return(&model.ProductVariant{}, nil)

		mockRepo1.On("CreateVariantWithName", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("string")).Return(&model.ProductVariant{}, nil)

		mockRepo1.On("CreateProductVariantDetailWithModel", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("*model.ProductVariantDetail")).Return(&model.ProductVariantDetail{}, nil)

		res, err := s.AddVariantDetails(uint(1), uint(1), req)

		assert.Nil(t, res)
		assert.NotNil(t, err)
	})

	t.Run("Should return error", func(t *testing.T) {

		gormDB := testutil.MockDB()
		mockRepo1 := new(mocks.ProductRepository)
		mockRepo2 := new(mocks.ReviewRepository)
		mockRepo3 := new(mocks.ProductVariantDetailRepository)
		mockRepo4 := new(mocks.SellerRepository)
		mockRepo5 := new(mocks.SocialGraphRepository)
		mockRepo6 := new(mocks.NotificationRepository)
		cfg := &service.ProductConfig{
			DB:                gormDB,
			ProductRepo:       mockRepo1,
			ReviewRepo:        mockRepo2,
			ProductVarDetRepo: mockRepo3,
			SellerRepo:        mockRepo4,
			SocialGraphRepo:   mockRepo5,
			NotificationRepo:  mockRepo6,
		}
		s := service.NewProductService(cfg)

		mockVN1 := ""
		mockVN2 := ""
		mockPVD2 := &dto.ProductVariantDetail{}
		mockPVDArr := []*dto.ProductVariantDetail{mockPVD2}

		req := &dto.AddVariantAndDetails{Variant1Name: &mockVN1, Variant2Name: &mockVN2, ProductVariantDetails: mockPVDArr}

		mockRepo4.On("FindSellerByUserID", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("uint")).Return(&model.Seller{}, nil)

		mockRepo1.On("FindProductByID", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("uint")).Return(&model.Product{}, nil)

		mockRepo1.On("DeleteNullProductVariantDetailsByID", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("uint")).Return(nil)

		mockRepo1.On("GetVariantByName", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("string")).Return(nil, errors.New(""))

		mockRepo1.On("CreateVariantWithName", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("string")).Return(nil, errors.New(""))

		mockRepo1.On("CreateProductVariantDetailWithModel", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("*model.ProductVariantDetail")).Return(&model.ProductVariantDetail{}, nil)

		res, err := s.AddVariantDetails(uint(1), uint(1), req)

		assert.Nil(t, res)
		assert.NotNil(t, err)
	})
}

func TestProductService_AddProductPhoto(t *testing.T) {
	t.Run("Should return response body", func(t *testing.T) {

		gormDB := testutil.MockDB()
		mockRepo1 := new(mocks.ProductRepository)
		mockRepo2 := new(mocks.ReviewRepository)
		mockRepo3 := new(mocks.ProductVariantDetailRepository)
		mockRepo4 := new(mocks.SellerRepository)
		mockRepo5 := new(mocks.SocialGraphRepository)
		mockRepo6 := new(mocks.NotificationRepository)
		cfg := &service.ProductConfig{
			DB:                gormDB,
			ProductRepo:       mockRepo1,
			ReviewRepo:        mockRepo2,
			ProductVarDetRepo: mockRepo3,
			SellerRepo:        mockRepo4,
			SocialGraphRepo:   mockRepo5,
			NotificationRepo:  mockRepo6,
		}
		s := service.NewProductService(cfg)

		req := &dto.ProductPhotoReq{}

		mockRepo4.On("FindSellerByUserID", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("uint")).Return(&model.Seller{}, nil)

		mockRepo1.On("FindProductByID", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("uint")).Return(&model.Product{}, nil)

		mockRepo1.On("CreateProductPhotos", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("uint"), mock.AnythingOfType("*dto.ProductPhotoReq")).Return([]*model.ProductPhoto{}, nil)

		res, err := s.AddProductPhoto(uint(1), uint(1), req)

		assert.Nil(t, err)
		assert.NotNil(t, res)
	})

	t.Run("Should return error", func(t *testing.T) {

		gormDB := testutil.MockDB()
		mockRepo1 := new(mocks.ProductRepository)
		mockRepo2 := new(mocks.ReviewRepository)
		mockRepo3 := new(mocks.ProductVariantDetailRepository)
		mockRepo4 := new(mocks.SellerRepository)
		mockRepo5 := new(mocks.SocialGraphRepository)
		mockRepo6 := new(mocks.NotificationRepository)
		cfg := &service.ProductConfig{
			DB:                gormDB,
			ProductRepo:       mockRepo1,
			ReviewRepo:        mockRepo2,
			ProductVarDetRepo: mockRepo3,
			SellerRepo:        mockRepo4,
			SocialGraphRepo:   mockRepo5,
			NotificationRepo:  mockRepo6,
		}
		s := service.NewProductService(cfg)

		req := &dto.ProductPhotoReq{}

		mockRepo4.On("FindSellerByUserID", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("uint")).Return(nil, errors.New(""))

		mockRepo1.On("FindProductByID", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("uint")).Return(&model.Product{}, nil)

		mockRepo1.On("CreateProductPhotos", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("uint"), mock.AnythingOfType("*dto.ProductPhotoReq")).Return([]*model.ProductPhoto{}, nil)

		res, err := s.AddProductPhoto(uint(1), uint(1), req)

		assert.Nil(t, res)
		assert.NotNil(t, err)
	})

	t.Run("Should return error", func(t *testing.T) {

		gormDB := testutil.MockDB()
		mockRepo1 := new(mocks.ProductRepository)
		mockRepo2 := new(mocks.ReviewRepository)
		mockRepo3 := new(mocks.ProductVariantDetailRepository)
		mockRepo4 := new(mocks.SellerRepository)
		mockRepo5 := new(mocks.SocialGraphRepository)
		mockRepo6 := new(mocks.NotificationRepository)
		cfg := &service.ProductConfig{
			DB:                gormDB,
			ProductRepo:       mockRepo1,
			ReviewRepo:        mockRepo2,
			ProductVarDetRepo: mockRepo3,
			SellerRepo:        mockRepo4,
			SocialGraphRepo:   mockRepo5,
			NotificationRepo:  mockRepo6,
		}
		s := service.NewProductService(cfg)

		req := &dto.ProductPhotoReq{}

		mockRepo4.On("FindSellerByUserID", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("uint")).Return(&model.Seller{ID: 1}, nil)

		mockRepo1.On("FindProductByID", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("uint")).Return(&model.Product{SellerID: 2}, nil)

		mockRepo1.On("CreateProductPhotos", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("uint"), mock.AnythingOfType("*dto.ProductPhotoReq")).Return([]*model.ProductPhoto{}, nil)

		res, err := s.AddProductPhoto(uint(1), uint(1), req)

		assert.Nil(t, res)
		assert.NotNil(t, err)
	})

	t.Run("Should return error", func(t *testing.T) {

		gormDB := testutil.MockDB()
		mockRepo1 := new(mocks.ProductRepository)
		mockRepo2 := new(mocks.ReviewRepository)
		mockRepo3 := new(mocks.ProductVariantDetailRepository)
		mockRepo4 := new(mocks.SellerRepository)
		mockRepo5 := new(mocks.SocialGraphRepository)
		mockRepo6 := new(mocks.NotificationRepository)
		cfg := &service.ProductConfig{
			DB:                gormDB,
			ProductRepo:       mockRepo1,
			ReviewRepo:        mockRepo2,
			ProductVarDetRepo: mockRepo3,
			SellerRepo:        mockRepo4,
			SocialGraphRepo:   mockRepo5,
			NotificationRepo:  mockRepo6,
		}
		s := service.NewProductService(cfg)

		req := &dto.ProductPhotoReq{}

		mockRepo4.On("FindSellerByUserID", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("uint")).Return(&model.Seller{}, nil)

		mockRepo1.On("FindProductByID", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("uint")).Return(&model.Product{}, nil)

		mockRepo1.On("CreateProductPhotos", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("uint"), mock.AnythingOfType("*dto.ProductPhotoReq")).Return(nil, errors.New(""))

		res, err := s.AddProductPhoto(uint(1), uint(1), req)

		assert.Nil(t, res)
		assert.NotNil(t, err)
	})
}

func TestProductService_DeleteProductPhoto(t *testing.T) {
	t.Run("Should return response body", func(t *testing.T) {

		gormDB := testutil.MockDB()
		mockRepo1 := new(mocks.ProductRepository)
		mockRepo2 := new(mocks.ReviewRepository)
		mockRepo3 := new(mocks.ProductVariantDetailRepository)
		mockRepo4 := new(mocks.SellerRepository)
		mockRepo5 := new(mocks.SocialGraphRepository)
		mockRepo6 := new(mocks.NotificationRepository)
		cfg := &service.ProductConfig{
			DB:                gormDB,
			ProductRepo:       mockRepo1,
			ReviewRepo:        mockRepo2,
			ProductVarDetRepo: mockRepo3,
			SellerRepo:        mockRepo4,
			SocialGraphRepo:   mockRepo5,
			NotificationRepo:  mockRepo6,
		}
		s := service.NewProductService(cfg)

		req := &dto.DeleteProductPhoto{}

		mockRepo4.On("FindSellerByUserID", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("uint")).Return(&model.Seller{}, nil)

		mockRepo1.On("FindProductByID", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("uint")).Return(&model.Product{}, nil)

		mockRepo1.On("DeleteProductPhotos", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("*dto.DeleteProductPhoto")).Return([]*model.ProductPhoto{}, nil)

		res, err := s.DeleteProductPhoto(uint(1), uint(1), req)

		assert.Nil(t, err)
		assert.NotNil(t, res)
	})

	t.Run("Should return error", func(t *testing.T) {

		gormDB := testutil.MockDB()
		mockRepo1 := new(mocks.ProductRepository)
		mockRepo2 := new(mocks.ReviewRepository)
		mockRepo3 := new(mocks.ProductVariantDetailRepository)
		mockRepo4 := new(mocks.SellerRepository)
		mockRepo5 := new(mocks.SocialGraphRepository)
		mockRepo6 := new(mocks.NotificationRepository)
		cfg := &service.ProductConfig{
			DB:                gormDB,
			ProductRepo:       mockRepo1,
			ReviewRepo:        mockRepo2,
			ProductVarDetRepo: mockRepo3,
			SellerRepo:        mockRepo4,
			SocialGraphRepo:   mockRepo5,
			NotificationRepo:  mockRepo6,
		}
		s := service.NewProductService(cfg)

		req := &dto.DeleteProductPhoto{}

		mockRepo4.On("FindSellerByUserID", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("uint")).Return(nil, errors.New(""))

		mockRepo1.On("FindProductByID", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("uint")).Return(&model.Product{}, nil)

		mockRepo1.On("DeleteProductPhotos", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("*dto.DeleteProductPhoto")).Return([]*model.ProductPhoto{}, nil)

		res, err := s.DeleteProductPhoto(uint(1), uint(1), req)

		assert.Nil(t, res)
		assert.NotNil(t, err)
	})

	t.Run("Should return error", func(t *testing.T) {

		gormDB := testutil.MockDB()
		mockRepo1 := new(mocks.ProductRepository)
		mockRepo2 := new(mocks.ReviewRepository)
		mockRepo3 := new(mocks.ProductVariantDetailRepository)
		mockRepo4 := new(mocks.SellerRepository)
		mockRepo5 := new(mocks.SocialGraphRepository)
		mockRepo6 := new(mocks.NotificationRepository)
		cfg := &service.ProductConfig{
			DB:                gormDB,
			ProductRepo:       mockRepo1,
			ReviewRepo:        mockRepo2,
			ProductVarDetRepo: mockRepo3,
			SellerRepo:        mockRepo4,
			SocialGraphRepo:   mockRepo5,
			NotificationRepo:  mockRepo6,
		}
		s := service.NewProductService(cfg)

		req := &dto.DeleteProductPhoto{}

		mockRepo4.On("FindSellerByUserID", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("uint")).Return(&model.Seller{ID: 1}, nil)

		mockRepo1.On("FindProductByID", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("uint")).Return(&model.Product{SellerID: 2}, nil)

		mockRepo1.On("DeleteProductPhotos", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("*dto.DeleteProductPhoto")).Return([]*model.ProductPhoto{}, nil)

		res, err := s.DeleteProductPhoto(uint(1), uint(1), req)

		assert.Nil(t, res)
		assert.NotNil(t, err)
	})

	t.Run("Should return error", func(t *testing.T) {

		gormDB := testutil.MockDB()
		mockRepo1 := new(mocks.ProductRepository)
		mockRepo2 := new(mocks.ReviewRepository)
		mockRepo3 := new(mocks.ProductVariantDetailRepository)
		mockRepo4 := new(mocks.SellerRepository)
		mockRepo5 := new(mocks.SocialGraphRepository)
		mockRepo6 := new(mocks.NotificationRepository)
		cfg := &service.ProductConfig{
			DB:                gormDB,
			ProductRepo:       mockRepo1,
			ReviewRepo:        mockRepo2,
			ProductVarDetRepo: mockRepo3,
			SellerRepo:        mockRepo4,
			SocialGraphRepo:   mockRepo5,
			NotificationRepo:  mockRepo6,
		}
		s := service.NewProductService(cfg)

		req := &dto.DeleteProductPhoto{}

		mockRepo4.On("FindSellerByUserID", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("uint")).Return(&model.Seller{}, nil)

		mockRepo1.On("FindProductByID", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("uint")).Return(&model.Product{}, nil)

		mockRepo1.On("DeleteProductPhotos", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("*dto.DeleteProductPhoto")).Return(nil, errors.New(""))

		res, err := s.DeleteProductPhoto(uint(1), uint(1), req)

		assert.Nil(t, res)
		assert.NotNil(t, err)
	})
}

func TestProductService_DeleteProduct(t *testing.T) {
	t.Run("Should return response body", func(t *testing.T) {

		gormDB := testutil.MockDB()
		mockRepo1 := new(mocks.ProductRepository)
		mockRepo2 := new(mocks.ReviewRepository)
		mockRepo3 := new(mocks.ProductVariantDetailRepository)
		mockRepo4 := new(mocks.SellerRepository)
		mockRepo5 := new(mocks.SocialGraphRepository)
		mockRepo6 := new(mocks.NotificationRepository)
		cfg := &service.ProductConfig{
			DB:                gormDB,
			ProductRepo:       mockRepo1,
			ReviewRepo:        mockRepo2,
			ProductVarDetRepo: mockRepo3,
			SellerRepo:        mockRepo4,
			SocialGraphRepo:   mockRepo5,
			NotificationRepo:  mockRepo6,
		}
		s := service.NewProductService(cfg)

		mockRepo4.On("FindSellerByUserID", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("uint")).Return(&model.Seller{}, nil)

		mockRepo1.On("FindProductByID", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("uint")).Return(&model.Product{}, nil)

		mockRepo1.On("DeleteProduct", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("uint")).Return(&model.Product{}, nil)

		res, err := s.DeleteProduct(uint(1), uint(1))

		assert.Nil(t, err)
		assert.NotNil(t, res)
	})

	t.Run("Should return error", func(t *testing.T) {

		gormDB := testutil.MockDB()
		mockRepo1 := new(mocks.ProductRepository)
		mockRepo2 := new(mocks.ReviewRepository)
		mockRepo3 := new(mocks.ProductVariantDetailRepository)
		mockRepo4 := new(mocks.SellerRepository)
		mockRepo5 := new(mocks.SocialGraphRepository)
		mockRepo6 := new(mocks.NotificationRepository)
		cfg := &service.ProductConfig{
			DB:                gormDB,
			ProductRepo:       mockRepo1,
			ReviewRepo:        mockRepo2,
			ProductVarDetRepo: mockRepo3,
			SellerRepo:        mockRepo4,
			SocialGraphRepo:   mockRepo5,
			NotificationRepo:  mockRepo6,
		}
		s := service.NewProductService(cfg)

		mockRepo4.On("FindSellerByUserID", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("uint")).Return(nil, errors.New(""))

		mockRepo1.On("FindProductByID", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("uint")).Return(&model.Product{}, nil)

		mockRepo1.On("DeleteProduct", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("uint")).Return(&model.Product{}, nil)

		res, err := s.DeleteProduct(uint(1), uint(1))

		assert.Nil(t, res)
		assert.NotNil(t, err)
	})

	t.Run("Should return error", func(t *testing.T) {

		gormDB := testutil.MockDB()
		mockRepo1 := new(mocks.ProductRepository)
		mockRepo2 := new(mocks.ReviewRepository)
		mockRepo3 := new(mocks.ProductVariantDetailRepository)
		mockRepo4 := new(mocks.SellerRepository)
		mockRepo5 := new(mocks.SocialGraphRepository)
		mockRepo6 := new(mocks.NotificationRepository)
		cfg := &service.ProductConfig{
			DB:                gormDB,
			ProductRepo:       mockRepo1,
			ReviewRepo:        mockRepo2,
			ProductVarDetRepo: mockRepo3,
			SellerRepo:        mockRepo4,
			SocialGraphRepo:   mockRepo5,
			NotificationRepo:  mockRepo6,
		}
		s := service.NewProductService(cfg)

		mockRepo4.On("FindSellerByUserID", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("uint")).Return(&model.Seller{ID: 1}, nil)

		mockRepo1.On("FindProductByID", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("uint")).Return(&model.Product{SellerID: 2}, nil)

		mockRepo1.On("DeleteProduct", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("uint")).Return(&model.Product{}, nil)

		res, err := s.DeleteProduct(uint(1), uint(1))

		assert.Nil(t, res)
		assert.NotNil(t, err)
	})

	t.Run("Should return error", func(t *testing.T) {

		gormDB := testutil.MockDB()
		mockRepo1 := new(mocks.ProductRepository)
		mockRepo2 := new(mocks.ReviewRepository)
		mockRepo3 := new(mocks.ProductVariantDetailRepository)
		mockRepo4 := new(mocks.SellerRepository)
		mockRepo5 := new(mocks.SocialGraphRepository)
		mockRepo6 := new(mocks.NotificationRepository)
		cfg := &service.ProductConfig{
			DB:                gormDB,
			ProductRepo:       mockRepo1,
			ReviewRepo:        mockRepo2,
			ProductVarDetRepo: mockRepo3,
			SellerRepo:        mockRepo4,
			SocialGraphRepo:   mockRepo5,
			NotificationRepo:  mockRepo6,
		}
		s := service.NewProductService(cfg)

		mockRepo4.On("FindSellerByUserID", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("uint")).Return(&model.Seller{}, nil)

		mockRepo1.On("FindProductByID", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("uint")).Return(&model.Product{}, nil)

		mockRepo1.On("DeleteProduct", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("uint")).Return(nil, errors.New(""))

		res, err := s.DeleteProduct(uint(1), uint(1))

		assert.Nil(t, res)
		assert.NotNil(t, err)
	})
}
