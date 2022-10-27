package service_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"seadeals-backend/dto"
	"seadeals-backend/mocks"
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

		expectedRes := []*dto.ProductRes{}
		mockRepo3.On("GetProductsBySellerID", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("*dto.SellerProductSearchQuery"), mock.AnythingOfType("uint")).Return([]*dto.SellerProductsCustomTable{}, int64(0), int64(0), nil)

		res, _, _, err := s.GetProductsBySellerID(&dto.SellerProductSearchQuery{}, uint(1))

		assert.Nil(t, err)
		assert.Equal(t, expectedRes, res)
	})

}
