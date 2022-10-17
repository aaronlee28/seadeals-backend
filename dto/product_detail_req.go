package dto

type ProductDetailsReq struct {
	Description     string `json:"description" binding:"required"`
	VideoURL        string `json:"video_url"`
	IsHazardous     *bool  `json:"is_hazardous" binding:"required"`
	ConditionStatus string `json:"condition_status" binding:"required"`
	Length          int    `json:"length" binding:"required"`
	Width           int    `json:"width" binding:"required"`
	Height          int    `json:"height" binding:"required"`
	Weight          int    `json:"weight" binding:"required"`
}
