package service_test

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"seadeals-backend/dto"
	"seadeals-backend/mocks"
	"seadeals-backend/model"
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
