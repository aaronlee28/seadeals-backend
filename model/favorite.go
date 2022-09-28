package model

type Favorite struct {
	ID         uint     `json:"id"`
	IsFavorite bool     `json:"is_favorite"`
	UserID     uint     `json:"user_id"`
	User       *User    `json:"user"`
	ProductID  uint     `json:"product_id"`
	Product    *Product `json:"product"`
}
