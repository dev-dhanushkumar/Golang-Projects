package services

import (
	"errors"

	"personal-expense-splitting-settlement/internal/dto"
	"personal-expense-splitting-settlement/internal/models"
	"personal-expense-splitting-settlement/internal/repository"

	"github.com/google/uuid"
)

// SettlementService interface defines business logic for settlement operations
type SettlementService interface {
	CreateSettlement(payerID uuid.UUID, req *dto.CreateSettlementRequest) (*models.Settlement, error)
	GetSettlement(settlementID, userID uuid.UUID) (*models.Settlement, error)
	GetUserSettlements(userID uuid.UUID, limit, offset int) ([]models.Settlement, error)
	GetSettlementsBetweenUsers(user1ID, user2ID uuid.UUID, limit, offset int) ([]models.Settlement, error)
	GetGroupSettlements(groupID, userID uuid.UUID, limit, offset int) ([]models.Settlement, error)
	UpdateSettlement(settlementID, userID uuid.UUID, req *dto.UpdateSettlementRequest) (*models.Settlement, error)
	ConfirmSettlement(settlementID, userID uuid.UUID) (*models.Settlement, error)
	DeleteSettlement(settlementID, userID uuid.UUID) error
}

type settlementService struct {
	settlementRepo repository.SettlementRepository
	userRepo       repository.UserRepository
	groupRepo      repository.GroupRepository
}

// NewSettlementService creates a new settlement service instance
func NewSettlementService(
	settlementRepo repository.SettlementRepository,
	userRepo repository.UserRepository,
	groupRepo repository.GroupRepository,
) SettlementService {
	return &settlementService{
		settlementRepo: settlementRepo,
		userRepo:       userRepo,
		groupRepo:      groupRepo,
	}
}

// CreateSettlement creates a new settlement
func (s *settlementService) CreateSettlement(payerID uuid.UUID, req *dto.CreateSettlementRequest) (*models.Settlement, error) {
	// Verify payee exists
	_, err := s.userRepo.FindByID(req.PayeeID)
	if err != nil {
		return nil, errors.New("payee not found")
	}

	// Prevent self-settlement
	if payerID == req.PayeeID {
		return nil, errors.New("cannot create settlement with yourself")
	}

	// If group settlement, verify both users are members
	if req.GroupID != nil {
		isPayerMember, err := s.groupRepo.IsUserMember(*req.GroupID, payerID)
		if err != nil {
			return nil, err
		}
		if !isPayerMember {
			return nil, errors.New("payer must be a member of the group")
		}

		isPayeeMember, err := s.groupRepo.IsUserMember(*req.GroupID, req.PayeeID)
		if err != nil {
			return nil, err
		}
		if !isPayeeMember {
			return nil, errors.New("payee must be a member of the group")
		}
	}

	settlement := &models.Settlement{
		PayerID:       payerID,
		PayeeID:       req.PayeeID,
		Amount:        req.Amount,
		PaymentMethod: req.PaymentMethod,
		Notes:         req.Notes,
		GroupID:       req.GroupID,
		IsConfirmed:   false,
	}

	if !settlement.IsValidPaymentMethod() {
		return nil, errors.New("invalid payment method")
	}

	if err := s.settlementRepo.Create(settlement); err != nil {
		return nil, err
	}

	// Reload settlement with relationships
	return s.settlementRepo.FindByID(settlement.ID)
}

// GetSettlement retrieves a settlement by ID
func (s *settlementService) GetSettlement(settlementID, userID uuid.UUID) (*models.Settlement, error) {
	settlement, err := s.settlementRepo.FindByID(settlementID)
	if err != nil {
		return nil, err
	}

	// Verify user is involved in the settlement
	if settlement.PayerID != userID && settlement.PayeeID != userID {
		return nil, errors.New("you are not involved in this settlement")
	}

	return settlement, nil
}

// GetUserSettlements retrieves all settlements for a user
func (s *settlementService) GetUserSettlements(userID uuid.UUID, limit, offset int) ([]models.Settlement, error) {
	if limit == 0 {
		limit = 20 // Default limit
	}
	return s.settlementRepo.FindByUserID(userID, limit, offset)
}

// GetSettlementsBetweenUsers retrieves settlements between two users
func (s *settlementService) GetSettlementsBetweenUsers(user1ID, user2ID uuid.UUID, limit, offset int) ([]models.Settlement, error) {
	if limit == 0 {
		limit = 20
	}
	return s.settlementRepo.FindByUsers(user1ID, user2ID, limit, offset)
}

// GetGroupSettlements retrieves all settlements for a group
func (s *settlementService) GetGroupSettlements(groupID, userID uuid.UUID, limit, offset int) ([]models.Settlement, error) {
	// Verify user is a group member
	isMember, err := s.groupRepo.IsUserMember(groupID, userID)
	if err != nil {
		return nil, err
	}
	if !isMember {
		return nil, errors.New("you are not a member of this group")
	}

	if limit == 0 {
		limit = 20
	}

	return s.settlementRepo.FindByGroupID(groupID, limit, offset)
}

// UpdateSettlement updates a settlement
func (s *settlementService) UpdateSettlement(settlementID, userID uuid.UUID, req *dto.UpdateSettlementRequest) (*models.Settlement, error) {
	settlement, err := s.settlementRepo.FindByID(settlementID)
	if err != nil {
		return nil, err
	}

	// Only payer can update before confirmation
	if settlement.PayerID != userID {
		return nil, errors.New("only the payer can update this settlement")
	}

	// Cannot update confirmed settlement
	if settlement.IsConfirmed {
		return nil, errors.New("cannot update a confirmed settlement")
	}

	// Update fields if provided
	if req.PaymentMethod != nil {
		settlement.PaymentMethod = *req.PaymentMethod
		if !settlement.IsValidPaymentMethod() {
			return nil, errors.New("invalid payment method")
		}
	}
	if req.Notes != nil {
		settlement.Notes = *req.Notes
	}

	if err := s.settlementRepo.Update(settlement); err != nil {
		return nil, err
	}

	// Reload settlement with relationships
	return s.settlementRepo.FindByID(settlementID)
}

// ConfirmSettlement confirms a settlement (only payee can confirm)
func (s *settlementService) ConfirmSettlement(settlementID, userID uuid.UUID) (*models.Settlement, error) {
	settlement, err := s.settlementRepo.FindByID(settlementID)
	if err != nil {
		return nil, err
	}

	if !settlement.CanConfirm(userID) {
		return nil, errors.New("only the payee can confirm this settlement, and it must not already be confirmed")
	}

	if err := s.settlementRepo.ConfirmSettlement(settlementID, userID); err != nil {
		return nil, err
	}

	// Reload settlement
	return s.settlementRepo.FindByID(settlementID)
}

// DeleteSettlement deletes a settlement
func (s *settlementService) DeleteSettlement(settlementID, userID uuid.UUID) error {
	settlement, err := s.settlementRepo.FindByID(settlementID)
	if err != nil {
		return err
	}

	if !settlement.CanDelete(userID) {
		return errors.New("only the payer can delete this settlement, and it must not be confirmed")
	}

	return s.settlementRepo.Delete(settlementID)
}
