package router

import (
	"digital-wallet-api/internal/handler"
	"digital-wallet-api/internal/middleware"

	"github.com/gin-gonic/gin"
)

type RouterConfig struct {
	AuthHandler        *handler.AuthHandler
	WalletHandler      *handler.WalletHandler
	TransactionHandler *handler.TransactionHandler
	BudgetHandler      *handler.BudgetHandler
	JWTSecret          string
}

func SetupRouter(config RouterConfig) *gin.Engine {
	router := gin.New()

	// Global middleware
	router.Use(gin.Recovery())
	router.Use(middleware.ErrorMiddleware())
	router.Use(middleware.LoggerMiddleware())
	router.Use(middleware.CORSMiddleware())

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"service": "digital-wallet-api",
		})
	})

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// Auth routes (public)
		auth := v1.Group("/auth")
		{
			auth.POST("/register", config.AuthHandler.Register)
			auth.POST("/login", config.AuthHandler.Login)

			// Protected auth routes
			authProtected := auth.Group("")
			authProtected.Use(middleware.AuthMiddleware(config.JWTSecret))
			{
				authProtected.GET("/profile", config.AuthHandler.GetProfile)
				authProtected.PUT("/profile", config.AuthHandler.UpdateProfile)
			}
		}

		// Protected routes
		protected := v1.Group("")
		protected.Use(middleware.AuthMiddleware(config.JWTSecret))
		{
			// Wallet routes
			wallet := protected.Group("/wallet")
			{
				wallet.GET("", config.WalletHandler.GetWallet)
				wallet.GET("/balance", config.WalletHandler.GetBalance)
			}

			// Transaction routes
			transactions := protected.Group("/transactions")
			{
				transactions.POST("/credit", config.TransactionHandler.Credit)
				transactions.POST("/debit", config.TransactionHandler.Debit)
				transactions.POST("/transfer", config.TransactionHandler.Transfer)
				transactions.GET("", config.TransactionHandler.GetTransactions)
				transactions.GET("/:id", config.TransactionHandler.GetTransaction)
				transactions.GET("/summary", config.TransactionHandler.GetSummary)
			}

			// Budget routes
			budgets := protected.Group("/budgets")
			{
				budgets.POST("", config.BudgetHandler.CreateBudget)
				budgets.GET("", config.BudgetHandler.GetBudgets)
				budgets.GET("/alerts", config.BudgetHandler.GetBudgetAlerts)
				budgets.GET("/:id", config.BudgetHandler.GetBudget)
				budgets.PUT("/:id", config.BudgetHandler.UpdateBudget)
				budgets.DELETE("/:id", config.BudgetHandler.DeleteBudget)
			}
		}
	}

	return router
}
