package repository

import (
	"gorm.io/gorm"
	"seadeals-backend/apperror"
	"seadeals-backend/model"
)

type AddressRepository interface {
	CreateAddress(*gorm.DB, *model.Address) (*model.Address, error)
	GetAddressesByUserID(*gorm.DB, uint) ([]*model.UserAddress, error)
	GetAddressesByID(tx *gorm.DB, id, userID uint) (*model.Address, error)
	UpdateAddress(*gorm.DB, *model.Address) (*model.Address, error)
}

type addressRepository struct{}

func NewAddressRepository() AddressRepository {
	return &addressRepository{}
}

func (a *addressRepository) CreateAddress(tx *gorm.DB, newAddress *model.Address) (*model.Address, error) {
	result := tx.Create(&newAddress)
	if result.Error != nil {
		return nil, apperror.InternalServerError("cannot create new Address")
	}

	return newAddress, result.Error
}

func (a *addressRepository) GetAddressesByUserID(tx *gorm.DB, userID uint) ([]*model.UserAddress, error) {
	var addresses []*model.UserAddress
	result := tx.Where("user_id = ?", userID).Preload("Address").Find(&addresses)
	if result.Error != nil {
		return nil, apperror.InternalServerError("cannot fetch addresses")
	}

	return addresses, result.Error
}

func (a *addressRepository) GetAddressesByID(tx *gorm.DB, id, userID uint) (*model.Address, error) {
	var address *model.Address
	result := tx.First(&address, id)
	if result.Error != nil {
		return nil, apperror.InternalServerError("cannot fetch addresses")
	}

	result = tx.Where("address_id = ? AND user_id = ?", id, userID).First(&address.UserAddress)
	if result.Error != nil {
		return nil, apperror.InternalServerError("cannot fetch someone address")
	}

	return address, result.Error
}

func (a *addressRepository) UpdateAddress(tx *gorm.DB, newAddress *model.Address) (*model.Address, error) {
	result := tx.Model(&newAddress).Updates(&newAddress)
	if result.Error != nil {
		return nil, apperror.InternalServerError("cannot update address")
	}

	return newAddress, result.Error
}
