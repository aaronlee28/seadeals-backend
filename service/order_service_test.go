package service_test

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"seadeals-backend/mocks"
	"seadeals-backend/model"
	"seadeals-backend/service"
	"seadeals-backend/testutil"
	"testing"
)

func TestOrderService_GetDetailOrderForReceipt(t *testing.T) {
	t.Run("Should return response body", func(t *testing.T) {

		gormDB := testutil.MockDB()
		mockRepo1 := new(mocks.AccountHolderRepository)
		mockRepo2 := new(mocks.AddressRepository)
		mockRepo3 := new(mocks.OrderRepository)
		mockRepo4 := new(mocks.CourierRepository)
		mockRepo5 := new(mocks.TransactionRepository)
		mockRepo6 := new(mocks.VoucherRepository)
		mockRepo7 := new(mocks.DeliveryRepository)
		mockRepo8 := new(mocks.SellerRepository)
		mockRepo9 := new(mocks.WalletRepository)
		mockRepo10 := new(mocks.WalletTransactionRepository)
		mockRepo11 := new(mocks.ProductVariantDetailRepository)
		mockRepo12 := new(mocks.ProductRepository)
		mockRepo13 := new(mocks.SeaLabsPayTransactionHolderRepository)
		mockRepo14 := new(mocks.ComplaintRepository)
		mockRepo15 := new(mocks.ComplaintPhotoRepository)
		mockRepo16 := new(mocks.NotificationRepository)
		cfg := &service.OrderServiceConfig{
			DB:                        gormDB,
			AccountHolderRepo:         mockRepo1,
			AddressRepository:         mockRepo2,
			OrderRepository:           mockRepo3,
			CourierRepository:         mockRepo4,
			SellerRepository:          mockRepo8,
			VoucherRepo:               mockRepo6,
			DeliveryRepo:              mockRepo7,
			TransactionRepo:           mockRepo5,
			WalletRepository:          mockRepo9,
			WalletTransRepo:           mockRepo10,
			ProductVarDetRepo:         mockRepo11,
			ProductRepo:               mockRepo12,
			SeaLabsPayTransHolderRepo: mockRepo13,
			ComplainRepo:              mockRepo14,
			ComplaintPhotoRepo:        mockRepo15,
			NotificationRepo:          mockRepo16,
		}

		s := service.NewOrderService(cfg)

		mockVoucher := &model.Voucher{AmountType: "quantity"}
		mockDelivery2 := &model.Delivery{Total: 1}
		mockSeller := &model.Seller{Name: ""}
		mockOrder := &model.Order{Total: 1, Delivery: mockDelivery2, Seller: mockSeller}
		mockOrders := []*model.Order{mockOrder}
		mockTransaction := &model.Transaction{Voucher: mockVoucher, Orders: mockOrders}
		mockCourier := &model.Courier{Name: ""}
		mockDelivery := &model.Delivery{Total: 1, Courier: mockCourier}
		mockSellerName := &model.Seller{Name: ""}
		mockUser := &model.User{FullName: ""}
		mockPD := &model.ProductDetail{Weight: 1}
		mockProduct := &model.Product{Name: "", ProductDetail: mockPD}
		mockPV1 := &model.ProductVariant{}
		mockPV2 := &model.ProductVariant{}
		Pv1Val := "1"
		Pv2Val := "1"
		mockPVD := &model.ProductVariantDetail{Product: mockProduct, ProductVariant1: mockPV1, ProductVariant2: mockPV2, Variant1Value: &Pv1Val, Variant2Value: &Pv2Val}
		mockOrderItems := &model.OrderItem{ProductVariantDetail: mockPVD}
		mockOrderItemsArr := []*model.OrderItem{mockOrderItems}

		mockRepo3.On("GetOrderDetailForReceipt", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("uint")).Return(&model.Order{UserID: 1, Transaction: mockTransaction, Total: 1, Delivery: mockDelivery, Seller: mockSellerName, User: mockUser, OrderItems: mockOrderItemsArr, Voucher: mockVoucher}, nil)

		mockRepo5.On("GetPriceBeforeGlobalDisc", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("uint")).Return(float64(1), nil)

		res, err := s.GetDetailOrderForReceipt(uint(1), uint(1))

		assert.Nil(t, err)
		assert.NotNil(t, res)
	})

	t.Run("Should return response body", func(t *testing.T) {

		gormDB := testutil.MockDB()
		mockRepo1 := new(mocks.AccountHolderRepository)
		mockRepo2 := new(mocks.AddressRepository)
		mockRepo3 := new(mocks.OrderRepository)
		mockRepo4 := new(mocks.CourierRepository)
		mockRepo5 := new(mocks.TransactionRepository)
		mockRepo6 := new(mocks.VoucherRepository)
		mockRepo7 := new(mocks.DeliveryRepository)
		mockRepo8 := new(mocks.SellerRepository)
		mockRepo9 := new(mocks.WalletRepository)
		mockRepo10 := new(mocks.WalletTransactionRepository)
		mockRepo11 := new(mocks.ProductVariantDetailRepository)
		mockRepo12 := new(mocks.ProductRepository)
		mockRepo13 := new(mocks.SeaLabsPayTransactionHolderRepository)
		mockRepo14 := new(mocks.ComplaintRepository)
		mockRepo15 := new(mocks.ComplaintPhotoRepository)
		mockRepo16 := new(mocks.NotificationRepository)
		cfg := &service.OrderServiceConfig{
			DB:                        gormDB,
			AccountHolderRepo:         mockRepo1,
			AddressRepository:         mockRepo2,
			OrderRepository:           mockRepo3,
			CourierRepository:         mockRepo4,
			SellerRepository:          mockRepo8,
			VoucherRepo:               mockRepo6,
			DeliveryRepo:              mockRepo7,
			TransactionRepo:           mockRepo5,
			WalletRepository:          mockRepo9,
			WalletTransRepo:           mockRepo10,
			ProductVarDetRepo:         mockRepo11,
			ProductRepo:               mockRepo12,
			SeaLabsPayTransHolderRepo: mockRepo13,
			ComplainRepo:              mockRepo14,
			ComplaintPhotoRepo:        mockRepo15,
			NotificationRepo:          mockRepo16,
		}

		s := service.NewOrderService(cfg)

		mockVoucher := &model.Voucher{AmountType: "percentage", Amount: 1}
		mockDelivery2 := &model.Delivery{Total: 1}
		mockSeller := &model.Seller{Name: ""}
		mockVoucher2 := &model.Voucher{AmountType: "percentage", Amount: 1}
		voucherID := uint(2)
		mockTrans := &model.Transaction{VoucherID: &voucherID, Voucher: mockVoucher2}
		mockOrder := &model.Order{Total: 1, Delivery: mockDelivery2, Seller: mockSeller, Transaction: mockTrans, Voucher: mockVoucher2}
		mockOrders := []*model.Order{mockOrder}
		mockTransaction := &model.Transaction{Voucher: mockVoucher, VoucherID: &voucherID, Orders: mockOrders}
		mockCourier := &model.Courier{Name: ""}
		mockDelivery := &model.Delivery{Total: 1, Courier: mockCourier}
		mockSellerName := &model.Seller{Name: ""}
		mockUser := &model.User{FullName: ""}
		mockPD := &model.ProductDetail{Weight: 1}
		mockProduct := &model.Product{Name: "", ProductDetail: mockPD}
		mockPV1 := &model.ProductVariant{}
		mockPV2 := &model.ProductVariant{}
		Pv1Val := "1"
		Pv2Val := "1"
		mockPVD := &model.ProductVariantDetail{Product: mockProduct, ProductVariant1: mockPV1, ProductVariant2: mockPV2, Variant1Value: &Pv1Val, Variant2Value: &Pv2Val}
		mockOrderItems := &model.OrderItem{ProductVariantDetail: mockPVD}
		mockOrderItemsArr := []*model.OrderItem{mockOrderItems}

		mockRepo3.On("GetOrderDetailForReceipt", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("uint")).Return(&model.Order{UserID: 1, Transaction: mockTransaction, Total: 1, Delivery: mockDelivery, Seller: mockSellerName, User: mockUser, OrderItems: mockOrderItemsArr, Voucher: mockVoucher}, nil)

		mockRepo5.On("GetPriceBeforeGlobalDisc", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("uint")).Return(float64(1), nil)

		res, err := s.GetDetailOrderForReceipt(uint(1), uint(1))

		assert.Nil(t, err)
		assert.NotNil(t, res)
	})

	t.Run("Should return error", func(t *testing.T) {

		gormDB := testutil.MockDB()
		mockRepo1 := new(mocks.AccountHolderRepository)
		mockRepo2 := new(mocks.AddressRepository)
		mockRepo3 := new(mocks.OrderRepository)
		mockRepo4 := new(mocks.CourierRepository)
		mockRepo5 := new(mocks.TransactionRepository)
		mockRepo6 := new(mocks.VoucherRepository)
		mockRepo7 := new(mocks.DeliveryRepository)
		mockRepo8 := new(mocks.SellerRepository)
		mockRepo9 := new(mocks.WalletRepository)
		mockRepo10 := new(mocks.WalletTransactionRepository)
		mockRepo11 := new(mocks.ProductVariantDetailRepository)
		mockRepo12 := new(mocks.ProductRepository)
		mockRepo13 := new(mocks.SeaLabsPayTransactionHolderRepository)
		mockRepo14 := new(mocks.ComplaintRepository)
		mockRepo15 := new(mocks.ComplaintPhotoRepository)
		mockRepo16 := new(mocks.NotificationRepository)
		cfg := &service.OrderServiceConfig{
			DB:                        gormDB,
			AccountHolderRepo:         mockRepo1,
			AddressRepository:         mockRepo2,
			OrderRepository:           mockRepo3,
			CourierRepository:         mockRepo4,
			SellerRepository:          mockRepo8,
			VoucherRepo:               mockRepo6,
			DeliveryRepo:              mockRepo7,
			TransactionRepo:           mockRepo5,
			WalletRepository:          mockRepo9,
			WalletTransRepo:           mockRepo10,
			ProductVarDetRepo:         mockRepo11,
			ProductRepo:               mockRepo12,
			SeaLabsPayTransHolderRepo: mockRepo13,
			ComplainRepo:              mockRepo14,
			ComplaintPhotoRepo:        mockRepo15,
			NotificationRepo:          mockRepo16,
		}

		s := service.NewOrderService(cfg)

		mockRepo3.On("GetOrderDetailForReceipt", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("uint")).Return(nil, errors.New(""))

		mockRepo5.On("GetPriceBeforeGlobalDisc", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("uint")).Return(float64(1), nil)

		res, err := s.GetDetailOrderForReceipt(uint(1), uint(1))

		assert.Nil(t, res)
		assert.NotNil(t, err)
	})
	t.Run("Should return error", func(t *testing.T) {

		gormDB := testutil.MockDB()
		mockRepo1 := new(mocks.AccountHolderRepository)
		mockRepo2 := new(mocks.AddressRepository)
		mockRepo3 := new(mocks.OrderRepository)
		mockRepo4 := new(mocks.CourierRepository)
		mockRepo5 := new(mocks.TransactionRepository)
		mockRepo6 := new(mocks.VoucherRepository)
		mockRepo7 := new(mocks.DeliveryRepository)
		mockRepo8 := new(mocks.SellerRepository)
		mockRepo9 := new(mocks.WalletRepository)
		mockRepo10 := new(mocks.WalletTransactionRepository)
		mockRepo11 := new(mocks.ProductVariantDetailRepository)
		mockRepo12 := new(mocks.ProductRepository)
		mockRepo13 := new(mocks.SeaLabsPayTransactionHolderRepository)
		mockRepo14 := new(mocks.ComplaintRepository)
		mockRepo15 := new(mocks.ComplaintPhotoRepository)
		mockRepo16 := new(mocks.NotificationRepository)
		cfg := &service.OrderServiceConfig{
			DB:                        gormDB,
			AccountHolderRepo:         mockRepo1,
			AddressRepository:         mockRepo2,
			OrderRepository:           mockRepo3,
			CourierRepository:         mockRepo4,
			SellerRepository:          mockRepo8,
			VoucherRepo:               mockRepo6,
			DeliveryRepo:              mockRepo7,
			TransactionRepo:           mockRepo5,
			WalletRepository:          mockRepo9,
			WalletTransRepo:           mockRepo10,
			ProductVarDetRepo:         mockRepo11,
			ProductRepo:               mockRepo12,
			SeaLabsPayTransHolderRepo: mockRepo13,
			ComplainRepo:              mockRepo14,
			ComplaintPhotoRepo:        mockRepo15,
			NotificationRepo:          mockRepo16,
		}

		s := service.NewOrderService(cfg)

		mockVoucher := &model.Voucher{AmountType: "quantity"}
		mockDelivery2 := &model.Delivery{Total: 1}
		mockSeller := &model.Seller{Name: ""}
		mockOrder := &model.Order{Total: 1, Delivery: mockDelivery2, Seller: mockSeller}
		mockOrders := []*model.Order{mockOrder}
		mockTransaction := &model.Transaction{Voucher: mockVoucher, Orders: mockOrders}
		mockCourier := &model.Courier{Name: ""}
		mockDelivery := &model.Delivery{Total: 1, Courier: mockCourier}
		mockSellerName := &model.Seller{Name: ""}
		mockUser := &model.User{FullName: ""}
		mockPD := &model.ProductDetail{Weight: 1}
		mockProduct := &model.Product{Name: "", ProductDetail: mockPD}
		mockPV1 := &model.ProductVariant{}
		mockPV2 := &model.ProductVariant{}
		Pv1Val := "1"
		Pv2Val := "1"
		mockPVD := &model.ProductVariantDetail{Product: mockProduct, ProductVariant1: mockPV1, ProductVariant2: mockPV2, Variant1Value: &Pv1Val, Variant2Value: &Pv2Val}
		mockOrderItems := &model.OrderItem{ProductVariantDetail: mockPVD}
		mockOrderItemsArr := []*model.OrderItem{mockOrderItems}

		mockRepo3.On("GetOrderDetailForReceipt", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("uint")).Return(&model.Order{UserID: 1, Transaction: mockTransaction, Total: 1, Delivery: mockDelivery, Seller: mockSellerName, User: mockUser, OrderItems: mockOrderItemsArr, Voucher: mockVoucher}, nil)

		mockRepo5.On("GetPriceBeforeGlobalDisc", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("uint")).Return(float64(1), nil)

		res, err := s.GetDetailOrderForReceipt(uint(1), uint(2))

		assert.Nil(t, res)
		assert.NotNil(t, err)
	})
	t.Run("Should return error", func(t *testing.T) {

		gormDB := testutil.MockDB()
		mockRepo1 := new(mocks.AccountHolderRepository)
		mockRepo2 := new(mocks.AddressRepository)
		mockRepo3 := new(mocks.OrderRepository)
		mockRepo4 := new(mocks.CourierRepository)
		mockRepo5 := new(mocks.TransactionRepository)
		mockRepo6 := new(mocks.VoucherRepository)
		mockRepo7 := new(mocks.DeliveryRepository)
		mockRepo8 := new(mocks.SellerRepository)
		mockRepo9 := new(mocks.WalletRepository)
		mockRepo10 := new(mocks.WalletTransactionRepository)
		mockRepo11 := new(mocks.ProductVariantDetailRepository)
		mockRepo12 := new(mocks.ProductRepository)
		mockRepo13 := new(mocks.SeaLabsPayTransactionHolderRepository)
		mockRepo14 := new(mocks.ComplaintRepository)
		mockRepo15 := new(mocks.ComplaintPhotoRepository)
		mockRepo16 := new(mocks.NotificationRepository)
		cfg := &service.OrderServiceConfig{
			DB:                        gormDB,
			AccountHolderRepo:         mockRepo1,
			AddressRepository:         mockRepo2,
			OrderRepository:           mockRepo3,
			CourierRepository:         mockRepo4,
			SellerRepository:          mockRepo8,
			VoucherRepo:               mockRepo6,
			DeliveryRepo:              mockRepo7,
			TransactionRepo:           mockRepo5,
			WalletRepository:          mockRepo9,
			WalletTransRepo:           mockRepo10,
			ProductVarDetRepo:         mockRepo11,
			ProductRepo:               mockRepo12,
			SeaLabsPayTransHolderRepo: mockRepo13,
			ComplainRepo:              mockRepo14,
			ComplaintPhotoRepo:        mockRepo15,
			NotificationRepo:          mockRepo16,
		}

		s := service.NewOrderService(cfg)

		mockVoucher := &model.Voucher{AmountType: "quantity"}
		mockDelivery2 := &model.Delivery{Total: 1}
		mockSeller := &model.Seller{Name: ""}
		mockOrder := &model.Order{Total: 1, Delivery: mockDelivery2, Seller: mockSeller}
		mockOrders := []*model.Order{mockOrder}
		mockTransaction := &model.Transaction{Voucher: mockVoucher, Orders: mockOrders}
		mockCourier := &model.Courier{Name: ""}
		mockDelivery := &model.Delivery{Total: 1, Courier: mockCourier}
		mockSellerName := &model.Seller{Name: ""}
		mockUser := &model.User{FullName: ""}
		mockPD := &model.ProductDetail{Weight: 1}
		mockProduct := &model.Product{Name: "", ProductDetail: mockPD}
		mockPV1 := &model.ProductVariant{}
		mockPV2 := &model.ProductVariant{}
		Pv1Val := "1"
		Pv2Val := "1"
		mockPVD := &model.ProductVariantDetail{Product: mockProduct, ProductVariant1: mockPV1, ProductVariant2: mockPV2, Variant1Value: &Pv1Val, Variant2Value: &Pv2Val}
		mockOrderItems := &model.OrderItem{ProductVariantDetail: mockPVD}
		mockOrderItemsArr := []*model.OrderItem{mockOrderItems}

		mockRepo3.On("GetOrderDetailForReceipt", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("uint")).Return(&model.Order{UserID: 1, Transaction: mockTransaction, Total: 1, Delivery: mockDelivery, Seller: mockSellerName, User: mockUser, OrderItems: mockOrderItemsArr, Voucher: mockVoucher}, nil)

		mockRepo5.On("GetPriceBeforeGlobalDisc", mock.AnythingOfType(testutil.GormDBPointerType), mock.AnythingOfType("uint")).Return(float64(0), errors.New(""))

		res, err := s.GetDetailOrderForReceipt(uint(1), uint(1))

		assert.Nil(t, res)
		assert.NotNil(t, err)
	})

}
