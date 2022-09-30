package service

import (
	"math"
	"seadeals-backend/apperror"
	"seadeals-backend/config"
	"seadeals-backend/dto"
	"seadeals-backend/helper"
	"seadeals-backend/model"
	"seadeals-backend/repository"
	"strconv"
	"time"

	"github.com/mailjet/mailjet-apiv3-go"
	"gorm.io/gorm"
)

type WalletService interface {
	UserWalletData(id uint) (*dto.WalletDataRes, error)
	TransactionDetails(userID uint, transactionID uint) (*dto.TransactionDetailsRes, error)
	PaginatedTransactions(q *repository.Query, userID uint) (*dto.PaginatedTransactionsRes, error)
	GetWalletTransactionsByUserID(q *dto.WalletTransactionsQuery, userID uint) ([]*model.WalletTransaction, int64, int64, error)

	WalletPin(userID uint, pin string) error
	RequestPinChangeWithEmail(userID uint) (*mailjet.ResultsV31, string, error)
	ValidateRequestIsValid(userID uint, key string) (string, error)
	ValidateCodeToRequestByEmail(userID uint, req *dto.CodeKeyRequestByEmailReq) (string, error)
	ChangeWalletPinByEmail(userID uint, req *dto.ChangePinByEmailReq) (*model.Wallet, error)
	ValidateWalletPin(userID uint, pin string) (bool, error)

	GetWalletStatus(userID uint) (string, error)
	CheckoutCart(userID uint, req *dto.CheckoutCartReq) (*dto.CheckoutCartRes, error)
}

type walletService struct {
	db               *gorm.DB
	walletRepository repository.WalletRepository
	userRepository   repository.UserRepository
	walletTransRepo  repository.WalletTransactionRepository
}

type WalletServiceConfig struct {
	DB               *gorm.DB
	WalletRepository repository.WalletRepository
	UserRepository   repository.UserRepository
	WalletTransRepo  repository.WalletTransactionRepository
}

func NewWalletService(c *WalletServiceConfig) WalletService {
	return &walletService{
		db:               c.DB,
		walletRepository: c.WalletRepository,
		userRepository:   c.UserRepository,
		walletTransRepo:  c.WalletTransRepo,
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
		Total:         t.Total,
		PaymentMethod: t.PaymentMethod,
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

	var ts []dto.TransactionsRes
	l, t, err := w.walletRepository.PaginatedTransactions(tx, q, userID)
	if err != nil {
		return nil, err
	}

	for _, transaction := range *t {
		tr := new(dto.TransactionsRes).FromTransaction(&transaction)
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

	if len(transactions) <= 0 {
		err = apperror.NotFoundError("No transactions were made")
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

func (w *walletService) RequestPinChangeWithEmail(userID uint) (*mailjet.ResultsV31, string, error) {
	tx := w.db.Begin()
	var err error
	defer helper.CommitOrRollback(tx, &err)

	user, err := w.userRepository.GetUserByID(tx, userID)
	if err != nil {
		return nil, "", err
	}

	wallet, err := w.walletRepository.GetWalletByUserID(tx, userID)
	if err != nil {
		return nil, "", err
	}

	if wallet.Pin == nil {
		err = apperror.NotFoundError("Pin is not setup yet")
		return nil, "", err
	}

	randomString := helper.RandomString(12)
	code := helper.RandomString(6)
	err = w.walletRepository.RequestChangePinByEmail(user.ID, randomString, code)
	if err != nil {
		return nil, "", err
	}

	mailjetClient := mailjet.NewMailjetClient(config.Config.MailJetPublicKey, config.Config.MailJetSecretKey)
	html := "<p>Berikut adalah kode untuk reset pin kamu:</p><h3>" + code + "</h3>"
	messagesInfo := []mailjet.InfoMessagesV31{
		{
			From: &mailjet.RecipientV31{
				Email: "seadeals04@gmail.com",
				Name:  "SeaDeals No Reply",
			},
			To: &mailjet.RecipientsV31{
				mailjet.RecipientV31{
					Email: user.Email,
					Name:  user.FullName,
				},
			},
			Subject:  "Wallet Pin Reset Request",
			TextPart: "request password for user" + user.FullName,
			HTMLPart: html,
			Priority: 0,
			CustomID: config.Config.AppName,
		},
	}
	messages := mailjet.MessagesV31{
		Info: messagesInfo,
	}

	res, err := mailjetClient.SendMailV31(&messages)
	if err != nil {
		return nil, "", err
	}
	return res, randomString, nil
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

func (w *walletService) ValidateWalletPin(userID uint, pin string) (bool, error) {
	tx := w.db.Begin()
	var err error
	defer helper.CommitOrRollback(tx, &err)

	if len(pin) != 6 {
		err = apperror.BadRequestError("Pin has to be 6 digits long")
		return false, err
	}

	err = w.walletRepository.ValidateWalletPin(tx, userID, pin)
	if err != nil {
		return false, err
	}

	return true, nil
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

func (w *walletService) CheckoutCart(userID uint, req *dto.CheckoutCartReq) (*dto.CheckoutCartRes, error) {
	tx := w.db.Begin()
	var err error
	defer helper.CommitOrRollback(tx, &err)

	globalVoucher, err := w.walletRepository.GetVoucher(tx, req.GlobalVoucherCode)
	if err != nil {
		return nil, err
	}
	timeNow := time.Now()

	if globalVoucher != nil {
		if timeNow.After(globalVoucher.EndDate) || timeNow.Before(globalVoucher.StartDate) {
			return nil, apperror.InternalServerError("Level 3 Voucher invalid")
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
		PaymentMethod: "wallet",
		Status:        "Waiting for Seller",
	}

	transaction, err = w.walletRepository.CreateTransaction(tx, transaction)
	if err != nil {
		return nil, err
	}

	var totalTransaction float64
	for _, item := range req.Cart {
		//check voucher if voucher still valid
		var voucher *model.Voucher
		voucher, err = w.walletRepository.GetVoucher(tx, item.VoucherCode)
		if err != nil {
			return nil, err
		}

		var order *model.Order
		if voucher != nil {
			if timeNow.After(voucher.EndDate) || timeNow.Before(voucher.StartDate) {
				return nil, apperror.InternalServerError("Level 2 Voucher invalid")
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
			newStock := cartItem.ProductVariantDetail.Stock - int(cartItem.Quantity)
			if newStock < 0 {
				err = apperror.InternalServerError(cartItem.ProductVariantDetail.Product.Name + "is out of stock")
				return nil, err
			}

			if cartItem.ProductVariantDetail.Product.Promotion != nil {
				totalOrderItem = (cartItem.ProductVariantDetail.Price - cartItem.ProductVariantDetail.Product.Promotion.Amount) * float64(cartItem.Quantity)
			} else {
				totalOrderItem = cartItem.ProductVariantDetail.Price * float64(cartItem.Quantity)
			}
			totalOrder += totalOrderItem

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
		if voucher != nil {
			totalOrder -= voucher.Amount
		}
		//update order price with map - voucher id
		err = w.walletRepository.UpdateOrder(tx, order, totalOrder)
		if err != nil {
			return nil, err
		}

		totalTransaction += totalOrder
	}
	//total transaction - voucher
	//4. check user wallet balance is sufficient
	wallet, err := w.walletRepository.GetWalletByUserID(tx, userID)
	if err != nil {
		return nil, err
	}
	if globalVoucher != nil {
		totalTransaction -= globalVoucher.Amount
	}

	if wallet.Balance-totalTransaction < 0 {
		err = apperror.InternalServerError("Insufficient Balance")
		return nil, err
	}
	//5. update transaction
	transaction.Total = totalTransaction
	err = w.walletRepository.UpdateTransaction(tx, transaction)
	if err != nil {
		return nil, err
	}

	if req.PaymentMethod == "wallet" {
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
