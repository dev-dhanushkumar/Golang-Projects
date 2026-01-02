package models

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type Wallet struct {
	BaseModel
	UserID   uuid.UUID       `gorm:"type:uuid;not null;uniqueIndex" json:"user_id"`
	Balance  decimal.Decimal `gorm:"type:decimal(15,2);default:0.00" json:"balance"`
	Currency string          `gorm:"size:3;default:'USD'" json:"currency"`
	IsActive bool            `gorm:"default:true" json:"is_active"`

	// Relationships
	User         User          `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"user,omitempty"`
	Transactions []Transaction `gorm:"foreignKey:WalletID;constraint:OnDelete:CASCADE" json:"transactions,omitempty"`
}

func (Wallet) TableName() string {
	return "wallets"
}

// BeforeCreate hook
func (w *Wallet) BeforeCreate(tx *gorm.DB) error {
	if w.ID == uuid.Nil {
		w.ID = uuid.New()
	}
	return nil
}

// Credit adds money to wallet
func (w *Wallet) Credit(amount decimal.Decimal) {
	w.Balance = w.Balance.Add(amount)
}

// Debit subtracts money from wallet
func (w *Wallet) Debit(amount decimal.Decimal) error {
	if w.Balance.LessThan(amount) {
		return ErrInsufficientBalance
	}
	w.Balance = w.Balance.Sub(amount)
	return nil
}

// HasSufficientBalance checks if wallet has enough balance
func (w *Wallet) HasSufficientBalance(amount decimal.Decimal) bool {
	return w.Balance.GreaterThanOrEqual(amount)
}
