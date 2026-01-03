package dto

import (
	"time"

	"github.com/google/uuid"
)

// FriendRequestDTO represents a request to send a friend request
type FriendRequestDTO struct {
	FriendEmail string `json:"friend_email" binding:"required,email"`
}

// FriendResponseDTO represents the response data for a friendship
type FriendResponseDTO struct {
	ID          uuid.UUID `json:"id"`
	FriendID    uuid.UUID `json:"friend_id"`
	FriendEmail string    `json:"friend_email"`
	FriendName  string    `json:"friend_name"`
	Status      string    `json:"status"`
	RequestedBy uuid.UUID `json:"requested_by"`
	IsRequester bool      `json:"is_requester"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// FriendListResponseDTO represents a list of friends
type FriendListResponseDTO struct {
	Friends []FriendResponseDTO `json:"friends"`
	Total   int                 `json:"total"`
}

// PendingFriendRequestDTO represents a pending friend request with full user details
type PendingFriendRequestDTO struct {
	ID           uuid.UUID `json:"id"`
	RequesterID  uuid.UUID `json:"requester_id"`
	Email        string    `json:"email"`
	FirstName    string    `json:"first_name"`
	LastName     string    `json:"last_name"`
	ProfileImage *string   `json:"profile_image,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
}

// PendingRequestsResponseDTO represents the response for pending requests
type PendingRequestsResponseDTO struct {
	Sent     []PendingFriendRequestDTO `json:"sent"`
	Received []PendingFriendRequestDTO `json:"received"`
}
