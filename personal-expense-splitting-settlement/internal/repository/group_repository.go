package repository

import (
	"errors"
	"time"

	"personal-expense-splitting-settlement/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// GroupRepository interface defines methods for group data operations
type GroupRepository interface {
	Create(group *models.Group) error
	FindByID(id uuid.UUID) (*models.Group, error)
	FindByUserID(userID uuid.UUID) ([]models.Group, error)
	Update(group *models.Group) error
	Delete(id uuid.UUID) error

	// Member operations
	AddMember(member *models.GroupMember) error
	RemoveMember(groupID, userID uuid.UUID) error
	GetMembers(groupID uuid.UUID) ([]models.GroupMember, error)
	GetMember(groupID, userID uuid.UUID) (*models.GroupMember, error)
	UpdateMemberRole(groupID, userID uuid.UUID, role models.MemberRole) error
	IsUserMember(groupID, userID uuid.UUID) (bool, error)
	IsUserAdmin(groupID, userID uuid.UUID) (bool, error)
}

type groupRepository struct {
	db *gorm.DB
}

// NewGroupRepository creates a new group repository instance
func NewGroupRepository(db *gorm.DB) GroupRepository {
	return &groupRepository{db: db}
}

// Create creates a new group
func (r *groupRepository) Create(group *models.Group) error {
	return r.db.Create(group).Error
}

// FindByID finds a group by ID with members
func (r *groupRepository) FindByID(id uuid.UUID) (*models.Group, error) {
	var group models.Group
	err := r.db.
		Preload("Creator").
		Preload("Members", "left_at IS NULL").
		Preload("Members.User").
		Where("id = ?", id).
		First(&group).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("group not found")
		}
		return nil, err
	}
	return &group, nil
}

// FindByUserID finds all groups where the user is a member
func (r *groupRepository) FindByUserID(userID uuid.UUID) ([]models.Group, error) {
	var groups []models.Group

	err := r.db.
		Joins("JOIN group_members ON group_members.group_id = groups.id").
		Where("group_members.user_id = ? AND group_members.left_at IS NULL", userID).
		Preload("Creator").
		Preload("Members", "left_at IS NULL").
		Preload("Members.User").
		Order("groups.created_at DESC").
		Find(&groups).Error

	if err != nil {
		return nil, err
	}
	return groups, nil
}

// Update updates a group's information
func (r *groupRepository) Update(group *models.Group) error {
	return r.db.Model(group).Updates(group).Error
}

// Delete soft deletes a group
func (r *groupRepository) Delete(id uuid.UUID) error {
	return r.db.Where("id = ?", id).Delete(&models.Group{}).Error
}

// AddMember adds a member to a group
func (r *groupRepository) AddMember(member *models.GroupMember) error {
	// Check if user is already an active member
	var existingMember models.GroupMember
	err := r.db.Where("group_id = ? AND user_id = ? AND left_at IS NULL",
		member.GroupID, member.UserID).First(&existingMember).Error

	if err == nil {
		return errors.New("user is already a member of this group")
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	// Create new membership
	return r.db.Create(member).Error
}

// RemoveMember removes a member from a group (sets left_at)
func (r *groupRepository) RemoveMember(groupID, userID uuid.UUID) error {
	now := time.Now()
	result := r.db.Model(&models.GroupMember{}).
		Where("group_id = ? AND user_id = ? AND left_at IS NULL", groupID, userID).
		Update("left_at", now)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("member not found or already removed")
	}

	return nil
}

// GetMembers retrieves all active members of a group
func (r *groupRepository) GetMembers(groupID uuid.UUID) ([]models.GroupMember, error) {
	var members []models.GroupMember
	err := r.db.
		Preload("User").
		Where("group_id = ? AND left_at IS NULL", groupID).
		Order("joined_at ASC").
		Find(&members).Error

	if err != nil {
		return nil, err
	}
	return members, nil
}

// GetMember retrieves a specific member of a group
func (r *groupRepository) GetMember(groupID, userID uuid.UUID) (*models.GroupMember, error) {
	var member models.GroupMember
	err := r.db.
		Preload("User").
		Where("group_id = ? AND user_id = ? AND left_at IS NULL", groupID, userID).
		First(&member).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("member not found")
		}
		return nil, err
	}
	return &member, nil
}

// UpdateMemberRole updates a member's role in the group
func (r *groupRepository) UpdateMemberRole(groupID, userID uuid.UUID, role models.MemberRole) error {
	result := r.db.Model(&models.GroupMember{}).
		Where("group_id = ? AND user_id = ? AND left_at IS NULL", groupID, userID).
		Update("role", role)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("member not found")
	}

	return nil
}

// IsUserMember checks if a user is an active member of a group
func (r *groupRepository) IsUserMember(groupID, userID uuid.UUID) (bool, error) {
	var count int64
	err := r.db.Model(&models.GroupMember{}).
		Where("group_id = ? AND user_id = ? AND left_at IS NULL", groupID, userID).
		Count(&count).Error

	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// IsUserAdmin checks if a user is an admin of a group
func (r *groupRepository) IsUserAdmin(groupID, userID uuid.UUID) (bool, error) {
	var count int64
	err := r.db.Model(&models.GroupMember{}).
		Where("group_id = ? AND user_id = ? AND role = ? AND left_at IS NULL",
			groupID, userID, models.MemberRoleAdmin).
		Count(&count).Error

	if err != nil {
		return false, err
	}
	return count > 0, nil
}
