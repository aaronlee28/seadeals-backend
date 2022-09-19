package service

import (
	"gorm.io/gorm"
	"seadeals-backend/apperror"
	"seadeals-backend/dto"
	"seadeals-backend/model"
	"seadeals-backend/repository"
)

type AddressService interface {
	CreateAddress(*dto.CreateAddressReq, uint) (*dto.CreateAddressRes, error)
	UpdateAddress(*dto.UpdateAddressReq) (*dto.UpdateAddressRes, error)
	GetAddressesByUserID(uint) ([]*model.Address, error)
}

type addressService struct {
	db                *gorm.DB
	addressRepository repository.AddressRepository
}

type AddressServiceConfig struct {
	DB                *gorm.DB
	AddressRepository repository.AddressRepository
}

func NewAddressService(config *AddressServiceConfig) AddressService {
	return &addressService{
		db:                config.DB,
		addressRepository: config.AddressRepository,
	}
}

func (a *addressService) CreateAddress(req *dto.CreateAddressReq, userID uint) (*dto.CreateAddressRes, error) {
	tx := a.db.Begin()
	newAddress := &model.Address{
		Address:       req.Address,
		Zipcode:       req.Zipcode,
		SubDistrictID: req.SubDistrictID,
		UserID:        userID,
	}
	address, err := a.addressRepository.CreateAddress(tx, newAddress)
	if err != nil {
		return nil, err
	}
	tx.Commit()

	dtoAddress := &dto.CreateAddressRes{
		ID:            address.ID,
		Address:       address.Address,
		Zipcode:       address.Zipcode,
		SubDistrictID: address.SubDistrictID,
		UserID:        address.UserID,
	}

	return dtoAddress, nil
}

func (a *addressService) UpdateAddress(req *dto.UpdateAddressReq) (*dto.UpdateAddressRes, error) {
	tx := a.db.Begin()
	currentAddress, err := a.addressRepository.GetAddressDetail(tx, req.ID)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	if currentAddress.UserID != req.UserID {
		return nil, apperror.ForbiddenError("cannot update another user address")
	}

	newAddress := &model.Address{
		Address:       req.Address,
		Zipcode:       req.Zipcode,
		SubDistrictID: req.SubDistrictID,
		UserID:        req.UserID,
		ID:            req.ID,
	}
	address, err := a.addressRepository.UpdateAddress(tx, newAddress)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	tx.Commit()

	dtoAddress := &dto.UpdateAddressRes{
		ID:            address.ID,
		Address:       address.Address,
		Zipcode:       address.Zipcode,
		SubDistrictID: address.SubDistrictID,
		UserID:        address.UserID,
	}

	return dtoAddress, nil
}

func (a *addressService) GetAddressesByUserID(userID uint) ([]*model.Address, error) {
	tx := a.db.Begin()
	addresses, err := a.addressRepository.GetAddressesByUserID(tx, userID)
	if err != nil {
		return nil, err
	}
	tx.Commit()
	return addresses, nil
}
