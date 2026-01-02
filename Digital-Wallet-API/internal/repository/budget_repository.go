package repository

import (
	"digital-wallet-api/internal/models"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BudgetRepository interface {
	Create(budget *models.Budget) error
	FindByID(id uuid.UUID) (*models.Budget, error)
	FindByUserID(userID uuid.UUID) ([]models.Budget, error)
	FindActiveByUserID(userID uuid.UUID) ([]models.Budget, error)
	Update(budget *models.Budget) error
	Delete(id uuid.UUID) error
	UpdateSpentAmount(budgetID uuid.UUID, amount interface{}) error
}

type budgetRepository struct {
	db *gorm.DB
}

func NewBudgetRepository(db *gorm.DB) BudgetRepository {
	return &budgetRepository{db: db}
}

func (r *budgetRepository) Create(budget *models.Budget) error {
	return r.db.Create(budget).Error
}

func (r *budgetRepository) FindByID(id uuid.UUID) (*models.Budget, error) {
	var budget models.Budget
	err := r.db.Preload("Category").Where("id = ?", id).First(&budget).Error
	if err != nil {
		return nil, err
	}
	return &budget, nil
}

func (r *budgetRepository) FindByUserID(userID uuid.UUID) ([]models.Budget, error) {
	var budgets []models.Budget
	err := r.db.Preload("Category").Where("user_id = ?", userID).Order("created_at DESC").Find(&budgets).Error
	return budgets, err
}

func (r *budgetRepository) FindActiveByUserID(userID uuid.UUID) ([]models.Budget, error) {
	var budgets []models.Budget
	now := time.Now()
	err := r.db.Preload("Category").
		Where("user_id = ? AND is_active = ? AND start_date <= ? AND end_date >= ?", userID, true, now, now).
		Find(&budgets).Error
	return budgets, err
}

func (r *budgetRepository) Update(budget *models.Budget) error {
	return r.db.Save(budget).Error
}

func (r *budgetRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.Budget{}, "id = ?", id).Error
}

func (r *budgetRepository) UpdateSpentAmount(budgetID uuid.UUID, amount interface{}) error {
	return r.db.Model(&models.Budget{}).Where("id = ?", budgetID).Update("spent_amount", amount).Error
}
