package repository

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"seadeals-backend/apperror"
	"seadeals-backend/model"
)

type ComplaintPhotoRepository interface {
	CreateComplaintPhotos(tx *gorm.DB, photos []*model.ComplaintPhoto) ([]*model.ComplaintPhoto, error)
}

type complaintPhotoRepository struct{}

func NewComplaintPhotoRepository() ComplaintPhotoRepository {
	return &complaintPhotoRepository{}
}

func (c *complaintPhotoRepository) CreateComplaintPhotos(tx *gorm.DB, photos []*model.ComplaintPhoto) ([]*model.ComplaintPhoto, error) {
	result := tx.Model(&photos).Clauses(clause.Returning{}).Create(&photos)
	if result.Error != nil {
		return nil, apperror.InternalServerError("Tidak bisa memasukan data complaint's photos")
	}
	return photos, nil
}
