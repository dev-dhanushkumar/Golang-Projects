package middleware

import (
	"digital-wallet-api/pkg/logger"
	"digital-wallet-api/pkg/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ErrorMiddleware handles panics and errors
func ErrorMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				logger.Error("Panic recovered", map[string]interface{}{
					"error":      err,
					"path":       c.Request.URL.Path,
					"method":     c.Request.Method,
					"request_id": c.GetString("request_id"),
				})

				utils.InternalServerError(c, "Internal server error", nil)
				c.Abort()
			}
		}()

		c.Next()
	}
}

// CORSMiddleware handles CORS
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}
