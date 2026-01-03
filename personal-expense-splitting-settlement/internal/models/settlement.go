package models

import (
	"time"

	"github.com/google/uuid"
)

// PaymentMethod represents different payment methods for settlements
type PaymentMethod string

const (
	PaymentMethodCash         PaymentMethod = "cash"
	PaymentMethodBankTransfer PaymentMethod = "bank_transfer"
	PaymentMethodUPI          PaymentMethod = "upi"
	PaymentMethodPaypal       PaymentMethod = "paypal"
	PaymentMethodVenmo        PaymentMethod = "venmo"
	PaymentMethodCreditCard   PaymentMethod = "credit_card"
	PaymentMethodDebitCard    PaymentMethod = "debit_card"
	PaymentMethodOther        PaymentMethod = "other"
)

// Settlement represents a debt settlement between two users
type Settlement struct {
	ID            uuid.UUID     `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	PayerID       uuid.UUID     `gorm:"type:uuid;not null;column:payer_id" json:"payer_id"`
	PayeeID       uuid.UUID     `gorm:"type:uuid;not null;column:payee_id" json:"payee_id"`
	Amount        float64       `gorm:"type:decimal(12,2);not null" json:"amount"`
	PaymentMethod PaymentMethod `gorm:"type:varchar(50);not null;default:'cash'" json:"payment_method"`
	Notes         string        `gorm:"type:text" json:"notes"`
	IsConfirmed   bool          `gorm:"not null;default:false;column:is_confirmed" json:"is_confirmed"`
	ConfirmedAt   *time.Time    `gorm:"column:confirmed_at" json:"confirmed_at,omitempty"`
	GroupID       *uuid.UUID    `gorm:"type:uuid;column:group_id" json:"group_id,omitempty"`
	CreatedAt     time.Time     `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt     time.Time     `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt     *time.Time    `gorm:"column:deleted_at" json:"deleted_at,omitempty"`

	// Relationships
	Payer User   `gorm:"foreignKey:PayerID" json:"payer,omitempty"`
	Payee User   `gorm:"foreignKey:PayeeID" json:"payee,omitempty"`
	Group *Group `gorm:"foreignKey:GroupID" json:"group,omitempty"`
}

// TableName specifies the table name for Settlement model
func (Settlement) TableName() string {
	return "settlements"
}

// IsValidPaymentMethod checks if the payment method is valid
func (s *Settlement) IsValidPaymentMethod() bool {
	validMethods := []PaymentMethod{
		PaymentMethodCash,
		PaymentMethodBankTransfer,
		PaymentMethodUPI,
		PaymentMethodPaypal,
		PaymentMethodVenmo,
		PaymentMethodCreditCard,
		PaymentMethodDebitCard,
		PaymentMethodOther,
	}

	for _, method := range validMethods {
		if s.PaymentMethod == method {
			return true
		}
	}
	return false
}

// CanConfirm checks if the settlement can be confirmed by the given user
func (s *Settlement) CanConfirm(userID uuid.UUID) bool {
	return s.PayeeID == userID && !s.IsConfirmed
}

// CanDelete checks if the settlement can be deleted by the given user
func (s *Settlement) CanDelete(userID uuid.UUID) bool {
	return s.PayerID == userID && !s.IsConfirmed
}

// Confirm confirms the settlement
func (s *Settlement) Confirm() {
	s.IsConfirmed = true
	now := time.Now()
	s.ConfirmedAt = &now
}

// Unconfirm unconfirms the settlement
func (s *Settlement) Unconfirm() {
	s.IsConfirmed = false
	s.ConfirmedAt = nil
}
