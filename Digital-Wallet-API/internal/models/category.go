package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CategoryType string

const (
	CategoryTypeIncome  CategoryType = "income"
	CategoryTypeExpense CategoryType = "expense"
)

type Category struct {
	BaseModel
	UserID    *uuid.UUID   `gorm:"type:uuid;index" json:"user_id,omitempty"` // Nullable for default categories
	Name      string       `gorm:"not null;size:100" json:"name" validate:"required,min=2,max=100"`
	Type      CategoryType `gorm:"not null;size:20" json:"type" validate:"required,oneof=income expense"`
	Icon      string       `gorm:"size:50" json:"icon"`
	IsDefault bool         `gorm:"default:false" json:"is_default"` // System-wide categories

	// Relationships
	User         *User         `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"user,omitempty"`
	Transactions []Transaction `gorm:"foreignKey:CategoryID;constraint:OnDelete:SET NULL" json:"transactions,omitempty"`
	Budgets      []Budget      `gorm:"foreignKey:CategoryID;constraint:OnDelete:CASCADE" json:"budgets,omitempty"`
}

func (Category) TableName() string {
	return "categories"
}

// BeforeCreate hook
func (c *Category) BeforeCreate(tx *gorm.DB) error {
	if c.ID == uuid.Nil {
		c.ID = uuid.New()
	}
	return nil
}

// DefaultCategories returns system-wide default categories
func DefaultCategories() []Category {
	return []Category{
		{Name: "Salary", Type: CategoryTypeIncome, Icon: "💰", IsDefault: true},
		{Name: "Freelance", Type: CategoryTypeIncome, Icon: "💼", IsDefault: true},
		{Name: "Investment", Type: CategoryTypeIncome, Icon: "📈", IsDefault: true},
		{Name: "Food & Dining", Type: CategoryTypeExpense, Icon: "🍔", IsDefault: true},
		{Name: "Transportation", Type: CategoryTypeExpense, Icon: "🚗", IsDefault: true},
		{Name: "Shopping", Type: CategoryTypeExpense, Icon: "🛍️", IsDefault: true},
		{Name: "Entertainment", Type: CategoryTypeExpense, Icon: "🎬", IsDefault: true},
		{Name: "Bills & Utilities", Type: CategoryTypeExpense, Icon: "📱", IsDefault: true},
		{Name: "Healthcare", Type: CategoryTypeExpense, Icon: "🏥", IsDefault: true},
		{Name: "Education", Type: CategoryTypeExpense, Icon: "📚", IsDefault: true},
		{Name: "Other", Type: CategoryTypeExpense, Icon: "📌", IsDefault: true},
	}
}
