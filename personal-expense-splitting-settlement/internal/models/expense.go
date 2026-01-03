package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ExpenseCategory represents the category of an expense
type ExpenseCategory string

const (
	ExpenseCategoryGeneral       ExpenseCategory = "general"
	ExpenseCategoryFood          ExpenseCategory = "food"
	ExpenseCategoryTransport     ExpenseCategory = "transport"
	ExpenseCategoryEntertainment ExpenseCategory = "entertainment"
	ExpenseCategoryUtilities     ExpenseCategory = "utilities"
	ExpenseCategoryShopping      ExpenseCategory = "shopping"
	ExpenseCategoryHealthcare    ExpenseCategory = "healthcare"
	ExpenseCategoryEducation     ExpenseCategory = "education"
	ExpenseCategoryTravel        ExpenseCategory = "travel"
	ExpenseCategoryOther         ExpenseCategory = "other"
)

// Expense represents an expense record
type Expense struct {
	ID          uuid.UUID       `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Description string          `gorm:"type:varchar(500);not null" json:"description"`
	Amount      float64         `gorm:"type:decimal(12,2);not null" json:"amount"`
	Category    ExpenseCategory `gorm:"type:varchar(50);not null;default:'general'" json:"category"`
	ReceiptURL  string          `gorm:"type:text" json:"receipt_url"`
	Date        time.Time       `gorm:"type:date;not null;default:CURRENT_DATE" json:"date"`
	CreatedBy   uuid.UUID       `gorm:"type:uuid;not null;column:created_by" json:"created_by"`
	GroupID     *uuid.UUID      `gorm:"type:uuid;column:group_id" json:"group_id"`
	CreatedAt   time.Time       `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time       `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt   gorm.DeletedAt  `gorm:"index" json:"deleted_at,omitempty"`

	// Relationships
	Creator      User                 `gorm:"foreignKey:CreatedBy" json:"creator,omitempty"`
	Group        *Group               `gorm:"foreignKey:GroupID" json:"group,omitempty"`
	Participants []ExpenseParticipant `gorm:"foreignKey:ExpenseID" json:"participants,omitempty"`
}

// TableName specifies the table name for Expense model
func (Expense) TableName() string {
	return "expenses"
}

// IsValidCategory checks if the expense category is valid
func (e *Expense) IsValidCategory() bool {
	validCategories := []ExpenseCategory{
		ExpenseCategoryGeneral,
		ExpenseCategoryFood,
		ExpenseCategoryTransport,
		ExpenseCategoryEntertainment,
		ExpenseCategoryUtilities,
		ExpenseCategoryShopping,
		ExpenseCategoryHealthcare,
		ExpenseCategoryEducation,
		ExpenseCategoryTravel,
		ExpenseCategoryOther,
	}
	for _, validCategory := range validCategories {
		if e.Category == validCategory {
			return true
		}
	}
	return false
}

// GetTotalPaid returns the total amount paid by all participants
func (e *Expense) GetTotalPaid() float64 {
	total := 0.0
	for _, participant := range e.Participants {
		total += participant.PaidAmount
	}
	return total
}

// GetTotalOwed returns the total amount owed by all participants
func (e *Expense) GetTotalOwed() float64 {
	total := 0.0
	for _, participant := range e.Participants {
		total += participant.OwedAmount
	}
	return total
}

// GetParticipantCount returns the number of participants
func (e *Expense) GetParticipantCount() int {
	return len(e.Participants)
}

// IsGroupExpense checks if this is a group expense
func (e *Expense) IsGroupExpense() bool {
	return e.GroupID != nil
}

// IsPaidBy checks if the expense was paid by the given user
func (e *Expense) IsPaidBy(userID uuid.UUID) bool {
	for _, participant := range e.Participants {
		if participant.UserID == userID && participant.PaidAmount > 0 {
			return true
		}
	}
	return false
}

// GetParticipant returns the participant for a given user
func (e *Expense) GetParticipant(userID uuid.UUID) *ExpenseParticipant {
	for i := range e.Participants {
		if e.Participants[i].UserID == userID {
			return &e.Participants[i]
		}
	}
	return nil
}

// IsFullySettled checks if all participants have settled their dues
func (e *Expense) IsFullySettled() bool {
	for _, participant := range e.Participants {
		if !participant.IsSettled && participant.OwedAmount > 0 {
			return false
		}
	}
	return true
}
