package dto

type ProductPhotoReq struct {
	ProductPhoto []ProductPhoto `json:"product_photo" binding:"required"`
}
