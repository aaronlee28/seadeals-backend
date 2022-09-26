package service

import (
	"github.com/golang-jwt/jwt/v4"
	"gorm.io/gorm"
	"net/mail"
	"regexp"
	"seadeals-backend/apperror"
	"seadeals-backend/config"
	"seadeals-backend/dto"
	"seadeals-backend/model"
	"seadeals-backend/repository"
	"strings"
	"time"
)

type UserService interface {
	Register(req *dto.RegisterRequest) (*dto.RegisterResponse, *gorm.DB, error)
	CheckGoogleAccount(email string) (*model.User, error)
	RegisterAsSeller(req *dto.RegisterAsSellerReq) (*model.Seller, string, error)
}

type userService struct {
	db               *gorm.DB
	userRepository   repository.UserRepository
	userRoleRepo     repository.UserRoleRepository
	walletRepository repository.WalletRepository
	appConfig        config.AppConfig
}

type UserServiceConfig struct {
	DB               *gorm.DB
	UserRepository   repository.UserRepository
	UserRoleRepo     repository.UserRoleRepository
	WalletRepository repository.WalletRepository
	AppConfig        config.AppConfig
}

func NewUserService(c *UserServiceConfig) UserService {
	return &userService{
		db:               c.DB,
		userRepository:   c.UserRepository,
		userRoleRepo:     c.UserRoleRepo,
		walletRepository: c.WalletRepository,
		appConfig:        c.AppConfig,
	}
}

func (u *userService) generateJWTToken(user *dto.UserJWT, role string, idExp int64, jwtType string) (string, error) {
	// 1 minutes times JWTExpireInMinutes
	unixTime := time.Now().Unix()
	tokenExp := unixTime + idExp

	timeExpire := jwt.NumericDate{Time: time.Unix(tokenExp, 0)}
	timeNow := jwt.NumericDate{Time: time.Now()}
	accessClaims := &idTokenClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: &timeExpire,
			IssuedAt:  &timeNow,
			Issuer:    u.appConfig.AppName,
		},
		User:  user,
		Scope: role,
		Type:  jwtType,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	tokenString, _ := token.SignedString(u.appConfig.JWTSecret)

	return tokenString, nil
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
	_, err := mail.ParseAddress(email)
	if err != nil {
		return nil, apperror.BadRequestError("Email is not valid")
	}

	tx := u.db.Begin()
	isEmailExist, err := u.userRepository.HasExistEmail(tx, email)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	if !isEmailExist {
		return nil, apperror.NotFoundError("email doesn't exist")
	}

	user, err := u.userRepository.GetUserByEmail(tx, email)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	tx.Commit()
	return user, nil
}

func (u *userService) RegisterAsSeller(req *dto.RegisterAsSellerReq) (*model.Seller, string, error) {
	tx := u.db.Begin()

	user, err := u.userRepository.GetUserByID(tx, req.UserID)
	if err != nil {
		tx.Rollback()
		return nil, "", err
	}

	address, err := u.userRepository.GetUserMainAddress(tx, user.ID)
	if err != nil {
		tx.Rollback()
		return nil, "", err
	}

	newUserRole := &model.UserRole{
		UserID: user.ID,
		RoleID: 3,
	}
	_, err = u.userRoleRepo.CreateRoleToUser(tx, newUserRole)
	if err != nil {
		tx.Rollback()
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
		tx.Rollback()
		return nil, "", err
	}

	wallet, err := u.walletRepository.GetWalletByUserID(tx, user.ID)
	if err != nil {
		tx.Rollback()
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
		tx.Rollback()
		return nil, "", err
	}
	var roles []string
	for _, role := range userRoles {
		roles = append(roles, role.Role.Name)
	}
	rolesString := strings.Join(roles[:], " ")
	accessToken, err := u.generateJWTToken(userJWT, rolesString, config.Config.JWTExpiredInMinuteTime*60, dto.JWTAccessToken)

	tx.Commit()
	return createdSeller, accessToken, nil
}
