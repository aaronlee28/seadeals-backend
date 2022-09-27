package service

import (
	"gorm.io/gorm"
	"seadeals-backend/apperror"
	"seadeals-backend/dto"
	"seadeals-backend/model"
	"seadeals-backend/repository"
)

type AddressService interface {
	CreateAddress(*dto.CreateAddressReq, uint) (*model.Address, error)
	UpdateAddress(*dto.UpdateAddressReq) (*model.Address, error)
	GetAddressesByUserID(uint) ([]*dto.GetAddressRes, error)
}

type addressService struct {
	db                *gorm.DB
	addressRepository repository.AddressRepository
	userAddressRepo   repository.UserAddressRepository
}

type AddressServiceConfig struct {
	DB                *gorm.DB
	AddressRepository repository.AddressRepository
	UserAddressRepo   repository.UserAddressRepository
}

func NewAddressService(config *AddressServiceConfig) AddressService {
	return &addressService{
		db:                config.DB,
		addressRepository: config.AddressRepository,
		userAddressRepo:   config.UserAddressRepo,
	}
}

func (a *addressService) CreateAddress(req *dto.CreateAddressReq, userID uint) (*model.Address, error) {
	tx := a.db.Begin()
	newAddress := &model.Address{
		CityID:      req.CityID,
		ProvinceID:  req.ProvinceID,
		Province:    req.Province,
		City:        req.City,
		Type:        req.Type,
		PostalCode:  req.PostalCode,
		SubDistrict: req.SubDistrict,
		Address:     req.Address,
	}
	address, err := a.addressRepository.CreateAddress(tx, newAddress)
	if err != nil {
		return nil, err
	}

	newUserAddress := &model.UserAddress{
		AddressID: address.ID,
		UserID:    userID,
	}
	userAddress, err := a.userAddressRepo.CreateUserAddress(tx, newUserAddress)
	if err != nil {
		return nil, err
	}

	address.UserAddress = userAddress

	tx.Commit()
	return address, nil
}

func (a *addressService) UpdateAddress(req *dto.UpdateAddressReq) (*model.Address, error) {
	tx := a.db.Begin()
	currentAddress, err := a.addressRepository.GetAddressesByID(tx, req.ID, req.UserID)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	if currentAddress.UserAddress.UserID != req.UserID {
		return nil, apperror.ForbiddenError("cannot update another user address")
	}

	newAddress := &model.Address{
		ID:          req.ID,
		CityID:      req.CityID,
		ProvinceID:  req.ProvinceID,
		Province:    req.Province,
		City:        req.City,
		Type:        req.Type,
		PostalCode:  req.PostalCode,
		SubDistrict: req.SubDistrict,
		Address:     req.Address,
	}
	address, err := a.addressRepository.UpdateAddress(tx, newAddress)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	tx.Commit()
	return address, nil
}

func (a *addressService) GetAddressesByUserID(userID uint) ([]*dto.GetAddressRes, error) {
	tx := a.db.Begin()
	userAddresses, err := a.userAddressRepo.GetUserAddressesByUserID(tx, userID)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	var addresses []*dto.GetAddressRes
	for _, userAddress := range userAddresses {
		addresses = append(addresses, new(dto.GetAddressRes).From(userAddress.Address))
	}

	tx.Commit()
	return addresses, nil
}
