package dto

import "time"

type RegisterRequest struct {
	Email       string `json:"email" validate:"required,email"`
	Password    string `json:"password" validate:"required,min=8,max=50,password"`
	FirstName   string `json:"first_name" validate:"required,min=2,max=255"`
	LastName    string `json:"last_name" validate:"required,min=2,max=255"`
	PhoneNumber string `json:"phone_number" validate:"omitempty,min=10,max=20"`
}

type RegisterResponse struct {
	User   UserResponse  `json:"user"`
	Tokens TokenResponse `json:"tokens"`
}

type UserResponse struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	CreatedAt time.Time `json:"created_at"`
}

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	User   UserResponse  `json:"user"`
	Tokens TokenResponse `json:"tokens"`
}

type UpdateProfileRequest struct {
	FirstName       *string `json:"first_name" validate:"omitempty,min=2,max=255"`
	LastName        *string `json:"last_name" validate:"omitempty,min=2,max=255"`
	DefaultCurrency *string `json:"default_currency" validate:"omitempty,len=3"`
	ProfileImageURL *string `json:"profile_image_url" validate:"omitempty,url"`
}
