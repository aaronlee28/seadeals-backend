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
	"time"
)

func TestVoucherService_DeleteVoucherByID(t *testing.T) {
	t.Run("Should return response body", func(t *testing.T) {

		gormDB := testutil.MockDB()
		mockRepo1 := new(mocks.VoucherRepository)
		mockRepo2 := new(mocks.SellerRepository)

		cfg := &service.VoucherServiceConfig{
			DB:          gormDB,
			VoucherRepo: mockRepo1,
			SellerRepo:  mockRepo2,
		}
		s := service.NewVoucherService(cfg)

		sellerMock := &model.Seller{UserID: 1}

		currentTime := time.Now()
		sD := currentTime.Add(time.Hour * 2)

		mockRepo1.On("FindVoucherDetailByID", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("uint")).Return(&model.Voucher{Seller: sellerMock, StartDate: sD}, nil)

		mockRepo1.On("DeleteVoucherByID", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("uint")).Return(false, nil)
		res, err := s.DeleteVoucherByID(uint(1), uint(1))

		assert.Nil(t, err)
		assert.Equal(t, false, res)
	})

	t.Run("Should return error", func(t *testing.T) {

		gormDB := testutil.MockDB()
		mockRepo1 := new(mocks.VoucherRepository)
		mockRepo2 := new(mocks.SellerRepository)

		cfg := &service.VoucherServiceConfig{
			DB:          gormDB,
			VoucherRepo: mockRepo1,
			SellerRepo:  mockRepo2,
		}
		s := service.NewVoucherService(cfg)

		mockRepo1.On("FindVoucherDetailByID", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("uint")).Return(nil, errors.New(""))

		mockRepo1.On("DeleteVoucherByID", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("uint")).Return(false, nil)
		res, err := s.DeleteVoucherByID(uint(1), uint(1))

		assert.False(t, res)
		assert.NotNil(t, err)
	})

	t.Run("Should return error", func(t *testing.T) {

		gormDB := testutil.MockDB()
		mockRepo1 := new(mocks.VoucherRepository)
		mockRepo2 := new(mocks.SellerRepository)

		cfg := &service.VoucherServiceConfig{
			DB:          gormDB,
			VoucherRepo: mockRepo1,
			SellerRepo:  mockRepo2,
		}
		s := service.NewVoucherService(cfg)

		sellerMock := &model.Seller{UserID: 2}

		currentTime := time.Now()
		sD := currentTime.Add(time.Hour * 2)

		mockRepo1.On("FindVoucherDetailByID", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("uint")).Return(&model.Voucher{Seller: sellerMock, StartDate: sD}, nil)

		mockRepo1.On("DeleteVoucherByID", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("uint")).Return(false, nil)
		res, err := s.DeleteVoucherByID(uint(1), uint(1))

		assert.False(t, res)
		assert.NotNil(t, err)
	})

	t.Run("Should return error", func(t *testing.T) {

		gormDB := testutil.MockDB()
		mockRepo1 := new(mocks.VoucherRepository)
		mockRepo2 := new(mocks.SellerRepository)

		cfg := &service.VoucherServiceConfig{
			DB:          gormDB,
			VoucherRepo: mockRepo1,
			SellerRepo:  mockRepo2,
		}
		s := service.NewVoucherService(cfg)

		sellerMock := &model.Seller{UserID: 1}

		currentTime := time.Now()
		sD := currentTime.Add(-time.Hour * 2)

		mockRepo1.On("FindVoucherDetailByID", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("uint")).Return(&model.Voucher{Seller: sellerMock, StartDate: sD}, nil)

		mockRepo1.On("DeleteVoucherByID", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("uint")).Return(true, nil)
		res, err := s.DeleteVoucherByID(uint(1), uint(1))

		assert.False(t, res)
		assert.NotNil(t, err)
	})

	t.Run("Should return error", func(t *testing.T) {

		gormDB := testutil.MockDB()
		mockRepo1 := new(mocks.VoucherRepository)
		mockRepo2 := new(mocks.SellerRepository)

		cfg := &service.VoucherServiceConfig{
			DB:          gormDB,
			VoucherRepo: mockRepo1,
			SellerRepo:  mockRepo2,
		}
		s := service.NewVoucherService(cfg)

		sellerMock := &model.Seller{UserID: 1}

		currentTime := time.Now()
		sD := currentTime.Add(time.Hour * 2)

		mockRepo1.On("FindVoucherDetailByID", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("uint")).Return(&model.Voucher{Seller: sellerMock, StartDate: sD}, nil)

		mockRepo1.On("DeleteVoucherByID", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("uint")).Return(false, errors.New(""))
		res, err := s.DeleteVoucherByID(uint(1), uint(1))

		assert.False(t, res)
		assert.NotNil(t, err)
	})
}

func TestVoucherService_GetVouchersBySellerID(t *testing.T) {
	t.Run("Should return response body", func(t *testing.T) {

		gormDB := testutil.MockDB()
		mockRepo1 := new(mocks.VoucherRepository)
		mockRepo2 := new(mocks.SellerRepository)

		cfg := &service.VoucherServiceConfig{
			DB:          gormDB,
			VoucherRepo: mockRepo1,
			SellerRepo:  mockRepo2,
		}
		s := service.NewVoucherService(cfg)

		expectedRes := []*dto.GetVoucherRes{}

		mockRepo1.On("GetVouchersBySellerID", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("uint")).Return([]*model.Voucher{}, nil)

		res, err := s.GetVouchersBySellerID(uint(1))

		assert.Nil(t, err)
		assert.Equal(t, expectedRes, res)
	})

	t.Run("Should return error", func(t *testing.T) {

		gormDB := testutil.MockDB()
		mockRepo1 := new(mocks.VoucherRepository)
		mockRepo2 := new(mocks.SellerRepository)

		cfg := &service.VoucherServiceConfig{
			DB:          gormDB,
			VoucherRepo: mockRepo1,
			SellerRepo:  mockRepo2,
		}
		s := service.NewVoucherService(cfg)

		mockRepo1.On("GetVouchersBySellerID", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("uint")).Return(nil, errors.New(""))

		res, err := s.GetVouchersBySellerID(uint(1))

		assert.Nil(t, res)
		assert.NotNil(t, err)
	})
}

func TestVoucherService_GetAvailableGlobalVouchers(t *testing.T) {
	t.Run("Should return response body", func(t *testing.T) {

		gormDB := testutil.MockDB()
		mockRepo1 := new(mocks.VoucherRepository)
		mockRepo2 := new(mocks.SellerRepository)

		cfg := &service.VoucherServiceConfig{
			DB:          gormDB,
			VoucherRepo: mockRepo1,
			SellerRepo:  mockRepo2,
		}
		s := service.NewVoucherService(cfg)
		expectedRes := []*dto.GetVoucherRes{}

		mockRepo1.On("GetAvailableGlobalVouchers", mock.AnythingOfType(testutil.GormDBPointerType)).Return([]*model.Voucher{}, nil)

		res, err := s.GetAvailableGlobalVouchers()

		assert.Nil(t, err)
		assert.Equal(t, expectedRes, res)
	})

	t.Run("Should return error", func(t *testing.T) {

		gormDB := testutil.MockDB()
		mockRepo1 := new(mocks.VoucherRepository)
		mockRepo2 := new(mocks.SellerRepository)

		cfg := &service.VoucherServiceConfig{
			DB:          gormDB,
			VoucherRepo: mockRepo1,
			SellerRepo:  mockRepo2,
		}
		s := service.NewVoucherService(cfg)

		mockRepo1.On("GetAvailableGlobalVouchers", mock.AnythingOfType(testutil.GormDBPointerType)).Return(nil, errors.New(""))

		res, err := s.GetAvailableGlobalVouchers()

		assert.Nil(t, res)
		assert.NotNil(t, err)
	})
}
