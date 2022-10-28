// Code generated by mockery v2.14.0. DO NOT EDIT.

package mocks

import (
	dto "seadeals-backend/dto"
	model "seadeals-backend/model"

	mock "github.com/stretchr/testify/mock"
)

// DeliveryService is an autogenerated mock type for the DeliveryService type
type DeliveryService struct {
	mock.Mock
}

// DeliverOrder provides a mock function with given fields: req, userID
func (_m *DeliveryService) DeliverOrder(req *dto.DeliverOrderReq, userID uint) (*model.Delivery, error) {
	ret := _m.Called(req, userID)

	var r0 *model.Delivery
	if rf, ok := ret.Get(0).(func(*dto.DeliverOrderReq, uint) *model.Delivery); ok {
		r0 = rf(req, userID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Delivery)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*dto.DeliverOrderReq, uint) error); ok {
		r1 = rf(req, userID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetSellerPrintSettings provides a mock function with given fields: sellerID
func (_m *DeliveryService) GetSellerPrintSettings(sellerID uint) (*dto.DeliverSettingsPrintRes, error) {
	ret := _m.Called(sellerID)

	var r0 *dto.DeliverSettingsPrintRes
	if rf, ok := ret.Get(0).(func(uint) *dto.DeliverSettingsPrintRes); ok {
		r0 = rf(sellerID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*dto.DeliverSettingsPrintRes)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(uint) error); ok {
		r1 = rf(sellerID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdatePrintSettings provides a mock function with given fields: req, sellerID
func (_m *DeliveryService) UpdatePrintSettings(req *dto.DeliverSettingsPrint, sellerID uint) (*dto.DeliverSettingsPrintRes, error) {
	ret := _m.Called(req, sellerID)

	var r0 *dto.DeliverSettingsPrintRes
	if rf, ok := ret.Get(0).(func(*dto.DeliverSettingsPrint, uint) *dto.DeliverSettingsPrintRes); ok {
		r0 = rf(req, sellerID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*dto.DeliverSettingsPrintRes)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*dto.DeliverSettingsPrint, uint) error); ok {
		r1 = rf(req, sellerID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewDeliveryService interface {
	mock.TestingT
	Cleanup(func())
}

// NewDeliveryService creates a new instance of DeliveryService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewDeliveryService(t mockConstructorTestingTNewDeliveryService) *DeliveryService {
	mock := &DeliveryService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
