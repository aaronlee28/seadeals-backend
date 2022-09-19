package repository

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"seadeals-backend/apperror"
	"seadeals-backend/model"
)

type UserRepository interface {
	Register(*gorm.DB, *model.User) (*model.User, error)
}

type userRepository struct {
}

type UserRepositoryConfig struct {
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func NewUserRepository(c *UserRepositoryConfig) UserRepository {
	return &userRepository{}
}

func (u *userRepository) Register(tx *gorm.DB, user *model.User) (*model.User, error) {
	var err error

	sameEmail := tx.Model(&model.User{}).Where("email LIKE ?", user.Email).First(&model.User{})
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
