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

func TestDeliveryService_DeliverOrder(t *testing.T) {
	t.Run("Should return response body", func(t *testing.T) {

		gormDB := testutil.MockDB()
		mockRepo1 := new(mocks.DeliveryRepository)
		mockRepo2 := new(mocks.DeliveryActivityRepository)
		mockRepo3 := new(mocks.SellerRepository)
		mockRepo4 := new(mocks.AddressRepository)
		mockRepo5 := new(mocks.OrderRepository)
		mockRepo6 := new(mocks.NotificationRepository)

		cfg := &service.DeliveryServiceConfig{
			DB:                     gormDB,
			DeliveryRepository:     mockRepo1,
			DeliverActivityRepo:    mockRepo2,
			SellerRepository:       mockRepo3,
			AddressRepository:      mockRepo4,
			OrderRepository:        mockRepo5,
			NotificationRepository: mockRepo6,
		}
		s := service.NewDeliveryService(cfg)
		req := dto.DeliverOrderReq{OrderID: 1}
		expectedRes := &model.Delivery{}

		mockRepo5.On("GetOrderDetailByID", mock.AnythingOfType(testutil.GormDBPointerType), uint(1)).Return(&model.Order{Status: dto.OrderWaitingSeller, ID: 1}, nil)

		mockRepo3.On("FindSellerByUserID", mock.AnythingOfType(testutil.GormDBPointerType), uint(1)).Return(&model.Seller{}, nil)

		mockRepo1.On("GetDeliveryByOrderID", mock.AnythingOfType(testutil.GormDBPointerType), uint(1)).Return(&model.Delivery{ID: 1}, nil)

		mockRepo1.On("UpdateDeliveryStatus", mock.AnythingOfType(testutil.GormDBPointerType), uint(1), "on delivery").Return(&model.Delivery{}, nil)

		mockRepo2.On("CreateActivity", mock.AnythingOfType(testutil.GormDBPointerType), uint(1), mock.AnythingOfType("string")).Return(&model.DeliveryActivity{}, nil)

		mockRepo5.On("UpdateOrderStatus", mock.AnythingOfType(testutil.GormDBPointerType), uint(1), "on delivery").Return(&model.Order{}, nil)

		mockRepo6.On("AddToNotificationFromModel", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("*model.Notification"))

		res, err := s.DeliverOrder(&req, uint(1))

		assert.Nil(t, err)
		assert.Equal(t, expectedRes, res)
	})

	t.Run("Should return response body", func(t *testing.T) {

		gormDB := testutil.MockDB()
		mockRepo1 := new(mocks.DeliveryRepository)
		mockRepo2 := new(mocks.DeliveryActivityRepository)
		mockRepo3 := new(mocks.SellerRepository)
		mockRepo4 := new(mocks.AddressRepository)
		mockRepo5 := new(mocks.OrderRepository)
		mockRepo6 := new(mocks.NotificationRepository)

		cfg := &service.DeliveryServiceConfig{
			DB:                     gormDB,
			DeliveryRepository:     mockRepo1,
			DeliverActivityRepo:    mockRepo2,
			SellerRepository:       mockRepo3,
			AddressRepository:      mockRepo4,
			OrderRepository:        mockRepo5,
			NotificationRepository: mockRepo6,
		}
		s := service.NewDeliveryService(cfg)
		req := dto.DeliverOrderReq{OrderID: 1}

		mockRepo5.On("GetOrderDetailByID", mock.AnythingOfType(testutil.GormDBPointerType), uint(1)).Return(nil, errors.New(""))

		mockRepo3.On("FindSellerByUserID", mock.AnythingOfType(testutil.GormDBPointerType), uint(1)).Return(&model.Seller{}, nil)

		mockRepo1.On("GetDeliveryByOrderID", mock.AnythingOfType(testutil.GormDBPointerType), uint(1)).Return(&model.Delivery{ID: 1}, nil)

		mockRepo1.On("UpdateDeliveryStatus", mock.AnythingOfType(testutil.GormDBPointerType), uint(1), "on delivery").Return(&model.Delivery{}, nil)

		mockRepo2.On("CreateActivity", mock.AnythingOfType(testutil.GormDBPointerType), uint(1), mock.AnythingOfType("string")).Return(&model.DeliveryActivity{}, nil)

		mockRepo5.On("UpdateOrderStatus", mock.AnythingOfType(testutil.GormDBPointerType), uint(1), "on delivery").Return(&model.Order{}, nil)

		mockRepo6.On("AddToNotificationFromModel", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("*model.Notification"))

		res, err := s.DeliverOrder(&req, uint(1))

		assert.Nil(t, res)
		assert.NotNil(t, err)
	})

	t.Run("Should return response body", func(t *testing.T) {

		gormDB := testutil.MockDB()
		mockRepo1 := new(mocks.DeliveryRepository)
		mockRepo2 := new(mocks.DeliveryActivityRepository)
		mockRepo3 := new(mocks.SellerRepository)
		mockRepo4 := new(mocks.AddressRepository)
		mockRepo5 := new(mocks.OrderRepository)
		mockRepo6 := new(mocks.NotificationRepository)

		cfg := &service.DeliveryServiceConfig{
			DB:                     gormDB,
			DeliveryRepository:     mockRepo1,
			DeliverActivityRepo:    mockRepo2,
			SellerRepository:       mockRepo3,
			AddressRepository:      mockRepo4,
			OrderRepository:        mockRepo5,
			NotificationRepository: mockRepo6,
		}
		s := service.NewDeliveryService(cfg)
		req := dto.DeliverOrderReq{OrderID: 1}

		mockRepo5.On("GetOrderDetailByID", mock.AnythingOfType(testutil.GormDBPointerType), uint(1)).Return(&model.Order{Status: dto.OrderWaitingSeller, ID: 1}, nil)

		mockRepo3.On("FindSellerByUserID", mock.AnythingOfType(testutil.GormDBPointerType), uint(1)).Return(nil, errors.New(""))

		mockRepo1.On("GetDeliveryByOrderID", mock.AnythingOfType(testutil.GormDBPointerType), uint(1)).Return(&model.Delivery{ID: 1}, nil)

		mockRepo1.On("UpdateDeliveryStatus", mock.AnythingOfType(testutil.GormDBPointerType), uint(1), "on delivery").Return(&model.Delivery{}, nil)

		mockRepo2.On("CreateActivity", mock.AnythingOfType(testutil.GormDBPointerType), uint(1), mock.AnythingOfType("string")).Return(&model.DeliveryActivity{}, nil)

		mockRepo5.On("UpdateOrderStatus", mock.AnythingOfType(testutil.GormDBPointerType), uint(1), "on delivery").Return(&model.Order{}, nil)

		mockRepo6.On("AddToNotificationFromModel", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("*model.Notification"))

		res, err := s.DeliverOrder(&req, uint(1))

		assert.Nil(t, res)
		assert.NotNil(t, err)
	})

	t.Run("Should return response body", func(t *testing.T) {

		gormDB := testutil.MockDB()
		mockRepo1 := new(mocks.DeliveryRepository)
		mockRepo2 := new(mocks.DeliveryActivityRepository)
		mockRepo3 := new(mocks.SellerRepository)
		mockRepo4 := new(mocks.AddressRepository)
		mockRepo5 := new(mocks.OrderRepository)
		mockRepo6 := new(mocks.NotificationRepository)

		cfg := &service.DeliveryServiceConfig{
			DB:                     gormDB,
			DeliveryRepository:     mockRepo1,
			DeliverActivityRepo:    mockRepo2,
			SellerRepository:       mockRepo3,
			AddressRepository:      mockRepo4,
			OrderRepository:        mockRepo5,
			NotificationRepository: mockRepo6,
		}
		s := service.NewDeliveryService(cfg)
		req := dto.DeliverOrderReq{OrderID: 1}

		mockRepo5.On("GetOrderDetailByID", mock.AnythingOfType(testutil.GormDBPointerType), uint(1)).Return(&model.Order{Status: dto.OrderWaitingSeller, ID: 1}, nil)

		mockRepo3.On("FindSellerByUserID", mock.AnythingOfType(testutil.GormDBPointerType), uint(1)).Return(&model.Seller{}, nil)

		mockRepo1.On("GetDeliveryByOrderID", mock.AnythingOfType(testutil.GormDBPointerType), uint(1)).Return(nil, errors.New(""))

		mockRepo1.On("UpdateDeliveryStatus", mock.AnythingOfType(testutil.GormDBPointerType), uint(1), "on delivery").Return(&model.Delivery{}, nil)

		mockRepo2.On("CreateActivity", mock.AnythingOfType(testutil.GormDBPointerType), uint(1), mock.AnythingOfType("string")).Return(&model.DeliveryActivity{}, nil)

		mockRepo5.On("UpdateOrderStatus", mock.AnythingOfType(testutil.GormDBPointerType), uint(1), "on delivery").Return(&model.Order{}, nil)

		mockRepo6.On("AddToNotificationFromModel", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("*model.Notification"))

		res, err := s.DeliverOrder(&req, uint(1))

		assert.Nil(t, res)
		assert.NotNil(t, err)
	})

	t.Run("Should return response body", func(t *testing.T) {

		gormDB := testutil.MockDB()
		mockRepo1 := new(mocks.DeliveryRepository)
		mockRepo2 := new(mocks.DeliveryActivityRepository)
		mockRepo3 := new(mocks.SellerRepository)
		mockRepo4 := new(mocks.AddressRepository)
		mockRepo5 := new(mocks.OrderRepository)
		mockRepo6 := new(mocks.NotificationRepository)

		cfg := &service.DeliveryServiceConfig{
			DB:                     gormDB,
			DeliveryRepository:     mockRepo1,
			DeliverActivityRepo:    mockRepo2,
			SellerRepository:       mockRepo3,
			AddressRepository:      mockRepo4,
			OrderRepository:        mockRepo5,
			NotificationRepository: mockRepo6,
		}
		s := service.NewDeliveryService(cfg)
		req := dto.DeliverOrderReq{OrderID: 1}

		mockRepo5.On("GetOrderDetailByID", mock.AnythingOfType(testutil.GormDBPointerType), uint(1)).Return(&model.Order{Status: dto.OrderWaitingSeller, ID: 1}, nil)

		mockRepo3.On("FindSellerByUserID", mock.AnythingOfType(testutil.GormDBPointerType), uint(1)).Return(&model.Seller{}, nil)

		mockRepo1.On("GetDeliveryByOrderID", mock.AnythingOfType(testutil.GormDBPointerType), uint(1)).Return(&model.Delivery{ID: 1}, nil)

		mockRepo1.On("UpdateDeliveryStatus", mock.AnythingOfType(testutil.GormDBPointerType), uint(1), "on delivery").Return(nil, errors.New(""))

		mockRepo2.On("CreateActivity", mock.AnythingOfType(testutil.GormDBPointerType), uint(1), mock.AnythingOfType("string")).Return(&model.DeliveryActivity{}, nil)

		mockRepo5.On("UpdateOrderStatus", mock.AnythingOfType(testutil.GormDBPointerType), uint(1), "on delivery").Return(&model.Order{}, nil)

		mockRepo6.On("AddToNotificationFromModel", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("*model.Notification"))

		res, err := s.DeliverOrder(&req, uint(1))

		assert.Nil(t, res)
		assert.NotNil(t, err)
	})

	t.Run("Should return response body", func(t *testing.T) {

		gormDB := testutil.MockDB()
		mockRepo1 := new(mocks.DeliveryRepository)
		mockRepo2 := new(mocks.DeliveryActivityRepository)
		mockRepo3 := new(mocks.SellerRepository)
		mockRepo4 := new(mocks.AddressRepository)
		mockRepo5 := new(mocks.OrderRepository)
		mockRepo6 := new(mocks.NotificationRepository)

		cfg := &service.DeliveryServiceConfig{
			DB:                     gormDB,
			DeliveryRepository:     mockRepo1,
			DeliverActivityRepo:    mockRepo2,
			SellerRepository:       mockRepo3,
			AddressRepository:      mockRepo4,
			OrderRepository:        mockRepo5,
			NotificationRepository: mockRepo6,
		}
		s := service.NewDeliveryService(cfg)
		req := dto.DeliverOrderReq{OrderID: 1}

		mockRepo5.On("GetOrderDetailByID", mock.AnythingOfType(testutil.GormDBPointerType), uint(1)).Return(&model.Order{Status: dto.OrderWaitingSeller, ID: 1}, nil)

		mockRepo3.On("FindSellerByUserID", mock.AnythingOfType(testutil.GormDBPointerType), uint(1)).Return(&model.Seller{}, nil)

		mockRepo1.On("GetDeliveryByOrderID", mock.AnythingOfType(testutil.GormDBPointerType), uint(1)).Return(&model.Delivery{ID: 1}, nil)

		mockRepo1.On("UpdateDeliveryStatus", mock.AnythingOfType(testutil.GormDBPointerType), uint(1), "on delivery").Return(&model.Delivery{}, nil)

		mockRepo2.On("CreateActivity", mock.AnythingOfType(testutil.GormDBPointerType), uint(1), mock.AnythingOfType("string")).Return(nil, errors.New(""))

		mockRepo5.On("UpdateOrderStatus", mock.AnythingOfType(testutil.GormDBPointerType), uint(1), "on delivery").Return(&model.Order{}, nil)

		mockRepo6.On("AddToNotificationFromModel", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("*model.Notification"))

		res, err := s.DeliverOrder(&req, uint(1))

		assert.Nil(t, res)
		assert.NotNil(t, err)
	})

	t.Run("Should return response body", func(t *testing.T) {

		gormDB := testutil.MockDB()
		mockRepo1 := new(mocks.DeliveryRepository)
		mockRepo2 := new(mocks.DeliveryActivityRepository)
		mockRepo3 := new(mocks.SellerRepository)
		mockRepo4 := new(mocks.AddressRepository)
		mockRepo5 := new(mocks.OrderRepository)
		mockRepo6 := new(mocks.NotificationRepository)

		cfg := &service.DeliveryServiceConfig{
			DB:                     gormDB,
			DeliveryRepository:     mockRepo1,
			DeliverActivityRepo:    mockRepo2,
			SellerRepository:       mockRepo3,
			AddressRepository:      mockRepo4,
			OrderRepository:        mockRepo5,
			NotificationRepository: mockRepo6,
		}
		s := service.NewDeliveryService(cfg)
		req := dto.DeliverOrderReq{OrderID: 1}

		mockRepo5.On("GetOrderDetailByID", mock.AnythingOfType(testutil.GormDBPointerType), uint(1)).Return(&model.Order{Status: dto.OrderWaitingSeller, ID: 1}, nil)

		mockRepo3.On("FindSellerByUserID", mock.AnythingOfType(testutil.GormDBPointerType), uint(1)).Return(&model.Seller{}, nil)

		mockRepo1.On("GetDeliveryByOrderID", mock.AnythingOfType(testutil.GormDBPointerType), uint(1)).Return(&model.Delivery{ID: 1}, nil)

		mockRepo1.On("UpdateDeliveryStatus", mock.AnythingOfType(testutil.GormDBPointerType), uint(1), "on delivery").Return(&model.Delivery{}, nil)

		mockRepo2.On("CreateActivity", mock.AnythingOfType(testutil.GormDBPointerType), uint(1), mock.AnythingOfType("string")).Return(&model.DeliveryActivity{}, nil)

		mockRepo5.On("UpdateOrderStatus", mock.AnythingOfType(testutil.GormDBPointerType), uint(1), "on delivery").Return(nil, errors.New(""))

		mockRepo6.On("AddToNotificationFromModel", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("*model.Notification"))

		res, err := s.DeliverOrder(&req, uint(1))

		assert.Nil(t, res)
		assert.NotNil(t, err)
	})

	t.Run("Should return response body", func(t *testing.T) {

		gormDB := testutil.MockDB()
		mockRepo1 := new(mocks.DeliveryRepository)
		mockRepo2 := new(mocks.DeliveryActivityRepository)
		mockRepo3 := new(mocks.SellerRepository)
		mockRepo4 := new(mocks.AddressRepository)
		mockRepo5 := new(mocks.OrderRepository)
		mockRepo6 := new(mocks.NotificationRepository)

		cfg := &service.DeliveryServiceConfig{
			DB:                     gormDB,
			DeliveryRepository:     mockRepo1,
			DeliverActivityRepo:    mockRepo2,
			SellerRepository:       mockRepo3,
			AddressRepository:      mockRepo4,
			OrderRepository:        mockRepo5,
			NotificationRepository: mockRepo6,
		}
		s := service.NewDeliveryService(cfg)
		req := dto.DeliverOrderReq{OrderID: 1}

		mockRepo5.On("GetOrderDetailByID", mock.AnythingOfType(testutil.GormDBPointerType), uint(1)).Return(&model.Order{Status: dto.DeliveryFailed, ID: 1}, nil)

		mockRepo3.On("FindSellerByUserID", mock.AnythingOfType(testutil.GormDBPointerType), uint(1)).Return(&model.Seller{}, nil)

		mockRepo1.On("GetDeliveryByOrderID", mock.AnythingOfType(testutil.GormDBPointerType), uint(1)).Return(&model.Delivery{ID: 1}, nil)

		mockRepo1.On("UpdateDeliveryStatus", mock.AnythingOfType(testutil.GormDBPointerType), uint(1), "on delivery").Return(&model.Delivery{}, nil)

		mockRepo2.On("CreateActivity", mock.AnythingOfType(testutil.GormDBPointerType), uint(1), mock.AnythingOfType("string")).Return(&model.DeliveryActivity{}, nil)

		mockRepo5.On("UpdateOrderStatus", mock.AnythingOfType(testutil.GormDBPointerType), uint(1), "on delivery").Return(&model.Order{}, nil)

		mockRepo6.On("AddToNotificationFromModel", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("*model.Notification"))

		res, err := s.DeliverOrder(&req, uint(1))

		assert.Nil(t, res)
		assert.NotNil(t, err)
	})

	t.Run("Should return response body", func(t *testing.T) {

		gormDB := testutil.MockDB()
		mockRepo1 := new(mocks.DeliveryRepository)
		mockRepo2 := new(mocks.DeliveryActivityRepository)
		mockRepo3 := new(mocks.SellerRepository)
		mockRepo4 := new(mocks.AddressRepository)
		mockRepo5 := new(mocks.OrderRepository)
		mockRepo6 := new(mocks.NotificationRepository)

		cfg := &service.DeliveryServiceConfig{
			DB:                     gormDB,
			DeliveryRepository:     mockRepo1,
			DeliverActivityRepo:    mockRepo2,
			SellerRepository:       mockRepo3,
			AddressRepository:      mockRepo4,
			OrderRepository:        mockRepo5,
			NotificationRepository: mockRepo6,
		}
		s := service.NewDeliveryService(cfg)
		req := dto.DeliverOrderReq{OrderID: 1}

		mockRepo5.On("GetOrderDetailByID", mock.AnythingOfType(testutil.GormDBPointerType), uint(1)).Return(&model.Order{Status: dto.OrderWaitingSeller, ID: 1}, nil)

		mockRepo3.On("FindSellerByUserID", mock.AnythingOfType(testutil.GormDBPointerType), uint(1)).Return(&model.Seller{ID: 10}, nil)

		mockRepo1.On("GetDeliveryByOrderID", mock.AnythingOfType(testutil.GormDBPointerType), uint(1)).Return(&model.Delivery{ID: 1}, nil)

		mockRepo1.On("UpdateDeliveryStatus", mock.AnythingOfType(testutil.GormDBPointerType), uint(1), "on delivery").Return(&model.Delivery{}, nil)

		mockRepo2.On("CreateActivity", mock.AnythingOfType(testutil.GormDBPointerType), uint(1), mock.AnythingOfType("string")).Return(&model.DeliveryActivity{}, nil)

		mockRepo5.On("UpdateOrderStatus", mock.AnythingOfType(testutil.GormDBPointerType), uint(1), "on delivery").Return(&model.Order{}, nil)

		mockRepo6.On("AddToNotificationFromModel", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("*model.Notification"))

		res, err := s.DeliverOrder(&req, uint(1))

		assert.Nil(t, res)
		assert.NotNil(t, err)
	})
}

func TestDeliveryService_UpdatePrintSettings(t *testing.T) {

}
