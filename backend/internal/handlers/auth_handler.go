package handlers

import (
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/HanzChrisrome/org-man-app/internal/config"
	"github.com/HanzChrisrome/org-man-app/internal/services"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService *services.AuthService
}

func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

type RefreshRequest struct {
	RefreshTokenID string `json:"refresh_token_id"`
	RefreshToken   string `json:"refresh_token"`
}

// Register godoc
// @Summary Register
// @Description Create a new member or executive account based on the source dashboard.
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body RegisterRequestDoc true "Register payload"
// @Success 201 {object} RegisterResponse
// @Failure 400 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req config.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	user, ownerType, err := h.authService.Register(c.Request.Context(), req)
	if err != nil {
		handleAuthError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"user":       user,
		"owner_type": ownerType,
	})
}

// Login godoc
// @Summary Login
// @Description Authenticate an executive and return a token pair. Member login is not supported.
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body LoginRequestDoc true "Login payload"
// @Success 200 {object} LoginResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req config.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	userAgent := c.Request.UserAgent()
	ip := c.ClientIP()

	userID, ownerType, tokenPair, err := h.authService.LoginWithToken(c.Request.Context(), req, userAgent, ip)
	if err != nil {
		handleAuthError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user_id":    userID,
		"owner_type": ownerType,
		"token":      tokenPair,
	})
}

// Refresh godoc
// @Summary Refresh access token
// @Description Refresh an access token using refresh_token_id and refresh_token.
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body RefreshRequest true "Refresh payload"
// @Success 200 {object} RefreshResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/refresh [post]
func (h *AuthHandler) Refresh(c *gin.Context) {
	var req RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	refreshTokenID, refreshTokenSecret := normalizeRefreshInput(req.RefreshTokenID, req.RefreshToken)
	if refreshTokenID == "" || refreshTokenSecret == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "refresh_token_id and refresh_token are required"})
		return
	}

	tokenPair, err := h.authService.RefreshAccessToken(c.Request.Context(), refreshTokenID, refreshTokenSecret)
	if err != nil {
		handleAuthError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": tokenPair,
	})
}

// Logout godoc
// @Summary Logout
// @Description Revoke an active session.
// @Tags Auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} MessageResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/logout [post]
func (h *AuthHandler) Logout(c *gin.Context) {
	refreshTokenID, exists := c.Get("refresh_token_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "session context missing"})
		return
	}

	idStr, ok := refreshTokenID.(string)
	if !ok || idStr == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid session context"})
		return
	}

	if err := h.authService.Logout(c.Request.Context(), idStr); err != nil {
		handleAuthError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "logged out"})
}

// RevokeSession godoc
// @Summary Revoke session
// @Description Revoke a session via workspace or admin action.
// @Tags Auth
// @Accept json
// @Produce json
// @Param id path string true "Refresh Token ID"
// @Success 200 {object} MessageResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/sessions/{id}/revoke [post]
func (h *AuthHandler) RevokeSession(c *gin.Context) {
	refreshTokenID := c.Param("id")
	if refreshTokenID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "refresh_token_id is required"})
		return
	}

	if err := h.authService.RevokeSession(c.Request.Context(), refreshTokenID); err != nil {
		handleAuthError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "session revoked"})
}

func normalizeRefreshInput(refreshTokenID, refreshToken string) (string, string) {
	if refreshTokenID != "" {
		parts := strings.SplitN(refreshToken, ".", 2)
		if len(parts) == 2 {
			return refreshTokenID, parts[1]
		}
		return refreshTokenID, refreshToken
	}

	parts := strings.SplitN(refreshToken, ".", 2)
	if len(parts) != 2 {
		return "", ""
	}
	return parts[0], parts[1]
}

func handleAuthError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, config.ErrInvalidCredentials):
		c.JSON(http.StatusUnauthorized, gin.H{"error": config.ErrInvalidCredentials.Error()})
	case errors.Is(err, config.ErrWeakPassword), errors.Is(err, config.ErrInvalidInput):
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	case errors.Is(err, config.ErrUserAlreadyExists):
		c.JSON(http.StatusConflict, gin.H{"error": config.ErrUserAlreadyExists.Error()})
	case errors.Is(err, config.ErrSessionNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": config.ErrSessionNotFound.Error()})
	case errors.Is(err, config.ErrSessionExpired), errors.Is(err, config.ErrInvalidToken):
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
	case errors.Is(err, config.ErrInternal):
		log.Printf("Internal Server Error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": config.ErrInternal.Error()})
	default:
		log.Printf("Internal Server Error: %+v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": config.ErrInternal.Error()})
	}
}
