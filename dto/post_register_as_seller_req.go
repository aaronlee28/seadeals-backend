package dto

type RegisterAsSellerReq struct {
	ShopName    string `json:"shop_name" binding:"required"`
	Description string `json:"description" binding:"required"`
}
