package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// SuccessResponse sends a success response
func SuccessResponse(ctx *gin.Context, statusCode int, message string, data interface{}) {
	ctx.JSON(statusCode, Response{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// ErrorResponse sends an error response
func ErrorResponse(ctx *gin.Context, statusCode int, message string, err error) {
	response := Response{
		Success: false,
		Message: message,
	}

	if err != nil {
		response.Error = err.Error()
	}

	ctx.JSON(statusCode, response)
}

func ValidationErrorResponse(ctx *gin.Context, errors interface{}) {
	ctx.JSON(http.StatusBadRequest, Response{
		Success: false,
		Message: "Validation failed",
		Data:    errors,
	})
}

// 201 created response
func Created(ctx *gin.Context, message string, data interface{}) {
	SuccessResponse(ctx, http.StatusCreated, message, data)
}

// 201 created response
func OK(ctx *gin.Context, message string, data interface{}) {
	SuccessResponse(ctx, http.StatusOK, message, data)
}

// BadRequest sends a 400 Bad Request response
func BadRequest(c *gin.Context, message string, err error) {
	ErrorResponse(c, http.StatusBadRequest, message, err)
}

// Unauthorized sends a 401 Unauthorized response
func Unauthorized(c *gin.Context, message string, err error) {
	ErrorResponse(c, http.StatusUnauthorized, message, err)
}

// Forbidden sends a 403 Forbidden response
func Forbidden(c *gin.Context, message string, err error) {
	ErrorResponse(c, http.StatusForbidden, message, err)
}

// NotFound sends a 404 Not Found response
func NotFound(c *gin.Context, message string, err error) {
	ErrorResponse(c, http.StatusNotFound, message, err)
}

// InternalServerError sends a 500 Internal Server Error response
func InternalServerError(c *gin.Context, message string, err error) {
	ErrorResponse(c, http.StatusInternalServerError, message, err)
}
