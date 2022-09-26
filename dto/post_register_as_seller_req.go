package dto

type RegisterAsSellerReq struct {
	UserID      uint   `json:"user_id" binding:"required"`
	ShopName    string `json:"shop_name" binding:"required"`
	Description string `json:"description" binding:"required"`
}
