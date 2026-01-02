package repository

import (
	"digital-wallet-api/internal/models"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type TransactionRepository interface {
	Create(transaction *models.Transaction) error
	FindByID(id uuid.UUID) (*models.Transaction, error)
	FindByWalletID(walletID uuid.UUID, limit, offset int) ([]models.Transaction, int64, error)
	FindByReferenceID(referenceID string) (*models.Transaction, error)
	GetSummary(walletID uuid.UUID, startDate, endDate time.Time) (totalCredit, totalDebit decimal.Decimal, count int64, err error)
}

type transactionRepository struct {
	db *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) TransactionRepository {
	return &transactionRepository{db: db}
}

func (r *transactionRepository) Create(transaction *models.Transaction) error {
	return r.db.Create(transaction).Error
}

func (r *transactionRepository) FindByID(id uuid.UUID) (*models.Transaction, error) {
	var transaction models.Transaction
	err := r.db.Preload("Category").Where("id = ?", id).First(&transaction).Error
	if err != nil {
		return nil, err
	}
	return &transaction, nil
}

func (r *transactionRepository) FindByWalletID(walletID uuid.UUID, limit, offset int) ([]models.Transaction, int64, error) {
	var transactions []models.Transaction
	var total int64

	query := r.db.Model(&models.Transaction{}).Where("wallet_id = ?", walletID)

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	err := query.
		Preload("Category").
		Order("transaction_date DESC").
		Limit(limit).
		Offset(offset).
		Find(&transactions).Error

	return transactions, total, err
}

func (r *transactionRepository) FindByReferenceID(referenceID string) (*models.Transaction, error) {
	var transaction models.Transaction
	err := r.db.Where("reference_id = ?", referenceID).First(&transaction).Error
	if err != nil {
		return nil, err
	}
	return &transaction, nil
}

func (r *transactionRepository) GetSummary(walletID uuid.UUID, startDate, endDate time.Time) (totalCredit, totalDebit decimal.Decimal, count int64, err error) {
	type Result struct {
		Type  string
		Total decimal.Decimal
		Count int64
	}

	var results []Result
	err = r.db.Model(&models.Transaction{}).
		Select("type, SUM(amount) as total, COUNT(*) as count").
		Where("wallet_id = ? AND transaction_date BETWEEN ? AND ?", walletID, startDate, endDate).
		Group("type").
		Scan(&results).Error

	if err != nil {
		return decimal.Zero, decimal.Zero, 0, err
	}

	totalCredit = decimal.Zero
	totalDebit = decimal.Zero
	count = 0

	for _, result := range results {
		if result.Type == string(models.TransactionTypeCredit) {
			totalCredit = result.Total
		} else if result.Type == string(models.TransactionTypeDebit) {
			totalDebit = result.Total
		}
		count += result.Count
	}

	return totalCredit, totalDebit, count, nil
}
