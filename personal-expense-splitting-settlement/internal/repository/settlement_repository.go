package repository

import (
	"personal-expense-splitting-settlement/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// SettlementRepository defines methods for settlement data access
type SettlementRepository interface {
	Create(settlement *models.Settlement) error
	FindByID(id uuid.UUID) (*models.Settlement, error)
	FindByUserID(userID uuid.UUID, limit, offset int) ([]models.Settlement, error)
	FindByUsers(user1ID, user2ID uuid.UUID, limit, offset int) ([]models.Settlement, error)
	FindByGroupID(groupID uuid.UUID, limit, offset int) ([]models.Settlement, error)
	Update(settlement *models.Settlement) error
	ConfirmSettlement(id, userID uuid.UUID) error
	Delete(id uuid.UUID) error
	CountByUserID(userID uuid.UUID) (int64, error)
}

type settlementRepository struct {
	db *gorm.DB
}

// NewSettlementRepository creates a new settlement repository instance
func NewSettlementRepository(db *gorm.DB) SettlementRepository {
	return &settlementRepository{db: db}
}

// Create creates a new settlement
func (r *settlementRepository) Create(settlement *models.Settlement) error {
	return r.db.Create(settlement).Error
}

// FindByID finds a settlement by ID with related data
func (r *settlementRepository) FindByID(id uuid.UUID) (*models.Settlement, error) {
	var settlement models.Settlement
	err := r.db.Preload("Payer").
		Preload("Payee").
		Preload("Group").
		First(&settlement, "id = ? AND deleted_at IS NULL", id).Error
	if err != nil {
		return nil, err
	}
	return &settlement, nil
}

// FindByUserID finds all settlements involving a user (as payer or payee)
func (r *settlementRepository) FindByUserID(userID uuid.UUID, limit, offset int) ([]models.Settlement, error) {
	var settlements []models.Settlement
	err := r.db.Where("(payer_id = ? OR payee_id = ?) AND deleted_at IS NULL", userID, userID).
		Preload("Payer").
		Preload("Payee").
		Preload("Group").
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&settlements).Error
	if err != nil {
		return nil, err
	}
	return settlements, nil
}

// FindByUsers finds all settlements between two specific users
func (r *settlementRepository) FindByUsers(user1ID, user2ID uuid.UUID, limit, offset int) ([]models.Settlement, error) {
	var settlements []models.Settlement
	err := r.db.Where(
		"((payer_id = ? AND payee_id = ?) OR (payer_id = ? AND payee_id = ?)) AND deleted_at IS NULL",
		user1ID, user2ID, user2ID, user1ID,
	).
		Preload("Payer").
		Preload("Payee").
		Preload("Group").
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&settlements).Error
	if err != nil {
		return nil, err
	}
	return settlements, nil
}

// FindByGroupID finds all settlements for a group
func (r *settlementRepository) FindByGroupID(groupID uuid.UUID, limit, offset int) ([]models.Settlement, error) {
	var settlements []models.Settlement
	err := r.db.Where("group_id = ? AND deleted_at IS NULL", groupID).
		Preload("Payer").
		Preload("Payee").
		Preload("Group").
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&settlements).Error
	if err != nil {
		return nil, err
	}
	return settlements, nil
}

// Update updates a settlement
func (r *settlementRepository) Update(settlement *models.Settlement) error {
	return r.db.Save(settlement).Error
}

// ConfirmSettlement confirms a settlement by ID if the user is the payee
func (r *settlementRepository) ConfirmSettlement(id, userID uuid.UUID) error {
	return r.db.Model(&models.Settlement{}).
		Where("id = ? AND payee_id = ? AND is_confirmed = ? AND deleted_at IS NULL", id, userID, false).
		Update("is_confirmed", true).Error
}

// Delete soft deletes a settlement
func (r *settlementRepository) Delete(id uuid.UUID) error {
	return r.db.Model(&models.Settlement{}).
		Where("id = ?", id).
		Update("deleted_at", gorm.Expr("CURRENT_TIMESTAMP")).Error
}

// CountByUserID counts settlements for a user
func (r *settlementRepository) CountByUserID(userID uuid.UUID) (int64, error) {
	var count int64
	err := r.db.Model(&models.Settlement{}).
		Where("(payer_id = ? OR payee_id = ?) AND deleted_at IS NULL", userID, userID).
		Count(&count).Error
	return count, err
}
