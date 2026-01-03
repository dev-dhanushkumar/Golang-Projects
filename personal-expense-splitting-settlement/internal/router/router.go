package router

import (
	"personal-expense-splitting-settlement/internal/handler"
	"personal-expense-splitting-settlement/internal/middleware"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type RouterConfig struct {
	AuthHandler       *handler.AuthHandler
	FriendshipHandler *handler.FriendshipHandler
	GroupHandler      *handler.GroupHandler
	ExpenseHandler    *handler.ExpenseHandler
	SettlementHandler *handler.SettlementHandler
	BalanceHandler    *handler.BalanceHandler
	JWTSecret         string
	Logger            *zap.SugaredLogger
}

func SetupRouter(config RouterConfig) *gin.Engine {
	router := gin.New()

	// Global Middleware
	router.Use(gin.Recovery())
	router.Use(middleware.LoggerMiddleware(config.Logger))

	// Health Check
	router.GET("/health", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"status":  "ok",
			"service": "Personal expense splitting and settlement",
		})
	})

	v1 := router.Group("/api/v1")
	{
		auth := v1.Group("/auth")
		{
			auth.POST("/register", config.AuthHandler.Register)
			auth.POST("/login", config.AuthHandler.Login)
			auth.POST("/refresh", config.AuthHandler.Refresh)

			// Protected routes
			protected := auth.Group("/")
			protected.Use(middleware.AuthMiddleware(config.JWTSecret))
			{
				protected.POST("/logout", config.AuthHandler.Logout)
				protected.GET("/me", config.AuthHandler.GetMe)
				protected.GET("/sessions", config.AuthHandler.GetSessions)
			}
		}

		// User endpoints
		users := v1.Group("/users")
		users.Use(middleware.AuthMiddleware(config.JWTSecret))
		{
			users.GET("/me", config.AuthHandler.GetMe)
			users.PATCH("/me", config.AuthHandler.UpdateProfile)
			users.GET("/me/balance-summary", config.BalanceHandler.GetBalanceSummary)
			users.GET("/me/balances", config.BalanceHandler.GetUserBalances)
		}

		// Friendship endpoints
		friends := v1.Group("/friends")
		friends.Use(middleware.AuthMiddleware(config.JWTSecret))
		{
			friends.POST("/request", config.FriendshipHandler.SendFriendRequest)
			friends.POST("/:id/accept", config.FriendshipHandler.AcceptFriendRequest)
			friends.POST("/:id/reject", config.FriendshipHandler.RejectFriendRequest)
			friends.POST("/:id/block", config.FriendshipHandler.BlockUser)
			friends.DELETE("/:id", config.FriendshipHandler.RemoveFriend)
			friends.GET("", config.FriendshipHandler.GetFriends)
			friends.GET("/pending", config.FriendshipHandler.GetPendingRequests)
		}

		// Group endpoints
		groups := v1.Group("/groups")
		groups.Use(middleware.AuthMiddleware(config.JWTSecret))
		{
			groups.POST("", config.GroupHandler.CreateGroup)
			groups.GET("", config.GroupHandler.GetUserGroups)
			groups.GET("/:id", config.GroupHandler.GetGroup)
			groups.PATCH("/:id", config.GroupHandler.UpdateGroup)
			groups.DELETE("/:id", config.GroupHandler.DeleteGroup)
			groups.POST("/:id/members", config.GroupHandler.AddMember)
			groups.DELETE("/:id/members/:user_id", config.GroupHandler.RemoveMember)
			groups.PATCH("/:id/members/:user_id", config.GroupHandler.UpdateMemberRole)
			// Group expenses - nested under groups
			groups.GET("/:id/expenses", config.ExpenseHandler.GetGroupExpenses)
			// Group balances and settlements
			groups.GET("/:id/balances", config.BalanceHandler.GetGroupBalances)
			groups.GET("/:id/settlement-suggestions", config.BalanceHandler.GetGroupSettlementSuggestions)
			groups.GET("/:id/settlements", config.SettlementHandler.GetGroupSettlements)
		}

		// Expense endpoints
		expenses := v1.Group("/expenses")
		expenses.Use(middleware.AuthMiddleware(config.JWTSecret))
		{
			expenses.POST("", config.ExpenseHandler.CreateExpense)
			expenses.GET("", config.ExpenseHandler.GetUserExpenses)
			expenses.GET("/filter", config.ExpenseHandler.GetExpensesWithFilters)

			// Settlement endpoints
			settlements := v1.Group("/settlements")
			settlements.Use(middleware.AuthMiddleware(config.JWTSecret))
			{
				settlements.POST("", config.SettlementHandler.CreateSettlement)
				settlements.GET("", config.SettlementHandler.GetUserSettlements)
				settlements.GET("/between", config.SettlementHandler.GetSettlementsBetweenUsers)
				settlements.GET("/suggestions", config.BalanceHandler.GetSettlementSuggestions)
				settlements.GET("/:id", config.SettlementHandler.GetSettlement)
				settlements.PATCH("/:id", config.SettlementHandler.UpdateSettlement)
				settlements.PATCH("/:id/confirm", config.SettlementHandler.ConfirmSettlement)
				settlements.DELETE("/:id", config.SettlementHandler.DeleteSettlement)
			}
			expenses.GET("/:id", config.ExpenseHandler.GetExpense)
			expenses.PATCH("/:id", config.ExpenseHandler.UpdateExpense)
			expenses.DELETE("/:id", config.ExpenseHandler.DeleteExpense)
		}
	}

	return router
}
