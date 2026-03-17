package routes

import (
	"github.com/HanzChrisrome/org-man-app/internal/handlers"
	"github.com/HanzChrisrome/org-man-app/internal/middleware"
	"github.com/HanzChrisrome/org-man-app/internal/repositories"
	"github.com/HanzChrisrome/org-man-app/internal/utils"
	"github.com/gin-gonic/gin"
)

func Register(router *gin.Engine, healthHandler *handlers.HealthHandler, authHandler *handlers.AuthHandler, memberHandler *handlers.MemberHandler, executiveHandler *handlers.ExecutiveHandler, jwtManager utils.JWTManager, sessionRepo repositories.SessionRepository) {
	router.GET("/health", healthHandler.Health)

	api := router.Group("/api")
	{
		api.POST("/register", authHandler.Register)
		api.POST("/login", authHandler.Login)
		api.POST("/refresh", authHandler.Refresh)

		// Protected routes
		protected := api.Group("")
		protected.Use(middleware.Auth(jwtManager, sessionRepo))
		{
			protected.POST("/logout", authHandler.Logout)
			protected.POST("/sessions/:id/revoke", authHandler.RevokeSession)

			// Member management (Executives only)
			members := protected.Group("/members")
			members.Use(middleware.RequireExecutive())
			{
				members.POST("", memberHandler.CreateMember)
				members.GET("", memberHandler.ListMembers)
				members.GET("/:id", memberHandler.GetMemberByID)
				members.PUT("/:id", memberHandler.UpdateMember)
				members.DELETE("/:id", memberHandler.DeleteMember)
			}

			executives := protected.Group("/executives")
			executives.Use(middleware.RequireExecutive())
			{
				executives.POST("", executiveHandler.CreateExecutive)
				executives.GET("", executiveHandler.ListExecutives)
				executives.GET("/:id", executiveHandler.GetExecutiveByID)
				executives.PUT("/:id", executiveHandler.UpdateExecutive)
				executives.DELETE("/:id", executiveHandler.DeleteExecutive)
			}
		}
	}
}
