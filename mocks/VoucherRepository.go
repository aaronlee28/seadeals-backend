// Code generated by mockery v2.14.0. DO NOT EDIT.

package mocks

import (
	model "seadeals-backend/model"

	mock "github.com/stretchr/testify/mock"
	gorm "gorm.io/gorm"
)

// VoucherRepository is an autogenerated mock type for the VoucherRepository type
type VoucherRepository struct {
	mock.Mock
}

// CreateVoucher provides a mock function with given fields: tx, v
func (_m *VoucherRepository) CreateVoucher(tx *gorm.DB, v *model.Voucher) (*model.Voucher, error) {
	ret := _m.Called(tx, v)

	var r0 *model.Voucher
	if rf, ok := ret.Get(0).(func(*gorm.DB, *model.Voucher) *model.Voucher); ok {
		r0 = rf(tx, v)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Voucher)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*gorm.DB, *model.Voucher) error); ok {
		r1 = rf(tx, v)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DeleteVoucherByID provides a mock function with given fields: tx, id
func (_m *VoucherRepository) DeleteVoucherByID(tx *gorm.DB, id uint) (bool, error) {
	ret := _m.Called(tx, id)

	var r0 bool
	if rf, ok := ret.Get(0).(func(*gorm.DB, uint) bool); ok {
		r0 = rf(tx, id)
	} else {
		r0 = ret.Get(0).(bool)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*gorm.DB, uint) error); ok {
		r1 = rf(tx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FindVoucherByCode provides a mock function with given fields: tx, code
func (_m *VoucherRepository) FindVoucherByCode(tx *gorm.DB, code string) (*model.Voucher, error) {
	ret := _m.Called(tx, code)

	var r0 *model.Voucher
	if rf, ok := ret.Get(0).(func(*gorm.DB, string) *model.Voucher); ok {
		r0 = rf(tx, code)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Voucher)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*gorm.DB, string) error); ok {
		r1 = rf(tx, code)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FindVoucherByID provides a mock function with given fields: tx, id
func (_m *VoucherRepository) FindVoucherByID(tx *gorm.DB, id uint) (*model.Voucher, error) {
	ret := _m.Called(tx, id)

	var r0 *model.Voucher
	if rf, ok := ret.Get(0).(func(*gorm.DB, uint) *model.Voucher); ok {
		r0 = rf(tx, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Voucher)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*gorm.DB, uint) error); ok {
		r1 = rf(tx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FindVoucherBySellerID provides a mock function with given fields: tx, sellerID, qp
func (_m *VoucherRepository) FindVoucherBySellerID(tx *gorm.DB, sellerID uint, qp *model.VoucherQueryParam) ([]*model.Voucher, int64, error) {
	ret := _m.Called(tx, sellerID, qp)

	var r0 []*model.Voucher
	if rf, ok := ret.Get(0).(func(*gorm.DB, uint, *model.VoucherQueryParam) []*model.Voucher); ok {
		r0 = rf(tx, sellerID, qp)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*model.Voucher)
		}
	}

	var r1 int64
	if rf, ok := ret.Get(1).(func(*gorm.DB, uint, *model.VoucherQueryParam) int64); ok {
		r1 = rf(tx, sellerID, qp)
	} else {
		r1 = ret.Get(1).(int64)
	}

	var r2 error
	if rf, ok := ret.Get(2).(func(*gorm.DB, uint, *model.VoucherQueryParam) error); ok {
		r2 = rf(tx, sellerID, qp)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// FindVoucherDetailByID provides a mock function with given fields: tx, id
func (_m *VoucherRepository) FindVoucherDetailByID(tx *gorm.DB, id uint) (*model.Voucher, error) {
	ret := _m.Called(tx, id)

	var r0 *model.Voucher
	if rf, ok := ret.Get(0).(func(*gorm.DB, uint) *model.Voucher); ok {
		r0 = rf(tx, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Voucher)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*gorm.DB, uint) error); ok {
		r1 = rf(tx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetAvailableGlobalVouchers provides a mock function with given fields: tx
func (_m *VoucherRepository) GetAvailableGlobalVouchers(tx *gorm.DB) ([]*model.Voucher, error) {
	ret := _m.Called(tx)

	var r0 []*model.Voucher
	if rf, ok := ret.Get(0).(func(*gorm.DB) []*model.Voucher); ok {
		r0 = rf(tx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*model.Voucher)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*gorm.DB) error); ok {
		r1 = rf(tx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetVouchersBySellerID provides a mock function with given fields: tx, sellerID
func (_m *VoucherRepository) GetVouchersBySellerID(tx *gorm.DB, sellerID uint) ([]*model.Voucher, error) {
	ret := _m.Called(tx, sellerID)

	var r0 []*model.Voucher
	if rf, ok := ret.Get(0).(func(*gorm.DB, uint) []*model.Voucher); ok {
		r0 = rf(tx, sellerID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*model.Voucher)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*gorm.DB, uint) error); ok {
		r1 = rf(tx, sellerID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateVoucher provides a mock function with given fields: tx, v, id
func (_m *VoucherRepository) UpdateVoucher(tx *gorm.DB, v *model.Voucher, id uint) (*model.Voucher, error) {
	ret := _m.Called(tx, v, id)

	var r0 *model.Voucher
	if rf, ok := ret.Get(0).(func(*gorm.DB, *model.Voucher, uint) *model.Voucher); ok {
		r0 = rf(tx, v, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Voucher)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*gorm.DB, *model.Voucher, uint) error); ok {
		r1 = rf(tx, v, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewVoucherRepository interface {
	mock.TestingT
	Cleanup(func())
}

// NewVoucherRepository creates a new instance of VoucherRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewVoucherRepository(t mockConstructorTestingTNewVoucherRepository) *VoucherRepository {
	mock := &VoucherRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
