package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// CreateCategoryRequest represents category creation request
type CreateCategoryRequest struct {
	Name string `json:"name" validate:"required,min=2,max=100"`
	Type string `json:"type" validate:"required,oneof=income expense"`
	Icon string `json:"icon" validate:"omitempty,max=50"`
}

// UpdateCategoryRequest represents category update request
type UpdateCategoryRequest struct {
	Name string `json:"name" validate:"omitempty,min=2,max=100"`
	Icon string `json:"icon" validate:"omitempty,max=50"`
}

// CategoryResponse represents category data in response
type CategoryResponse struct {
	ID        uuid.UUID  `json:"id"`
	UserID    *uuid.UUID `json:"user_id,omitempty"`
	Name      string     `json:"name"`
	Type      string     `json:"type"`
	Icon      string     `json:"icon"`
	IsDefault bool       `json:"is_default"`
	CreatedAt time.Time  `json:"created_at"`
}

// CreateBudgetRequest represents budget creation request
type CreateBudgetRequest struct {
	CategoryID *uuid.UUID      `json:"category_id,omitempty"`
	Amount     decimal.Decimal `json:"amount" validate:"required,gt=0"`
	Period     string          `json:"period" validate:"required,oneof=weekly monthly"`
	StartDate  time.Time       `json:"start_date" validate:"required"`
}

// UpdateBudgetRequest represents budget update request
type UpdateBudgetRequest struct {
	Amount   decimal.Decimal `json:"amount" validate:"omitempty,gt=0"`
	IsActive *bool           `json:"is_active,omitempty"`
}

// BudgetResponse represents budget data in response
type BudgetResponse struct {
	ID              uuid.UUID       `json:"id"`
	UserID          uuid.UUID       `json:"user_id"`
	CategoryID      *uuid.UUID      `json:"category_id,omitempty"`
	CategoryName    string          `json:"category_name,omitempty"`
	Amount          decimal.Decimal `json:"amount"`
	SpentAmount     decimal.Decimal `json:"spent_amount"`
	RemainingAmount decimal.Decimal `json:"remaining_amount"`
	PercentageUsed  float64         `json:"percentage_used"`
	Period          string          `json:"period"`
	StartDate       time.Time       `json:"start_date"`
	EndDate         time.Time       `json:"end_date"`
	IsActive        bool            `json:"is_active"`
	IsExceeded      bool            `json:"is_exceeded"`
	IsNearLimit     bool            `json:"is_near_limit"`
	CreatedAt       time.Time       `json:"created_at"`
}

// BudgetAlert represents a budget alert
type BudgetAlert struct {
	BudgetID       uuid.UUID       `json:"budget_id"`
	CategoryName   string          `json:"category_name,omitempty"`
	Amount         decimal.Decimal `json:"amount"`
	SpentAmount    decimal.Decimal `json:"spent_amount"`
	PercentageUsed float64         `json:"percentage_used"`
	Status         string          `json:"status"` // "exceeded", "near_limit"
	Message        string          `json:"message"`
}
