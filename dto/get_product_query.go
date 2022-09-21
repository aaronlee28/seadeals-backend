package dto

import (
	"strconv"
)

type SellerProductSearchQuery struct {
	SortBy string `json:"sortBy"`
	Sort   string `json:"sort"`
	Search string `json:"search"`
	Page   int    `json:"page"`
	Limit  int    `json:"limit"`
}

func (_ *SellerProductSearchQuery) FromQuery(t map[string]string) (*SellerProductSearchQuery, error) {
	limit, _ := strconv.ParseUint(t["limit"], 10, 32)
	if limit == 0 {
		limit = 20
	}

	page, _ := strconv.ParseUint(t["page"], 10, 32)
	if page == 0 {
		page = 1
	}

	return &SellerProductSearchQuery{
		Search: t["s"],
		SortBy: t["sortBy"],
		Sort:   t["sort"],
		Limit:  int(limit),
		Page:   int(page),
	}, nil
}
