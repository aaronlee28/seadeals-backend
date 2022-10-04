package service

import (
	"gorm.io/gorm"
	"seadeals-backend/helper"
	"seadeals-backend/model"
	"seadeals-backend/repository"
)

type OrderService interface {
	GetOrderBySellerID(userID uint, query *repository.OrderQuery) ([]*model.Order, int64, int64, error)
}

type orderService struct {
	db               *gorm.DB
	orderRepository  repository.OrderRepository
	sellerRepository repository.SellerRepository
}

type OrderServiceConfig struct {
	DB               *gorm.DB
	OrderRepository  repository.OrderRepository
	SellerRepository repository.SellerRepository
}

func NewOrderService(c *OrderServiceConfig) OrderService {
	return &orderService{
		db:               c.DB,
		orderRepository:  c.OrderRepository,
		sellerRepository: c.SellerRepository,
	}
}

func (o *orderService) GetOrderBySellerID(userID uint, query *repository.OrderQuery) ([]*model.Order, int64, int64, error) {
	tx := o.db.Begin()
	var err error
	defer helper.CommitOrRollback(tx, &err)

	seller, err := o.sellerRepository.FindSellerByUserID(tx, userID)
	if err != nil {
		return nil, 0, 0, err
	}

	orders, totalPage, totalData, err := o.orderRepository.GetOrderBySellerID(tx, seller.ID, query)
	if err != nil {
		return nil, 0, 0, err
	}

	return orders, totalPage, totalData, nil
}
