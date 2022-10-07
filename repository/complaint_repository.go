package repository

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"seadeals-backend/apperror"
	"seadeals-backend/model"
)

type ComplaintRepository interface {
	CreateComplaint(tx *gorm.DB, orderID uint, description string) (*model.Complaint, error)
}

type complaintRepository struct{}

func NewComplaintRepository() ComplaintRepository {
	return &complaintRepository{}
}

func (c *complaintRepository) CreateComplaint(tx *gorm.DB, orderID uint, description string) (*model.Complaint, error) {
	var complaint = &model.Complaint{
		OrderID:     orderID,
		Description: description,
	}
	result := tx.Model(&complaint).Clauses(clause.Returning{}).Create(&complaint)
	if result.Error != nil {
		return nil, apperror.InternalServerError("Tidak bisa memasukan data complaint")
	}
	return complaint, nil
}
