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

func TestSocialGraphService_FollowToSeller(t *testing.T) {
	t.Run("Should return response body", func(t *testing.T) {

		gormDB := testutil.MockDB()
		mockRepo1 := new(mocks.SocialGraphRepository)
		cfg := &service.SocialGraphServiceConfig{
			DB:                    gormDB,
			SocialGraphRepository: mockRepo1,
		}
		s := service.NewSocialGraphService(cfg)
		expectedRes := &model.SocialGraph{}
		mockRepo1.On("FollowToSeller", mock.AnythingOfType(testutil.GormDBPointerType), uint(1), uint(1)).Return(&model.SocialGraph{}, nil)

		res, err := s.FollowToSeller(uint(1), uint(1))

		assert.Nil(t, err)
		assert.Equal(t, expectedRes, res)
	})

	t.Run("Should return error", func(t *testing.T) {

		gormDB := testutil.MockDB()
		mockRepo1 := new(mocks.SocialGraphRepository)
		cfg := &service.SocialGraphServiceConfig{
			DB:                    gormDB,
			SocialGraphRepository: mockRepo1,
		}
		s := service.NewSocialGraphService(cfg)
		mockRepo1.On("FollowToSeller", mock.AnythingOfType(testutil.GormDBPointerType), uint(1), uint(1)).Return(nil, errors.New(""))

		res, err := s.FollowToSeller(uint(1), uint(1))

		assert.Nil(t, res)
		assert.NotNil(t, err)
	})
}
