package service_test

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"seadeals-backend/dto"
	"seadeals-backend/mocks"
	"seadeals-backend/model"
	"seadeals-backend/repository"
	"seadeals-backend/service"
	"seadeals-backend/testutil"
	"testing"
)

func TestCartItemService_DeleteCartItem(t *testing.T) {
	t.Run("Should return response body", func(t *testing.T) {

		gormDB := testutil.MockDB()
		mockRepo1 := new(mocks.CartItemRepository)
		mockRepo2 := new(mocks.ProductVariantDetailRepository)
		cfg := &service.CartItemServiceConfig{
			DB:                 gormDB,
			CartItemRepository: mockRepo1,
			ProductVarDetRepo:  mockRepo2,
		}
		s := service.NewCartItemService(cfg)
		expectedRes := &model.CartItem{}
		mockRepo1.On("DeleteCartItem", mock.AnythingOfType(testutil.GormDBPointerType), uint(1), uint(1)).Return(&model.CartItem{}, nil)

		res, err := s.DeleteCartItem(uint(1), uint(1))

		assert.Nil(t, err)
		assert.Equal(t, expectedRes, res)
	})

	t.Run("Should return error", func(t *testing.T) {
		gormDB := testutil.MockDB()
		mockRepo1 := new(mocks.CartItemRepository)
		mockRepo2 := new(mocks.ProductVariantDetailRepository)
		cfg := &service.CartItemServiceConfig{
			DB:                 gormDB,
			CartItemRepository: mockRepo1,
			ProductVarDetRepo:  mockRepo2,
		}
		s := service.NewCartItemService(cfg)
		mockRepo1.On("DeleteCartItem", mock.AnythingOfType(testutil.GormDBPointerType), uint(1), uint(1)).Return(nil, errors.New(""))

		res, err := s.DeleteCartItem(uint(1), uint(1))

		assert.Nil(t, res)
		assert.NotNil(t, err)
	})
}

func TestCartItemService_AddToCart(t *testing.T) {
	t.Run("Should return response body", func(t *testing.T) {

		gormDB := testutil.MockDB()
		mockRepo1 := new(mocks.CartItemRepository)
		mockRepo2 := new(mocks.ProductVariantDetailRepository)
		cfg := &service.CartItemServiceConfig{
			DB:                 gormDB,
			CartItemRepository: mockRepo1,
			ProductVarDetRepo:  mockRepo2,
		}
		s := service.NewCartItemService(cfg)
		req := dto.AddToCartReq{}
		expectedRes := &model.CartItem{}

		mockSeller := model.Seller{ID: 50}
		mockProduct := model.Product{SellerID: 2, Seller: &mockSeller}
		mockRepo2.On("GetProductVariantDetailByID", mock.AnythingOfType(testutil.GormDBPointerType), uint(0)).Return(&model.ProductVariantDetail{
			ID:        1,
			ProductID: 1,
			Product:   &mockProduct,
		}, nil)

		mockCartItem := model.CartItem{UserID: 2}

		mockRepo1.On("AddToCart", mock.AnythingOfType(testutil.GormDBPointerType), &mockCartItem).Return(&model.CartItem{}, nil)

		res, err := s.AddToCart(uint(2), &req)

		assert.Nil(t, err)
		assert.Equal(t, expectedRes, res)
	})

	t.Run("Should return error when get product variant detailed by id returns error", func(t *testing.T) {
		gormDB := testutil.MockDB()
		mockRepo1 := new(mocks.CartItemRepository)
		mockRepo2 := new(mocks.ProductVariantDetailRepository)
		cfg := &service.CartItemServiceConfig{
			DB:                 gormDB,
			CartItemRepository: mockRepo1,
			ProductVarDetRepo:  mockRepo2,
		}
		s := service.NewCartItemService(cfg)
		req := dto.AddToCartReq{}
		mockRepo2.On("GetProductVariantDetailByID", mock.AnythingOfType(testutil.GormDBPointerType), uint(0)).Return(nil, errors.New(""))

		res, err := s.AddToCart(uint(2), &req)

		assert.Nil(t, res)
		assert.NotNil(t, err)
	})

	t.Run("Should return error when add to cart returns error", func(t *testing.T) {
		gormDB := testutil.MockDB()
		mockRepo1 := new(mocks.CartItemRepository)
		mockRepo2 := new(mocks.ProductVariantDetailRepository)
		cfg := &service.CartItemServiceConfig{
			DB:                 gormDB,
			CartItemRepository: mockRepo1,
			ProductVarDetRepo:  mockRepo2,
		}
		s := service.NewCartItemService(cfg)
		req := dto.AddToCartReq{}
		mockSeller := model.Seller{ID: 50}
		mockProduct := model.Product{SellerID: 2, Seller: &mockSeller}
		mockRepo2.On("GetProductVariantDetailByID", mock.AnythingOfType(testutil.GormDBPointerType), uint(0)).Return(&model.ProductVariantDetail{
			ID:        1,
			ProductID: 1,
			Product:   &mockProduct,
		}, nil)

		mockCartItem := model.CartItem{UserID: 2}

		mockRepo1.On("AddToCart", mock.AnythingOfType(testutil.GormDBPointerType), &mockCartItem).Return(nil, errors.New(""))

		res, err := s.AddToCart(uint(2), &req)

		assert.Nil(t, res)
		assert.NotNil(t, err)
	})

	t.Run("Should return error trying to buy their own product", func(t *testing.T) {
		gormDB := testutil.MockDB()
		mockRepo1 := new(mocks.CartItemRepository)
		mockRepo2 := new(mocks.ProductVariantDetailRepository)
		cfg := &service.CartItemServiceConfig{
			DB:                 gormDB,
			CartItemRepository: mockRepo1,
			ProductVarDetRepo:  mockRepo2,
		}
		s := service.NewCartItemService(cfg)
		req := dto.AddToCartReq{}
		mockSeller := model.Seller{ID: 1, UserID: 2}
		mockProduct := model.Product{SellerID: 1, Seller: &mockSeller}
		mockRepo2.On("GetProductVariantDetailByID", mock.AnythingOfType(testutil.GormDBPointerType), uint(0)).Return(&model.ProductVariantDetail{
			ID:        1,
			ProductID: 1,
			Product:   &mockProduct,
		}, nil)

		mockCartItem := model.CartItem{UserID: 2}

		mockRepo1.On("AddToCart", mock.AnythingOfType(testutil.GormDBPointerType), &mockCartItem).Return(&model.CartItem{}, nil)

		res, err := s.AddToCart(uint(2), &req)

		assert.Nil(t, res)
		assert.NotNil(t, err)
	})

	t.Run("Should return error when stock is less than quantity", func(t *testing.T) {
		gormDB := testutil.MockDB()
		mockRepo1 := new(mocks.CartItemRepository)
		mockRepo2 := new(mocks.ProductVariantDetailRepository)
		cfg := &service.CartItemServiceConfig{
			DB:                 gormDB,
			CartItemRepository: mockRepo1,
			ProductVarDetRepo:  mockRepo2,
		}
		s := service.NewCartItemService(cfg)
		req := dto.AddToCartReq{}
		mockSeller := model.Seller{ID: 1, UserID: 5}
		mockProduct := model.Product{SellerID: 1, Seller: &mockSeller}
		mockRepo2.On("GetProductVariantDetailByID", mock.AnythingOfType(testutil.GormDBPointerType), uint(0)).Return(&model.ProductVariantDetail{
			ID:        1,
			ProductID: 1,
			Product:   &mockProduct,
			Stock:     1,
		}, nil)

		mockCartItem := model.CartItem{UserID: 2}

		mockRepo1.On("AddToCart", mock.AnythingOfType(testutil.GormDBPointerType), &mockCartItem).Return(&model.CartItem{Quantity: 5}, nil)

		res, err := s.AddToCart(uint(2), &req)

		assert.Nil(t, res)
		assert.NotNil(t, err)
	})

	t.Run("Should return error when buying less than minimum", func(t *testing.T) {
		gormDB := testutil.MockDB()
		mockRepo1 := new(mocks.CartItemRepository)
		mockRepo2 := new(mocks.ProductVariantDetailRepository)
		cfg := &service.CartItemServiceConfig{
			DB:                 gormDB,
			CartItemRepository: mockRepo1,
			ProductVarDetRepo:  mockRepo2,
		}
		s := service.NewCartItemService(cfg)
		req := dto.AddToCartReq{}
		mockSeller := model.Seller{ID: 1, UserID: 5}
		mockProduct := model.Product{SellerID: 1, Seller: &mockSeller, MinQuantity: 10}
		mockRepo2.On("GetProductVariantDetailByID", mock.AnythingOfType(testutil.GormDBPointerType), uint(0)).Return(&model.ProductVariantDetail{
			ID:        1,
			ProductID: 1,
			Product:   &mockProduct,
			Stock:     10,
		}, nil)

		mockCartItem := model.CartItem{UserID: 2}

		mockRepo1.On("AddToCart", mock.AnythingOfType(testutil.GormDBPointerType), &mockCartItem).Return(&model.CartItem{Quantity: 5}, nil)

		res, err := s.AddToCart(uint(2), &req)

		assert.Nil(t, res)
		assert.NotNil(t, err)
	})
}

func TestCartItemService_UpdateCart(t *testing.T) {
	t.Run("Should return response body", func(t *testing.T) {

		gormDB := testutil.MockDB()
		mockRepo1 := new(mocks.CartItemRepository)
		mockRepo2 := new(mocks.ProductVariantDetailRepository)
		cfg := &service.CartItemServiceConfig{
			DB:                 gormDB,
			CartItemRepository: mockRepo1,
			ProductVarDetRepo:  mockRepo2,
		}
		s := service.NewCartItemService(cfg)
		req := dto.UpdateCartItemReq{CartItemID: 16, CurrentQuantity: 10}
		//expectedRes := &model.CartItem{}

		mockRepo1.On("UpdateCart", mock.AnythingOfType(testutil.GormDBPointerType), &req, uint(1)).Return(&model.CartItem{ID: 16, Quantity: 100, ProductVariantDetailID: 1}, nil)

		mockProduct := &model.Product{MaxQuantity: 1000}

		mockRepo2.On("GetProductVariantDetailByID", mock.AnythingOfType(testutil.GormDBPointerType), uint(1)).Return(&model.ProductVariantDetail{Stock: 100, Product: mockProduct}, nil)

		_, err := s.UpdateCart(uint(1), &req)

		assert.Nil(t, err)
		//assert.Equal(t, expectedRes, &res)
	})

	t.Run("Should return error when update cart returns error", func(t *testing.T) {
		gormDB := testutil.MockDB()
		mockRepo1 := new(mocks.CartItemRepository)
		mockRepo2 := new(mocks.ProductVariantDetailRepository)
		cfg := &service.CartItemServiceConfig{
			DB:                 gormDB,
			CartItemRepository: mockRepo1,
			ProductVarDetRepo:  mockRepo2,
		}
		s := service.NewCartItemService(cfg)
		req := dto.UpdateCartItemReq{CartItemID: 16, CurrentQuantity: 1}

		mockRepo1.On("UpdateCart", mock.AnythingOfType(testutil.GormDBPointerType), &req, uint(1)).Return(nil, errors.New(""))

		mockRepo2.On("GetProductVariantDetailByID", mock.AnythingOfType(testutil.GormDBPointerType), uint(1)).Return(&model.ProductVariantDetail{}, nil)

		res, err := s.UpdateCart(uint(1), &req)

		assert.Nil(t, res)
		assert.NotNil(t, err)
	})

	t.Run("Should return error when get product variant detail id returns error", func(t *testing.T) {
		gormDB := testutil.MockDB()
		mockRepo1 := new(mocks.CartItemRepository)
		mockRepo2 := new(mocks.ProductVariantDetailRepository)
		cfg := &service.CartItemServiceConfig{
			DB:                 gormDB,
			CartItemRepository: mockRepo1,
			ProductVarDetRepo:  mockRepo2,
		}
		s := service.NewCartItemService(cfg)
		req := dto.UpdateCartItemReq{CartItemID: 16, CurrentQuantity: 1}

		mockRepo1.On("UpdateCart", mock.AnythingOfType(testutil.GormDBPointerType), &req, uint(1)).Return(&model.CartItem{ID: 16, Quantity: 1, ProductVariantDetailID: 1}, nil)

		mockRepo2.On("GetProductVariantDetailByID", mock.AnythingOfType(testutil.GormDBPointerType), uint(1)).Return(nil, errors.New(""))

		res, err := s.UpdateCart(uint(1), &req)

		assert.Nil(t, res)
		assert.NotNil(t, err)
	})

	t.Run("Should return error when stock is less than quantity", func(t *testing.T) {
		gormDB := testutil.MockDB()
		mockRepo1 := new(mocks.CartItemRepository)
		mockRepo2 := new(mocks.ProductVariantDetailRepository)
		cfg := &service.CartItemServiceConfig{
			DB:                 gormDB,
			CartItemRepository: mockRepo1,
			ProductVarDetRepo:  mockRepo2,
		}
		s := service.NewCartItemService(cfg)
		req := dto.UpdateCartItemReq{CartItemID: 16, CurrentQuantity: 1}

		mockRepo1.On("UpdateCart", mock.AnythingOfType(testutil.GormDBPointerType), &req, uint(1)).Return(&model.CartItem{ID: 16, Quantity: 10, ProductVariantDetailID: 1}, nil)
		mockProduct := &model.Product{MaxQuantity: 1000}

		mockRepo2.On("GetProductVariantDetailByID", mock.AnythingOfType(testutil.GormDBPointerType), uint(1)).Return(&model.ProductVariantDetail{Stock: 1, Product: mockProduct}, nil)

		res, err := s.UpdateCart(uint(1), &req)

		assert.Nil(t, res)
		assert.NotNil(t, err)
	})

	t.Run("Should return error when min quantity is greater than quantity", func(t *testing.T) {
		gormDB := testutil.MockDB()
		mockRepo1 := new(mocks.CartItemRepository)
		mockRepo2 := new(mocks.ProductVariantDetailRepository)
		cfg := &service.CartItemServiceConfig{
			DB:                 gormDB,
			CartItemRepository: mockRepo1,
			ProductVarDetRepo:  mockRepo2,
		}
		s := service.NewCartItemService(cfg)
		req := dto.UpdateCartItemReq{CartItemID: 16, CurrentQuantity: 1}

		mockRepo1.On("UpdateCart", mock.AnythingOfType(testutil.GormDBPointerType), &req, uint(1)).Return(&model.CartItem{ID: 16, Quantity: 10, ProductVariantDetailID: 1}, nil)
		mockProduct := &model.Product{MinQuantity: 1000}

		mockRepo2.On("GetProductVariantDetailByID", mock.AnythingOfType(testutil.GormDBPointerType), uint(1)).Return(&model.ProductVariantDetail{Stock: 100, Product: mockProduct}, nil)

		res, err := s.UpdateCart(uint(1), &req)

		assert.Nil(t, res)
		assert.NotNil(t, err)
	})
}

func TestCartItemService_GetCartItems(t *testing.T) {
	t.Run("Should return response body", func(t *testing.T) {

		gormDB := testutil.MockDB()
		mockRepo1 := new(mocks.CartItemRepository)
		mockRepo2 := new(mocks.ProductVariantDetailRepository)
		cfg := &service.CartItemServiceConfig{
			DB:                 gormDB,
			CartItemRepository: mockRepo1,
			ProductVarDetRepo:  mockRepo2,
		}
		s := service.NewCartItemService(cfg)
		expectedRes := []*dto.CartItemRes{}
		mockQuery := &repository.Query{}
		mockRepo1.On("GetCartItem", mock.AnythingOfType(testutil.GormDBPointerType), mockQuery, uint(1)).Return([]*model.CartItem{}, int64(1), int64(1), nil)

		res1, _, _, err := s.GetCartItems(mockQuery, uint(1))

		assert.Nil(t, err)
		assert.Equal(t, expectedRes, res1)
	})
	//
	//t.Run("Should return response body when orderItems is returned", func(t *testing.T) {
	//
	//	gormDB := testutil.MockDB()
	//	mockRepo1 := new(mocks.CartItemRepository)
	//	mockRepo2 := new(mocks.ProductVariantDetailRepository)
	//	cfg := &service.CartItemServiceConfig{
	//		DB:                 gormDB,
	//		CartItemRepository: mockRepo1,
	//		ProductVarDetRepo:  mockRepo2,
	//	}
	//	s := service.NewCartItemService(cfg)
	//	expectedRes := []*dto.CartItemRes{}
	//	mockQuery := &repository.Query{}
	//	mockCartItem := &model.CartItem{
	//		ID:                     1,
	//		ProductVariantDetailID: 1,
	//		ProductVariantDetail:   nil,
	//		UserID:                 0,
	//		User:                   nil,
	//		Quantity:               0,
	//	}
	//	mockRepo1.On("GetCartItem", mock.AnythingOfType(testutil.GormDBPointerType), mockQuery, uint(1)).Return(
	//		[]*model.CartItem{mockCartItem}, int64(1), int64(1), nil)
	//
	//	res1, _, _, err := s.GetCartItems(mockQuery, uint(1))
	//
	//	assert.Nil(t, err)
	//	assert.Equal(t, expectedRes, res1)
	//})

	t.Run("Should return error", func(t *testing.T) {
		gormDB := testutil.MockDB()
		mockRepo1 := new(mocks.CartItemRepository)
		mockRepo2 := new(mocks.ProductVariantDetailRepository)
		cfg := &service.CartItemServiceConfig{
			DB:                 gormDB,
			CartItemRepository: mockRepo1,
			ProductVarDetRepo:  mockRepo2,
		}
		s := service.NewCartItemService(cfg)
		mockQuery := &repository.Query{}
		mockRepo1.On("GetCartItem", mock.AnythingOfType(testutil.GormDBPointerType), mockQuery, uint(1)).Return(nil, int64(0), int64(0), errors.New(""))

		res1, _, _, err := s.GetCartItems(mockQuery, uint(1))

		assert.Nil(t, res1)
		assert.NotNil(t, err)
	})
}
