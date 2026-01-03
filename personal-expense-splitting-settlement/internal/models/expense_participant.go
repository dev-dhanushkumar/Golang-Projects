package models

import (
	"time"

	"github.com/google/uuid"
)

// ExpenseParticipant represents a user's participation in an expense
type ExpenseParticipant struct {
	ID         uuid.UUID  `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	ExpenseID  uuid.UUID  `gorm:"type:uuid;not null;column:expense_id" json:"expense_id"`
	UserID     uuid.UUID  `gorm:"type:uuid;not null;column:user_id" json:"user_id"`
	PaidAmount float64    `gorm:"type:decimal(12,2);not null;default:0" json:"paid_amount"`
	OwedAmount float64    `gorm:"type:decimal(12,2);not null;default:0" json:"owed_amount"`
	IsSettled  bool       `gorm:"not null;default:false;column:is_settled" json:"is_settled"`
	SettledAt  *time.Time `gorm:"column:settled_at" json:"settled_at,omitempty"`
	CreatedAt  time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt  time.Time  `gorm:"autoUpdateTime" json:"updated_at"`

	// Relationships
	Expense Expense `gorm:"foreignKey:ExpenseID" json:"expense,omitempty"`
	User    User    `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// TableName specifies the table name for ExpenseParticipant model
func (ExpenseParticipant) TableName() string {
	return "expense_participants"
}

// GetNetAmount returns the net amount (paid - owed)
// Positive value means user is owed money
// Negative value means user owes money
func (ep *ExpenseParticipant) GetNetAmount() float64 {
	return ep.PaidAmount - ep.OwedAmount
}

// IsCreditor checks if the participant is a creditor (paid more than owed)
func (ep *ExpenseParticipant) IsCreditor() bool {
	return ep.GetNetAmount() > 0
}

// IsDebtor checks if the participant is a debtor (owes more than paid)
func (ep *ExpenseParticipant) IsDebtor() bool {
	return ep.GetNetAmount() < 0
}

// IsBalanced checks if the participant has paid exactly what they owe
func (ep *ExpenseParticipant) IsBalanced() bool {
	return ep.GetNetAmount() == 0
}

// MarkAsSettled marks the participant's share as settled
func (ep *ExpenseParticipant) MarkAsSettled() {
	ep.IsSettled = true
	now := time.Now()
	ep.SettledAt = &now
}

// MarkAsUnsettled marks the participant's share as unsettled
func (ep *ExpenseParticipant) MarkAsUnsettled() {
	ep.IsSettled = false
	ep.SettledAt = nil
}

// GetAmountToSettle returns the amount that needs to be settled
// Only returns positive values (amount owed)
func (ep *ExpenseParticipant) GetAmountToSettle() float64 {
	netAmount := ep.GetNetAmount()
	if netAmount < 0 {
		return -netAmount // Return positive value
	}
	return 0
}
