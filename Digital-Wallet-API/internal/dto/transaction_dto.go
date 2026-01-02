package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// CreditRequest represents a credit transaction request
type CreditRequest struct {
	Amount      decimal.Decimal `json:"amount" validate:"required,gt=0"`
	CategoryID  *uuid.UUID      `json:"category_id,omitempty"`
	Description string          `json:"description" validate:"omitempty,max=500"`
}

// DebitRequest represents a debit transaction request
type DebitRequest struct {
	Amount      decimal.Decimal `json:"amount" validate:"required,gt=0"`
	CategoryID  *uuid.UUID      `json:"category_id" validate:"required"`
	Description string          `json:"description" validate:"omitempty,max=500"`
}

// TransferRequest represents a transfer request
type TransferRequest struct {
	ToUserID    uuid.UUID       `json:"to_user_id" validate:"required"`
	Amount      decimal.Decimal `json:"amount" validate:"required,gt=0"`
	Description string          `json:"description" validate:"omitempty,max=500"`
}

// TransactionResponse represents transaction data in response
type TransactionResponse struct {
	ID              uuid.UUID       `json:"id"`
	WalletID        uuid.UUID       `json:"wallet_id"`
	CategoryID      *uuid.UUID      `json:"category_id,omitempty"`
	CategoryName    string          `json:"category_name,omitempty"`
	Type            string          `json:"type"`
	Amount          decimal.Decimal `json:"amount"`
	BalanceAfter    decimal.Decimal `json:"balance_after"`
	Description     string          `json:"description"`
	ReferenceID     string          `json:"reference_id"`
	Status          string          `json:"status"`
	TransactionDate time.Time       `json:"transaction_date"`
	CreatedAt       time.Time       `json:"created_at"`
}

// TransactionListResponse represents paginated transaction list
type TransactionListResponse struct {
	Transactions []TransactionResponse `json:"transactions"`
	Total        int64                 `json:"total"`
	Page         int                   `json:"page"`
	PageSize     int                   `json:"page_size"`
}

// TransactionSummary represents transaction summary
type TransactionSummary struct {
	TotalCredit  decimal.Decimal `json:"total_credit"`
	TotalDebit   decimal.Decimal `json:"total_debit"`
	NetAmount    decimal.Decimal `json:"net_amount"`
	Transactions int64           `json:"transaction_count"`
}

// TransferResponse represents transfer data in response
type TransferResponse struct {
	ID           uuid.UUID       `json:"id"`
	FromWalletID uuid.UUID       `json:"from_wallet_id"`
	ToWalletID   uuid.UUID       `json:"to_wallet_id"`
	Amount       decimal.Decimal `json:"amount"`
	Description  string          `json:"description"`
	ReferenceID  string          `json:"reference_id"`
	Status       string          `json:"status"`
	CreatedAt    time.Time       `json:"created_at"`
}
