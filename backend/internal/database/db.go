package database

import (
	"context"
	"log"
	"time"

	"github.com/jackc/pgx/v5"
)

func NewConnection(databaseURL string) *pgx.Conn {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := pgx.Connect(ctx, databaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	return conn
}
