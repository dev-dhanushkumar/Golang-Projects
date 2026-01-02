package handler

import (
	"digital-wallet-api/internal/dto"
	"digital-wallet-api/internal/service"
	"digital-wallet-api/pkg/utils"
	"digital-wallet-api/pkg/validator"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AuthHandler struct {
	authService service.AuthService
}

func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// Register godoc
// @Summary Register a new user
// @Description Register a new user with email and password
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.RegisterRequest true "Registration request"
// @Success 201 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Router /api/v1/auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "Invalid request body", err)
		return
	}

	// Validate request
	if err := validator.Validate(req); err != nil {
		utils.ValidationErrorResponse(c, err.Error())
		return
	}

	user, err := h.authService.Register(req)
	if err != nil {
		utils.BadRequest(c, "Registration failed", err)
		return
	}

	utils.Created(c, "User registered successfully", user)
}

// Login godoc
// @Summary Login user
// @Description Login user with email and password
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.LoginRequest true "Login request"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Router /api/v1/auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "Invalid request body", err)
		return
	}

	// Validate request
	if err := validator.Validate(req); err != nil {
		utils.ValidationErrorResponse(c, err.Error())
		return
	}

	loginResp, err := h.authService.Login(req)
	if err != nil {
		utils.Unauthorized(c, "Login failed", err)
		return
	}

	utils.OK(c, "Login successful", loginResp)
}

// GetProfile godoc
// @Summary Get user profile
// @Description Get authenticated user profile
// @Tags auth
// @Produce json
// @Security BearerAuth
// @Success 200 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Router /api/v1/auth/profile [get]
func (h *AuthHandler) GetProfile(c *gin.Context) {
	userID := c.GetString("user_id")
	uid, err := uuid.Parse(userID)
	if err != nil {
		utils.BadRequest(c, "Invalid user ID", err)
		return
	}

	profile, err := h.authService.GetProfile(uid)
	if err != nil {
		utils.NotFound(c, "User not found", err)
		return
	}

	utils.OK(c, "Profile retrieved successfully", profile)
}

// UpdateProfile godoc
// @Summary Update user profile
// @Description Update authenticated user profile
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.UpdateProfileRequest true "Update profile request"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Router /api/v1/auth/profile [put]
func (h *AuthHandler) UpdateProfile(c *gin.Context) {
	userID := c.GetString("user_id")
	uid, err := uuid.Parse(userID)
	if err != nil {
		utils.BadRequest(c, "Invalid user ID", err)
		return
	}

	var req dto.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "Invalid request body", err)
		return
	}

	// Validate request
	if err := validator.Validate(req); err != nil {
		utils.ValidationErrorResponse(c, err.Error())
		return
	}

	profile, err := h.authService.UpdateProfile(uid, req)
	if err != nil {
		utils.BadRequest(c, "Failed to update profile", err)
		return
	}

	utils.OK(c, "Profile updated successfully", profile)
}
