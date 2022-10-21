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

func TestAddressService_CreateAddress(t *testing.T) {
	t.Run("Should return response body", func(t *testing.T) {

		gormDB := testutil.MockDB()
		mockRepo := new(mocks.AddressRepository)
		cfg := &service.AddressServiceConfig{
			DB:                gormDB,
			AddressRepository: mockRepo,
		}
		s := service.NewAddressService(cfg)
		req := dto.CreateAddressReq{}
		expectedRes := &model.Address{}
		mockRepo.On("CreateAddress", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("*dto.CreateAddressReq"), uint(1)).Return(&model.Address{}, nil)

		res, err := s.CreateAddress(&req, 1)

		assert.Nil(t, err)
		assert.Equal(t, expectedRes, res)
	})

	t.Run("Should return error", func(t *testing.T) {
		gormDB := testutil.MockDB()
		mockRepo := new(mocks.AddressRepository)
		cfg := &service.AddressServiceConfig{
			DB:                gormDB,
			AddressRepository: mockRepo,
		}
		s := service.NewAddressService(cfg)
		req := dto.CreateAddressReq{}
		mockRepo.On("CreateAddress", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("*dto.CreateAddressReq"), uint(1)).Return(nil, errors.New(""))

		res, err := s.CreateAddress(&req, 1)

		assert.Nil(t, res)
		assert.NotNil(t, err)
	})
}

func TestAddressService_UpdateAddress(t *testing.T) {
	t.Run("Should return response body", func(t *testing.T) {

		gormDB := testutil.MockDB()
		mockRepo := new(mocks.AddressRepository)
		cfg := &service.AddressServiceConfig{
			DB:                gormDB,
			AddressRepository: mockRepo,
		}
		s := service.NewAddressService(cfg)
		req := dto.UpdateAddressReq{}
		expectedRes := &model.Address{}
		mockRepo.On("UpdateAddress", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("*model.Address")).Return(&model.Address{}, nil)

		res, err := s.UpdateAddress(&req, 1)

		assert.Nil(t, err)
		assert.Equal(t, expectedRes, res)
	})

	t.Run("Should return error", func(t *testing.T) {
		gormDB := testutil.MockDB()
		mockRepo := new(mocks.AddressRepository)
		cfg := &service.AddressServiceConfig{
			DB:                gormDB,
			AddressRepository: mockRepo,
		}
		s := service.NewAddressService(cfg)
		req := dto.UpdateAddressReq{}
		mockRepo.On("UpdateAddress", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("*model.Address")).Return(nil, errors.New(""))

		res, err := s.UpdateAddress(&req, 1)

		assert.Nil(t, res)
		assert.NotNil(t, err)
	})
}

func TestAddressService_GetAddressesByUserID(t *testing.T) {
	t.Run("Should return response body", func(t *testing.T) {

		gormDB := testutil.MockDB()
		mockRepo := new(mocks.AddressRepository)
		cfg := &service.AddressServiceConfig{
			DB:                gormDB,
			AddressRepository: mockRepo,
		}
		s := service.NewAddressService(cfg)
		var expectedRes = []*dto.GetAddressRes{}
		mockRepo.On("GetAddressesByUserID", mock.AnythingOfType(testutil.GormDBPointerType), uint(1)).Return([]*model.Address{}, nil)

		res, err := s.GetAddressesByUserID(1)

		assert.Nil(t, err)
		assert.Equal(t, expectedRes, res)
	})

	t.Run("Should return error", func(t *testing.T) {
		gormDB := testutil.MockDB()
		mockRepo := new(mocks.AddressRepository)
		cfg := &service.AddressServiceConfig{
			DB:                gormDB,
			AddressRepository: mockRepo,
		}
		s := service.NewAddressService(cfg)
		mockRepo.On("GetAddressesByUserID", mock.AnythingOfType(testutil.GormDBPointerType), uint(1)).Return(nil, errors.New(""))

		res, err := s.GetAddressesByUserID(1)

		assert.Nil(t, res)
		assert.NotNil(t, err)
	})
}

func TestAddressService_GetUserMainAddress(t *testing.T) {
	t.Run("Should return response body", func(t *testing.T) {

		gormDB := testutil.MockDB()
		mockRepo := new(mocks.AddressRepository)
		cfg := &service.AddressServiceConfig{
			DB:                gormDB,
			AddressRepository: mockRepo,
		}
		s := service.NewAddressService(cfg)
		expectedRes := &dto.GetAddressRes{}
		mockRepo.On("GetUserMainAddress", mock.AnythingOfType(testutil.GormDBPointerType), uint(1)).Return(&model.Address{}, nil)

		res, err := s.GetUserMainAddress(1)

		assert.Nil(t, err)
		assert.Equal(t, expectedRes, res)
	})

	t.Run("Should return error", func(t *testing.T) {
		gormDB := testutil.MockDB()
		mockRepo := new(mocks.AddressRepository)
		cfg := &service.AddressServiceConfig{
			DB:                gormDB,
			AddressRepository: mockRepo,
		}
		s := service.NewAddressService(cfg)
		mockRepo.On("GetUserMainAddress", mock.AnythingOfType(testutil.GormDBPointerType), uint(1)).Return(nil, errors.New(""))

		res, err := s.GetUserMainAddress(1)

		assert.Nil(t, res)
		assert.NotNil(t, err)
	})
}

func TestAddressService_ChangeMainAddress(t *testing.T) {
	t.Run("Should return response body", func(t *testing.T) {

		gormDB := testutil.MockDB()
		mockRepo := new(mocks.AddressRepository)
		cfg := &service.AddressServiceConfig{
			DB:                gormDB,
			AddressRepository: mockRepo,
		}
		s := service.NewAddressService(cfg)
		expectedRes := &dto.GetAddressRes{}
		mockRepo.On("ChangeMainAddress", mock.AnythingOfType(testutil.GormDBPointerType), uint(1), uint(1)).Return(&model.Address{}, nil)

		res, err := s.ChangeMainAddress(1, 1)

		assert.Nil(t, err)
		assert.Equal(t, expectedRes, res)
	})

	t.Run("Should return error", func(t *testing.T) {
		gormDB := testutil.MockDB()
		mockRepo := new(mocks.AddressRepository)
		cfg := &service.AddressServiceConfig{
			DB:                gormDB,
			AddressRepository: mockRepo,
		}
		s := service.NewAddressService(cfg)
		mockRepo.On("ChangeMainAddress", mock.AnythingOfType(testutil.GormDBPointerType), uint(1), uint(1)).Return(nil, errors.New(""))

		res, err := s.ChangeMainAddress(1, 1)

		assert.Nil(t, res)
		assert.NotNil(t, err)
	})
}
