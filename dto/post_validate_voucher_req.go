package dto

type PostValidateVoucherReq struct {
	SellerID uint    `json:"seller_id" binding:"required,numeric"`
	Code     string  `json:"code" binding:"required,alphanum"`
	Price    float64 `json:"price" binding:"required,numeric,gte=0"`
	Quantity uint    `json:"quantity" binding:"required,numeric,gte=1"`
}
