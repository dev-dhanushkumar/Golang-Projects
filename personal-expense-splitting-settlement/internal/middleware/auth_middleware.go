package middleware

import (
	"personal-expense-splitting-settlement/pkg/utils"
	"strings"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware validate JWT token and set user content
func AuthMiddleware(jwtSecret string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			utils.Unauthorized(ctx, "Authorization header required", nil)
			ctx.Abort()
			return
		}

		// Extract token from "Bearer <token>"
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			utils.Unauthorized(ctx, "Invalid authorization header format", nil)
			ctx.Abort()
			return
		}

		tokenString := parts[1]

		// Validate token
		claims, err := utils.ValidateToken(tokenString, jwtSecret)
		if err != nil {
			utils.Unauthorized(ctx, "INvalid or expired token", err)
			ctx.Abort()
			return
		}

		// Set User context
		ctx.Set("user_id", claims.UserId)
		ctx.Set("session_id", claims.SessionId)
		ctx.Next()
	}
}
