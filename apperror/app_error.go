package apperror

import (
	"net/http"
)

type AppError struct {
	StatusCode int    `json:"statusCode"`
	Status     string `json:"status"`
	Message    string `json:"message"`
}

func (err AppError) Error() string {
	return err.Message
}

func BadRequestError(message string) AppError {
	return AppError{
		Status:     "BAD_REQUEST_ERROR",
		Message:    message,
		StatusCode: http.StatusBadRequest,
	}
}

func ForbiddenError(message string) AppError {
	return AppError{
		Status:     "FORBIDDEN_ERROR",
		Message:    message,
		StatusCode: http.StatusForbidden,
	}
}

func NotFoundError(message string) AppError {
	return AppError{
		Status:     "NOT_FOUND_ERROR",
		Message:    message,
		StatusCode: http.StatusNotFound,
	}
}

func UnauthorizedError(message string) AppError {
	return AppError{
		Status:     "UNAUTHORIZED_ERROR",
		Message:    message,
		StatusCode: http.StatusUnauthorized,
	}
}

func InternalServerError(message string) AppError {
	return AppError{
		Status:     "INTERNAL_SERVER_ERROR",
		Message:    message,
		StatusCode: http.StatusInternalServerError,
	}
}
