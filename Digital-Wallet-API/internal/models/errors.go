package models

import "errors"

var (
	ErrInsufficientBalance = errors.New("insufficient wallet balance")
	ErrInvalidAmount       = errors.New("invalid amount: must be greater than zero")
	ErrSameWalletTransfer  = errors.New("cannot transfer to the same wallet")
	ErrWalletNotFound      = errors.New("wallet not found")
	ErrUserNotFound        = errors.New("user not found")
	ErrTransactionFailed   = errors.New("transaction failed")
	ErrDuplicateReference  = errors.New("duplicate reference ID")
	ErrInvalidCredentials  = errors.New("invalid credentials")
	ErrUserAlreadyExists   = errors.New("user already exists")
	ErrCategoryNotFound    = errors.New("category not found")
	ErrBudgetNotFound      = errors.New("budget not found")
	ErrUnauthorized        = errors.New("unauthorized")
)
