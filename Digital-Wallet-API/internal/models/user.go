package models

import (
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	BaseModel
	Email        string `gorm:"uniqueIndex;not null;size:255" json:"email" validate:"required,email"`
	PasswordHash string `gorm:"not null;size:255" json:"-"` // Don't expose in JSON
	FullName     string `gorm:"not null;size:255" json:"full_name" validate:"required,min=2,max=255"`
	Phone        string `gorm:"size:20" json:"phone"`
	IsActive     bool   `gorm:"default:true" json:"is_active"`

	// Relationships
	Wallet     *Wallet    `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"wallet,omitempty"`
	Categories []Category `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"categories,omitempty"`
	Budgets    []Budget   `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"budgets,omitempty"`
}

// TableName specifies the table name
func (User) TableName() string {
	return "users"
}

// HashPassword hashes the password before saving
func (u *User) HashPassword(password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.PasswordHash = string(hashedPassword)
	return nil
}

// CheckPassword verifies the password
func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password))
	return err == nil
}

// BeforeCreate hook - automatically called before creating user
func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	// Additional validations can be added here
	return nil
}
