package service

import (
	"crypto/hmac"
	"crypto/sha256"
	"fmt"
	"gorm.io/gorm"
	"seadeals-backend/dto"
	"strconv"
)

type SealabsPayService interface {
	CreateSignature(req *dto.SeaDealspayReq) string
}

type sealabsPayService struct {
	db *gorm.DB
}

type SealabsServiceConfig struct {
	DB *gorm.DB
}

func NewSealabsPayService(c *SealabsServiceConfig) SealabsPayService {
	return &sealabsPayService{
		db: c.DB,
	}
}

func (s *sealabsPayService) CreateSignature(req *dto.SeaDealspayReq) string {
	amountString := strconv.Itoa(req.Amount)
	signatureString := req.CardNumber + ":" + amountString + ":sea_baiwrg"
	apiKey := "da486c418f604b17a77acc633c951158"
	h := hmac.New(sha256.New, []byte(apiKey))
	h.Write([]byte(signatureString))
	generatedSignatureString := fmt.Sprintf("%x", sha256.Sum256([]byte(signatureString)))
	return generatedSignatureString
}
