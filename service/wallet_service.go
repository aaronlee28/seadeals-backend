package service

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
	"math"
	"seadeals-backend/apperror"
	"seadeals-backend/config"
	"seadeals-backend/dto"
	"seadeals-backend/helper"
	"seadeals-backend/model"
	"seadeals-backend/repository"
	"strconv"
	"strings"
	"time"

	"gorm.io/gorm"
)

type WalletService interface {
	UserWalletData(id uint) (*dto.WalletDataRes, error)
	TransactionDetails(userID uint, transactionID uint) (*dto.TransactionDetailsRes, error)
	PaginatedTransactions(q *repository.Query, userID uint) (*dto.PaginatedTransactionsRes, error)
	GetWalletTransactionsByUserID(q *dto.WalletTransactionsQuery, userID uint) ([]*model.WalletTransaction, int64, int64, error)

	WalletPin(userID uint, pin string) error
	RequestPinChangeWithEmail(userID uint) (string, string, error)
	ValidateRequestIsValid(userID uint, key string) (string, error)
	ValidateCodeToRequestByEmail(userID uint, req *dto.CodeKeyRequestByEmailReq) (string, error)
	ChangeWalletPinByEmail(userID uint, req *dto.ChangePinByEmailReq) (*model.Wallet, error)
	ValidateWalletPin(user *dto.UserJWT, pin string) (string, bool, error)

	GetWalletStatus(userID uint) (string, error)
	PayOrderWithWallet(userID uint, req *dto.CheckoutCartReq) (*dto.CheckoutCartRes, error)
}

type walletService struct {
	db                *gorm.DB
	addressRepository repository.AddressRepository
	walletRepository  repository.WalletRepository
	courierRepository repository.CourierRepository
	deliveryRepo      repository.DeliveryRepository
	deliveryActRepo   repository.DeliveryActivityRepository
	userRepository    repository.UserRepository
	walletTransRepo   repository.WalletTransactionRepository
	userRoleRepo      repository.UserRoleRepository
	sellerRepository  repository.SellerRepository
	accountHolderRepo repository.AccountHolderRepository
}

type WalletServiceConfig struct {
	DB                *gorm.DB
	AddressRepository repository.AddressRepository
	WalletRepository  repository.WalletRepository
	CourierRepository repository.CourierRepository
	DeliveryRepo      repository.DeliveryRepository
	DeliveryActRepo   repository.DeliveryActivityRepository
	UserRepository    repository.UserRepository
	WalletTransRepo   repository.WalletTransactionRepository
	UserRoleRepo      repository.UserRoleRepository
	SellerRepository  repository.SellerRepository
	AccountHolderRepo repository.AccountHolderRepository
}

func NewWalletService(c *WalletServiceConfig) WalletService {
	return &walletService{
		db:                c.DB,
		addressRepository: c.AddressRepository,
		walletRepository:  c.WalletRepository,
		courierRepository: c.CourierRepository,
		deliveryRepo:      c.DeliveryRepo,
		deliveryActRepo:   c.DeliveryActRepo,
		userRepository:    c.UserRepository,
		walletTransRepo:   c.WalletTransRepo,
		userRoleRepo:      c.UserRoleRepo,
		sellerRepository:  c.SellerRepository,
		accountHolderRepo: c.AccountHolderRepo,
	}
}

func (w *walletService) UserWalletData(id uint) (*dto.WalletDataRes, error) {
	tx := w.db.Begin()
	var err error
	defer helper.CommitOrRollback(tx, &err)

	wallet, err := w.walletRepository.GetWalletByUserID(tx, id)
	if err != nil {
		return nil, err
	}
	transactions, err := w.walletRepository.GetTransactionsByUserID(tx, id)
	var status string
	if wallet.Pin == nil {
		status = "Pin has not been set"
	} else {
		status = "Pin has been set"
	}
	walletData := &dto.WalletDataRes{
		UserID:       2,
		Balance:      wallet.Balance,
		Status:       &status,
		Transactions: transactions,
	}

	return walletData, nil
}

func (w *walletService) TransactionDetails(userID uint, transactionID uint) (*dto.TransactionDetailsRes, error) {
	tx := w.db.Begin()
	var err error
	defer helper.CommitOrRollback(tx, &err)

	t, err := w.walletRepository.TransactionDetails(tx, transactionID)
	if err != nil {
		return nil, err
	}

	if t.UserID != userID {
		err = apperror.UnauthorizedError("Cannot access another user transactions")
		return nil, err
	}

	transaction := &dto.TransactionDetailsRes{
		Id:            t.ID,
		VoucherID:     t.VoucherID,
		Voucher:       t.Voucher,
		Total:         t.Total,
		PaymentMethod: t.PaymentMethod,
		Orders:        t.Orders,
		CreatedAt:     t.CreatedAt,
		UpdatedAt:     t.UpdatedAt,
	}

	return transaction, nil
}

func (w *walletService) PaginatedTransactions(q *repository.Query, userID uint) (*dto.PaginatedTransactionsRes, error) {
	if q.Limit == "" {
		q.Limit = "10"
	}
	if q.Page == "" {
		q.Page = "1"
	}

	tx := w.db.Begin()
	var err error
	defer helper.CommitOrRollback(tx, &err)

	status, err := w.walletRepository.GetWalletStatus(tx, userID)
	if err != nil {
		return nil, err
	}
	if status == repository.WalletBlocked {
		return nil, apperror.BadRequestError("Wallet is currently blocked")
	}

	var ts = make([]dto.TransactionsRes, 0)
	l, t, err := w.walletRepository.PaginatedTransactions(tx, q, userID)
	if err != nil {
		return nil, err
	}

	for _, transaction := range t {
		tr := new(dto.TransactionsRes).FromTransaction(transaction)
		ts = append(ts, *tr)
	}
	limit, _ := strconv.Atoi(q.Limit)
	page, _ := strconv.Atoi(q.Page)
	totalPage := float64(l) / float64(limit)
	paginatedTransactions := dto.PaginatedTransactionsRes{
		TotalLength:  l,
		TotalPage:    int(math.Ceil(totalPage)),
		CurrentPage:  page,
		Limit:        limit,
		Transactions: ts,
	}

	return &paginatedTransactions, nil
}

func (w *walletService) GetWalletTransactionsByUserID(q *dto.WalletTransactionsQuery, userID uint) ([]*model.WalletTransaction, int64, int64, error) {
	tx := w.db.Begin()
	var err error
	defer helper.CommitOrRollback(tx, &err)

	wallet, err := w.walletRepository.GetWalletByUserID(tx, userID)
	if err != nil {
		return nil, 0, 0, err
	}

	transactions, totalPage, totalData, err := w.walletTransRepo.GetTransactionsByWalletID(tx, q, wallet.ID)
	if err != nil {
		return nil, 0, 0, err
	}

	return transactions, totalPage, totalData, nil
}

func (w *walletService) WalletPin(userID uint, pin string) error {
	tx := w.db.Begin()
	var err error
	defer helper.CommitOrRollback(tx, &err)

	if len(pin) != 6 {
		err = apperror.BadRequestError("Pin has to be 6 digits long")
		return err
	}
	err = w.walletRepository.WalletPin(tx, userID, pin)
	if err != nil {
		return err
	}

	return nil
}

func (w *walletService) RequestPinChangeWithEmail(userID uint) (string, string, error) {
	tx := w.db.Begin()
	var err error
	defer helper.CommitOrRollback(tx, &err)

	user, err := w.userRepository.GetUserByID(tx, userID)
	if err != nil {
		return "", "", err
	}

	wallet, err := w.walletRepository.GetWalletByUserID(tx, userID)
	if err != nil {
		return "", "", err
	}

	if wallet.Pin == nil {
		err = apperror.NotFoundError("Pin is not setup yet")
		return "", "", err
	}

	randomString := helper.RandomString(12)
	code := helper.RandomString(6)
	err = w.walletRepository.RequestChangePinByEmail(user.ID, randomString, code)
	if err != nil {
		return "", "", err
	}

	html := "<p>Berikut adalah kode untuk reset pin kamu:</p><h3>" + code + "</h3>"

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("ap-southeast-1")},
	)

	sender := config.Config.AWSMail
	textBody := "This email is for SeaDeals wallet pin email verification"
	subject := "SeaDeals Wallet PIN Verification"
	charSet := "UTF-8"

	svc := ses.New(sess)

	input := &ses.SendEmailInput{
		Destination: &ses.Destination{
			CcAddresses: []*string{},
			ToAddresses: []*string{
				aws.String(user.Email),
			},
		},
		Message: &ses.Message{
			Body: &ses.Body{
				Html: &ses.Content{
					Charset: aws.String(charSet),
					Data:    aws.String(html),
				},
				Text: &ses.Content{
					Charset: aws.String(charSet),
					Data:    aws.String(textBody),
				},
			},
			Subject: &ses.Content{
				Charset: aws.String(charSet),
				Data:    aws.String(subject),
			},
		},
		Source: aws.String(sender),
	}

	_, err = svc.SendEmail(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case ses.ErrCodeMessageRejected:
				return "", "", apperror.InternalServerError("Email is rejected")
			case ses.ErrCodeMailFromDomainNotVerifiedException:
				return "", "", apperror.InternalServerError("Email is not verified")
			case ses.ErrCodeConfigurationSetDoesNotExistException:
				return "", "", apperror.InternalServerError("Configuration is not yet set")
			default:
				return "", "", aerr
			}
		} else {
			return "", "", err
		}
	}

	return user.Email, randomString, nil
}

func (w *walletService) ValidateRequestIsValid(userID uint, key string) (string, error) {
	err := w.walletRepository.ValidateRequestIsValid(userID, key)
	if err != nil {
		return "Request is invalid", err
	}

	return "Request is valid", nil
}

func (w *walletService) ValidateCodeToRequestByEmail(userID uint, req *dto.CodeKeyRequestByEmailReq) (string, error) {
	err := w.walletRepository.ValidateRequestByEmailCodeIsValid(userID, req)
	if err != nil {
		return "Request is invalid", err
	}

	return "Request is valid", nil
}

func (w *walletService) ChangeWalletPinByEmail(userID uint, req *dto.ChangePinByEmailReq) (*model.Wallet, error) {
	tx := w.db.Begin()
	var err error
	defer helper.CommitOrRollback(tx, &err)

	if len(req.Pin) != 6 {
		err = apperror.BadRequestError("Pin has to be 6 digits long")
		return nil, err
	}

	wallet, err := w.walletRepository.GetWalletByUserID(tx, userID)
	if err != nil {
		return nil, err
	}

	if wallet.Pin == nil {
		err = apperror.NotFoundError("Pin is not setup yet")
		return nil, err
	}

	result, err := w.walletRepository.ChangeWalletPinByEmail(tx, userID, wallet.ID, req)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (w *walletService) ValidateWalletPin(user *dto.UserJWT, pin string) (string, bool, error) {
	tx := w.db.Begin()
	var err error
	defer helper.CommitOrRollback(tx, &err)

	if len(pin) != 6 {
		err = apperror.BadRequestError("Pin has to be 6 digits long")
		return "", false, err
	}

	err = w.walletRepository.ValidateWalletPin(tx, user.UserID, pin)
	if err != nil {
		return "", false, err
	}

	userRoles, err := w.userRoleRepo.GetRolesByUserID(tx, user.UserID)
	if err != nil {
		return "", false, err
	}
	var roles []string
	for _, role := range userRoles {
		roles = append(roles, role.Role.Name)
	}
	rolesString := strings.Join(roles[:], " ")
	rolesString += " level1"

	idToken, err := helper.GenerateJWTToken(user, rolesString, config.Config.JWTExpiredInMinuteTime*60, dto.JWTAccessToken)
	if err != nil {
		return "", false, err
	}

	return idToken, true, nil
}

func (w *walletService) GetWalletStatus(userID uint) (string, error) {
	tx := w.db.Begin()
	var err error
	defer helper.CommitOrRollback(tx, &err)

	status, err := w.walletRepository.GetWalletStatus(tx, userID)
	if err != nil {
		return "", err
	}

	return status, nil
}

func (w *walletService) PayOrderWithWallet(userID uint, req *dto.CheckoutCartReq) (*dto.CheckoutCartRes, error) {
	tx := w.db.Begin()
	var err error
	defer helper.CommitOrRollback(tx, &err)

	status, err := w.walletRepository.GetWalletStatus(tx, userID)
	if err != nil {
		return nil, err
	}
	if status == repository.WalletBlocked {
		err = apperror.BadRequestError("Wallet is currently blocked")
		return nil, err
	}

	if len(req.Cart) <= 0 {
		err = apperror.BadRequestError("Checkout setidaknya harus terdapat satu barang")
		return nil, err
	}

	globalVoucher, err := w.walletRepository.GetVoucher(tx, req.GlobalVoucherCode)
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
	//create transaction
	var transaction = &model.Transaction{
		UserID:        userID,
		VoucherID:     voucherID,
		Total:         0,
		PaymentMethod: dto.Wallet,
		Status:        dto.TransactionPayed,
	}

	transaction, err = w.walletRepository.CreateTransaction(tx, transaction)
	if err != nil {
		return nil, err
	}

	var totalOrderPrice float64
	var totalDelivery float64
	var sellerIDs []uint
	for _, item := range req.Cart {
		for _, id := range sellerIDs {
			if id == item.SellerID {
				err = apperror.BadRequestError("Tidak bisa membuat 2 order dengan seller yang sama dalam satu transaksi")
				return nil, err
			}
		}

		//check voucher if voucher still valid
		var voucher *model.Voucher
		voucher, err = w.walletRepository.GetVoucher(tx, item.VoucherCode)
		if err != nil {
			return nil, err
		}

		var order *model.Order
		if voucher != nil {
			if timeNow.After(voucher.EndDate) || timeNow.Before(voucher.StartDate) {
				err = apperror.InternalServerError("Level 2 Voucher invalid")
				return nil, err
			}
			order, err = w.walletRepository.CreateOrder(tx, item.SellerID, &voucher.ID, transaction.ID, userID)

			if err != nil {
				return nil, err
			}

		} else {
			//create order before order_items
			order, err = w.walletRepository.CreateOrder(tx, item.SellerID, nil, transaction.ID, userID)
			if err != nil {
				return nil, err
			}
		}
		var totalOrder float64
		var totalWeight int

		for _, id := range item.CartItemID {
			var totalOrderItem float64
			var cartItem *model.CartItem
			cartItem, err = w.walletRepository.GetCartItem(tx, id)
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

			// update stock
			err = w.walletRepository.UpdateStock(tx, cartItem.ProductVariantDetail, uint(newStock))
			if err != nil {
				return nil, err
			}

			//1. create order item and remove cart
			err = w.walletRepository.CreateOrderItemAndRemoveFromCart(tx, cartItem.ProductVariantDetailID, cartItem.ProductVariantDetail.Product, order.ID, userID, cartItem.Quantity, totalOrderItem, cartItem)
			if err != nil {
				return nil, err
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
		seller, err = w.sellerRepository.FindSellerByID(tx, item.SellerID)
		if err != nil {
			return nil, err
		}

		// Create delivery
		var courier *model.Courier
		courier, err = w.courierRepository.GetCourierDetailByID(tx, item.CourierID)
		if err != nil {
			return nil, err
		}

		var buyerAddress *model.Address
		buyerAddress, err = w.addressRepository.CheckUserAddress(tx, req.BuyerAddressID, userID)
		if err != nil {
			return nil, err
		}

		deliveryReq := &dto.DeliveryCalculateReq{
			OriginCity:      seller.Address.CityID,
			DestinationCity: buyerAddress.CityID,
			Weight:          strconv.Itoa(totalWeight),
			Courier:         courier.Code,
		}

		var deliveryResult = &dto.DeliveryCalculateReturn{}
		deliveryResult, err = helper.CalculateDeliveryPrice(deliveryReq)
		if err != nil {
			return nil, err
		}

		delivery := &model.Delivery{
			Address:         buyerAddress.Address + ", " + buyerAddress.Province + ", " + buyerAddress.City + ", " + buyerAddress.SubDistrict + ", " + buyerAddress.PostalCode,
			Status:          dto.DeliveryWaitingForSeller,
			DeliveryNumber:  helper.RandomString(10),
			Total:           float64(deliveryResult.Total),
			Eta:             deliveryResult.Eta,
			OrderID:         order.ID,
			CourierID:       courier.ID,
			CityDestination: buyerAddress.City,
			Weight:          uint(totalWeight),
		}
		newDelivery := &model.Delivery{}
		newDelivery, err = w.deliveryRepo.CreateDelivery(tx, delivery)
		if err != nil {
			return nil, err
		}
		_, err = w.deliveryActRepo.CreateActivity(tx, newDelivery.ID, "Process dibuat dan menunggu pembayaran dari buyer")
		if err != nil {
			return nil, err
		}

		//update order price with map - voucher id
		order.Total = totalOrder
		order.Status = dto.OrderWaitingSeller
		err = w.walletRepository.UpdateOrder(tx, order)
		if err != nil {
			return nil, err
		}

		totalOrderPrice += totalOrder
		totalDelivery += delivery.Total
		sellerIDs = append(sellerIDs, item.SellerID)

		accountHolder := &model.AccountHolder{
			UserID:   userID,
			OrderID:  order.ID,
			SellerID: seller.ID,
			Total:    totalOrder,
			HasTaken: false,
		}
		_, err = w.accountHolderRepo.SendToAccountHolder(tx, accountHolder)
		if err != nil {
			return nil, err
		}
	}
	//total transaction - voucher
	//4. check user wallet balance is sufficient
	wallet, err := w.walletRepository.GetWalletByUserID(tx, userID)
	if err != nil {
		return nil, err
	}
	if globalVoucher != nil && globalVoucher.SellerID == nil && globalVoucher.MinSpending <= totalOrderPrice {
		if globalVoucher.AmountType == "percentage" {
			totalOrderPrice -= (globalVoucher.Amount / 100) * totalOrderPrice
		} else {
			totalOrderPrice -= globalVoucher.Amount
		}
	} else if globalVoucher != nil {
		err = apperror.BadRequestError("Order tidak memenuhi kriteria voucher global")
		return nil, err
	}
	if totalOrderPrice < 0 {
		totalOrderPrice = 0
	}

	var totalTransaction float64
	totalTransaction = totalOrderPrice + totalDelivery

	if wallet.Balance-totalTransaction < 0 {
		err = apperror.InternalServerError("Insufficient Balance")
		return nil, err
	}
	//5. update transaction
	transaction.Total = totalTransaction
	fmt.Println(transaction.VoucherID)
	err = w.walletRepository.UpdateTransaction(tx, transaction)
	if err != nil {
		return nil, err
	}

	if req.PaymentMethod == dto.Wallet {
		err = w.walletRepository.CreateWalletTransaction(tx, wallet.ID, transaction)
		if err != nil {
			return nil, err
		}
		err = w.walletRepository.UpdateWalletBalance(tx, wallet, totalTransaction)
		if err != nil {
			return nil, err
		}
	}
	//6. create response
	transRes := dto.CheckoutCartRes{
		UserID:        userID,
		TransactionID: transaction.ID,
		Total:         transaction.Total,
		PaymentMethod: transaction.PaymentMethod,
		CreatedAt:     transaction.CreatedAt,
	}
	return &transRes, nil
}
