package middleware

import (
	"net/http"
	"strings"

	"github.com/HanzChrisrome/org-man-app/internal/repositories"
	"github.com/HanzChrisrome/org-man-app/internal/utils"
	"github.com/gin-gonic/gin"
)

func Auth(jwtManager utils.JWTManager, sessionRepo repositories.SessionRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header must be Bearer token"})
			c.Abort()
			return
		}

		accessToken := parts[1]
		claims, err := jwtManager.ValidateAccessToken(accessToken)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired access token"})
			c.Abort()
			return
		}

		// Optional: Verify session is not revoked if SID is present
		if claims.SessionID != "" {
			session, err := sessionRepo.GetByID(c.Request.Context(), claims.SessionID)
			if err != nil || session.RevokedAt != nil {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "session revoked or not found"})
				c.Abort()
				return
			}
		}

		// Set user identity, session ID, and owner type in context
		c.Set("user_id", claims.UserID)
		c.Set("session_id", claims.SessionID)
		c.Set("owner_type", claims.OwnerType)
		c.Next()
	}
}

func RequireExecutive() gin.HandlerFunc {
	return func(c *gin.Context) {
		ownerType, exists := c.Get("owner_type")
		if !exists || ownerType != "executive" {
			c.JSON(http.StatusForbidden, gin.H{"error": "forbidden: executive access required"})
			c.Abort()
			return
		}
		c.Next()
	}
}
