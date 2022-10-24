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

func TestSellerAvailableCourService_CreateOrUpdateCourier(t *testing.T) {
	t.Run("Should return response body", func(t *testing.T) {

		gormDB := testutil.MockDB()
		mockRepo1 := new(mocks.SellerAvailableCourierRepository)
		mockRepo2 := new(mocks.SellerRepository)
		cfg := &service.SellerAvailableCourServiceConfig{
			DB:                  gormDB,
			SellerAvailCourRepo: mockRepo1,
			SellerRepository:    mockRepo2,
		}
		s := service.NewSellerAvailableCourService(cfg)
		expectedRes := &model.SellerAvailableCourier{}
		mockRepo2.On("FindSellerByUserID", mock.AnythingOfType(testutil.GormDBPointerType), uint(1)).Return(&model.Seller{}, nil)

		mockBool := new(*bool)
		mockReq := &dto.AddDeliveryReq{
			CourierID:  0,
			IsSelected: *mockBool,
			SlaDay:     0,
		}
		mockRepo1.On("AddSellerAvailableDeliveryMethod", mock.AnythingOfType(testutil.GormDBPointerType), mockReq, uint(0)).Return(&model.SellerAvailableCourier{
			SellerID:   0,
			CourierID:  0,
			IsSelected: *mockBool,
			SlaDay:     0,
		}, nil)

		res, err := s.CreateOrUpdateCourier(&dto.AddDeliveryReq{
			CourierID:  0,
			IsSelected: *mockBool,
			SlaDay:     0,
		}, uint(1))

		assert.Nil(t, err)
		assert.Equal(t, expectedRes, res)
	})

	t.Run("Should return error", func(t *testing.T) {

		gormDB := testutil.MockDB()
		mockRepo1 := new(mocks.SellerAvailableCourierRepository)
		mockRepo2 := new(mocks.SellerRepository)
		cfg := &service.SellerAvailableCourServiceConfig{
			DB:                  gormDB,
			SellerAvailCourRepo: mockRepo1,
			SellerRepository:    mockRepo2,
		}
		s := service.NewSellerAvailableCourService(cfg)
		mockRepo2.On("FindSellerByUserID", mock.AnythingOfType(testutil.GormDBPointerType), uint(1)).Return(nil, errors.New(""))

		mockBool := new(*bool)
		mockReq := &dto.AddDeliveryReq{
			CourierID:  0,
			IsSelected: *mockBool,
			SlaDay:     0,
		}
		mockRepo1.On("AddSellerAvailableDeliveryMethod", mock.AnythingOfType(testutil.GormDBPointerType), mockReq, uint(0)).Return(&model.SellerAvailableCourier{
			SellerID:   0,
			CourierID:  0,
			IsSelected: *mockBool,
			SlaDay:     0,
		}, nil)

		res, err := s.CreateOrUpdateCourier(&dto.AddDeliveryReq{
			CourierID:  0,
			IsSelected: *mockBool,
			SlaDay:     0,
		}, uint(1))

		assert.Nil(t, res)
		assert.NotNil(t, err)
	})

	t.Run("Should return error", func(t *testing.T) {

		gormDB := testutil.MockDB()
		mockRepo1 := new(mocks.SellerAvailableCourierRepository)
		mockRepo2 := new(mocks.SellerRepository)
		cfg := &service.SellerAvailableCourServiceConfig{
			DB:                  gormDB,
			SellerAvailCourRepo: mockRepo1,
			SellerRepository:    mockRepo2,
		}
		s := service.NewSellerAvailableCourService(cfg)
		mockRepo2.On("FindSellerByUserID", mock.AnythingOfType(testutil.GormDBPointerType), uint(1)).Return(&model.Seller{}, nil)

		mockBool := new(*bool)
		mockReq := &dto.AddDeliveryReq{
			CourierID:  0,
			IsSelected: *mockBool,
			SlaDay:     0,
		}
		mockRepo1.On("AddSellerAvailableDeliveryMethod", mock.AnythingOfType(testutil.GormDBPointerType), mockReq, uint(0)).Return(nil, errors.New(""))

		res, err := s.CreateOrUpdateCourier(&dto.AddDeliveryReq{
			CourierID:  0,
			IsSelected: *mockBool,
			SlaDay:     0,
		}, uint(1))

		assert.Nil(t, res)
		assert.NotNil(t, err)
	})
}

func TestSellerAvailableCourService_GetAvailableCourierForSeller(t *testing.T) {
	t.Run("Should return response body", func(t *testing.T) {

		gormDB := testutil.MockDB()
		mockRepo1 := new(mocks.SellerAvailableCourierRepository)
		mockRepo2 := new(mocks.SellerRepository)
		cfg := &service.SellerAvailableCourServiceConfig{
			DB:                  gormDB,
			SellerAvailCourRepo: mockRepo1,
			SellerRepository:    mockRepo2,
		}
		s := service.NewSellerAvailableCourService(cfg)
		expectedRes := []*model.SellerAvailableCourier{}

		mockRepo2.On("FindSellerByUserID", mock.AnythingOfType(testutil.GormDBPointerType), uint(1)).Return(&model.Seller{}, nil)

		mockRepo1.On("GetAllSellerAvailableCourier", mock.AnythingOfType(testutil.GormDBPointerType), uint(0)).Return([]*model.SellerAvailableCourier{}, nil)
		res, err := s.GetAvailableCourierForSeller(uint(1))

		assert.Nil(t, err)
		assert.Equal(t, expectedRes, res)
	})

	t.Run("Should return error", func(t *testing.T) {

		gormDB := testutil.MockDB()
		mockRepo1 := new(mocks.SellerAvailableCourierRepository)
		mockRepo2 := new(mocks.SellerRepository)
		cfg := &service.SellerAvailableCourServiceConfig{
			DB:                  gormDB,
			SellerAvailCourRepo: mockRepo1,
			SellerRepository:    mockRepo2,
		}
		s := service.NewSellerAvailableCourService(cfg)

		mockRepo2.On("FindSellerByUserID", mock.AnythingOfType(testutil.GormDBPointerType), uint(1)).Return(nil, errors.New(""))

		mockRepo1.On("GetAllSellerAvailableCourier", mock.AnythingOfType(testutil.GormDBPointerType), uint(0)).Return([]*model.SellerAvailableCourier{}, nil)
		res, err := s.GetAvailableCourierForSeller(uint(1))

		assert.Nil(t, res)
		assert.NotNil(t, err)
	})

	t.Run("Should return error", func(t *testing.T) {

		gormDB := testutil.MockDB()
		mockRepo1 := new(mocks.SellerAvailableCourierRepository)
		mockRepo2 := new(mocks.SellerRepository)
		cfg := &service.SellerAvailableCourServiceConfig{
			DB:                  gormDB,
			SellerAvailCourRepo: mockRepo1,
			SellerRepository:    mockRepo2,
		}
		s := service.NewSellerAvailableCourService(cfg)

		mockRepo2.On("FindSellerByUserID", mock.AnythingOfType(testutil.GormDBPointerType), uint(1)).Return(&model.Seller{}, nil)

		mockRepo1.On("GetAllSellerAvailableCourier", mock.AnythingOfType(testutil.GormDBPointerType), uint(0)).Return(nil, errors.New(""))
		res, err := s.GetAvailableCourierForSeller(uint(1))

		assert.Nil(t, res)
		assert.NotNil(t, err)
	})
}

func TestSellerAvailableCourService_GetAvailableCourierForBuyer(t *testing.T) {
	t.Run("Should return response body", func(t *testing.T) {

		gormDB := testutil.MockDB()
		mockRepo1 := new(mocks.SellerAvailableCourierRepository)
		mockRepo2 := new(mocks.SellerRepository)
		cfg := &service.SellerAvailableCourServiceConfig{
			DB:                  gormDB,
			SellerAvailCourRepo: mockRepo1,
			SellerRepository:    mockRepo2,
		}
		s := service.NewSellerAvailableCourService(cfg)
		expectedRes := []*model.SellerAvailableCourier{}

		mockRepo1.On("GetAllSellerAvailableCourier", mock.AnythingOfType(testutil.GormDBPointerType), uint(1)).Return([]*model.SellerAvailableCourier{}, nil)

		res, err := s.GetAvailableCourierForBuyer(uint(1))

		assert.Nil(t, err)
		assert.Equal(t, expectedRes, res)
	})

	t.Run("Should return error", func(t *testing.T) {
		gormDB := testutil.MockDB()
		mockRepo1 := new(mocks.SellerAvailableCourierRepository)
		mockRepo2 := new(mocks.SellerRepository)
		cfg := &service.SellerAvailableCourServiceConfig{
			DB:                  gormDB,
			SellerAvailCourRepo: mockRepo1,
			SellerRepository:    mockRepo2,
		}
		s := service.NewSellerAvailableCourService(cfg)

		mockRepo1.On("GetAllSellerAvailableCourier", mock.AnythingOfType(testutil.GormDBPointerType), uint(1)).Return(nil, errors.New(""))

		res, err := s.GetAvailableCourierForBuyer(uint(1))

		assert.Nil(t, res)
		assert.NotNil(t, err)
	})
}
