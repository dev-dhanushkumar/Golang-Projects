package repository

import (
	"personal-expense-splitting-settlement/internal/models"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SessionRepository interface {
	CreateSession(session *models.UserSession) error
	GetSessionByRefreshToken(hash string) (*models.UserSession, error)
	RevokeSession(id uuid.UUID) error
	DeleteExpiredSession() error
	GetActiveSessionByUserID(userID uuid.UUID) ([]models.UserSession, error)
}

type sessionRepository struct {
	db *gorm.DB
}

// Create New session
func NewSesionRepository(db *gorm.DB) SessionRepository {
	return &sessionRepository{
		db: db,
	}
}

func (r *sessionRepository) CreateSession(session *models.UserSession) error {
	return r.db.Create(session).Error
}

func (r *sessionRepository) GetSessionByRefreshToken(hash string) (*models.UserSession, error) {
	var session models.UserSession

	err := r.db.Where("refresh_token_hash = ? AND revoked_at IS NULL AND expire_at > ?",
		hash, time.Now()).First(&session).Error

	if err != nil {
		return nil, err
	}

	return &session, nil
}

func (r *sessionRepository) RevokeSession(id uuid.UUID) error {
	return r.db.Model(&models.UserSession{}).
		Where("id = ?", id).
		Update("revoked_at", time.Now()).Error
}

func (r *sessionRepository) DeleteExpiredSession() error {
	return r.db.Where("expire_at <= ? OR revoked_at IS NOT NULL", time.Now()).
		Delete(&models.UserSession{}).Error
}

func (r *sessionRepository) GetActiveSessionByUserID(userID uuid.UUID) ([]models.UserSession, error) {
	var sessions []models.UserSession

	err := r.db.Where("user_id = ? AND revoked_at IS NULL AND expire_at > ?",
		userID, time.Now()).Find(&sessions).Error
	return sessions, err
}
