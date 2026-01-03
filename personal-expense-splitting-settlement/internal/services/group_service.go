package services

import (
	"errors"
	"time"

	"personal-expense-splitting-settlement/internal/dto"
	"personal-expense-splitting-settlement/internal/models"
	"personal-expense-splitting-settlement/internal/repository"

	"github.com/google/uuid"
)

// GroupService interface defines business logic for group operations
type GroupService interface {
	CreateGroup(userID uuid.UUID, req *dto.CreateGroupRequest) (*models.Group, error)
	GetGroup(groupID uuid.UUID) (*models.Group, error)
	GetUserGroups(userID uuid.UUID) ([]models.Group, error)
	UpdateGroup(groupID, userID uuid.UUID, req *dto.UpdateGroupRequest) (*models.Group, error)
	DeleteGroup(groupID, userID uuid.UUID) error

	// Member operations
	AddMember(groupID, requestingUserID uuid.UUID, req *dto.AddMemberRequest) (*models.GroupMember, error)
	RemoveMember(groupID, requestingUserID, targetUserID uuid.UUID) error
	UpdateMemberRole(groupID, requestingUserID, targetUserID uuid.UUID, role models.MemberRole) error
	GetMembers(groupID uuid.UUID) ([]models.GroupMember, error)
}

type groupService struct {
	groupRepo      repository.GroupRepository
	userRepo       repository.UserRepository
	friendshipRepo repository.FriendshipRepository
}

// NewGroupService creates a new group service instance
func NewGroupService(
	groupRepo repository.GroupRepository,
	userRepo repository.UserRepository,
	friendshipRepo repository.FriendshipRepository,
) GroupService {
	return &groupService{
		groupRepo:      groupRepo,
		userRepo:       userRepo,
		friendshipRepo: friendshipRepo,
	}
}

// CreateGroup creates a new group with the creator as admin
func (s *groupService) CreateGroup(userID uuid.UUID, req *dto.CreateGroupRequest) (*models.Group, error) {
	// Validate group type
	group := &models.Group{
		Name:        req.Name,
		Description: req.Description,
		Type:        req.Type,
		ImageURL:    req.ImageURL,
		CreatedBy:   userID,
	}

	if !group.IsValidType() {
		return nil, errors.New("invalid group type")
	}

	// Create group
	if err := s.groupRepo.Create(group); err != nil {
		return nil, err
	}

	// Add creator as admin
	creatorMember := &models.GroupMember{
		GroupID:  group.ID,
		UserID:   userID,
		Role:     models.MemberRoleAdmin,
		JoinedAt: time.Now(),
	}

	if err := s.groupRepo.AddMember(creatorMember); err != nil {
		return nil, err
	}

	// Add other members if provided
	if len(req.MemberIDs) > 0 {
		for _, memberID := range req.MemberIDs {
			// Skip if memberID is the creator
			if memberID == userID {
				continue
			}

			// Verify user exists
			_, err := s.userRepo.FindByID(memberID)
			if err != nil {
				continue // Skip invalid users
			}

			// Add member with default role
			member := &models.GroupMember{
				GroupID:  group.ID,
				UserID:   memberID,
				Role:     models.MemberRoleMember,
				JoinedAt: time.Now(),
			}

			_ = s.groupRepo.AddMember(member) // Ignore errors for individual members
		}
	}

	// Reload group with members
	return s.groupRepo.FindByID(group.ID)
}

// GetGroup retrieves a group by ID
func (s *groupService) GetGroup(groupID uuid.UUID) (*models.Group, error) {
	return s.groupRepo.FindByID(groupID)
}

// GetUserGroups retrieves all groups for a user
func (s *groupService) GetUserGroups(userID uuid.UUID) ([]models.Group, error) {
	return s.groupRepo.FindByUserID(userID)
}

// UpdateGroup updates group information
func (s *groupService) UpdateGroup(groupID, userID uuid.UUID, req *dto.UpdateGroupRequest) (*models.Group, error) {
	// Check if user is admin
	isAdmin, err := s.groupRepo.IsUserAdmin(groupID, userID)
	if err != nil {
		return nil, err
	}

	if !isAdmin {
		return nil, errors.New("only admins can update group information")
	}

	// Get existing group
	group, err := s.groupRepo.FindByID(groupID)
	if err != nil {
		return nil, err
	}

	// Update fields if provided
	if req.Name != nil {
		group.Name = *req.Name
	}
	if req.Description != nil {
		group.Description = *req.Description
	}
	if req.Type != nil {
		group.Type = *req.Type
		if !group.IsValidType() {
			return nil, errors.New("invalid group type")
		}
	}
	if req.ImageURL != nil {
		group.ImageURL = *req.ImageURL
	}

	// Save updates
	if err := s.groupRepo.Update(group); err != nil {
		return nil, err
	}

	// Reload group with members
	return s.groupRepo.FindByID(groupID)
}

// DeleteGroup deletes a group (only by admin)
func (s *groupService) DeleteGroup(groupID, userID uuid.UUID) error {
	// Check if user is admin
	isAdmin, err := s.groupRepo.IsUserAdmin(groupID, userID)
	if err != nil {
		return err
	}

	if !isAdmin {
		return errors.New("only admins can delete groups")
	}

	return s.groupRepo.Delete(groupID)
}

// AddMember adds a member to a group
func (s *groupService) AddMember(groupID, requestingUserID uuid.UUID, req *dto.AddMemberRequest) (*models.GroupMember, error) {
	// Check if requesting user is admin
	isAdmin, err := s.groupRepo.IsUserAdmin(groupID, requestingUserID)
	if err != nil {
		return nil, err
	}

	if !isAdmin {
		return nil, errors.New("only admins can add members")
	}

	// Verify target user exists
	targetUser, err := s.userRepo.FindByID(req.UserID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	// Check if users are friends (optional validation - can be removed if not required)
	areFriends, err := s.friendshipRepo.AreFriends(requestingUserID, req.UserID)
	if err != nil {
		return nil, err
	}

	if !areFriends && requestingUserID != req.UserID {
		// Allow adding yourself or friends only
		return nil, errors.New("can only add friends to the group")
	}

	// Set default role if not provided
	role := req.Role
	if role == "" {
		role = models.MemberRoleMember
	}

	// Create member
	member := &models.GroupMember{
		GroupID:  groupID,
		UserID:   targetUser.ID,
		Role:     role,
		JoinedAt: time.Now(),
	}

	if err := s.groupRepo.AddMember(member); err != nil {
		return nil, err
	}

	// Reload member with user details
	return s.groupRepo.GetMember(groupID, req.UserID)
}

// RemoveMember removes a member from a group
func (s *groupService) RemoveMember(groupID, requestingUserID, targetUserID uuid.UUID) error {
	// Users can remove themselves
	if requestingUserID == targetUserID {
		return s.groupRepo.RemoveMember(groupID, targetUserID)
	}

	// Otherwise, must be admin
	isAdmin, err := s.groupRepo.IsUserAdmin(groupID, requestingUserID)
	if err != nil {
		return err
	}

	if !isAdmin {
		return errors.New("only admins can remove other members")
	}

	// Cannot remove if target is the last admin
	targetMember, err := s.groupRepo.GetMember(groupID, targetUserID)
	if err != nil {
		return err
	}

	if targetMember.IsAdmin() {
		// Count total admins
		members, err := s.groupRepo.GetMembers(groupID)
		if err != nil {
			return err
		}

		adminCount := 0
		for _, m := range members {
			if m.IsAdmin() {
				adminCount++
			}
		}

		if adminCount <= 1 {
			return errors.New("cannot remove the last admin from the group")
		}
	}

	return s.groupRepo.RemoveMember(groupID, targetUserID)
}

// UpdateMemberRole updates a member's role
func (s *groupService) UpdateMemberRole(groupID, requestingUserID, targetUserID uuid.UUID, role models.MemberRole) error {
	// Must be admin to change roles
	isAdmin, err := s.groupRepo.IsUserAdmin(groupID, requestingUserID)
	if err != nil {
		return err
	}

	if !isAdmin {
		return errors.New("only admins can change member roles")
	}

	// Cannot change your own role
	if requestingUserID == targetUserID {
		return errors.New("cannot change your own role")
	}

	// If demoting an admin, ensure there's at least one admin left
	targetMember, err := s.groupRepo.GetMember(groupID, targetUserID)
	if err != nil {
		return err
	}

	if targetMember.IsAdmin() && role == models.MemberRoleMember {
		// Count total admins
		members, err := s.groupRepo.GetMembers(groupID)
		if err != nil {
			return err
		}

		adminCount := 0
		for _, m := range members {
			if m.IsAdmin() {
				adminCount++
			}
		}

		if adminCount <= 1 {
			return errors.New("cannot demote the last admin")
		}
	}

	return s.groupRepo.UpdateMemberRole(groupID, targetUserID, role)
}

// GetMembers retrieves all members of a group
func (s *groupService) GetMembers(groupID uuid.UUID) ([]models.GroupMember, error) {
	return s.groupRepo.GetMembers(groupID)
}
