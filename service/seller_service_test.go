package service_test

import (
	"testing"
)

func TestSellerService_FindSellerByID(t *testing.T) {
	//t.Run("Should return response body", func(t *testing.T) {
	//
	//	gormDB := testutil.MockDB()
	//	mockRepo1 := new(mocks.ProductRepository)
	//	mockRepo2 := new(mocks.SellerRepository)
	//	mockRepo3 := new(mocks.ReviewRepository)
	//	mockRepo4 := new(mocks.SocialGraphRepository)
	//	cfg := &service.SellerServiceConfig{
	//		DB:              gormDB,
	//		ProductRepo:     mockRepo1,
	//		SellerRepo:      mockRepo2,
	//		ReviewRepo:      mockRepo3,
	//		SocialGraphRepo: mockRepo4,
	//	}
	//	s := service.NewSellerService(cfg)
	//
	//	expectedRes := &dto.GetSellerRes{}
	//
	//	mockRepo2.On("FindSellerDetailByID", mock.AnythingOfType(testutil.GormDBPointerType), uint(1), uint(1)).Return(&model.Seller{
	//		Name:        "test",
	//		Slug:        "test",
	//		UserID:      1,
	//		Description: "test",
	//		AddressID:   2,
	//		PictureURL:  "test",
	//		BannerURL:   "test",
	//		AllowPrint:  false,
	//	}, nil)
	//
	//	mockRepo3.On("GetReviewsAvgAndCountBySellerID", mock.AnythingOfType(testutil.GormDBPointerType), uint(1)).Return(float64(1), 1, nil)
	//
	//	mockRepo4.On("GetFollowerCountBySellerID", mock.AnythingOfType(testutil.GormDBPointerType), uint(1)).Return(int64(1), nil)
	//
	//	mockRepo4.On("GetFollowingCountByUserID", mock.AnythingOfType(testutil.GormDBPointerType), uint(1)).Return(int64(1), nil)
	//
	//	mockRepo1.On("GetProductCountBySellerID", mock.AnythingOfType(testutil.GormDBPointerType), uint(1)).Return(int64(1), nil)
	//
	//	res, err := s.FindSellerByID(uint(1), uint(1))
	//
	//	assert.Nil(t, err)
	//	assert.Equal(t, expectedRes, res)
	//})
}
