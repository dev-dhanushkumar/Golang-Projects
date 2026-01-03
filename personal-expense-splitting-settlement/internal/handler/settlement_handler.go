package handler

import (
	"net/http"
	"strconv"

	"personal-expense-splitting-settlement/internal/dto"
	"personal-expense-splitting-settlement/internal/services"
	"personal-expense-splitting-settlement/pkg/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// SettlementHandler handles HTTP requests for settlement operations
type SettlementHandler struct {
	settlementService services.SettlementService
}

// NewSettlementHandler creates a new settlement handler instance
func NewSettlementHandler(settlementService services.SettlementService) *SettlementHandler {
	return &SettlementHandler{
		settlementService: settlementService,
	}
}

// CreateSettlement godoc
// @Summary      Create settlement
// @Description  Record a payment/settlement between users
// @Tags         settlements
// @Accept       json
// @Produce      json
// @Param        request body dto.CreateSettlementRequest true "Settlement details"
// @Success      201 {object} utils.Response{data=dto.SettlementResponse}
// @Failure      400 {object} utils.Response
// @Failure      401 {object} utils.Response
// @Failure      500 {object} utils.Response
// @Security     Bearer
// @Router       /settlements [post]
func (h *SettlementHandler) CreateSettlement(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated", nil)
		return
	}

	var req dto.CreateSettlementRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	settlement, err := h.settlementService.CreateSettlement(userID.(uuid.UUID), &req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Failed to create settlement", err)
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "Settlement created successfully", dto.ToSettlementResponse(settlement))
}

// GetSettlement godoc
// @Summary      Get settlement by ID
// @Description  Retrieve a specific settlement
// @Tags         settlements
// @Produce      json
// @Param        id path string true "Settlement ID"
// @Success      200 {object} utils.Response{data=dto.SettlementResponse}
// @Failure      400 {object} utils.Response
// @Failure      401 {object} utils.Response
// @Failure      404 {object} utils.Response
// @Security     Bearer
// @Router       /settlements/{id} [get]
func (h *SettlementHandler) GetSettlement(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated", nil)
		return
	}

	settlementID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid settlement ID", err)
		return
	}

	settlement, err := h.settlementService.GetSettlement(settlementID, userID.(uuid.UUID))
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Failed to retrieve settlement", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Settlement retrieved successfully", dto.ToSettlementResponse(settlement))
}

// GetUserSettlements godoc
// @Summary      Get user settlements
// @Description  Retrieve all settlements for the authenticated user
// @Tags         settlements
// @Produce      json
// @Param        limit query int false "Limit" default(20)
// @Param        offset query int false "Offset" default(0)
// @Success      200 {object} utils.Response{data=[]dto.SettlementResponse}
// @Failure      401 {object} utils.Response
// @Failure      500 {object} utils.Response
// @Security     Bearer
// @Router       /settlements [get]
func (h *SettlementHandler) GetUserSettlements(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated", nil)
		return
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	settlements, err := h.settlementService.GetUserSettlements(userID.(uuid.UUID), limit, offset)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve settlements", err)
		return
	}

	responses := make([]dto.SettlementResponse, len(settlements))
	for i, settlement := range settlements {
		responses[i] = dto.ToSettlementResponse(&settlement)
	}

	utils.SuccessResponse(c, http.StatusOK, "Settlements retrieved successfully", responses)
}

// GetSettlementsBetweenUsers godoc
// @Summary      Get settlements between users
// @Description  Retrieve settlements between two specific users
// @Tags         settlements
// @Produce      json
// @Param        user_id query string true "Other user ID"
// @Param        limit query int false "Limit" default(20)
// @Param        offset query int false "Offset" default(0)
// @Success      200 {object} utils.Response{data=[]dto.SettlementResponse}
// @Failure      400 {object} utils.Response
// @Failure      401 {object} utils.Response
// @Failure      500 {object} utils.Response
// @Security     Bearer
// @Router       /settlements/between [get]
func (h *SettlementHandler) GetSettlementsBetweenUsers(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated", nil)
		return
	}

	otherUserID, err := uuid.Parse(c.Query("other_user_id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid other_user_id parameter", err)
		return
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	settlements, err := h.settlementService.GetSettlementsBetweenUsers(userID.(uuid.UUID), otherUserID, limit, offset)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve settlements", err)
		return
	}

	responses := make([]dto.SettlementResponse, len(settlements))
	for i, settlement := range settlements {
		responses[i] = dto.ToSettlementResponse(&settlement)
	}

	utils.SuccessResponse(c, http.StatusOK, "Settlements retrieved successfully", responses)
}

// GetGroupSettlements godoc
// @Summary      Get group settlements
// @Description  Retrieve all settlements for a specific group
// @Tags         settlements
// @Produce      json
// @Param        id path string true "Group ID"
// @Param        limit query int false "Limit" default(20)
// @Param        offset query int false "Offset" default(0)
// @Success      200 {object} utils.Response{data=[]dto.SettlementResponse}
// @Failure      400 {object} utils.Response
// @Failure      401 {object} utils.Response
// @Failure      403 {object} utils.Response
// @Failure      500 {object} utils.Response
// @Security     Bearer
// @Router       /groups/{id}/settlements [get]
func (h *SettlementHandler) GetGroupSettlements(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated", nil)
		return
	}

	groupID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid group ID", err)
		return
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	settlements, err := h.settlementService.GetGroupSettlements(groupID, userID.(uuid.UUID), limit, offset)
	if err != nil {
		utils.ErrorResponse(c, http.StatusForbidden, "Failed to retrieve group settlements", err)
		return
	}

	responses := make([]dto.SettlementResponse, len(settlements))
	for i, settlement := range settlements {
		responses[i] = dto.ToSettlementResponse(&settlement)
	}

	utils.SuccessResponse(c, http.StatusOK, "Settlements retrieved successfully", responses)
}

// UpdateSettlement godoc
// @Summary      Update settlement
// @Description  Update settlement details (only payer can update before confirmation)
// @Tags         settlements
// @Accept       json
// @Produce      json
// @Param        id path string true "Settlement ID"
// @Param        request body dto.UpdateSettlementRequest true "Update details"
// @Success      200 {object} utils.Response{data=dto.SettlementResponse}
// @Failure      400 {object} utils.Response
// @Failure      401 {object} utils.Response
// @Failure      403 {object} utils.Response
// @Failure      404 {object} utils.Response
// @Security     Bearer
// @Router       /settlements/{id} [patch]
func (h *SettlementHandler) UpdateSettlement(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated", nil)
		return
	}

	settlementID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid settlement ID", err)
		return
	}

	var req dto.UpdateSettlementRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	settlement, err := h.settlementService.UpdateSettlement(settlementID, userID.(uuid.UUID), &req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusForbidden, "Failed to update settlement", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Settlement updated successfully", dto.ToSettlementResponse(settlement))
}

// ConfirmSettlement godoc
// @Summary      Confirm settlement
// @Description  Confirm receipt of payment (only payee can confirm)
// @Tags         settlements
// @Produce      json
// @Param        id path string true "Settlement ID"
// @Success      200 {object} utils.Response{data=dto.SettlementResponse}
// @Failure      400 {object} utils.Response
// @Failure      401 {object} utils.Response
// @Failure      403 {object} utils.Response
// @Failure      404 {object} utils.Response
// @Security     Bearer
// @Router       /settlements/{id}/confirm [patch]
func (h *SettlementHandler) ConfirmSettlement(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated", nil)
		return
	}

	settlementID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid settlement ID", err)
		return
	}

	settlement, err := h.settlementService.ConfirmSettlement(settlementID, userID.(uuid.UUID))
	if err != nil {
		utils.ErrorResponse(c, http.StatusForbidden, "Failed to confirm settlement", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Settlement confirmed successfully", dto.ToSettlementResponse(settlement))
}

// DeleteSettlement godoc
// @Summary      Delete settlement
// @Description  Delete a settlement (only payer can delete before confirmation)
// @Tags         settlements
// @Produce      json
// @Param        id path string true "Settlement ID"
// @Success      200 {object} utils.Response
// @Failure      400 {object} utils.Response
// @Failure      401 {object} utils.Response
// @Failure      403 {object} utils.Response
// @Failure      404 {object} utils.Response
// @Security     Bearer
// @Router       /settlements/{id} [delete]
func (h *SettlementHandler) DeleteSettlement(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated", nil)
		return
	}

	settlementID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid settlement ID", err)
		return
	}

	if err := h.settlementService.DeleteSettlement(settlementID, userID.(uuid.UUID)); err != nil {
		utils.ErrorResponse(c, http.StatusForbidden, "Failed to delete settlement", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Settlement deleted successfully", nil)
}
