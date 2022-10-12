package dto

import "seadeals-backend/model"

type GetSellerSummaryProductRes struct {
	ID                   uint                    `json:"id"`
	CategoryName         string                  `json:"category"`
	Name                 string                  `json:"name"`
	Slug                 string                  `json:"slug"`
	IsBulkEnabled        bool                    `json:"is_bulk_enabled"`
	SoldCount            int                     `json:"sold_count"`
	FavoriteCount        uint                    `json:"favorite_count"`
	IsArchived           bool                    `json:"is_archived"`
	ProductVariantDetail []*ProductVariantDetail `json:"product_variant_detail"`
	Photo                string                  `json:"photo"`
	Rating               float64                 `json:"rating"`
	TotalReview          int64                   `json:"total_review"`
}

func (_ *GetSellerSummaryProductRes) From(p *model.Product, rating float64, totalReview int64) *GetSellerSummaryProductRes {
	var pvd []*ProductVariantDetail
	for _, detail := range p.ProductVariantDetail {
		pvd = append(pvd, new(ProductVariantDetail).From(detail))
	}

	var photoURL string
	if len(p.ProductPhotos) > 0 {
		photoURL = p.ProductPhotos[0].PhotoURL
	}

	return &GetSellerSummaryProductRes{
		ID:                   p.ID,
		CategoryName:         p.Category.Name,
		Name:                 p.Name,
		Slug:                 p.Slug,
		IsBulkEnabled:        p.IsBulkEnabled,
		SoldCount:            p.SoldCount,
		FavoriteCount:        p.FavoriteCount,
		IsArchived:           p.IsArchived,
		ProductVariantDetail: pvd,
		Photo:                photoURL,
		Rating:               rating,
		TotalReview:          totalReview,
	}
}
