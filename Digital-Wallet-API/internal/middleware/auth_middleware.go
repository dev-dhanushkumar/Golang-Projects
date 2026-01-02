package middleware

import (
	"digital-wallet-api/pkg/utils"
	"strings"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware validates JWT token and sets user context
func AuthMiddleware(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			utils.Unauthorized(c, "Authorization header required", nil)
			c.Abort()
			return
		}

		// Extract token from "Bearer <token>"
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			utils.Unauthorized(c, "Invalid authorization header format", nil)
			c.Abort()
			return
		}

		tokenString := parts[1]

		// Validate token
		claims, err := utils.ValidateToken(tokenString, jwtSecret)
		if err != nil {
			utils.Unauthorized(c, "Invalid or expired token", err)
			c.Abort()
			return
		}

		// Set user context
		c.Set("user_id", claims.UserID.String())
		c.Set("email", claims.Email)

		c.Next()
	}
}
