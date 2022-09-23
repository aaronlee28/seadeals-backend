package service

import (
	"github.com/mailjet/mailjet-apiv3-go"
	"gorm.io/gorm"
	"math"
	"seadeals-backend/apperror"
	"seadeals-backend/config"
	"seadeals-backend/dto"
	"seadeals-backend/helper"
	"seadeals-backend/repository"
	"strconv"
)

type WalletService interface {
	UserWalletData(id uint) (*dto.WalletDataRes, error)
	TransactionDetails(id uint) (*dto.TransactionDetailsRes, error)
	PaginatedTransactions(q *repository.Query, userID uint) (*dto.PaginatedTransactionsRes, error)
	WalletPin(userID uint, pin string) error
	RequestPinChangeWithEmail(userID uint) (*mailjet.ResultsV31, error)
	ValidateWalletPin(userID uint, pin string) (bool, error)
	GetWalletStatus(userID uint) (string, error)
}

type walletService struct {
	db               *gorm.DB
	walletRepository repository.WalletRepository
	userRepository   repository.UserRepository
}

type WalletServiceConfig struct {
	DB               *gorm.DB
	WalletRepository repository.WalletRepository
	UserRepository   repository.UserRepository
}

func NewWalletService(c *WalletServiceConfig) WalletService {
	return &walletService{
		db:               c.DB,
		walletRepository: c.WalletRepository,
		userRepository:   c.UserRepository,
	}
}

func (w *walletService) UserWalletData(id uint) (*dto.WalletDataRes, error) {
	tx := w.db.Begin()
	wallet, err := w.walletRepository.GetWalletByUserID(tx, id)

	if err != nil {
		tx.Rollback()
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

func (w *walletService) TransactionDetails(id uint) (*dto.TransactionDetailsRes, error) {
	tx := w.db.Begin()
	t, err := w.walletRepository.TransactionDetails(tx, id)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	transaction := &dto.TransactionDetailsRes{
		Id:            t.Id,
		VoucherID:     t.VoucherID,
		Total:         t.Total,
		PaymentType:   t.PaymentType,
		PaymentMethod: t.PaymentMethod,
		CreatedAt:     t.CreatedAt,
		UpdatedAt:     t.UpdatedAt,
	}

	tx.Commit()
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
	var ts []dto.TransactionsRes
	l, t, err := w.walletRepository.PaginatedTransactions(tx, q, userID)
	if err != nil {
		tx.Rollback()
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

	tx.Commit()
	return &paginatedTransactions, nil
}

func (w *walletService) WalletPin(userID uint, pin string) error {
	tx := w.db.Begin()
	if len(pin) != 6 {
		return apperror.BadRequestError("Pin has to be 6 digits long")
	}
	err := w.walletRepository.WalletPin(tx, userID, pin)
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

func (w *walletService) RequestPinChangeWithEmail(userID uint) (*mailjet.ResultsV31, error) {
	tx := w.db.Begin()
	user, err := w.userRepository.GetUserByID(tx, userID)
	if err != nil {
		return nil, err
	}

	randomString := helper.RandomString(12)
	code := helper.RandomString(6)
	w.walletRepository.RequestChangePinByEmail(user.ID, randomString, code)

	mailjetClient := mailjet.NewMailjetClient(config.Config.MailJetPublicKey, config.Config.MailJetSecretKey)
	messagesInfo := []mailjet.InfoMessagesV31{
		{
			From: &mailjet.RecipientV31{
				Email: "SeaDeals-noreply@seadeals.com",
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
			HTMLPart: "<p>Here are the code to reset your pin:</p><h3>" + code + "</h3>" +
				"<p>you can open the link <a href=\"localhost:3000\\check?userID=" + strconv.FormatUint(uint64(user.ID), 10) + "key=" + randomString + "\">here</a></p>",
			Priority: 0,
			CustomID: config.Config.AppName,
		},
	}
	messages := mailjet.MessagesV31{
		Info: messagesInfo,
	}

	res, err := mailjetClient.SendMailV31(&messages)
	if err != nil {
		return nil, err
	}
	tx.Commit()
	return res, nil
}

func (w *walletService) ValidateWalletPin(userID uint, pin string) (bool, error) {
	tx := w.db.Begin()
	if len(pin) != 6 {
		return false, apperror.BadRequestError("Pin has to be 6 digits long")
	}

	err := w.walletRepository.ValidateWalletPin(tx, userID, pin)
	if err != nil {
		tx.Rollback()
		return false, err
	}

	tx.Commit()
	return true, nil
}

func (w *walletService) GetWalletStatus(userID uint) (string, error) {
	tx := w.db.Begin()

	status, err := w.walletRepository.GetWalletStatus(tx, userID)
	if err != nil {
		tx.Rollback()
		return "", err
	}

	tx.Commit()
	return status, nil
}
