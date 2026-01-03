package repository

import (
	"errors"

	"personal-expense-splitting-settlement/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ExpenseParticipantRepository defines methods for expense participant data access
type ExpenseParticipantRepository interface {
	Create(participant *models.ExpenseParticipant) error
	CreateBulk(participants []models.ExpenseParticipant) error
	FindByID(id uuid.UUID) (*models.ExpenseParticipant, error)
	FindByExpenseID(expenseID uuid.UUID) ([]models.ExpenseParticipant, error)
	FindByUserID(userID uuid.UUID) ([]models.ExpenseParticipant, error)
	FindByExpenseAndUser(expenseID, userID uuid.UUID) (*models.ExpenseParticipant, error)
	Update(participant *models.ExpenseParticipant) error
	UpdateSettlementStatus(id uuid.UUID, isSettled bool) error
	Delete(id uuid.UUID) error
	DeleteByExpenseID(expenseID uuid.UUID) error
	GetUnsettledByUser(userID uuid.UUID) ([]models.ExpenseParticipant, error)
}

type expenseParticipantRepository struct {
	db *gorm.DB
}

// NewExpenseParticipantRepository creates a new expense participant repository instance
func NewExpenseParticipantRepository(db *gorm.DB) ExpenseParticipantRepository {
	return &expenseParticipantRepository{db: db}
}

// Create creates a new expense participant
func (r *expenseParticipantRepository) Create(participant *models.ExpenseParticipant) error {
	return r.db.Create(participant).Error
}

// CreateBulk creates multiple expense participants efficiently
func (r *expenseParticipantRepository) CreateBulk(participants []models.ExpenseParticipant) error {
	if len(participants) == 0 {
		return errors.New("no participants to create")
	}
	return r.db.Create(&participants).Error
}

// FindByID finds an expense participant by ID
func (r *expenseParticipantRepository) FindByID(id uuid.UUID) (*models.ExpenseParticipant, error) {
	var participant models.ExpenseParticipant
	err := r.db.Preload("User").Preload("Expense").First(&participant, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &participant, nil
}

// FindByExpenseID finds all participants for an expense
func (r *expenseParticipantRepository) FindByExpenseID(expenseID uuid.UUID) ([]models.ExpenseParticipant, error) {
	var participants []models.ExpenseParticipant
	err := r.db.Where("expense_id = ?", expenseID).
		Preload("User").
		Find(&participants).Error
	if err != nil {
		return nil, err
	}
	return participants, nil
}

// FindByUserID finds all expense participations for a user
func (r *expenseParticipantRepository) FindByUserID(userID uuid.UUID) ([]models.ExpenseParticipant, error) {
	var participants []models.ExpenseParticipant
	err := r.db.Where("user_id = ?", userID).
		Preload("Expense").
		Preload("Expense.Creator").
		Preload("Expense.Group").
		Find(&participants).Error
	if err != nil {
		return nil, err
	}
	return participants, nil
}

// FindByExpenseAndUser finds a specific participant for an expense and user
func (r *expenseParticipantRepository) FindByExpenseAndUser(expenseID, userID uuid.UUID) (*models.ExpenseParticipant, error) {
	var participant models.ExpenseParticipant
	err := r.db.Where("expense_id = ? AND user_id = ?", expenseID, userID).
		Preload("User").
		Preload("Expense").
		First(&participant).Error
	if err != nil {
		return nil, err
	}
	return &participant, nil
}

// Update updates an expense participant
func (r *expenseParticipantRepository) Update(participant *models.ExpenseParticipant) error {
	return r.db.Save(participant).Error
}

// UpdateSettlementStatus updates the settlement status of a participant
func (r *expenseParticipantRepository) UpdateSettlementStatus(id uuid.UUID, isSettled bool) error {
	return r.db.Model(&models.ExpenseParticipant{}).
		Where("id = ?", id).
		Update("is_settled", isSettled).Error
}

// Delete deletes an expense participant
func (r *expenseParticipantRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.ExpenseParticipant{}, "id = ?", id).Error
}

// DeleteByExpenseID deletes all participants for an expense
func (r *expenseParticipantRepository) DeleteByExpenseID(expenseID uuid.UUID) error {
	return r.db.Where("expense_id = ?", expenseID).Delete(&models.ExpenseParticipant{}).Error
}

// GetUnsettledByUser gets all unsettled expense participations for a user
func (r *expenseParticipantRepository) GetUnsettledByUser(userID uuid.UUID) ([]models.ExpenseParticipant, error) {
	var participants []models.ExpenseParticipant
	err := r.db.Where("user_id = ? AND is_settled = ? AND owed_amount > 0", userID, false).
		Preload("Expense").
		Preload("Expense.Creator").
		Preload("Expense.Group").
		Preload("User").
		Find(&participants).Error
	if err != nil {
		return nil, err
	}
	return participants, nil
}
