package repository

import (
	"gorm.io/gorm"
	"seadeals-backend/apperror"
	"seadeals-backend/model"
)

type UserRoleRepository interface {
	CreateRoleToUser(*gorm.DB, *model.UserRole) (*model.UserRole, error)
	GetRolesByUserID(*gorm.DB, uint) ([]*model.UserRole, error)
}

type userRoleRepository struct {
}

type UserRoleRepositoryConfig struct {
}

func NewUserRoleRepository(c *UserRoleRepositoryConfig) UserRoleRepository {
	return &userRoleRepository{}
}

func (u *userRoleRepository) CreateRoleToUser(tx *gorm.DB, userRole *model.UserRole) (*model.UserRole, error) {
	result := tx.Create(&userRole)
	if result.Error != nil {
		return nil, apperror.InternalServerError("cannot create role")
	}

	return userRole, result.Error
}

func (u *userRoleRepository) GetRolesByUserID(tx *gorm.DB, userID uint) ([]*model.UserRole, error) {
	var userRoles []*model.UserRole
	result := tx.Where("user_id = ?", userID).Find(&userRoles)
	if result.Error != nil {
		return nil, apperror.InternalServerError("unable to get user roles")
	}

	return userRoles, nil
}
