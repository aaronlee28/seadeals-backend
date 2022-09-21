package dto

type SearchedSortFilterProduct struct {
	TotalLength     int                  `json:"total_length"`
	SearchedProduct []SearchedProductRes `json:"searched_product"`
}
