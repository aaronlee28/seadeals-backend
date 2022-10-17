package dto

type CreateUpdateReview struct {
	ProductID   uint    `json:"product_id" binding:"required"`
	Rating      uint    `json:"rating" binding:"required"`
	Description *string `json:"description" binding:"required"`
	ImageURL    *string `json:"image_url" binding:"required"`
	ImageName   *string `json:"image_name" binding:"required"`
}
