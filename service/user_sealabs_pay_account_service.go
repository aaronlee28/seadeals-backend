package service

import (
	"fmt"
	"gorm.io/gorm"
	"seadeals-backend/apperror"
	"seadeals-backend/config"
	"seadeals-backend/dto"
	"seadeals-backend/helper"
	"seadeals-backend/model"
	"seadeals-backend/repository"
	"strconv"
	"time"
)

type UserSeaPayAccountServ interface {
	RegisterSeaLabsPayAccount(req *dto.RegisterSeaLabsPayReq, userID uint) (*dto.RegisterSeaLabsPayRes, error)
	CheckSeaLabsAccountExists(req *dto.CheckSeaLabsPayReq, userID uint) (*dto.CheckSeaLabsPayRes, error)
	UpdateSeaLabsAccountToMain(req *dto.UpdateSeaLabsPayToMainReq, userID uint) (*model.UserSealabsPayAccount, error)
	GetSeaLabsAccountByUserID(userID uint) ([]*model.UserSealabsPayAccount, error)

	PayWithSeaLabsPay(userID uint, req *dto.CheckoutCartReq) (string, *model.SeaLabsPayTransactionHolder, error)
	PayWithSeaLabsPayCallback(txnID uint, status string) (*model.Transaction, error)
	TopUpWithSeaLabsPay(amount float64, userID uint, accountNumber string) (*model.SeaLabsPayTopUpHolder, string, error)
	TopUpWithSeaLabsPayCallback(txnID uint, status string) (*model.WalletTransaction, error)
}

type userSeaPayAccountServ struct {
	db                          *gorm.DB
	accountHolderRepo           repository.AccountHolderRepository
	addressRepository           repository.AddressRepository
	userSeaPayAccountRepo       repository.UserSeaPayAccountRepo
	courierRepository           repository.CourierRepository
	deliveryRepo                repository.DeliveryRepository
	deliveryActRepo             repository.DeliveryActivityRepository
	orderRepo                   repository.OrderRepository
	seaLabsPayTopUpHolderRepo   repository.SeaLabsPayTopUpHolderRepository
	sellerRepository            repository.SellerRepository
	seaLabsPayTransactionHolder repository.SeaLabsPayTransactionHolderRepository
	walletRepository            repository.WalletRepository
	walletTransactionRepo       repository.WalletTransactionRepository
}

type UserSeaPayAccountServConfig struct {
	DB                          *gorm.DB
	AccountHolderRepo           repository.AccountHolderRepository
	AddressRepository           repository.AddressRepository
	UserSeaPayAccountRepo       repository.UserSeaPayAccountRepo
	CourierRepository           repository.CourierRepository
	DeliveryRepo                repository.DeliveryRepository
	DeliveryActRepo             repository.DeliveryActivityRepository
	OrderRepo                   repository.OrderRepository
	SeaLabsPayTopUpHolderRepo   repository.SeaLabsPayTopUpHolderRepository
	SeaLabsPayTransactionHolder repository.SeaLabsPayTransactionHolderRepository
	SellerRepository            repository.SellerRepository
	WalletRepository            repository.WalletRepository
	WalletTransactionRepo       repository.WalletTransactionRepository
}

func NewUserSeaPayAccountServ(c *UserSeaPayAccountServConfig) UserSeaPayAccountServ {
	return &userSeaPayAccountServ{
		db:                          c.DB,
		accountHolderRepo:           c.AccountHolderRepo,
		addressRepository:           c.AddressRepository,
		userSeaPayAccountRepo:       c.UserSeaPayAccountRepo,
		courierRepository:           c.CourierRepository,
		deliveryRepo:                c.DeliveryRepo,
		deliveryActRepo:             c.DeliveryActRepo,
		orderRepo:                   c.OrderRepo,
		seaLabsPayTopUpHolderRepo:   c.SeaLabsPayTopUpHolderRepo,
		sellerRepository:            c.SellerRepository,
		seaLabsPayTransactionHolder: c.SeaLabsPayTransactionHolder,
		walletRepository:            c.WalletRepository,
		walletTransactionRepo:       c.WalletTransactionRepo,
	}
}

func (u *userSeaPayAccountServ) CheckSeaLabsAccountExists(req *dto.CheckSeaLabsPayReq, userID uint) (*dto.CheckSeaLabsPayRes, error) {
	tx := u.db.Begin()
	var err error
	defer helper.CommitOrRollback(tx, &err)

	merchantCode := config.Config.SeaLabsPayMerchantCode
	apiKey := config.Config.SeaLabsPayAPIKey
	combinedString := req.AccountNumber + ":" + strconv.Itoa(1) + ":" + merchantCode

	sign := helper.GenerateHMACSHA256(combinedString, apiKey)
	_, _, err = helper.TransactionToSeaLabsPay(req.AccountNumber, strconv.Itoa(1), sign, "/order/pay/sea-labs-pay/callback", "trx")
	fmt.Println("hahahh", err)
	if err != nil {
		if err.Error() == "user not found" {
			return nil, apperror.BadRequestError("Invalid Sea Labs Pay Account")
		} else if err.Error() != "insufficient fund to create transaction" {
			return nil, err
		}
	}

	hasExists, err := u.userSeaPayAccountRepo.HasExistsSeaLabsPayAccountWith(tx, userID, req.AccountNumber)
	if err != nil {
		return nil, err
	}

	response := &dto.CheckSeaLabsPayRes{IsExists: hasExists}
	return response, nil
}

func (u *userSeaPayAccountServ) RegisterSeaLabsPayAccount(req *dto.RegisterSeaLabsPayReq, userID uint) (*dto.RegisterSeaLabsPayRes, error) {
	tx := u.db.Begin()
	var err error
	defer helper.CommitOrRollback(tx, &err)

	merchantCode := config.Config.SeaLabsPayMerchantCode
	apiKey := config.Config.SeaLabsPayAPIKey
	combinedString := req.AccountNumber + ":" + strconv.Itoa(1) + ":" + merchantCode

	sign := helper.GenerateHMACSHA256(combinedString, apiKey)
	_, _, err = helper.TransactionToSeaLabsPay(req.AccountNumber, strconv.Itoa(1), sign, "/order/pay/sea-labs-pay/callback", "trx")
	if err != nil {
		if err.Error() == "user not found" {
			return nil, apperror.BadRequestError("Invalid Sea Labs Pay Account")
		} else if err.Error() != "insufficient fund to create transaction" {
			return nil, err
		}
	}

	hasExists, err := u.userSeaPayAccountRepo.HasExistsSeaLabsPayAccountWith(tx, userID, req.AccountNumber)
	if err != nil {
		return nil, err
	}

	if hasExists {
		err = apperror.BadRequestError("Sea Labs PayWithSeaLabsPay Account is already registered")
		return nil, err
	}

	seaLabsPayAccount, err := u.userSeaPayAccountRepo.RegisterSeaLabsPayAccount(tx, req, userID)
	if err != nil {
		return nil, err
	}

	response := &dto.RegisterSeaLabsPayRes{
		Status:         "Completed",
		SeaLabsAccount: seaLabsPayAccount,
	}
	return response, nil
}

func (u *userSeaPayAccountServ) UpdateSeaLabsAccountToMain(req *dto.UpdateSeaLabsPayToMainReq, userID uint) (*model.UserSealabsPayAccount, error) {
	tx := u.db.Begin()
	var err error
	defer helper.CommitOrRollback(tx, &err)

	updatedData, err := u.userSeaPayAccountRepo.UpdateSeaLabsPayAccountToMain(tx, req, userID)
	if err != nil {
		return nil, err
	}

	return updatedData, nil
}

func (u *userSeaPayAccountServ) GetSeaLabsAccountByUserID(userID uint) ([]*model.UserSealabsPayAccount, error) {
	tx := u.db.Begin()
	var err error
	defer helper.CommitOrRollback(tx, &err)

	accounts, err := u.userSeaPayAccountRepo.GetSeaLabsPayAccountByUserID(tx, userID)
	if err != nil {
		return nil, err
	}

	return accounts, nil
}

func (u *userSeaPayAccountServ) PayWithSeaLabsPay(userID uint, req *dto.CheckoutCartReq) (string, *model.SeaLabsPayTransactionHolder, error) {
	tx := u.db.Begin()
	var err error
	defer helper.CommitOrRollback(tx, &err)

	var hasAccount bool
	hasAccount, err = u.userSeaPayAccountRepo.HasExistsSeaLabsPayAccountWith(tx, userID, req.AccountNumber)
	if err != nil {
		return "", nil, err
	}
	if !hasAccount {
		err = apperror.BadRequestError("That sea labs pay account is not registered in your account")
		return "", nil, err
	}

	if len(req.Cart) <= 0 {
		err = apperror.BadRequestError("Checkout setidaknya harus terdapat satu barang")
		return "", nil, err
	}

	var globalVoucher = &model.Voucher{}
	globalVoucher, err = u.walletRepository.GetVoucher(tx, req.GlobalVoucherCode)
	if err != nil {
		return "", nil, err
	}
	timeNow := time.Now()

	if globalVoucher != nil {
		if timeNow.After(globalVoucher.EndDate) || timeNow.Before(globalVoucher.StartDate) {
			err = apperror.InternalServerError("Level 3 Voucher invalid")
			return "", nil, err
		}
	}

	//create transaction
	var voucherID *uint
	if globalVoucher != nil {
		voucherID = &globalVoucher.ID
	}
	var transaction = &model.Transaction{
		UserID:        userID,
		VoucherID:     voucherID,
		Total:         0,
		PaymentMethod: dto.SeaLabsPay,
		Status:        dto.TransactionWaitingPayment,
	}

	transaction, err = u.walletRepository.CreateTransaction(tx, transaction)
	if err != nil {
		return "", nil, err
	}

	var totalOrderPrice float64
	var totalDelivery float64
	var sellerIDs []uint
	for _, item := range req.Cart {
		//check voucher if voucher still valid
		for _, id := range sellerIDs {
			if id == item.SellerID {
				err = apperror.BadRequestError("Tidak bisa membuat 2 order dengan seller yang sama dalam satu transaksi")
				return "", nil, err
			}
		}
		var voucher *model.Voucher
		voucher, err = u.walletRepository.GetVoucher(tx, item.VoucherCode)
		if err != nil {
			return "", nil, err
		}
		var order *model.Order
		if voucher != nil {
			if timeNow.After(voucher.EndDate) || timeNow.Before(voucher.StartDate) {
				err = apperror.InternalServerError("Level 2 Voucher invalid")
				return "", nil, err
			}

			order, err = u.walletRepository.CreateOrder(tx, item.SellerID, &voucher.ID, transaction.ID, userID)
			if err != nil {
				return "", nil, err
			}
		} else {
			//create order before order_items
			order, err = u.walletRepository.CreateOrder(tx, item.SellerID, nil, transaction.ID, userID)
			if err != nil {
				return "", nil, err
			}
		}
		var orderSubtotal float64
		var totalWeight int

		for _, id := range item.CartItemID {
			var totalOrderItem float64
			var cartItem *model.CartItem
			cartItem, err = u.walletRepository.GetCartItem(tx, id)
			if err != nil {
				return "", nil, err
			}

			if cartItem.ProductVariantDetail.Product.SellerID != item.SellerID {
				err = apperror.BadRequestError("That cart item is not belong to that seller")
				return "", nil, err
			}
			//check stock
			newStock := cartItem.ProductVariantDetail.Stock - cartItem.Quantity
			if newStock < 0 {
				err = apperror.InternalServerError(cartItem.ProductVariantDetail.Product.Name + "is out of stock")
				return "", nil, err
			}

			if cartItem.ProductVariantDetail.Product.Promotion != nil && cartItem.ProductVariantDetail.Product.Promotion.MaxOrder >= cartItem.Quantity {
				totalOrderItem = (cartItem.ProductVariantDetail.Price - cartItem.ProductVariantDetail.Product.Promotion.Amount) * float64(cartItem.Quantity)
			} else {
				totalOrderItem = cartItem.ProductVariantDetail.Price * float64(cartItem.Quantity)
			}
			if totalOrderItem < 0 {
				totalOrderItem = 0
			}
			orderSubtotal += totalOrderItem

			// Get weight
			totalWeight += int(cartItem.Quantity) * cartItem.ProductVariantDetail.Product.ProductDetail.Weight
			if totalWeight > 20000 {
				err = apperror.BadRequestError(cartItem.ProductVariantDetail.Product.Name + " exceeded weight limit of 20000")
				return "", nil, err
			}

			// update stock
			err = u.walletRepository.UpdateStock(tx, cartItem.ProductVariantDetail, newStock)
			if err != nil {
				return "", nil, err
			}

			//1. create order item and remove cart
			err = u.walletRepository.CreateOrderItemAndRemoveFromCart(tx, cartItem.ProductVariantDetailID, cartItem.ProductVariantDetail.Product, order.ID, userID, cartItem.Quantity, totalOrderItem, cartItem)
			if err != nil {
				return "", nil, err
			}
		}

		//order - voucher
		if voucher != nil && voucher.MinSpending <= orderSubtotal {
			if voucher.AmountType == "percentage" {
				orderSubtotal -= (voucher.Amount / 100) * orderSubtotal
			} else {
				orderSubtotal -= voucher.Amount
			}
		} else if voucher != nil {
			err = apperror.BadRequestError("Order tidak memenuhi kriteria voucher " + voucher.Name)
			return "", nil, err
		}
		if orderSubtotal < 0 {
			orderSubtotal = 0
		}

		var seller *model.Seller
		seller, err = u.sellerRepository.FindSellerByID(tx, item.SellerID)
		if err != nil {
			return "", nil, err
		}

		// Create delivery
		var courier *model.Courier
		courier, err = u.courierRepository.GetCourierDetailByID(tx, item.CourierID)
		if err != nil {
			return "", nil, err
		}

		var buyerAddress *model.Address
		buyerAddress, err = u.addressRepository.CheckUserAddress(tx, req.BuyerAddressID, userID)
		if err != nil {
			return "", nil, err
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
			return "", nil, err
		}

		delivery := &model.Delivery{
			Address:         buyerAddress.Address + ", " + buyerAddress.Province + ", " + buyerAddress.City + ", " + buyerAddress.SubDistrict + ", " + buyerAddress.PostalCode,
			Status:          dto.DeliveryWaitingForPayment,
			DeliveryNumber:  helper.RandomString(10),
			Eta:             deliveryResult.Eta,
			Total:           float64(deliveryResult.Total),
			OrderID:         order.ID,
			CourierID:       courier.ID,
			CityDestination: buyerAddress.City,
			Weight:          uint(totalWeight),
		}
		newDelivery := &model.Delivery{}
		newDelivery, err = u.deliveryRepo.CreateDelivery(tx, delivery)
		if err != nil {
			return "", nil, err
		}
		_, err = u.deliveryActRepo.CreateActivity(tx, newDelivery.ID, "Process dibuat dan menunggu pembayaran dari buyer")
		if err != nil {
			return "", nil, err
		}

		// Update order price with map - voucher id
		order.Total = orderSubtotal
		order.Status = dto.OrderWaitingPayment
		err = u.walletRepository.UpdateOrder(tx, order)
		if err != nil {
			return "", nil, err
		}
		totalOrderPrice += orderSubtotal
		totalDelivery += delivery.Total
		sellerIDs = append(sellerIDs, item.SellerID)

		accountHolder := &model.AccountHolder{
			UserID:   userID,
			OrderID:  order.ID,
			SellerID: seller.ID,
			Total:    orderSubtotal,
			HasTaken: false,
		}
		_, err = u.accountHolderRepo.SendToAccountHolder(tx, accountHolder)
		if err != nil {
			return "", nil, err
		}
	}

	if globalVoucher != nil && globalVoucher.SellerID == nil && globalVoucher.MinSpending <= totalOrderPrice {
		if globalVoucher.AmountType == "percentage" {
			totalOrderPrice -= (globalVoucher.Amount / 100) * totalOrderPrice
		} else {
			totalOrderPrice -= globalVoucher.Amount
		}
	} else if globalVoucher != nil {
		err = apperror.BadRequestError("Order tidak memenuhi kriteria voucher global")
		return "", nil, err
	}
	if totalOrderPrice < 0 {
		totalOrderPrice = 0
	}

	var totalTransaction float64
	totalTransaction = totalOrderPrice + totalDelivery

	transaction.Total = totalTransaction
	err = u.walletRepository.UpdateTransaction(tx, transaction)
	if err != nil {
		return "", nil, err
	}

	merchantCode := config.Config.SeaLabsPayMerchantCode
	apiKey := config.Config.SeaLabsPayAPIKey
	combinedString := req.AccountNumber + ":" + strconv.Itoa(int(totalTransaction)) + ":" + merchantCode

	sign := helper.GenerateHMACSHA256(combinedString, apiKey)
	redirectURL, txnID, err := helper.TransactionToSeaLabsPay(req.AccountNumber, strconv.Itoa(int(totalTransaction)), sign, "/order/pay/sea-labs-pay/callback", "trx")
	if err != nil {
		return "", nil, err
	}

	newData := &model.SeaLabsPayTransactionHolder{
		UserID:        userID,
		TransactionID: transaction.ID,
		TxnID:         txnID,
		Sign:          sign,
		Total:         totalTransaction,
	}

	holder, err := u.seaLabsPayTransactionHolder.CreateTransactionHolder(tx, newData)
	if err != nil {
		return "", nil, err
	}

	return redirectURL, holder, nil
}

func (u *userSeaPayAccountServ) PayWithSeaLabsPayCallback(txnID uint, status string) (*model.Transaction, error) {
	tx := u.db.Begin()
	var err error
	defer helper.CommitOrRollback(tx, &err)

	transactionHolder, err := u.seaLabsPayTransactionHolder.UpdateTransactionHolder(tx, txnID, status)
	if err != nil {
		return nil, err
	}

	var transaction *model.Transaction
	if status == dto.TXN_PAID {
		transactionModel := &model.Transaction{
			ID:            transactionHolder.TransactionID,
			PaymentMethod: dto.SeaLabsPay,
			Status:        dto.TransactionPayed,
		}
		err = u.walletRepository.UpdateTransaction(tx, transactionModel)
		if err != nil {
			return nil, err
		}

		var orders []*model.Order
		orders, err = u.orderRepo.UpdateOrderStatusByTransID(tx, transactionHolder.TransactionID, dto.OrderWaitingSeller)
		if err != nil {
			return nil, err
		}
		for _, order := range orders {
			_, _ = u.deliveryRepo.UpdateDeliveryStatus(tx, order.Delivery.ID, dto.DeliveryWaitingForSeller)
			_, _ = u.deliveryActRepo.CreateActivity(tx, order.Delivery.ID, "Pembayaran diterima dan menunggu Seller")
		}
	}

	return transaction, nil
}

func (u *userSeaPayAccountServ) TopUpWithSeaLabsPay(amount float64, userID uint, accountNumber string) (*model.SeaLabsPayTopUpHolder, string, error) {
	tx := u.db.Begin()
	var err error
	defer helper.CommitOrRollback(tx, &err)

	status, err := u.walletRepository.GetWalletStatus(tx, userID)
	if err != nil {
		return nil, "", err
	}
	if status == repository.WalletBlocked {
		return nil, "", apperror.BadRequestError("Wallet is currently blocked")
	}

	hasAccount, err := u.userSeaPayAccountRepo.HasExistsSeaLabsPayAccountWith(tx, userID, accountNumber)
	if err != nil {
		return nil, "", err
	}
	if !hasAccount {
		err = apperror.BadRequestError("That sea labs pay account is not registered in your account")
		return nil, "", err
	}

	merchantCode := config.Config.SeaLabsPayMerchantCode
	apiKey := config.Config.SeaLabsPayAPIKey
	amountString := strconv.Itoa(int(amount))
	combinedString := accountNumber + ":" + amountString + ":" + merchantCode

	sign := helper.GenerateHMACSHA256(combinedString, apiKey)
	redirectURL, txnId, err := helper.TransactionToSeaLabsPay(accountNumber, amountString, sign, "/user/wallet/top-up/sea-labs-pay/callback", "topup")
	if err != nil {
		return nil, "", err
	}

	newData := &model.SeaLabsPayTopUpHolder{
		UserID: userID,
		TxnID:  txnId,
		Total:  amount,
		Sign:   sign,
	}

	holder, err := u.seaLabsPayTopUpHolderRepo.CreateTopUpHolder(tx, newData)
	if err != nil {
		return nil, "", err
	}
	return holder, redirectURL, nil
}

func (u *userSeaPayAccountServ) TopUpWithSeaLabsPayCallback(txnID uint, status string) (*model.WalletTransaction, error) {
	tx := u.db.Begin()
	var err error
	defer helper.CommitOrRollback(tx, &err)

	topUpHolder, err := u.seaLabsPayTopUpHolderRepo.UpdateTopUpHolder(tx, txnID, status)
	if err != nil {
		return nil, err
	}

	var transaction *model.WalletTransaction
	if status == dto.TXN_PAID {
		var wallet *model.Wallet
		wallet, err = u.walletRepository.GetWalletByUserID(tx, topUpHolder.UserID)
		if err != nil {
			return nil, err
		}

		_, err = u.walletRepository.TopUp(tx, wallet, topUpHolder.Total)
		if err != nil {
			return nil, err
		}

		transactionModel := &model.WalletTransaction{
			WalletID:      wallet.ID,
			Total:         topUpHolder.Total,
			PaymentMethod: dto.SeaLabsPay,
			PaymentType:   "CREDIT",
			Description:   "Top up from Sea Labs Pay",
		}
		transaction, err = u.walletTransactionRepo.CreateTransaction(tx, transactionModel)
		if err != nil {
			return nil, err
		}
	}

	return transaction, nil
}
