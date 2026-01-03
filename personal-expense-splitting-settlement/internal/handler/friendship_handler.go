package handler

import (
	"net/http"
	"personal-expense-splitting-settlement/internal/dto"
	"personal-expense-splitting-settlement/internal/services"
	"personal-expense-splitting-settlement/pkg/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// FriendshipHandler handles HTTP requests for friendship operations
type FriendshipHandler struct {
	friendshipService services.FriendshipService
}

// NewFriendshipHandler creates a new instance of FriendshipHandler
func NewFriendshipHandler(friendshipService services.FriendshipService) *FriendshipHandler {
	return &FriendshipHandler{
		friendshipService: friendshipService,
	}
}

// SendFriendRequest handles POST /api/v1/friends/request
func (h *FriendshipHandler) SendFriendRequest(c *gin.Context) {
	var req dto.FriendRequestDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	// Get current user ID from context (set by auth middleware)
	userID, exists := c.Get("user_id")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}

	currentUserID := userID.(uuid.UUID)

	err := h.friendshipService.SendFriendRequest(currentUserID, req.FriendEmail)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error(), err)
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "Friend request sent successfully", nil)
}

// AcceptFriendRequest handles POST /api/v1/friends/:id/accept
func (h *FriendshipHandler) AcceptFriendRequest(c *gin.Context) {
	friendshipIDStr := c.Param("id")
	friendshipID, err := uuid.Parse(friendshipIDStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid friendship ID", err)
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}

	currentUserID := userID.(uuid.UUID)

	err = h.friendshipService.AcceptFriendRequest(currentUserID, friendshipID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error(), err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Friend request accepted successfully", nil)
}

// RejectFriendRequest handles POST /api/v1/friends/:id/reject
func (h *FriendshipHandler) RejectFriendRequest(c *gin.Context) {
	friendshipIDStr := c.Param("id")
	friendshipID, err := uuid.Parse(friendshipIDStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid friendship ID", err)
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}

	currentUserID := userID.(uuid.UUID)

	err = h.friendshipService.RejectFriendRequest(currentUserID, friendshipID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error(), err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Friend request rejected successfully", nil)
}

// BlockUser handles POST /api/v1/friends/:id/block
func (h *FriendshipHandler) BlockUser(c *gin.Context) {
	friendshipIDStr := c.Param("id")
	friendshipID, err := uuid.Parse(friendshipIDStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid friendship ID", err)
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}

	currentUserID := userID.(uuid.UUID)

	err = h.friendshipService.BlockUser(currentUserID, friendshipID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error(), err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "User blocked successfully", nil)
}

// RemoveFriend handles DELETE /api/v1/friends/:id
func (h *FriendshipHandler) RemoveFriend(c *gin.Context) {
	friendshipIDStr := c.Param("id")
	friendshipID, err := uuid.Parse(friendshipIDStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid friendship ID", err)
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}

	currentUserID := userID.(uuid.UUID)

	err = h.friendshipService.RemoveFriend(currentUserID, friendshipID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error(), err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Friend removed successfully", nil)
}

// GetFriends handles GET /api/v1/friends
func (h *FriendshipHandler) GetFriends(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}

	currentUserID := userID.(uuid.UUID)

	friends, err := h.friendshipService.GetFriends(currentUserID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve friends", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Friends retrieved successfully", friends)
}

// GetPendingRequests handles GET /api/v1/friends/pending
func (h *FriendshipHandler) GetPendingRequests(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}

	currentUserID := userID.(uuid.UUID)

	pendingRequests, err := h.friendshipService.GetPendingRequests(currentUserID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve pending requests", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Pending requests retrieved successfully", pendingRequests)
}
