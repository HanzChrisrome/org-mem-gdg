package repositories

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/HanzChrisrome/org-man-app/internal/config"
	"github.com/jackc/pgx/v5"
)

type PostgresUserRepository struct {
	conn *pgx.Conn
}

func NewPostgresUserRepository(conn *pgx.Conn) *PostgresUserRepository {
	return &PostgresUserRepository{conn: conn}
}

func (r *PostgresUserRepository) GetByID(ctx context.Context, id string) (*config.Member, error) {
	query := `SELECT member_id, name, email, student_id, password_hash, created_at, last_updated
	          FROM members WHERE member_id = $1`

	user := &config.Member{}
	err := r.conn.QueryRow(ctx, query, id).Scan(
		&user.ID, &user.Name, &user.Email, &user.StudentID, &user.PasswordHash, &user.CreatedAt, &user.LastUpdated,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, config.ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to get user by id: %w", err)
	}

	return user, nil
}

func (r *PostgresUserRepository) GetByEmail(ctx context.Context, email string) (*config.Member, error) {
	query := `SELECT member_id, name, email, student_id, password_hash, created_at, last_updated
	          FROM members WHERE email = $1`

	user := &config.Member{}
	err := r.conn.QueryRow(ctx, query, email).Scan(
		&user.ID, &user.Name, &user.Email, &user.StudentID, &user.PasswordHash, &user.CreatedAt, &user.LastUpdated,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, config.ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}

	return user, nil
}

func (r *PostgresUserRepository) GetByStudentID(ctx context.Context, studentID string) (*config.Member, error) {
	query := `SELECT member_id, name, email, student_id, password_hash, created_at, last_updated
	          FROM members WHERE student_id = $1`

	user := &config.Member{}
	err := r.conn.QueryRow(ctx, query, studentID).Scan(
		&user.ID, &user.Name, &user.Email, &user.StudentID, &user.PasswordHash, &user.CreatedAt, &user.LastUpdated,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, config.ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to get user by student id: %w", err)
	}

	return user, nil
}

func (r *PostgresUserRepository) Create(ctx context.Context, user *config.Member) error {
	query := `INSERT INTO members (name, email, student_id, password_hash, created_at, last_updated)
	          VALUES ($1, $2, $3, $4, $5, $6) RETURNING member_id`

	now := time.Now()
	user.CreatedAt = now
	user.LastUpdated = now

	err := r.conn.QueryRow(ctx, query, user.Name, user.Email, user.StudentID, user.PasswordHash, user.CreatedAt, user.LastUpdated).Scan(&user.ID)

	if err != nil {
		// pgx standard handles unique constraints by returning certain error patterns but for now generic catch
		// and simple check for duplicate keyword
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

func (r *PostgresUserRepository) Exists(ctx context.Context, email, studentID string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM members WHERE email = $1 OR student_id = $2)`
	var exists bool
	err := r.conn.QueryRow(ctx, query, email, studentID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check user existence: %w", err)
	}
	return exists, nil
}
