package dto

type PatchProduct struct {
	Name          string `json:"name"`
	CategoryID    uint   `json:"category_id"`
	IsBulkEnabled bool   `json:"is_bulk_enabled"`
	MinQuantity   uint   `json:"min_quantity"`
	MaxQuantity   uint   `json:"max_quantity"`
}
