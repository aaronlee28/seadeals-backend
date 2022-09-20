package dto

import (
	"seadeals-backend/apperror"
	"strconv"
)

type CartItemQuery struct {
	Limit int `json:"limit"`
}

func (_ *CartItemQuery) FromQuery(t map[string]string) (*CartItemQuery, error) {
	limit, err := strconv.ParseUint(t["limit"], 10, 32)
	if limit == 0 {
		limit = 5
	}
	if err != nil {
		return nil, apperror.BadRequestError("invalid limit format")
	}

	return &CartItemQuery{Limit: int(limit)}, nil
}
