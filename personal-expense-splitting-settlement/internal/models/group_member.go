package models

import (
	"time"

	"github.com/google/uuid"
)

// MemberRole represents the role of a group member
type MemberRole string

const (
	MemberRoleAdmin  MemberRole = "admin"
	MemberRoleMember MemberRole = "member"
)

// GroupMember represents a user's membership in a group
type GroupMember struct {
	ID        uuid.UUID  `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	GroupID   uuid.UUID  `gorm:"type:uuid;not null;column:group_id" json:"group_id"`
	UserID    uuid.UUID  `gorm:"type:uuid;not null;column:user_id" json:"user_id"`
	Role      MemberRole `gorm:"type:varchar(50);not null;default:'member'" json:"role"`
	JoinedAt  time.Time  `gorm:"autoCreateTime;column:joined_at" json:"joined_at"`
	LeftAt    *time.Time `gorm:"column:left_at" json:"left_at,omitempty"`
	CreatedAt time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time  `gorm:"autoUpdateTime" json:"updated_at"`

	// Relationships
	Group Group `gorm:"foreignKey:GroupID" json:"group,omitempty"`
	User  User  `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// TableName specifies the table name for GroupMember model
func (GroupMember) TableName() string {
	return "group_members"
}

// IsAdmin checks if the member has admin role
func (gm *GroupMember) IsAdmin() bool {
	return gm.Role == MemberRoleAdmin
}

// IsMember checks if the member has member role
func (gm *GroupMember) IsMember() bool {
	return gm.Role == MemberRoleMember
}

// HasLeft checks if the member has left the group
func (gm *GroupMember) HasLeft() bool {
	return gm.LeftAt != nil
}

// IsActive checks if the member is still active in the group
func (gm *GroupMember) IsActive() bool {
	return !gm.HasLeft()
}

// IsValidRole checks if the role is valid
func (gm *GroupMember) IsValidRole() bool {
	return gm.Role == MemberRoleAdmin || gm.Role == MemberRoleMember
}

// Leave marks the member as having left the group
func (gm *GroupMember) Leave() {
	now := time.Now()
	gm.LeftAt = &now
}

// PromoteToAdmin promotes the member to admin role
func (gm *GroupMember) PromoteToAdmin() {
	gm.Role = MemberRoleAdmin
}

// DemoteToMember demotes the member to regular member role
func (gm *GroupMember) DemoteToMember() {
	gm.Role = MemberRoleMember
}
