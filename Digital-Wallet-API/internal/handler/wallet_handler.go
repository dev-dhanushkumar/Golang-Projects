package handler

import (
	"digital-wallet-api/internal/service"
	"digital-wallet-api/pkg/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type WalletHandler struct {
	walletService service.WalletService
}

func NewWalletHandler(walletService service.WalletService) *WalletHandler {
	return &WalletHandler{
		walletService: walletService,
	}
}

// GetWallet godoc
// @Summary Get user wallet
// @Description Get authenticated user's wallet information
// @Tags wallet
// @Produce json
// @Security BearerAuth
// @Success 200 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Router /api/v1/wallet [get]
func (h *WalletHandler) GetWallet(c *gin.Context) {
	userID := c.GetString("user_id")
	uid, err := uuid.Parse(userID)
	if err != nil {
		utils.BadRequest(c, "Invalid user ID", err)
		return
	}

	wallet, err := h.walletService.GetWallet(uid)
	if err != nil {
		utils.NotFound(c, "Wallet not found", err)
		return
	}

	utils.OK(c, "Wallet retrieved successfully", wallet)
}

// GetBalance godoc
// @Summary Get wallet balance
// @Description Get authenticated user's wallet balance
// @Tags wallet
// @Produce json
// @Security BearerAuth
// @Success 200 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Router /api/v1/wallet/balance [get]
func (h *WalletHandler) GetBalance(c *gin.Context) {
	userID := c.GetString("user_id")
	uid, err := uuid.Parse(userID)
	if err != nil {
		utils.BadRequest(c, "Invalid user ID", err)
		return
	}

	balance, err := h.walletService.GetBalance(uid)
	if err != nil {
		utils.NotFound(c, "Wallet not found", err)
		return
	}

	utils.OK(c, "Balance retrieved successfully", balance)
}
