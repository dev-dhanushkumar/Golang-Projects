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
	}

	return router
}
