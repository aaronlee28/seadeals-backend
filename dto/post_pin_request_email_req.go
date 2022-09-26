package dto

type KeyRequestByEmailReq struct {
	Key string `json:"key" binding:"required"`
}

type CodeKeyRequestByEmailReq struct {
	Key  string `json:"key" binding:"required"`
	Code string `json:"code" binding:"required"`
}

type ChangePinByEmailReq struct {
	Key  string `json:"key" binding:"required"`
	Code string `json:"code" binding:"required"`
	Pin  string `json:"pin" binding:"required"`
}
