package services

import (
	"errors"
	"fmt"
	"personal-expense-splitting-settlement/internal/dto"
	"personal-expense-splitting-settlement/internal/models"
	"personal-expense-splitting-settlement/internal/repository"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// FriendshipService handles business logic for friendships
type FriendshipService interface {
	SendFriendRequest(currentUserID uuid.UUID, friendEmail string) error
	AcceptFriendRequest(currentUserID uuid.UUID, friendshipID uuid.UUID) error
	RejectFriendRequest(currentUserID uuid.UUID, friendshipID uuid.UUID) error
	BlockUser(currentUserID uuid.UUID, friendshipID uuid.UUID) error
	RemoveFriend(currentUserID uuid.UUID, friendshipID uuid.UUID) error
	GetFriends(currentUserID uuid.UUID) (*dto.FriendListResponseDTO, error)
	GetPendingRequests(currentUserID uuid.UUID) (*dto.PendingRequestsResponseDTO, error)
}

type friendshipService struct {
	friendshipRepo repository.FriendshipRepository
	userRepo       repository.UserRepository
}

// NewFriendshipService creates a new instance of FriendshipService
func NewFriendshipService(friendshipRepo repository.FriendshipRepository, userRepo repository.UserRepository) FriendshipService {
	return &friendshipService{
		friendshipRepo: friendshipRepo,
		userRepo:       userRepo,
	}
}

// SendFriendRequest sends a friend request to another user
func (s *friendshipService) SendFriendRequest(currentUserID uuid.UUID, friendEmail string) error {
	// Find the user by email
	friend, err := s.userRepo.FindByEmail(friendEmail)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.ErrUserNotFound
		}
		return err
	}

	// Check if user is trying to add themselves
	if friend.ID == currentUserID {
		return errors.New("cannot send friend request to yourself")
	}

	// Check if friendship already exists
	exists, err := s.friendshipRepo.CheckFriendshipExists(currentUserID, friend.ID)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("friendship or friend request already exists")
	}

	// Create friendship record
	userID1, userID2 := currentUserID, friend.ID
	if userID1.String() > userID2.String() {
		userID1, userID2 = userID2, userID1
	}

	friendship := &models.Friendship{
		UserID1:     userID1,
		UserID2:     userID2,
		Status:      models.FriendshipStatusPending,
		RequestedBy: currentUserID,
	}

	return s.friendshipRepo.Create(friendship)
}

// AcceptFriendRequest accepts a pending friend request
func (s *friendshipService) AcceptFriendRequest(currentUserID uuid.UUID, friendshipID uuid.UUID) error {
	friendship, err := s.friendshipRepo.FindByID(friendshipID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("friend request not found")
		}
		return err
	}

	// Verify the current user is the recipient (not the requester)
	if friendship.RequestedBy == currentUserID {
		return errors.New("cannot accept your own friend request")
	}

	// Verify the current user is one of the participants
	if friendship.UserID1 != currentUserID && friendship.UserID2 != currentUserID {
		return errors.New("you are not authorized to accept this friend request")
	}

	// Verify the friendship is still pending
	if friendship.Status != models.FriendshipStatusPending {
		return fmt.Errorf("friend request is already %s", friendship.Status)
	}

	// Update status to accepted
	friendship.Status = models.FriendshipStatusAccepted
	return s.friendshipRepo.Update(friendship)
}

// RejectFriendRequest rejects a pending friend request
func (s *friendshipService) RejectFriendRequest(currentUserID uuid.UUID, friendshipID uuid.UUID) error {
	friendship, err := s.friendshipRepo.FindByID(friendshipID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("friend request not found")
		}
		return err
	}

	// Verify the current user is the recipient (not the requester)
	if friendship.RequestedBy == currentUserID {
		return errors.New("cannot reject your own friend request")
	}

	// Verify the current user is one of the participants
	if friendship.UserID1 != currentUserID && friendship.UserID2 != currentUserID {
		return errors.New("you are not authorized to reject this friend request")
	}

	// Verify the friendship is still pending
	if friendship.Status != models.FriendshipStatusPending {
		return fmt.Errorf("friend request is already %s", friendship.Status)
	}

	// Update status to rejected
	friendship.Status = models.FriendshipStatusRejected
	return s.friendshipRepo.Update(friendship)
}

// BlockUser blocks a user (can be used on any friendship status)
func (s *friendshipService) BlockUser(currentUserID uuid.UUID, friendshipID uuid.UUID) error {
	friendship, err := s.friendshipRepo.FindByID(friendshipID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("friendship not found")
		}
		return err
	}

	// Verify the current user is one of the participants
	if friendship.UserID1 != currentUserID && friendship.UserID2 != currentUserID {
		return errors.New("you are not authorized to block this user")
	}

	// Update status to blocked
	friendship.Status = models.FriendshipStatusBlocked
	friendship.RequestedBy = currentUserID // The blocker becomes the requester
	return s.friendshipRepo.Update(friendship)
}

// RemoveFriend removes a friend (deletes the friendship)
func (s *friendshipService) RemoveFriend(currentUserID uuid.UUID, friendshipID uuid.UUID) error {
	friendship, err := s.friendshipRepo.FindByID(friendshipID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("friendship not found")
		}
		return err
	}

	// Verify the current user is one of the participants
	if friendship.UserID1 != currentUserID && friendship.UserID2 != currentUserID {
		return errors.New("you are not authorized to remove this friendship")
	}

	return s.friendshipRepo.Delete(friendshipID)
}

// GetFriends gets all accepted friends of a user
func (s *friendshipService) GetFriends(currentUserID uuid.UUID) (*dto.FriendListResponseDTO, error) {
	friendships, err := s.friendshipRepo.GetFriendsByUserID(currentUserID, models.FriendshipStatusAccepted)
	if err != nil {
		return nil, err
	}

	friends := make([]dto.FriendResponseDTO, 0, len(friendships))
	for _, friendship := range friendships {
		// Determine which user is the friend
		var friend *models.User
		if friendship.UserID1 == currentUserID {
			friend = friendship.User2
		} else {
			friend = friendship.User1
		}

		if friend == nil {
			continue
		}

		friends = append(friends, dto.FriendResponseDTO{
			ID:          friendship.ID,
			FriendID:    friend.ID,
			FriendEmail: friend.Email,
			FriendName:  friend.FirstName + " " + friend.LastName,
			Status:      string(friendship.Status),
			RequestedBy: friendship.RequestedBy,
			IsRequester: friendship.RequestedBy == currentUserID,
			CreatedAt:   friendship.CreatedAt,
			UpdatedAt:   friendship.UpdatedAt,
		})
	}

	return &dto.FriendListResponseDTO{
		Friends: friends,
		Total:   len(friends),
	}, nil
}

// GetPendingRequests gets all pending friend requests (sent and received)
func (s *friendshipService) GetPendingRequests(currentUserID uuid.UUID) (*dto.PendingRequestsResponseDTO, error) {
	// Get sent requests
	sentFriendships, err := s.friendshipRepo.GetPendingRequestsSent(currentUserID)
	if err != nil {
		return nil, err
	}

	// Get received requests
	receivedFriendships, err := s.friendshipRepo.GetPendingRequestsReceived(currentUserID)
	if err != nil {
		return nil, err
	}

	// Build sent requests response
	sent := make([]dto.PendingFriendRequestDTO, 0, len(sentFriendships))
	for _, friendship := range sentFriendships {
		// Determine which user is the recipient
		var recipient *models.User
		if friendship.UserID1 == currentUserID {
			recipient = friendship.User2
		} else {
			recipient = friendship.User1
		}

		if recipient == nil {
			continue
		}

		sent = append(sent, dto.PendingFriendRequestDTO{
			ID:           friendship.ID,
			RequesterID:  recipient.ID,
			Email:        recipient.Email,
			FirstName:    recipient.FirstName,
			LastName:     recipient.LastName,
			ProfileImage: &recipient.ProfileImageURL,
			CreatedAt:    friendship.CreatedAt,
		})
	}

	// Build received requests response
	received := make([]dto.PendingFriendRequestDTO, 0, len(receivedFriendships))
	for _, friendship := range receivedFriendships {
		requester := friendship.Requester
		if requester == nil {
			continue
		}

		received = append(received, dto.PendingFriendRequestDTO{
			ID:           friendship.ID,
			RequesterID:  requester.ID,
			Email:        requester.Email,
			FirstName:    requester.FirstName,
			LastName:     requester.LastName,
			ProfileImage: &requester.ProfileImageURL,
			CreatedAt:    friendship.CreatedAt,
		})
	}

	return &dto.PendingRequestsResponseDTO{
		Sent:     sent,
		Received: received,
	}, nil
}
