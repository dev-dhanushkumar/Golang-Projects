package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type BudgetPeriod string

const (
	BudgetPeriodWeekly  BudgetPeriod = "weekly"
	BudgetPeriodMonthly BudgetPeriod = "monthly"
)

type Budget struct {
	BaseModel
	UserID      uuid.UUID       `gorm:"type:uuid;not null;index" json:"user_id"`
	CategoryID  *uuid.UUID      `gorm:"type:uuid;index" json:"category_id,omitempty"` // Null means overall budget
	Amount      decimal.Decimal `gorm:"type:decimal(15,2);not null" json:"amount" validate:"required,gt=0"`
	SpentAmount decimal.Decimal `gorm:"type:decimal(15,2);default:0.00" json:"spent_amount"`
	Period      BudgetPeriod    `gorm:"not null;size:20" json:"period" validate:"required,oneof=weekly monthly"`
	StartDate   time.Time       `gorm:"not null;index" json:"start_date"`
	EndDate     time.Time       `gorm:"not null;index" json:"end_date"`
	IsActive    bool            `gorm:"default:true" json:"is_active"`

	// Relationships
	User     User      `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"user,omitempty"`
	Category *Category `gorm:"foreignKey:CategoryID;constraint:OnDelete:CASCADE" json:"category,omitempty"`
}

func (Budget) TableName() string {
	return "budgets"
}

// BeforeCreate hook
func (b *Budget) BeforeCreate(tx *gorm.DB) error {
	if b.ID == uuid.Nil {
		b.ID = uuid.New()
	}
	return nil
}

// RemainingAmount returns remaining budget
func (b *Budget) RemainingAmount() decimal.Decimal {
	return b.Amount.Sub(b.SpentAmount)
}

// PercentageUsed returns budget usage percentage
func (b *Budget) PercentageUsed() float64 {
	if b.Amount.IsZero() {
		return 0
	}
	percentage := b.SpentAmount.Div(b.Amount).Mul(decimal.NewFromInt(100))
	result, _ := percentage.Float64()
	return result
}

// IsExceeded checks if budget is exceeded
func (b *Budget) IsExceeded() bool {
	return b.SpentAmount.GreaterThan(b.Amount)
}

// IsNearLimit checks if budget is near limit (80%)
func (b *Budget) IsNearLimit() bool {
	return b.PercentageUsed() >= 80.0
}

// IsExpired checks if budget period has ended
func (b *Budget) IsExpired() bool {
	return time.Now().After(b.EndDate)
}
