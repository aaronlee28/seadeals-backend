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

func TestProductCategoryService_FindCategories(t *testing.T) {
	t.Run("Should return response body", func(t *testing.T) {

		gormDB := testutil.MockDB()
		mockRepo := new(mocks.ProductCategoryRepository)
		cfg := &service.ProductCategoryServiceConfig{
			DB:                        gormDB,
			ProductCategoryRepository: mockRepo,
		}
		s := service.NewProductCategoryService(cfg)

		req := &model.CategoryQuery{}
		expectedRes := []*model.ProductCategory{}
		mockRepo.On("FindCategories", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("*model.CategoryQuery")).Return([]*model.ProductCategory{}, int64(0), int64(0), nil)

		res, _, _, err := s.FindCategories(req)

		assert.Nil(t, err)
		assert.Equal(t, expectedRes, res)
	})
	t.Run("Should return error", func(t *testing.T) {

		gormDB := testutil.MockDB()
		mockRepo := new(mocks.ProductCategoryRepository)
		cfg := &service.ProductCategoryServiceConfig{
			DB:                        gormDB,
			ProductCategoryRepository: mockRepo,
		}
		s := service.NewProductCategoryService(cfg)

		req := &model.CategoryQuery{}
		mockRepo.On("FindCategories", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("*model.CategoryQuery")).Return(nil, int64(0), int64(0), errors.New(""))

		res, _, _, err := s.FindCategories(req)

		assert.Nil(t, res)
		assert.NotNil(t, err)
	})
}
