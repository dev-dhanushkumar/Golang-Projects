package handler

import (
	"digital-wallet-api/internal/dto"
	"digital-wallet-api/internal/service"
	"digital-wallet-api/pkg/utils"
	"digital-wallet-api/pkg/validator"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type TransactionHandler struct {
	transactionService service.TransactionService
}

func NewTransactionHandler(transactionService service.TransactionService) *TransactionHandler {
	return &TransactionHandler{
		transactionService: transactionService,
	}
}

// Credit godoc
// @Summary Add money to wallet
// @Description Credit money to user's wallet
// @Tags transactions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.CreditRequest true "Credit request"
// @Success 201 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Router /api/v1/transactions/credit [post]
func (h *TransactionHandler) Credit(c *gin.Context) {
	userID := c.GetString("user_id")
	uid, err := uuid.Parse(userID)
	if err != nil {
		utils.BadRequest(c, "Invalid user ID", err)
		return
	}

	var req dto.CreditRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "Invalid request body", err)
		return
	}

	if err := validator.Validate(req); err != nil {
		utils.ValidationErrorResponse(c, err.Error())
		return
	}

	transaction, err := h.transactionService.Credit(uid, req)
	if err != nil {
		utils.BadRequest(c, "Credit transaction failed", err)
		return
	}

	utils.Created(c, "Money credited successfully", transaction)
}

// Debit godoc
// @Summary Spend money from wallet
// @Description Debit money from user's wallet
// @Tags transactions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.DebitRequest true "Debit request"
// @Success 201 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Router /api/v1/transactions/debit [post]
func (h *TransactionHandler) Debit(c *gin.Context) {
	userID := c.GetString("user_id")
	uid, err := uuid.Parse(userID)
	if err != nil {
		utils.BadRequest(c, "Invalid user ID", err)
		return
	}

	var req dto.DebitRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "Invalid request body", err)
		return
	}

	if err := validator.Validate(req); err != nil {
		utils.ValidationErrorResponse(c, err.Error())
		return
	}

	transaction, err := h.transactionService.Debit(uid, req)
	if err != nil {
		utils.BadRequest(c, "Debit transaction failed", err)
		return
	}

	utils.Created(c, "Money debited successfully", transaction)
}

// Transfer godoc
// @Summary Transfer money to another user
// @Description Transfer money from user's wallet to another user
// @Tags transactions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.TransferRequest true "Transfer request"
// @Success 201 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Router /api/v1/transactions/transfer [post]
func (h *TransactionHandler) Transfer(c *gin.Context) {
	userID := c.GetString("user_id")
	uid, err := uuid.Parse(userID)
	if err != nil {
		utils.BadRequest(c, "Invalid user ID", err)
		return
	}

	var req dto.TransferRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "Invalid request body", err)
		return
	}

	if err := validator.Validate(req); err != nil {
		utils.ValidationErrorResponse(c, err.Error())
		return
	}

	transfer, err := h.transactionService.Transfer(uid, req)
	if err != nil {
		utils.BadRequest(c, "Transfer failed", err)
		return
	}

	utils.Created(c, "Transfer completed successfully", transfer)
}

// GetTransactions godoc
// @Summary Get transaction history
// @Description Get paginated transaction history
// @Tags transactions
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(20)
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Router /api/v1/transactions [get]
func (h *TransactionHandler) GetTransactions(c *gin.Context) {
	userID := c.GetString("user_id")
	uid, err := uuid.Parse(userID)
	if err != nil {
		utils.BadRequest(c, "Invalid user ID", err)
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	transactions, err := h.transactionService.GetTransactions(uid, page, pageSize)
	if err != nil {
		utils.InternalServerError(c, "Failed to get transactions", err)
		return
	}

	utils.OK(c, "Transactions retrieved successfully", transactions)
}

// GetTransaction godoc
// @Summary Get transaction details
// @Description Get specific transaction details
// @Tags transactions
// @Produce json
// @Security BearerAuth
// @Param id path string true "Transaction ID"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Router /api/v1/transactions/{id} [get]
func (h *TransactionHandler) GetTransaction(c *gin.Context) {
	userID := c.GetString("user_id")
	uid, err := uuid.Parse(userID)
	if err != nil {
		utils.BadRequest(c, "Invalid user ID", err)
		return
	}

	transactionID := c.Param("id")
	txID, err := uuid.Parse(transactionID)
	if err != nil {
		utils.BadRequest(c, "Invalid transaction ID", err)
		return
	}

	transaction, err := h.transactionService.GetTransaction(uid, txID)
	if err != nil {
		utils.NotFound(c, "Transaction not found", err)
		return
	}

	utils.OK(c, "Transaction retrieved successfully", transaction)
}

// GetSummary godoc
// @Summary Get transaction summary
// @Description Get transaction summary for a date range
// @Tags transactions
// @Produce json
// @Security BearerAuth
// @Param start_date query string false "Start date (YYYY-MM-DD)"
// @Param end_date query string false "End date (YYYY-MM-DD)"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Router /api/v1/transactions/summary [get]
func (h *TransactionHandler) GetSummary(c *gin.Context) {
	userID := c.GetString("user_id")
	uid, err := uuid.Parse(userID)
	if err != nil {
		utils.BadRequest(c, "Invalid user ID", err)
		return
	}

	// Parse dates
	startDateStr := c.DefaultQuery("start_date", time.Now().AddDate(0, -1, 0).Format("2006-01-02"))
	endDateStr := c.DefaultQuery("end_date", time.Now().Format("2006-01-02"))

	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		utils.BadRequest(c, "Invalid start date format", err)
		return
	}

	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		utils.BadRequest(c, "Invalid end date format", err)
		return
	}

	// Add end of day to end date
	endDate = endDate.Add(23*time.Hour + 59*time.Minute + 59*time.Second)

	summary, err := h.transactionService.GetSummary(uid, startDate, endDate)
	if err != nil {
		utils.InternalServerError(c, "Failed to get summary", err)
		return
	}

	utils.OK(c, "Summary retrieved successfully", summary)
}
