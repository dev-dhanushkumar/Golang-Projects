package service

import (
	"digital-wallet-api/internal/dto"
	"digital-wallet-api/internal/models"
	"digital-wallet-api/internal/repository"
	"digital-wallet-api/pkg/utils"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type AuthService interface {
	Register(req dto.RegisterRequest) (*dto.UserResponse, error)
	Login(req dto.LoginRequest) (*dto.LoginResponse, error)
	GetProfile(userID uuid.UUID) (*dto.UserResponse, error)
	UpdateProfile(userID uuid.UUID, req dto.UpdateProfileRequest) (*dto.UserResponse, error)
}

type authService struct {
	userRepo   repository.UserRepository
	walletRepo repository.WalletRepository
	jwtSecret  string
	jwtExpiry  time.Duration
}

func NewAuthService(userRepo repository.UserRepository, walletRepo repository.WalletRepository, jwtSecret string, jwtExpiry time.Duration) AuthService {
	return &authService{
		userRepo:   userRepo,
		walletRepo: walletRepo,
		jwtSecret:  jwtSecret,
		jwtExpiry:  jwtExpiry,
	}
}

func (s *authService) Register(req dto.RegisterRequest) (*dto.UserResponse, error) {
	// Check if user already exists
	exists, err := s.userRepo.Exists(req.Email)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, models.ErrUserAlreadyExists
	}

	// Create user
	user := &models.User{
		Email:    req.Email,
		FullName: req.FullName,
		Phone:    req.Phone,
		IsActive: true,
	}

	// Hash password
	if err := user.HashPassword(req.Password); err != nil {
		return nil, err
	}

	// Save user
	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	// Create wallet for user
	wallet := &models.Wallet{
		UserID:   user.ID,
		Balance:  decimal.Zero,
		Currency: "USD",
		IsActive: true,
	}

	if err := s.walletRepo.Create(wallet); err != nil {
		return nil, err
	}

	return &dto.UserResponse{
		ID:       user.ID,
		Email:    user.Email,
		FullName: user.FullName,
		Phone:    user.Phone,
		IsActive: user.IsActive,
	}, nil
}

func (s *authService) Login(req dto.LoginRequest) (*dto.LoginResponse, error) {
	// Find user by email
	user, err := s.userRepo.FindByEmail(req.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, models.ErrInvalidCredentials
		}
		return nil, err
	}

	// Check password
	if !user.CheckPassword(req.Password) {
		return nil, models.ErrInvalidCredentials
	}

	// Check if user is active
	if !user.IsActive {
		return nil, models.ErrUnauthorized
	}

	// Generate JWT token
	token, err := utils.GenerateToken(user.ID, user.Email, s.jwtSecret, s.jwtExpiry)
	if err != nil {
		return nil, err
	}

	return &dto.LoginResponse{
		Token:     token,
		UserID:    user.ID,
		Email:     user.Email,
		FullName:  user.FullName,
		ExpiresAt: time.Now().Add(s.jwtExpiry).Unix(),
	}, nil
}

func (s *authService) GetProfile(userID uuid.UUID) (*dto.UserResponse, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, models.ErrUserNotFound
		}
		return nil, err
	}

	return &dto.UserResponse{
		ID:       user.ID,
		Email:    user.Email,
		FullName: user.FullName,
		Phone:    user.Phone,
		IsActive: user.IsActive,
	}, nil
}

func (s *authService) UpdateProfile(userID uuid.UUID, req dto.UpdateProfileRequest) (*dto.UserResponse, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, models.ErrUserNotFound
		}
		return nil, err
	}

	// Update fields
	if req.FullName != "" {
		user.FullName = req.FullName
	}
	if req.Phone != "" {
		user.Phone = req.Phone
	}

	if err := s.userRepo.Update(user); err != nil {
		return nil, err
	}

	return &dto.UserResponse{
		ID:       user.ID,
		Email:    user.Email,
		FullName: user.FullName,
		Phone:    user.Phone,
		IsActive: user.IsActive,
	}, nil
}
