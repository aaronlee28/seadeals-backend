package service

import (
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
	RegisterSeaLabsPayAccount(req *dto.RegisterSeaLabsPayReq) (*dto.RegisterSeaLabsPayRes, error)
	CheckSeaLabsAccountExists(req *dto.CheckSeaLabsPayReq) (*dto.CheckSeaLabsPayRes, error)
	UpdateSeaLabsAccountToMain(req *dto.UpdateSeaLabsPayToMainReq) (*model.UserSealabsPayAccount, error)
	GetSeaLabsAccountByUserID(userID uint) ([]*model.UserSealabsPayAccount, error)

	PayWithSeaLabsPay(userID uint, req *dto.CheckoutCartReq) (string, *model.SeaLabsPayTransactionHolder, error)
	PayWithSeaLabsPayCallback(txnID uint, status string) (*model.Transaction, error)
	TopUpWithSeaLabsPay(amount float64, userID uint, accountNumber string) (*model.SeaLabsPayTopUpHolder, string, error)
	TopUpWithSeaLabsPayCallback(txnID uint, status string) (*model.WalletTransaction, error)
}

type userSeaPayAccountServ struct {
	db                          *gorm.DB
	userSeaPayAccountRepo       repository.UserSeaPayAccountRepo
	seaLabsPayTopUpHolderRepo   repository.SeaLabsPayTopUpHolderRepository
	seaLabsPayTransactionHolder repository.SeaLabsPayTransactionHolderRepository
	walletRepository            repository.WalletRepository
	walletTransactionRepo       repository.WalletTransactionRepository
}

type UserSeaPayAccountServConfig struct {
	DB                          *gorm.DB
	UserSeaPayAccountRepo       repository.UserSeaPayAccountRepo
	SeaLabsPayTopUpHolderRepo   repository.SeaLabsPayTopUpHolderRepository
	SeaLabsPayTransactionHolder repository.SeaLabsPayTransactionHolderRepository
	WalletRepository            repository.WalletRepository
	WalletTransactionRepo       repository.WalletTransactionRepository
}

func NewUserSeaPayAccountServ(c *UserSeaPayAccountServConfig) UserSeaPayAccountServ {
	return &userSeaPayAccountServ{
		db:                          c.DB,
		userSeaPayAccountRepo:       c.UserSeaPayAccountRepo,
		seaLabsPayTopUpHolderRepo:   c.SeaLabsPayTopUpHolderRepo,
		seaLabsPayTransactionHolder: c.SeaLabsPayTransactionHolder,
		walletRepository:            c.WalletRepository,
		walletTransactionRepo:       c.WalletTransactionRepo,
	}
}

func (u *userSeaPayAccountServ) CheckSeaLabsAccountExists(req *dto.CheckSeaLabsPayReq) (*dto.CheckSeaLabsPayRes, error) {
	tx := u.db.Begin()
	var err error
	defer helper.CommitOrRollback(tx, &err)

	hasExists, err := u.userSeaPayAccountRepo.HasExistsSeaLabsPayAccountWith(tx, req.UserID, req.AccountNumber)
	if err != nil {
		return nil, err
	}

	response := &dto.CheckSeaLabsPayRes{IsExists: hasExists}
	return response, nil
}

func (u *userSeaPayAccountServ) RegisterSeaLabsPayAccount(req *dto.RegisterSeaLabsPayReq) (*dto.RegisterSeaLabsPayRes, error) {
	tx := u.db.Begin()
	var err error
	defer helper.CommitOrRollback(tx, &err)

	hasExists, err := u.userSeaPayAccountRepo.HasExistsSeaLabsPayAccountWith(tx, req.UserID, req.AccountNumber)
	if err != nil {
		return nil, err
	}

	if hasExists {
		err = apperror.BadRequestError("Sea Labs PayWithSeaLabsPay Account is already registered")
		return nil, err
	}

	seaLabsPayAccount, err := u.userSeaPayAccountRepo.RegisterSeaLabsPayAccount(tx, req)
	if err != nil {
		return nil, err
	}

	response := &dto.RegisterSeaLabsPayRes{
		Status:         "Completed",
		SeaLabsAccount: seaLabsPayAccount,
	}
	return response, nil
}

func (u *userSeaPayAccountServ) UpdateSeaLabsAccountToMain(req *dto.UpdateSeaLabsPayToMainReq) (*model.UserSealabsPayAccount, error) {
	tx := u.db.Begin()
	var err error
	defer helper.CommitOrRollback(tx, &err)

	updatedData, err := u.userSeaPayAccountRepo.UpdateSeaLabsPayAccountToMain(tx, req)
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

	hasAccount, err := u.userSeaPayAccountRepo.HasExistsSeaLabsPayAccountWith(tx, userID, req.AccountNumber)
	if err != nil {
		return "", nil, err
	}
	if !hasAccount {
		err = apperror.BadRequestError("That sea labs pay account is not registered in your account")
		return "", nil, err
	}

	globalVoucher, err := u.walletRepository.GetVoucher(tx, req.GlobalVoucherCode)
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
		PaymentMethod: dto.SEA_LABS_PAY,
		Status:        "Waiting for payment",
	}

	transaction, err = u.walletRepository.CreateTransaction(tx, transaction)
	if err != nil {
		return "", nil, err
	}

	var totalTransaction float64
	for _, item := range req.Cart {
		//check voucher if voucher still valid
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
		var totalOrder float64

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

			if cartItem.ProductVariantDetail.Product.Promotion != nil {
				totalOrderItem = (cartItem.ProductVariantDetail.Price - cartItem.ProductVariantDetail.Product.Promotion.Amount) * float64(cartItem.Quantity)
			} else {
				totalOrderItem = cartItem.ProductVariantDetail.Price * float64(cartItem.Quantity)
			}
			totalOrder += totalOrderItem

			// update stock
			err = u.walletRepository.UpdateStock(tx, cartItem.ProductVariantDetail, uint(newStock))
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
		if voucher != nil {
			totalOrder -= voucher.Amount
		}
		//update order price with map - voucher id
		err = u.walletRepository.UpdateOrder(tx, order, totalOrder)
		if err != nil {
			return "", nil, err
		}
		totalTransaction += totalOrder
	}

	transaction.Total = totalTransaction
	err = u.walletRepository.UpdateTransaction(tx, transaction)
	if err != nil {
		return "", nil, err
	}

	merchantCode := config.Config.SeaLabsPayMerchantCode
	apiKey := config.Config.SeaLabsPayAPIKey
	combinedString := req.AccountNumber + ":" + strconv.Itoa(int(totalTransaction)) + ":" + merchantCode

	sign := helper.GenerateHMACSHA256(combinedString, apiKey)
	redirectURL, txnID, err := helper.TransactionToSeaLabsPay(req.AccountNumber, strconv.Itoa(int(totalTransaction)), sign, "/order/pay/sea-labs-pay/callback")
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
			PaymentMethod: dto.SEA_LABS_PAY,
			Status:        "Waiting for seller",
		}
		err = u.walletRepository.UpdateTransaction(tx, transactionModel)
		if err != nil {
			return nil, err
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
	redirectURL, txnId, err := helper.TransactionToSeaLabsPay(accountNumber, amountString, sign, "/user/wallet/top-up/sea-labs-pay/callback")
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
			PaymentMethod: dto.SEA_LABS_PAY,
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
