package service_test

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
	"seadeals-backend/dto"
	"seadeals-backend/mocks"
	"seadeals-backend/model"
	"seadeals-backend/service"
	"seadeals-backend/testutil"
	"testing"
	"time"
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

func TestPromotionService_CreatePromotion(t *testing.T) {
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
		req := dto.CreatePromotionArrayReq{}

		var expectedRes []*dto.CreatePromotionRes

		mockRepo2.On("FindSellerByUserID", mock.AnythingOfType(testutil.GormDBPointerType), uint(1)).Return(&model.Seller{}, nil)

		mockRepo3.On("GetProductDetail", mock.AnythingOfType(testutil.GormDBPointerType), uint(1)).Return(&model.Seller{}, nil)

		mockRepo1.On("CreatePromotion", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("dto.CreatePromotionReq"), uint(1)).Return(&model.Seller{}, nil)

		mockRepo4.On("GetFollowerUserID", mock.AnythingOfType(testutil.GormDBPointerType), uint(1)).Return([]*model.SocialGraph{}, nil)

		mockRepo5.On("AddToNotificationFromModel", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("*model.Notification"))

		res, err := s.CreatePromotion(uint(1), &req)

		assert.Nil(t, err)
		assert.Equal(t, expectedRes, res)
	})

	t.Run("Should return error", func(t *testing.T) {

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
		req := dto.CreatePromotionArrayReq{}

		mockRepo2.On("FindSellerByUserID", mock.AnythingOfType(testutil.GormDBPointerType), uint(1)).Return(nil, errors.New(""))

		mockRepo3.On("GetProductDetail", mock.AnythingOfType(testutil.GormDBPointerType), uint(1)).Return(&model.Seller{}, nil)

		mockRepo1.On("CreatePromotion", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("dto.CreatePromotionReq"), uint(1)).Return(&model.Seller{}, nil)

		mockRepo4.On("GetFollowerUserID", mock.AnythingOfType(testutil.GormDBPointerType), uint(1)).Return([]*model.SocialGraph{}, nil)

		mockRepo5.On("AddToNotificationFromModel", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("*model.Notification"))

		res, err := s.CreatePromotion(uint(1), &req)

		assert.Nil(t, res)
		assert.NotNil(t, err)
	})

	t.Run("Should return error", func(t *testing.T) {

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
		createPromotionReqArrayMock := []dto.CreatePromotionReq{}
		createPromotionMock := dto.CreatePromotionReq{
			ProductID:   0,
			Name:        "",
			Description: "",
			StartDate:   time.Time{},
			EndDate:     time.Time{},
			Quota:       0,
			MaxOrder:    0,
			AmountType:  "",
			Amount:      0,
			BannerURL:   "",
		}
		createPromotionReqArrayMock = append(createPromotionReqArrayMock, createPromotionMock)
		req := dto.CreatePromotionArrayReq{createPromotionReqArrayMock}

		mockRepo2.On("FindSellerByUserID", mock.AnythingOfType(testutil.GormDBPointerType), uint(1)).Return(&model.Seller{}, nil)

		mockRepo3.On("GetProductDetail", mock.AnythingOfType(testutil.GormDBPointerType), uint(1)).Return(&model.Seller{}, nil)

		mockRepo1.On("CreatePromotion", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("dto.CreatePromotionReq"), uint(1)).Return(&model.Seller{}, nil)

		mockRepo4.On("GetFollowerUserID", mock.AnythingOfType(testutil.GormDBPointerType), uint(1)).Return([]*model.SocialGraph{}, nil)

		mockRepo5.On("AddToNotificationFromModel", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("*model.Notification"))

		res, err := s.CreatePromotion(uint(1), &req)

		assert.Nil(t, res)
		assert.NotNil(t, err)
	})

	t.Run("Should return error", func(t *testing.T) {

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
		createPromotionReqArrayMock := []dto.CreatePromotionReq{}
		createPromotionMock := dto.CreatePromotionReq{
			ProductID:   0,
			Name:        "",
			Description: "",
			StartDate:   time.Time{},
			EndDate:     time.Time{},
			Quota:       0,
			MaxOrder:    0,
			AmountType:  "percentage",
			Amount:      101,
			BannerURL:   "",
		}
		createPromotionReqArrayMock = append(createPromotionReqArrayMock, createPromotionMock)
		req := dto.CreatePromotionArrayReq{createPromotionReqArrayMock}

		mockRepo2.On("FindSellerByUserID", mock.AnythingOfType(testutil.GormDBPointerType), uint(1)).Return(&model.Seller{}, nil)

		mockRepo3.On("GetProductDetail", mock.AnythingOfType(testutil.GormDBPointerType), uint(1)).Return(&model.Seller{}, nil)

		mockRepo1.On("CreatePromotion", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("dto.CreatePromotionReq"), uint(1)).Return(&model.Seller{}, nil)

		mockRepo4.On("GetFollowerUserID", mock.AnythingOfType(testutil.GormDBPointerType), uint(1)).Return([]*model.SocialGraph{}, nil)

		mockRepo5.On("AddToNotificationFromModel", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("*model.Notification"))

		res, err := s.CreatePromotion(uint(1), &req)

		assert.Nil(t, res)
		assert.NotNil(t, err)
	})

	t.Run("Should return error", func(t *testing.T) {

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
		createPromotionReqArrayMock := []dto.CreatePromotionReq{}
		currentTime := time.Now()
		newT := currentTime.Add(-time.Hour * 1)

		createPromotionMock := dto.CreatePromotionReq{
			ProductID:   0,
			Name:        "",
			Description: "",
			StartDate:   time.Now(),
			EndDate:     newT,
			Quota:       0,
			MaxOrder:    0,
			AmountType:  "nominal",
			Amount:      15000,
			BannerURL:   "",
		}
		createPromotionReqArrayMock = append(createPromotionReqArrayMock, createPromotionMock)
		req := dto.CreatePromotionArrayReq{createPromotionReqArrayMock}

		mockRepo2.On("FindSellerByUserID", mock.AnythingOfType(testutil.GormDBPointerType), uint(1)).Return(&model.Seller{}, nil)

		mockRepo3.On("GetProductDetail", mock.AnythingOfType(testutil.GormDBPointerType), uint(1)).Return(&model.Seller{}, nil)

		mockRepo1.On("CreatePromotion", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("dto.CreatePromotionReq"), uint(1)).Return(&model.Seller{}, nil)

		mockRepo4.On("GetFollowerUserID", mock.AnythingOfType(testutil.GormDBPointerType), uint(1)).Return([]*model.SocialGraph{}, nil)

		mockRepo5.On("AddToNotificationFromModel", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("*model.Notification"))

		res, err := s.CreatePromotion(uint(1), &req)

		assert.Nil(t, res)
		assert.NotNil(t, err)
	})
	t.Run("Should return error", func(t *testing.T) {

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
		createPromotionReqArrayMock := []dto.CreatePromotionReq{}
		createPromotionMock := dto.CreatePromotionReq{
			ProductID:   0,
			Name:        "",
			Description: "",
			StartDate:   time.Time{},
			EndDate:     time.Time{},
			Quota:       0,
			MaxOrder:    0,
			AmountType:  "percentage",
			Amount:      10,
			BannerURL:   "",
		}
		createPromotionReqArrayMock = append(createPromotionReqArrayMock, createPromotionMock)
		req := dto.CreatePromotionArrayReq{createPromotionReqArrayMock}

		mockRepo2.On("FindSellerByUserID", mock.AnythingOfType(testutil.GormDBPointerType), uint(1)).Return(&model.Seller{}, nil)

		mockRepo3.On("GetProductDetail", mock.AnythingOfType(testutil.GormDBPointerType), uint(1)).Return(nil, errors.New(""))

		mockRepo1.On("CreatePromotion", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("dto.CreatePromotionReq"), uint(1)).Return(&model.Seller{}, nil)

		mockRepo4.On("GetFollowerUserID", mock.AnythingOfType(testutil.GormDBPointerType), uint(1)).Return([]*model.SocialGraph{}, nil)

		mockRepo5.On("AddToNotificationFromModel", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("*model.Notification"))

		res, err := s.CreatePromotion(uint(1), &req)

		assert.Nil(t, res)
		assert.NotNil(t, err)
	})

	t.Run("Should return error", func(t *testing.T) {

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
		createPromotionReqArrayMock := []dto.CreatePromotionReq{}
		currentTime := time.Now()
		sD := currentTime.Add(time.Hour * 1)
		eD := currentTime.Add(time.Hour * 2)
		createPromotionMock := dto.CreatePromotionReq{
			ProductID:   0,
			Name:        "",
			Description: "",
			StartDate:   sD,
			EndDate:     eD,
			Quota:       0,
			MaxOrder:    0,
			AmountType:  "percentage",
			Amount:      10,
			BannerURL:   "",
		}
		createPromotionReqArrayMock = append(createPromotionReqArrayMock, createPromotionMock)
		req := dto.CreatePromotionArrayReq{createPromotionReqArrayMock}

		mockRepo2.On("FindSellerByUserID", mock.AnythingOfType(testutil.GormDBPointerType), uint(0)).Return(&model.Seller{}, nil)

		mockRepo3.On("GetProductDetail", mock.AnythingOfType(testutil.GormDBPointerType), uint(0)).Return(nil, errors.New(""))

		mockRepo1.On("CreatePromotion", mock.AnythingOfType(testutil.GormDBPointerType), &createPromotionMock, uint(0)).Return(&model.Promotion{}, nil)

		mockRepo4.On("GetFollowerUserID", mock.AnythingOfType(testutil.GormDBPointerType), uint(0)).Return([]*model.SocialGraph{}, nil)

		mockRepo5.On("AddToNotificationFromModel", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("*model.Notification"))

		res, err := s.CreatePromotion(uint(0), &req)

		assert.NotNil(t, err)
		assert.Nil(t, res)
	})
	t.Run("Should return error", func(t *testing.T) {

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
		createPromotionReqArrayMock := []dto.CreatePromotionReq{}
		currentTime := time.Now()
		sD := currentTime.Add(time.Hour * 2)
		eD := currentTime.Add(time.Hour * 1)
		createPromotionMock := dto.CreatePromotionReq{
			ProductID:   0,
			Name:        "",
			Description: "",
			StartDate:   sD,
			EndDate:     eD,
			Quota:       0,
			MaxOrder:    0,
			AmountType:  "percentage",
			Amount:      10,
			BannerURL:   "",
		}
		createPromotionReqArrayMock = append(createPromotionReqArrayMock, createPromotionMock)
		req := dto.CreatePromotionArrayReq{createPromotionReqArrayMock}

		mockRepo2.On("FindSellerByUserID", mock.AnythingOfType(testutil.GormDBPointerType), uint(0)).Return(&model.Seller{}, nil)

		mockRepo3.On("GetProductDetail", mock.AnythingOfType(testutil.GormDBPointerType), uint(0)).Return(nil, errors.New(""))

		mockRepo1.On("CreatePromotion", mock.AnythingOfType(testutil.GormDBPointerType), &createPromotionMock, uint(0)).Return(&model.Promotion{}, nil)

		mockRepo4.On("GetFollowerUserID", mock.AnythingOfType(testutil.GormDBPointerType), uint(0)).Return([]*model.SocialGraph{}, nil)

		mockRepo5.On("AddToNotificationFromModel", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("*model.Notification"))

		res, err := s.CreatePromotion(uint(0), &req)

		assert.NotNil(t, err)
		assert.Nil(t, res)
	})
	t.Run("Should return error", func(t *testing.T) {

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
		createPromotionReqArrayMock := []dto.CreatePromotionReq{}
		currentTime := time.Now()
		sD := currentTime.Add(time.Hour * 2)
		eD := currentTime.Add(time.Hour * 1)
		createPromotionMock := dto.CreatePromotionReq{
			ProductID:   0,
			Name:        "",
			Description: "",
			StartDate:   sD,
			EndDate:     eD,
			Quota:       0,
			MaxOrder:    0,
			AmountType:  "percentage",
			Amount:      10,
			BannerURL:   "",
		}
		createPromotionReqArrayMock = append(createPromotionReqArrayMock, createPromotionMock)
		req := dto.CreatePromotionArrayReq{createPromotionReqArrayMock}

		mockRepo2.On("FindSellerByUserID", mock.AnythingOfType(testutil.GormDBPointerType), uint(0)).Return(&model.Seller{ID: 1}, nil)

		mockRepo3.On("GetProductDetail", mock.AnythingOfType(testutil.GormDBPointerType), uint(0)).Return(&model.Product{SellerID: 2}, nil)

		mockRepo1.On("CreatePromotion", mock.AnythingOfType(testutil.GormDBPointerType), &createPromotionMock, uint(0)).Return(&model.Promotion{}, nil)

		mockRepo4.On("GetFollowerUserID", mock.AnythingOfType(testutil.GormDBPointerType), uint(0)).Return([]*model.SocialGraph{}, nil)

		mockRepo5.On("AddToNotificationFromModel", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("*model.Notification"))

		res, err := s.CreatePromotion(uint(0), &req)

		assert.NotNil(t, err)
		assert.Nil(t, res)
	})
	t.Run("Should return error", func(t *testing.T) {

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
		createPromotionReqArrayMock := []dto.CreatePromotionReq{}
		currentTime := time.Now()
		sD := currentTime.Add(time.Hour * 2)
		eD := currentTime.Add(time.Hour * 1)
		createPromotionMock := dto.CreatePromotionReq{
			ProductID:   0,
			Name:        "",
			Description: "",
			StartDate:   sD,
			EndDate:     eD,
			Quota:       0,
			MaxOrder:    0,
			AmountType:  "percentage",
			Amount:      10,
			BannerURL:   "",
		}
		createPromotionReqArrayMock = append(createPromotionReqArrayMock, createPromotionMock)
		req := dto.CreatePromotionArrayReq{createPromotionReqArrayMock}

		mockRepo2.On("FindSellerByUserID", mock.AnythingOfType(testutil.GormDBPointerType), uint(0)).Return(&model.Seller{}, nil)

		mockRepo3.On("GetProductDetail", mock.AnythingOfType(testutil.GormDBPointerType), uint(0)).Return(&model.Product{}, nil)

		mockRepo1.On("CreatePromotion", mock.AnythingOfType(testutil.GormDBPointerType), &createPromotionMock, uint(0)).Return(nil, errors.New(""))

		mockRepo4.On("GetFollowerUserID", mock.AnythingOfType(testutil.GormDBPointerType), uint(0)).Return([]*model.SocialGraph{}, nil)

		mockRepo5.On("AddToNotificationFromModel", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("*model.Notification"))

		res, err := s.CreatePromotion(uint(0), &req)

		assert.NotNil(t, err)
		assert.Nil(t, res)
	})
	t.Run("Should return error", func(t *testing.T) {

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
		createPromotionReqArrayMock := []dto.CreatePromotionReq{}
		currentTime := time.Now()
		sD := currentTime.Add(time.Hour * 1)
		eD := currentTime.Add(time.Hour * 2)
		createPromotionMock := dto.CreatePromotionReq{
			ProductID:   0,
			Name:        "",
			Description: "",
			StartDate:   sD,
			EndDate:     eD,
			Quota:       0,
			MaxOrder:    0,
			AmountType:  "percentage",
			Amount:      10,
			BannerURL:   "",
		}
		createPromotionReqArrayMock = append(createPromotionReqArrayMock, createPromotionMock)
		req := dto.CreatePromotionArrayReq{createPromotionReqArrayMock}

		mockRepo2.On("FindSellerByUserID", mock.AnythingOfType(testutil.GormDBPointerType), uint(0)).Return(&model.Seller{ID: 1}, nil)

		mockRepo3.On("GetProductDetail", mock.AnythingOfType(testutil.GormDBPointerType), uint(0)).Return(&model.Product{SellerID: 2}, nil)

		mockRepo1.On("CreatePromotion", mock.AnythingOfType(testutil.GormDBPointerType), &createPromotionMock, uint(1)).Return(&model.Promotion{}, nil)

		socialGraphMock := &model.SocialGraph{
			Model:    gorm.Model{},
			ID:       0,
			IsFollow: false,
			UserID:   0,
			User:     nil,
			SellerID: 0,
			Seller:   nil,
		}

		mockRepo4.On("GetFollowerUserID", mock.AnythingOfType(testutil.GormDBPointerType), uint(1)).Return([]*model.SocialGraph{socialGraphMock}, nil)

		mockRepo5.On("AddToNotificationFromModel", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("*model.Notification"))

		res, err := s.CreatePromotion(uint(0), &req)

		assert.Nil(t, res)
		assert.NotNil(t, err)
	})
}

func TestPromotionService_UpdatePromotion(t *testing.T) {
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
		var expectedRes []*dto.PatchPromotionRes
		req := &dto.PatchPromotionArrayReq{}

		mockRepo1.On("ViewDetailPromotionByID", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("uint")).Return(&model.Promotion{}, nil)

		mockRepo1.On("UpdatePromotion", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("uint")).Return(&model.Promotion{}, nil)

		res, err := s.UpdatePromotion(req, uint(1))

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
		var expectedRes []*dto.PatchPromotionRes
		req := &dto.PatchPromotionArrayReq{}

		mockRepo1.On("ViewDetailPromotionByID", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("uint")).Return(&model.Promotion{}, nil)

		mockRepo1.On("UpdatePromotion", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("uint")).Return(&model.Promotion{}, nil)

		res, err := s.UpdatePromotion(req, uint(1))

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

		//var expectedRes []*dto.PatchPromotionRes
		req := &dto.PatchPromotionArrayReq{}
		patchPromotionReqMock := dto.PatchPromotionReq{
			PromotionID: 0,
			Name:        "",
			Description: "",
			StartDate:   time.Time{},
			EndDate:     time.Time{},
			Quota:       0,
			MaxOrder:    0,
			AmountType:  "",
			Amount:      0,
		}

		req.PatchPromotion = append(req.PatchPromotion, patchPromotionReqMock)

		sellerMock := &model.Seller{UserID: 1}
		productMock := &model.Product{Seller: sellerMock}
		mockRepo1.On("ViewDetailPromotionByID", mock.AnythingOfType(testutil.GormDBPointerType), uint(0)).Return(&model.Promotion{Product: productMock}, nil)

		mockRepo1.On("UpdatePromotion", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("uint"), mock.AnythingOfType("*model.Promotion")).Return(&model.Promotion{}, nil)

		_, err := s.UpdatePromotion(req, uint(1))

		assert.Nil(t, err)
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

		//var expectedRes []*dto.PatchPromotionRes
		req := &dto.PatchPromotionArrayReq{}
		patchPromotionReqMock := dto.PatchPromotionReq{
			PromotionID: 0,
			Name:        "",
			Description: "",
			StartDate:   time.Time{},
			EndDate:     time.Time{},
			Quota:       0,
			MaxOrder:    0,
			AmountType:  "",
			Amount:      0,
		}

		req.PatchPromotion = append(req.PatchPromotion, patchPromotionReqMock)

		sellerMock := &model.Seller{UserID: 2}
		productMock := &model.Product{Seller: sellerMock}
		mockRepo1.On("ViewDetailPromotionByID", mock.AnythingOfType(testutil.GormDBPointerType), uint(0)).Return(&model.Promotion{Product: productMock}, nil)

		mockRepo1.On("UpdatePromotion", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("uint"), mock.AnythingOfType("*model.Promotion")).Return(&model.Promotion{}, nil)

		_, err := s.UpdatePromotion(req, uint(1))

		assert.NotNil(t, err)
	})

}
