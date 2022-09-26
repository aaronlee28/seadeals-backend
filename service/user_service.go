package service

import (
	"gorm.io/gorm"
	"net/mail"
	"regexp"
	"seadeals-backend/apperror"
	"seadeals-backend/dto"
	"seadeals-backend/model"
	"seadeals-backend/repository"
	"time"
)

type UserService interface {
	Register(req *dto.RegisterRequest) (*dto.RegisterResponse, *gorm.DB, error)
	CheckGoogleAccount(req *dto.GoogleLogin) (*model.User, error)
	RegisterAsSeller(req *dto.RegisterAsSellerReq) (*model.Seller, error)
}

type userService struct {
	db               *gorm.DB
	userRepository   repository.UserRepository
	userRoleRepo     repository.UserRoleRepository
	walletRepository repository.WalletRepository
}

type UserServiceConfig struct {
	DB               *gorm.DB
	UserRepository   repository.UserRepository
	UserRoleRepo     repository.UserRoleRepository
	WalletRepository repository.WalletRepository
}

func NewUserService(c *UserServiceConfig) UserService {
	return &userService{
		db:               c.DB,
		userRepository:   c.UserRepository,
		userRoleRepo:     c.UserRoleRepo,
		walletRepository: c.WalletRepository,
	}
}

func (u *userService) Register(req *dto.RegisterRequest) (*dto.RegisterResponse, *gorm.DB, error) {
	_, err := mail.ParseAddress(req.Email)
	if err != nil {
		return nil, nil, apperror.BadRequestError("Email is not valid")
	}

	isMatch, _ := regexp.MatchString(req.Username, req.Password)
	if isMatch {
		return nil, nil, apperror.BadRequestError("Password cannot contain username")
	}

	tx := u.db.Begin()
	birthDate, _ := time.Parse("2006-01-02", req.BirthDate)
	newUser := &model.User{
		FullName:  req.FullName,
		Username:  req.Username,
		Email:     req.Email,
		Phone:     req.Phone,
		Password:  req.Password,
		Gender:    req.Gender,
		BirthDate: birthDate,
	}
	user, err := u.userRepository.Register(tx, newUser)
	if err != nil {
		tx.Rollback()
		return nil, nil, err
	}

	newWallet := &model.Wallet{
		UserID:       user.ID,
		Balance:      0,
		Pin:          nil,
		Status:       model.WalletActive,
		BlockedUntil: nil,
	}
	wallet, err := u.walletRepository.CreateWallet(tx, newWallet)
	if err != nil {
		tx.Rollback()
		return nil, nil, err
	}

	newUserRole := &model.UserRole{
		UserID: user.ID,
		RoleID: 1,
	}
	_, err = u.userRoleRepo.CreateRoleToUser(tx, newUserRole)
	if err != nil {
		tx.Rollback()
		return nil, nil, err
	}

	userResponse := &dto.RegisterResponse{
		ID:       user.ID,
		FullName: user.FullName,
		Username: user.Username,
		Email:    user.Email,
		Phone:    user.Phone,
		Role:     model.UserRoleName,
		Wallet: model.Wallet{
			ID:      wallet.ID,
			Balance: wallet.Balance,
		},
	}

	return userResponse, tx, nil
}

func (u *userService) CheckGoogleAccount(req *dto.GoogleLogin) (*model.User, error) {
	_, err := mail.ParseAddress(req.Email)
	if err != nil {
		return nil, apperror.BadRequestError("Email is not valid")
	}

	tx := u.db.Begin()
	isEmailExist, err := u.userRepository.HasExistEmail(tx, req.Email)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	if !isEmailExist {
		return nil, apperror.NotFoundError("email doesn't exist")
	}

	user, err := u.userRepository.GetUserByEmail(tx, req.Email)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	tx.Commit()
	return user, nil
}

func (u *userService) RegisterAsSeller(req *dto.RegisterAsSellerReq) (*model.Seller, error) {
	tx := u.db.Begin()

	user, err := u.userRepository.GetUserByID(tx, req.UserID)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	address, err := u.userRepository.GetUserMainAddress(tx, user.ID)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	newSeller := &model.Seller{
		Name:        req.ShopName,
		Slug:        "",
		UserID:      user.ID,
		Description: req.Description,
		AddressID:   address.ID,
		PictureURL:  "",
		BannerURL:   "",
	}

	createdSeller, err := u.userRepository.RegisterAsSeller(tx, newSeller)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	tx.Commit()
	return createdSeller, nil
}
