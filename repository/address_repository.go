package repository

import (
	"gorm.io/gorm"
	"seadeals-backend/apperror"
	"seadeals-backend/model"
)

type AddressRepository interface {
	CreateAddress(*gorm.DB, *model.Address) (*model.Address, error)
	GetAddressesByUserID(*gorm.DB, uint) ([]*model.Address, error)
	GetAddressesByID(tx *gorm.DB, id, userID uint) (*model.Address, error)
	UpdateAddress(*gorm.DB, *model.Address) (*model.Address, error)
	ChangeMainAddress(tx *gorm.DB, ID, userID uint) (*model.Address, error)
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

func (a *addressRepository) GetAddressesByUserID(tx *gorm.DB, userID uint) ([]*model.Address, error) {
	var addresses []*model.Address
	result := tx.Where("user_id = ?", userID).Find(&addresses)
	if result.Error != nil {
		return nil, apperror.InternalServerError("cannot fetch addresses")
	}

	return addresses, result.Error
}

func (a *addressRepository) GetAddressesByID(tx *gorm.DB, id, userID uint) (*model.Address, error) {
	var address *model.Address
	result := tx.Where("user_id = ?", userID).First(&address, id)
	if result.Error != nil {
		return nil, apperror.InternalServerError("cannot fetch addresses")
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

func (a *addressRepository) ChangeMainAddress(tx *gorm.DB, ID, userID uint) (*model.Address, error) {
	ud := &model.Address{
		ID:     ID,
		IsMain: true,
	}

	result := tx.Model(&model.Address{}).Where("user_id = ? AND is_main = true", userID).Update("is_main", false)
	if result.Error != nil {
		return nil, result.Error
	}

	result = tx.Where("user_id = ?", userID).Updates(&ud).First(&ud, ID)
	if result.Error != nil {
		return nil, result.Error
	}

	return ud, nil
}
