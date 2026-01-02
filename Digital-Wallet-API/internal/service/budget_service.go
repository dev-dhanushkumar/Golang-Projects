package service

import (
	"digital-wallet-api/internal/dto"
	"digital-wallet-api/internal/models"
	"digital-wallet-api/internal/repository"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type BudgetService interface {
	CreateBudget(userID uuid.UUID, req dto.CreateBudgetRequest) (*dto.BudgetResponse, error)
	GetBudgets(userID uuid.UUID) ([]dto.BudgetResponse, error)
	GetBudget(userID uuid.UUID, budgetID uuid.UUID) (*dto.BudgetResponse, error)
	UpdateBudget(userID uuid.UUID, budgetID uuid.UUID, req dto.UpdateBudgetRequest) (*dto.BudgetResponse, error)
	DeleteBudget(userID uuid.UUID, budgetID uuid.UUID) error
	GetBudgetAlerts(userID uuid.UUID) ([]dto.BudgetAlert, error)
}

type budgetService struct {
	budgetRepo   repository.BudgetRepository
	categoryRepo repository.CategoryRepository
}

func NewBudgetService(budgetRepo repository.BudgetRepository, categoryRepo repository.CategoryRepository) BudgetService {
	return &budgetService{
		budgetRepo:   budgetRepo,
		categoryRepo: categoryRepo,
	}
}

func (s *budgetService) CreateBudget(userID uuid.UUID, req dto.CreateBudgetRequest) (*dto.BudgetResponse, error) {
	// Validate category if provided
	if req.CategoryID != nil {
		category, err := s.categoryRepo.FindByID(*req.CategoryID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, models.ErrCategoryNotFound
			}
			return nil, err
		}

		// Verify ownership
		if category.UserID != nil && *category.UserID != userID && !category.IsDefault {
			return nil, models.ErrUnauthorized
		}
	}

	// Calculate end date based on period
	var endDate time.Time
	if req.Period == string(models.BudgetPeriodWeekly) {
		endDate = req.StartDate.AddDate(0, 0, 7)
	} else {
		endDate = req.StartDate.AddDate(0, 1, 0)
	}

	budget := &models.Budget{
		UserID:      userID,
		CategoryID:  req.CategoryID,
		Amount:      req.Amount,
		SpentAmount: decimal.Zero,
		Period:      models.BudgetPeriod(req.Period),
		StartDate:   req.StartDate,
		EndDate:     endDate,
		IsActive:    true,
	}

	if err := s.budgetRepo.Create(budget); err != nil {
		return nil, err
	}

	return s.toBudgetResponse(budget), nil
}

func (s *budgetService) GetBudgets(userID uuid.UUID) ([]dto.BudgetResponse, error) {
	budgets, err := s.budgetRepo.FindByUserID(userID)
	if err != nil {
		return nil, err
	}

	var responses []dto.BudgetResponse
	for _, budget := range budgets {
		responses = append(responses, *s.toBudgetResponse(&budget))
	}

	return responses, nil
}

func (s *budgetService) GetBudget(userID uuid.UUID, budgetID uuid.UUID) (*dto.BudgetResponse, error) {
	budget, err := s.budgetRepo.FindByID(budgetID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, models.ErrBudgetNotFound
		}
		return nil, err
	}

	// Verify ownership
	if budget.UserID != userID {
		return nil, models.ErrUnauthorized
	}

	return s.toBudgetResponse(budget), nil
}

func (s *budgetService) UpdateBudget(userID uuid.UUID, budgetID uuid.UUID, req dto.UpdateBudgetRequest) (*dto.BudgetResponse, error) {
	budget, err := s.budgetRepo.FindByID(budgetID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, models.ErrBudgetNotFound
		}
		return nil, err
	}

	// Verify ownership
	if budget.UserID != userID {
		return nil, models.ErrUnauthorized
	}

	// Update fields
	if !req.Amount.IsZero() {
		budget.Amount = req.Amount
	}
	if req.IsActive != nil {
		budget.IsActive = *req.IsActive
	}

	if err := s.budgetRepo.Update(budget); err != nil {
		return nil, err
	}

	return s.toBudgetResponse(budget), nil
}

func (s *budgetService) DeleteBudget(userID uuid.UUID, budgetID uuid.UUID) error {
	budget, err := s.budgetRepo.FindByID(budgetID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.ErrBudgetNotFound
		}
		return err
	}

	// Verify ownership
	if budget.UserID != userID {
		return models.ErrUnauthorized
	}

	return s.budgetRepo.Delete(budgetID)
}

func (s *budgetService) GetBudgetAlerts(userID uuid.UUID) ([]dto.BudgetAlert, error) {
	budgets, err := s.budgetRepo.FindActiveByUserID(userID)
	if err != nil {
		return nil, err
	}

	var alerts []dto.BudgetAlert
	for _, budget := range budgets {
		if budget.IsExceeded() {
			alert := dto.BudgetAlert{
				BudgetID:       budget.ID,
				Amount:         budget.Amount,
				SpentAmount:    budget.SpentAmount,
				PercentageUsed: budget.PercentageUsed(),
				Status:         "exceeded",
				Message:        "Budget exceeded",
			}
			if budget.Category != nil {
				alert.CategoryName = budget.Category.Name
			}
			alerts = append(alerts, alert)
		} else if budget.IsNearLimit() {
			alert := dto.BudgetAlert{
				BudgetID:       budget.ID,
				Amount:         budget.Amount,
				SpentAmount:    budget.SpentAmount,
				PercentageUsed: budget.PercentageUsed(),
				Status:         "near_limit",
				Message:        "Budget near limit (80%)",
			}
			if budget.Category != nil {
				alert.CategoryName = budget.Category.Name
			}
			alerts = append(alerts, alert)
		}
	}

	return alerts, nil
}

func (s *budgetService) toBudgetResponse(budget *models.Budget) *dto.BudgetResponse {
	response := &dto.BudgetResponse{
		ID:              budget.ID,
		UserID:          budget.UserID,
		CategoryID:      budget.CategoryID,
		Amount:          budget.Amount,
		SpentAmount:     budget.SpentAmount,
		RemainingAmount: budget.RemainingAmount(),
		PercentageUsed:  budget.PercentageUsed(),
		Period:          string(budget.Period),
		StartDate:       budget.StartDate,
		EndDate:         budget.EndDate,
		IsActive:        budget.IsActive,
		IsExceeded:      budget.IsExceeded(),
		IsNearLimit:     budget.IsNearLimit(),
		CreatedAt:       budget.CreatedAt,
	}

	if budget.Category != nil {
		response.CategoryName = budget.Category.Name
	}

	return response
}
