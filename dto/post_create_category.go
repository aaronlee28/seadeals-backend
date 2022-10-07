package dto

type CreateCategory struct {
	Name     string `json:"name" binding:"required"`
	Slug     string `json:"slug" binding:"required"`
	IconURL  string `json:"icon_url" binding:"required"`
	ParentID *uint  `json:"parent_id"`
}
