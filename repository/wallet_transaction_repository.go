package repository

import (
	"gorm.io/gorm"
	"seadeals-backend/apperror"
	"seadeals-backend/dto"
	"seadeals-backend/model"
	"strconv"
)

type WalletTransactionRepository interface {
	CreateTransaction(tx *gorm.DB, model *model.WalletTransaction) (*model.WalletTransaction, error)
	GetTransactionsByWalletID(tx *gorm.DB, query *dto.WalletTransactionsQuery, walletID uint) ([]*model.WalletTransaction, int64, int64, error)
}

type walletTransactionRepository struct{}

func NewWalletTransactionRepository() WalletTransactionRepository {
	return &walletTransactionRepository{}
}

func (w *walletTransactionRepository) CreateTransaction(tx *gorm.DB, model *model.WalletTransaction) (*model.WalletTransaction, error) {
	result := tx.Create(&model)
	if result.Error != nil {
		return nil, apperror.InternalServerError("Cannot create wallet transaction")
	}
	return model, nil
}

func (w *walletTransactionRepository) GetTransactionsByWalletID(tx *gorm.DB, query *dto.WalletTransactionsQuery, walletID uint) ([]*model.WalletTransaction, int64, int64, error) {
	var transactions = make([]*model.WalletTransaction, 0)
	result := tx.Model(&transactions).Where("wallet_id = ?", walletID)

	orderByString := query.SortBy
	if query.SortBy == "total" {
		orderByString = "total"
	} else {
		orderByString = "created_at"
	}

	if query.SortBy == "" {
		if query.Sort != "asc" {
			orderByString += " desc"
		}
	} else if query.Sort == "desc" {
		orderByString += " desc"
	}

	var totalData int64
	result = result.Order(orderByString).Order("id")
	table := tx.Table("(?) as s1", result).Count(&totalData)

	limit, _ := strconv.Atoi(query.Limit)
	if limit == 0 {
		limit = 20
	}
	table = table.Limit(limit)

	page, _ := strconv.Atoi(query.Page)
	if page == 0 {
		page = 1
	}
	table = table.Offset((page - 1) * limit)

	table = table.Unscoped().Find(&transactions)
	if table.Error != nil {
		return nil, 0, 0, table.Error
	}

	totalPage := totalData / int64(limit)
	if totalData%int64(limit) != 0 {
		totalPage += 1
	}

	return transactions, totalPage, totalData, nil
}
