package repository

import (
	"digital-wallet-api/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type WalletRepository interface {
	Create(wallet *models.Wallet) error
	FindByID(id uuid.UUID) (*models.Wallet, error)
	FindByUserID(userID uuid.UUID) (*models.Wallet, error)
	Update(wallet *models.Wallet) error
	UpdateBalance(walletID uuid.UUID, newBalance interface{}) error
}

type walletRepository struct {
	db *gorm.DB
}

func NewWalletRepository(db *gorm.DB) WalletRepository {
	return &walletRepository{db: db}
}

func (r *walletRepository) Create(wallet *models.Wallet) error {
	return r.db.Create(wallet).Error
}

func (r *walletRepository) FindByID(id uuid.UUID) (*models.Wallet, error) {
	var wallet models.Wallet
	err := r.db.Where("id = ?", id).First(&wallet).Error
	if err != nil {
		return nil, err
	}
	return &wallet, nil
}

func (r *walletRepository) FindByUserID(userID uuid.UUID) (*models.Wallet, error) {
	var wallet models.Wallet
	err := r.db.Where("user_id = ?", userID).First(&wallet).Error
	if err != nil {
		return nil, err
	}
	return &wallet, nil
}

func (r *walletRepository) Update(wallet *models.Wallet) error {
	return r.db.Save(wallet).Error
}

func (r *walletRepository) UpdateBalance(walletID uuid.UUID, newBalance interface{}) error {
	return r.db.Model(&models.Wallet{}).Where("id = ?", walletID).Update("balance", newBalance).Error
}
