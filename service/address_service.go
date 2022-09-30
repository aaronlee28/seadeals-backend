package service

import (
	"gorm.io/gorm"
	"seadeals-backend/dto"
	"seadeals-backend/helper"
	"seadeals-backend/model"
	"seadeals-backend/repository"
)

type AddressService interface {
	CreateAddress(*dto.CreateAddressReq, uint) (*model.Address, error)
	UpdateAddress(*dto.UpdateAddressReq) (*model.Address, error)
	GetAddressesByUserID(uint) ([]*dto.GetAddressRes, error)
	GetUserMainAddress(userID uint) (*dto.GetAddressRes, error)
	ChangeMainAddress(ID, userID uint) (*dto.GetAddressRes, error)
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

func (a *addressService) CreateAddress(req *dto.CreateAddressReq, userID uint) (*model.Address, error) {
	tx := a.db.Begin()
	var err error
	defer helper.CommitOrRollback(tx, &err)

	newAddress := &model.Address{
		UserID:      userID,
		CityID:      req.CityID,
		ProvinceID:  req.ProvinceID,
		Province:    req.Province,
		City:        req.City,
		Type:        req.Type,
		PostalCode:  req.PostalCode,
		SubDistrict: req.SubDistrict,
		Address:     req.Address,
		IsMain:      false,
	}
	address, err := a.addressRepository.CreateAddress(tx, newAddress)
	if err != nil {
		return nil, err
	}

	return address, nil
}

func (a *addressService) UpdateAddress(req *dto.UpdateAddressReq) (*model.Address, error) {
	tx := a.db.Begin()
	var err error
	defer helper.CommitOrRollback(tx, &err)

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
		return nil, err
	}

	return address, nil
}

func (a *addressService) GetAddressesByUserID(userID uint) ([]*dto.GetAddressRes, error) {
	tx := a.db.Begin()
	var err error
	defer helper.CommitOrRollback(tx, &err)

	addresses, err := a.addressRepository.GetAddressesByUserID(tx, userID)
	if err != nil {
		return nil, err
	}

	var res []*dto.GetAddressRes
	for _, addr := range addresses {
		res = append(res, new(dto.GetAddressRes).From(addr))
	}

	return res, nil
}

func (a *addressService) GetUserMainAddress(userID uint) (*dto.GetAddressRes, error) {
	tx := a.db.Begin()
	var err error
	defer helper.CommitOrRollback(tx, &err)

	addr, err := a.addressRepository.GetUserMainAddress(tx, userID)
	if err != nil {
		return nil, err
	}

	res := new(dto.GetAddressRes).From(addr)
	return res, nil
}

func (a *addressService) ChangeMainAddress(ID, userID uint) (*dto.GetAddressRes, error) {
	tx := a.db.Begin()
	var err error
	defer helper.CommitOrRollback(tx, &err)

	addr, err := a.addressRepository.ChangeMainAddress(tx, ID, userID)
	if err != nil {
		return nil, err
	}
	
	res := new(dto.GetAddressRes).From(addr)
	return res, nil
}
