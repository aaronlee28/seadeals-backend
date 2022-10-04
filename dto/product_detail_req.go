package dto

type ProductDetailReq struct {
	Description     string  `json:"description" binding:"required"`
	VideoURL        string  `json:"video_url"`
	IsHazardous     *bool   `json:"is_hazardous" binding:"required"`
	ConditionStatus string  `json:"condition_status" binding:"required"`
	Length          float64 `json:"length" binding:"required"`
	Width           float64 `json:"width" binding:"required"`
	Height          float64 `json:"height" binding:"required"`
	Weight          float64 `json:"weight" binding:"required"`
}
