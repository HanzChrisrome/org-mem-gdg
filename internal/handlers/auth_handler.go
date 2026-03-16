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

type refreshRequest struct {
	SessionID    string `json:"session_id"`
	RefreshToken string `json:"refresh_token"`
}

type logoutRequest struct {
	SessionID string `json:"session_id"`
}

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

func (h *AuthHandler) Login(c *gin.Context) {
	var req config.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	userAgent := c.Request.UserAgent()
	ip := c.ClientIP()

	user, tokenPair, err := h.authService.LoginWithToken(c.Request.Context(), req, userAgent, ip)
	if err != nil {
		handleAuthError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user":  user,
		"token": tokenPair,
	})
}

func (h *AuthHandler) Refresh(c *gin.Context) {
	var req refreshRequest
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

func (h *AuthHandler) Logout(c *gin.Context) {
	var req logoutRequest
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
		log.Printf("Internal Server Error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": config.ErrInternal.Error()})
	}
}
