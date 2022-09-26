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

func (a *addressService) CreateAddress(req *dto.CreateAddressReq, userID uint) (*dto.CreateAddressRes, error) {
	tx := a.db.Begin()
	newAddress := &model.Address{
		Address:       req.Address,
		Zipcode:       req.Zipcode,
		SubDistrictID: req.SubDistrictID,
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

	dtoAddress := &dto.CreateAddressRes{
		ID:            address.ID,
		Address:       address.Address,
		Zipcode:       address.Zipcode,
		SubDistrictID: address.SubDistrictID,
		UserID:        userAddress.UserID,
	}

	tx.Commit()
	return dtoAddress, nil
}

func (a *addressService) UpdateAddress(req *dto.UpdateAddressReq) (*dto.UpdateAddressRes, error) {
	tx := a.db.Begin()
	currentAddress, err := a.addressRepository.GetAddressDetail(tx, req.ID)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	if currentAddress.UserAddress.UserID != req.UserID {
		return nil, apperror.ForbiddenError("cannot update another user address")
	}

	newAddress := &model.Address{
		Address:       req.Address,
		Zipcode:       req.Zipcode,
		SubDistrictID: req.SubDistrictID,
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
	}

	return dtoAddress, nil
}

func (a *addressService) GetAddressesByUserID(userID uint) ([]*model.Address, error) {
	tx := a.db.Begin()
	userAddresses, err := a.userAddressRepo.GetUserAddressesByUserID(tx, userID)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	var addresses []*model.Address
	for _, userAddress := range userAddresses {
		addresses = append(addresses, userAddress.Address)
	}

	tx.Commit()
	return addresses, nil
}
