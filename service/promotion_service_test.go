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

func TestPromotionService_GetPromotionByUserID(t *testing.T) {
	t.Run("Should return response body", func(t *testing.T) {

		gormDB := testutil.MockDB()
		mockRepo1 := new(mocks.PromotionRepository)
		mockRepo2 := new(mocks.SellerRepository)
		mockRepo3 := new(mocks.ProductRepository)
		mockRepo4 := new(mocks.SocialGraphRepository)
		mockRepo5 := new(mocks.NotificationRepository)
		cfg := &service.PromotionServiceConfig{
			DB:                  gormDB,
			PromotionRepository: mockRepo1,
			SellerRepo:          mockRepo2,
			ProductRepo:         mockRepo3,
			SocialGraphRepo:     mockRepo4,
			NotificationRepo:    mockRepo5,
		}
		s := service.NewPromotionService(cfg)
		expectedRes := []*dto.GetPromotionRes{}
		mockRepo2.On("FindSellerByUserID", mock.AnythingOfType(testutil.GormDBPointerType), uint(1)).Return(&model.Seller{}, nil)

		mockRepo1.On("GetPromotionBySellerID", mock.AnythingOfType(testutil.GormDBPointerType), uint(0)).Return([]*model.Promotion{}, nil)

		mockRepo3.On("GetProductPhotoURL", mock.AnythingOfType(testutil.GormDBPointerType), uint(1)).Return("", nil)

		res, err := s.GetPromotionByUserID(uint(1))

		assert.Nil(t, err)
		assert.Equal(t, expectedRes, res)
	})

	t.Run("Should return response body", func(t *testing.T) {

		gormDB := testutil.MockDB()
		mockRepo1 := new(mocks.PromotionRepository)
		mockRepo2 := new(mocks.SellerRepository)
		mockRepo3 := new(mocks.ProductRepository)
		mockRepo4 := new(mocks.SocialGraphRepository)
		mockRepo5 := new(mocks.NotificationRepository)
		cfg := &service.PromotionServiceConfig{
			DB:                  gormDB,
			PromotionRepository: mockRepo1,
			SellerRepo:          mockRepo2,
			ProductRepo:         mockRepo3,
			SocialGraphRepo:     mockRepo4,
			NotificationRepo:    mockRepo5,
		}
		s := service.NewPromotionService(cfg)

		mockRepo2.On("FindSellerByUserID", mock.AnythingOfType(testutil.GormDBPointerType), uint(1)).Return(nil, errors.New(""))

		mockRepo1.On("GetPromotionBySellerID", mock.AnythingOfType(testutil.GormDBPointerType), uint(0)).Return([]*model.Promotion{}, nil)

		mockRepo3.On("GetProductPhotoURL", mock.AnythingOfType(testutil.GormDBPointerType), uint(1)).Return("", nil)

		res, err := s.GetPromotionByUserID(uint(1))

		assert.Nil(t, res)
		assert.NotNil(t, err)
	})

	t.Run("Should return response body", func(t *testing.T) {

		gormDB := testutil.MockDB()
		mockRepo1 := new(mocks.PromotionRepository)
		mockRepo2 := new(mocks.SellerRepository)
		mockRepo3 := new(mocks.ProductRepository)
		mockRepo4 := new(mocks.SocialGraphRepository)
		mockRepo5 := new(mocks.NotificationRepository)
		cfg := &service.PromotionServiceConfig{
			DB:                  gormDB,
			PromotionRepository: mockRepo1,
			SellerRepo:          mockRepo2,
			ProductRepo:         mockRepo3,
			SocialGraphRepo:     mockRepo4,
			NotificationRepo:    mockRepo5,
		}
		s := service.NewPromotionService(cfg)

		mockRepo2.On("FindSellerByUserID", mock.AnythingOfType(testutil.GormDBPointerType), uint(1)).Return(&model.Seller{}, nil)

		mockRepo1.On("GetPromotionBySellerID", mock.AnythingOfType(testutil.GormDBPointerType), uint(0)).Return(nil, errors.New(""))

		mockRepo3.On("GetProductPhotoURL", mock.AnythingOfType(testutil.GormDBPointerType), uint(1)).Return("", nil)

		res, err := s.GetPromotionByUserID(uint(1))

		assert.Nil(t, res)
		assert.NotNil(t, err)
	})
}
