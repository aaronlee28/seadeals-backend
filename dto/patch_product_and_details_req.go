package dto

type PatchProductAndDetailsReq struct {
	Product       *PatchProduct          `json:"product"`
	ProductDetail *PatchProductDetailReq `json:"product_details" binding:"omitempty"`
}
