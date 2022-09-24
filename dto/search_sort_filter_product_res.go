package dto

type SearchedSortFilterProduct struct {
	TotalLength     int                   `json:"total_data"`
	SearchedProduct []*SearchedProductRes `json:"products"`
}
