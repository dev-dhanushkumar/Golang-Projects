package service

import (
	"digital-wallet-api/internal/dto"
	"digital-wallet-api/internal/models"
	"digital-wallet-api/internal/repository"
	"digital-wallet-api/pkg/logger"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TransactionService interface {
	Credit(userID uuid.UUID, req dto.CreditRequest) (*dto.TransactionResponse, error)
	Debit(userID uuid.UUID, req dto.DebitRequest) (*dto.TransactionResponse, error)
	Transfer(userID uuid.UUID, req dto.TransferRequest) (*dto.TransferResponse, error)
	GetTransactions(userID uuid.UUID, page, pageSize int) (*dto.TransactionListResponse, error)
	GetTransaction(userID uuid.UUID, transactionID uuid.UUID) (*dto.TransactionResponse, error)
	GetSummary(userID uuid.UUID, startDate, endDate time.Time) (*dto.TransactionSummary, error)
}

type transactionService struct {
	walletRepo      repository.WalletRepository
	transactionRepo repository.TransactionRepository
	db              *gorm.DB
}

func NewTransactionService(
	walletRepo repository.WalletRepository,
	transactionRepo repository.TransactionRepository,
	db *gorm.DB,
) TransactionService {
	return &transactionService{
		walletRepo:      walletRepo,
		transactionRepo: transactionRepo,
		db:              db,
	}
}

func (s *transactionService) Credit(userID uuid.UUID, req dto.CreditRequest) (*dto.TransactionResponse, error) {
	// Start transaction
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Get wallet
	wallet, err := s.walletRepo.FindByUserID(userID)
	if err != nil {
		tx.Rollback()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, models.ErrWalletNotFound
		}
		return nil, err
	}

	// Update balance
	wallet.Credit(req.Amount)
	if err := tx.Save(wallet).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// Create transaction
	transaction := &models.Transaction{
		WalletID:        wallet.ID,
		CategoryID:      req.CategoryID,
		Type:            models.TransactionTypeCredit,
		Amount:          req.Amount,
		BalanceAfter:    wallet.Balance,
		Description:     req.Description,
		Status:          models.TransactionStatusCompleted,
		TransactionDate: time.Now(),
	}

	if err := tx.Create(transaction).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	logger.Info("Credit transaction completed", map[string]interface{}{
		"user_id":        userID,
		"wallet_id":      wallet.ID,
		"amount":         req.Amount,
		"transaction_id": transaction.ID,
	})

	return s.toTransactionResponse(transaction), nil
}

func (s *transactionService) Debit(userID uuid.UUID, req dto.DebitRequest) (*dto.TransactionResponse, error) {
	// Start transaction
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Get wallet
	wallet, err := s.walletRepo.FindByUserID(userID)
	if err != nil {
		tx.Rollback()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, models.ErrWalletNotFound
		}
		return nil, err
	}

	// Check sufficient balance
	if !wallet.HasSufficientBalance(req.Amount) {
		tx.Rollback()
		return nil, models.ErrInsufficientBalance
	}

	// Update balance
	if err := wallet.Debit(req.Amount); err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Save(wallet).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// Create transaction
	transaction := &models.Transaction{
		WalletID:        wallet.ID,
		CategoryID:      req.CategoryID,
		Type:            models.TransactionTypeDebit,
		Amount:          req.Amount,
		BalanceAfter:    wallet.Balance,
		Description:     req.Description,
		Status:          models.TransactionStatusCompleted,
		TransactionDate: time.Now(),
	}

	if err := tx.Create(transaction).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	logger.Info("Debit transaction completed", map[string]interface{}{
		"user_id":        userID,
		"wallet_id":      wallet.ID,
		"amount":         req.Amount,
		"transaction_id": transaction.ID,
	})

	return s.toTransactionResponse(transaction), nil
}

func (s *transactionService) Transfer(userID uuid.UUID, req dto.TransferRequest) (*dto.TransferResponse, error) {
	// Start transaction
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Get sender wallet
	fromWallet, err := s.walletRepo.FindByUserID(userID)
	if err != nil {
		tx.Rollback()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, models.ErrWalletNotFound
		}
		return nil, err
	}

	// Get receiver wallet
	toWallet, err := s.walletRepo.FindByUserID(req.ToUserID)
	if err != nil {
		tx.Rollback()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("receiver wallet not found")
		}
		return nil, err
	}

	// Check if same wallet
	if fromWallet.ID == toWallet.ID {
		tx.Rollback()
		return nil, models.ErrSameWalletTransfer
	}

	// Check sufficient balance
	if !fromWallet.HasSufficientBalance(req.Amount) {
		tx.Rollback()
		return nil, models.ErrInsufficientBalance
	}

	// Debit from sender
	if err := fromWallet.Debit(req.Amount); err != nil {
		tx.Rollback()
		return nil, err
	}
	if err := tx.Save(fromWallet).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// Credit to receiver
	toWallet.Credit(req.Amount)
	if err := tx.Save(toWallet).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// Create transfer record
	transfer := &models.Transfer{
		FromWalletID: fromWallet.ID,
		ToWalletID:   toWallet.ID,
		Amount:       req.Amount,
		Description:  req.Description,
		Status:       models.TransactionStatusCompleted,
	}

	if err := tx.Create(transfer).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// Create transactions for both wallets
	// Debit transaction for sender
	debitTx := &models.Transaction{
		WalletID:        fromWallet.ID,
		Type:            models.TransactionTypeDebit,
		Amount:          req.Amount,
		BalanceAfter:    fromWallet.Balance,
		Description:     "Transfer to user",
		ReferenceID:     transfer.ReferenceID,
		Status:          models.TransactionStatusCompleted,
		TransactionDate: time.Now(),
	}
	if err := tx.Create(debitTx).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// Credit transaction for receiver
	creditTx := &models.Transaction{
		WalletID:        toWallet.ID,
		Type:            models.TransactionTypeCredit,
		Amount:          req.Amount,
		BalanceAfter:    toWallet.Balance,
		Description:     "Transfer from user",
		ReferenceID:     transfer.ReferenceID,
		Status:          models.TransactionStatusCompleted,
		TransactionDate: time.Now(),
	}
	if err := tx.Create(creditTx).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	logger.Info("Transfer completed", map[string]interface{}{
		"from_user_id": userID,
		"to_user_id":   req.ToUserID,
		"amount":       req.Amount,
		"transfer_id":  transfer.ID,
	})

	return &dto.TransferResponse{
		ID:           transfer.ID,
		FromWalletID: transfer.FromWalletID,
		ToWalletID:   transfer.ToWalletID,
		Amount:       transfer.Amount,
		Description:  transfer.Description,
		ReferenceID:  transfer.ReferenceID,
		Status:       string(transfer.Status),
		CreatedAt:    transfer.CreatedAt,
	}, nil
}

func (s *transactionService) GetTransactions(userID uuid.UUID, page, pageSize int) (*dto.TransactionListResponse, error) {
	// Get wallet
	wallet, err := s.walletRepo.FindByUserID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, models.ErrWalletNotFound
		}
		return nil, err
	}

	// Calculate offset
	offset := (page - 1) * pageSize

	// Get transactions
	transactions, total, err := s.transactionRepo.FindByWalletID(wallet.ID, pageSize, offset)
	if err != nil {
		return nil, err
	}

	// Convert to response
	var txResponses []dto.TransactionResponse
	for _, tx := range transactions {
		txResponses = append(txResponses, *s.toTransactionResponse(&tx))
	}

	return &dto.TransactionListResponse{
		Transactions: txResponses,
		Total:        total,
		Page:         page,
		PageSize:     pageSize,
	}, nil
}

func (s *transactionService) GetTransaction(userID uuid.UUID, transactionID uuid.UUID) (*dto.TransactionResponse, error) {
	// Get wallet
	wallet, err := s.walletRepo.FindByUserID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, models.ErrWalletNotFound
		}
		return nil, err
	}

	// Get transaction
	transaction, err := s.transactionRepo.FindByID(transactionID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("transaction not found")
		}
		return nil, err
	}

	// Verify ownership
	if transaction.WalletID != wallet.ID {
		return nil, models.ErrUnauthorized
	}

	return s.toTransactionResponse(transaction), nil
}

func (s *transactionService) GetSummary(userID uuid.UUID, startDate, endDate time.Time) (*dto.TransactionSummary, error) {
	// Get wallet
	wallet, err := s.walletRepo.FindByUserID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, models.ErrWalletNotFound
		}
		return nil, err
	}

	// Get summary
	totalCredit, totalDebit, count, err := s.transactionRepo.GetSummary(wallet.ID, startDate, endDate)
	if err != nil {
		return nil, err
	}

	netAmount := totalCredit.Sub(totalDebit)

	return &dto.TransactionSummary{
		TotalCredit:  totalCredit,
		TotalDebit:   totalDebit,
		NetAmount:    netAmount,
		Transactions: count,
	}, nil
}

func (s *transactionService) toTransactionResponse(tx *models.Transaction) *dto.TransactionResponse {
	response := &dto.TransactionResponse{
		ID:              tx.ID,
		WalletID:        tx.WalletID,
		CategoryID:      tx.CategoryID,
		Type:            string(tx.Type),
		Amount:          tx.Amount,
		BalanceAfter:    tx.BalanceAfter,
		Description:     tx.Description,
		ReferenceID:     tx.ReferenceID,
		Status:          string(tx.Status),
		TransactionDate: tx.TransactionDate,
		CreatedAt:       tx.CreatedAt,
	}

	if tx.Category != nil {
		response.CategoryName = tx.Category.Name
	}

	return response
}
