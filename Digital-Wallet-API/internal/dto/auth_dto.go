package dto

import "github.com/google/uuid"

// RegisterRequest represents user registration request
type RegisterRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6,max=50"`
	FullName string `json:"full_name" validate:"required,min=2,max=255"`
	Phone    string `json:"phone" validate:"omitempty,min=10,max=20"`
}

// LoginRequest represents user login request
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// LoginResponse represents user login response
type LoginResponse struct {
	Token     string    `json:"token"`
	UserID    uuid.UUID `json:"user_id"`
	Email     string    `json:"email"`
	FullName  string    `json:"full_name"`
	ExpiresAt int64     `json:"expires_at"`
}

// UpdateProfileRequest represents profile update request
type UpdateProfileRequest struct {
	FullName string `json:"full_name" validate:"omitempty,min=2,max=255"`
	Phone    string `json:"phone" validate:"omitempty,min=10,max=20"`
}

// UserResponse represents user data in response
type UserResponse struct {
	ID       uuid.UUID `json:"id"`
	Email    string    `json:"email"`
	FullName string    `json:"full_name"`
	Phone    string    `json:"phone"`
	IsActive bool      `json:"is_active"`
}
