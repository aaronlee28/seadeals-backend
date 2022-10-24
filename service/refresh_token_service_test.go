package service_test

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"seadeals-backend/mocks"
	"seadeals-backend/service"
	"seadeals-backend/testutil"
	"testing"
)

func TestRefreshTokenService_CheckIfTokenExist(t *testing.T) {
	t.Run("Should return response body", func(t *testing.T) {

		gormDB := testutil.MockDB()
		mockRepo := new(mocks.RefreshTokenRepository)
		cfg := &service.RefreshTokenServiceConfig{
			DB:               gormDB,
			RefreshTokenRepo: mockRepo,
		}
		s := service.NewRefreshTokenService(cfg)

		mockRepo.On("CheckIfTokenExist", mock.AnythingOfType(testutil.GormDBPointerType), "").Return(true, uint(1), nil)

		_, _, err := s.CheckIfTokenExist("")

		assert.Nil(t, err)
	})

	t.Run("Should return error ", func(t *testing.T) {
		gormDB := testutil.MockDB()
		mockRepo := new(mocks.RefreshTokenRepository)
		cfg := &service.RefreshTokenServiceConfig{
			DB:               gormDB,
			RefreshTokenRepo: mockRepo,
		}
		s := service.NewRefreshTokenService(cfg)

		mockRepo.On("CheckIfTokenExist", mock.AnythingOfType(testutil.GormDBPointerType), "").Return(false, uint(1), errors.New(""))

		res, _, err := s.CheckIfTokenExist("")

		assert.False(t, res)
		assert.NotNil(t, err)
	})
}
