package dto

type DeliverSettingsPrint struct {
	AllowPrint *bool `json:"allow_print" binding:"required"`
}
