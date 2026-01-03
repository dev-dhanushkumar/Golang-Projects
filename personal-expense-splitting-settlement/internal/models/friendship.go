package models

import (
	"time"

	"github.com/google/uuid"
)

// FriendshipStatus represents the status of a friendship
type FriendshipStatus string

const (
	FriendshipStatusPending  FriendshipStatus = "pending"
	FriendshipStatusAccepted FriendshipStatus = "accepted"
	FriendshipStatusRejected FriendshipStatus = "rejected"
	FriendshipStatusBlocked  FriendshipStatus = "blocked"
)

// Friendship represents a friendship relationship between two users
type Friendship struct {
	ID          uuid.UUID        `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	UserID1     uuid.UUID        `gorm:"column:user_id_1;type:uuid;not null" json:"user_id_1"`
	UserID2     uuid.UUID        `gorm:"column:user_id_2;type:uuid;not null" json:"user_id_2"`
	Status      FriendshipStatus `gorm:"type:varchar(20);not null" json:"status"`
	RequestedBy uuid.UUID        `gorm:"column:requested_by;type:uuid;not null" json:"requested_by"`
	CreatedAt   time.Time        `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time        `gorm:"autoUpdateTime" json:"updated_at"`

	// Virtual fields for relationships
	User1     *User `gorm:"foreignKey:UserID1" json:"user_1,omitempty"`
	User2     *User `gorm:"foreignKey:UserID2" json:"user_2,omitempty"`
	Requester *User `gorm:"foreignKey:RequestedBy" json:"requester,omitempty"`
}

// TableName specifies the table name for Friendship model
func (Friendship) TableName() string {
	return "friendships"
}

// IsAccepted checks if the friendship is accepted
func (f *Friendship) IsAccepted() bool {
	return f.Status == FriendshipStatusAccepted
}

// IsPending checks if the friendship is pending
func (f *Friendship) IsPending() bool {
	return f.Status == FriendshipStatusPending
}

// IsBlocked checks if the friendship is blocked
func (f *Friendship) IsBlocked() bool {
	return f.Status == FriendshipStatusBlocked
}

// GetOtherUserID returns the other user's ID in the friendship
func (f *Friendship) GetOtherUserID(currentUserID uuid.UUID) uuid.UUID {
	if f.UserID1 == currentUserID {
		return f.UserID2
	}
	return f.UserID1
}

// IsRequester checks if the given user is the one who requested the friendship
func (f *Friendship) IsRequester(userID uuid.UUID) bool {
	return f.RequestedBy == userID
}
