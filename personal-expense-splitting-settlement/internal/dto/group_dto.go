package dto

import (
	"time"

	"personal-expense-splitting-settlement/internal/models"

	"github.com/google/uuid"
)

// CreateGroupRequest represents the request to create a new group
type CreateGroupRequest struct {
	Name        string           `json:"name" validate:"required,min=1,max=255"`
	Description string           `json:"description" validate:"omitempty,max=1000"`
	Type        models.GroupType `json:"type" validate:"required,oneof=general trip home couple event project other"`
	ImageURL    string           `json:"image_url" validate:"omitempty,url"`
	MemberIDs   []uuid.UUID      `json:"member_ids" validate:"omitempty,dive,uuid"`
}

// UpdateGroupRequest represents the request to update a group
type UpdateGroupRequest struct {
	Name        *string           `json:"name" validate:"omitempty,min=1,max=255"`
	Description *string           `json:"description" validate:"omitempty,max=1000"`
	Type        *models.GroupType `json:"type" validate:"omitempty,oneof=general trip home couple event project other"`
	ImageURL    *string           `json:"image_url" validate:"omitempty,url"`
}

// AddMemberRequest represents the request to add a member to a group
type AddMemberRequest struct {
	UserID uuid.UUID         `json:"user_id" validate:"required,uuid"`
	Role   models.MemberRole `json:"role" validate:"omitempty,oneof=admin member"`
}

// UpdateMemberRoleRequest represents the request to update a member's role
type UpdateMemberRoleRequest struct {
	Role models.MemberRole `json:"role" validate:"required,oneof=admin member"`
}

// GroupMemberResponse represents the response for a group member
type GroupMemberResponse struct {
	ID       uuid.UUID         `json:"id"`
	UserID   uuid.UUID         `json:"user_id"`
	Name     string            `json:"name"`
	Email    string            `json:"email"`
	Role     models.MemberRole `json:"role"`
	JoinedAt time.Time         `json:"joined_at"`
	LeftAt   *time.Time        `json:"left_at,omitempty"`
}

// GroupResponse represents the basic group response
type GroupResponse struct {
	ID          uuid.UUID        `json:"id"`
	Name        string           `json:"name"`
	Description string           `json:"description"`
	Type        models.GroupType `json:"type"`
	ImageURL    string           `json:"image_url"`
	CreatedBy   uuid.UUID        `json:"created_by"`
	MemberCount int              `json:"member_count"`
	CreatedAt   time.Time        `json:"created_at"`
	UpdatedAt   time.Time        `json:"updated_at"`
}

// GroupDetailResponse represents the detailed group response
type GroupDetailResponse struct {
	ID          uuid.UUID             `json:"id"`
	Name        string                `json:"name"`
	Description string                `json:"description"`
	Type        models.GroupType      `json:"type"`
	ImageURL    string                `json:"image_url"`
	CreatedBy   uuid.UUID             `json:"created_by"`
	Members     []GroupMemberResponse `json:"members"`
	CreatedAt   time.Time             `json:"created_at"`
	UpdatedAt   time.Time             `json:"updated_at"`
}

// GroupListResponse represents the response for a list of groups
type GroupListResponse struct {
	Groups []GroupResponse `json:"groups"`
	Total  int             `json:"total"`
}

// ToGroupResponse converts a Group model to GroupResponse
func ToGroupResponse(group *models.Group) GroupResponse {
	return GroupResponse{
		ID:          group.ID,
		Name:        group.Name,
		Description: group.Description,
		Type:        group.Type,
		ImageURL:    group.ImageURL,
		CreatedBy:   group.CreatedBy,
		MemberCount: group.GetActiveMemberCount(),
		CreatedAt:   group.CreatedAt,
		UpdatedAt:   group.UpdatedAt,
	}
}

// ToGroupDetailResponse converts a Group model to GroupDetailResponse
func ToGroupDetailResponse(group *models.Group) GroupDetailResponse {
	members := make([]GroupMemberResponse, 0)
	for _, member := range group.Members {
		if !member.HasLeft() {
			members = append(members, GroupMemberResponse{
				ID:       member.ID,
				UserID:   member.UserID,
				Name:     member.User.FirstName + " " + member.User.LastName,
				Email:    member.User.Email,
				Role:     member.Role,
				JoinedAt: member.JoinedAt,
				LeftAt:   member.LeftAt,
			})
		}
	}

	return GroupDetailResponse{
		ID:          group.ID,
		Name:        group.Name,
		Description: group.Description,
		Type:        group.Type,
		ImageURL:    group.ImageURL,
		CreatedBy:   group.CreatedBy,
		Members:     members,
		CreatedAt:   group.CreatedAt,
		UpdatedAt:   group.UpdatedAt,
	}
}

// ToGroupMemberResponse converts a GroupMember model to GroupMemberResponse
func ToGroupMemberResponse(member *models.GroupMember) GroupMemberResponse {
	return GroupMemberResponse{
		ID:       member.ID,
		UserID:   member.UserID,
		Name:     member.User.FirstName + " " + member.User.LastName,
		Email:    member.User.Email,
		Role:     member.Role,
		JoinedAt: member.JoinedAt,
		LeftAt:   member.LeftAt,
	}
}
