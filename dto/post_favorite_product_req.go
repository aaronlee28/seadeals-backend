package dto

type FavoriteProductReq struct {
	ProductID uint `json:"product_id" binding:"required"`
}
