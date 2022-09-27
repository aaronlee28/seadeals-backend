package service

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"gorm.io/gorm"
	"io"
	"net/http"
	"net/url"
	"seadeals-backend/apperror"
	"seadeals-backend/config"
	"seadeals-backend/dto"
	"seadeals-backend/model"
	"seadeals-backend/repository"
	"strconv"
	"strings"
	"time"
)

type UserSeaPayAccountServ interface {
	RegisterSeaLabsPayAccount(req *dto.RegisterSeaLabsPayReq) (*dto.RegisterSeaLabsPayRes, error)
	CheckSeaLabsAccountExists(req *dto.CheckSeaLabsPayReq) (*dto.CheckSeaLabsPayRes, error)
	UpdateSeaLabsAccountToMain(req *dto.UpdateSeaLabsPayToMainReq) (*model.UserSealabsPayAccount, error)
	GetSeaLabsAccountByUserID(userID uint) ([]*model.UserSealabsPayAccount, error)

	PayWithSeaLabsPay(amount int, accountNumber string) (string, error)
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

func generateHMACSHA256(value string, key string) string {
	h := hmac.New(sha256.New, []byte(key))
	h.Write([]byte(value))
	buf := h.Sum(nil)
	sign := hex.EncodeToString(buf)
	return sign
}

func transactionToSeaLabsPay(accountNumber string, amount string, sign string, callback string) (string, uint, error) {
	client := &http.Client{
		Timeout: time.Second * 10,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	data := url.Values{}
	data.Set("card_number", accountNumber)
	data.Set("amount", amount)
	data.Set("merchant_code", config.Config.SeaLabsPayMerchantCode)
	data.Set("redirect_url", "https://www.google.com")
	data.Set("callback_url", config.Config.NgrokURL+callback)
	data.Set("signature", sign)
	encodeData := data.Encode()

	fmt.Println(sign)
	req, err := http.NewRequest(http.MethodPost, config.Config.SeaLabsPayTransactionURL, strings.NewReader(encodeData))
	if err != nil {
		return "", 0, err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))
	response, err := client.Do(req)
	if err != nil {
		fmt.Println("dua")
		return "", 0, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println("Error Closing Client")
		}
	}(response.Body)

	if response.StatusCode == http.StatusSeeOther {
		redirectUrl, err := response.Location()
		if err != nil {
			return "", 0, err
		}

		TxnID, error := strconv.ParseUint(redirectUrl.Query().Get("txn_id"), 10, 64)
		if error != nil {
			return "", 0, error
		}
		return redirectUrl.String(), uint(TxnID), nil
	}
	return "", 0, apperror.BadRequestError("Cannot send data to sea labs pay API")
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

func (u *userSeaPayAccountServ) PayWithSeaLabsPay(amount int, accountNumber string) (string, error) {
	merchantCode := config.Config.SeaLabsPayMerchantCode
	apiKey := config.Config.SeaLabsPayAPIKey
	combinedString := accountNumber + ":" + strconv.Itoa(amount) + ":" + merchantCode

	sign := generateHMACSHA256(combinedString, apiKey)
	redirectURL, _, err := transactionToSeaLabsPay(accountNumber, strconv.Itoa(amount), sign, "")
	if err != nil {
		return "", err
	}

	return redirectURL, nil
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

	sign := generateHMACSHA256(combinedString, apiKey)
	redirectURL, txnId, err := transactionToSeaLabsPay(accountNumber, amountString, sign, "/user/wallet/top-up/sea-labs-pay/callback")
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
			PaymentMethod: "",
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
