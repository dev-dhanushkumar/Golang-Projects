package services

import (
	"errors"
	"math"
	"sort"
	"time"

	"personal-expense-splitting-settlement/internal/dto"
	"personal-expense-splitting-settlement/internal/repository"

	"github.com/google/uuid"
)

// BalanceService interface defines business logic for balance calculations
type BalanceService interface {
	GetBalanceSummary(userID uuid.UUID) (*dto.BalanceSummaryResponse, error)
	GetUserBalances(userID uuid.UUID) (*dto.BalanceSummaryResponse, error)
	GetGroupBalances(groupID, userID uuid.UUID) (*dto.GroupBalancesResponse, error)
	GetSettlementSuggestions(userID uuid.UUID) (*dto.SettlementSuggestionsResponse, error)
	GetGroupSettlementSuggestions(groupID, userID uuid.UUID) (*dto.SettlementSuggestionsResponse, error)
}

type balanceService struct {
	balanceRepo repository.BalanceRepository
	groupRepo   repository.GroupRepository
	userRepo    repository.UserRepository
}

// NewBalanceService creates a new balance service instance
func NewBalanceService(
	balanceRepo repository.BalanceRepository,
	groupRepo repository.GroupRepository,
	userRepo repository.UserRepository,
) BalanceService {
	return &balanceService{
		balanceRepo: balanceRepo,
		groupRepo:   groupRepo,
		userRepo:    userRepo,
	}
}

// GetBalanceSummary returns a summary of all balances for a user
func (s *balanceService) GetBalanceSummary(userID uuid.UUID) (*dto.BalanceSummaryResponse, error) {
	return s.GetUserBalances(userID)
}

// GetUserBalances calculates detailed balances for a user with all their connections
func (s *balanceService) GetUserBalances(userID uuid.UUID) (*dto.BalanceSummaryResponse, error) {
	balances, err := s.balanceRepo.CalculateAllUserBalances(userID)
	if err != nil {
		return nil, err
	}

	// Calculate summary
	var totalOwed, totalOwing float64
	for _, balance := range balances {
		if balance.Amount > 0 {
			totalOwed += balance.Amount
		} else {
			totalOwing += math.Abs(balance.Amount)
		}
	}

	return &dto.BalanceSummaryResponse{
		TotalOwed:   totalOwed,
		TotalOwing:  totalOwing,
		NetBalance:  totalOwed - totalOwing,
		Balances:    balances,
		LastUpdated: time.Now(),
	}, nil
}

// GetGroupBalances calculates balances for all members in a group
func (s *balanceService) GetGroupBalances(groupID, userID uuid.UUID) (*dto.GroupBalancesResponse, error) {
	// Verify user is a group member
	isMember, err := s.groupRepo.IsUserMember(groupID, userID)
	if err != nil {
		return nil, err
	}
	if !isMember {
		return nil, errors.New("you are not a member of this group")
	}

	// Get group details
	group, err := s.groupRepo.FindByID(groupID)
	if err != nil {
		return nil, err
	}

	// Calculate balances
	balances, err := s.balanceRepo.CalculateGroupBalances(groupID)
	if err != nil {
		return nil, err
	}

	// Calculate total expense
	var totalExpense float64
	for _, balance := range balances {
		totalExpense += balance.TotalPaid
	}

	return &dto.GroupBalancesResponse{
		GroupID:      groupID,
		GroupName:    group.Name,
		TotalExpense: totalExpense,
		Balances:     balances,
		LastUpdated:  time.Now(),
	}, nil
}

// GetSettlementSuggestions generates optimal settlement suggestions to minimize transactions
func (s *balanceService) GetSettlementSuggestions(userID uuid.UUID) (*dto.SettlementSuggestionsResponse, error) {
	// Get all balances for the user
	balances, err := s.balanceRepo.CalculateAllUserBalances(userID)
	if err != nil {
		return nil, err
	}

	if len(balances) == 0 {
		return &dto.SettlementSuggestionsResponse{
			Suggestions: []dto.SettlementSuggestion{},
			TotalAmount: 0,
			Message:     "You're all settled up!",
		}, nil
	}

	// Get requesting user info
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, err
	}
	userName := user.FirstName + " " + user.LastName

	// Convert balances to suggestions
	// For individual user balances:
	// - Positive amount means other user owes the requesting user
	// - Negative amount means requesting user owes the other user
	var suggestions []dto.SettlementSuggestion
	var totalAmount float64

	for _, balance := range balances {
		if math.Abs(balance.Amount) > 0.01 { // Small threshold for floating point
			if balance.Amount > 0 {
				// Other user owes the requesting user
				suggestions = append(suggestions, dto.SettlementSuggestion{
					From:     balance.UserID,
					FromName: balance.UserName,
					To:       userID,
					ToName:   userName,
					Amount:   math.Round(balance.Amount*100) / 100,
				})
				totalAmount += balance.Amount
			} else {
				// Requesting user owes the other user
				suggestions = append(suggestions, dto.SettlementSuggestion{
					From:     userID,
					FromName: userName,
					To:       balance.UserID,
					ToName:   balance.UserName,
					Amount:   math.Round(math.Abs(balance.Amount)*100) / 100,
				})
				totalAmount += math.Abs(balance.Amount)
			}
		}
	}

	message := "Here are suggested settlements to simplify your balances"
	if len(suggestions) == 0 {
		message = "You're all settled up!"
	}

	return &dto.SettlementSuggestionsResponse{
		Suggestions: suggestions,
		TotalAmount: totalAmount,
		Message:     message,
	}, nil
}

// GetGroupSettlementSuggestions generates settlement suggestions for a group
func (s *balanceService) GetGroupSettlementSuggestions(groupID, userID uuid.UUID) (*dto.SettlementSuggestionsResponse, error) {
	// Verify user is a group member
	isMember, err := s.groupRepo.IsUserMember(groupID, userID)
	if err != nil {
		return nil, err
	}
	if !isMember {
		return nil, errors.New("you are not a member of this group")
	}

	// Get group balances
	groupBalances, err := s.balanceRepo.CalculateGroupBalances(groupID)
	if err != nil {
		return nil, err
	}

	// Convert group balances to balance items for simplification
	balanceItems := make([]dto.BalanceItem, len(groupBalances))
	for i, gb := range groupBalances {
		balanceItems[i] = dto.BalanceItem{
			UserID:   gb.UserID,
			UserName: gb.UserName,
			Amount:   gb.NetBalance,
		}
	}

	// Simplify debts
	suggestions := simplifyDebts(balanceItems)

	// Calculate total amount
	var totalAmount float64
	for _, suggestion := range suggestions {
		totalAmount += suggestion.Amount
	}

	message := "Here are suggested settlements for the group"
	if len(suggestions) == 0 {
		message = "Group is all settled up!"
	}

	return &dto.SettlementSuggestionsResponse{
		Suggestions: suggestions,
		TotalAmount: totalAmount,
		Message:     message,
	}, nil
}

// simplifyDebts uses a greedy algorithm to minimize the number of transactions
// Algorithm: Always match the largest creditor with the largest debtor
func simplifyDebts(balances []dto.BalanceItem) []dto.SettlementSuggestion {
	if len(balances) == 0 {
		return []dto.SettlementSuggestion{}
	}

	// Separate creditors (positive balance) and debtors (negative balance)
	type personBalance struct {
		UserID   uuid.UUID
		UserName string
		Amount   float64
	}

	var creditors, debtors []personBalance
	for _, balance := range balances {
		if balance.Amount > 0.01 { // Small threshold to avoid floating point issues
			creditors = append(creditors, personBalance{
				UserID:   balance.UserID,
				UserName: balance.UserName,
				Amount:   balance.Amount,
			})
		} else if balance.Amount < -0.01 {
			debtors = append(debtors, personBalance{
				UserID:   balance.UserID,
				UserName: balance.UserName,
				Amount:   math.Abs(balance.Amount),
			})
		}
	}

	// Sort creditors and debtors by amount (descending)
	sort.Slice(creditors, func(i, j int) bool {
		return creditors[i].Amount > creditors[j].Amount
	})
	sort.Slice(debtors, func(i, j int) bool {
		return debtors[i].Amount > debtors[j].Amount
	})

	// Generate suggestions
	var suggestions []dto.SettlementSuggestion
	i, j := 0, 0

	for i < len(creditors) && j < len(debtors) {
		creditor := &creditors[i]
		debtor := &debtors[j]

		// Settle the minimum of what creditor is owed and debtor owes
		amount := math.Min(creditor.Amount, debtor.Amount)

		if amount > 0.01 { // Only create suggestion if amount is significant
			suggestions = append(suggestions, dto.SettlementSuggestion{
				From:     debtor.UserID,
				FromName: debtor.UserName,
				To:       creditor.UserID,
				ToName:   creditor.UserName,
				Amount:   math.Round(amount*100) / 100, // Round to 2 decimal places
			})
		}

		// Update balances
		creditor.Amount -= amount
		debtor.Amount -= amount

		// Move to next if current is settled
		if creditor.Amount < 0.01 {
			i++
		}
		if debtor.Amount < 0.01 {
			j++
		}
	}

	return suggestions
}
