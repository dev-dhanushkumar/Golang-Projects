package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type TransactionType string
type TransactionStatus string

const (
	TransactionTypeCredit TransactionType = "credit"
	TransactionTypeDebit  TransactionType = "debit"

	TransactionStatusCompleted TransactionStatus = "completed"
	TransactionStatusPending   TransactionStatus = "pending"
	TransactionStatusFailed    TransactionStatus = "failed"
)

type Transaction struct {
	BaseModel
	WalletID        uuid.UUID         `gorm:"type:uuid;not null;index" json:"wallet_id"`
	CategoryID      *uuid.UUID        `gorm:"type:uuid;index" json:"category_id,omitempty"`
	Type            TransactionType   `gorm:"not null;size:20;index" json:"type" validate:"required,oneof=credit debit"`
	Amount          decimal.Decimal   `gorm:"type:decimal(15,2);not null" json:"amount" validate:"required,gt=0"`
	BalanceAfter    decimal.Decimal   `gorm:"type:decimal(15,2);not null" json:"balance_after"`
	Description     string            `gorm:"type:text" json:"description"`
	ReferenceID     string            `gorm:"uniqueIndex;size:100" json:"reference_id"` // For idempotency
	Status          TransactionStatus `gorm:"size:20;default:'completed';index" json:"status"`
	TransactionDate time.Time         `gorm:"index" json:"transaction_date"`
	Metadata        string            `gorm:"type:jsonb" json:"metadata,omitempty"` // For PostgreSQL JSONB

	// Relationships
	Wallet   Wallet    `gorm:"foreignKey:WalletID;constraint:OnDelete:CASCADE" json:"wallet,omitempty"`
	Category *Category `gorm:"foreignKey:CategoryID;constraint:OnDelete:SET NULL" json:"category,omitempty"`
}

func (Transaction) TableName() string {
	return "transactions"
}

// BeforeCreate hook - set transaction date and ID
func (t *Transaction) BeforeCreate(tx *gorm.DB) error {
	if t.ID == uuid.Nil {
		t.ID = uuid.New()
	}
	if t.TransactionDate.IsZero() {
		t.TransactionDate = time.Now()
	}
	if t.ReferenceID == "" {
		t.ReferenceID = uuid.New().String()
	}
	return nil
}

// IsCredit checks if transaction is credit
func (t *Transaction) IsCredit() bool {
	return t.Type == TransactionTypeCredit
}

// IsDebit checks if transaction is debit
func (t *Transaction) IsDebit() bool {
	return t.Type == TransactionTypeDebit
}
