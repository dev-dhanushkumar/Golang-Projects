package middleware

import (
	"bytes"
	"io"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// responseWriter is a custom response writer to capture response body
type responseWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w *responseWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

// LoggerMiddleware logs all HTTP requests and responses with detailed information
func LoggerMiddleware(logger *zap.SugaredLogger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Generate unique request ID
		requestID := uuid.New().String()
		c.Set("request_id", requestID)

		// Capture request start time
		startTime := time.Now()

		// Read and restore request body
		var requestBody []byte
		if c.Request.Body != nil {
			requestBody, _ = io.ReadAll(c.Request.Body)
			c.Request.Body = io.NopCloser(bytes.NewBuffer(requestBody))
		}

		// Wrap response writer to capture response
		responseWriter := &responseWriter{
			ResponseWriter: c.Writer,
			body:           &bytes.Buffer{},
		}
		c.Writer = responseWriter

		// Log incoming request
		logger.Infow("Incoming request",
			"request_id", requestID,
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"query", c.Request.URL.RawQuery,
			"ip", c.ClientIP(),
			"user_agent", c.Request.UserAgent(),
			"content_type", c.ContentType(),
			"request_body", string(requestBody),
		)

		// Process request
		c.Next()

		// Calculate request duration
		duration := time.Since(startTime)

		// Get response body (limit to 1000 characters to avoid huge logs)
		responseBody := responseWriter.body.String()
		if len(responseBody) > 1000 {
			responseBody = responseBody[:1000] + "... (truncated)"
		}

		// Log response details
		logLevel := logger.Infow
		if c.Writer.Status() >= 500 {
			logLevel = logger.Errorw
		} else if c.Writer.Status() >= 400 {
			logLevel = logger.Warnw
		}

		logLevel("Request completed",
			"request_id", requestID,
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"status", c.Writer.Status(),
			"duration_ms", duration.Milliseconds(),
			"duration", duration.String(),
			"response_size", c.Writer.Size(),
			"response_body", responseBody,
			"ip", c.ClientIP(),
		)

		// Log errors if any occurred
		if len(c.Errors) > 0 {
			for _, e := range c.Errors {
				logger.Errorw("Request error",
					"request_id", requestID,
					"error", e.Error(),
					"type", e.Type,
				)
			}
		}
	}
}
