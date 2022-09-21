package service

import (
	"github.com/golang-jwt/jwt/v4"
	"gorm.io/gorm"
	"os"
	"seadeals-backend/apperror"
	"seadeals-backend/config"
	"seadeals-backend/dto"
	"seadeals-backend/model"
	"seadeals-backend/repository"
	"strings"
	"time"
)

type AuthService interface {
	AuthAfterRegister(*model.User, *model.Wallet, *gorm.DB) (string, string, error)
	SignInWithGoogle(*model.User) (string, string, error)
	SignIn(*dto.SignInReq) (string, string, error)
	SignOut(uint) error
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

type idTokenClaims struct {
	jwt.RegisteredClaims
	User  *dto.UserJWT `json:"user"`
	Scope string       `json:"scope"`
}

func (a *authService) generateJWTToken(user *dto.UserJWT, role string) (string, string, error) {
	// 1 minutes times JWTExpireInMinutes
	var idExp = a.appConfig.JWTExpiredInMinuteTime * 60
	unixTime := time.Now().Unix()
	tokenExp := unixTime + idExp

	timeExpire := jwt.NumericDate{Time: time.Unix(tokenExp, 0)}
	timeNow := jwt.NumericDate{Time: time.Now()}
	accessClaims := &idTokenClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: &timeExpire,
			IssuedAt:  &timeNow,
			Issuer:    a.appConfig.AppName,
		},
		User:  user,
		Scope: role,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	tokenString, _ := token.SignedString(a.appConfig.JWTSecret)

	// one day
	idExp = 60 * 60 * 24
	unixTime = time.Now().Unix()
	tokenExp = unixTime + idExp
	timeExpire = jwt.NumericDate{Time: time.Unix(tokenExp, 0)}
	refreshClaim := &idTokenClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: &timeExpire,
			IssuedAt:  &timeNow,
			Issuer:    a.appConfig.AppName,
		},
		User:  user,
		Scope: role,
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaim)
	refreshTokenString, _ := refreshToken.SignedString(a.appConfig.JWTSecret)

	return tokenString, refreshTokenString, nil
}

func (a *authService) AuthAfterRegister(user *model.User, wallet *model.Wallet, tx *gorm.DB) (string, string, error) {
	userJWT := &dto.UserJWT{
		UserID:   user.ID,
		Email:    user.Email,
		Username: user.Username,
		WalletID: wallet.ID,
	}
	token, refreshToken, err := a.generateJWTToken(userJWT, model.UserRoleName)
	if os.Getenv("ENV") == "testing" {
		token = "test"
		refreshToken = "test"
	}

	err = a.refreshTokenRepo.CreateRefreshToken(tx, user.ID, refreshToken)
	if err != nil {
		tx.Rollback()
		return "", "", apperror.InternalServerError("Cannot add refresh token")
	}

	tx.Commit()
	return token, refreshToken, err
}

func (a *authService) SignInWithGoogle(user *model.User) (string, string, error) {
	tx := a.db.Begin()
	wallet, err := a.walletRepository.GetWalletByUserID(tx, user.ID)
	if err != nil {
		tx.Rollback()
		return "", "", err
	}
	userJWT := &dto.UserJWT{
		UserID:   user.ID,
		Email:    user.Email,
		Username: user.Username,
		WalletID: wallet.ID,
	}

	userRoles, err := a.userRoleRepo.GetRolesByUserID(tx, user.ID)
	if err != nil {
		tx.Rollback()
		return "", "", err
	}
	var roles []string
	for _, role := range userRoles {
		roles = append(roles, role.Role.Name)
	}
	rolesString := strings.Join(roles[:], " ")

	token, refreshToken, err := a.generateJWTToken(userJWT, rolesString)
	if os.Getenv("ENV") == "testing" {
		token = "test"
		refreshToken = "test"
	}
	err = a.refreshTokenRepo.CreateRefreshToken(tx, user.ID, refreshToken)
	if err != nil {
		tx.Rollback()
		return "", "", apperror.InternalServerError("Cannot add refresh token")
	}

	tx.Commit()
	return token, refreshToken, err
}

func (a *authService) SignIn(req *dto.SignInReq) (string, string, error) {
	tx := a.db.Begin()
	user, err := a.userRepository.MatchingCredential(tx, req.Email, req.Password)
	if err != nil || user == nil {
		tx.Rollback()
		return "", "", err
	}
	wallet, err := a.walletRepository.GetWalletByUserID(tx, user.ID)
	if err != nil {
		tx.Rollback()
		return "", "", err
	}
	userJWT := &dto.UserJWT{
		UserID:   user.ID,
		Email:    user.Email,
		WalletID: wallet.ID,
	}

	userRoles, err := a.userRoleRepo.GetRolesByUserID(tx, user.ID)
	if err != nil {
		tx.Rollback()
		return "", "", err
	}
	var roles []string
	for _, role := range userRoles {
		roles = append(roles, role.Role.Name)
	}
	rolesString := strings.Join(roles[:], " ")
	token, refreshToken, err := a.generateJWTToken(userJWT, rolesString)

	if os.Getenv("ENV") == "testing" {
		token = "test"
		refreshToken = "test"
	}
	err = a.refreshTokenRepo.CreateRefreshToken(tx, user.ID, refreshToken)
	if err != nil {
		tx.Rollback()
		return "", "", err
	}

	tx.Commit()
	return token, refreshToken, err
}

func (a *authService) SignOut(userID uint) error {
	tx := a.db.Begin()
	err := a.refreshTokenRepo.DeleteRefreshToken(tx, userID)
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}
