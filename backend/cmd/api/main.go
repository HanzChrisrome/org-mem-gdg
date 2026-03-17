package main

import (
	"context"
	"log"
	"net/http"
	"time"

	_ "github.com/HanzChrisrome/org-man-app/docs"
	"github.com/HanzChrisrome/org-man-app/internal/config"
	"github.com/HanzChrisrome/org-man-app/internal/database"
	"github.com/HanzChrisrome/org-man-app/internal/handlers"
	"github.com/HanzChrisrome/org-man-app/internal/middleware"
	"github.com/HanzChrisrome/org-man-app/internal/repositories"
	"github.com/HanzChrisrome/org-man-app/internal/routes"
	"github.com/HanzChrisrome/org-man-app/internal/services"
	"github.com/HanzChrisrome/org-man-app/internal/utils"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Org Mem GDG API
// @version 1.0
// @description API documentation for org-mem-gdg backend
// @BasePath /

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type 'Bearer ' followed by a space and your JWT token.
func main() {
	cfg := config.LoadConfig()

	conn := database.NewConnection(cfg.DatabaseURL)
	defer conn.Close(context.Background())

	// Composition Root - Wired for future handler injection
	userRepo := repositories.NewPostgresUserRepository(conn)
	execRepo := repositories.NewPostgresExecutiveRepository(conn)
	sessionRepo := repositories.NewPostgresSessionRepository(conn)
	hasher := utils.NewBcryptHasher(cfg.BcryptCost)
	validator := utils.NewPasswordValidator(cfg.MinPassLen)
	jwtManager := utils.NewHMACJWTManager(cfg.JWTSecret, cfg.JWTIssuer, cfg.AccessTokenTTLMinutes)
	sessionManager := utils.NewDefaultSessionManager(cfg.RefreshTokenTTLHours)
	authService := services.NewAuthService(userRepo, execRepo, sessionRepo, hasher, validator, jwtManager, sessionManager)
	memberService := services.NewMemberService(userRepo, hasher, validator)

	var version string
	startupCtx, cancelStartup := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelStartup()
	err := conn.QueryRow(startupCtx, "SELECT version()").Scan(&version)
	if err != nil {
		log.Fatalf("Query failed: %v", err)
	}

	log.Println("Connected to:", version)

	router := gin.Default()

	// Add CORS middleware
	router.Use(middleware.CORS())

	// Swagger endpoint
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	healthHandler := handlers.NewHealthHandler()
	authHandler := handlers.NewAuthHandler(authService)
	memberHandler := handlers.NewMemberHandler(memberService)
	routes.Register(router, healthHandler, authHandler, memberHandler, jwtManager, sessionRepo)

	log.Printf("Server running on :%s", cfg.Port)

	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	err = srv.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		log.Fatalf("Server failed: %v", err)
	}
}
