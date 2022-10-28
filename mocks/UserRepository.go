// Code generated by mockery v2.14.0. DO NOT EDIT.

package mocks

import (
	dto "seadeals-backend/dto"

	gorm "gorm.io/gorm"

	mock "github.com/stretchr/testify/mock"

	model "seadeals-backend/model"
)

// UserRepository is an autogenerated mock type for the UserRepository type
type UserRepository struct {
	mock.Mock
}

// ChangeUserDetailsLessPassword provides a mock function with given fields: tx, userID, details
func (_m *UserRepository) ChangeUserDetailsLessPassword(tx *gorm.DB, userID uint, details *dto.ChangeUserDetails) (*model.User, error) {
	ret := _m.Called(tx, userID, details)

	var r0 *model.User
	if rf, ok := ret.Get(0).(func(*gorm.DB, uint, *dto.ChangeUserDetails) *model.User); ok {
		r0 = rf(tx, userID, details)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.User)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*gorm.DB, uint, *dto.ChangeUserDetails) error); ok {
		r1 = rf(tx, userID, details)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ChangeUserPassword provides a mock function with given fields: tx, userID, newPassword
func (_m *UserRepository) ChangeUserPassword(tx *gorm.DB, userID uint, newPassword string) error {
	ret := _m.Called(tx, userID, newPassword)

	var r0 error
	if rf, ok := ret.Get(0).(func(*gorm.DB, uint, string) error); ok {
		r0 = rf(tx, userID, newPassword)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetUserByEmail provides a mock function with given fields: _a0, _a1
func (_m *UserRepository) GetUserByEmail(_a0 *gorm.DB, _a1 string) (*model.User, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *model.User
	if rf, ok := ret.Get(0).(func(*gorm.DB, string) *model.User); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.User)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*gorm.DB, string) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetUserByID provides a mock function with given fields: tx, userID
func (_m *UserRepository) GetUserByID(tx *gorm.DB, userID uint) (*model.User, error) {
	ret := _m.Called(tx, userID)

	var r0 *model.User
	if rf, ok := ret.Get(0).(func(*gorm.DB, uint) *model.User); ok {
		r0 = rf(tx, userID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.User)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*gorm.DB, uint) error); ok {
		r1 = rf(tx, userID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetUserDetailsByID provides a mock function with given fields: tx, userID
func (_m *UserRepository) GetUserDetailsByID(tx *gorm.DB, userID uint) (*model.User, error) {
	ret := _m.Called(tx, userID)

	var r0 *model.User
	if rf, ok := ret.Get(0).(func(*gorm.DB, uint) *model.User); ok {
		r0 = rf(tx, userID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.User)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*gorm.DB, uint) error); ok {
		r1 = rf(tx, userID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// HasExistEmail provides a mock function with given fields: _a0, _a1
func (_m *UserRepository) HasExistEmail(_a0 *gorm.DB, _a1 string) (bool, error) {
	ret := _m.Called(_a0, _a1)

	var r0 bool
	if rf, ok := ret.Get(0).(func(*gorm.DB, string) bool); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Get(0).(bool)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*gorm.DB, string) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MatchingCredential provides a mock function with given fields: _a0, _a1, _a2
func (_m *UserRepository) MatchingCredential(_a0 *gorm.DB, _a1 string, _a2 string) (*model.User, error) {
	ret := _m.Called(_a0, _a1, _a2)

	var r0 *model.User
	if rf, ok := ret.Get(0).(func(*gorm.DB, string, string) *model.User); ok {
		r0 = rf(_a0, _a1, _a2)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.User)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*gorm.DB, string, string) error); ok {
		r1 = rf(_a0, _a1, _a2)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Register provides a mock function with given fields: _a0, _a1
func (_m *UserRepository) Register(_a0 *gorm.DB, _a1 *model.User) (*model.User, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *model.User
	if rf, ok := ret.Get(0).(func(*gorm.DB, *model.User) *model.User); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.User)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*gorm.DB, *model.User) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// RegisterAsSeller provides a mock function with given fields: db, _a1
func (_m *UserRepository) RegisterAsSeller(db *gorm.DB, _a1 *model.Seller) (*model.Seller, error) {
	ret := _m.Called(db, _a1)

	var r0 *model.Seller
	if rf, ok := ret.Get(0).(func(*gorm.DB, *model.Seller) *model.Seller); ok {
		r0 = rf(db, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Seller)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*gorm.DB, *model.Seller) error); ok {
		r1 = rf(db, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewUserRepository interface {
	mock.TestingT
	Cleanup(func())
}

// NewUserRepository creates a new instance of UserRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewUserRepository(t mockConstructorTestingTNewUserRepository) *UserRepository {
	mock := &UserRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
