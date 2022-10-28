package service

import (
	"gorm.io/gorm"
	"net/mail"
	"regexp"
	"seadeals-backend/apperror"
	"seadeals-backend/config"
	"seadeals-backend/dto"
	"seadeals-backend/helper"
	"seadeals-backend/model"
	"seadeals-backend/repository"
	"strings"
	"time"
)

type UserService interface {
	Register(req *dto.RegisterRequest) (*dto.RegisterResponse, *gorm.DB, error)
	CheckGoogleAccount(email string) (*model.User, error)
	RegisterAsSeller(req *dto.RegisterAsSellerReq, userID uint) (*model.Seller, string, error)
	UserDetails(userID uint) (*dto.UserDetailsRes, error)
	ChangeUserDetailsLessPassword(userID uint, req *dto.ChangeUserDetails) (*model.User, error)
	ChangeUserPassword(userID uint, req *dto.ChangePasswordReq) error
}

type userService struct {
	db               *gorm.DB
	userRepository   repository.UserRepository
	userRoleRepo     repository.UserRoleRepository
	addressRepo      repository.AddressRepository
	walletRepository repository.WalletRepository
	appConfig        config.AppConfig
}

type UserServiceConfig struct {
	DB               *gorm.DB
	UserRepository   repository.UserRepository
	UserRoleRepo     repository.UserRoleRepository
	AddressRepo      repository.AddressRepository
	WalletRepository repository.WalletRepository
	AppConfig        config.AppConfig
}

func NewUserService(c *UserServiceConfig) UserService {
	return &userService{
		db:               c.DB,
		userRepository:   c.UserRepository,
		userRoleRepo:     c.UserRoleRepo,
		addressRepo:      c.AddressRepo,
		walletRepository: c.WalletRepository,
		appConfig:        c.AppConfig,
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

func (u *userService) CheckGoogleAccount(email string) (*model.User, error) {
	tx := u.db.Begin()
	var err error
	defer helper.CommitOrRollback(tx, &err)

	_, err = mail.ParseAddress(email)
	if err != nil {
		err = apperror.BadRequestError("Email is not valid")
		return nil, err
	}

	isEmailExist, err := u.userRepository.HasExistEmail(tx, email)
	if err != nil {
		return nil, err
	}

	if !isEmailExist {
		err = apperror.NotFoundError("email doesn't exist")
		return nil, err
	}

	user, err := u.userRepository.GetUserByEmail(tx, email)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (u *userService) RegisterAsSeller(req *dto.RegisterAsSellerReq, userID uint) (*model.Seller, string, error) {
	tx := u.db.Begin()
	var err error
	defer helper.CommitOrRollback(tx, &err)

	user, err := u.userRepository.GetUserByID(tx, userID)
	if err != nil {
		return nil, "", err
	}

	address, err := u.addressRepo.GetUserMainAddress(tx, user.ID)
	if err != nil {
		return nil, "", err
	}

	newUserRole := &model.UserRole{
		UserID: user.ID,
		RoleID: 3,
	}
	_, err = u.userRoleRepo.CreateRoleToUser(tx, newUserRole)
	if err != nil {
		return nil, "", err
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
		return nil, "", err
	}

	wallet, err := u.walletRepository.GetWalletByUserID(tx, user.ID)
	if err != nil {
		return nil, "", err
	}

	userJWT := &dto.UserJWT{
		UserID:   user.ID,
		Email:    user.Email,
		Username: user.Username,
		WalletID: wallet.ID,
	}

	userRoles, err := u.userRoleRepo.GetRolesByUserID(tx, user.ID)
	if err != nil {
		return nil, "", err
	}
	var roles []string
	for _, role := range userRoles {
		roles = append(roles, role.Role.Name)
	}
	rolesString := strings.Join(roles[:], " ")
	accessToken, err := helper.GenerateJWTToken(userJWT, rolesString, config.Config.JWTExpiredInMinuteTime*60, dto.JWTAccessToken)

	return createdSeller, accessToken, nil
}

func (u *userService) UserDetails(userID uint) (*dto.UserDetailsRes, error) {
	tx := u.db.Begin()
	var err error
	defer helper.CommitOrRollback(tx, &err)

	var userDetails *model.User
	userDetails, err = u.userRepository.GetUserDetailsByID(tx, userID)

	if err != nil {
		return nil, err
	}
	res := &dto.UserDetailsRes{
		Username:  userDetails.Username,
		FullName:  userDetails.FullName,
		Email:     userDetails.Email,
		Phone:     userDetails.Phone,
		Gender:    userDetails.Gender,
		BirthDate: userDetails.BirthDate,
	}
	if userDetails.AvatarURL != nil {
		res.AvatarURL = *userDetails.AvatarURL
	}
	return res, nil
}

func (u *userService) ChangeUserDetailsLessPassword(userID uint, req *dto.ChangeUserDetails) (*model.User, error) {
	tx := u.db.Begin()
	var err error
	defer helper.CommitOrRollback(tx, &err)

	var userDetails *model.User
	userDetails, err = u.userRepository.ChangeUserDetailsLessPassword(tx, userID, req)

	if err != nil {
		return nil, err
	}

	return userDetails, nil
}

func (u *userService) ChangeUserPassword(userID uint, req *dto.ChangePasswordReq) error {
	tx := u.db.Begin()
	var err error
	defer helper.CommitOrRollback(tx, &err)

	var userDetails *model.User
	userDetails, err = u.userRepository.GetUserDetailsByID(tx, userID)

	user, err := u.userRepository.MatchingCredential(tx, req.Email, req.CurrentPassword)
	if err != nil || user == nil {
		return apperror.BadRequestError("incorrect current password")
	}
	if req.CurrentPassword == req.NewPassword {
		return apperror.InternalServerError("new password must be different")
	}

	isMatch, _ := regexp.MatchString(req.NewPassword, req.RepeatNewPassword)
	if !isMatch {
		return apperror.BadRequestError("Password does not match")
	}

	isMatchUsername, _ := regexp.MatchString(userDetails.Username, req.NewPassword)
	if isMatchUsername {
		return apperror.BadRequestError("Password cannot contain username")
	}

	err = u.userRepository.ChangeUserPassword(tx, userID, req.NewPassword)

	if err != nil {
		return err
	}

	return nil
}
