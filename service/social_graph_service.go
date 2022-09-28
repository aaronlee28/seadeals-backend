package service

import (
	"gorm.io/gorm"
	"seadeals-backend/model"
	"seadeals-backend/repository"
)

type SocialGraphService interface {
	FollowToSeller(userID uint, sellerID uint) (*model.SocialGraph, error)
}

type socialGraphService struct {
	db                    *gorm.DB
	socialGraphRepository repository.SocialGraphRepository
}

type SocialGraphServiceConfig struct {
	DB                    *gorm.DB
	SocialGraphRepository repository.SocialGraphRepository
}

func NewSocialGraphService(c *SocialGraphServiceConfig) SocialGraphService {
	return &socialGraphService{
		db:                    c.DB,
		socialGraphRepository: c.SocialGraphRepository,
	}
}

func (s *socialGraphService) FollowToSeller(userID uint, sellerID uint) (*model.SocialGraph, error) {
	tx := s.db.Begin()

	favorite, err := s.socialGraphRepository.FollowToSeller(tx, userID, sellerID)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	tx.Commit()
	return favorite, nil
}
