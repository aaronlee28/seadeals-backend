package service_test

import (
	"github.com/stretchr/testify/assert"
	"seadeals-backend/dto"
	"seadeals-backend/service"
	"seadeals-backend/testutil"
	"testing"
)

func TestNewSealabsPayService(t *testing.T) {
	t.Run("Should return response body", func(t *testing.T) {

		gormDB := testutil.MockDB()

		cfg := &service.SealabsServiceConfig{
			DB: gormDB,
		}
		s := service.NewSealabsPayService(cfg)

		res := s.CreateSignature(&dto.SeaDealspayReq{})

		assert.NotNil(t, res)
	})
}
