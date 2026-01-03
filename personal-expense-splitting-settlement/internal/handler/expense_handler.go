package handler

import (
	"net/http"

	"personal-expense-splitting-settlement/internal/dto"
	"personal-expense-splitting-settlement/internal/services"
	"personal-expense-splitting-settlement/pkg/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// ExpenseHandler handles expense-related HTTP requests
type ExpenseHandler struct {
	expenseService services.ExpenseService
}

// NewExpenseHandler creates a new expense handler instance
func NewExpenseHandler(expenseService services.ExpenseService) *ExpenseHandler {
	return &ExpenseHandler{
		expenseService: expenseService,
	}
}

// CreateExpense creates a new expense
// @Summary Create a new expense
// @Description Create a new expense with participants and split method
// @Tags expenses
// @Accept json
// @Produce json
// @Param expense body dto.CreateExpenseRequest true "Expense details"
// @Success 201 {object} dto.ExpenseDetailResponse
// @Failure 400 {object} utils.ErrorResponse
// @Failure 401 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /expenses [post]
func (h *ExpenseHandler) CreateExpense(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated", nil)
		return
	}

	var req dto.CreateExpenseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	expense, err := h.expenseService.CreateExpense(userID.(uuid.UUID), &req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error(), err)
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "Expense created successfully", dto.ToExpenseDetailResponse(expense))
}

// GetExpense retrieves a specific expense by ID
// @Summary Get expense details
// @Description Get detailed information about a specific expense
// @Tags expenses
// @Produce json
// @Param id path string true "Expense ID"
// @Success 200 {object} dto.ExpenseDetailResponse
// @Failure 400 {object} utils.ErrorResponse
// @Failure 401 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Router /expenses/{id} [get]
func (h *ExpenseHandler) GetExpense(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated", nil)
		return
	}

	expenseIDStr := c.Param("id")
	expenseID, err := uuid.Parse(expenseIDStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid expense ID", err)
		return
	}

	expense, err := h.expenseService.GetExpense(expenseID, userID.(uuid.UUID))
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, err.Error(), err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Expense retrieved successfully", dto.ToExpenseDetailResponse(expense))
}

// GetUserExpenses retrieves all expenses for the authenticated user
// @Summary Get user expenses
// @Description Get all expenses where the user is a participant
// @Tags expenses
// @Produce json
// @Param limit query int false "Limit" default(20)
// @Param offset query int false "Offset" default(0)
// @Success 200 {array} dto.ExpenseResponse
// @Failure 400 {object} utils.ErrorResponse
// @Failure 401 {object} utils.ErrorResponse
// @Router /expenses [get]
func (h *ExpenseHandler) GetUserExpenses(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated", nil)
		return
	}

	var req struct {
		Limit  int `form:"limit"`
		Offset int `form:"offset"`
	}

	if err := c.ShouldBindQuery(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid query parameters", err)
		return
	}

	expenses, err := h.expenseService.GetUserExpenses(userID.(uuid.UUID), req.Limit, req.Offset)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve expenses", err)
		return
	}

	response := make([]dto.ExpenseResponse, len(expenses))
	for i, expense := range expenses {
		response[i] = dto.ToExpenseResponse(&expense)
	}

	utils.SuccessResponse(c, http.StatusOK, "Expenses retrieved successfully", response)
}

// GetGroupExpenses retrieves all expenses for a group
// @Summary Get group expenses
// @Description Get all expenses for a specific group
// @Tags expenses
// @Produce json
// @Param id path string true "Group ID"
// @Param limit query int false "Limit" default(20)
// @Param offset query int false "Offset" default(0)
// @Success 200 {array} dto.ExpenseResponse
// @Failure 400 {object} utils.ErrorResponse
// @Failure 401 {object} utils.ErrorResponse
// @Failure 403 {object} utils.ErrorResponse
// @Router /groups/{id}/expenses [get]
func (h *ExpenseHandler) GetGroupExpenses(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated", nil)
		return
	}

	groupIDStr := c.Param("id")
	groupID, err := uuid.Parse(groupIDStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid group ID", err)
		return
	}

	var req struct {
		Limit  int `form:"limit"`
		Offset int `form:"offset"`
	}

	if err := c.ShouldBindQuery(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid query parameters", err)
		return
	}

	expenses, err := h.expenseService.GetGroupExpenses(groupID, userID.(uuid.UUID), req.Limit, req.Offset)
	if err != nil {
		utils.ErrorResponse(c, http.StatusForbidden, err.Error(), err)
		return
	}

	response := make([]dto.ExpenseResponse, len(expenses))
	for i, expense := range expenses {
		response[i] = dto.ToExpenseResponse(&expense)
	}

	utils.SuccessResponse(c, http.StatusOK, "Expenses retrieved successfully", response)
}

// GetExpensesWithFilters retrieves expenses with filters
// @Summary Get expenses with filters
// @Description Get expenses filtered by group, category, and date range
// @Tags expenses
// @Produce json
// @Param group_id query string false "Group ID"
// @Param category query string false "Category"
// @Param start_date query string false "Start date (YYYY-MM-DD)"
// @Param end_date query string false "End date (YYYY-MM-DD)"
// @Param limit query int false "Limit" default(20)
// @Param offset query int false "Offset" default(0)
// @Success 200 {array} dto.ExpenseResponse
// @Failure 400 {object} utils.ErrorResponse
// @Failure 401 {object} utils.ErrorResponse
// @Router /expenses/filter [get]
func (h *ExpenseHandler) GetExpensesWithFilters(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated", nil)
		return
	}

	var req dto.ExpenseListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid query parameters", err)
		return
	}

	expenses, err := h.expenseService.GetExpensesWithFilters(userID.(uuid.UUID), &req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve expenses", err)
		return
	}

	response := make([]dto.ExpenseResponse, len(expenses))
	for i, expense := range expenses {
		response[i] = dto.ToExpenseResponse(&expense)
	}

	utils.SuccessResponse(c, http.StatusOK, "Expenses retrieved successfully", response)
}

// UpdateExpense updates an expense
// @Summary Update expense
// @Description Update expense details (only creator can update)
// @Tags expenses
// @Accept json
// @Produce json
// @Param id path string true "Expense ID"
// @Param expense body dto.UpdateExpenseRequest true "Update details"
// @Success 200 {object} dto.ExpenseDetailResponse
// @Failure 400 {object} utils.ErrorResponse
// @Failure 401 {object} utils.ErrorResponse
// @Failure 403 {object} utils.ErrorResponse
// @Router /expenses/{id} [patch]
func (h *ExpenseHandler) UpdateExpense(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated", nil)
		return
	}

	expenseIDStr := c.Param("id")
	expenseID, err := uuid.Parse(expenseIDStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid expense ID", err)
		return
	}

	var req dto.UpdateExpenseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	expense, err := h.expenseService.UpdateExpense(expenseID, userID.(uuid.UUID), &req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusForbidden, err.Error(), err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Expense updated successfully", dto.ToExpenseDetailResponse(expense))
}

// DeleteExpense deletes an expense
// @Summary Delete expense
// @Description Delete an expense (only creator can delete, cannot delete if settlements made)
// @Tags expenses
// @Produce json
// @Param id path string true "Expense ID"
// @Success 200 {object} utils.SuccessResponse
// @Failure 400 {object} utils.ErrorResponse
// @Failure 401 {object} utils.ErrorResponse
// @Failure 403 {object} utils.ErrorResponse
// @Router /expenses/{id} [delete]
func (h *ExpenseHandler) DeleteExpense(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated", nil)
		return
	}

	expenseIDStr := c.Param("id")
	expenseID, err := uuid.Parse(expenseIDStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid expense ID", err)
		return
	}

	if err := h.expenseService.DeleteExpense(expenseID, userID.(uuid.UUID)); err != nil {
		utils.ErrorResponse(c, http.StatusForbidden, err.Error(), err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Expense deleted successfully", nil)
}
