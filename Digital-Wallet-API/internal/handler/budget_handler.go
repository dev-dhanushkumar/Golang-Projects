package handler

import (
	"digital-wallet-api/internal/dto"
	"digital-wallet-api/internal/service"
	"digital-wallet-api/pkg/utils"
	"digital-wallet-api/pkg/validator"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type BudgetHandler struct {
	budgetService service.BudgetService
}

func NewBudgetHandler(budgetService service.BudgetService) *BudgetHandler {
	return &BudgetHandler{
		budgetService: budgetService,
	}
}

// CreateBudget godoc
// @Summary Create a new budget
// @Description Create a new budget for a category
// @Tags budgets
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.CreateBudgetRequest true "Create budget request"
// @Success 201 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Router /api/v1/budgets [post]
func (h *BudgetHandler) CreateBudget(c *gin.Context) {
	userID := c.GetString("user_id")
	uid, err := uuid.Parse(userID)
	if err != nil {
		utils.BadRequest(c, "Invalid user ID", err)
		return
	}

	var req dto.CreateBudgetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "Invalid request body", err)
		return
	}

	if err := validator.Validate(req); err != nil {
		utils.ValidationErrorResponse(c, err.Error())
		return
	}

	budget, err := h.budgetService.CreateBudget(uid, req)
	if err != nil {
		utils.BadRequest(c, "Failed to create budget", err)
		return
	}

	utils.Created(c, "Budget created successfully", budget)
}

// GetBudgets godoc
// @Summary Get all budgets
// @Description Get all budgets for authenticated user
// @Tags budgets
// @Produce json
// @Security BearerAuth
// @Success 200 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Router /api/v1/budgets [get]
func (h *BudgetHandler) GetBudgets(c *gin.Context) {
	userID := c.GetString("user_id")
	uid, err := uuid.Parse(userID)
	if err != nil {
		utils.BadRequest(c, "Invalid user ID", err)
		return
	}

	budgets, err := h.budgetService.GetBudgets(uid)
	if err != nil {
		utils.InternalServerError(c, "Failed to get budgets", err)
		return
	}

	utils.OK(c, "Budgets retrieved successfully", budgets)
}

// GetBudget godoc
// @Summary Get budget details
// @Description Get specific budget details
// @Tags budgets
// @Produce json
// @Security BearerAuth
// @Param id path string true "Budget ID"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Router /api/v1/budgets/{id} [get]
func (h *BudgetHandler) GetBudget(c *gin.Context) {
	userID := c.GetString("user_id")
	uid, err := uuid.Parse(userID)
	if err != nil {
		utils.BadRequest(c, "Invalid user ID", err)
		return
	}

	budgetID := c.Param("id")
	bid, err := uuid.Parse(budgetID)
	if err != nil {
		utils.BadRequest(c, "Invalid budget ID", err)
		return
	}

	budget, err := h.budgetService.GetBudget(uid, bid)
	if err != nil {
		utils.NotFound(c, "Budget not found", err)
		return
	}

	utils.OK(c, "Budget retrieved successfully", budget)
}

// UpdateBudget godoc
// @Summary Update budget
// @Description Update budget details
// @Tags budgets
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Budget ID"
// @Param request body dto.UpdateBudgetRequest true "Update budget request"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Router /api/v1/budgets/{id} [put]
func (h *BudgetHandler) UpdateBudget(c *gin.Context) {
	userID := c.GetString("user_id")
	uid, err := uuid.Parse(userID)
	if err != nil {
		utils.BadRequest(c, "Invalid user ID", err)
		return
	}

	budgetID := c.Param("id")
	bid, err := uuid.Parse(budgetID)
	if err != nil {
		utils.BadRequest(c, "Invalid budget ID", err)
		return
	}

	var req dto.UpdateBudgetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "Invalid request body", err)
		return
	}

	if err := validator.Validate(req); err != nil {
		utils.ValidationErrorResponse(c, err.Error())
		return
	}

	budget, err := h.budgetService.UpdateBudget(uid, bid, req)
	if err != nil {
		utils.BadRequest(c, "Failed to update budget", err)
		return
	}

	utils.OK(c, "Budget updated successfully", budget)
}

// DeleteBudget godoc
// @Summary Delete budget
// @Description Delete a budget
// @Tags budgets
// @Produce json
// @Security BearerAuth
// @Param id path string true "Budget ID"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Router /api/v1/budgets/{id} [delete]
func (h *BudgetHandler) DeleteBudget(c *gin.Context) {
	userID := c.GetString("user_id")
	uid, err := uuid.Parse(userID)
	if err != nil {
		utils.BadRequest(c, "Invalid user ID", err)
		return
	}

	budgetID := c.Param("id")
	bid, err := uuid.Parse(budgetID)
	if err != nil {
		utils.BadRequest(c, "Invalid budget ID", err)
		return
	}

	if err := h.budgetService.DeleteBudget(uid, bid); err != nil {
		utils.BadRequest(c, "Failed to delete budget", err)
		return
	}

	utils.OK(c, "Budget deleted successfully", nil)
}

// GetBudgetAlerts godoc
// @Summary Get budget alerts
// @Description Get alerts for budgets that are exceeded or near limit
// @Tags budgets
// @Produce json
// @Security BearerAuth
// @Success 200 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Router /api/v1/budgets/alerts [get]
func (h *BudgetHandler) GetBudgetAlerts(c *gin.Context) {
	userID := c.GetString("user_id")
	uid, err := uuid.Parse(userID)
	if err != nil {
		utils.BadRequest(c, "Invalid user ID", err)
		return
	}

	alerts, err := h.budgetService.GetBudgetAlerts(uid)
	if err != nil {
		utils.InternalServerError(c, "Failed to get budget alerts", err)
		return
	}

	utils.OK(c, "Budget alerts retrieved successfully", alerts)
}
