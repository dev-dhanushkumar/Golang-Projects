package services

import (
	"errors"
	"fmt"
	"time"

	"personal-expense-splitting-settlement/internal/dto"
	"personal-expense-splitting-settlement/internal/models"
	"personal-expense-splitting-settlement/internal/repository"

	"github.com/google/uuid"
)

// ExpenseService interface defines business logic for expense operations
type ExpenseService interface {
	CreateExpense(userID uuid.UUID, req *dto.CreateExpenseRequest) (*models.Expense, error)
	GetExpense(expenseID, userID uuid.UUID) (*models.Expense, error)
	GetUserExpenses(userID uuid.UUID, limit, offset int) ([]models.Expense, error)
	GetGroupExpenses(groupID, userID uuid.UUID, limit, offset int) ([]models.Expense, error)
	GetExpensesWithFilters(userID uuid.UUID, req *dto.ExpenseListRequest) ([]models.Expense, error)
	UpdateExpense(expenseID, userID uuid.UUID, req *dto.UpdateExpenseRequest) (*models.Expense, error)
	DeleteExpense(expenseID, userID uuid.UUID) error
}

type expenseService struct {
	expenseRepo     repository.ExpenseRepository
	participantRepo repository.ExpenseParticipantRepository
	groupRepo       repository.GroupRepository
	userRepo        repository.UserRepository
}

// NewExpenseService creates a new expense service instance
func NewExpenseService(
	expenseRepo repository.ExpenseRepository,
	participantRepo repository.ExpenseParticipantRepository,
	groupRepo repository.GroupRepository,
	userRepo repository.UserRepository,
) ExpenseService {
	return &expenseService{
		expenseRepo:     expenseRepo,
		participantRepo: participantRepo,
		groupRepo:       groupRepo,
		userRepo:        userRepo,
	}
}

// CreateExpense creates a new expense with participants
func (s *expenseService) CreateExpense(userID uuid.UUID, req *dto.CreateExpenseRequest) (*models.Expense, error) {
	// Validate category
	expense := &models.Expense{
		Description: req.Description,
		Amount:      req.Amount,
		Category:    req.Category,
		ReceiptURL:  req.ReceiptURL,
		CreatedBy:   userID,
		GroupID:     req.GroupID,
	}

	// Set date
	if req.Date != nil {
		expense.Date = req.Date.Time
	} else {
		expense.Date = time.Now()
	}

	// Validate date is not in future
	if expense.Date.After(time.Now()) {
		return nil, errors.New("expense date cannot be in the future")
	}

	if !expense.IsValidCategory() {
		return nil, errors.New("invalid expense category")
	}

	// If group expense, validate group membership
	if req.GroupID != nil {
		isMember, err := s.groupRepo.IsUserMember(*req.GroupID, userID)
		if err != nil {
			return nil, err
		}
		if !isMember {
			return nil, errors.New("you must be a member of the group to add expenses")
		}

		// Validate all participants are group members
		for _, participant := range req.Participants {
			isMember, err := s.groupRepo.IsUserMember(*req.GroupID, participant.UserID)
			if err != nil {
				return nil, err
			}
			if !isMember {
				return nil, fmt.Errorf("user %s is not a member of the group", participant.UserID)
			}
		}
	}

	// Verify all participant users exist
	for _, participant := range req.Participants {
		_, err := s.userRepo.FindByID(participant.UserID)
		if err != nil {
			return nil, fmt.Errorf("user not found: %s", participant.UserID)
		}
	}

	// Calculate split amounts based on split method
	participants, err := s.calculateSplits(req)
	if err != nil {
		return nil, err
	}

	// Create expense
	if err := s.expenseRepo.Create(expense); err != nil {
		return nil, err
	}

	// Create participants
	participantModels := make([]models.ExpenseParticipant, len(participants))
	for i, p := range participants {
		participantModels[i] = models.ExpenseParticipant{
			ExpenseID:  expense.ID,
			UserID:     p.UserID,
			PaidAmount: p.PaidAmount,
			OwedAmount: *p.OwedAmount,
			IsSettled:  false,
		}
	}

	if err := s.participantRepo.CreateBulk(participantModels); err != nil {
		return nil, err
	}

	// Reload expense with participants
	return s.expenseRepo.FindByID(expense.ID)
}

// calculateSplits calculates the owed amounts for each participant based on split method
func (s *expenseService) calculateSplits(req *dto.CreateExpenseRequest) ([]dto.ParticipantSplit, error) {
	totalAmount := req.Amount
	participants := req.Participants

	switch req.SplitMethod {
	case dto.SplitMethodEqual:
		return s.calculateEqualSplit(totalAmount, participants)
	case dto.SplitMethodExact:
		return s.calculateExactSplit(totalAmount, participants)
	case dto.SplitMethodPercentage:
		return s.calculatePercentageSplit(totalAmount, participants)
	case dto.SplitMethodShares:
		return s.calculateShareSplit(totalAmount, participants)
	default:
		return nil, errors.New("invalid split method")
	}
}

// calculateEqualSplit splits the amount equally among all participants
func (s *expenseService) calculateEqualSplit(totalAmount float64, participants []dto.ParticipantSplit) ([]dto.ParticipantSplit, error) {
	if len(participants) == 0 {
		return nil, errors.New("at least one participant is required")
	}

	perPersonAmount := totalAmount / float64(len(participants))

	// Validate total paid equals total amount
	totalPaid := 0.0
	for _, p := range participants {
		totalPaid += p.PaidAmount
	}

	if totalPaid != totalAmount {
		return nil, fmt.Errorf("total paid amount (%.2f) must equal expense amount (%.2f)", totalPaid, totalAmount)
	}

	result := make([]dto.ParticipantSplit, len(participants))
	for i, p := range participants {
		result[i] = dto.ParticipantSplit{
			UserID:     p.UserID,
			PaidAmount: p.PaidAmount,
			OwedAmount: &perPersonAmount,
		}
	}

	return result, nil
}

// calculateExactSplit uses exact amounts specified for each participant
func (s *expenseService) calculateExactSplit(totalAmount float64, participants []dto.ParticipantSplit) ([]dto.ParticipantSplit, error) {
	if len(participants) == 0 {
		return nil, errors.New("at least one participant is required")
	}

	totalOwed := 0.0
	totalPaid := 0.0

	for _, p := range participants {
		if p.OwedAmount == nil {
			return nil, errors.New("owed_amount is required for exact split method")
		}
		totalOwed += *p.OwedAmount
		totalPaid += p.PaidAmount
	}

	// Allow small rounding differences (0.01)
	if totalOwed < totalAmount-0.01 || totalOwed > totalAmount+0.01 {
		return nil, fmt.Errorf("total owed amount (%.2f) must equal expense amount (%.2f)", totalOwed, totalAmount)
	}

	if totalPaid < totalAmount-0.01 || totalPaid > totalAmount+0.01 {
		return nil, fmt.Errorf("total paid amount (%.2f) must equal expense amount (%.2f)", totalPaid, totalAmount)
	}

	return participants, nil
}

// calculatePercentageSplit splits based on percentages
func (s *expenseService) calculatePercentageSplit(totalAmount float64, participants []dto.ParticipantSplit) ([]dto.ParticipantSplit, error) {
	if len(participants) == 0 {
		return nil, errors.New("at least one participant is required")
	}

	totalPercentage := 0.0
	totalPaid := 0.0

	for _, p := range participants {
		if p.Percentage == nil {
			return nil, errors.New("percentage is required for percentage split method")
		}
		totalPercentage += *p.Percentage
		totalPaid += p.PaidAmount
	}

	// Allow small rounding differences
	if totalPercentage < 99.99 || totalPercentage > 100.01 {
		return nil, fmt.Errorf("total percentage (%.2f) must equal 100", totalPercentage)
	}

	if totalPaid < totalAmount-0.01 || totalPaid > totalAmount+0.01 {
		return nil, fmt.Errorf("total paid amount (%.2f) must equal expense amount (%.2f)", totalPaid, totalAmount)
	}

	result := make([]dto.ParticipantSplit, len(participants))
	for i, p := range participants {
		owedAmount := (totalAmount * (*p.Percentage)) / 100.0
		result[i] = dto.ParticipantSplit{
			UserID:     p.UserID,
			PaidAmount: p.PaidAmount,
			OwedAmount: &owedAmount,
		}
	}

	return result, nil
}

// calculateShareSplit splits based on shares
func (s *expenseService) calculateShareSplit(totalAmount float64, participants []dto.ParticipantSplit) ([]dto.ParticipantSplit, error) {
	if len(participants) == 0 {
		return nil, errors.New("at least one participant is required")
	}

	totalShares := 0
	totalPaid := 0.0

	for _, p := range participants {
		if p.Shares == nil {
			return nil, errors.New("shares is required for shares split method")
		}
		totalShares += *p.Shares
		totalPaid += p.PaidAmount
	}

	if totalShares == 0 {
		return nil, errors.New("total shares must be greater than 0")
	}

	if totalPaid < totalAmount-0.01 || totalPaid > totalAmount+0.01 {
		return nil, fmt.Errorf("total paid amount (%.2f) must equal expense amount (%.2f)", totalPaid, totalAmount)
	}

	result := make([]dto.ParticipantSplit, len(participants))
	for i, p := range participants {
		owedAmount := (totalAmount * float64(*p.Shares)) / float64(totalShares)
		result[i] = dto.ParticipantSplit{
			UserID:     p.UserID,
			PaidAmount: p.PaidAmount,
			OwedAmount: &owedAmount,
		}
	}

	return result, nil
}

// GetExpense retrieves an expense by ID
func (s *expenseService) GetExpense(expenseID, userID uuid.UUID) (*models.Expense, error) {
	expense, err := s.expenseRepo.FindByID(expenseID)
	if err != nil {
		return nil, err
	}

	// Verify user is a participant
	isParticipant := false
	for _, p := range expense.Participants {
		if p.UserID == userID {
			isParticipant = true
			break
		}
	}

	if !isParticipant {
		return nil, errors.New("you are not a participant in this expense")
	}

	return expense, nil
}

// GetUserExpenses retrieves all expenses for a user
func (s *expenseService) GetUserExpenses(userID uuid.UUID, limit, offset int) ([]models.Expense, error) {
	if limit == 0 {
		limit = 20 // Default limit
	}
	return s.expenseRepo.FindByUserID(userID, limit, offset)
}

// GetGroupExpenses retrieves all expenses for a group
func (s *expenseService) GetGroupExpenses(groupID, userID uuid.UUID, limit, offset int) ([]models.Expense, error) {
	// Verify user is a group member
	isMember, err := s.groupRepo.IsUserMember(groupID, userID)
	if err != nil {
		return nil, err
	}

	if !isMember {
		return nil, errors.New("you are not a member of this group")
	}

	if limit == 0 {
		limit = 20 // Default limit
	}

	return s.expenseRepo.FindByGroupID(groupID, limit, offset)
}

// GetExpensesWithFilters retrieves expenses with filters
func (s *expenseService) GetExpensesWithFilters(userID uuid.UUID, req *dto.ExpenseListRequest) ([]models.Expense, error) {
	limit := req.Limit
	if limit == 0 {
		limit = 20 // Default limit
	}

	// Parse group ID if provided
	var groupID *uuid.UUID
	if req.GroupID != "" {
		parsed, err := uuid.Parse(req.GroupID)
		if err != nil {
			return nil, fmt.Errorf("invalid group_id: %w", err)
		}
		groupID = &parsed
	}

	// Parse category if provided
	var category *models.ExpenseCategory
	if req.Category != "" {
		cat := models.ExpenseCategory(req.Category)
		category = &cat
	}

	// Parse dates if provided
	var startDate, endDate *time.Time
	if req.StartDate != "" {
		parsed, err := time.Parse("2006-01-02", req.StartDate)
		if err != nil {
			return nil, fmt.Errorf("invalid start_date format, expected YYYY-MM-DD: %w", err)
		}
		startDate = &parsed
	}
	if req.EndDate != "" {
		parsed, err := time.Parse("2006-01-02", req.EndDate)
		if err != nil {
			return nil, fmt.Errorf("invalid end_date format, expected YYYY-MM-DD: %w", err)
		}
		endDate = &parsed
	}

	return s.expenseRepo.FindWithFilters(
		userID,
		groupID,
		category,
		startDate,
		endDate,
		limit,
		req.Offset,
	)
}

// UpdateExpense updates an expense
func (s *expenseService) UpdateExpense(expenseID, userID uuid.UUID, req *dto.UpdateExpenseRequest) (*models.Expense, error) {
	expense, err := s.expenseRepo.FindByID(expenseID)
	if err != nil {
		return nil, err
	}

	// Only creator can update
	if expense.CreatedBy != userID {
		return nil, errors.New("only the creator can update this expense")
	}

	// Update fields if provided
	if req.Description != nil {
		expense.Description = *req.Description
	}
	if req.Amount != nil {
		return nil, errors.New("cannot update expense amount after creation")
	}
	if req.Category != nil {
		expense.Category = *req.Category
		if !expense.IsValidCategory() {
			return nil, errors.New("invalid expense category")
		}
	}
	if req.ReceiptURL != nil {
		expense.ReceiptURL = *req.ReceiptURL
	}
	if req.Date != nil {
		if req.Date.Time.After(time.Now()) {
			return nil, errors.New("expense date cannot be in the future")
		}
		expense.Date = req.Date.Time
	}

	// Save updates
	if err := s.expenseRepo.Update(expense); err != nil {
		return nil, err
	}

	// Reload expense with participants
	return s.expenseRepo.FindByID(expenseID)
}

// DeleteExpense deletes an expense
func (s *expenseService) DeleteExpense(expenseID, userID uuid.UUID) error {
	expense, err := s.expenseRepo.FindByID(expenseID)
	if err != nil {
		return err
	}

	// Only creator can delete
	if expense.CreatedBy != userID {
		return errors.New("only the creator can delete this expense")
	}

	// Check if any settlements have been made
	for _, p := range expense.Participants {
		if p.IsSettled {
			return errors.New("cannot delete expense with settled participants")
		}
	}

	return s.expenseRepo.Delete(expenseID)
}
