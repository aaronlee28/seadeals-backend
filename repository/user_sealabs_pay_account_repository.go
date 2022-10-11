package repository

import (
	"gorm.io/gorm"
	"seadeals-backend/apperror"
	"seadeals-backend/dto"
	"seadeals-backend/model"
	"time"
)

type UserSeaPayAccountRepo interface {
	HasExistsSeaLabsPayAccountWith(tx *gorm.DB, userID uint, accountNumber string) (bool, error)
	RegisterSeaLabsPayAccount(tx *gorm.DB, req *dto.RegisterSeaLabsPayReq, userID uint) (*model.UserSealabsPayAccount, error)
	UpdateSeaLabsPayAccountToMain(tx *gorm.DB, req *dto.UpdateSeaLabsPayToMainReq, userID uint) (*model.UserSealabsPayAccount, error)
	GetSeaLabsPayAccountByUserID(tx *gorm.DB, userID uint) ([]*model.UserSealabsPayAccount, error)
}

type userSeaPayAccountRepo struct{}

func NewSeaPayAccountRepo() UserSeaPayAccountRepo {
	return &userSeaPayAccountRepo{}
}

func (u *userSeaPayAccountRepo) RegisterSeaLabsPayAccount(tx *gorm.DB, req *dto.RegisterSeaLabsPayReq, userID uint) (*model.UserSealabsPayAccount, error) {
	var newSeaAccount = &model.UserSealabsPayAccount{}
	var isMain = false
	result := tx.Model(&newSeaAccount).Where("user_id = ?", userID).Where("is_main IS TRUE").First(&newSeaAccount)
	if result.Error != nil {
		if result.Error != gorm.ErrRecordNotFound {
			return nil, apperror.InternalServerError("Cannot find main sea labs pay account ")
		}
		isMain = true
	}

	newSeaAccount = &model.UserSealabsPayAccount{
		UserID:        userID,
		ActiveDate:    time.Now(),
		AccountNumber: req.AccountNumber,
		Name:          req.Name,
		IsMain:        isMain,
	}
	result = tx.Create(&newSeaAccount)
	if result.Error != nil {
		return nil, apperror.InternalServerError("Cannot register sealabs pay account")
	}
	return newSeaAccount, nil
}

func (u *userSeaPayAccountRepo) UpdateSeaLabsPayAccountToMain(tx *gorm.DB, req *dto.UpdateSeaLabsPayToMainReq, userID uint) (*model.UserSealabsPayAccount, error) {
	var updateSeaAccount = &model.UserSealabsPayAccount{}
	result := tx.Model(&updateSeaAccount).Where("user_id = ?", userID).Where("account_number LIKE ?", req.AccountNumber).First(&updateSeaAccount)
	if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
		return nil, apperror.InternalServerError("Cannot find sea labs pay account")
	}

	if result.Error == gorm.ErrRecordNotFound {
		return nil, apperror.InternalServerError("That Sea labs pay account hasn't been registered to your SeaDeals account")
	}

	updateSeaAccount.ID = 0
	result = tx.Model(&updateSeaAccount).Where("user_id = ?", userID).Where("is_main IS TRUE").First(&updateSeaAccount)
	if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
		return nil, apperror.InternalServerError("Cannot find main sea labs pay account ")
	}

	if result.Error == nil {
		result = tx.Model(&updateSeaAccount).Update("is_main", false)
		if result.Error != nil {
			return nil, apperror.InternalServerError("Cannot update old main sea labs account ")
		}
	}

	updateSeaAccount.IsMain = true
	result = tx.Where("user_id = ?", userID).Where("account_number LIKE ?", req.AccountNumber).Updates(&updateSeaAccount)
	if result.Error != nil {
		return nil, apperror.InternalServerError("Cannot update sea labs pay account to main")
	}
	return updateSeaAccount, nil
}

func (u *userSeaPayAccountRepo) HasExistsSeaLabsPayAccountWith(tx *gorm.DB, userID uint, accountNumber string) (bool, error) {
	var account *model.UserSealabsPayAccount
	result := tx.Model(&account).Where("user_id = ?", userID).Where("account_number LIKE ?", accountNumber).First(&account)
	if result.Error == nil {
		return true, nil
	}

	if result.Error != gorm.ErrRecordNotFound {
		return false, apperror.InternalServerError("Cannot fetch sea labs pay account")
	}

	return false, nil
}

func (u *userSeaPayAccountRepo) GetSeaLabsPayAccountByUserID(tx *gorm.DB, userID uint) ([]*model.UserSealabsPayAccount, error) {
	var accounts []*model.UserSealabsPayAccount
	result := tx.Model(&accounts).Where("user_id = ?", userID).Find(&accounts)
	if result.Error != nil {
		return nil, apperror.InternalServerError("cannot find accounts")
	}
	return accounts, nil
}
