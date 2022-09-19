package dto

import "net/http"

type AppResponse struct {
	StatusCode int    `json:"statusCode"`
	Status     string `json:"status"`
	Data       any    `json:"data"`
}

func StatusOKResponse(data any) AppResponse {
	return AppResponse{
		StatusCode: http.StatusOK,
		Status:     "OK",
		Data:       data,
	}
}

func StatusCreatedResponse(data any) AppResponse {
	return AppResponse{
		StatusCode: http.StatusCreated,
		Status:     "OK",
		Data:       data,
	}
}
