package handler

import (
	"fmt"
	"net/http"
	"personal-expense-splitting-settlement/internal/dto"
	"personal-expense-splitting-settlement/internal/services"
	"personal-expense-splitting-settlement/pkg/utils"
	"personal-expense-splitting-settlement/pkg/validator"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AuthHandler struct {
	authService    services.AuthService
	sessionService services.SessionService
}

func NewAuthHandler(authService services.AuthService, sesionService services.SessionService) *AuthHandler {
	return &AuthHandler{
		authService:    authService,
		sessionService: sesionService,
	}
}

// Register New User
func (h *AuthHandler) Register(ctx *gin.Context) {
	var req dto.RegisterRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(ctx, "Invalid request body", err)
		return
	}

	// validate request
	if err := validator.Validate(req); err != nil {
		utils.ValidationErrorResponse(ctx, err.Error())
		return
	}

	user, err := h.authService.Register(req)
	if err != nil {
		utils.BadRequest(ctx, "Registeration failed", err)
		return
	}

	userId, _ := uuid.Parse(user.User.ID)

	accessToken, refreshToken, err := h.sessionService.CreateNewSession(
		userId,
		ctx.ClientIP(),
		ctx.Request.UserAgent(),
	)

	if err != nil {
		utils.Unauthorized(ctx, "Login Failed", err)
	}

	// Set Access and Refresh Token
	user.Tokens.AccessToken = accessToken
	user.Tokens.RefreshToken = refreshToken

	utils.Created(ctx, "User registered successfully", user)
}

// Login User
func (h *AuthHandler) Login(ctx *gin.Context) {
	var req dto.LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(ctx, "Invalid request body", err)
		return
	}

	// Validate request
	if err := validator.Validate(req); err != nil {
		utils.ValidationErrorResponse(ctx, err.Error())
		return
	}
	ip := ctx.ClientIP()
	userAgent := ctx.Request.UserAgent()

	loginResp, err := h.authService.Login(req)
	if err != nil {
		utils.Unauthorized(ctx, "Login Failed", err)
		return
	}

	userId, _ := uuid.Parse(loginResp.User.ID)

	accessToken, refreshToken, err := h.sessionService.CreateNewSession(userId, ip, userAgent)
	if err != nil {
		utils.Unauthorized(ctx, "Login Failed", err)
	}

	// Set Access and Refresh Token
	loginResp.Tokens.AccessToken = accessToken
	loginResp.Tokens.RefreshToken = refreshToken

	utils.OK(ctx, "Login Successfully", loginResp)
}

// Refresh Handler POST /api/v1/auth/refresh
func (h *AuthHandler) Refresh(ctx *gin.Context) {
	var req dto.RefreshRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(ctx, "Invalid request Body", err)
		return
	}

	// Call Service to perform rotation
	accessToken, refreshTOken, err := h.sessionService.RefreshSession(
		req.RefreshToken,
		ctx.ClientIP(),
		ctx.Request.UserAgent(),
	)

	if err != nil {
		utils.Unauthorized(ctx, "Sesion is no longer valid", err)
		return
	}

	res := &dto.TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshTOken,
		ExpiresIn:    3600,
	}

	utils.OK(ctx, "Token Rotation perform successfully", res)

}

// Logout Handler POST /api/v1/auth/logout
func (h *AuthHandler) Logout(ctx *gin.Context) {
	sessionIDRaw, exists := ctx.Get("session_id")
	if !exists {
		utils.InternalServerError(ctx, "Sesion not found in context", fmt.Errorf("SessionID Not set in context"))
		return
	}

	sesionID, ok := sessionIDRaw.(uuid.UUID)
	if !ok {
		utils.InternalServerError(ctx, "error", fmt.Errorf("Invalid session ID format"))
		return
	}

	err := h.sessionService.TerminateSession(sesionID)
	if err != nil {
		utils.InternalServerError(ctx, "error", fmt.Errorf("Failed to logout"))
		return
	}

	// Return 204 No content
	ctx.Status(http.StatusNoContent)
}

// GetMe handler GET /api/v1/auth/me
func (h *AuthHandler) GetMe(ctx *gin.Context) {
	// Get UserID from Auth middleware
	userIDRaw, exists := ctx.Get("user_id")
	if !exists {
		utils.Unauthorized(ctx, "User not found in context", nil)
		return
	}

	userID, ok := userIDRaw.(uuid.UUID)
	if !ok {
		utils.InternalServerError(ctx, "Invalid user ID format", nil)
		return
	}

	userProfile, err := h.authService.GetProfile(userID)
	if err != nil {
		utils.NotFound(ctx, "User profile not found", err)
		return
	}

	utils.OK(ctx, "User profile retrived successfully", userProfile)
}

func (h *AuthHandler) GetSessions(ctx *gin.Context) {
	// Get UserID from Auth middleware
	userIDRaw, exists := ctx.Get("user_id")
	if !exists {
		utils.Unauthorized(ctx, "User not found in context", nil)
		return
	}

	userID, ok := userIDRaw.(uuid.UUID)
	if !ok {
		utils.InternalServerError(ctx, "Invalid user ID format", nil)
		return
	}

	sessions, err := h.sessionService.GetUserSessions(userID)
	if err != nil {
		utils.InternalServerError(ctx, "Failed to fetch sessions", err)
		return
	}

	utils.OK(ctx, "Active sessions retrived seccessfully", sessions)
}

// UpdateProfile handler PATCH /api/v1/users/me
func (h *AuthHandler) UpdateProfile(ctx *gin.Context) {
	// Get UserID from Auth middleware
	userIDRaw, exists := ctx.Get("user_id")
	if !exists {
		utils.Unauthorized(ctx, "User not found in context", nil)
		return
	}

	userID, ok := userIDRaw.(uuid.UUID)
	if !ok {
		utils.InternalServerError(ctx, "Invalid user ID format", nil)
		return
	}

	var req dto.UpdateProfileRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(ctx, "Invalid request body", err)
		return
	}

	// Validate request
	if err := validator.Validate(req); err != nil {
		utils.ValidationErrorResponse(ctx, err.Error())
		return
	}

	userProfile, err := h.authService.UpdateProfile(userID, req)
	if err != nil {
		utils.InternalServerError(ctx, "Failed to update profile", err)
		return
	}

	utils.OK(ctx, "Profile updated successfully", userProfile)
}
