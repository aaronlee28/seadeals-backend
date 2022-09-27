package repository

import (
	"gorm.io/gorm"
	"seadeals-backend/apperror"
	"seadeals-backend/model"
)

type UserAddressRepository interface {
	CreateUserAddress(*gorm.DB, *model.UserAddress) (*model.UserAddress, error)
	GetUserAddressesByUserID(*gorm.DB, uint) ([]*model.UserAddress, error)
}

type userAddressRepository struct{}

func NewUserAddressRepository() UserAddressRepository {
	return &userAddressRepository{}
}

func (u *userAddressRepository) CreateUserAddress(tx *gorm.DB, userAddress *model.UserAddress) (*model.UserAddress, error) {
	result := tx.Model(userAddress).Where("user_id = ? AND is_main IS TRUE", userAddress.UserID).First(&model.UserAddress{})
	if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
		return nil, apperror.InternalServerError("Cannot use database to find user addresses")
	}
	userAddress.IsMain = true
	if result.Error != gorm.ErrRecordNotFound {
		userAddress.IsMain = false
	}

	result = tx.Create(&userAddress)
	if result.Error != nil {
		return nil, apperror.InternalServerError("cannot create user address")
	}

	return userAddress, result.Error
}

func (u *userAddressRepository) GetUserAddressesByUserID(tx *gorm.DB, userID uint) ([]*model.UserAddress, error) {
	var addresses []*model.UserAddress
	result := tx.Where("user_id = ?", userID).Preload("Address.UserAddress").Find(&addresses)
	if result.Error != nil {
		return nil, apperror.InternalServerError("cannot fetch addresses")
	}

	return addresses, result.Error
}
