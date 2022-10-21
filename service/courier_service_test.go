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

func TestCourierService_GetAllCouriers(t *testing.T) {
	t.Run("Should return response body", func(t *testing.T) {

		gormDB := testutil.MockDB()
		mockRepo := new(mocks.CourierRepository)
		cfg := &service.CourierServiceConfig{
			DB:                gormDB,
			CourierRepository: mockRepo,
		}
		s := service.NewCourierService(cfg)

		expectedRes := []*model.Courier{}
		mockRepo.On("GetAllCouriers", mock.AnythingOfType(testutil.GormDBPointerType)).Return([]*model.Courier{}, nil)

		res, err := s.GetAllCouriers()

		assert.Nil(t, err)
		assert.Equal(t, expectedRes, res)
	})

	t.Run("Should return error", func(t *testing.T) {

		gormDB := testutil.MockDB()
		mockRepo := new(mocks.CourierRepository)
		cfg := &service.CourierServiceConfig{
			DB:                gormDB,
			CourierRepository: mockRepo,
		}
		s := service.NewCourierService(cfg)

		mockRepo.On("GetAllCouriers", mock.AnythingOfType(testutil.GormDBPointerType)).Return(nil, errors.New(""))

		res, err := s.GetAllCouriers()

		assert.Nil(t, res)
		assert.NotNil(t, err)
	})
}
