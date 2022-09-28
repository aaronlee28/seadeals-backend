package service

import (
	"gorm.io/gorm"
	"seadeals-backend/dto"
	"seadeals-backend/repository"
)

type PromotionService interface {
	GetPromotionByID(id uint) (*dto.GetPromotionRes, error)
}

type promotionService struct {
	db                  *gorm.DB
	promotionRepository repository.PromotionRepository
}

type PromotionServiceConfig struct {
	DB                  *gorm.DB
	PromotionRepository repository.PromotionRepository
}

func NewPromotionService(c *PromotionServiceConfig) PromotionService {
	return &promotionService{
		db:                  c.DB,
		promotionRepository: c.PromotionRepository,
	}
}

func (p *promotionService) GetPromotionByID(id uint) (*dto.GetPromotionRes, error) {
	tx := p.db.Begin()
	promotion, err := Ï
}
