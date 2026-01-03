package repository

import (
	"errors"
	"time"

	"personal-expense-splitting-settlement/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ExpenseRepository interface defines methods for expense data operations
type ExpenseRepository interface {
	Create(expense *models.Expense) error
	FindByID(id uuid.UUID) (*models.Expense, error)
	FindByUserID(userID uuid.UUID, limit, offset int) ([]models.Expense, error)
	FindByGroupID(groupID uuid.UUID, limit, offset int) ([]models.Expense, error)
	FindWithFilters(userID uuid.UUID, groupID *uuid.UUID, category *models.ExpenseCategory, startDate, endDate *time.Time, limit, offset int) ([]models.Expense, error)
	Update(expense *models.Expense) error
	Delete(id uuid.UUID) error
	CountByUserID(userID uuid.UUID) (int64, error)
	CountByGroupID(groupID uuid.UUID) (int64, error)
}

type expenseRepository struct {
	db *gorm.DB
}

// NewExpenseRepository creates a new expense repository instance
func NewExpenseRepository(db *gorm.DB) ExpenseRepository {
	return &expenseRepository{db: db}
}

// Create creates a new expense
func (r *expenseRepository) Create(expense *models.Expense) error {
	return r.db.Create(expense).Error
}

// FindByID finds an expense by ID with participants
func (r *expenseRepository) FindByID(id uuid.UUID) (*models.Expense, error) {
	var expense models.Expense
	err := r.db.
		Preload("Creator").
		Preload("Group").
		Preload("Participants").
		Preload("Participants.User").
		Where("id = ?", id).
		First(&expense).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("expense not found")
		}
		return nil, err
	}
	return &expense, nil
}

// FindByUserID finds all expenses where the user is a participant
func (r *expenseRepository) FindByUserID(userID uuid.UUID, limit, offset int) ([]models.Expense, error) {
	var expenses []models.Expense

	query := r.db.
		Joins("JOIN expense_participants ON expense_participants.expense_id = expenses.id").
		Where("expense_participants.user_id = ?", userID).
		Preload("Creator").
		Preload("Group").
		Preload("Participants").
		Preload("Participants.User").
		Order("expenses.date DESC, expenses.created_at DESC")

	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	err := query.Find(&expenses).Error
	if err != nil {
		return nil, err
	}
	return expenses, nil
}

// FindByGroupID finds all expenses for a group
func (r *expenseRepository) FindByGroupID(groupID uuid.UUID, limit, offset int) ([]models.Expense, error) {
	var expenses []models.Expense

	query := r.db.
		Where("group_id = ?", groupID).
		Preload("Creator").
		Preload("Group").
		Preload("Participants").
		Preload("Participants.User").
		Order("date DESC, created_at DESC")

	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	err := query.Find(&expenses).Error
	if err != nil {
		return nil, err
	}
	return expenses, nil
}

// FindWithFilters finds expenses with various filters
func (r *expenseRepository) FindWithFilters(
	userID uuid.UUID,
	groupID *uuid.UUID,
	category *models.ExpenseCategory,
	startDate, endDate *time.Time,
	limit, offset int,
) ([]models.Expense, error) {
	var expenses []models.Expense

	query := r.db.
		Joins("JOIN expense_participants ON expense_participants.expense_id = expenses.id").
		Where("expense_participants.user_id = ?", userID)

	if groupID != nil {
		query = query.Where("expenses.group_id = ?", *groupID)
	}

	if category != nil {
		query = query.Where("expenses.category = ?", *category)
	}

	if startDate != nil {
		query = query.Where("expenses.date >= ?", *startDate)
	}

	if endDate != nil {
		query = query.Where("expenses.date <= ?", *endDate)
	}

	query = query.
		Preload("Creator").
		Preload("Group").
		Preload("Participants").
		Preload("Participants.User").
		Order("expenses.date DESC, expenses.created_at DESC")

	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	err := query.Find(&expenses).Error
	if err != nil {
		return nil, err
	}
	return expenses, nil
}

// Update updates an expense's information
func (r *expenseRepository) Update(expense *models.Expense) error {
	return r.db.Model(expense).Updates(expense).Error
}

// Delete soft deletes an expense
func (r *expenseRepository) Delete(id uuid.UUID) error {
	return r.db.Where("id = ?", id).Delete(&models.Expense{}).Error
}

// CountByUserID counts expenses for a user
func (r *expenseRepository) CountByUserID(userID uuid.UUID) (int64, error) {
	var count int64
	err := r.db.Model(&models.Expense{}).
		Joins("JOIN expense_participants ON expense_participants.expense_id = expenses.id").
		Where("expense_participants.user_id = ?", userID).
		Count(&count).Error
	return count, err
}

// CountByGroupID counts expenses for a group
func (r *expenseRepository) CountByGroupID(groupID uuid.UUID) (int64, error) {
	var count int64
	err := r.db.Model(&models.Expense{}).
		Where("group_id = ?", groupID).
		Count(&count).Error
	return count, err
}
