package models

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type Transfer struct {
	BaseModel
	FromWalletID uuid.UUID         `gorm:"type:uuid;not null;index" json:"from_wallet_id"`
	ToWalletID   uuid.UUID         `gorm:"type:uuid;not null;index" json:"to_wallet_id"`
	Amount       decimal.Decimal   `gorm:"type:decimal(15,2);not null" json:"amount" validate:"required,gt=0"`
	Description  string            `gorm:"type:text" json:"description"`
	ReferenceID  string            `gorm:"uniqueIndex;size:100" json:"reference_id"`
	Status       TransactionStatus `gorm:"size:20;default:'completed'" json:"status"`

	// Relationships
	FromWallet Wallet `gorm:"foreignKey:FromWalletID;constraint:OnDelete:CASCADE" json:"from_wallet,omitempty"`
	ToWallet   Wallet `gorm:"foreignKey:ToWalletID;constraint:OnDelete:CASCADE" json:"to_wallet,omitempty"`
}

func (Transfer) TableName() string {
	return "transfers"
}

// BeforeCreate hook
func (t *Transfer) BeforeCreate(tx *gorm.DB) error {
	if t.ID == uuid.Nil {
		t.ID = uuid.New()
	}
	if t.ReferenceID == "" {
		t.ReferenceID = uuid.New().String()
	}
	return nil
}

// Validate checks if transfer is valid
func (t *Transfer) Validate() error {
	if t.FromWalletID == t.ToWalletID {
		return ErrSameWalletTransfer
	}
	if t.Amount.LessThanOrEqual(decimal.Zero) {
		return ErrInvalidAmount
	}
	return nil
}
