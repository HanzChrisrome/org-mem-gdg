package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/HanzChrisrome/org-man-app/internal/config"
	"github.com/HanzChrisrome/org-man-app/internal/services"
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

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var req config.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	user, err := h.authService.RegisterMember(r.Context(), req)
	if err != nil {
		handleAuthError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, map[string]any{
		"user": user,
	})
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var req config.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	userAgent := r.UserAgent()
	ip := r.RemoteAddr

	user, tokenPair, err := h.authService.LoginWithToken(r.Context(), req, userAgent, ip)
	if err != nil {
		handleAuthError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"user":  user,
		"token": tokenPair,
	})
}

func (h *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var req refreshRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	sessionID, refreshToken := normalizeRefreshInput(req.SessionID, req.RefreshToken)
	if sessionID == "" || refreshToken == "" {
		writeJSONError(w, http.StatusBadRequest, "session_id and refresh_token are required")
		return
	}

	tokenPair, err := h.authService.RefreshAccessToken(r.Context(), sessionID, refreshToken)
	if err != nil {
		handleAuthError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"token": tokenPair,
	})
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var req logoutRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.SessionID == "" {
		writeJSONError(w, http.StatusBadRequest, "session_id is required")
		return
	}

	if err := h.authService.Logout(r.Context(), req.SessionID); err != nil {
		handleAuthError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"message": "logged out"})
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

func handleAuthError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, config.ErrInvalidCredentials):
		writeJSONError(w, http.StatusUnauthorized, config.ErrInvalidCredentials.Error())
	case errors.Is(err, config.ErrWeakPassword), errors.Is(err, config.ErrInvalidInput):
		writeJSONError(w, http.StatusBadRequest, err.Error())
	case errors.Is(err, config.ErrUserAlreadyExists):
		writeJSONError(w, http.StatusConflict, config.ErrUserAlreadyExists.Error())
	case errors.Is(err, config.ErrSessionNotFound):
		writeJSONError(w, http.StatusNotFound, config.ErrSessionNotFound.Error())
	case errors.Is(err, config.ErrSessionExpired), errors.Is(err, config.ErrInvalidToken):
		writeJSONError(w, http.StatusUnauthorized, err.Error())
	case errors.Is(err, config.ErrInternal):
		log.Printf("Internal Server Error: %v", err)
		writeJSONError(w, http.StatusInternalServerError, config.ErrInternal.Error())
	default:
		log.Printf("Internal Server Error: %v", err)
		writeJSONError(w, http.StatusInternalServerError, config.ErrInternal.Error())
	}
}

func writeJSONError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}
