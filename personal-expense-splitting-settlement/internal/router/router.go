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
			// Balance summary endpoint will be added later when expense module is implemented
			// users.GET("/me/balance-summary", config.AuthHandler.GetBalanceSummary)
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
		}

		// Expense endpoints
		expenses := v1.Group("/expenses")
		expenses.Use(middleware.AuthMiddleware(config.JWTSecret))
		{
			expenses.POST("", config.ExpenseHandler.CreateExpense)
			expenses.GET("", config.ExpenseHandler.GetUserExpenses)
			expenses.GET("/filter", config.ExpenseHandler.GetExpensesWithFilters)
			expenses.GET("/:id", config.ExpenseHandler.GetExpense)
			expenses.PATCH("/:id", config.ExpenseHandler.UpdateExpense)
			expenses.DELETE("/:id", config.ExpenseHandler.DeleteExpense)
		}
	}

	return router
}
