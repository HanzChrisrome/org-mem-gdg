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
	SessionID    string `json:"session_id"`
	RefreshToken string `json:"refresh_token"`
}

type LogoutRequest struct {
	SessionID string `json:"session_id"`
}

// Register godoc
// @Summary Register member
// @Description Create a new member account.
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body config.RegisterRequest true "Register payload"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 409 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req config.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	user, err := h.authService.RegisterMember(c.Request.Context(), req)
	if err != nil {
		handleAuthError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"user": user,
	})
}

// Login godoc
// @Summary Login
// @Description Authenticate a member and return a token pair.
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body config.LoginRequest true "Login payload"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
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
// @Description Refresh an access token using session_id and refresh_token.
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body RefreshRequest true "Refresh payload"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/refresh [post]
func (h *AuthHandler) Refresh(c *gin.Context) {
	var req RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	sessionID, refreshToken := normalizeRefreshInput(req.SessionID, req.RefreshToken)
	if sessionID == "" || refreshToken == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "session_id and refresh_token are required"})
		return
	}

	tokenPair, err := h.authService.RefreshAccessToken(c.Request.Context(), sessionID, refreshToken)
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
// @Param request body LogoutRequest true "Logout payload"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/logout [post]
func (h *AuthHandler) Logout(c *gin.Context) {
	var req LogoutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	if req.SessionID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "session_id is required"})
		return
	}

	if err := h.authService.Logout(c.Request.Context(), req.SessionID); err != nil {
		handleAuthError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "logged out"})
}

func normalizeRefreshInput(sessionID, refreshToken string) (string, string) {
	if sessionID != "" {
		parts := strings.SplitN(refreshToken, ".", 2)
		if len(parts) == 2 {
			return sessionID, parts[1]
		}
		return sessionID, refreshToken
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
