package repositories

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/HanzChrisrome/org-man-app/internal/config"
	"github.com/jackc/pgx/v5"
)

type PostgresSessionRepository struct {
	conn *pgx.Conn
}

func NewPostgresSessionRepository(conn *pgx.Conn) *PostgresSessionRepository {
	return &PostgresSessionRepository{conn: conn}
}

func (r *PostgresSessionRepository) Create(ctx context.Context, session *config.Session) error {
	query := `INSERT INTO sessions (session_id, owner_id, owner_type, refresh_token_hash, user_agent, ip_address, expires_at, created_at, revoked_at)
		  VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`

	_, err := r.conn.Exec(
		ctx,
		query,
		session.ID,
		session.OwnerID,
		session.OwnerType,
		session.RefreshTokenHash,
		session.UserAgent,
		session.IPAddress,
		session.ExpiresAt,
		session.CreatedAt,
		session.RevokedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create session: %w", err)
	}

	return nil
}

func (r *PostgresSessionRepository) GetByID(ctx context.Context, id string) (*config.Session, error) {
	query := `SELECT session_id, owner_id, owner_type, refresh_token_hash, user_agent, ip_address, expires_at, created_at, revoked_at
		  FROM sessions WHERE session_id = $1`

	var session config.Session
	err := r.conn.QueryRow(ctx, query, id).Scan(
		&session.ID,
		&session.OwnerID,
		&session.OwnerType,
		&session.RefreshTokenHash,
		&session.UserAgent,
		&session.IPAddress,
		&session.ExpiresAt,
		&session.CreatedAt,
		&session.RevokedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, config.ErrSessionNotFound
		}
		return nil, fmt.Errorf("failed to get session by id: %w", err)
	}

	return &session, nil
}

func (r *PostgresSessionRepository) Update(ctx context.Context, session *config.Session) error {
	query := `UPDATE sessions SET refresh_token_hash = $2, ip_address = $3, user_agent = $4, expires_at = $5 WHERE session_id = $1`

	_, err := r.conn.Exec(
		ctx,
		query,
		session.ID,
		session.RefreshTokenHash,
		session.IPAddress,
		session.UserAgent,
		session.ExpiresAt,
	)
	if err != nil {
		return fmt.Errorf("failed to update session: %w", err)
	}

	return nil
}

func (r *PostgresSessionRepository) Revoke(ctx context.Context, id string, revokedAt time.Time) error {
	query := `UPDATE sessions SET revoked_at = $2 WHERE session_id = $1`

	cmd, err := r.conn.Exec(ctx, query, id, revokedAt)
	if err != nil {
		return fmt.Errorf("failed to revoke session: %w", err)
	}

	if cmd.RowsAffected() == 0 {
		return config.ErrSessionNotFound
	}

	return nil
}
