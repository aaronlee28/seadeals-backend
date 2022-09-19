package service

import (
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"gorm.io/gorm"
	"os"
	"seadeals-backend/apperror"
	"seadeals-backend/config"
	"seadeals-backend/dto"
	"seadeals-backend/model"
	"seadeals-backend/repository"
	"time"
)

type AuthService interface {
	AuthAfterRegister(user *model.User, wallet *model.Wallet, tx *gorm.DB) (string, string, error)
}

type authService struct {
	db               *gorm.DB
	refreshTokenRepo repository.RefreshTokenRepository
	appConfig        config.AppConfig
}

type AuthSConfig struct {
	DB               *gorm.DB
	RefreshTokenRepo repository.RefreshTokenRepository
	AppConfig        config.AppConfig
}

func NewAuthService(config *AuthSConfig) AuthService {
	return &authService{
		db:               config.DB,
		refreshTokenRepo: config.RefreshTokenRepo,
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
	token, refreshToken, err := a.generateJWTToken(userJWT, "user")
	if os.Getenv("ENV") == "testing" {
		fmt.Println("aku disini")
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
