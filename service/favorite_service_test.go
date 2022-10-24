package service_test

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"seadeals-backend/mocks"
	"seadeals-backend/model"
	"seadeals-backend/service"
	"seadeals-backend/testutil"
	"testing"
)

func TestFavoriteService_FavoriteToProduct(t *testing.T) {
	t.Run("Should return response body", func(t *testing.T) {

		gormDB := testutil.MockDB()
		mockRepo1 := new(mocks.FavoriteRepository)
		mockRepo2 := new(mocks.ProductRepository)
		cfg := &service.FavoriteServiceConfig{
			DB:                 gormDB,
			FavoriteRepository: mockRepo1,
			ProductRepository:  mockRepo2,
		}

		s := service.NewFavoriteService(cfg)
		expectedRes := &model.Favorite{}
		mockRepo1.On("FavoriteToProduct", mock.AnythingOfType(testutil.GormDBPointerType), uint(1), uint(1)).Return(&model.Favorite{}, nil)

		mockRepo2.On("UpdateProductFavoriteCount", mock.AnythingOfType(testutil.GormDBPointerType), uint(1), false).Return(&model.Product{}, nil)

		res, _, err := s.FavoriteToProduct(uint(1), uint(1))

		assert.Nil(t, err)
		assert.Equal(t, expectedRes, res)
	})

	t.Run("Should return error", func(t *testing.T) {
		gormDB := testutil.MockDB()
		mockRepo1 := new(mocks.FavoriteRepository)
		mockRepo2 := new(mocks.ProductRepository)
		cfg := &service.FavoriteServiceConfig{
			DB:                 gormDB,
			FavoriteRepository: mockRepo1,
			ProductRepository:  mockRepo2,
		}

		s := service.NewFavoriteService(cfg)
		mockRepo1.On("FavoriteToProduct", mock.AnythingOfType(testutil.GormDBPointerType), uint(1), uint(1)).Return(nil, errors.New(""))

		mockRepo2.On("UpdateProductFavoriteCount", mock.AnythingOfType(testutil.GormDBPointerType), uint(1), false).Return(&model.Product{}, nil)

		res, _, err := s.FavoriteToProduct(uint(1), uint(1))

		assert.Nil(t, res)
		assert.NotNil(t, err)
	})

	t.Run("Should return error", func(t *testing.T) {
		gormDB := testutil.MockDB()
		mockRepo1 := new(mocks.FavoriteRepository)
		mockRepo2 := new(mocks.ProductRepository)
		cfg := &service.FavoriteServiceConfig{
			DB:                 gormDB,
			FavoriteRepository: mockRepo1,
			ProductRepository:  mockRepo2,
		}

		s := service.NewFavoriteService(cfg)
		mockRepo1.On("FavoriteToProduct", mock.AnythingOfType(testutil.GormDBPointerType), uint(1), uint(1)).Return(&model.Favorite{}, nil)

		mockRepo2.On("UpdateProductFavoriteCount", mock.AnythingOfType(testutil.GormDBPointerType), uint(1), false).Return(nil, errors.New(""))

		res, _, err := s.FavoriteToProduct(uint(1), uint(1))

		assert.Nil(t, res)
		assert.NotNil(t, err)
	})
}

func TestFavoriteService_GetUserFavoriteCount(t *testing.T) {
	t.Run("Should return response body", func(t *testing.T) {

		gormDB := testutil.MockDB()
		mockRepo1 := new(mocks.FavoriteRepository)
		mockRepo2 := new(mocks.ProductRepository)
		cfg := &service.FavoriteServiceConfig{
			DB:                 gormDB,
			FavoriteRepository: mockRepo1,
			ProductRepository:  mockRepo2,
		}

		s := service.NewFavoriteService(cfg)
		expectedRes := uint(1)
		mockRepo1.On("GetUserFavoriteCount", mock.AnythingOfType(testutil.GormDBPointerType), uint(1)).Return(uint(1), nil)

		res, err := s.GetUserFavoriteCount(uint(1))

		assert.Nil(t, err)
		assert.Equal(t, expectedRes, res)
	})

	t.Run("Should return error", func(t *testing.T) {

		gormDB := testutil.MockDB()
		mockRepo1 := new(mocks.FavoriteRepository)
		mockRepo2 := new(mocks.ProductRepository)
		cfg := &service.FavoriteServiceConfig{
			DB:                 gormDB,
			FavoriteRepository: mockRepo1,
			ProductRepository:  mockRepo2,
		}

		s := service.NewFavoriteService(cfg)
		mockRepo1.On("GetUserFavoriteCount", mock.AnythingOfType(testutil.GormDBPointerType), uint(1)).Return(uint(0), errors.New(""))

		res, err := s.GetUserFavoriteCount(uint(1))

		assert.Equal(t, uint(0), res)
		assert.NotNil(t, err)
	})

}
