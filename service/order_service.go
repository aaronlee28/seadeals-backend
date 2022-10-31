package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/robfig/cron/v3"
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
	GetDetailOrderForReceipt(orderID uint, userID uint) (*dto.Receipt, error)
	GetDetailOrderForThermal(orderID uint, userID uint) (*dto.Thermal, error)
	GetOrderBySellerID(userID uint, query *repository.OrderQuery) ([]*dto.OrderListRes, int64, int64, error)
	GetOrderByUserID(userID uint, query *repository.OrderQuery) ([]*dto.OrderListRes, int64, int64, error)

	CancelOrderBySeller(orderID uint, userID uint) (*model.Order, error)
	RequestRefundByBuyer(req *dto.CreateComplaintReq, userID uint) (*dto.CreateComplaintRes, error)
	AcceptRefundRequest(req *dto.RejectAcceptRefundReq, userID uint) (*dto.RejectAcceptRefundRes, error)
	RejectRefundRequest(req *dto.RejectAcceptRefundReq, userID uint) (*dto.RejectAcceptRefundRes, error)
	FinishOrder(req *dto.FinishOrderReq, userID uint) (*model.Order, error)

	RunCronJobs()
	GetTotalPredictedPrice(req *dto.PredictedPriceReq, userID uint) (*dto.TotalPredictedPriceRes, error)
	GetOrderByID(userID uint, orderID uint) (*dto.OrderListRes, error)
}

type orderService struct {
	db                        *gorm.DB
	accountHolderRepo         repository.AccountHolderRepository
	addressRepository         repository.AddressRepository
	orderRepository           repository.OrderRepository
	courierRepository         repository.CourierRepository
	transactionRepo           repository.TransactionRepository
	voucherRepo               repository.VoucherRepository
	deliveryRepo              repository.DeliveryRepository
	sellerRepository          repository.SellerRepository
	walletRepository          repository.WalletRepository
	walletTransRepo           repository.WalletTransactionRepository
	productVarDetRepo         repository.ProductVariantDetailRepository
	productRepo               repository.ProductRepository
	seaLabsPayTransHolderRepo repository.SeaLabsPayTransactionHolderRepository
	complaintRepo             repository.ComplaintRepository
	complaintPhotoRepo        repository.ComplaintPhotoRepository
	notificationRepo          repository.NotificationRepository
}

type OrderServiceConfig struct {
	DB                        *gorm.DB
	AccountHolderRepo         repository.AccountHolderRepository
	AddressRepository         repository.AddressRepository
	OrderRepository           repository.OrderRepository
	CourierRepository         repository.CourierRepository
	SellerRepository          repository.SellerRepository
	VoucherRepo               repository.VoucherRepository
	DeliveryRepo              repository.DeliveryRepository
	TransactionRepo           repository.TransactionRepository
	WalletRepository          repository.WalletRepository
	WalletTransRepo           repository.WalletTransactionRepository
	ProductVarDetRepo         repository.ProductVariantDetailRepository
	ProductRepo               repository.ProductRepository
	SeaLabsPayTransHolderRepo repository.SeaLabsPayTransactionHolderRepository
	ComplainRepo              repository.ComplaintRepository
	ComplaintPhotoRepo        repository.ComplaintPhotoRepository
	NotificationRepo          repository.NotificationRepository
}

func NewOrderService(c *OrderServiceConfig) OrderService {
	return &orderService{
		db:                        c.DB,
		accountHolderRepo:         c.AccountHolderRepo,
		addressRepository:         c.AddressRepository,
		orderRepository:           c.OrderRepository,
		courierRepository:         c.CourierRepository,
		sellerRepository:          c.SellerRepository,
		voucherRepo:               c.VoucherRepo,
		deliveryRepo:              c.DeliveryRepo,
		transactionRepo:           c.TransactionRepo,
		walletRepository:          c.WalletRepository,
		walletTransRepo:           c.WalletTransRepo,
		productVarDetRepo:         c.ProductVarDetRepo,
		productRepo:               c.ProductRepo,
		seaLabsPayTransHolderRepo: c.SeaLabsPayTransHolderRepo,
		complaintRepo:             c.ComplainRepo,
		complaintPhotoRepo:        c.ComplaintPhotoRepo,
		notificationRepo:          c.NotificationRepo,
	}
}

func refundMoneyToSeaLabsPay(URL string, jsonStr []byte) error {
	client := &http.Client{}
	bearer := "Bearer " + config.Config.SeaLabsPayAPIKey
	httpReq, err := http.NewRequest("POST", URL, bytes.NewBuffer(jsonStr))
	if err != nil {
		return err
	}
	httpReq.Header.Add("Authorization", bearer)
	httpReq.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(httpReq)
	if err != nil {
		return err
	}
	fmt.Println("CODE:", resp.StatusCode)
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
			fmt.Println(j)
		}
		return apperror.BadRequestError(j.Message)
	}
	return nil
}

func (o *orderService) GetDetailOrderForReceipt(orderID uint, userID uint) (*dto.Receipt, error) {
	tx := o.db.Begin()
	var err error
	defer helper.CommitOrRollback(tx, &err)

	order, err := o.orderRepository.GetOrderDetailForReceipt(tx, orderID)
	if err != nil {
		return nil, err
	}

	if order.UserID != userID {
		err = apperror.UnauthorizedError("Tidak bisa melihat invoice user lain")
		return nil, err
	}

	totalPriceBeforeDisc, err := o.transactionRepo.GetPriceBeforeGlobalDisc(tx, order.TransactionID)
	if err != nil {
		return nil, err
	}

	var totalQuantity uint
	var totalOrderBeforeDisc float64
	var orderItems []*dto.OrderItemReceipt
	for _, item := range order.OrderItems {
		totalQuantity += item.Quantity
		totalOrderBeforeDisc += math.Floor(item.Subtotal)

		var variantDetail string
		if item.ProductVariantDetail.ProductVariant1 != nil {
			variantDetail += *item.ProductVariantDetail.Variant1Value
		}
		if item.ProductVariantDetail.ProductVariant2 != nil {
			variantDetail += ", " + *item.ProductVariantDetail.Variant2Value
		}

		orderItem := &dto.OrderItemReceipt{
			Name:         item.ProductVariantDetail.Product.Name,
			Weight:       uint(item.ProductVariantDetail.Product.ProductDetail.Weight),
			Quantity:     item.Quantity,
			PricePerItem: math.Floor(item.ProductVariantDetail.Price),
			Discount:     math.Floor(item.ProductVariantDetail.Price*float64(item.Quantity) - item.Subtotal),
			Subtotal:     math.Floor(item.Subtotal),
			Variant:      variantDetail,
		}
		orderItems = append(orderItems, orderItem)
	}

	var voucher = &dto.ShopVoucherReceipt{}
	if order.Voucher != nil {
		totalReduced := totalOrderBeforeDisc - order.Total
		if totalReduced < 0 {
			totalReduced = 0
		}
		voucher = &dto.ShopVoucherReceipt{
			Type:        order.Voucher.AmountType,
			Name:        order.Voucher.Name,
			Amount:      order.Voucher.Amount,
			TotalReduce: math.Floor(totalReduced),
		}
	}

	var globalVoucherForOrder = &dto.GlobalVoucherForOrderReceipt{}
	if order.Transaction.Voucher != nil {
		var totalReduced float64
		var globalVoucher = order.Transaction.Voucher
		if globalVoucher.AmountType == "percentage" {
			totalReduced = (globalVoucher.Amount / 100) * order.Total
		} else {
			totalReduced = (order.Total / totalPriceBeforeDisc) * globalVoucher.Amount
			fmt.Println("order Total : ", order.Total)
			fmt.Println("total transaction : ", totalPriceBeforeDisc)
			fmt.Println((order.Total / totalPriceBeforeDisc) * order.Total)

		}
		globalVoucherForOrder = &dto.GlobalVoucherForOrderReceipt{
			Type:        order.Transaction.Voucher.AmountType,
			Name:        order.Transaction.Voucher.Name,
			Amount:      order.Transaction.Voucher.Amount,
			TotalReduce: math.Floor(totalReduced),
		}
	}

	var globalVouchers []*dto.GlobalDiscountReceipt
	var orderPayments []*dto.OrderPaymentReceipt
	var totalTransaction float64
	var total float64
	for _, o2 := range order.Transaction.Orders {
		var totalReduced float64

		if order.Transaction.VoucherID != nil && order.Transaction.Voucher.AmountType == "percentage" {

			totalReduced = (order.Transaction.Voucher.Amount / 100) * o2.Total
			globalVoucher := &dto.GlobalDiscountReceipt{
				SellerName:   o2.Seller.Name,
				Name:         order.Transaction.Voucher.Name,
				Type:         order.Transaction.Voucher.AmountType,
				Amount:       order.Transaction.Voucher.Amount,
				TotalReduced: math.Floor(totalReduced),
			}
			globalVouchers = append(globalVouchers, globalVoucher)
		}

		if order.Transaction.VoucherID != nil && order.Transaction.Voucher.AmountType == "nominal" {

			totalReduced = (o2.Total / totalPriceBeforeDisc) * order.Transaction.Voucher.Amount
			globalVoucher := &dto.GlobalDiscountReceipt{
				SellerName:   o2.Seller.Name,
				Name:         order.Transaction.Voucher.Name,
				Type:         order.Transaction.Voucher.AmountType,
				Amount:       order.Transaction.Voucher.Amount,
				TotalReduced: math.Floor(totalReduced),
			}
			globalVouchers = append(globalVouchers, globalVoucher)
		}

		orderPayment := &dto.OrderPaymentReceipt{
			SellerName: o2.Seller.Name,
			TotalOrder: math.Floor(o2.Total + o2.Delivery.Total),
		}
		orderPayments = append(orderPayments, orderPayment)
		totalTransaction += math.Floor(o2.Total + o2.Delivery.Total)
		total += math.Floor(o2.Total + o2.Delivery.Total - totalReduced)
	}

	if order.Transaction.Voucher != nil && order.Transaction.Voucher.AmountType == "quantity" {
		globalVoucher := &dto.GlobalDiscountReceipt{
			SellerName:   "Sea Deals",
			Name:         order.Transaction.Voucher.Name,
			Type:         order.Transaction.Voucher.AmountType,
			Amount:       order.Transaction.Voucher.Amount,
			TotalReduced: order.Transaction.Voucher.Amount,
		}

		globalVouchers = append(globalVouchers, globalVoucher)
		total -= order.Transaction.Voucher.Amount
		if total < 0 {
			total = 0
		}
	}

	var orderRes = &dto.Receipt{
		SellerName: order.Seller.Name,
		Buyer: dto.BuyerReceipt{
			Name:       order.User.FullName,
			BoughtDate: order.CreatedAt,
			Address:    order.Delivery.Address,
		},
		OrderDetail: dto.OrderDetailReceipt{
			TotalQuantity:         totalQuantity,
			TotalOrder:            totalOrderBeforeDisc,
			DeliveryPrice:         order.Delivery.Total,
			Total:                 math.Floor(order.Total + order.Delivery.Total - globalVoucherForOrder.TotalReduce),
			GlobalVoucherForOrder: globalVoucherForOrder,
			ShopVoucher:           voucher,
			OrderItems:            orderItems,
		},
		Transaction: dto.TransactionReceipt{
			TotalTransaction: totalTransaction,
			GlobalDiscount:   globalVouchers,
			OrderPayments:    orderPayments,
			Total:            math.Floor(total),
		},
		Courier: dto.CourierReceipt{
			Name:    order.Delivery.Courier.Name,
			Service: "Regular",
		},
		PaymentMethod: order.Transaction.PaymentMethod,
	}
	return orderRes, nil
}

func (o *orderService) GetDetailOrderForThermal(orderID uint, userID uint) (*dto.Thermal, error) {
	tx := o.db.Begin()
	var err error
	defer helper.CommitOrRollback(tx, &err)

	seller, err := o.sellerRepository.FindSellerByUserID(tx, userID)
	if err != nil {
		return nil, err
	}
	order, err := o.orderRepository.GetOrderDetailForThermal(tx, orderID)
	if err != nil {
		return nil, err
	}
	if seller.ID != order.SellerID {
		err = apperror.UnauthorizedError("Tidak bisa melihat Thermal seller lain")
		return nil, err
	}
	var products []*dto.ProductDetailThermal
	for _, item := range order.OrderItems {
		var variantDetail string
		if item.ProductVariantDetail.ProductVariant1 != nil {
			variantDetail += *item.ProductVariantDetail.Variant1Value
		}
		if item.ProductVariantDetail.ProductVariant2 != nil {
			variantDetail += ", " + *item.ProductVariantDetail.Variant2Value
		}
		product := &dto.ProductDetailThermal{
			ID:       item.ProductVariantDetail.ID,
			Name:     item.ProductVariantDetail.Product.Name,
			Variant:  variantDetail,
			Quantity: item.Quantity,
		}
		products = append(products, product)
	}
	var orderRes = &dto.Thermal{
		Buyer: dto.BuyerThermal{
			Name:    order.User.FullName,
			Address: order.Delivery.Address,
			City:    order.Delivery.CityDestination,
		},
		Courier: dto.CourierThermal{
			Name: order.Delivery.Courier.Name,
			Code: order.Delivery.Courier.Code,
		},
		SellerName:     order.Seller.Name,
		TotalWeight:    order.Delivery.Weight,
		DeliveryNumber: order.Delivery.DeliveryNumber,
		Price:          order.Delivery.Total,
		OriginCity:     order.Seller.Address.City,
		IssuedAt:       order.CreatedAt,
		Products:       products,
	}

	return orderRes, nil
}

func (o *orderService) GetOrderBySellerID(userID uint, query *repository.OrderQuery) ([]*dto.OrderListRes, int64, int64, error) {
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
	var orderRes = make([]*dto.OrderListRes, 0)
	for _, order := range orders {
		var voucher *dto.VoucherOrderList
		var voucherID uint

		var payedAt *time.Time
		if order.Transaction.Status == dto.TransactionPayed {
			payedAt = &order.Transaction.UpdatedAt
		}
		var orderItems []*dto.OrderItemOrderList
		var priceBeforeDisc float64
		for _, item := range order.OrderItems {
			var variantDetail string
			if item.ProductVariantDetail.ProductVariant1 != nil {
				variantDetail += *item.ProductVariantDetail.Variant1Value
			}
			if item.ProductVariantDetail.ProductVariant2 != nil {
				variantDetail += ", " + *item.ProductVariantDetail.Variant2Value
			}

			var productImageURL string
			if len(item.ProductVariantDetail.Product.ProductPhotos) > 0 {
				productImageURL = item.ProductVariantDetail.Product.ProductPhotos[0].PhotoURL
			}
			var orderItemRes = &dto.OrderItemOrderList{
				ID:                     item.ID,
				ProductVariantDetailID: item.ProductVariantDetailID,
				ProductDetail: dto.ProductDetailOrderList{
					ID:         item.ProductVariantDetail.Product.ID,
					Name:       item.ProductVariantDetail.Product.Name,
					CategoryID: item.ProductVariantDetail.Product.CategoryID,
					Category:   item.ProductVariantDetail.Product.Category.Name,
					Slug:       item.ProductVariantDetail.Product.Slug,
					PhotoURL:   productImageURL,
					Variant:    variantDetail,
					Price:      item.ProductVariantDetail.Price,
				},
				Quantity: item.Quantity,
				Subtotal: item.Subtotal,
			}
			priceBeforeDisc += item.Subtotal
			orderItems = append(orderItems, orderItemRes)
		}
		if order.VoucherID != nil && *order.VoucherID != 0 {
			voucherID = *order.VoucherID
			voucher = &dto.VoucherOrderList{
				Code:          order.Voucher.Code,
				VoucherType:   order.Voucher.AmountType,
				Amount:        order.Voucher.Amount,
				AmountReduced: priceBeforeDisc - order.Total,
			}
		}
		var orderDelivery *dto.DeliveryOrderList
		var deliveryTotal float64
		var deliveryID uint
		if order.Delivery != nil {
			var orderDeliveryActivity []*dto.DeliveryActivityOrderList
			for _, activity := range order.Delivery.DeliveryActivity {
				var deliveryActivity = &dto.DeliveryActivityOrderList{
					Description: activity.Description,
					CreatedAt:   activity.CreatedAt,
				}
				orderDeliveryActivity = append(orderDeliveryActivity, deliveryActivity)
			}
			orderDelivery = &dto.DeliveryOrderList{
				DestinationAddress: order.Delivery.Address,
				Status:             order.Delivery.Status,
				DeliveryNumber:     order.Delivery.DeliveryNumber,
				ETA:                order.Delivery.Eta,
				CourierID:          order.Delivery.CourierID,
				Courier:            order.Delivery.Courier.Name,
				Activity:           orderDeliveryActivity,
			}
			deliveryTotal = order.Delivery.Total
			deliveryID = order.Delivery.ID
		}
		var res = &dto.OrderListRes{
			ID:        order.ID,
			BuyerName: order.User.FullName,
			SellerID:  order.SellerID,
			Seller: dto.SellerOrderList{
				Name: order.Seller.Name,
			},
			VoucherID:     voucherID,
			Voucher:       voucher,
			TransactionID: order.TransactionID,
			Transaction: dto.TransactionOrderList{
				PaymentMethod: order.Transaction.PaymentMethod,
				Total:         order.Transaction.Total,
				Status:        order.Transaction.Status,
				PayedAt:       payedAt,
			},
			TotalOrderPrice:          priceBeforeDisc,
			TotalOrderPriceAfterDisc: order.Total,
			TotalDelivery:            deliveryTotal,
			Status:                   order.Status,
			OrderItems:               orderItems,
			DeliveryID:               deliveryID,
			Delivery:                 orderDelivery,
			Complaint:                order.Complaint,
			UpdatedAt:                order.UpdatedAt,
		}
		orderRes = append(orderRes, res)
	}
	return orderRes, totalPage, totalData, nil
}

func (o *orderService) GetOrderByUserID(userID uint, query *repository.OrderQuery) ([]*dto.OrderListRes, int64, int64, error) {
	tx := o.db.Begin()
	var err error
	defer helper.CommitOrRollback(tx, &err)

	orders, totalPage, totalData, err := o.orderRepository.GetOrderByUserID(tx, userID, query)
	if err != nil {
		return nil, 0, 0, err
	}
	var orderRes []*dto.OrderListRes
	for _, order := range orders {
		var hasReviewEveryItem = true
		var voucher *dto.VoucherOrderList
		var voucherID uint

		var payedAt *time.Time
		if order.Transaction.Status == dto.TransactionPayed {
			payedAt = &order.Transaction.UpdatedAt
		}
		var orderItems []*dto.OrderItemOrderList
		var priceBeforeDisc float64
		for _, item := range order.OrderItems {
			var variantDetail string
			if item.ProductVariantDetail.ProductVariant1 != nil {
				variantDetail += *item.ProductVariantDetail.Variant1Value
			}
			if item.ProductVariantDetail.ProductVariant2 != nil {
				variantDetail += ", " + *item.ProductVariantDetail.Variant2Value
			}

			var productImageURL string
			if len(item.ProductVariantDetail.Product.ProductPhotos) > 0 {
				productImageURL = item.ProductVariantDetail.Product.ProductPhotos[0].PhotoURL
			}
			var review *dto.ReviewOrderList
			if item.ProductVariantDetail.Product.Review != nil {
				review = &dto.ReviewOrderList{
					ID:          item.ProductVariantDetail.Product.Review.ID,
					Rating:      item.ProductVariantDetail.Product.Review.Rating,
					Description: item.ProductVariantDetail.Product.Review.Description,
					ImageUrl:    item.ProductVariantDetail.Product.Review.ImageURL,
				}
			} else {
				hasReviewEveryItem = false
			}
			var orderItemRes = &dto.OrderItemOrderList{
				ID:                     item.ID,
				ProductVariantDetailID: item.ProductVariantDetailID,
				ProductDetail: dto.ProductDetailOrderList{
					ID:           item.ProductVariantDetail.Product.ID,
					Name:         item.ProductVariantDetail.Product.Name,
					CategoryID:   item.ProductVariantDetail.Product.CategoryID,
					Category:     item.ProductVariantDetail.Product.Category.Name,
					Slug:         item.ProductVariantDetail.Product.Slug,
					PhotoURL:     productImageURL,
					Variant:      variantDetail,
					Price:        item.ProductVariantDetail.Price,
					ReviewByUser: review,
				},
				Quantity: item.Quantity,
				Subtotal: item.Subtotal,
			}
			priceBeforeDisc += item.Subtotal
			orderItems = append(orderItems, orderItemRes)
		}
		if order.VoucherID != nil && *order.VoucherID != 0 {
			voucherID = *order.VoucherID
			voucher = &dto.VoucherOrderList{
				Code:          order.Voucher.Code,
				VoucherType:   order.Voucher.AmountType,
				Amount:        order.Voucher.Amount,
				AmountReduced: priceBeforeDisc - order.Total,
			}
		}
		var orderDelivery *dto.DeliveryOrderList
		var deliveryTotal float64
		var deliveryID uint
		if order.Delivery != nil {
			var orderDeliveryActivity []*dto.DeliveryActivityOrderList
			for _, activity := range order.Delivery.DeliveryActivity {
				var deliveryActivity = &dto.DeliveryActivityOrderList{
					Description: activity.Description,
					CreatedAt:   activity.CreatedAt,
				}
				orderDeliveryActivity = append(orderDeliveryActivity, deliveryActivity)
			}
			orderDelivery = &dto.DeliveryOrderList{
				DestinationAddress: order.Delivery.Address,
				Status:             order.Delivery.Status,
				DeliveryNumber:     order.Delivery.DeliveryNumber,
				ETA:                order.Delivery.Eta,
				CourierID:          order.Delivery.CourierID,
				Courier:            order.Delivery.Courier.Name,
				Activity:           orderDeliveryActivity,
			}
			deliveryTotal = order.Delivery.Total
			deliveryID = order.Delivery.ID
		}
		var res = &dto.OrderListRes{
			ID:        order.ID,
			BuyerName: order.User.FullName,
			SellerID:  order.SellerID,
			Seller: dto.SellerOrderList{
				Name: order.Seller.Name,
			},
			VoucherID:     voucherID,
			Voucher:       voucher,
			TransactionID: order.TransactionID,
			Transaction: dto.TransactionOrderList{
				PaymentMethod: order.Transaction.PaymentMethod,
				Total:         order.Transaction.Total,
				Status:        order.Transaction.Status,
				PayedAt:       payedAt,
			},
			TotalOrderPrice:          priceBeforeDisc,
			TotalOrderPriceAfterDisc: order.Total,
			TotalDelivery:            deliveryTotal,
			Status:                   order.Status,
			HasReviewedAllItem:       hasReviewEveryItem,
			OrderItems:               orderItems,
			DeliveryID:               deliveryID,
			Delivery:                 orderDelivery,
			Complaint:                order.Complaint,
			UpdatedAt:                order.UpdatedAt,
		}
		orderRes = append(orderRes, res)
	}
	return orderRes, totalPage, totalData, nil
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
		err = apperror.BadRequestError("Cannot cancel order that is currently " + order.Status)
		return nil, err
	}

	seller, err := o.sellerRepository.FindSellerByUserID(tx, userID)
	if err != nil {
		return nil, err
	}
	if order.SellerID != seller.ID {
		err = apperror.BadRequestError("Cannot cancel another seller order")
		return nil, err
	}

	var priceBeforeGlobalDisc float64
	var voucher *model.Voucher
	var delivery *model.Delivery
	var amountRefunded = order.Total
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
			amountReduced := (order.Total / priceBeforeGlobalDisc) * voucher.Amount
			amountRefunded = order.Total - amountReduced
		}
	}
	delivery, err = o.deliveryRepo.GetDeliveryByOrderID(tx, order.ID)
	if err != nil {
		return nil, err
	}
	amountRefunded += delivery.Total

	var buyerWallet *model.Wallet
	var transHolder *model.SeaLabsPayTransactionHolder
	var req *http.Request
	var resp *http.Response
	if order.Transaction.PaymentMethod == dto.Wallet {
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
			PaymentMethod: dto.Wallet,
			PaymentType:   "CREDIT",
			Description:   "Refund from transaction ID " + strconv.FormatUint(uint64(order.TransactionID), 10),
			CreatedAt:     time.Time{},
		}
		_, err = o.walletTransRepo.CreateTransaction(tx, walletTrans)
		if err != nil {
			return nil, err
		}
	} else if order.Transaction.PaymentMethod == dto.SeaLabsPay {
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

	_, err = o.accountHolderRepo.TakeMoneyFromAccountHolderByOrderID(tx, order.ID)
	if err != nil {
		return nil, err
	}

	refundedOrder, err := o.orderRepository.UpdateOrderStatus(tx, orderID, dto.OrderRefunded)
	if err != nil {
		return nil, err
	}

	newNotification := &model.Notification{
		UserID:   order.UserID,
		SellerID: order.SellerID,
		Title:    dto.NotificationSellerMembatalkanPesanan,
		Detail:   "Seller membatalkan pesanan",
	}

	o.notificationRepo.AddToNotificationFromModel(tx, newNotification)

	return refundedOrder, nil
}

func (o *orderService) RequestRefundByBuyer(req *dto.CreateComplaintReq, userID uint) (*dto.CreateComplaintRes, error) {
	tx := o.db.Begin()
	var err error
	defer helper.CommitOrRollback(tx, &err)

	order, err := o.orderRepository.GetOrderDetailByID(tx, req.OrderID)
	if err != nil {
		return nil, err
	}
	if order.UserID != userID {
		err = apperror.BadRequestError("Tidak bisa membatalkan order user lain")
		return nil, err
	}
	if order.Status != dto.OrderDelivered {
		err = apperror.BadRequestError("Cannot refund order that is currently " + order.Status)
		return nil, err
	}

	updatedOrder, err := o.orderRepository.UpdateOrderStatus(tx, req.OrderID, dto.OrderComplained)
	if err != nil {
		return nil, err
	}
	complaint, err := o.complaintRepo.CreateComplaint(tx, req.OrderID, req.Description)
	if err != nil {
		return nil, err
	}

	var photos []*model.ComplaintPhoto
	for _, photo := range req.Photos {
		var data = &model.ComplaintPhoto{
			ComplaintID: complaint.ID,
			PhotoURL:    photo.PhotoURL,
			PhotoName:   photo.PhotoName,
		}
		photos = append(photos, data)
	}
	complaintPhoto, err := o.complaintPhotoRepo.CreateComplaintPhotos(tx, photos)
	if err != nil {
		return nil, err
	}

	res := &dto.CreateComplaintRes{
		Order:           updatedOrder,
		ComplaintPhotos: complaintPhoto,
		Description:     complaint.Description,
	}

	return res, nil
}

func (o *orderService) AcceptRefundRequest(req *dto.RejectAcceptRefundReq, userID uint) (*dto.RejectAcceptRefundRes, error) {
	tx := o.db.Begin()
	var err error
	defer helper.CommitOrRollback(tx, &err)

	seller, err := o.sellerRepository.FindSellerByUserID(tx, userID)
	if err != nil {
		return nil, err
	}
	order, err := o.orderRepository.GetOrderDetailByID(tx, req.OrderID)
	if err != nil {
		return nil, err
	}
	if order.SellerID != seller.ID {
		err = apperror.BadRequestError("Tidak bisa menyetujui refund request seller lain")
		return nil, err
	}
	if order.Status != dto.OrderComplained {
		err = apperror.BadRequestError("Cannot accept refund order that is currently " + order.Status)
		return nil, err
	}

	var priceBeforeGlobalDisc float64
	var voucher *model.Voucher
	var delivery *model.Delivery
	var amountRefunded = order.Total
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
			amountReduced := (order.Total / priceBeforeGlobalDisc) * voucher.Amount
			amountRefunded = order.Total - amountReduced
		}
	}
	delivery, err = o.deliveryRepo.GetDeliveryByOrderID(tx, order.ID)
	if err != nil {
		return nil, err
	}
	amountRefunded += delivery.Total

	var buyerWallet *model.Wallet
	var transHolder *model.SeaLabsPayTransactionHolder
	if order.Transaction.PaymentMethod == dto.Wallet {
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
			PaymentMethod: dto.Wallet,
			PaymentType:   "CREDIT",
			Description:   "Refund from transaction ID " + strconv.FormatUint(uint64(order.TransactionID), 10),
			CreatedAt:     time.Time{},
		}
		_, err = o.walletTransRepo.CreateTransaction(tx, walletTrans)
		if err != nil {
			return nil, err
		}
	} else if order.Transaction.PaymentMethod == dto.SeaLabsPay {
		transHolder, err = o.seaLabsPayTransHolderRepo.GetTransHolderFromTransactionID(tx, order.TransactionID)
		if err != nil {
			return nil, err
		}

		URL := config.Config.SeaLabsPayRefundURL
		var jsonStr = []byte(`{"reason":"Seller cancel the order", "amount":` + strconv.Itoa(int(amountRefunded)) + `, "txn_id":` + strconv.Itoa(int(transHolder.TxnID)) + `}`)

		err = refundMoneyToSeaLabsPay(URL, jsonStr)
		if err != nil {
			return nil, err
		}
	}

	for _, item := range order.OrderItems {
		_, err = o.productVarDetRepo.AddProductVariantStock(tx, item.ProductVariantDetailID, item.Quantity)
		if err != nil {
			return nil, err
		}
	}

	_, err = o.accountHolderRepo.TakeMoneyFromAccountHolderByOrderID(tx, order.ID)
	if err != nil {
		return nil, err
	}

	refundedOrder, err := o.orderRepository.UpdateOrderStatus(tx, req.OrderID, dto.OrderRefunded)
	if err != nil {
		return nil, err
	}

	response := &dto.RejectAcceptRefundRes{
		Order:          refundedOrder,
		AmountRefunded: amountRefunded,
	}
	newNotification := &model.Notification{
		UserID:   order.UserID,
		SellerID: order.SellerID,
		Title:    dto.NotificationSellerMenyetujuiRefund,
		Detail:   "Seller menyetujui refund request",
	}

	o.notificationRepo.AddToNotificationFromModel(tx, newNotification)

	return response, nil
}

func (o *orderService) RejectRefundRequest(req *dto.RejectAcceptRefundReq, userID uint) (*dto.RejectAcceptRefundRes, error) {
	tx := o.db.Begin()
	var err error
	defer helper.CommitOrRollback(tx, &err)

	seller, err := o.sellerRepository.FindSellerByUserID(tx, userID)
	if err != nil {
		return nil, err
	}
	order, err := o.orderRepository.GetOrderDetailByID(tx, req.OrderID)
	if err != nil {
		return nil, err
	}
	if order.SellerID != seller.ID {
		err = apperror.BadRequestError("Tidak bisa menolak refund request seller lain")
		return nil, err
	}
	if order.Status != dto.OrderComplained {
		err = apperror.BadRequestError("Cannot reject refund order that is currently " + order.Status)
		return nil, err
	}

	for _, item := range order.OrderItems {
		_, err = o.productVarDetRepo.AddProductVariantStock(tx, item.ProductVariantDetailID, item.Quantity)
		if err != nil {
			return nil, err
		}
	}

	// ADD GET HOLDING ACCOUNT MONEY HERE
	accountHolder, err := o.accountHolderRepo.TakeMoneyFromAccountHolderByOrderID(tx, order.ID)
	if err != nil {
		return nil, err
	}
	wallet, err := o.walletRepository.GetWalletByUserID(tx, seller.UserID)
	if err != nil {
		return nil, err
	}
	_, err = o.walletRepository.TopUp(tx, wallet, accountHolder.Total)
	if err != nil {
		return nil, err
	}

	doneOrder, err := o.orderRepository.UpdateOrderStatus(tx, req.OrderID, dto.OrderDone)
	if err != nil {
		return nil, err
	}

	response := &dto.RejectAcceptRefundRes{
		Order:          doneOrder,
		AmountRefunded: 0,
	}
	newNotification := &model.Notification{
		UserID:   order.UserID,
		SellerID: order.SellerID,
		Title:    dto.NotificationSellerMenolakRefund,
		Detail:   "Seller menolak refund request",
	}

	o.notificationRepo.AddToNotificationFromModel(tx, newNotification)
	return response, nil
}

func (o *orderService) FinishOrder(req *dto.FinishOrderReq, userID uint) (*model.Order, error) {
	tx := o.db.Begin()
	var err error
	defer helper.CommitOrRollback(tx, &err)

	order, err := o.orderRepository.GetOrderDetailByID(tx, req.OrderID)
	if err != nil {
		return nil, err
	}
	if order.UserID != userID {
		err = apperror.BadRequestError("Tidak bisa menyelesaikan order user lain")
		return nil, err
	}
	if order.Status != dto.OrderDelivered {
		err = apperror.BadRequestError("Tidak bisa menyelesaikan order yang sedang dalam proses " + order.Status)
		return nil, err
	}

	accountHolder, err := o.accountHolderRepo.TakeMoneyFromAccountHolderByOrderID(tx, order.ID)
	if err != nil {
		return nil, err
	}
	seller, err := o.sellerRepository.FindSellerByID(tx, order.SellerID)
	if err != nil {
		return nil, err
	}
	wallet, err := o.walletRepository.GetWalletByUserID(tx, seller.UserID)
	if err != nil {
		return nil, err
	}
	_, err = o.walletRepository.TopUp(tx, wallet, accountHolder.Total)
	if err != nil {
		return nil, err
	}
	transWalletRepo := &model.WalletTransaction{
		WalletID:      wallet.ID,
		TransactionID: &order.TransactionID,
		Total:         accountHolder.Total,
		PaymentMethod: "wallet",
		PaymentType:   "CREDIT",
		Description:   "Pembayaran dari order ID " + strconv.FormatUint(uint64(order.ID), 10),
	}
	_, err = o.walletTransRepo.CreateTransaction(tx, transWalletRepo)
	if err != nil {
		return nil, err
	}

	doneOrder, err := o.orderRepository.UpdateOrderStatus(tx, req.OrderID, dto.OrderDone)
	if err != nil {
		return nil, err
	}
	for _, item := range order.OrderItems {
		_, err = o.productRepo.AddProductSoldCount(tx, item.ProductVariantDetail.ProductID, int(item.Quantity))
		if err != nil {
			return nil, err
		}
	}

	newNotification := &model.Notification{
		UserID:   order.UserID,
		SellerID: order.SellerID,
		Title:    dto.NotificationPesananSelesai,
		Detail:   "Order produk telah diselesaikan",
	}
	o.notificationRepo.AddToNotificationFromModel(tx, newNotification)

	return doneOrder, nil
}

func (o *orderService) RunCronJobs() {

	c := cron.New(cron.WithLocation(time.UTC))
	_, _ = c.AddFunc(config.Config.IntervalCron, func() {
		deliveries, _ := o.deliveryRepo.FindAndUpdateOngoingToDelivered()
		tx := o.db.Begin()
		for _, delivery := range deliveries {
			order, _ := o.orderRepository.UpdateOrderStatus(tx, delivery.OrderID, dto.OrderDelivered)
			newNotification := &model.Notification{
				UserID:   order.UserID,
				SellerID: order.SellerID,
				Title:    dto.NotificationPesananSampai,
				Detail:   "Order dengan ID " + strconv.FormatUint(uint64(order.ID), 10) + " sampai Tujuan",
			}
			o.notificationRepo.AddToNotificationFromModelForCron(newNotification)
		}
		tx.Commit()
	})

	_, _ = c.AddFunc(config.Config.IntervalCron, func() {
		orders := o.orderRepository.FindAndUpdateDeliveredOrderToDone()
		for _, order := range orders {
			tx := o.db.Begin()
			accountHolder, _ := o.accountHolderRepo.TakeMoneyFromAccountHolderByOrderID(tx, order.ID)
			orderDetail, _ := o.orderRepository.GetOrderDetailByID(tx, order.ID)

			seller, _ := o.sellerRepository.FindSellerByID(tx, order.SellerID)
			wallet, _ := o.walletRepository.GetWalletByUserID(tx, seller.UserID)
			_, _ = o.walletRepository.TopUp(tx, wallet, accountHolder.Total)
			for _, item := range orderDetail.OrderItems {
				_, _ = o.productRepo.AddProductSoldCount(tx, item.ProductVariantDetail.ProductID, int(item.Quantity))
			}
			transWalletRepo := &model.WalletTransaction{
				WalletID:      wallet.ID,
				TransactionID: &order.TransactionID,
				Total:         accountHolder.Total,
				PaymentMethod: "wallet",
				PaymentType:   "CREDIT",
				Description:   "Pembayaran dari order ID " + strconv.FormatUint(uint64(order.ID), 10),
			}
			_, _ = o.walletTransRepo.CreateTransaction(tx, transWalletRepo)

			tx.Commit()

			newNotification := &model.Notification{
				UserID:   order.UserID,
				SellerID: order.SellerID,
				Title:    dto.NotificationPesananSelesai,
				Detail:   "Order dengan ID " + strconv.FormatUint(uint64(order.ID), 10) + " selesai",
			}
			o.notificationRepo.AddToNotificationFromModelForCron(newNotification)
		}
	})

	_, _ = c.AddFunc(config.Config.IntervalCron, func() {
		orders := o.orderRepository.FindAndUpdateWaitingForSellerToRefunded()
		for _, order := range orders {
			tx := o.db.Begin()
			orderDetail, _ := o.orderRepository.GetOrderDetailByID(tx, order.ID)
			if orderDetail.Transaction.PaymentMethod == dto.Wallet {
				var wallet *model.Wallet
				var orderItems []*model.OrderItem
				var amountRefunded = orderDetail.Total
				if orderDetail.Transaction.VoucherID != nil {
					priceBeforeGlobalDisc, _ := o.transactionRepo.GetPriceBeforeGlobalDisc(tx, orderDetail.TransactionID)
					voucher, _ := o.voucherRepo.FindVoucherDetailByID(tx, *orderDetail.Transaction.VoucherID)
					if voucher.AmountType == "percentage" {
						amountRefunded = orderDetail.Total - ((voucher.Amount / 100) * orderDetail.Total)
					} else {
						amountReduced := (orderDetail.Total / priceBeforeGlobalDisc) * voucher.Amount
						amountRefunded = orderDetail.Total - amountReduced
					}
				}
				delivery, _ := o.deliveryRepo.GetDeliveryByOrderID(tx, orderDetail.ID)
				amountRefunded += delivery.Total
				tx.Commit()

				wallet = o.orderRepository.RefundToWalletByUserID(orderDetail.UserID, amountRefunded)
				o.orderRepository.AddToWalletTransaction(wallet.ID, amountRefunded)
				orderItems = o.orderRepository.GetOrderItemsByOrderID(orderDetail.ID)
				for _, orderItem := range orderItems {
					o.orderRepository.UpdateStockByProductVariantDetailID(orderItem.ProductVariantDetailID, orderItem.Quantity)
				}
			} else {
				var amountRefunded = orderDetail.Total
				if orderDetail.Transaction.VoucherID != nil {
					priceBeforeGlobalDisc, _ := o.transactionRepo.GetPriceBeforeGlobalDisc(tx, orderDetail.TransactionID)
					voucher, _ := o.voucherRepo.FindVoucherDetailByID(tx, *orderDetail.Transaction.VoucherID)
					if voucher.AmountType == "percentage" {
						amountRefunded = orderDetail.Total - ((voucher.Amount / 100) * orderDetail.Total)
					} else {
						amountReduced := (orderDetail.Total / priceBeforeGlobalDisc) * orderDetail.Total
						amountRefunded = orderDetail.Total - amountReduced
					}
				}

				delivery, _ := o.deliveryRepo.GetDeliveryByOrderID(tx, orderDetail.ID)
				amountRefunded += delivery.Total

				transHolder, err := o.seaLabsPayTransHolderRepo.GetTransHolderFromTransactionID(tx, orderDetail.TransactionID)
				URL := config.Config.SeaLabsPayRefundURL
				var jsonStr = []byte(`{"reason":"Seller cancel the order", "amount":` + strconv.Itoa(int(amountRefunded)) + `, "txn_id":` + strconv.Itoa(int(transHolder.TxnID)) + `}`)

				err = refundMoneyToSeaLabsPay(URL, jsonStr)
				orderItems := o.orderRepository.GetOrderItemsByOrderID(orderDetail.ID)
				for _, orderItem := range orderItems {
					o.orderRepository.UpdateStockByProductVariantDetailID(orderItem.ProductVariantDetailID, orderItem.Quantity)
				}
				if err != nil {
					fmt.Println("ERROR: ", err)
					tx.Rollback()
				}
				tx.Commit()
			}
		}
	})

	c.Start()

}

func (o *orderService) GetTotalPredictedPrice(req *dto.PredictedPriceReq, userID uint) (*dto.TotalPredictedPriceRes, error) {
	tx := o.db.Begin()
	var err error
	defer helper.CommitOrRollback(tx, &err)

	if len(req.Cart) <= 0 {
		err = apperror.BadRequestError("Checkout setidaknya harus terdapat satu barang")
		return nil, err
	}

	var res = &dto.TotalPredictedPriceRes{}

	globalVoucher, err := o.walletRepository.GetVoucher(tx, req.GlobalVoucherCode)
	if err != nil {
		return nil, err
	}
	timeNow := time.Now()
	if globalVoucher != nil {
		if timeNow.After(globalVoucher.EndDate) || timeNow.Before(globalVoucher.StartDate) {
			err = apperror.InternalServerError("Level 3 Voucher invalid")
			return nil, err
		}
	}

	var voucherID *uint
	if globalVoucher != nil {
		voucherID = &globalVoucher.ID
	}
	var ordersPrices []*dto.PredictedPriceRes
	res.GlobalVoucherID = voucherID
	var totalAllOrderPrices float64
	var totalDelivery float64
	var sellerIDs []uint

	for _, item := range req.Cart {
		for _, id := range sellerIDs {
			if id == item.SellerID {
				err = apperror.BadRequestError("Tidak bisa membuat 2 order dengan seller yang sama dalam satu transaksi")
				return nil, err
			}
		}

		var predictedPrice = &dto.PredictedPriceRes{}
		var voucher *model.Voucher
		voucher, err = o.walletRepository.GetVoucher(tx, item.VoucherCode)
		if err != nil {
			return nil, err
		}

		predictedPrice.SellerID = item.SellerID

		if voucher != nil {
			if timeNow.After(voucher.EndDate) || timeNow.Before(voucher.StartDate) {
				err = apperror.InternalServerError("Level 2 Voucher invalid")
				return nil, err
			}
			predictedPrice.VoucherID = &voucher.ID
		} else {
			predictedPrice.VoucherID = nil
		}

		var totalOrder float64
		var totalWeight int

		for _, id := range item.CartItemID {

			var totalOrderItem float64
			var cartItem *model.CartItem
			cartItem, err = o.walletRepository.GetCartItem(tx, id)
			if err != nil {
				return nil, err
			}

			if cartItem.ProductVariantDetail.Product.SellerID != item.SellerID {
				err = apperror.BadRequestError("That cart item does not belong to that seller")
				return nil, err
			}

			//check stock
			newStock := cartItem.ProductVariantDetail.Stock - cartItem.Quantity
			if newStock < 0 {
				err = apperror.InternalServerError(cartItem.ProductVariantDetail.Product.Name + "is out of stock")
				return nil, err
			}

			if cartItem.ProductVariantDetail.Product.Promotion != nil && cartItem.ProductVariantDetail.Product.Promotion.MaxOrder >= cartItem.Quantity {
				totalOrderItem = (cartItem.ProductVariantDetail.Price - cartItem.ProductVariantDetail.Product.Promotion.Amount) * float64(cartItem.Quantity)
			} else {
				totalOrderItem = cartItem.ProductVariantDetail.Price * float64(cartItem.Quantity)
			}
			if totalOrderItem < 0 {
				totalOrderItem = 0
			}
			totalOrder += totalOrderItem

			// Get weight
			totalWeight += int(cartItem.Quantity) * cartItem.ProductVariantDetail.Product.ProductDetail.Weight
			if totalWeight > 20000 {
				err = apperror.BadRequestError(cartItem.ProductVariantDetail.Product.Name + " exceeded weight limit of 20000")
				return nil, apperror.BadRequestError(cartItem.ProductVariantDetail.Product.Name + " exceeded weight limit of 20000")
			}

		}
		//order - voucher
		if voucher != nil && voucher.MinSpending <= totalOrder {
			if voucher.AmountType == "percentage" {
				totalOrder -= (voucher.Amount / 100) * totalOrder
			} else {
				totalOrder -= voucher.Amount
			}
		} else if voucher != nil {
			err = apperror.BadRequestError("Order tidak memenuhi kriteria voucher " + voucher.Name)
			return nil, err
		}
		if totalOrder < 0 {
			totalOrder = 0
		}

		var seller *model.Seller
		seller, err = o.sellerRepository.FindSellerByID(tx, item.SellerID)
		if err != nil {
			return nil, err
		}

		// Check delivery
		var courier *model.Courier
		courier, err = o.courierRepository.GetCourierDetailByID(tx, item.CourierID)
		if err != nil {
			return nil, err
		}
		var buyerAddress *model.Address
		buyerAddress, err = o.addressRepository.CheckUserAddress(tx, req.BuyerAddressID, userID)
		if err != nil {
			return nil, err
		}

		deliveryReq := &dto.DeliveryCalculateReq{
			OriginCity:      seller.Address.CityID,
			DestinationCity: buyerAddress.CityID,
			Weight:          strconv.Itoa(totalWeight),
			Courier:         courier.Code,
		}
		var deliveryCalcResult *dto.DeliveryCalculateReturn

		deliveryCalcResult, err = helper.CalculateDeliveryPrice(deliveryReq)
		if err != nil {
			return nil, err
		}

		predictedPrice.DeliveryPrice = float64(deliveryCalcResult.Total)
		predictedPrice.TotalOrder = totalOrder
		predictedPrice.PredictedPrice = totalOrder + float64(deliveryCalcResult.Total)

		ordersPrices = append(ordersPrices, predictedPrice)
		totalAllOrderPrices += predictedPrice.TotalOrder
		totalDelivery += float64(deliveryCalcResult.Total)
		sellerIDs = append(sellerIDs, item.SellerID)
	}

	res.PredictedPrices = ordersPrices

	if globalVoucher != nil && globalVoucher.SellerID == nil && globalVoucher.MinSpending <= totalAllOrderPrices {
		if globalVoucher.AmountType == "percentage" {
			totalAllOrderPrices -= (globalVoucher.Amount / 100) * totalAllOrderPrices
		} else {
			totalAllOrderPrices -= globalVoucher.Amount
		}
	} else if globalVoucher != nil {
		err = apperror.BadRequestError("Order tidak memenuhi kriteria voucher global")
		return nil, err
	}
	if totalAllOrderPrices < 0 {
		totalAllOrderPrices = 0
	}

	res.TotalPredictedPrice = totalAllOrderPrices + totalDelivery
	return res, nil
}

func (o *orderService) GetOrderByID(userID uint, orderID uint) (*dto.OrderListRes, error) {
	tx := o.db.Begin()
	var err error
	defer helper.CommitOrRollback(tx, &err)

	order, err := o.orderRepository.GetOrderByID(tx, userID, orderID)
	if err != nil {
		return nil, err
	}
	if order.UserID != userID {
		err = apperror.BadRequestError("Tidak bisa membatalkan order user lain")
		return nil, err
	}
	var hasReviewEveryItem = true
	var voucher *dto.VoucherOrderList
	var voucherID uint

	var payedAt *time.Time
	if order.Transaction.Status == dto.TransactionPayed {
		payedAt = &order.Transaction.UpdatedAt
	}

	var orderItems []*dto.OrderItemOrderList
	var priceBeforeDisc float64
	for _, item := range order.OrderItems {
		var variantDetail string
		if item.ProductVariantDetail.ProductVariant1 != nil {
			variantDetail += *item.ProductVariantDetail.Variant1Value
		}
		if item.ProductVariantDetail.ProductVariant2 != nil {
			variantDetail += ", " + *item.ProductVariantDetail.Variant2Value
		}

		var productImageURL string
		if len(item.ProductVariantDetail.Product.ProductPhotos) > 0 {
			productImageURL = item.ProductVariantDetail.Product.ProductPhotos[0].PhotoURL
		}

		var review *dto.ReviewOrderList
		if item.ProductVariantDetail.Product.Review != nil {
			review = &dto.ReviewOrderList{
				ID:          item.ProductVariantDetail.Product.Review.ID,
				Rating:      item.ProductVariantDetail.Product.Review.Rating,
				Description: item.ProductVariantDetail.Product.Review.Description,
				ImageUrl:    item.ProductVariantDetail.Product.Review.ImageURL,
			}
		} else {
			hasReviewEveryItem = false
		}

		var orderItemRes = &dto.OrderItemOrderList{
			ID:                     item.ID,
			ProductVariantDetailID: item.ProductVariantDetailID,
			ProductDetail: dto.ProductDetailOrderList{
				ID:           item.ProductVariantDetail.Product.ID,
				Name:         item.ProductVariantDetail.Product.Name,
				CategoryID:   item.ProductVariantDetail.Product.CategoryID,
				Category:     item.ProductVariantDetail.Product.Category.Name,
				Slug:         item.ProductVariantDetail.Product.Slug,
				PhotoURL:     productImageURL,
				Variant:      variantDetail,
				Price:        item.ProductVariantDetail.Price,
				ReviewByUser: review,
			},
			Quantity: item.Quantity,
			Subtotal: item.Subtotal,
		}
		priceBeforeDisc += item.Subtotal
		orderItems = append(orderItems, orderItemRes)
	}

	if order.VoucherID != nil && *order.VoucherID != 0 {
		voucherID = *order.VoucherID
		voucher = &dto.VoucherOrderList{
			Code:          order.Voucher.Code,
			VoucherType:   order.Voucher.AmountType,
			Amount:        order.Voucher.Amount,
			AmountReduced: priceBeforeDisc - order.Total,
		}
	}

	var orderDelivery *dto.DeliveryOrderList
	var deliveryTotal float64
	var deliveryID uint
	if order.Delivery != nil {
		var orderDeliveryActivity []*dto.DeliveryActivityOrderList
		for _, activity := range order.Delivery.DeliveryActivity {
			var deliveryActivity = &dto.DeliveryActivityOrderList{
				Description: activity.Description,
				CreatedAt:   activity.CreatedAt,
			}
			orderDeliveryActivity = append(orderDeliveryActivity, deliveryActivity)
		}

		orderDelivery = &dto.DeliveryOrderList{
			DestinationAddress: order.Delivery.Address,
			Status:             order.Delivery.Status,
			DeliveryNumber:     order.Delivery.DeliveryNumber,
			ETA:                order.Delivery.Eta,
			CourierID:          order.Delivery.CourierID,
			Courier:            order.Delivery.Courier.Name,
			Activity:           orderDeliveryActivity,
		}
		deliveryTotal = order.Delivery.Total
		deliveryID = order.Delivery.ID
	}

	var res = &dto.OrderListRes{
		ID:        order.ID,
		BuyerName: order.User.FullName,
		SellerID:  order.SellerID,
		Seller: dto.SellerOrderList{
			Name: order.Seller.Name,
		},
		VoucherID:     voucherID,
		Voucher:       voucher,
		TransactionID: order.TransactionID,
		Transaction: dto.TransactionOrderList{
			PaymentMethod: order.Transaction.PaymentMethod,
			Total:         order.Transaction.Total,
			Status:        order.Transaction.Status,
			PayedAt:       payedAt,
		},
		TotalOrderPrice:          priceBeforeDisc,
		TotalOrderPriceAfterDisc: order.Total,
		TotalDelivery:            deliveryTotal,
		Status:                   order.Status,
		HasReviewedAllItem:       hasReviewEveryItem,
		OrderItems:               orderItems,
		DeliveryID:               deliveryID,
		Delivery:                 orderDelivery,
		Complaint:                order.Complaint,
		UpdatedAt:                order.UpdatedAt,
	}
	return res, nil
}
