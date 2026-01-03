package services

import (
	"errors"
	"personal-expense-splitting-settlement/internal/dto"
	"personal-expense-splitting-settlement/internal/models"
	"personal-expense-splitting-settlement/internal/repository"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AuthService interface {
	Register(req dto.RegisterRequest) (*dto.RegisterResponse, error)
	Login(req dto.LoginRequest) (*dto.LoginResponse, error)
	GetProfile(userID uuid.UUID) (*dto.UserResponse, error)
	UpdateProfile(userID uuid.UUID, req dto.UpdateProfileRequest) (*dto.UserResponse, error)
}

type authService struct {
	userRepo   repository.UserRepository
	jwtSecret  string
	jwtExpiry  time.Duration
	dataSecret string
}

func NewAuthService(userRepo repository.UserRepository, jwtSecret string, jwtExpiry time.Duration, dataSecret string) AuthService {
	return &authService{
		userRepo:   userRepo,
		jwtSecret:  jwtSecret,
		jwtExpiry:  jwtExpiry,
		dataSecret: dataSecret,
	}
}

func (s *authService) Register(req dto.RegisterRequest) (*dto.RegisterResponse, error) {
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
		Email:     req.Email,
		FirstName: req.FirstName,
		LastName:  req.LastName,
	}

	// Hash Password
	if err := user.HashPassword(req.Password); err != nil {
		return nil, err
	}

	// Encrypt Phone number
	if err := user.EncryptPhone(req.PhoneNumber, s.dataSecret); err != nil {
		return nil, err
	}

	// Save User
	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	getUser, err := s.userRepo.FindByEmail(req.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, models.ErrInvalidCredentials
		}
		return nil, err
	}

	// Generate JWT token
	// token, err := utils.GenerateToken(getUser.ID, s.jwtSecret, s.jwtExpiry)
	// if err != nil {
	// 	return nil, err
	// }

	res := &dto.RegisterResponse{
		User: dto.UserResponse{
			ID:        getUser.ID.String(),
			Email:     getUser.Email,
			FirstName: getUser.FirstName,
			LastName:  getUser.LastName,
			CreatedAt: getUser.CreatedAt,
		},
		Tokens: dto.TokenResponse{
			ExpiresIn: int(s.jwtExpiry.Seconds()),
		},
	}

	return res, nil
}

func (s *authService) Login(req dto.LoginRequest) (*dto.LoginResponse, error) {
	// Find User by email
	user, err := s.userRepo.FindByEmail(req.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, models.ErrInvalidCredentials
		}
		return nil, err
	}

	// Check Password
	if !user.CheckPassword(req.Password) {
		return nil, models.ErrInvalidCredentials
	}

	// Check if User is active
	if !user.IsActive {
		return nil, models.ErrUnauthorized
	}

	// Update last login timestamp
	now := time.Now()
	user.LastLoginAt = &now
	if err := s.userRepo.Update(user); err != nil {
		// Log error but don't fail login
	}

	return &dto.LoginResponse{
		User: dto.UserResponse{
			ID:        user.ID.String(),
			Email:     user.Email,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			CreatedAt: user.CreatedAt,
		},
		Tokens: dto.TokenResponse{
			ExpiresIn: int(s.jwtExpiry.Seconds()),
		},
	}, nil
}

func (s *authService) GetProfile(userID uuid.UUID) (*dto.UserResponse, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, err
	}

	return &dto.UserResponse{
		ID:        userID.String(),
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		CreatedAt: user.CreatedAt,
	}, nil
}

func (s *authService) UpdateProfile(userID uuid.UUID, req dto.UpdateProfileRequest) (*dto.UserResponse, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, err
	}

	// Update fields if provided
	if req.FirstName != nil {
		user.FirstName = *req.FirstName
	}
	if req.LastName != nil {
		user.LastName = *req.LastName
	}
	if req.DefaultCurrency != nil {
		user.DefaultCurrency = *req.DefaultCurrency
	}
	if req.ProfileImageURL != nil {
		user.ProfileImageURL = *req.ProfileImageURL
	}

	if err := s.userRepo.Update(user); err != nil {
		return nil, err
	}

	return &dto.UserResponse{
		ID:        userID.String(),
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		CreatedAt: user.CreatedAt,
	}, nil
}
