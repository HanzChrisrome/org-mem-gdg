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

	pool := database.NewConnection(cfg.DatabaseURL)
	defer pool.Close()

	// Composition Root - Wired for future handler injection
	userRepo := repositories.NewPostgresUserRepository(pool)
	execRepo := repositories.NewPostgresExecutiveRepository(pool)
	sessionRepo := repositories.NewPostgresSessionRepository(pool)
	hasher := utils.NewBcryptHasher(cfg.BcryptCost)
	validator := utils.NewPasswordValidator(cfg.MinPassLen)
	jwtManager := utils.NewHMACJWTManager(cfg.JWTSecret, cfg.JWTIssuer, cfg.AccessTokenTTLMinutes)
	sessionManager := utils.NewDefaultSessionManager(cfg.RefreshTokenTTLHours, cfg.MaxSessionTTLHours)
	authService := services.NewAuthService(userRepo, execRepo, sessionRepo, hasher, validator, jwtManager, sessionManager)
	memberService := services.NewMemberService(userRepo, hasher, validator)
	executiveService := services.NewExecutiveService(execRepo, hasher, validator)

	var version string
	startupCtx, cancelStartup := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelStartup()
	err := pool.QueryRow(startupCtx, "SELECT version()").Scan(&version)
	if err != nil {
		log.Fatalf("Query failed: %v", err)
	}

	log.Println("Connected to:", version)

	// Pre-flight check: Verify critical auth columns exist to prevent 500/401 logic failures
	preflightCtx, cancelPreflight := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelPreflight()
	var exists bool
	preflightQuery := `
		SELECT EXISTS (
			SELECT 1 FROM information_schema.columns
			WHERE table_name = 'executives' AND column_name = 'password_hash'
		) AND EXISTS (
			SELECT 1 FROM information_schema.columns
			WHERE table_name = 'sessions' AND column_name = 'refresh_token_id'
		)`
	err = pool.QueryRow(preflightCtx, preflightQuery).Scan(&exists)
	if err != nil || !exists {
		log.Fatal("Startup pre-flight check failed: required database columns (executives.password_hash or sessions.refresh_token_id) are missing. Check migrations.")
	}
	log.Println("Database pre-flight check passed.")

	// Optional: Pool stats logging
	go func() {
		ticker := time.NewTicker(30 * time.Second)
		for range ticker.C {
			s := pool.Stat()
			log.Printf("[DB Pool] Total: %d, Acquired: %d, Idle: %d", s.TotalConns(), s.AcquiredConns(), s.IdleConns())
		}
	}()

	router := gin.Default()

	// Global Middlewares
	router.Use(middleware.CORS())
	router.Use(middleware.PayloadTracer(cfg))
	router.Use(middleware.RateLimiter(cfg.RateLimitRPS, cfg.RateLimitBurst))

	// Swagger endpoint
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	healthHandler := handlers.NewHealthHandler()
	authHandler := handlers.NewAuthHandler(authService)
	memberHandler := handlers.NewMemberHandler(memberService)
	executiveHandler := handlers.NewExecutiveHandler(executiveService)
	routes.Register(router, healthHandler, authHandler, memberHandler, executiveHandler, jwtManager, sessionRepo)

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
