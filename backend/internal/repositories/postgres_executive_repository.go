package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/HanzChrisrome/org-man-app/internal/config"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresExecutiveRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresExecutiveRepository(pool *pgxpool.Pool) *PostgresExecutiveRepository {
	return &PostgresExecutiveRepository{pool: pool}
}

func (r *PostgresExecutiveRepository) GetByID(ctx context.Context, id string) (*config.Executive, error) {
	query := `SELECT executive_id, name, email, student_id, role_id, password_hash, created_at, last_updated
	          FROM executives WHERE executive_id = $1`

	exec := &config.Executive{}
	var roleID sql.NullInt64
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&exec.ID, &exec.Name, &exec.Email, &exec.StudentID, &roleID, &exec.PasswordHash, &exec.CreatedAt, &exec.LastUpdated,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, config.ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to get executive by id: %w", err)
	}

	if roleID.Valid {
		exec.RoleID = int(roleID.Int64)
	} else {
		exec.RoleID = 0
	}

	return exec, nil
}

func (r *PostgresExecutiveRepository) GetByEmail(ctx context.Context, email string) (*config.Executive, error) {
	query := `SELECT executive_id, name, email, student_id, role_id, password_hash, created_at, last_updated
	          FROM executives WHERE email = $1`

	exec := &config.Executive{}
	var roleID sql.NullInt64
	err := r.pool.QueryRow(ctx, query, email).Scan(
		&exec.ID, &exec.Name, &exec.Email, &exec.StudentID, &roleID, &exec.PasswordHash, &exec.CreatedAt, &exec.LastUpdated,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, config.ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to get executive by email: %w", err)
	}

	if roleID.Valid {
		exec.RoleID = int(roleID.Int64)
	} else {
		exec.RoleID = 0
	}

	return exec, nil
}

func (r *PostgresExecutiveRepository) GetByStudentID(ctx context.Context, studentID string) (*config.Executive, error) {
	query := `SELECT executive_id, name, email, student_id, role_id, password_hash, created_at, last_updated
	          FROM executives WHERE student_id = $1`

	exec := &config.Executive{}
	var roleID sql.NullInt64
	err := r.pool.QueryRow(ctx, query, studentID).Scan(
		&exec.ID, &exec.Name, &exec.Email, &exec.StudentID, &roleID, &exec.PasswordHash, &exec.CreatedAt, &exec.LastUpdated,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, config.ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to get executive by student id: %w", err)
	}

	if roleID.Valid {
		exec.RoleID = int(roleID.Int64)
	} else {
		exec.RoleID = 0
	}

	return exec, nil
}

func (r *PostgresExecutiveRepository) Exists(ctx context.Context, email, studentID string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM executives WHERE email = $1 OR student_id = $2)`
	var exists bool
	err := r.pool.QueryRow(ctx, query, email, studentID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check executive existence: %w", err)
	}
	return exists, nil
}

func (r *PostgresExecutiveRepository) Create(ctx context.Context, exec *config.Executive) error {
	query := `INSERT INTO executives (name, email, student_id, role_id, password_hash, created_at, last_updated)
	          VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING executive_id`

	now := time.Now()
	exec.CreatedAt = now
	exec.LastUpdated = now

	var roleIDValue interface{}
	if exec.RoleID == 0 {
		roleIDValue = nil
	} else {
		roleIDValue = exec.RoleID
	}

	err := r.pool.QueryRow(ctx, query, exec.Name, exec.Email, exec.StudentID, roleIDValue, exec.PasswordHash, exec.CreatedAt, exec.LastUpdated).Scan(&exec.ID)
	if err != nil {
		return fmt.Errorf("failed to create executive: %w", err)
	}

	return nil
}

func (r *PostgresExecutiveRepository) List(ctx context.Context) ([]config.Executive, error) {
	query := `SELECT executive_id, name, email, student_id, role_id, password_hash, created_at, last_updated
	          FROM executives ORDER BY created_at DESC`

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to list executives: %w", err)
	}
	defer rows.Close()

	executives := make([]config.Executive, 0)
	for rows.Next() {
		exec := config.Executive{}
		var roleID sql.NullInt64
		err := rows.Scan(
			&exec.ID, &exec.Name, &exec.Email, &exec.StudentID, &roleID, &exec.PasswordHash, &exec.CreatedAt, &exec.LastUpdated,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan executive row: %w", err)
		}

		if roleID.Valid {
			exec.RoleID = int(roleID.Int64)
		} else {
			exec.RoleID = 0
		}

		executives = append(executives, exec)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return executives, nil
}

func (r *PostgresExecutiveRepository) Update(ctx context.Context, exec *config.Executive) error {
	query := `UPDATE executives SET name = $2, email = $3, student_id = $4, role_id = $5, password_hash = $6, last_updated = $7 WHERE executive_id = $1`

	now := time.Now()
	exec.LastUpdated = now

	var roleIDValue interface{}
	if exec.RoleID == 0 {
		roleIDValue = nil
	} else {
		roleIDValue = exec.RoleID
	}

	result, err := r.pool.Exec(ctx, query, exec.ID, exec.Name, exec.Email, exec.StudentID, roleIDValue, exec.PasswordHash, exec.LastUpdated)
	if err != nil {
		return fmt.Errorf("failed to update executive: %w", err)
	}

	if result.RowsAffected() == 0 {
		return config.ErrUserNotFound
	}

	return nil
}

func (r *PostgresExecutiveRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM executives WHERE executive_id = $1`

	result, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete executive: %w", err)
	}

	if result.RowsAffected() == 0 {
		return config.ErrUserNotFound
	}

	return nil
}
