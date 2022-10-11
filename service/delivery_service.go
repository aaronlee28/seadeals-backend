package service

import (
	"gorm.io/gorm"
	"seadeals-backend/apperror"
	"seadeals-backend/dto"
	"seadeals-backend/helper"
	"seadeals-backend/model"
	"seadeals-backend/repository"
)

type DeliveryService interface {
	DeliverOrder(req *dto.DeliverOrderReq, userID uint) (*model.Delivery, error)
}

type deliveryService struct {
	db                   *gorm.DB
	deliveryRepository   repository.DeliveryRepository
	deliveryActivityRepo repository.DeliveryActivityRepository
	sellerRepository     repository.SellerRepository
	addressRepository    repository.AddressRepository
	orderRepository      repository.OrderRepository
}

type DeliveryServiceConfig struct {
	DB                  *gorm.DB
	DeliveryRepository  repository.DeliveryRepository
	DeliverActivityRepo repository.DeliveryActivityRepository
	SellerRepository    repository.SellerRepository
	AddressRepository   repository.AddressRepository
	OrderRepository     repository.OrderRepository
}

func NewDeliveryService(c *DeliveryServiceConfig) DeliveryService {
	return &deliveryService{
		db:                   c.DB,
		deliveryRepository:   c.DeliveryRepository,
		deliveryActivityRepo: c.DeliverActivityRepo,
		sellerRepository:     c.SellerRepository,
		addressRepository:    c.AddressRepository,
		orderRepository:      c.OrderRepository,
	}
}

func (d *deliveryService) DeliverOrder(req *dto.DeliverOrderReq, userID uint) (*model.Delivery, error) {
	tx := d.db.Begin()
	var err error
	defer helper.CommitOrRollback(tx, &err)

	order, err := d.orderRepository.GetOrderDetailByID(tx, req.OrderID)
	if err != nil {
		return nil, err
	}
	if order.Status != dto.OrderWaitingSeller {
		return nil, apperror.BadRequestError("Cannot deliver order that has status " + order.Status)
	}

	seller, err := d.sellerRepository.FindSellerByUserID(tx, userID)
	if err != nil {
		return nil, err
	}
	if seller.ID != order.SellerID {
		return nil, apperror.BadRequestError("You cannot deliver another seller order")
	}

	delivery, err := d.deliveryRepository.GetDeliveryByOrderID(tx, order.ID)
	if err != nil {
		return nil, err
	}

	updatedDelivery, err := d.deliveryRepository.UpdateDeliveryStatus(tx, delivery.ID, dto.DeliveryOngoing)
	if err != nil {
		return nil, err
	}
	_, err = d.deliveryActivityRepo.CreateActivity(tx, delivery.ID, "Delivery is being delivered by "+helper.RandomDriver())
	if err != nil {
		return nil, err
	}

	_, err = d.orderRepository.UpdateOrderStatus(tx, order.ID, dto.OrderOnDelivery)
	if err != nil {
		return nil, err
	}

	return updatedDelivery, nil
}
