package dto

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// WalletResponse represents wallet data in response
type WalletResponse struct {
	ID       uuid.UUID       `json:"id"`
	UserID   uuid.UUID       `json:"user_id"`
	Balance  decimal.Decimal `json:"balance"`
	Currency string          `json:"currency"`
	IsActive bool            `json:"is_active"`
}

// BalanceResponse represents balance information
type BalanceResponse struct {
	Balance  decimal.Decimal `json:"balance"`
	Currency string          `json:"currency"`
}
