package dto

type PatchProductDetailReq struct {
	Description     string `json:"description"`
	VideoURL        string `json:"video_url"`
	IsHazardous     bool   `json:"is_hazardous"`
	ConditionStatus string `json:"condition_status"`
	Length          int    `json:"length"`
	Width           int    `json:"width"`
	Height          int    `json:"height"`
	Weight          int    `json:"weight"`
}
