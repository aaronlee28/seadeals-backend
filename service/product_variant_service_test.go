package service_test

import (
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
		expectedRes := &dto.ProductVariantPriceRes{}
		mockRepo1.On("DeleteCartItem", mock.AnythingOfType(testutil.GormDBPointerType), uint(1), uint(1)).Return(&model.CartItem{}, nil)

		res, err := s.DeleteCartItem(uint(1), uint(1))

		assert.Nil(t, err)
		assert.Equal(t, expectedRes, res)
	})
}
