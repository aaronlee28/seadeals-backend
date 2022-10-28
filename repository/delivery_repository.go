package repository

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"seadeals-backend/apperror"
	"seadeals-backend/config"
	"seadeals-backend/db"
	"seadeals-backend/dto"
	"seadeals-backend/model"
	"time"
)

type DeliveryRepository interface {
	GetDeliveryByOrderID(tx *gorm.DB, orderID uint) (*model.Delivery, error)

	CreateDelivery(tx *gorm.DB, delivery *model.Delivery) (*model.Delivery, error)
	UpdateDeliveryStatus(tx *gorm.DB, deliveryID uint, status string) (*model.Delivery, error)
	UpdateDeliveryStatusByOrderID(tx *gorm.DB, orderID uint, status string) (*model.Delivery, error)
	FindAndUpdateOngoingToDelivered() ([]*model.Delivery, error)
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

func (d *deliveryRepository) UpdateDeliveryStatusByOrderID(tx *gorm.DB, orderID uint, status string) (*model.Delivery, error) {
	var delivery = &model.Delivery{}
	result := tx.Model(&delivery).Clauses(clause.Returning{}).Where("order_id = ?", orderID).Update("status", status)
	if result.Error != nil {
		return nil, apperror.InternalServerError("Cannot update delivery")
	}
	return delivery, nil
}

func (d *deliveryRepository) FindAndUpdateOngoingToDelivered() ([]*model.Delivery, error) {
	tx := db.Get().Begin()
	var deliveries []*model.Delivery
	result := tx.Model(&deliveries).Clauses(clause.Returning{}).Where("status = ?", dto.DeliveryOngoing).Where("? >= updated_at at time zone '"+config.Config.TZ+"' + interval '"+config.Config.Interval.OnDeliveryToDelivered+"' * eta", time.Now()).Find(&deliveries).Update("status", dto.DeliveryDone)
	if result.Error != nil {
		return nil, apperror.InternalServerError("Cannot update delivery")
	}
	tx.Commit()

	return deliveries, nil
}
