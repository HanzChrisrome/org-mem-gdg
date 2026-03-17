package repositories

import (
	"context"
	"errors"
	"fmt"

	"github.com/HanzChrisrome/org-man-app/internal/config"
	"github.com/jackc/pgx/v5"
)

type PostgresExecutiveRepository struct {
	conn *pgx.Conn
}

func NewPostgresExecutiveRepository(conn *pgx.Conn) *PostgresExecutiveRepository {
	return &PostgresExecutiveRepository{conn: conn}
}

func (r *PostgresExecutiveRepository) GetByID(ctx context.Context, id string) (*config.Executive, error) {
	query := `SELECT executive_id, name, email, student_id, role_id, password_hash, created_at, last_updated
	          FROM executives WHERE executive_id = $1`

	exec := &config.Executive{}
	err := r.conn.QueryRow(ctx, query, id).Scan(
		&exec.ID, &exec.Name, &exec.Email, &exec.StudentID, &exec.RoleID, &exec.PasswordHash, &exec.CreatedAt, &exec.LastUpdated,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, config.ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to get executive by id: %w", err)
	}

	return exec, nil
}

func (r *PostgresExecutiveRepository) GetByEmail(ctx context.Context, email string) (*config.Executive, error) {
	query := `SELECT executive_id, name, email, student_id, role_id, password_hash, created_at, last_updated
	          FROM executives WHERE email = $1`

	exec := &config.Executive{}
	err := r.conn.QueryRow(ctx, query, email).Scan(
		&exec.ID, &exec.Name, &exec.Email, &exec.StudentID, &exec.RoleID, &exec.PasswordHash, &exec.CreatedAt, &exec.LastUpdated,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, config.ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to get executive by email: %w", err)
	}

	return exec, nil
}

func (r *PostgresExecutiveRepository) GetByStudentID(ctx context.Context, studentID string) (*config.Executive, error) {
	query := `SELECT executive_id, name, email, student_id, role_id, password_hash, created_at, last_updated
	          FROM executives WHERE student_id = $1`

	exec := &config.Executive{}
	err := r.conn.QueryRow(ctx, query, studentID).Scan(
		&exec.ID, &exec.Name, &exec.Email, &exec.StudentID, &exec.RoleID, &exec.PasswordHash, &exec.CreatedAt, &exec.LastUpdated,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, config.ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to get executive by student id: %w", err)
	}

	return exec, nil
}

func (r *PostgresExecutiveRepository) Exists(ctx context.Context, email, studentID string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM executives WHERE email = $1 OR student_id = $2)`
	var exists bool
	err := r.conn.QueryRow(ctx, query, email, studentID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check executive existence: %w", err)
	}
	return exists, nil
}
