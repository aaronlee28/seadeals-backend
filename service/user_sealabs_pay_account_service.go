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
	RegisterSeaLabsPayAccount(req *dto.RegisterSeaLabsPayReq) (*dto.RegisterSeaLabsPayRes, error)
	CheckSeaLabsAccountExists(req *dto.CheckSeaLabsPayReq) (*dto.CheckSeaLabsPayRes, error)
	UpdateSeaLabsAccountToMain(req *dto.UpdateSeaLabsPayToMainReq) (*model.UserSealabsPayAccount, error)
	GetSeaLabsAccountByUserID(userID uint) ([]*model.UserSealabsPayAccount, error)

	PayWithSeaLabsPay(userID uint, req *dto.CheckoutCartReq) (string, int, error)
	TopUpWithSeaLabsPay(amount float64, userID uint, accountNumber string) (*model.SeaLabsPayTopUpHolder, string, error)
	TopUpWithSeaLabsPayCallback(txnID uint, status string) (*model.WalletTransaction, error)
}

type userSeaPayAccountServ struct {
	db                        *gorm.DB
	userSeaPayAccountRepo     repository.UserSeaPayAccountRepo
	seaLabsPayTopUpHolderRepo repository.SeaLabsPayTopUpHolderRepository
	walletRepository          repository.WalletRepository
	walletTransactionRepo     repository.WalletTransactionRepository
}

type UserSeaPayAccountServConfig struct {
	DB                        *gorm.DB
	UserSeaPayAccountRepo     repository.UserSeaPayAccountRepo
	SeaLabsPayTopUpHolderRepo repository.SeaLabsPayTopUpHolderRepository
	WalletRepository          repository.WalletRepository
	WalletTransactionRepo     repository.WalletTransactionRepository
}

func NewUserSeaPayAccountServ(c *UserSeaPayAccountServConfig) UserSeaPayAccountServ {
	return &userSeaPayAccountServ{
		db:                        c.DB,
		userSeaPayAccountRepo:     c.UserSeaPayAccountRepo,
		seaLabsPayTopUpHolderRepo: c.SeaLabsPayTopUpHolderRepo,
		walletRepository:          c.WalletRepository,
		walletTransactionRepo:     c.WalletTransactionRepo,
	}
}

func (u *userSeaPayAccountServ) CheckSeaLabsAccountExists(req *dto.CheckSeaLabsPayReq) (*dto.CheckSeaLabsPayRes, error) {
	tx := u.db.Begin()

	hasExists, err := u.userSeaPayAccountRepo.HasExistsSeaLabsPayAccountWith(tx, req.UserID, req.AccountNumber)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	response := &dto.CheckSeaLabsPayRes{IsExists: hasExists}

	tx.Commit()
	return response, nil
}

func (u *userSeaPayAccountServ) RegisterSeaLabsPayAccount(req *dto.RegisterSeaLabsPayReq) (*dto.RegisterSeaLabsPayRes, error) {
	tx := u.db.Begin()

	hasExists, err := u.userSeaPayAccountRepo.HasExistsSeaLabsPayAccountWith(tx, req.UserID, req.AccountNumber)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	if hasExists {
		tx.Rollback()
		return nil, apperror.BadRequestError("Sea Labs PayWithSeaLabsPay Account is already registered")
	}

	seaLabsPayAccount, err := u.userSeaPayAccountRepo.RegisterSeaLabsPayAccount(tx, req)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	response := &dto.RegisterSeaLabsPayRes{
		Status:         "Completed",
		SeaLabsAccount: seaLabsPayAccount,
	}

	tx.Commit()
	return response, nil
}

func (u *userSeaPayAccountServ) UpdateSeaLabsAccountToMain(req *dto.UpdateSeaLabsPayToMainReq) (*model.UserSealabsPayAccount, error) {
	tx := u.db.Begin()

	updatedData, err := u.userSeaPayAccountRepo.UpdateSeaLabsPayAccountToMain(tx, req)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	tx.Commit()
	return updatedData, nil
}

func (u *userSeaPayAccountServ) GetSeaLabsAccountByUserID(userID uint) ([]*model.UserSealabsPayAccount, error) {
	tx := u.db.Begin()

	accounts, err := u.userSeaPayAccountRepo.GetSeaLabsPayAccountByUserID(tx, userID)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	tx.Commit()
	return accounts, nil
}

func (u *userSeaPayAccountServ) PayWithSeaLabsPay(userID uint, req *dto.CheckoutCartReq) (string, int, error) {
	tx := u.db.Begin()

	hasAccount, err := u.userSeaPayAccountRepo.HasExistsSeaLabsPayAccountWith(tx, userID, req.AccountNumber)
	if err != nil {
		tx.Rollback()
		return "", 0, err
	}
	if !hasAccount {
		tx.Rollback()
		return "", 0, apperror.BadRequestError("That sea labs pay account is not registered in your account")
	}

	globalVoucher, err := u.walletRepository.GetVoucher(tx, req.GlobalVoucherCode)
	if err != nil {
		tx.Rollback()
		return "", 0, err
	}
	timeNow := time.Now()

	if globalVoucher != nil {
		if timeNow.After(globalVoucher.EndDate) || timeNow.Before(globalVoucher.StartDate) {
			return "", 0, apperror.InternalServerError("Level 3 Voucher invalid")
		}
	}

	//create transaction
	var transaction *model.Transaction
	var err5 error
	if globalVoucher != nil {
		transaction, err5 = u.walletRepository.CreateTransaction(tx, userID, &globalVoucher.ID)
		if err5 != nil {
			tx.Rollback()
			return "", 0, err5
		}

	} else {
		transaction, err5 = u.walletRepository.CreateTransaction(tx, userID, nil)
		if err5 != nil {
			tx.Rollback()
			return "", 0, err5
		}
	}

	var totalTransaction float64

	for _, item := range req.Cart {
		//check voucher if voucher still valid
		voucher, err1 := u.walletRepository.GetVoucher(tx, item.VoucherCode)
		if err1 != nil {
			tx.Rollback()
			return "", 0, err1
		}
		var order *model.Order
		var err6 error
		if voucher != nil {
			if timeNow.After(voucher.EndDate) || timeNow.Before(voucher.StartDate) {
				return "", 0, apperror.InternalServerError("Level 2 Voucher invalid")
			}
			order, err6 = u.walletRepository.CreateOrder(tx, item.SellerID, &voucher.ID, transaction.Id, userID)

			if err6 != nil {
				tx.Rollback()
				return "", 0, err6
			}

		} else {
			//create order before order_items
			order, err6 = u.walletRepository.CreateOrder(tx, item.SellerID, nil, transaction.Id, userID)

			if err6 != nil {
				tx.Rollback()
				return "", 0, err6
			}
		}
		var totalOrder float64

		for _, id := range item.CartItemID {
			var totalOrderItem float64
			cartItem, err2 := u.walletRepository.GetCartItem(tx, id)
			if err2 != nil {
				tx.Rollback()
				return "", 0, err2
			}

			if cartItem.ProductVariantDetail.Product.SellerID != item.SellerID {
				tx.Rollback()
				return "", 0, apperror.BadRequestError("That cart item is not belong to that seller")
			}
			//check stock
			newStock := cartItem.ProductVariantDetail.Stock - int(cartItem.Quantity)
			if newStock < 0 {
				tx.Rollback()
				return "", 0, apperror.InternalServerError(cartItem.ProductVariantDetail.Product.Name + "is out of stock")
			}
			fmt.Println("stock", id)
			if cartItem.ProductVariantDetail.Product.Promotion != nil {
				totalOrderItem = (cartItem.ProductVariantDetail.Price - cartItem.ProductVariantDetail.Product.Promotion.Amount) * float64(cartItem.Quantity)
			} else {
				totalOrderItem = cartItem.ProductVariantDetail.Price * float64(cartItem.Quantity)
			}
			totalOrder += totalOrderItem

			// update stock
			err10 := u.walletRepository.UpdateStock(tx, cartItem.ProductVariantDetail, uint(newStock))
			if err10 != nil {
				tx.Rollback()
				return "", 0, err10
			}

			//1. create order item and remove cart
			err3 := u.walletRepository.CreateOrderItemAndRemoveFromCart(tx, cartItem.ProductVariantDetailID, cartItem.ProductVariantDetail.Product, order.ID, userID, cartItem.Quantity, totalOrderItem, cartItem)
			if err3 != nil {
				tx.Rollback()
				return "", 0, err3
			}

		}

		//order - voucher
		if voucher != nil {
			totalOrder -= voucher.Amount
		}
		//update order price with map - voucher id
		err7 := u.walletRepository.UpdateOrder(tx, order, totalOrder)
		if err7 != nil {
			tx.Rollback()
			return "", 0, err7
		}

		totalTransaction += totalOrder
	}

	merchantCode := config.Config.SeaLabsPayMerchantCode
	apiKey := config.Config.SeaLabsPayAPIKey
	combinedString := req.AccountNumber + ":" + strconv.Itoa(int(totalTransaction)) + ":" + merchantCode

	sign := helper.GenerateHMACSHA256(combinedString, apiKey)
	redirectURL, _, err := helper.TransactionToSeaLabsPay(req.AccountNumber, strconv.Itoa(int(totalTransaction)), sign, "")
	if err != nil {
		tx.Rollback()
		return "", 0, err
	}

	tx.Commit()
	return redirectURL, int(totalTransaction), nil
}

func (u *userSeaPayAccountServ) TopUpWithSeaLabsPay(amount float64, userID uint, accountNumber string) (*model.SeaLabsPayTopUpHolder, string, error) {
	tx := u.db.Begin()

	hasAccount, err := u.userSeaPayAccountRepo.HasExistsSeaLabsPayAccountWith(tx, userID, accountNumber)
	if err != nil {
		return nil, "", err
	}
	if !hasAccount {
		return nil, "", apperror.BadRequestError("That sea labs pay account is not registered in your account")
	}

	merchantCode := config.Config.SeaLabsPayMerchantCode
	apiKey := config.Config.SeaLabsPayAPIKey
	amountString := strconv.Itoa(int(amount))
	combinedString := accountNumber + ":" + amountString + ":" + merchantCode

	sign := helper.GenerateHMACSHA256(combinedString, apiKey)
	redirectURL, txnId, err := helper.TransactionToSeaLabsPay(accountNumber, amountString, sign, "/user/wallet/top-up/sea-labs-pay/callback")
	if err != nil {
		tx.Rollback()
		return nil, "", err
	}

	newData := &model.SeaLabsPayTopUpHolder{
		UserID: userID,
		TxnID:  txnId,
		Total:  amount,
	}
	holder, err := u.seaLabsPayTopUpHolderRepo.CreateTopUpHolder(tx, newData)
	if err != nil {
		tx.Rollback()
		return nil, "", err
	}

	tx.Commit()
	return holder, redirectURL, nil
}

func (u *userSeaPayAccountServ) TopUpWithSeaLabsPayCallback(txnID uint, status string) (*model.WalletTransaction, error) {
	tx := u.db.Begin()
	topUpHolder, topUpHolderError := u.seaLabsPayTopUpHolderRepo.UpdateTopUpHolder(tx, txnID, status)
	if topUpHolderError != nil {
		tx.Rollback()
		return nil, topUpHolderError
	}

	var transaction *model.WalletTransaction
	if status == dto.TXN_PAID {
		wallet, err := u.walletRepository.GetWalletByUserID(tx, topUpHolder.UserID)
		if err != nil {
			tx.Rollback()
			return nil, err
		}

		_, err = u.walletRepository.TopUp(tx, wallet, topUpHolder.Total)
		if err != nil {
			tx.Rollback()
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
			tx.Rollback()
			return nil, err
		}
	}

	tx.Commit()
	return transaction, nil
}
