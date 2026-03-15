package database

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5"
)

func NewConnection(databaseURL string) *pgx.Conn {
	conn, err := pgx.Connect(context.Background(), databaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	return conn
}
