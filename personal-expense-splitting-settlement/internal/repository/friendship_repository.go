package repository

import (
	"personal-expense-splitting-settlement/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// FriendshipRepository handles database operations for friendships
type FriendshipRepository interface {
	Create(friendship *models.Friendship) error
	FindByID(id uuid.UUID) (*models.Friendship, error)
	FindByUsers(userID1, userID2 uuid.UUID) (*models.Friendship, error)
	Update(friendship *models.Friendship) error
	Delete(id uuid.UUID) error
	GetFriendsByUserID(userID uuid.UUID, status models.FriendshipStatus) ([]models.Friendship, error)
	GetPendingRequestsSent(userID uuid.UUID) ([]models.Friendship, error)
	GetPendingRequestsReceived(userID uuid.UUID) ([]models.Friendship, error)
	CheckFriendshipExists(userID1, userID2 uuid.UUID) (bool, error)
	AreFriends(userID1, userID2 uuid.UUID) (bool, error)
}

type friendshipRepository struct {
	db *gorm.DB
}

// NewFriendshipRepository creates a new instance of FriendshipRepository
func NewFriendshipRepository(db *gorm.DB) FriendshipRepository {
	return &friendshipRepository{db: db}
}

// Create creates a new friendship record
func (r *friendshipRepository) Create(friendship *models.Friendship) error {
	// Ensure user_id_1 < user_id_2 for consistency
	if friendship.UserID1.String() > friendship.UserID2.String() {
		friendship.UserID1, friendship.UserID2 = friendship.UserID2, friendship.UserID1
	}
	return r.db.Create(friendship).Error
}

// FindByID finds a friendship by its ID
func (r *friendshipRepository) FindByID(id uuid.UUID) (*models.Friendship, error) {
	var friendship models.Friendship
	err := r.db.Preload("User1").Preload("User2").Preload("Requester").
		Where("id = ?", id).First(&friendship).Error
	if err != nil {
		return nil, err
	}
	return &friendship, nil
}

// FindByUsers finds a friendship between two users (order independent)
func (r *friendshipRepository) FindByUsers(userID1, userID2 uuid.UUID) (*models.Friendship, error) {
	var friendship models.Friendship

	// Normalize user IDs (smaller UUID first)
	if userID1.String() > userID2.String() {
		userID1, userID2 = userID2, userID1
	}

	err := r.db.Preload("User1").Preload("User2").Preload("Requester").
		Where("user_id_1 = ? AND user_id_2 = ?", userID1, userID2).
		First(&friendship).Error
	if err != nil {
		return nil, err
	}
	return &friendship, nil
}

// Update updates a friendship record
func (r *friendshipRepository) Update(friendship *models.Friendship) error {
	return r.db.Save(friendship).Error
}

// Delete deletes a friendship by ID
func (r *friendshipRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.Friendship{}, "id = ?", id).Error
}

// GetFriendsByUserID gets all friends of a user with a specific status
func (r *friendshipRepository) GetFriendsByUserID(userID uuid.UUID, status models.FriendshipStatus) ([]models.Friendship, error) {
	var friendships []models.Friendship
	err := r.db.Preload("User1").Preload("User2").Preload("Requester").
		Where("(user_id_1 = ? OR user_id_2 = ?) AND status = ?", userID, userID, status).
		Order("created_at DESC").
		Find(&friendships).Error
	return friendships, err
}

// GetPendingRequestsSent gets all pending friend requests sent by a user
func (r *friendshipRepository) GetPendingRequestsSent(userID uuid.UUID) ([]models.Friendship, error) {
	var friendships []models.Friendship
	err := r.db.Preload("User1").Preload("User2").Preload("Requester").
		Where("requested_by = ? AND status = ?", userID, models.FriendshipStatusPending).
		Order("created_at DESC").
		Find(&friendships).Error
	return friendships, err
}

// GetPendingRequestsReceived gets all pending friend requests received by a user
func (r *friendshipRepository) GetPendingRequestsReceived(userID uuid.UUID) ([]models.Friendship, error) {
	var friendships []models.Friendship
	err := r.db.Preload("User1").Preload("User2").Preload("Requester").
		Where("(user_id_1 = ? OR user_id_2 = ?) AND requested_by != ? AND status = ?",
			userID, userID, userID, models.FriendshipStatusPending).
		Order("created_at DESC").
		Find(&friendships).Error
	return friendships, err
}

// CheckFriendshipExists checks if a friendship already exists between two users
func (r *friendshipRepository) CheckFriendshipExists(userID1, userID2 uuid.UUID) (bool, error) {
	// Normalize user IDs
	if userID1.String() > userID2.String() {
		userID1, userID2 = userID2, userID1
	}

	var count int64
	err := r.db.Model(&models.Friendship{}).
		Where("user_id_1 = ? AND user_id_2 = ?", userID1, userID2).
		Count(&count).Error
	return count > 0, err
}

// AreFriends checks if two users are accepted friends
func (r *friendshipRepository) AreFriends(userID1, userID2 uuid.UUID) (bool, error) {
	// Normalize user IDs
	if userID1.String() > userID2.String() {
		userID1, userID2 = userID2, userID1
	}

	var count int64
	err := r.db.Model(&models.Friendship{}).
		Where("user_id_1 = ? AND user_id_2 = ? AND status = ?",
			userID1, userID2, models.FriendshipStatusAccepted).
		Count(&count).Error
	return count > 0, err
}
