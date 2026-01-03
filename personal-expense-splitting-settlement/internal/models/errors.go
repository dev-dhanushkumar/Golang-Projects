package models

import "errors"

var (
	ErrInsufficientBalance = errors.New("insufficient wallet balance")
	ErrDuplicateReference  = errors.New("duplicate reference ID")
	ErrInvalidCredentials  = errors.New("invalid credentials")
	ErrUserAlreadyExists   = errors.New("user already exists")
	ErrCategoryNotFound    = errors.New("category not found")
	ErrBudgetNotFound      = errors.New("budget not found")
	ErrUnauthorized        = errors.New("unauthorized")
)
