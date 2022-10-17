package repository

import (
	"gorm.io/gorm"
	"seadeals-backend/apperror"
	"seadeals-backend/dto"
	"seadeals-backend/model"
)

type AddressRepository interface {
	GetAddressesByUserID(*gorm.DB, uint) ([]*model.Address, error)
	GetAddressesByID(tx *gorm.DB, id, userID uint) (*model.Address, error)

	UpdateAddress(*gorm.DB, *model.Address) (*model.Address, error)
	CreateAddress(tx *gorm.DB, req *dto.CreateAddressReq, userID uint) (*model.Address, error)

	CheckUserAddress(tx *gorm.DB, addressID uint, userID uint) (*model.Address, error)
	GetUserMainAddress(tx *gorm.DB, userID uint) (*model.Address, error)
	ChangeMainAddress(tx *gorm.DB, ID, userID uint) (*model.Address, error)
}

type addressRepository struct{}

func NewAddressRepository() AddressRepository {
	return &addressRepository{}
}

func (a *addressRepository) CreateAddress(tx *gorm.DB, req *dto.CreateAddressReq, userID uint) (*model.Address, error) {
	var newAddress = &model.Address{}
	var isMain = false
	result := tx.Model(&newAddress).Where("user_id = ?", userID).Where("is_main IS TRUE").First(&newAddress)
	if result.Error != nil {
		if result.Error != gorm.ErrRecordNotFound {
			return nil, apperror.InternalServerError("Cannot find main ")
		}
		isMain = true
	}

	newAddress = &model.Address{
		UserID:      userID,
		CityID:      req.CityID,
		ProvinceID:  req.ProvinceID,
		Province:    req.Province,
		City:        req.City,
		Type:        req.Type,
		PostalCode:  req.PostalCode,
		SubDistrict: req.SubDistrict,
		Address:     req.Address,
		IsMain:      isMain,
	}
	result = tx.Create(&newAddress)
	if result.Error != nil {
		return nil, apperror.InternalServerError("Cannot register new Address")
	}
	return newAddress, nil
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

func (a *addressRepository) CheckUserAddress(tx *gorm.DB, addressID uint, userID uint) (*model.Address, error) {
	var address = &model.Address{}
	address.ID = addressID
	result := tx.Model(&address).Unscoped().Where("user_id = ?", userID).First(&address)
	if result.Error != nil {
		if result.Error != gorm.ErrRecordNotFound {
			return nil, apperror.InternalServerError("Cannot find address")
		}
		return nil, apperror.BadRequestError("The user does not have this address")
	}

	return address, nil
}

func (a *addressRepository) GetUserMainAddress(tx *gorm.DB, userID uint) (*model.Address, error) {
	var address *model.Address
	result := tx.Model(&address).Where("user_id = ? AND is_main IS TRUE", userID).First(&address)
	if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
		return nil, apperror.InternalServerError("Cannot use database to find Address")
	}
	if result.Error == gorm.ErrRecordNotFound {
		return nil, apperror.NotFoundError("Main address not found")
	}

	return address, nil
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
