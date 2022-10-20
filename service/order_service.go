package service

import (
	"bytes"
	"encoding/json"
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
	GetOrderBySellerID(userID uint, query *repository.OrderQuery) ([]*dto.OrderListRes, int64, int64, error)
	GetOrderByUserID(userID uint, query *repository.OrderQuery) ([]*dto.OrderListRes, int64, int64, error)

	CancelOrderBySeller(orderID uint, userID uint) (*model.Order, error)
	RequestRefundByBuyer(req *dto.CreateComplaintReq, userID uint) (*dto.CreateComplaintRes, error)
	AcceptRefundRequest(req *dto.RejectAcceptRefundReq, userID uint) (*dto.RejectAcceptRefundRes, error)
	RejectRefundRequest(req *dto.RejectAcceptRefundReq, userID uint) (*dto.RejectAcceptRefundRes, error)
	FinishOrder(req *dto.FinishOrderReq, userID uint) (*model.Order, error)

	RunCronJobs()
	GetTotalPredictedPrice(req *dto.TotalPredictedPriceReq, userID uint) (*dto.TotalPredictedPriceRes, error)
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
		return apperror.BadRequestError(j.Message)
	}
	return nil
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
	var orderRes []*dto.OrderListRes
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

			var imageURL string
			if len(item.ProductVariantDetail.Product.ProductPhotos) > 0 {
				imageURL = item.ProductVariantDetail.Product.ProductPhotos[0].PhotoURL
			}

			var orderItemRes = &dto.OrderItemOrderList{
				ID:                     item.ID,
				ProductVariantDetailID: item.ProductVariantDetailID,
				ProductDetail: dto.ProductDetailOrderList{
					Name:       item.ProductVariantDetail.Product.Name,
					CategoryID: item.ProductVariantDetail.Product.CategoryID,
					Category:   item.ProductVariantDetail.Product.Category.Name,
					Slug:       item.ProductVariantDetail.Product.Slug,
					PhotoURL:   imageURL,
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
			ID:       order.ID,
			SellerID: order.SellerID,
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

			var imageURL string
			if len(item.ProductVariantDetail.Product.ProductPhotos) > 0 {
				imageURL = item.ProductVariantDetail.Product.ProductPhotos[0].PhotoURL
			}

			var orderItemRes = &dto.OrderItemOrderList{
				ID:                     item.ID,
				ProductVariantDetailID: item.ProductVariantDetailID,
				ProductDetail: dto.ProductDetailOrderList{
					Name:       item.ProductVariantDetail.Product.Name,
					CategoryID: item.ProductVariantDetail.Product.CategoryID,
					Category:   item.ProductVariantDetail.Product.Category.Name,
					Slug:       item.ProductVariantDetail.Product.Slug,
					PhotoURL:   imageURL,
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
			ID:       order.ID,
			SellerID: order.SellerID,
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
			amountReduced := (order.Total / priceBeforeGlobalDisc) * order.Total
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
			amountReduced := (order.Total / priceBeforeGlobalDisc) * order.Total
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
	_, _ = c.AddFunc("@daily", func() {
		deliveries, _ := o.deliveryRepo.CheckAndUpdateToDelivered()
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

	_, _ = c.AddFunc("@daily", func() {
		orders := o.orderRepository.CheckAndUpdateOnOrderDelivered()
		for _, order := range orders {
			tx := o.db.Begin()
			accountHolder, _ := o.accountHolderRepo.TakeMoneyFromAccountHolderByOrderID(tx, order.ID)

			seller, _ := o.sellerRepository.FindSellerByID(tx, order.SellerID)
			wallet, _ := o.walletRepository.GetWalletByUserID(tx, seller.UserID)
			_, _ = o.walletRepository.TopUp(tx, wallet, accountHolder.Total)
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

	_, _ = c.AddFunc("@daily", func() {
		orders := o.orderRepository.CheckAndUpdateWaitingForSeller()
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
						amountReduced := (orderDetail.Total / priceBeforeGlobalDisc) * orderDetail.Total
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
				if orderDetail.Transaction.VoucherID != nil {
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
						tx.Rollback()
					}
				}
				tx.Commit()
			}
		}
	})

	c.Start()

}

func (o *orderService) GetTotalPredictedPrice(req *dto.TotalPredictedPriceReq, userID uint) (*dto.TotalPredictedPriceRes, error) {
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
	}
	res.TotalPredictedPrice = totalAllOrderPrices + totalDelivery
	return res, nil
}
