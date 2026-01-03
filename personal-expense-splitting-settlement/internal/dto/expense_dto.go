package dto

import (
	"encoding/json"
	"fmt"
	"time"

	"personal-expense-splitting-settlement/internal/models"

	"github.com/google/uuid"
)

// DateOnly is a custom type for handling date-only (YYYY-MM-DD) JSON marshaling/unmarshaling
type DateOnly struct {
	Time time.Time
}

const dateFormat = "2006-01-02"

// UnmarshalJSON implements custom JSON unmarshaling for DateOnly
func (d *DateOnly) UnmarshalJSON(b []byte) error {
	s := string(b)
	// Remove quotes
	if len(s) > 2 && s[0] == '"' && s[len(s)-1] == '"' {
		s = s[1 : len(s)-1]
	}

	if s == "null" || s == "" {
		d.Time = time.Time{}
		return nil
	}

	t, err := time.Parse(dateFormat, s)
	if err != nil {
		return fmt.Errorf("invalid date format, expected YYYY-MM-DD: %w", err)
	}
	d.Time = t
	return nil
}

// MarshalJSON implements custom JSON marshaling for DateOnly
func (d DateOnly) MarshalJSON() ([]byte, error) {
	if d.Time.IsZero() {
		return []byte("null"), nil
	}
	return json.Marshal(d.Time.Format(dateFormat))
}

// UnmarshalText implements encoding.TextUnmarshaler for query parameter binding
func (d *DateOnly) UnmarshalText(data []byte) error {
	s := string(data)

	if s == "" {
		d.Time = time.Time{}
		return nil
	}

	t, err := time.Parse(dateFormat, s)
	if err != nil {
		return fmt.Errorf("invalid date format, expected YYYY-MM-DD: %w", err)
	}
	d.Time = t
	return nil
}

// MarshalText implements encoding.TextMarshaler
func (d DateOnly) MarshalText() ([]byte, error) {
	if d.Time.IsZero() {
		return []byte(""), nil
	}
	return []byte(d.Time.Format(dateFormat)), nil
}

// SplitMethod represents how the expense should be split
type SplitMethod string

const (
	SplitMethodEqual      SplitMethod = "equal"      // Split equally among all participants
	SplitMethodExact      SplitMethod = "exact"      // Exact amounts specified for each participant
	SplitMethodPercentage SplitMethod = "percentage" // Percentage-based split
	SplitMethodShares     SplitMethod = "shares"     // Share-based split
)

// ParticipantSplit represents how a participant should be split
type ParticipantSplit struct {
	UserID     uuid.UUID `json:"user_id" validate:"required,uuid"`
	PaidAmount float64   `json:"paid_amount" validate:"gte=0"`
	OwedAmount *float64  `json:"owed_amount,omitempty" validate:"omitempty,gt=0"`        // For exact split
	Percentage *float64  `json:"percentage,omitempty" validate:"omitempty,gt=0,lte=100"` // For percentage split
	Shares     *int      `json:"shares,omitempty" validate:"omitempty,gt=0"`             // For share-based split
}

// CreateExpenseRequest represents the request to create a new expense
type CreateExpenseRequest struct {
	Description  string                 `json:"description" validate:"required,min=1,max=500"`
	Amount       float64                `json:"amount" validate:"required,gt=0"`
	Category     models.ExpenseCategory `json:"category" validate:"required,oneof=general food transport entertainment utilities shopping healthcare education travel other"`
	ReceiptURL   string                 `json:"receipt_url" validate:"omitempty,url"`
	Date         *DateOnly              `json:"date" validate:"omitempty"`
	GroupID      *uuid.UUID             `json:"group_id" validate:"omitempty,uuid"`
	SplitMethod  SplitMethod            `json:"split_method" validate:"required,oneof=equal exact percentage shares"`
	Participants []ParticipantSplit     `json:"participants" validate:"required,min=1,dive"`
}

// UpdateExpenseRequest represents the request to update an expense
type UpdateExpenseRequest struct {
	Description *string                 `json:"description" validate:"omitempty,min=1,max=500"`
	Amount      *float64                `json:"amount" validate:"omitempty,gt=0"`
	Category    *models.ExpenseCategory `json:"category" validate:"omitempty,oneof=general food transport entertainment utilities shopping healthcare education travel other"`
	ReceiptURL  *string                 `json:"receipt_url" validate:"omitempty,url"`
	Date        *DateOnly               `json:"date" validate:"omitempty"`
}

// ExpenseParticipantResponse represents a participant in an expense
type ExpenseParticipantResponse struct {
	ID         uuid.UUID  `json:"id"`
	UserID     uuid.UUID  `json:"user_id"`
	Name       string     `json:"name"`
	Email      string     `json:"email"`
	PaidAmount float64    `json:"paid_amount"`
	OwedAmount float64    `json:"owed_amount"`
	NetAmount  float64    `json:"net_amount"` // paid - owed
	IsSettled  bool       `json:"is_settled"`
	SettledAt  *time.Time `json:"settled_at,omitempty"`
}

// ExpenseResponse represents the basic expense response
type ExpenseResponse struct {
	ID               uuid.UUID              `json:"id"`
	Description      string                 `json:"description"`
	Amount           float64                `json:"amount"`
	Category         models.ExpenseCategory `json:"category"`
	ReceiptURL       string                 `json:"receipt_url"`
	Date             time.Time              `json:"date"`
	CreatedBy        uuid.UUID              `json:"created_by"`
	GroupID          *uuid.UUID             `json:"group_id"`
	ParticipantCount int                    `json:"participant_count"`
	IsSettled        bool                   `json:"is_settled"`
	CreatedAt        time.Time              `json:"created_at"`
	UpdatedAt        time.Time              `json:"updated_at"`
}

// ExpenseDetailResponse represents the detailed expense response
type ExpenseDetailResponse struct {
	ID           uuid.UUID                    `json:"id"`
	Description  string                       `json:"description"`
	Amount       float64                      `json:"amount"`
	Category     models.ExpenseCategory       `json:"category"`
	ReceiptURL   string                       `json:"receipt_url"`
	Date         time.Time                    `json:"date"`
	CreatedBy    uuid.UUID                    `json:"created_by"`
	GroupID      *uuid.UUID                   `json:"group_id"`
	GroupName    string                       `json:"group_name,omitempty"`
	Participants []ExpenseParticipantResponse `json:"participants"`
	IsSettled    bool                         `json:"is_settled"`
	CreatedAt    time.Time                    `json:"created_at"`
	UpdatedAt    time.Time                    `json:"updated_at"`
}

// ExpenseListResponse represents the response for a list of expenses
type ExpenseListResponse struct {
	Expenses []ExpenseResponse `json:"expenses"`
	Total    int               `json:"total"`
}

// ExpenseListRequest represents filters for listing expenses
type ExpenseListRequest struct {
	GroupID   string `form:"group_id"`
	Category  string `form:"category"`
	StartDate string `form:"start_date"`
	EndDate   string `form:"end_date"`
	Limit     int    `form:"limit" validate:"omitempty,min=1,max=100"`
	Offset    int    `form:"offset" validate:"omitempty,min=0"`
}

// ToExpenseResponse converts an Expense model to ExpenseResponse
func ToExpenseResponse(expense *models.Expense) ExpenseResponse {
	return ExpenseResponse{
		ID:               expense.ID,
		Description:      expense.Description,
		Amount:           expense.Amount,
		Category:         expense.Category,
		ReceiptURL:       expense.ReceiptURL,
		Date:             expense.Date,
		CreatedBy:        expense.CreatedBy,
		GroupID:          expense.GroupID,
		ParticipantCount: expense.GetParticipantCount(),
		IsSettled:        expense.IsFullySettled(),
		CreatedAt:        expense.CreatedAt,
		UpdatedAt:        expense.UpdatedAt,
	}
}

// ToExpenseDetailResponse converts an Expense model to ExpenseDetailResponse
func ToExpenseDetailResponse(expense *models.Expense) ExpenseDetailResponse {
	participants := make([]ExpenseParticipantResponse, len(expense.Participants))
	for i, participant := range expense.Participants {
		participants[i] = ExpenseParticipantResponse{
			ID:         participant.ID,
			UserID:     participant.UserID,
			Name:       participant.User.FirstName + " " + participant.User.LastName,
			Email:      participant.User.Email,
			PaidAmount: participant.PaidAmount,
			OwedAmount: participant.OwedAmount,
			NetAmount:  participant.GetNetAmount(),
			IsSettled:  participant.IsSettled,
			SettledAt:  participant.SettledAt,
		}
	}

	groupName := ""
	if expense.Group != nil {
		groupName = expense.Group.Name
	}

	return ExpenseDetailResponse{
		ID:           expense.ID,
		Description:  expense.Description,
		Amount:       expense.Amount,
		Category:     expense.Category,
		ReceiptURL:   expense.ReceiptURL,
		Date:         expense.Date,
		CreatedBy:    expense.CreatedBy,
		GroupID:      expense.GroupID,
		GroupName:    groupName,
		Participants: participants,
		IsSettled:    expense.IsFullySettled(),
		CreatedAt:    expense.CreatedAt,
		UpdatedAt:    expense.UpdatedAt,
	}
}

// ToExpenseParticipantResponse converts an ExpenseParticipant model to ExpenseParticipantResponse
func ToExpenseParticipantResponse(participant *models.ExpenseParticipant) ExpenseParticipantResponse {
	return ExpenseParticipantResponse{
		ID:         participant.ID,
		UserID:     participant.UserID,
		Name:       participant.User.FirstName + " " + participant.User.LastName,
		Email:      participant.User.Email,
		PaidAmount: participant.PaidAmount,
		OwedAmount: participant.OwedAmount,
		NetAmount:  participant.GetNetAmount(),
		IsSettled:  participant.IsSettled,
		SettledAt:  participant.SettledAt,
	}
}
