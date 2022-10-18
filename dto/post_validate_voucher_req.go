package dto

type PostValidateVoucherReq struct {
	SellerID uint    `json:"seller_id" binding:"required,numeric"`
	Code     string  `json:"code" binding:"required,alphanum"`
	Spend    float64 `json:"spend" binding:"required,numeric,gte=0"`
}
