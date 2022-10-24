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

func TestProductVariantService_FindAllProductVariantByProductID(t *testing.T) {
	t.Run("Should return response body", func(t *testing.T) {

		gormDB := testutil.MockDB()
		mockRepo1 := new(mocks.ProductRepository)
		mockRepo2 := new(mocks.ProductVariantRepository)
		mockRepo3 := new(mocks.ProductVariantDetailRepository)
		cfg := &service.ProductVariantServiceConfig{
			DB:                 gormDB,
			ProductRepo:        mockRepo1,
			ProductVariantRepo: mockRepo2,
			ProductVarDetRepo:  mockRepo3,
		}
		s := service.NewProductVariantService(cfg)
		var mockArray []*model.ProductVariantDetail
		mockReturn := &model.ProductVariantDetail{}
		mockArray = append(mockArray, mockReturn)

		var mockArray2 []*dto.GetProductVariantRes
		mockReturn2 := &dto.GetProductVariantRes{}
		mockArray2 = append(mockArray2, mockReturn2)
		expectedRes := &dto.ProductVariantRes{ProductVariants: mockArray2}

		mockRepo2.On("FindAllProductVariantByProductID", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("uint")).Return(mockArray, nil)

		res, err := s.FindAllProductVariantByProductID(uint(1))

		assert.Nil(t, err)
		assert.Equal(t, expectedRes, res)
	})

	t.Run("Should return response body", func(t *testing.T) {

		gormDB := testutil.MockDB()
		mockRepo1 := new(mocks.ProductRepository)
		mockRepo2 := new(mocks.ProductVariantRepository)
		mockRepo3 := new(mocks.ProductVariantDetailRepository)
		cfg := &service.ProductVariantServiceConfig{
			DB:                 gormDB,
			ProductRepo:        mockRepo1,
			ProductVariantRepo: mockRepo2,
			ProductVarDetRepo:  mockRepo3,
		}
		s := service.NewProductVariantService(cfg)
		var mockArray []*model.ProductVariantDetail
		mockReturn1 := &model.ProductVariantDetail{Price: 50000}
		mockReturn12 := &model.ProductVariantDetail{Price: 40000}
		mockArray = append(mockArray, mockReturn1, mockReturn12)

		var mockArray2 []*dto.GetProductVariantRes
		mockReturn2 := &dto.GetProductVariantRes{Price: 50000}
		mockReturn22 := &dto.GetProductVariantRes{Price: 40000}
		mockArray2 = append(mockArray2, mockReturn2, mockReturn22)

		mockRepo2.On("FindAllProductVariantByProductID", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("uint")).Return(mockArray, nil)

		_, err := s.FindAllProductVariantByProductID(uint(1))

		assert.Nil(t, err)
	})

	t.Run("Should return response body", func(t *testing.T) {

		gormDB := testutil.MockDB()
		mockRepo1 := new(mocks.ProductRepository)
		mockRepo2 := new(mocks.ProductVariantRepository)
		mockRepo3 := new(mocks.ProductVariantDetailRepository)
		cfg := &service.ProductVariantServiceConfig{
			DB:                 gormDB,
			ProductRepo:        mockRepo1,
			ProductVariantRepo: mockRepo2,
			ProductVarDetRepo:  mockRepo3,
		}
		s := service.NewProductVariantService(cfg)
		var mockArray []*model.ProductVariantDetail
		mockReturn1 := &model.ProductVariantDetail{Price: 40000}
		mockReturn12 := &model.ProductVariantDetail{Price: 50000}
		mockArray = append(mockArray, mockReturn1, mockReturn12)

		var mockArray2 []*dto.GetProductVariantRes
		mockReturn2 := &dto.GetProductVariantRes{Price: 50000}
		mockReturn22 := &dto.GetProductVariantRes{Price: 40000}
		mockArray2 = append(mockArray2, mockReturn2, mockReturn22)

		mockRepo2.On("FindAllProductVariantByProductID", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("uint")).Return(mockArray, nil)

		_, err := s.FindAllProductVariantByProductID(uint(1))

		assert.Nil(t, err)
	})

	t.Run("Should return error", func(t *testing.T) {
		gormDB := testutil.MockDB()
		mockRepo1 := new(mocks.ProductRepository)
		mockRepo2 := new(mocks.ProductVariantRepository)
		mockRepo3 := new(mocks.ProductVariantDetailRepository)
		cfg := &service.ProductVariantServiceConfig{
			DB:                 gormDB,
			ProductRepo:        mockRepo1,
			ProductVariantRepo: mockRepo2,
			ProductVarDetRepo:  mockRepo3,
		}
		s := service.NewProductVariantService(cfg)
		var mockArray []*model.ProductVariantDetail
		mockReturn := &model.ProductVariantDetail{}
		mockArray = append(mockArray, mockReturn)

		var mockArray2 []*dto.GetProductVariantRes
		mockReturn2 := &dto.GetProductVariantRes{}
		mockArray2 = append(mockArray2, mockReturn2)

		mockRepo2.On("FindAllProductVariantByProductID", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("uint")).Return(nil, errors.New(""))

		res, err := s.FindAllProductVariantByProductID(uint(1))

		assert.Nil(t, res)
		assert.NotNil(t, err)
	})

}

func TestProductVariantService_GetVariantPriceAfterPromotionByProductID(t *testing.T) {
	t.Run("Should return response body", func(t *testing.T) {

		gormDB := testutil.MockDB()
		mockRepo1 := new(mocks.ProductRepository)
		mockRepo2 := new(mocks.ProductVariantRepository)
		mockRepo3 := new(mocks.ProductVariantDetailRepository)
		cfg := &service.ProductVariantServiceConfig{
			DB:                 gormDB,
			ProductRepo:        mockRepo1,
			ProductVariantRepo: mockRepo2,
			ProductVarDetRepo:  mockRepo3,
		}
		s := service.NewProductVariantService(cfg)

		mockProduct := &model.Promotion{}
		mockRepo1.On("GetProductDetail", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("uint")).Return(&model.Product{Promotion: mockProduct}, nil)

		_, err := s.GetVariantPriceAfterPromotionByProductID(1)

		assert.Nil(t, err)
	})

	t.Run("Should return response body", func(t *testing.T) {

		gormDB := testutil.MockDB()
		mockRepo1 := new(mocks.ProductRepository)
		mockRepo2 := new(mocks.ProductVariantRepository)
		mockRepo3 := new(mocks.ProductVariantDetailRepository)
		cfg := &service.ProductVariantServiceConfig{
			DB:                 gormDB,
			ProductRepo:        mockRepo1,
			ProductVariantRepo: mockRepo2,
			ProductVarDetRepo:  mockRepo3,
		}
		s := service.NewProductVariantService(cfg)

		mockProduct := &model.Promotion{}
		var mockArrayProductVariantDetail []*model.ProductVariantDetail
		mockProductVariantDetail := &model.ProductVariantDetail{}
		mockArrayProductVariantDetail = append(mockArrayProductVariantDetail, mockProductVariantDetail)
		mockRepo1.On("GetProductDetail", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("uint")).Return(&model.Product{Promotion: mockProduct, ProductVariantDetail: mockArrayProductVariantDetail}, nil)

		_, err := s.GetVariantPriceAfterPromotionByProductID(1)

		assert.Nil(t, err)
	})

	t.Run("Should return error when promotion is nil", func(t *testing.T) {
		gormDB := testutil.MockDB()
		mockRepo1 := new(mocks.ProductRepository)
		mockRepo2 := new(mocks.ProductVariantRepository)
		mockRepo3 := new(mocks.ProductVariantDetailRepository)
		cfg := &service.ProductVariantServiceConfig{
			DB:                 gormDB,
			ProductRepo:        mockRepo1,
			ProductVariantRepo: mockRepo2,
			ProductVarDetRepo:  mockRepo3,
		}
		s := service.NewProductVariantService(cfg)

		mockRepo1.On("GetProductDetail", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("uint")).Return(&model.Product{}, nil)

		res, err := s.GetVariantPriceAfterPromotionByProductID(1)

		assert.Nil(t, res)
		assert.NotNil(t, err)
	})

	t.Run("Should return error when promotion is nil", func(t *testing.T) {
		gormDB := testutil.MockDB()
		mockRepo1 := new(mocks.ProductRepository)
		mockRepo2 := new(mocks.ProductVariantRepository)
		mockRepo3 := new(mocks.ProductVariantDetailRepository)
		cfg := &service.ProductVariantServiceConfig{
			DB:                 gormDB,
			ProductRepo:        mockRepo1,
			ProductVariantRepo: mockRepo2,
			ProductVarDetRepo:  mockRepo3,
		}
		s := service.NewProductVariantService(cfg)

		mockRepo1.On("GetProductDetail", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("uint")).Return(nil, errors.New(""))

		res, err := s.GetVariantPriceAfterPromotionByProductID(1)

		assert.Nil(t, res)
		assert.NotNil(t, err)
	})
}
