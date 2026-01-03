package handler

import (
	"net/http"

	"personal-expense-splitting-settlement/internal/services"
	"personal-expense-splitting-settlement/pkg/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// BalanceHandler handles HTTP requests for balance calculations
type BalanceHandler struct {
	balanceService services.BalanceService
}

// NewBalanceHandler creates a new balance handler instance
func NewBalanceHandler(balanceService services.BalanceService) *BalanceHandler {
	return &BalanceHandler{
		balanceService: balanceService,
	}
}

// GetBalanceSummary godoc
// @Summary      Get balance summary
// @Description  Get total balance summary for the authenticated user
// @Tags         balances
// @Produce      json
// @Success      200 {object} utils.Response{data=dto.BalanceSummaryResponse}
// @Failure      401 {object} utils.Response
// @Failure      500 {object} utils.Response
// @Security     Bearer
// @Router       /users/me/balance-summary [get]
func (h *BalanceHandler) GetBalanceSummary(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated", nil)
		return
	}

	summary, err := h.balanceService.GetBalanceSummary(userID.(uuid.UUID))
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve balance summary", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Balance summary retrieved successfully", summary)
}

// GetUserBalances godoc
// @Summary      Get detailed user balances
// @Description  Get detailed balances with each person for the authenticated user
// @Tags         balances
// @Produce      json
// @Success      200 {object} utils.Response{data=dto.BalanceSummaryResponse}
// @Failure      401 {object} utils.Response
// @Failure      500 {object} utils.Response
// @Security     Bearer
// @Router       /users/me/balances [get]
func (h *BalanceHandler) GetUserBalances(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated", nil)
		return
	}

	balances, err := h.balanceService.GetUserBalances(userID.(uuid.UUID))
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve balances", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Balances retrieved successfully", balances)
}

// GetGroupBalances godoc
// @Summary      Get group balances
// @Description  Get balance breakdown for all members in a group
// @Tags         balances
// @Produce      json
// @Param        id path string true "Group ID"
// @Success      200 {object} utils.Response{data=dto.GroupBalancesResponse}
// @Failure      400 {object} utils.Response
// @Failure      401 {object} utils.Response
// @Failure      403 {object} utils.Response
// @Failure      500 {object} utils.Response
// @Security     Bearer
// @Router       /groups/{id}/balances [get]
func (h *BalanceHandler) GetGroupBalances(c *gin.Context) {
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

	balances, err := h.balanceService.GetGroupBalances(groupID, userID.(uuid.UUID))
	if err != nil {
		utils.ErrorResponse(c, http.StatusForbidden, "Failed to retrieve group balances", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Group balances retrieved successfully", balances)
}

// GetSettlementSuggestions godoc
// @Summary      Get settlement suggestions
// @Description  Get smart settlement suggestions to minimize transactions
// @Tags         balances
// @Produce      json
// @Success      200 {object} utils.Response{data=dto.SettlementSuggestionsResponse}
// @Failure      401 {object} utils.Response
// @Failure      500 {object} utils.Response
// @Security     Bearer
// @Router       /settlements/suggestions [get]
func (h *BalanceHandler) GetSettlementSuggestions(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated", nil)
		return
	}

	suggestions, err := h.balanceService.GetSettlementSuggestions(userID.(uuid.UUID))
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to generate settlement suggestions", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Settlement suggestions generated successfully", suggestions)
}

// GetGroupSettlementSuggestions godoc
// @Summary      Get group settlement suggestions
// @Description  Get smart settlement suggestions for a group to minimize transactions
// @Tags         balances
// @Produce      json
// @Param        id path string true "Group ID"
// @Success      200 {object} utils.Response{data=dto.SettlementSuggestionsResponse}
// @Failure      400 {object} utils.Response
// @Failure      401 {object} utils.Response
// @Failure      403 {object} utils.Response
// @Failure      500 {object} utils.Response
// @Security     Bearer
// @Router       /groups/{id}/settlement-suggestions [get]
func (h *BalanceHandler) GetGroupSettlementSuggestions(c *gin.Context) {
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

	suggestions, err := h.balanceService.GetGroupSettlementSuggestions(groupID, userID.(uuid.UUID))
	if err != nil {
		utils.ErrorResponse(c, http.StatusForbidden, "Failed to generate group settlement suggestions", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Group settlement suggestions generated successfully", suggestions)
}
