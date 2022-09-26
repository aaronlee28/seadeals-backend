package service

import (
	"gorm.io/gorm"
	"seadeals-backend/apperror"
	"seadeals-backend/dto"
	"seadeals-backend/model"
	"seadeals-backend/repository"
)

type UserSeaPayAccountServ interface {
	RegisterSeaLabsPayAccount(req *dto.RegisterSeaLabsPayReq) (*dto.RegisterSeaLabsPayRes, error)
	CheckSeaLabsAccountExists(req *dto.CheckSeaLabsPayReq) (*dto.CheckSeaLabsPayRes, error)
	UpdateSeaLabsAccountToMain(req *dto.UpdateSeaLabsPayToMainReq) (*model.UserSealabsPayAccount, error)
	GetSeaLabsAccountByUserID(userID uint) ([]*model.UserSealabsPayAccount, error)
}

type userSeaPayAccountServ struct {
	db                    *gorm.DB
	userSeaPayAccountRepo repository.UserSeaPayAccountRepo
}

type UserSeaPayAccountServConfig struct {
	DB                    *gorm.DB
	UserSeaPayAccountRepo repository.UserSeaPayAccountRepo
}

func NewUserSeaPayAccountServ(c *UserSeaPayAccountServConfig) UserSeaPayAccountServ {
	return &userSeaPayAccountServ{
		db:                    c.DB,
		userSeaPayAccountRepo: c.UserSeaPayAccountRepo,
	}
}

func (u *userSeaPayAccountServ) CheckSeaLabsAccountExists(req *dto.CheckSeaLabsPayReq) (*dto.CheckSeaLabsPayRes, error) {
	tx := u.db.Begin()

	hasExists, err := u.userSeaPayAccountRepo.HasExistsSeaLabsPayAccountWith(tx, req.UserID, req.AccountNumber)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	response := &dto.CheckSeaLabsPayRes{IsExists: hasExists}

	tx.Commit()
	return response, nil
}

func (u *userSeaPayAccountServ) RegisterSeaLabsPayAccount(req *dto.RegisterSeaLabsPayReq) (*dto.RegisterSeaLabsPayRes, error) {
	tx := u.db.Begin()

	hasExists, err := u.userSeaPayAccountRepo.HasExistsSeaLabsPayAccountWith(tx, req.UserID, req.AccountNumber)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	if hasExists {
		tx.Rollback()
		return nil, apperror.BadRequestError("Sea Labs Pay Account is already registered")
	}

	seaLabsPayAccount, err := u.userSeaPayAccountRepo.RegisterSeaLabsPayAccount(tx, req)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	response := &dto.RegisterSeaLabsPayRes{
		Status:         "Completed",
		SeaLabsAccount: seaLabsPayAccount,
	}

	tx.Commit()
	return response, nil
}

func (u *userSeaPayAccountServ) UpdateSeaLabsAccountToMain(req *dto.UpdateSeaLabsPayToMainReq) (*model.UserSealabsPayAccount, error) {
	tx := u.db.Begin()

	updatedData, err := u.userSeaPayAccountRepo.UpdateSeaLabsPayAccountToMain(tx, req)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	tx.Commit()
	return updatedData, nil
}

func (u *userSeaPayAccountServ) GetSeaLabsAccountByUserID(userID uint) ([]*model.UserSealabsPayAccount, error) {
	tx := u.db.Begin()

	accounts, err := u.userSeaPayAccountRepo.GetSeaLabsPayAccountByUserID(tx, userID)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	tx.Commit()
	return accounts, nil
}
