package service

import (
	"bytes"
	"encoding/json"
	"gorm.io/gorm"
	"math"
	"net/http"
	"seadeals-backend/apperror"
	"seadeals-backend/config"
	"seadeals-backend/dto"
	"seadeals-backend/helper"
	"seadeals-backend/model"
	"seadeals-backend/repository"
	"strconv"
	"time"
)

type OrderService interface {
	GetOrderBySellerID(userID uint, query *repository.OrderQuery) ([]*model.Order, int64, int64, error)
	CancelOrderBySeller(orderID uint, userID uint) (*model.Order, error)
}

type orderService struct {
	db                        *gorm.DB
	orderRepository           repository.OrderRepository
	transactionRepo           repository.TransactionRepository
	voucherRepo               repository.VoucherRepository
	sellerRepository          repository.SellerRepository
	walletRepository          repository.WalletRepository
	walletTransRepo           repository.WalletTransactionRepository
	productVarDetRepo         repository.ProductVariantDetailRepository
	seaLabsPayTransHolderRepo repository.SeaLabsPayTransactionHolderRepository
}

type OrderServiceConfig struct {
	DB                        *gorm.DB
	OrderRepository           repository.OrderRepository
	SellerRepository          repository.SellerRepository
	VoucherRepo               repository.VoucherRepository
	TransactionRepo           repository.TransactionRepository
	WalletRepository          repository.WalletRepository
	WalletTransRepo           repository.WalletTransactionRepository
	ProductVarDetRepo         repository.ProductVariantDetailRepository
	SeaLabsPayTransHolderRepo repository.SeaLabsPayTransactionHolderRepository
}

func NewOrderService(c *OrderServiceConfig) OrderService {
	return &orderService{
		db:                        c.DB,
		orderRepository:           c.OrderRepository,
		sellerRepository:          c.SellerRepository,
		voucherRepo:               c.VoucherRepo,
		transactionRepo:           c.TransactionRepo,
		walletRepository:          c.WalletRepository,
		walletTransRepo:           c.WalletTransRepo,
		productVarDetRepo:         c.ProductVarDetRepo,
		seaLabsPayTransHolderRepo: c.SeaLabsPayTransHolderRepo,
	}
}

func (o *orderService) GetOrderBySellerID(userID uint, query *repository.OrderQuery) ([]*model.Order, int64, int64, error) {
	tx := o.db.Begin()
	var err error
	defer helper.CommitOrRollback(tx, &err)

	seller, err := o.sellerRepository.FindSellerByUserID(tx, userID)
	if err != nil {
		return nil, 0, 0, err
	}

	orders, totalPage, totalData, err := o.orderRepository.GetOrderBySellerID(tx, seller.ID, query)
	if err != nil {
		return nil, 0, 0, err
	}

	return orders, totalPage, totalData, nil
}

func (o *orderService) CancelOrderBySeller(orderID uint, userID uint) (*model.Order, error) {
	tx := o.db.Begin()
	var err error
	defer helper.CommitOrRollback(tx, &err)

	order, err := o.orderRepository.GetOrderDetailByID(tx, orderID)
	if err != nil {
		return nil, err
	}
	if order.Status != dto.OrderWaitingSeller {
		return nil, apperror.BadRequestError("Cannot cancel order that is currently " + order.Status)
	}

	seller, err := o.sellerRepository.FindSellerByUserID(tx, userID)
	if err != nil {
		return nil, err
	}
	if order.SellerID != seller.ID {
		return nil, apperror.BadRequestError("Cannot cancel another seller order")
	}

	var priceBeforeGlobalDisc float64
	var voucher *model.Voucher
	var amountRefunded float64
	if order.Transaction.VoucherID != nil {
		priceBeforeGlobalDisc, err = o.transactionRepo.GetPriceBeforeGlobalDisc(tx, order.TransactionID)
		if err != nil {
			return nil, err
		}
		voucher, err = o.voucherRepo.FindVoucherDetailByID(tx, *order.Transaction.VoucherID)
		if err != nil {
			return nil, err
		}
		if voucher.AmountType == "percentage" {
			amountRefunded = order.Total - ((voucher.Amount / 100) * order.Total)
		} else {
			amountReduced := (order.Total / priceBeforeGlobalDisc) * order.Total
			amountRefunded = order.Total - amountReduced
		}
	}

	var buyerWallet *model.Wallet
	var transHolder *model.SeaLabsPayTransactionHolder
	var req *http.Request
	var resp *http.Response
	if order.Transaction.PaymentMethod == dto.WALLET {
		buyerWallet, err = o.walletRepository.GetWalletByUserID(tx, order.UserID)
		if err != nil {
			return nil, err
		}

		_, err = o.walletRepository.TopUp(tx, buyerWallet, order.Total)
		if err != nil {
			return nil, err
		}

		walletTrans := &model.WalletTransaction{
			WalletID:      buyerWallet.ID,
			TransactionID: &order.TransactionID,
			Total:         math.Floor(amountRefunded),
			PaymentMethod: dto.WALLET,
			PaymentType:   "CREDIT",
			Description:   "Refund from transaction ID " + strconv.FormatUint(uint64(order.TransactionID), 10),
			CreatedAt:     time.Time{},
		}
		_, err = o.walletTransRepo.CreateTransaction(tx, walletTrans)
		if err != nil {
			return nil, err
		}
	} else if order.Transaction.PaymentMethod == dto.SEA_LABS_PAY {
		transHolder, err = o.seaLabsPayTransHolderRepo.GetTransHolderFromTransactionID(tx, order.TransactionID)
		if err != nil {
			return nil, err
		}

		client := &http.Client{}
		URL := config.Config.SeaLabsPayRefundURL
		var jsonStr = []byte(`{"reason":"Seller cancel the order", "amount":` + strconv.Itoa(int(amountRefunded)) + `, "txn_id":` + strconv.Itoa(int(transHolder.TxnID)) + `}`)

		bearer := "Bearer " + config.Config.SeaLabsPayAPIKey
		req, err = http.NewRequest("POST", URL, bytes.NewBuffer(jsonStr))
		if err != nil {
			return nil, err
		}
		req.Header.Add("Authorization", bearer)
		req.Header.Set("Content-Type", "application/json")
		resp, err = client.Do(req)
		if err != nil {
			return nil, err
		}
		if resp.StatusCode != 200 {
			type seaLabsPayError struct {
				Code    string `json:"code"`
				Message string `json:"message"`
				Data    struct {
				} `json:"data"`
			}
			var j seaLabsPayError
			err = json.NewDecoder(resp.Body).Decode(&j)
			if err != nil {
				panic(err)
			}
			return nil, apperror.BadRequestError(j.Message)
		}
	}

	for _, item := range order.OrderItems {
		_, err = o.productVarDetRepo.AddProductVariantStock(tx, item.ProductVariantDetailID, item.Quantity)
		if err != nil {
			return nil, err
		}
	}

	refundedOrder, err := o.orderRepository.UpdateOrderStatus(tx, orderID, dto.OrderRefunded)
	if err != nil {
		return nil, err
	}

	return refundedOrder, nil
}
