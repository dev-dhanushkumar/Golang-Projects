package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserSession struct {
	ID               uuid.UUID  `gorm:"type:uuid;primaryKey"`
	UserID           uuid.UUID  `gorm:"type:uuid;index"`
	TokenHash        string     `gorm:"index" json:"-"`
	RefreshTokenHash string     `json:"-"`
	IPAddress        string     `json:"ip_address"`
	UserAgent        string     `json:"user_agent"`
	ExpireAt         time.Time  `json:"expires_at"`
	CreatedAt        time.Time  `json:"created_at"`
	RevokedAt        *time.Time `json:"revoked_at"`
}

// BeforeCreate hook - automatically called before creating user
func (s *UserSession) BeforeCreate(tx *gorm.DB) error {
	if s.ID == uuid.Nil {
		s.ID = uuid.New()
	}
	// Additional validations can be added here
	return nil
}
