package routes

import (
	"github.com/HanzChrisrome/org-man-app/internal/handlers"
	"github.com/gin-gonic/gin"
)

func Register(router *gin.Engine, healthHandler *handlers.HealthHandler, authHandler *handlers.AuthHandler) {
	router.GET("/health", healthHandler.Health)

	api := router.Group("/api")
	{
		api.POST("/register", authHandler.Register)
		api.POST("/login", authHandler.Login)
		api.POST("/refresh", authHandler.Refresh)
		api.POST("/logout", authHandler.Logout)
	}
}
