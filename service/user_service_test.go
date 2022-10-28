package service_test

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"seadeals-backend/config"
	"seadeals-backend/dto"
	"seadeals-backend/mocks"
	"seadeals-backend/model"
	"seadeals-backend/service"
	"seadeals-backend/testutil"
	"testing"
)

func TestUserService_Register(t *testing.T) {
	t.Run("Should return response body", func(t *testing.T) {

		gormDB := testutil.MockDB()
		mockRepo1 := new(mocks.UserRepository)
		mockRepo2 := new(mocks.UserRoleRepository)
		mockRepo3 := new(mocks.AddressRepository)
		mockRepo4 := new(mocks.WalletRepository)
		cfg := &service.UserServiceConfig{
			DB:               gormDB,
			UserRepository:   mockRepo1,
			UserRoleRepo:     mockRepo2,
			AddressRepo:      mockRepo3,
			WalletRepository: mockRepo4,
			AppConfig:        config.AppConfig{},
		}
		s := service.NewUserService(cfg)
		req := &dto.RegisterRequest{Email: "test@mail.com", Username: "13245", Password: "asdfgh"}

		mockRepo1.On("Register", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("*model.User")).Return(&model.User{}, nil)

		mockRepo4.On("CreateWallet", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("*model.Wallet")).Return(&model.Wallet{}, nil)

		mockRepo2.On("CreateRoleToUser", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("*model.UserRole")).Return(nil, nil)

		_, _, err := s.Register(req)

		assert.Nil(t, err)
	})

	t.Run("Should return error", func(t *testing.T) {

		gormDB := testutil.MockDB()
		mockRepo1 := new(mocks.UserRepository)
		mockRepo2 := new(mocks.UserRoleRepository)
		mockRepo3 := new(mocks.AddressRepository)
		mockRepo4 := new(mocks.WalletRepository)
		cfg := &service.UserServiceConfig{
			DB:               gormDB,
			UserRepository:   mockRepo1,
			UserRoleRepo:     mockRepo2,
			AddressRepo:      mockRepo3,
			WalletRepository: mockRepo4,
			AppConfig:        config.AppConfig{},
		}
		s := service.NewUserService(cfg)
		req := &dto.RegisterRequest{Email: "test@mail.com", Username: "test", Password: "test"}

		mockRepo1.On("Register", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("*model.User")).Return(&model.User{}, nil)

		mockRepo4.On("CreateWallet", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("*model.Wallet")).Return(&model.Wallet{}, nil)

		mockRepo2.On("CreateRoleToUser", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("*model.UserRole")).Return(nil, nil)

		res, _, err := s.Register(req)

		assert.Nil(t, res)
		assert.NotNil(t, err)
	})

	t.Run("Should return error", func(t *testing.T) {

		gormDB := testutil.MockDB()
		mockRepo1 := new(mocks.UserRepository)
		mockRepo2 := new(mocks.UserRoleRepository)
		mockRepo3 := new(mocks.AddressRepository)
		mockRepo4 := new(mocks.WalletRepository)
		cfg := &service.UserServiceConfig{
			DB:               gormDB,
			UserRepository:   mockRepo1,
			UserRoleRepo:     mockRepo2,
			AddressRepo:      mockRepo3,
			WalletRepository: mockRepo4,
			AppConfig:        config.AppConfig{},
		}
		s := service.NewUserService(cfg)
		req := &dto.RegisterRequest{Username: "1234", Password: "test"}

		mockRepo1.On("Register", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("*model.User")).Return(&model.User{}, nil)

		mockRepo4.On("CreateWallet", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("*model.Wallet")).Return(&model.Wallet{}, nil)

		mockRepo2.On("CreateRoleToUser", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("*model.UserRole")).Return(nil, nil)

		res, _, err := s.Register(req)

		assert.Nil(t, res)
		assert.NotNil(t, err)
	})

	t.Run("Should return error", func(t *testing.T) {

		gormDB := testutil.MockDB()
		mockRepo1 := new(mocks.UserRepository)
		mockRepo2 := new(mocks.UserRoleRepository)
		mockRepo3 := new(mocks.AddressRepository)
		mockRepo4 := new(mocks.WalletRepository)
		cfg := &service.UserServiceConfig{
			DB:               gormDB,
			UserRepository:   mockRepo1,
			UserRoleRepo:     mockRepo2,
			AddressRepo:      mockRepo3,
			WalletRepository: mockRepo4,
			AppConfig:        config.AppConfig{},
		}
		s := service.NewUserService(cfg)
		req := &dto.RegisterRequest{Email: "test@mail.com", Username: "1234", Password: "test"}

		mockRepo1.On("Register", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("*model.User")).Return(nil, errors.New(""))

		mockRepo4.On("CreateWallet", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("*model.Wallet")).Return(&model.Wallet{}, nil)

		mockRepo2.On("CreateRoleToUser", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("*model.UserRole")).Return(nil, nil)

		res, _, err := s.Register(req)

		assert.Nil(t, res)
		assert.NotNil(t, err)
	})

	t.Run("Should return error", func(t *testing.T) {

		gormDB := testutil.MockDB()
		mockRepo1 := new(mocks.UserRepository)
		mockRepo2 := new(mocks.UserRoleRepository)
		mockRepo3 := new(mocks.AddressRepository)
		mockRepo4 := new(mocks.WalletRepository)
		cfg := &service.UserServiceConfig{
			DB:               gormDB,
			UserRepository:   mockRepo1,
			UserRoleRepo:     mockRepo2,
			AddressRepo:      mockRepo3,
			WalletRepository: mockRepo4,
			AppConfig:        config.AppConfig{},
		}
		s := service.NewUserService(cfg)
		req := &dto.RegisterRequest{Email: "test@mail.com", Username: "1234", Password: "test"}

		mockRepo1.On("Register", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("*model.User")).Return(&model.User{}, nil)

		mockRepo4.On("CreateWallet", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("*model.Wallet")).Return(nil, errors.New(""))

		mockRepo2.On("CreateRoleToUser", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("*model.UserRole")).Return(nil, nil)

		res, _, err := s.Register(req)

		assert.Nil(t, res)
		assert.NotNil(t, err)
	})

	t.Run("Should return error", func(t *testing.T) {

		gormDB := testutil.MockDB()
		mockRepo1 := new(mocks.UserRepository)
		mockRepo2 := new(mocks.UserRoleRepository)
		mockRepo3 := new(mocks.AddressRepository)
		mockRepo4 := new(mocks.WalletRepository)
		cfg := &service.UserServiceConfig{
			DB:               gormDB,
			UserRepository:   mockRepo1,
			UserRoleRepo:     mockRepo2,
			AddressRepo:      mockRepo3,
			WalletRepository: mockRepo4,
			AppConfig:        config.AppConfig{},
		}
		s := service.NewUserService(cfg)
		req := &dto.RegisterRequest{Email: "test@mail.com", Username: "1234", Password: "test"}

		mockRepo1.On("Register", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("*model.User")).Return(&model.User{}, nil)

		mockRepo4.On("CreateWallet", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("*model.Wallet")).Return(&model.Wallet{}, nil)

		mockRepo2.On("CreateRoleToUser", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("*model.UserRole")).Return(nil, errors.New(""))

		res, _, err := s.Register(req)

		assert.Nil(t, res)
		assert.NotNil(t, err)
	})
}

func TestUserService_CheckGoogleAccount(t *testing.T) {
	t.Run("Should return response body", func(t *testing.T) {

		gormDB := testutil.MockDB()
		mockRepo1 := new(mocks.UserRepository)
		mockRepo2 := new(mocks.UserRoleRepository)
		mockRepo3 := new(mocks.AddressRepository)
		mockRepo4 := new(mocks.WalletRepository)
		cfg := &service.UserServiceConfig{
			DB:               gormDB,
			UserRepository:   mockRepo1,
			UserRoleRepo:     mockRepo2,
			AddressRepo:      mockRepo3,
			WalletRepository: mockRepo4,
			AppConfig:        config.AppConfig{},
		}
		s := service.NewUserService(cfg)
		expectedRes := &model.User{}

		mockRepo1.On("HasExistEmail", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("string")).Return(true, nil)

		mockRepo1.On("GetUserByEmail", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("string")).Return(&model.User{}, nil)

		res, err := s.CheckGoogleAccount("test@mail.com")

		assert.Nil(t, err)
		assert.Equal(t, expectedRes, res)

	})

	t.Run("Should return response body", func(t *testing.T) {

		gormDB := testutil.MockDB()
		mockRepo1 := new(mocks.UserRepository)
		mockRepo2 := new(mocks.UserRoleRepository)
		mockRepo3 := new(mocks.AddressRepository)
		mockRepo4 := new(mocks.WalletRepository)
		cfg := &service.UserServiceConfig{
			DB:               gormDB,
			UserRepository:   mockRepo1,
			UserRoleRepo:     mockRepo2,
			AddressRepo:      mockRepo3,
			WalletRepository: mockRepo4,
			AppConfig:        config.AppConfig{},
		}
		s := service.NewUserService(cfg)

		mockRepo1.On("HasExistEmail", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("string")).Return(true, nil)

		mockRepo1.On("GetUserByEmail", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("string")).Return(&model.User{Email: "test@mail.com"}, nil)

		res, err := s.CheckGoogleAccount("")

		assert.Nil(t, res)
		assert.NotNil(t, err)

	})

	t.Run("Should return response body", func(t *testing.T) {

		gormDB := testutil.MockDB()
		mockRepo1 := new(mocks.UserRepository)
		mockRepo2 := new(mocks.UserRoleRepository)
		mockRepo3 := new(mocks.AddressRepository)
		mockRepo4 := new(mocks.WalletRepository)
		cfg := &service.UserServiceConfig{
			DB:               gormDB,
			UserRepository:   mockRepo1,
			UserRoleRepo:     mockRepo2,
			AddressRepo:      mockRepo3,
			WalletRepository: mockRepo4,
			AppConfig:        config.AppConfig{},
		}
		s := service.NewUserService(cfg)

		mockRepo1.On("HasExistEmail", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("string")).Return(false, errors.New(""))

		mockRepo1.On("GetUserByEmail", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("string")).Return(&model.User{Email: "test@mail.com"}, nil)

		res, err := s.CheckGoogleAccount("test@mail.com")

		assert.Nil(t, res)
		assert.NotNil(t, err)

	})

	t.Run("Should return response body", func(t *testing.T) {

		gormDB := testutil.MockDB()
		mockRepo1 := new(mocks.UserRepository)
		mockRepo2 := new(mocks.UserRoleRepository)
		mockRepo3 := new(mocks.AddressRepository)
		mockRepo4 := new(mocks.WalletRepository)
		cfg := &service.UserServiceConfig{
			DB:               gormDB,
			UserRepository:   mockRepo1,
			UserRoleRepo:     mockRepo2,
			AddressRepo:      mockRepo3,
			WalletRepository: mockRepo4,
			AppConfig:        config.AppConfig{},
		}
		s := service.NewUserService(cfg)

		mockRepo1.On("HasExistEmail", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("string")).Return(false, nil)

		mockRepo1.On("GetUserByEmail", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("string")).Return(&model.User{Email: "test@mail.com"}, nil)

		res, err := s.CheckGoogleAccount("test@mail.com")

		assert.Nil(t, res)
		assert.NotNil(t, err)

	})

	t.Run("Should return response body", func(t *testing.T) {

		gormDB := testutil.MockDB()
		mockRepo1 := new(mocks.UserRepository)
		mockRepo2 := new(mocks.UserRoleRepository)
		mockRepo3 := new(mocks.AddressRepository)
		mockRepo4 := new(mocks.WalletRepository)
		cfg := &service.UserServiceConfig{
			DB:               gormDB,
			UserRepository:   mockRepo1,
			UserRoleRepo:     mockRepo2,
			AddressRepo:      mockRepo3,
			WalletRepository: mockRepo4,
			AppConfig:        config.AppConfig{},
		}
		s := service.NewUserService(cfg)

		mockRepo1.On("HasExistEmail", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("string")).Return(true, nil)

		mockRepo1.On("GetUserByEmail", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("string")).Return(nil, errors.New(""))

		res, err := s.CheckGoogleAccount("test@mail.com")

		assert.Nil(t, res)
		assert.NotNil(t, err)

	})
}

func TestUserService_RegisterAsSeller(t *testing.T) {
	t.Run("Should return response body", func(t *testing.T) {

		gormDB := testutil.MockDB()
		mockRepo1 := new(mocks.UserRepository)
		mockRepo2 := new(mocks.UserRoleRepository)
		mockRepo3 := new(mocks.AddressRepository)
		mockRepo4 := new(mocks.WalletRepository)
		cfg := &service.UserServiceConfig{
			DB:               gormDB,
			UserRepository:   mockRepo1,
			UserRoleRepo:     mockRepo2,
			AddressRepo:      mockRepo3,
			WalletRepository: mockRepo4,
			AppConfig:        config.AppConfig{},
		}
		s := service.NewUserService(cfg)
		req := &dto.RegisterAsSellerReq{}
		expectedRes := &model.Seller{}

		mockRepo1.On("GetUserByID", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("uint")).Return(&model.User{ID: 1}, nil)

		mockRepo3.On("GetUserMainAddress", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("uint")).Return(&model.Address{}, nil)

		mockRepo2.On("CreateRoleToUser", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("*model.UserRole")).Return(nil, nil)

		mockRepo1.On("RegisterAsSeller", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("*model.Seller")).Return(&model.Seller{}, nil)

		mockRepo4.On("GetWalletByUserID", mock.AnythingOfType(testutil.GormDBPointerType), uint(1)).Return(&model.Wallet{}, nil)

		mockRepo2.On("GetRolesByUserID", mock.AnythingOfType(testutil.GormDBPointerType), uint(1)).Return([]*model.UserRole{}, nil)

		res, _, err := s.RegisterAsSeller(req, uint(1))

		assert.Nil(t, err)
		assert.Equal(t, expectedRes, res)

	})

	t.Run("Should return response body", func(t *testing.T) {

		gormDB := testutil.MockDB()
		mockRepo1 := new(mocks.UserRepository)
		mockRepo2 := new(mocks.UserRoleRepository)
		mockRepo3 := new(mocks.AddressRepository)
		mockRepo4 := new(mocks.WalletRepository)
		cfg := &service.UserServiceConfig{
			DB:               gormDB,
			UserRepository:   mockRepo1,
			UserRoleRepo:     mockRepo2,
			AddressRepo:      mockRepo3,
			WalletRepository: mockRepo4,
			AppConfig:        config.AppConfig{},
		}
		s := service.NewUserService(cfg)
		req := &dto.RegisterAsSellerReq{}
		expectedRes := &model.Seller{}

		mockRepo1.On("GetUserByID", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("uint")).Return(&model.User{ID: 1}, nil)

		mockRepo3.On("GetUserMainAddress", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("uint")).Return(&model.Address{}, nil)

		mockRepo2.On("CreateRoleToUser", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("*model.UserRole")).Return(nil, nil)

		mockRepo1.On("RegisterAsSeller", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("*model.Seller")).Return(&model.Seller{}, nil)

		mockRepo4.On("GetWalletByUserID", mock.AnythingOfType(testutil.GormDBPointerType), uint(1)).Return(&model.Wallet{}, nil)

		mockUserRole := &model.UserRole{UserID: 1}

		mockRepo2.On("GetRolesByUserID", mock.AnythingOfType(testutil.GormDBPointerType), uint(1)).Return([]*model.UserRole{mockUserRole}, nil)

		res, _, err := s.RegisterAsSeller(req, uint(1))

		assert.Nil(t, err)
		assert.Equal(t, expectedRes, res)

	})
	t.Run("Should return error", func(t *testing.T) {

		gormDB := testutil.MockDB()
		mockRepo1 := new(mocks.UserRepository)
		mockRepo2 := new(mocks.UserRoleRepository)
		mockRepo3 := new(mocks.AddressRepository)
		mockRepo4 := new(mocks.WalletRepository)
		cfg := &service.UserServiceConfig{
			DB:               gormDB,
			UserRepository:   mockRepo1,
			UserRoleRepo:     mockRepo2,
			AddressRepo:      mockRepo3,
			WalletRepository: mockRepo4,
			AppConfig:        config.AppConfig{},
		}
		s := service.NewUserService(cfg)
		req := &dto.RegisterAsSellerReq{}

		mockRepo1.On("GetUserByID", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("uint")).Return(nil, errors.New(""))

		mockRepo3.On("GetUserMainAddress", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("uint")).Return(&model.Address{}, nil)

		mockRepo2.On("CreateRoleToUser", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("*model.UserRole")).Return(nil, nil)

		mockRepo1.On("RegisterAsSeller", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("*model.Seller")).Return(&model.Seller{}, nil)

		mockRepo4.On("GetWalletByUserID", mock.AnythingOfType(testutil.GormDBPointerType), uint(1)).Return(&model.Wallet{}, nil)

		mockRepo2.On("GetRolesByUserID", mock.AnythingOfType(testutil.GormDBPointerType), uint(1)).Return([]*model.UserRole{}, nil)

		res, _, err := s.RegisterAsSeller(req, uint(1))

		assert.Nil(t, res)
		assert.NotNil(t, err)
	})

	t.Run("Should return error", func(t *testing.T) {

		gormDB := testutil.MockDB()
		mockRepo1 := new(mocks.UserRepository)
		mockRepo2 := new(mocks.UserRoleRepository)
		mockRepo3 := new(mocks.AddressRepository)
		mockRepo4 := new(mocks.WalletRepository)
		cfg := &service.UserServiceConfig{
			DB:               gormDB,
			UserRepository:   mockRepo1,
			UserRoleRepo:     mockRepo2,
			AddressRepo:      mockRepo3,
			WalletRepository: mockRepo4,
			AppConfig:        config.AppConfig{},
		}
		s := service.NewUserService(cfg)
		req := &dto.RegisterAsSellerReq{}

		mockRepo1.On("GetUserByID", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("uint")).Return(&model.User{ID: 1}, nil)

		mockRepo3.On("GetUserMainAddress", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("uint")).Return(nil, errors.New(""))

		mockRepo2.On("CreateRoleToUser", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("*model.UserRole")).Return(nil, nil)

		mockRepo1.On("RegisterAsSeller", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("*model.Seller")).Return(&model.Seller{}, nil)

		mockRepo4.On("GetWalletByUserID", mock.AnythingOfType(testutil.GormDBPointerType), uint(1)).Return(&model.Wallet{}, nil)

		mockRepo2.On("GetRolesByUserID", mock.AnythingOfType(testutil.GormDBPointerType), uint(1)).Return([]*model.UserRole{}, nil)

		res, _, err := s.RegisterAsSeller(req, uint(1))

		assert.Nil(t, res)
		assert.NotNil(t, err)
	})

	t.Run("Should return error", func(t *testing.T) {

		gormDB := testutil.MockDB()
		mockRepo1 := new(mocks.UserRepository)
		mockRepo2 := new(mocks.UserRoleRepository)
		mockRepo3 := new(mocks.AddressRepository)
		mockRepo4 := new(mocks.WalletRepository)
		cfg := &service.UserServiceConfig{
			DB:               gormDB,
			UserRepository:   mockRepo1,
			UserRoleRepo:     mockRepo2,
			AddressRepo:      mockRepo3,
			WalletRepository: mockRepo4,
			AppConfig:        config.AppConfig{},
		}
		s := service.NewUserService(cfg)
		req := &dto.RegisterAsSellerReq{}

		mockRepo1.On("GetUserByID", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("uint")).Return(&model.User{ID: 1}, nil)

		mockRepo3.On("GetUserMainAddress", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("uint")).Return(&model.Address{}, nil)

		mockRepo2.On("CreateRoleToUser", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("*model.UserRole")).Return(nil, errors.New(""))

		mockRepo1.On("RegisterAsSeller", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("*model.Seller")).Return(&model.Seller{}, nil)

		mockRepo4.On("GetWalletByUserID", mock.AnythingOfType(testutil.GormDBPointerType), uint(1)).Return(&model.Wallet{}, nil)

		mockRepo2.On("GetRolesByUserID", mock.AnythingOfType(testutil.GormDBPointerType), uint(1)).Return([]*model.UserRole{}, nil)

		res, _, err := s.RegisterAsSeller(req, uint(1))

		assert.Nil(t, res)
		assert.NotNil(t, err)
	})

	t.Run("Should return error", func(t *testing.T) {

		gormDB := testutil.MockDB()
		mockRepo1 := new(mocks.UserRepository)
		mockRepo2 := new(mocks.UserRoleRepository)
		mockRepo3 := new(mocks.AddressRepository)
		mockRepo4 := new(mocks.WalletRepository)
		cfg := &service.UserServiceConfig{
			DB:               gormDB,
			UserRepository:   mockRepo1,
			UserRoleRepo:     mockRepo2,
			AddressRepo:      mockRepo3,
			WalletRepository: mockRepo4,
			AppConfig:        config.AppConfig{},
		}
		s := service.NewUserService(cfg)
		req := &dto.RegisterAsSellerReq{}

		mockRepo1.On("GetUserByID", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("uint")).Return(&model.User{ID: 1}, nil)

		mockRepo3.On("GetUserMainAddress", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("uint")).Return(&model.Address{}, nil)

		mockRepo2.On("CreateRoleToUser", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("*model.UserRole")).Return(nil, nil)

		mockRepo1.On("RegisterAsSeller", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("*model.Seller")).Return(nil, errors.New(""))

		mockRepo4.On("GetWalletByUserID", mock.AnythingOfType(testutil.GormDBPointerType), uint(1)).Return(&model.Wallet{}, nil)

		mockRepo2.On("GetRolesByUserID", mock.AnythingOfType(testutil.GormDBPointerType), uint(1)).Return([]*model.UserRole{}, nil)

		res, _, err := s.RegisterAsSeller(req, uint(1))

		assert.Nil(t, res)
		assert.NotNil(t, err)
	})

	t.Run("Should return error", func(t *testing.T) {

		gormDB := testutil.MockDB()
		mockRepo1 := new(mocks.UserRepository)
		mockRepo2 := new(mocks.UserRoleRepository)
		mockRepo3 := new(mocks.AddressRepository)
		mockRepo4 := new(mocks.WalletRepository)
		cfg := &service.UserServiceConfig{
			DB:               gormDB,
			UserRepository:   mockRepo1,
			UserRoleRepo:     mockRepo2,
			AddressRepo:      mockRepo3,
			WalletRepository: mockRepo4,
			AppConfig:        config.AppConfig{},
		}
		s := service.NewUserService(cfg)
		req := &dto.RegisterAsSellerReq{}

		mockRepo1.On("GetUserByID", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("uint")).Return(&model.User{ID: 1}, nil)

		mockRepo3.On("GetUserMainAddress", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("uint")).Return(&model.Address{}, nil)

		mockRepo2.On("CreateRoleToUser", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("*model.UserRole")).Return(nil, nil)

		mockRepo1.On("RegisterAsSeller", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("*model.Seller")).Return(&model.Seller{}, nil)

		mockRepo4.On("GetWalletByUserID", mock.AnythingOfType(testutil.GormDBPointerType), uint(1)).Return(nil, errors.New(""))

		mockRepo2.On("GetRolesByUserID", mock.AnythingOfType(testutil.GormDBPointerType), uint(1)).Return([]*model.UserRole{}, nil)

		res, _, err := s.RegisterAsSeller(req, uint(1))

		assert.Nil(t, res)
		assert.NotNil(t, err)
	})

	t.Run("Should return error", func(t *testing.T) {

		gormDB := testutil.MockDB()
		mockRepo1 := new(mocks.UserRepository)
		mockRepo2 := new(mocks.UserRoleRepository)
		mockRepo3 := new(mocks.AddressRepository)
		mockRepo4 := new(mocks.WalletRepository)
		cfg := &service.UserServiceConfig{
			DB:               gormDB,
			UserRepository:   mockRepo1,
			UserRoleRepo:     mockRepo2,
			AddressRepo:      mockRepo3,
			WalletRepository: mockRepo4,
			AppConfig:        config.AppConfig{},
		}
		s := service.NewUserService(cfg)
		req := &dto.RegisterAsSellerReq{}

		mockRepo1.On("GetUserByID", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("uint")).Return(&model.User{ID: 1}, nil)

		mockRepo3.On("GetUserMainAddress", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("uint")).Return(&model.Address{}, nil)

		mockRepo2.On("CreateRoleToUser", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("*model.UserRole")).Return(nil, nil)

		mockRepo1.On("RegisterAsSeller", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("*model.Seller")).Return(&model.Seller{}, nil)

		mockRepo4.On("GetWalletByUserID", mock.AnythingOfType(testutil.GormDBPointerType), uint(1)).Return(&model.Wallet{}, nil)

		mockRepo2.On("GetRolesByUserID", mock.AnythingOfType(testutil.GormDBPointerType), uint(1)).Return(nil, errors.New(""))

		res, _, err := s.RegisterAsSeller(req, uint(1))

		assert.Nil(t, res)
		assert.NotNil(t, err)
	})
}

func TestUserService_UserDetails(t *testing.T) {
	t.Run("Should return response body", func(t *testing.T) {

		gormDB := testutil.MockDB()
		mockRepo1 := new(mocks.UserRepository)
		mockRepo2 := new(mocks.UserRoleRepository)
		mockRepo3 := new(mocks.AddressRepository)
		mockRepo4 := new(mocks.WalletRepository)
		cfg := &service.UserServiceConfig{
			DB:               gormDB,
			UserRepository:   mockRepo1,
			UserRoleRepo:     mockRepo2,
			AddressRepo:      mockRepo3,
			WalletRepository: mockRepo4,
			AppConfig:        config.AppConfig{},
		}
		s := service.NewUserService(cfg)
		expectedRes := &dto.UserDetailsRes{}

		mockRepo1.On("GetUserDetailsByID", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("uint")).Return(&model.User{ID: 1}, nil)

		res, err := s.UserDetails(uint(1))

		assert.Nil(t, err)
		assert.Equal(t, expectedRes, res)

	})

	t.Run("Should return error", func(t *testing.T) {

		gormDB := testutil.MockDB()
		mockRepo1 := new(mocks.UserRepository)
		mockRepo2 := new(mocks.UserRoleRepository)
		mockRepo3 := new(mocks.AddressRepository)
		mockRepo4 := new(mocks.WalletRepository)
		cfg := &service.UserServiceConfig{
			DB:               gormDB,
			UserRepository:   mockRepo1,
			UserRoleRepo:     mockRepo2,
			AddressRepo:      mockRepo3,
			WalletRepository: mockRepo4,
			AppConfig:        config.AppConfig{},
		}
		s := service.NewUserService(cfg)

		mockRepo1.On("GetUserDetailsByID", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("uint")).Return(nil, errors.New(""))

		res, err := s.UserDetails(uint(1))

		assert.Nil(t, res)
		assert.NotNil(t, err)

	})

}
func TestUserService_ChangeUserDetailsLessPassword(t *testing.T) {
	t.Run("Should return response body", func(t *testing.T) {

		gormDB := testutil.MockDB()
		mockRepo1 := new(mocks.UserRepository)
		mockRepo2 := new(mocks.UserRoleRepository)
		mockRepo3 := new(mocks.AddressRepository)
		mockRepo4 := new(mocks.WalletRepository)
		cfg := &service.UserServiceConfig{
			DB:               gormDB,
			UserRepository:   mockRepo1,
			UserRoleRepo:     mockRepo2,
			AddressRepo:      mockRepo3,
			WalletRepository: mockRepo4,
			AppConfig:        config.AppConfig{},
		}
		s := service.NewUserService(cfg)
		req := &dto.ChangeUserDetails{}
		expectedRes := &model.User{}

		mockRepo1.On("ChangeUserDetailsLessPassword", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("uint"), mock.AnythingOfType("*dto.ChangeUserDetails")).Return(&model.User{}, nil)

		res, err := s.ChangeUserDetailsLessPassword(uint(1), req)

		assert.Nil(t, err)
		assert.Equal(t, expectedRes, res)

	})

	t.Run("Should return error", func(t *testing.T) {

		gormDB := testutil.MockDB()
		mockRepo1 := new(mocks.UserRepository)
		mockRepo2 := new(mocks.UserRoleRepository)
		mockRepo3 := new(mocks.AddressRepository)
		mockRepo4 := new(mocks.WalletRepository)
		cfg := &service.UserServiceConfig{
			DB:               gormDB,
			UserRepository:   mockRepo1,
			UserRoleRepo:     mockRepo2,
			AddressRepo:      mockRepo3,
			WalletRepository: mockRepo4,
			AppConfig:        config.AppConfig{},
		}
		s := service.NewUserService(cfg)
		req := &dto.ChangeUserDetails{}

		mockRepo1.On("ChangeUserDetailsLessPassword", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("uint"), mock.AnythingOfType("*dto.ChangeUserDetails")).Return(nil, errors.New("s"))

		res, err := s.ChangeUserDetailsLessPassword(uint(1), req)

		assert.Nil(t, res)
		assert.NotNil(t, err)

	})
}

func TestUserService_ChangeUserPassword(t *testing.T) {
	t.Run("Should return response body", func(t *testing.T) {

		gormDB := testutil.MockDB()
		mockRepo1 := new(mocks.UserRepository)
		mockRepo2 := new(mocks.UserRoleRepository)
		mockRepo3 := new(mocks.AddressRepository)
		mockRepo4 := new(mocks.WalletRepository)
		cfg := &service.UserServiceConfig{
			DB:               gormDB,
			UserRepository:   mockRepo1,
			UserRoleRepo:     mockRepo2,
			AddressRepo:      mockRepo3,
			WalletRepository: mockRepo4,
			AppConfig:        config.AppConfig{},
		}
		s := service.NewUserService(cfg)
		req := &dto.ChangePasswordReq{NewPassword: "456", RepeatNewPassword: "456"}

		mockRepo1.On("MatchingCredential", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(&model.User{Username: "asdf"}, nil)

		mockRepo1.On("GetUserDetailsByID", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("uint")).Return(&model.User{Username: "asdf"}, nil)

		mockRepo1.On("ChangeUserPassword", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("uint"), mock.AnythingOfType("string")).Return(nil)

		err := s.ChangeUserPassword(uint(1), req)

		assert.Nil(t, err)

	})
	t.Run("Should return response body", func(t *testing.T) {

		gormDB := testutil.MockDB()
		mockRepo1 := new(mocks.UserRepository)
		mockRepo2 := new(mocks.UserRoleRepository)
		mockRepo3 := new(mocks.AddressRepository)
		mockRepo4 := new(mocks.WalletRepository)
		cfg := &service.UserServiceConfig{
			DB:               gormDB,
			UserRepository:   mockRepo1,
			UserRoleRepo:     mockRepo2,
			AddressRepo:      mockRepo3,
			WalletRepository: mockRepo4,
			AppConfig:        config.AppConfig{},
		}
		s := service.NewUserService(cfg)
		req := &dto.ChangePasswordReq{NewPassword: "456", RepeatNewPassword: "456"}

		mockRepo1.On("MatchingCredential", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(nil, errors.New(""))

		mockRepo1.On("GetUserDetailsByID", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("uint")).Return(&model.User{Username: "asdf"}, nil)

		mockRepo1.On("ChangeUserPassword", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("uint"), mock.AnythingOfType("string")).Return(nil)

		err := s.ChangeUserPassword(uint(1), req)

		assert.NotNil(t, err)

	})

	t.Run("Should return error", func(t *testing.T) {

		gormDB := testutil.MockDB()
		mockRepo1 := new(mocks.UserRepository)
		mockRepo2 := new(mocks.UserRoleRepository)
		mockRepo3 := new(mocks.AddressRepository)
		mockRepo4 := new(mocks.WalletRepository)
		cfg := &service.UserServiceConfig{
			DB:               gormDB,
			UserRepository:   mockRepo1,
			UserRoleRepo:     mockRepo2,
			AddressRepo:      mockRepo3,
			WalletRepository: mockRepo4,
			AppConfig:        config.AppConfig{},
		}
		s := service.NewUserService(cfg)
		req := &dto.ChangePasswordReq{NewPassword: "456", RepeatNewPassword: "123"}

		mockRepo1.On("MatchingCredential", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(&model.User{Username: "asdf"}, nil)

		mockRepo1.On("GetUserDetailsByID", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("uint")).Return(&model.User{Username: "asdf"}, nil)

		mockRepo1.On("ChangeUserPassword", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("uint"), mock.AnythingOfType("string")).Return(nil)

		err := s.ChangeUserPassword(uint(1), req)

		assert.NotNil(t, err)
	})

	t.Run("Should return error", func(t *testing.T) {

		gormDB := testutil.MockDB()
		mockRepo1 := new(mocks.UserRepository)
		mockRepo2 := new(mocks.UserRoleRepository)
		mockRepo3 := new(mocks.AddressRepository)
		mockRepo4 := new(mocks.WalletRepository)
		cfg := &service.UserServiceConfig{
			DB:               gormDB,
			UserRepository:   mockRepo1,
			UserRoleRepo:     mockRepo2,
			AddressRepo:      mockRepo3,
			WalletRepository: mockRepo4,
			AppConfig:        config.AppConfig{},
		}
		s := service.NewUserService(cfg)
		req := &dto.ChangePasswordReq{NewPassword: "asdf", RepeatNewPassword: "asdf"}

		mockRepo1.On("MatchingCredential", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(&model.User{Username: "asdf"}, nil)

		mockRepo1.On("GetUserDetailsByID", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("uint")).Return(&model.User{Username: "asdf"}, nil)

		mockRepo1.On("ChangeUserPassword", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("uint"), mock.AnythingOfType("string")).Return(nil)

		err := s.ChangeUserPassword(uint(1), req)

		assert.NotNil(t, err)
	})

	t.Run("Should return error", func(t *testing.T) {

		gormDB := testutil.MockDB()
		mockRepo1 := new(mocks.UserRepository)
		mockRepo2 := new(mocks.UserRoleRepository)
		mockRepo3 := new(mocks.AddressRepository)
		mockRepo4 := new(mocks.WalletRepository)
		cfg := &service.UserServiceConfig{
			DB:               gormDB,
			UserRepository:   mockRepo1,
			UserRoleRepo:     mockRepo2,
			AddressRepo:      mockRepo3,
			WalletRepository: mockRepo4,
			AppConfig:        config.AppConfig{},
		}
		s := service.NewUserService(cfg)
		req := &dto.ChangePasswordReq{NewPassword: "456", RepeatNewPassword: "456"}

		mockRepo1.On("MatchingCredential", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(&model.User{Username: "asdf"}, nil)

		mockRepo1.On("GetUserDetailsByID", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("uint")).Return(&model.User{Username: "asdf"}, nil)

		mockRepo1.On("ChangeUserPassword", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("uint"), mock.AnythingOfType("string")).Return(errors.New(""))

		err := s.ChangeUserPassword(uint(1), req)

		assert.NotNil(t, err)
	})
}
