package dto

type ProductPhoto struct {
	PhotoURL string `json:"photo_url" binding:"required"`
	Name     string `json:"name" binding:"required"`
}
