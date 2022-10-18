package repository

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"seadeals-backend/apperror"
	"seadeals-backend/dto"
	"seadeals-backend/model"
)

type SellerAvailableCourierRepository interface {
	AddSellerAvailableDeliveryMethod(tx *gorm.DB, req *dto.AddDeliveryReq, sellerID uint) (*model.SellerAvailableCourier, error)
	GetAllSellerAvailableCourier(tx *gorm.DB, sellerID uint) ([]*model.SellerAvailableCourier, error)
}

type sellerAvailableCourierRepository struct{}

func NewSellerAvailableCourierRepository() SellerAvailableCourierRepository {
	return &sellerAvailableCourierRepository{}
}

func (s *sellerAvailableCourierRepository) AddSellerAvailableDeliveryMethod(tx *gorm.DB, req *dto.AddDeliveryReq, sellerID uint) (*model.SellerAvailableCourier, error) {
	var courier *model.SellerAvailableCourier
	result := tx.Model(&courier).Where("seller_id = ?", sellerID).Where("courier_id = ?", req.CourierID).First(&courier)
	if result.Error != nil {
		if result.Error != gorm.ErrRecordNotFound {
			return nil, apperror.InternalServerError("Cannot find seller available method")
		}
		courier.CourierID = req.CourierID
		courier.SellerID = sellerID
		isSelected := true
		courier.IsSelected = &isSelected
		if req.SlaDay != 0 {
			courier.SlaDay = req.SlaDay
		}
		result = tx.Model(&courier).Clauses(clause.Returning{}).Preload("Courier").Preload("Seller").Create(&courier).First(&courier)
		if result.Error != nil {
			return nil, apperror.InternalServerError("Cannot create seller available method")
		}
		return courier, nil
	}
	if req.SlaDay != 0 {
		courier.SlaDay = req.SlaDay
	}
	courier.IsSelected = req.IsSelected
	result = tx.Model(&courier).Clauses(clause.Returning{}).Preload("Courier").Preload("Seller").Updates(&courier).First(&courier)
	if result.Error != nil {
		return nil, apperror.InternalServerError("Cannot update seller available method")
	}
	return courier, nil
}

func (s *sellerAvailableCourierRepository) GetAllSellerAvailableCourier(tx *gorm.DB, sellerID uint) ([]*model.SellerAvailableCourier, error) {
	var couriers []*model.SellerAvailableCourier
	result := tx.Model(&couriers).Where("seller_id = ?", sellerID).Where("is_selected IS TRUE").Preload("Courier").Find(&couriers)
	if result.Error != nil {
		return nil, apperror.InternalServerError("Cannot find seller available method")
	}
	return couriers, nil
}
