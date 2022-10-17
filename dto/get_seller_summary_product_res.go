package dto

import (
	"seadeals-backend/model"
)

type GetSellerSummaryProductRes struct {
	ID                   uint                    `json:"id"`
	CategoryName         string                  `json:"category"`
	Name                 string                  `json:"name"`
	Slug                 string                  `json:"slug"`
	SoldCount            int                     `json:"sold_count"`
	FavoriteCount        uint                    `json:"favorite_count"`
	IsArchived           bool                    `json:"is_archived"`
	ProductVariantDetail []*GetProductVariantRes `json:"product_variant_detail"`
	Photo                string                  `json:"photo"`
	IsDeleted            bool                    `json:"is_deleted"`
}

func (_ *GetSellerSummaryProductRes) From(p *model.Product) *GetSellerSummaryProductRes {
	var pvd []*GetProductVariantRes
	for _, detail := range p.ProductVariantDetail {
		pvd = append(pvd, new(GetProductVariantRes).From(detail))
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
		SoldCount:            p.SoldCount,
		FavoriteCount:        p.FavoriteCount,
		IsArchived:           p.IsArchived,
		ProductVariantDetail: pvd,
		Photo:                photoURL,
		IsDeleted:            p.DeletedAt.Valid,
	}
}
