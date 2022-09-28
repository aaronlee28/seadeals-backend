package dto

type FollowSellerReq struct {
	SellerID uint `json:"seller_id" binding:"required"`
}
