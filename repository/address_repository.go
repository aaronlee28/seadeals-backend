package repository

import (
	"gorm.io/gorm"
	"seadeals-backend/apperror"
	"seadeals-backend/model"
)

type AddressRepository interface {
	CreateAddress(*gorm.DB, *model.Address) (*model.Address, error)
	GetAddressesByUserID(*gorm.DB, uint) ([]*model.Address, error)
	GetAddressDetail(*gorm.DB, uint) (*model.Address, error)
	UpdateAddress(*gorm.DB, *model.Address) (*model.Address, error)
}

type addressRepository struct {
}

type AddressRepositoryConfig struct {
}

func NewAddressRepository(c *AddressRepositoryConfig) AddressRepository {
	return &addressRepository{}
}

func (a *addressRepository) CreateAddress(tx *gorm.DB, newAddress *model.Address) (*model.Address, error) {
	result := tx.Create(&newAddress)
	if result.Error != nil {
		return nil, apperror.InternalServerError("cannot create new Address")
	}

	return newAddress, result.Error
}

func (a *addressRepository) GetAddressesByUserID(tx *gorm.DB, userID uint) ([]*model.Address, error) {
	var addresses []*model.Address
	result := tx.Where("user_id = ?", userID).Find(&addresses)
	if result.Error != nil {
		return nil, apperror.InternalServerError("cannot fetch addresses")
	}

	return addresses, result.Error
}

func (a *addressRepository) GetAddressDetail(tx *gorm.DB, id uint) (*model.Address, error) {
	var address = &model.Address{}
	address.ID = id
	result := tx.First(&address)

	if result.Error == gorm.ErrRecordNotFound {
		return nil, apperror.InternalServerError("no record of that id exists")
	}
	if result.Error != nil {
		return nil, apperror.InternalServerError("cannot fetch address")
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
