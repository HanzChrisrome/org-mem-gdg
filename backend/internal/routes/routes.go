package routes

import (
	"net/http"

	"github.com/HanzChrisrome/org-man-app/internal/handlers"
)

func Register(mux *http.ServeMux, healthHandler *handlers.HealthHandler, authHandler *handlers.AuthHandler) {
	mux.HandleFunc("/health", healthHandler.Health)

	mux.HandleFunc("/api/register", authHandler.Register)
	mux.HandleFunc("/api/login", authHandler.Login)
	mux.HandleFunc("/api/refresh", authHandler.Refresh)
	mux.HandleFunc("/api/logout", authHandler.Logout)
}
