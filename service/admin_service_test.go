package service_test

import (
	"errors"
	"github.com/icza/gog"
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

func TestAdminService_CreateGlobalVoucher(t *testing.T) {
	t.Run("Should return response body", func(t *testing.T) {

		gormDB := testutil.MockDB()
		mockRepo := new(mocks.AdminRepository)
		cfg := &service.AdminConfig{
			DB:        gormDB,
			AdminRepo: mockRepo,
		}
		s := service.NewAdminRService(cfg)
		req := dto.CreateGlobalVoucher{
			Name:        "",
			Code:        "",
			StartDate:   time.Time{},
			EndDate:     time.Time{},
			Quota:       0,
			AmountType:  "percentage",
			Amount:      0,
			MinSpending: 0,
		}
		expectedRes := &model.Voucher{}
		mockRepo.On("CreateGlobalVoucher", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("*model.Voucher")).Return(&model.Voucher{}, nil)

		res, err := s.CreateGlobalVoucher(&req)

		assert.Nil(t, err)
		assert.Equal(t, expectedRes, res)
	})

	t.Run("Should return error", func(t *testing.T) {
		gormDB := testutil.MockDB()
		mockRepo := new(mocks.AdminRepository)
		cfg := &service.AdminConfig{
			DB:        gormDB,
			AdminRepo: mockRepo,
		}
		s := service.NewAdminRService(cfg)
		req := dto.CreateGlobalVoucher{
			Name:        "",
			Code:        "",
			StartDate:   time.Time{},
			EndDate:     time.Time{},
			Quota:       0,
			AmountType:  "percentage",
			Amount:      0,
			MinSpending: 0,
		}
		mockRepo.On("CreateGlobalVoucher", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("*model.Voucher")).Return(nil, errors.New(""))

		res, err := s.CreateGlobalVoucher(&req)

		assert.Nil(t, res)
		assert.NotNil(t, err)
	})

	t.Run("Should return error if amountType != percentage", func(t *testing.T) {
		gormDB := testutil.MockDB()
		mockRepo := new(mocks.AdminRepository)
		cfg := &service.AdminConfig{
			DB:        gormDB,
			AdminRepo: mockRepo,
		}
		s := service.NewAdminRService(cfg)
		req := dto.CreateGlobalVoucher{
			Name:        "",
			Code:        "",
			StartDate:   time.Time{},
			EndDate:     time.Time{},
			Quota:       0,
			AmountType:  "test",
			Amount:      0,
			MinSpending: 0,
		}
		mockRepo.On("CreateGlobalVoucher", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("*model.Voucher")).Return(nil, errors.New(""))

		res, err := s.CreateGlobalVoucher(&req)

		assert.Nil(t, res)
		assert.NotNil(t, err)
	})

	t.Run("Should return error if amountType == percentage && >100", func(t *testing.T) {
		gormDB := testutil.MockDB()
		mockRepo := new(mocks.AdminRepository)
		cfg := &service.AdminConfig{
			DB:        gormDB,
			AdminRepo: mockRepo,
		}
		s := service.NewAdminRService(cfg)
		req := dto.CreateGlobalVoucher{
			Name:        "",
			Code:        "",
			StartDate:   time.Time{},
			EndDate:     time.Time{},
			Quota:       0,
			AmountType:  "percentage",
			Amount:      101,
			MinSpending: 0,
		}
		mockRepo.On("CreateGlobalVoucher", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("*model.Voucher")).Return(nil, errors.New(""))

		res, err := s.CreateGlobalVoucher(&req)

		assert.Nil(t, res)
		assert.NotNil(t, err)
	})
}

func TestAdminService_CreateCategory(t *testing.T) {
	t.Run("Should return response body when parentID == nil", func(t *testing.T) {

		gormDB := testutil.MockDB()
		mockRepo := new(mocks.AdminRepository)
		cfg := &service.AdminConfig{
			DB:        gormDB,
			AdminRepo: mockRepo,
		}
		s := service.NewAdminRService(cfg)
		req := dto.CreateCategory{}
		expectedRes := &model.ProductCategory{}

		mockRepo.On("GetCategoryByID", mock.AnythingOfType(testutil.GormDBPointerType), uint(1)).Return(&model.ProductCategory{}, nil)

		mockRepo.On("CreateCategory", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("*model.ProductCategory")).Return(&model.ProductCategory{}, nil)

		res, err := s.CreateCategory(&req)

		assert.Nil(t, err)
		assert.Equal(t, expectedRes, res)
	})

	t.Run("Should return response body when parentID != nil", func(t *testing.T) {

		gormDB := testutil.MockDB()
		mockRepo := new(mocks.AdminRepository)
		cfg := &service.AdminConfig{
			DB:        gormDB,
			AdminRepo: mockRepo,
		}
		s := service.NewAdminRService(cfg)
		pID := gog.Ptr(uint(1))
		req := dto.CreateCategory{
			Name:     "",
			Slug:     "",
			IconURL:  "",
			ParentID: pID,
		}
		expectedRes := &model.ProductCategory{}

		mockRepo.On("GetCategoryByID", mock.AnythingOfType(testutil.GormDBPointerType), uint(1)).Return(&model.ProductCategory{}, nil)

		mockRepo.On("CreateCategory", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("*model.ProductCategory")).Return(&model.ProductCategory{}, nil)

		res, err := s.CreateCategory(&req)

		assert.Nil(t, err)
		assert.Equal(t, expectedRes, res)
	})

	t.Run("Should return error body when parentID == nil", func(t *testing.T) {

		gormDB := testutil.MockDB()
		mockRepo := new(mocks.AdminRepository)
		cfg := &service.AdminConfig{
			DB:        gormDB,
			AdminRepo: mockRepo,
		}
		s := service.NewAdminRService(cfg)
		req := dto.CreateCategory{}

		mockRepo.On("GetCategoryByID", mock.AnythingOfType(testutil.GormDBPointerType), uint(1)).Return(nil, errors.New(""))

		mockRepo.On("CreateCategory", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("*model.ProductCategory")).Return(nil, errors.New(""))

		res, err := s.CreateCategory(&req)

		assert.Nil(t, res)
		assert.NotNil(t, err)
	})

	t.Run("Should return error body when parentID != nil", func(t *testing.T) {

		gormDB := testutil.MockDB()
		mockRepo := new(mocks.AdminRepository)
		cfg := &service.AdminConfig{
			DB:        gormDB,
			AdminRepo: mockRepo,
		}
		s := service.NewAdminRService(cfg)
		pID := gog.Ptr(uint(1))
		req := dto.CreateCategory{
			Name:     "",
			Slug:     "",
			IconURL:  "",
			ParentID: pID,
		}
		mockRepo.On("GetCategoryByID", mock.AnythingOfType(testutil.GormDBPointerType), uint(1)).Return(nil, errors.New(""))

		mockRepo.On("CreateCategory", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("*model.ProductCategory")).Return(nil, errors.New(""))

		res, err := s.CreateCategory(&req)

		assert.Nil(t, res)
		assert.NotNil(t, err)
	})
}
