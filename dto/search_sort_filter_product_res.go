package dto

type SearchedSortFilterProduct struct {
	TotalLength     int                   `json:"total_page"`
	SearchedProduct []*SearchedProductRes `json:"products"`
}
