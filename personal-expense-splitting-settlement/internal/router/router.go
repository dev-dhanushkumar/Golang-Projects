package router

import (
	"personal-expense-splitting-settlement/internal/handler"
	"personal-expense-splitting-settlement/internal/middleware"

	"github.com/gin-gonic/gin"
)

type RouterConfig struct {
	AuthHandler *handler.AuthHandler
	JWTSecret   string
}

func SetupRouter(config RouterConfig) *gin.Engine {
	router := gin.New()

	// Global Middleware
	router.Use(gin.Recovery())

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
	}

	return router
}
