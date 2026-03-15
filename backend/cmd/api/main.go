package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/HanzChrisrome/org-man-app/internal/config"
	"github.com/HanzChrisrome/org-man-app/internal/database"
)

func main() {
	cfg := config.LoadConfig()

	conn := database.NewConnection(cfg.DatabaseURL)
	defer conn.Close(context.Background())

	var version string
	err := conn.QueryRow(context.Background(), "SELECT version()").Scan(&version)
	if err != nil {
		log.Fatalf("Query failed: %v", err)
	}

	log.Println("Connected to:", version)

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("OK"))
	})

	log.Printf("Server running on :%s", cfg.Port)

	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}
