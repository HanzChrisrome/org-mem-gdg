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
func main() {
    cfg := config.LoadConfig()

    conn := database.NewConnection(cfg.DatabaseURL)
    defer conn.Close(context.Background())

    // Composition Root - Wired for future handler injection
    userRepo := repositories.NewPostgresUserRepository(conn)
    sessionRepo := repositories.NewPostgresSessionRepository(conn)
    hasher := utils.NewBcryptHasher(cfg.BcryptCost)
    validator := utils.NewPasswordValidator(cfg.MinPassLen)
    jwtManager := utils.NewHMACJWTManager(cfg.JWTSecret, cfg.JWTIssuer, cfg.AccessTokenTTLMinutes)
    sessionManager := utils.NewDefaultSessionManager(cfg.RefreshTokenTTLHours)
    authService := services.NewAuthService(userRepo, sessionRepo, hasher, validator, jwtManager, sessionManager)

    var version string
    err := conn.QueryRow(context.Background(), "SELECT version()").Scan(&version)
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
    routes.Register(router, healthHandler, authHandler)

    log.Printf("Server running on :%s", cfg.Port)

    srv := &http.Server{
        Addr:         ":" + cfg.Port,
        Handler:      router,
        ReadTimeout:  10 * time.Second,
        WriteTimeout: 10 * time.Second,
        IdleTimeout:  120 * time.Second,
    }

    log.Fatal(srv.ListenAndServe())
}
