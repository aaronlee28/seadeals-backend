package repository

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"seadeals-backend/apperror"
	"seadeals-backend/model"
)

type DeliveryRepository interface {
	GetDeliveryByOrderID(tx *gorm.DB, orderID uint) (*model.Delivery, error)

	CreateDelivery(tx *gorm.DB, delivery *model.Delivery) (*model.Delivery, error)
	UpdateDeliveryStatus(tx *gorm.DB, deliveryID uint, status string) (*model.Delivery, error)
}

type deliveryRepository struct{}

func NewDeliveryRepository() DeliveryRepository {
	return &deliveryRepository{}
}

func (d *deliveryRepository) GetDeliveryByOrderID(tx *gorm.DB, orderID uint) (*model.Delivery, error) {
	var delivery *model.Delivery
	result := tx.Model(&delivery).Where("order_id = ?", orderID).First(&delivery)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, apperror.BadRequestError("Order doesn't exists")
		}
		return nil, apperror.InternalServerError("Cannot find delivery")
	}
	return delivery, nil
}

func (d *deliveryRepository) CreateDelivery(tx *gorm.DB, delivery *model.Delivery) (*model.Delivery, error) {
	result := tx.Model(&delivery).Clauses(clause.Returning{}).Create(&delivery)
	if result.Error != nil {
		return nil, apperror.InternalServerError("Cannot create delivery")
	}
	return delivery, nil
}

func (d *deliveryRepository) UpdateDeliveryStatus(tx *gorm.DB, deliveryID uint, status string) (*model.Delivery, error) {
	var delivery = &model.Delivery{}
	delivery.ID = deliveryID
	result := tx.Model(&delivery).Clauses(clause.Returning{}).Update("status", status)
	if result.Error != nil {
		return nil, apperror.InternalServerError("Cannot update delivery")
	}
	return delivery, nil
}
