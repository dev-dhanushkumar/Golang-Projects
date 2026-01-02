package service

import (
	"digital-wallet-api/internal/dto"
	"digital-wallet-api/internal/models"
	"digital-wallet-api/internal/repository"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type WalletService interface {
	GetWallet(userID uuid.UUID) (*dto.WalletResponse, error)
	GetBalance(userID uuid.UUID) (*dto.BalanceResponse, error)
}

type walletService struct {
	walletRepo repository.WalletRepository
}

func NewWalletService(walletRepo repository.WalletRepository) WalletService {
	return &walletService{
		walletRepo: walletRepo,
	}
}

func (s *walletService) GetWallet(userID uuid.UUID) (*dto.WalletResponse, error) {
	wallet, err := s.walletRepo.FindByUserID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, models.ErrWalletNotFound
		}
		return nil, err
	}

	return &dto.WalletResponse{
		ID:       wallet.ID,
		UserID:   wallet.UserID,
		Balance:  wallet.Balance,
		Currency: wallet.Currency,
		IsActive: wallet.IsActive,
	}, nil
}

func (s *walletService) GetBalance(userID uuid.UUID) (*dto.BalanceResponse, error) {
	wallet, err := s.walletRepo.FindByUserID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, models.ErrWalletNotFound
		}
		return nil, err
	}

	return &dto.BalanceResponse{
		Balance:  wallet.Balance,
		Currency: wallet.Currency,
	}, nil
}
