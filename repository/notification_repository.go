package repository

import (
	"fmt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"seadeals-backend/apperror"
	"seadeals-backend/db"
	"seadeals-backend/model"
)

type NotificationRepository interface {
	AddToNotificationFromModelForCron(newNotification *model.Notification)
	AddToNotificationFromModel(tx *gorm.DB, newNotification *model.Notification)
}

type notificationRepository struct{}

func NewNotificationRepository() NotificationRepository {
	return &notificationRepository{}
}

func (n *notificationRepository) AddToNotificationOneToOne(tx *gorm.DB, userID uint, sellerID uint, title string, detail string) error {
	newNotification := &model.Notification{
		UserID:   userID,
		SellerID: sellerID,
		Title:    title,
		Detail:   detail,
	}

	result := tx.Clauses(clause.Returning{}).Create(&newNotification)
	if result.Error != nil {
		fmt.Println(result.Error)
		return apperror.InternalServerError("notification failed to create")
	}

	fmt.Println("Notification", newNotification)
	return nil
}

func (n *notificationRepository) AddToNotificationFromModelForCron(newNotification *model.Notification) {
	tx := db.Get().Begin()
	result := tx.Clauses(clause.Returning{}).Create(&newNotification)
	if result.Error != nil {
		tx.Rollback()
		fmt.Println("error:", result.Error)
		return
	}
	tx.Commit()

	fmt.Println("Notification", newNotification)

	return
}

func (n *notificationRepository) AddToNotificationFromModel(tx *gorm.DB, newNotification *model.Notification) {

	result := tx.Clauses(clause.Returning{}).Create(&newNotification)
	if result.Error != nil {
		fmt.Println("error:", result.Error)
		return
	}

	fmt.Println("Notification", newNotification)

	return
}
