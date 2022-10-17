package repository

import (
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"seadeals-backend/apperror"
	"seadeals-backend/dto"
	"seadeals-backend/model"
)

type UserRepository interface {
	Register(*gorm.DB, *model.User) (*model.User, error)
	HasExistEmail(*gorm.DB, string) (bool, error)
	GetUserByEmail(*gorm.DB, string) (*model.User, error)
	GetUserByID(tx *gorm.DB, userID uint) (*model.User, error)
	MatchingCredential(*gorm.DB, string, string) (*model.User, error)
	RegisterAsSeller(db *gorm.DB, model *model.Seller) (*model.Seller, error)
	GetUserDetailsByID(tx *gorm.DB, userID uint) (*model.User, error)
	ChangeUserDetailsLessPassword(tx *gorm.DB, userID uint, details *dto.ChangeUserDetails) (*model.User, error)
	ChangeUserPassword(tx *gorm.DB, userID uint, newPassword string) error
}

type userRepository struct{}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func NewUserRepository() UserRepository {
	return &userRepository{}
}

func (u *userRepository) Register(tx *gorm.DB, user *model.User) (*model.User, error) {
	var err error

	sameEmail := tx.Model(&model.User{}).Unscoped().Where("email LIKE ?", user.Email).First(&model.User{})
	if sameEmail.Error == nil {
		return nil, apperror.BadRequestError("Email has already exists")
	}

	// created user must be user role first cannot be defined
	user.Password, err = hashPassword(user.Password)
	if err != nil {
		return nil, apperror.BadRequestError("password format is invalid")
	}
	result := tx.Create(&user)
	if result.Error != nil {
		return nil, apperror.InternalServerError("cannot create new user")
	}
	result.Find(&user)

	// DO NOT PASS HASHED PASSWORD
	user.Password = ""
	return user, result.Error
}

func (u *userRepository) HasExistEmail(tx *gorm.DB, email string) (bool, error) {
	existedEmail := tx.Model(&model.User{}).Where("email LIKE ?", email).First(&model.User{})
	if existedEmail.Error == nil {
		return true, nil
	}

	if existedEmail.Error == gorm.ErrRecordNotFound {
		return false, nil
	}

	return false, existedEmail.Error
}

func (u *userRepository) GetUserByEmail(tx *gorm.DB, email string) (*model.User, error) {
	var user = &model.User{}
	result := tx.Model(&user).Where("email LIKE ?", email).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}

	return user, nil
}

func (u *userRepository) GetUserByID(tx *gorm.DB, userID uint) (*model.User, error) {
	var user = &model.User{ID: userID}
	result := tx.Model(&user).First(&user)
	if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
		return nil, apperror.InternalServerError("Cannot use database to find user")
	}
	if result.Error == gorm.ErrRecordNotFound {
		return nil, apperror.NotFoundError("Cannot found user")
	}

	return user, nil
}

func (u *userRepository) MatchingCredential(tx *gorm.DB, email string, password string) (*model.User, error) {
	var user model.User
	query := tx.Model(&user).Where("email = ?", email).First(&user)
	err := query.Error

	isNotFound := errors.Is(err, gorm.ErrRecordNotFound)
	if isNotFound {
		return nil, apperror.BadRequestError("Invalid email or password")
	}

	match := checkPasswordHash(password, user.Password)
	if !match {
		return nil, apperror.BadRequestError("Invalid email or password")
	}

	// do not show hashed password to service
	user.Password = ""
	return &user, err
}

func (u *userRepository) RegisterAsSeller(tx *gorm.DB, seller *model.Seller) (*model.Seller, error) {
	result := tx.Model(&model.Seller{}).Where("user_id = ?", seller.UserID).First(&model.Seller{})
	if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
		return nil, apperror.InternalServerError("Cannot search database for seller")
	}
	if result.Error == nil {
		return nil, apperror.BadRequestError("This user already registered as seller")
	}

	result = tx.Model(&model.Seller{}).Where("name ILIKE ?", seller.Name).First(&model.Seller{})
	if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
		return nil, apperror.InternalServerError("Cannot search database for seller")
	}
	if result.Error == nil {
		return nil, apperror.BadRequestError("Shop name is already registered, please choose another name for the shop")
	}

	result = tx.Clauses(clause.Returning{}).Preload("Address").Preload("User").Create(&seller).First(&seller)
	if result.Error != nil {
		return nil, apperror.InternalServerError("Cannot register user as seller")
	}
	return seller, nil
}

func (u *userRepository) GetUserDetailsByID(tx *gorm.DB, userID uint) (*model.User, error) {
	var userDetails *model.User
	result := tx.First(&userDetails, userID)
	if result.Error != nil {
		return nil, result.Error
	}
	return userDetails, nil

}

func (u *userRepository) ChangeUserDetailsLessPassword(tx *gorm.DB, userID uint, details *dto.ChangeUserDetails) (*model.User, error) {
	var userDetails *model.User

	result := tx.Model(&userDetails).Where("id = ?", userID).Updates(&details)
	result = tx.First(&userDetails, userID)
	fmt.Println("detailssss", details)
	if result.Error != nil {
		return nil, apperror.InternalServerError("Cannot update user profile")
	}
	return userDetails, nil
}

func (u *userRepository) ChangeUserPassword(tx *gorm.DB, userID uint, newPassword string) error {
	var userDetails *model.User

	hashedPassword, err := hashPassword(newPassword)

	if err != nil {
		return apperror.InternalServerError("Failed to hash password")
	}

	result := tx.Model(&userDetails).Where("id = ?", userID).Update("password", hashedPassword)
	if result.Error != nil {
		return apperror.InternalServerError("Cannot change user's password")
	}

	return nil
}
