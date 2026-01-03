package dto

import (
	"time"

	"personal-expense-splitting-settlement/internal/models"

	"github.com/google/uuid"
)

// CreateSettlementRequest represents the request to create a settlement
type CreateSettlementRequest struct {
	PayeeID       uuid.UUID            `json:"payee_id" validate:"required,uuid"`
	Amount        float64              `json:"amount" validate:"required,gt=0"`
	PaymentMethod models.PaymentMethod `json:"payment_method" validate:"required,oneof=cash bank_transfer upi paypal venmo credit_card debit_card other"`
	Notes         string               `json:"notes" validate:"omitempty,max=500"`
	GroupID       *uuid.UUID           `json:"group_id" validate:"omitempty,uuid"`
}

// UpdateSettlementRequest represents the request to update a settlement
type UpdateSettlementRequest struct {
	PaymentMethod *models.PaymentMethod `json:"payment_method" validate:"omitempty,oneof=cash bank_transfer upi paypal venmo credit_card debit_card other"`
	Notes         *string               `json:"notes" validate:"omitempty,max=500"`
}

// SettlementResponse represents a settlement in responses
type SettlementResponse struct {
	ID            uuid.UUID            `json:"id"`
	PayerID       uuid.UUID            `json:"payer_id"`
	PayerName     string               `json:"payer_name"`
	PayeeID       uuid.UUID            `json:"payee_id"`
	PayeeName     string               `json:"payee_name"`
	Amount        float64              `json:"amount"`
	PaymentMethod models.PaymentMethod `json:"payment_method"`
	Notes         string               `json:"notes"`
	IsConfirmed   bool                 `json:"is_confirmed"`
	ConfirmedAt   *time.Time           `json:"confirmed_at,omitempty"`
	GroupID       *uuid.UUID           `json:"group_id,omitempty"`
	GroupName     string               `json:"group_name,omitempty"`
	CreatedAt     time.Time            `json:"created_at"`
	UpdatedAt     time.Time            `json:"updated_at"`
}

// BalanceItem represents a balance between the user and another person
type BalanceItem struct {
	UserID    uuid.UUID `json:"user_id"`
	UserName  string    `json:"user_name"`
	UserEmail string    `json:"user_email"`
	Amount    float64   `json:"amount"` // Positive = they owe you, Negative = you owe them
}

// BalanceSummaryResponse represents the overall balance summary for a user
type BalanceSummaryResponse struct {
	TotalOwed   float64       `json:"total_owed"`  // Total amount others owe you
	TotalOwing  float64       `json:"total_owing"` // Total amount you owe others
	NetBalance  float64       `json:"net_balance"` // Total owed - Total owing
	Balances    []BalanceItem `json:"balances"`    // Individual balances
	LastUpdated time.Time     `json:"last_updated"`
}

// GroupBalanceItem represents a balance for a member in a group
type GroupBalanceItem struct {
	UserID     uuid.UUID `json:"user_id"`
	UserName   string    `json:"user_name"`
	TotalPaid  float64   `json:"total_paid"`  // Total amount paid by user
	TotalOwed  float64   `json:"total_owed"`  // Total amount owed by user
	NetBalance float64   `json:"net_balance"` // Paid - Owed
}

// GroupBalancesResponse represents balances within a group
type GroupBalancesResponse struct {
	GroupID      uuid.UUID          `json:"group_id"`
	GroupName    string             `json:"group_name"`
	TotalExpense float64            `json:"total_expense"`
	Balances     []GroupBalanceItem `json:"balances"`
	LastUpdated  time.Time          `json:"last_updated"`
}

// SettlementSuggestion represents a suggested settlement to simplify debts
type SettlementSuggestion struct {
	From     uuid.UUID `json:"from_user_id"`
	FromName string    `json:"from_user_name"`
	To       uuid.UUID `json:"to_user_id"`
	ToName   string    `json:"to_user_name"`
	Amount   float64   `json:"amount"`
}

// SettlementSuggestionsResponse represents settlement suggestions
type SettlementSuggestionsResponse struct {
	Suggestions []SettlementSuggestion `json:"suggestions"`
	TotalAmount float64                `json:"total_amount"`
	Message     string                 `json:"message"`
}

// ToSettlementResponse converts Settlement model to SettlementResponse
func ToSettlementResponse(s *models.Settlement) SettlementResponse {
	resp := SettlementResponse{
		ID:            s.ID,
		PayerID:       s.PayerID,
		PayeeID:       s.PayeeID,
		Amount:        s.Amount,
		PaymentMethod: s.PaymentMethod,
		Notes:         s.Notes,
		IsConfirmed:   s.IsConfirmed,
		ConfirmedAt:   s.ConfirmedAt,
		GroupID:       s.GroupID,
		CreatedAt:     s.CreatedAt,
		UpdatedAt:     s.UpdatedAt,
	}

	if s.Payer.ID != uuid.Nil {
		resp.PayerName = s.Payer.FirstName + " " + s.Payer.LastName
	}
	if s.Payee.ID != uuid.Nil {
		resp.PayeeName = s.Payee.FirstName + " " + s.Payee.LastName
	}
	if s.Group != nil {
		resp.GroupName = s.Group.Name
	}

	return resp
}
