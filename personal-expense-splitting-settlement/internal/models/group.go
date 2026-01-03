package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// GroupType represents the type of group
type GroupType string

const (
	GroupTypeGeneral GroupType = "general"
	GroupTypeTrip    GroupType = "trip"
	GroupTypeHome    GroupType = "home"
	GroupTypeCouple  GroupType = "couple"
	GroupTypeEvent   GroupType = "event"
	GroupTypeProject GroupType = "project"
	GroupTypeOther   GroupType = "other"
)

// Group represents a group of users who share expenses
type Group struct {
	ID          uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Name        string         `gorm:"type:varchar(255);not null" json:"name"`
	Description string         `gorm:"type:text" json:"description"`
	Type        GroupType      `gorm:"type:varchar(50);not null;default:'general'" json:"type"`
	ImageURL    string         `gorm:"type:text" json:"image_url"`
	CreatedBy   uuid.UUID      `gorm:"type:uuid;not null;column:created_by" json:"created_by"`
	CreatedAt   time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// Relationships
	Creator User          `gorm:"foreignKey:CreatedBy" json:"creator,omitempty"`
	Members []GroupMember `gorm:"foreignKey:GroupID" json:"members,omitempty"`
}

// TableName specifies the table name for Group model
func (Group) TableName() string {
	return "groups"
}

// IsValidType checks if the group type is valid
func (g *Group) IsValidType() bool {
	validTypes := []GroupType{
		GroupTypeGeneral,
		GroupTypeTrip,
		GroupTypeHome,
		GroupTypeCouple,
		GroupTypeEvent,
		GroupTypeProject,
		GroupTypeOther,
	}
	for _, validType := range validTypes {
		if g.Type == validType {
			return true
		}
	}
	return false
}

// GetActiveMemberCount returns the count of active members
func (g *Group) GetActiveMemberCount() int {
	count := 0
	for _, member := range g.Members {
		if !member.HasLeft() {
			count++
		}
	}
	return count
}

// GetAdminCount returns the count of active admins
func (g *Group) GetAdminCount() int {
	count := 0
	for _, member := range g.Members {
		if member.IsAdmin() && !member.HasLeft() {
			count++
		}
	}
	return count
}

// HasMember checks if a user is an active member of the group
func (g *Group) HasMember(userID uuid.UUID) bool {
	for _, member := range g.Members {
		if member.UserID == userID && !member.HasLeft() {
			return true
		}
	}
	return false
}

// IsUserAdmin checks if a user is an admin of the group
func (g *Group) IsUserAdmin(userID uuid.UUID) bool {
	for _, member := range g.Members {
		if member.UserID == userID && member.IsAdmin() && !member.HasLeft() {
			return true
		}
	}
	return false
}
