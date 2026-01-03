package services

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"personal-expense-splitting-settlement/internal/models"
	"personal-expense-splitting-settlement/internal/repository"
	"personal-expense-splitting-settlement/pkg/utils"
	"time"

	"github.com/google/uuid"
)

type SessionService interface {
	CreateNewSession(userID uuid.UUID, ip, userAgent string) (accessToken, refreshToken string, err error)
	RefreshSession(oldRefreshToken string, ip, userAgent string) (newAccessToken, newRefreshToken string, err error)
	TerminateSession(sessionID uuid.UUID) error
	GetUserSessions(userID uuid.UUID) ([]models.UserSession, error)
}

type sessionService struct {
	sessionRepo repository.SessionRepository
	jwtSecret   string
	jwtExpiry   time.Duration
}

func NewSessionService(repo repository.SessionRepository, jwtSecret string, jwtExpiry time.Duration) SessionService {
	return &sessionService{
		sessionRepo: repo,
		jwtSecret:   jwtSecret,
		jwtExpiry:   jwtExpiry,
	}
}

func (s *sessionService) CreateNewSession(userId uuid.UUID, ip, userAgent string) (string, string, error) {
	// Generate the Session ID
	sessionID := uuid.New()

	// Generate tokens
	accessToken, _ := utils.GenerateToken(userId, sessionID, s.jwtSecret, s.jwtExpiry)
	refreshToken, _ := utils.GenerateRandomString(32)

	// Hash token for secure
	accessHash := s.hashToken(accessToken)
	refreshHash := s.hashToken(refreshToken)

	session := &models.UserSession{
		ID:               sessionID,
		UserID:           userId,
		TokenHash:        accessHash,
		RefreshTokenHash: refreshHash,
		IPAddress:        ip,
		UserAgent:        userAgent,
		ExpireAt:         time.Now().Add(time.Hour * 24 * 7),
		CreatedAt:        time.Now(),
	}

	if err := s.sessionRepo.CreateSession(session); err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func (s *sessionService) RefreshSession(oldRefreshToken string, ip, userAgent string) (string, string, error) {
	oldHash := s.hashToken(oldRefreshToken)

	session, err := s.sessionRepo.GetSessionByRefreshToken(oldHash)
	if err != nil {
		return "", "", errors.New("invalid or expired session")
	}

	s.sessionRepo.RevokeSession(session.ID)

	return s.CreateNewSession(session.UserID, ip, userAgent)
}

func (s *sessionService) TerminateSession(sessionID uuid.UUID) error {
	return s.sessionRepo.RevokeSession(sessionID)
}

func (s *sessionService) hashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}
func (s *sessionService) GetUserSessions(userID uuid.UUID) ([]models.UserSession, error) {
	return s.sessionRepo.GetActiveSessionByUserID(userID)
}
