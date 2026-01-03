package handler

import (
	"net/http"

	"personal-expense-splitting-settlement/internal/dto"
	"personal-expense-splitting-settlement/internal/services"
	"personal-expense-splitting-settlement/pkg/utils"
	"personal-expense-splitting-settlement/pkg/validator"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// GroupHandler handles HTTP requests for group operations
type GroupHandler struct {
	groupService services.GroupService
}

// NewGroupHandler creates a new group handler instance
func NewGroupHandler(groupService services.GroupService) *GroupHandler {
	return &GroupHandler{
		groupService: groupService,
	}
}

// CreateGroup handles group creation
// POST /api/v1/groups
func (h *GroupHandler) CreateGroup(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated", nil)
		return
	}

	var req dto.CreateGroupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request body", nil)
		return
	}

	if err := validator.Validate(req); err != nil {
		utils.ValidationErrorResponse(c, err.Error())
		return
	}

	group, err := h.groupService.CreateGroup(userID.(uuid.UUID), &req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error(), err)
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "Group created successfully", dto.ToGroupDetailResponse(group))
}

// GetUserGroups retrieves all groups for the authenticated user
// GET /api/v1/groups
func (h *GroupHandler) GetUserGroups(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated", nil)
		return
	}

	groups, err := h.groupService.GetUserGroups(userID.(uuid.UUID))
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error(), err)
		return
	}

	groupResponses := make([]dto.GroupResponse, len(groups))
	for i, group := range groups {
		groupResponses[i] = dto.ToGroupResponse(&group)
	}

	utils.SuccessResponse(c, http.StatusOK, "Groups retrieved successfully", dto.GroupListResponse{
		Groups: groupResponses,
		Total:  len(groupResponses),
	})
}

// GetGroup retrieves a specific group by ID
// GET /api/v1/groups/:id
func (h *GroupHandler) GetGroup(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated", nil)
		return
	}

	groupID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid group ID", nil)
		return
	}

	group, err := h.groupService.GetGroup(groupID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, err.Error(), err)
		return
	}

	// Verify user is a member of the group
	if !group.HasMember(userID.(uuid.UUID)) {
		utils.ErrorResponse(c, http.StatusForbidden, "You are not a member of this group", nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Group retrieved successfully", dto.ToGroupDetailResponse(group))
}

// UpdateGroup updates group information
// PATCH /api/v1/groups/:id
func (h *GroupHandler) UpdateGroup(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated", nil)
		return
	}

	groupID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid group ID", nil)
		return
	}

	var req dto.UpdateGroupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request body", nil)
		return
	}

	if err := validator.Validate(req); err != nil {
		utils.ValidationErrorResponse(c, err.Error())
		return
	}

	group, err := h.groupService.UpdateGroup(groupID, userID.(uuid.UUID), &req)
	if err != nil {
		if err.Error() == "only admins can update group information" {
			utils.ErrorResponse(c, http.StatusForbidden, err.Error(), err)
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error(), err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Group updated successfully", dto.ToGroupDetailResponse(group))
}

// DeleteGroup deletes a group
// DELETE /api/v1/groups/:id
func (h *GroupHandler) DeleteGroup(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated", nil)
		return
	}

	groupID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid group ID", nil)
		return
	}

	err = h.groupService.DeleteGroup(groupID, userID.(uuid.UUID))
	if err != nil {
		if err.Error() == "only admins can delete groups" {
			utils.ErrorResponse(c, http.StatusForbidden, err.Error(), err)
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error(), err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Group deleted successfully", nil)
}

// AddMember adds a member to a group
// POST /api/v1/groups/:id/members
func (h *GroupHandler) AddMember(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated", nil)
		return
	}

	groupID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid group ID", nil)
		return
	}

	var req dto.AddMemberRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request body", nil)
		return
	}

	if err := validator.Validate(req); err != nil {
		utils.ValidationErrorResponse(c, err.Error())
		return
	}

	member, err := h.groupService.AddMember(groupID, userID.(uuid.UUID), &req)
	if err != nil {
		if err.Error() == "only admins can add members" {
			utils.ErrorResponse(c, http.StatusForbidden, err.Error(), err)
			return
		}
		if err.Error() == "can only add friends to the group" {
			utils.ErrorResponse(c, http.StatusBadRequest, err.Error(), err)
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error(), err)
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "Member added successfully", dto.ToGroupMemberResponse(member))
}

// RemoveMember removes a member from a group
// DELETE /api/v1/groups/:id/members/:user_id
func (h *GroupHandler) RemoveMember(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated", nil)
		return
	}

	groupID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid group ID", nil)
		return
	}

	targetUserID, err := uuid.Parse(c.Param("user_id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid user ID", nil)
		return
	}

	err = h.groupService.RemoveMember(groupID, userID.(uuid.UUID), targetUserID)
	if err != nil {
		if err.Error() == "only admins can remove other members" {
			utils.ErrorResponse(c, http.StatusForbidden, err.Error(), err)
			return
		}
		if err.Error() == "cannot remove the last admin from the group" {
			utils.ErrorResponse(c, http.StatusBadRequest, err.Error(), err)
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error(), err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Member removed successfully", nil)
}

// UpdateMemberRole updates a member's role in a group
// PATCH /api/v1/groups/:id/members/:user_id
func (h *GroupHandler) UpdateMemberRole(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated", nil)
		return
	}

	groupID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid group ID", nil)
		return
	}

	targetUserID, err := uuid.Parse(c.Param("user_id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid user ID", nil)
		return
	}

	var req dto.UpdateMemberRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request body", nil)
		return
	}

	if err := validator.Validate(req); err != nil {
		utils.ValidationErrorResponse(c, err.Error())
		return
	}

	err = h.groupService.UpdateMemberRole(groupID, userID.(uuid.UUID), targetUserID, req.Role)
	if err != nil {
		if err.Error() == "only admins can change member roles" || err.Error() == "cannot change your own role" {
			utils.ErrorResponse(c, http.StatusForbidden, err.Error(), err)
			return
		}
		if err.Error() == "cannot demote the last admin" {
			utils.ErrorResponse(c, http.StatusBadRequest, err.Error(), err)
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error(), err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Member role updated successfully", nil)
}
