package service

import (
	"gorm.io/gorm"
	"os"
	"seadeals-backend/config"
	"seadeals-backend/dto"
	"seadeals-backend/helper"
	"seadeals-backend/model"
	"seadeals-backend/repository"
	"strings"
)

type AuthService interface {
	AuthAfterRegister(*model.User, *model.Wallet, *gorm.DB) (string, string, error)
	SignInWithGoogle(*model.User) (string, string, error)
	SignIn(*dto.SignInReq) (string, string, error)
	SignOut(uint) error
	StepUpPassword(userID uint, req *dto.StepUpPasswordRes) error
}

type authService struct {
	db               *gorm.DB
	refreshTokenRepo repository.RefreshTokenRepository
	userRepository   repository.UserRepository
	userRoleRepo     repository.UserRoleRepository
	walletRepository repository.WalletRepository
	appConfig        config.AppConfig
}

type AuthSConfig struct {
	DB               *gorm.DB
	RefreshTokenRepo repository.RefreshTokenRepository
	UserRepository   repository.UserRepository
	UserRoleRepo     repository.UserRoleRepository
	WalletRepository repository.WalletRepository
	AppConfig        config.AppConfig
}

func NewAuthService(config *AuthSConfig) AuthService {
	return &authService{
		db:               config.DB,
		refreshTokenRepo: config.RefreshTokenRepo,
		userRepository:   config.UserRepository,
		userRoleRepo:     config.UserRoleRepo,
		walletRepository: config.WalletRepository,
		appConfig:        config.AppConfig,
	}
}

func (a *authService) AuthAfterRegister(user *model.User, wallet *model.Wallet, tx *gorm.DB) (string, string, error) {
	userRoles, err := a.userRoleRepo.GetRolesByUserID(tx, user.ID)
	if err != nil {
		return "", "", err
	}
	var roles []string
	for _, role := range userRoles {
		roles = append(roles, role.Role.Name)
	}
	rolesString := strings.Join(roles[:], " ")

	userJWT := &dto.UserJWT{
		Name:     user.FullName,
		UserID:   user.ID,
		Email:    user.Email,
		Username: user.Username,
		WalletID: wallet.ID,
	}
	token, err := helper.GenerateJWTToken(userJWT, rolesString, config.Config.JWTExpiredInMinuteTime*60, dto.JWTAccessToken)
	refreshToken, err := helper.GenerateJWTToken(userJWT, rolesString, 24*60*60, dto.JWTRefreshToken)
	if os.Getenv("ENV") == "testing" {
		token = "test"
		refreshToken = "test"
	}

	err = a.refreshTokenRepo.CreateRefreshToken(tx, user.ID, refreshToken)
	if err != nil {
		return "", "", err
	}

	tx.Commit()
	return token, refreshToken, err
}

func (a *authService) SignInWithGoogle(user *model.User) (string, string, error) {
	tx := a.db.Begin()
	var err error
	defer helper.CommitOrRollback(tx, &err)

	wallet, err := a.walletRepository.GetWalletByUserID(tx, user.ID)
	if err != nil {
		return "", "", err
	}
	userJWT := &dto.UserJWT{
		Name:      user.FullName,
		UserID:    user.ID,
		Email:     user.Email,
		Username:  user.Username,
		WalletID:  wallet.ID,
		AvatarURL: user.AvatarURL,
	}

	userRoles, err := a.userRoleRepo.GetRolesByUserID(tx, user.ID)
	if err != nil {
		return "", "", err
	}
	var roles []string
	for _, role := range userRoles {
		roles = append(roles, role.Role.Name)
	}
	rolesString := strings.Join(roles[:], " ")

	token, err := helper.GenerateJWTToken(userJWT, rolesString, config.Config.JWTExpiredInMinuteTime*60, dto.JWTAccessToken)
	refreshToken, err := helper.GenerateJWTToken(userJWT, rolesString, 24*60*60, dto.JWTRefreshToken)
	if os.Getenv("ENV") == "testing" {
		token = "test"
		refreshToken = "test"
	}
	err = a.refreshTokenRepo.CreateRefreshToken(tx, user.ID, refreshToken)
	if err != nil {
		return "", "", err
	}

	return token, refreshToken, err
}

func (a *authService) SignIn(req *dto.SignInReq) (string, string, error) {
	tx := a.db.Begin()
	var err error
	defer helper.CommitOrRollback(tx, &err)

	user, err := a.userRepository.MatchingCredential(tx, req.Email, req.Password)
	if err != nil || user == nil {
		return "", "", err
	}
	wallet, err := a.walletRepository.GetWalletByUserID(tx, user.ID)
	if err != nil {
		return "", "", err
	}

	userJWT := &dto.UserJWT{
		Name:      user.FullName,
		UserID:    user.ID,
		Email:     user.Email,
		Username:  user.Username,
		WalletID:  wallet.ID,
		AvatarURL: user.AvatarURL,
	}

	userRoles, err := a.userRoleRepo.GetRolesByUserID(tx, user.ID)
	if err != nil {
		return "", "", err
	}
	var roles []string
	for _, role := range userRoles {
		roles = append(roles, role.Role.Name)
	}
	rolesString := strings.Join(roles[:], " ")
	token, err := helper.GenerateJWTToken(userJWT, rolesString, config.Config.JWTExpiredInMinuteTime*60, dto.JWTAccessToken)
	refreshToken, err := helper.GenerateJWTToken(userJWT, rolesString, 24*60*60, dto.JWTRefreshToken)
	if os.Getenv("ENV") == "testing" {
		token = "test"
		refreshToken = "test"
	}
	err = a.refreshTokenRepo.CreateRefreshToken(tx, user.ID, refreshToken)
	if err != nil {
		return "", "", err
	}

	return token, refreshToken, err
}

func (a *authService) SignOut(userID uint) error {
	tx := a.db.Begin()
	var err error
	defer helper.CommitOrRollback(tx, &err)

	err = a.refreshTokenRepo.DeleteRefreshToken(tx, userID)
	if err != nil {
		return err
	}
	return nil
}
func (a *authService) StepUpPassword(userID uint, req *dto.StepUpPasswordRes) error {
	tx := a.db.Begin()
	var err error
	defer helper.CommitOrRollback(tx, &err)

	err = a.walletRepository.StepUpPassword(tx, userID, req.Password)
	if err != nil {
		return err
	}

	return nil
}
