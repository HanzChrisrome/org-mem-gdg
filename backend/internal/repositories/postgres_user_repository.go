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
	query := `SELECT member_id, name, email, student_id, course, contact_number, registration_status, password_hash, created_at, last_updated
	          FROM members WHERE member_id = $1`

	user := &config.Member{}
	var contactNumber, course *string
	err := r.conn.QueryRow(ctx, query, id).Scan(
		&user.ID, &user.Name, &user.Email, &user.StudentID, &course, &contactNumber, &user.RegistrationStatus, &user.PasswordHash, &user.CreatedAt, &user.LastUpdated,
	)

	if course != nil {
		user.Course = *course
	}
	if contactNumber != nil {
		user.ContactNumber = *contactNumber
	}

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, config.ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to get user by id: %w", err)
	}

	return user, nil
}

func (r *PostgresUserRepository) GetByEmail(ctx context.Context, email string) (*config.Member, error) {
	query := `SELECT member_id, name, email, student_id, course, contact_number, registration_status, password_hash, created_at, last_updated
	          FROM members WHERE email = $1`

	user := &config.Member{}
	var contactNumber, course *string
	err := r.conn.QueryRow(ctx, query, email).Scan(
		&user.ID, &user.Name, &user.Email, &user.StudentID, &course, &contactNumber, &user.RegistrationStatus, &user.PasswordHash, &user.CreatedAt, &user.LastUpdated,
	)

	if course != nil {
		user.Course = *course
	}
	if contactNumber != nil {
		user.ContactNumber = *contactNumber
	}

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, config.ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}

	return user, nil
}

func (r *PostgresUserRepository) GetByStudentID(ctx context.Context, studentID string) (*config.Member, error) {
	query := `SELECT member_id, name, email, student_id, course, contact_number, registration_status, password_hash, created_at, last_updated
	          FROM members WHERE student_id = $1`

	user := &config.Member{}
	var contactNumber, course *string
	err := r.conn.QueryRow(ctx, query, studentID).Scan(
		&user.ID, &user.Name, &user.Email, &user.StudentID, &course, &contactNumber, &user.RegistrationStatus, &user.PasswordHash, &user.CreatedAt, &user.LastUpdated,
	)

	if course != nil {
		user.Course = *course
	}
	if contactNumber != nil {
		user.ContactNumber = *contactNumber
	}

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, config.ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to get user by student id: %w", err)
	}

	return user, nil
}

func (r *PostgresUserRepository) Create(ctx context.Context, user *config.Member) error {
	query := `INSERT INTO members (name, email, student_id, course, contact_number, registration_status, password_hash, created_at, last_updated)
	          VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING member_id`

	now := time.Now()
	user.CreatedAt = now
	user.LastUpdated = now
	if user.RegistrationStatus == "" {
		user.RegistrationStatus = config.StatusPending
	}

	err := r.conn.QueryRow(ctx, query, user.Name, user.Email, user.StudentID, user.Course, user.ContactNumber, user.RegistrationStatus, user.PasswordHash, user.CreatedAt, user.LastUpdated).Scan(&user.ID)

	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

func (r *PostgresUserRepository) Update(ctx context.Context, user *config.Member) error {
	query := `UPDATE members SET name = $2, email = $3, student_id = $4, course = $5, contact_number = $6, registration_status = $7, last_updated = $8 WHERE member_id = $1`

	now := time.Now()
	user.LastUpdated = now

	_, err := r.conn.Exec(ctx, query, user.ID, user.Name, user.Email, user.StudentID, user.Course, user.ContactNumber, user.RegistrationStatus, user.LastUpdated)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	return nil
}

func (r *PostgresUserRepository) Delete(ctx context.Context, id string) error {
	query := `UPDATE members SET registration_status = $2, last_updated = $3 WHERE member_id = $1`

	_, err := r.conn.Exec(ctx, query, id, config.StatusInactive, time.Now())
	if err != nil {
		return fmt.Errorf("failed to soft delete user: %w", err)
	}

	return nil
}

func (r *PostgresUserRepository) List(ctx context.Context, searchTerm string, status string) ([]config.MemberWithPayment, error) {
	// JOIN logic to get latest payment summary
	query := `
		SELECT
			m.member_id, m.name, m.email, m.student_id, m.course, m.contact_number, m.registration_status, m.created_at, m.last_updated,
			p.payment_id, p.payment_status, p.submission_date, p.approval_date, e.name as approver_name
		FROM members m
		LEFT JOIN LATERAL (
			SELECT payment_id, payment_status, submission_date, approval_date, approved_by
			FROM payments
			WHERE member_id = m.member_id
			ORDER BY submission_date DESC, payment_id DESC
			LIMIT 1
		) p ON true
		LEFT JOIN executives e ON p.approved_by = e.executive_id
		WHERE ($1 = '' OR m.name ILIKE '%' || $1 || '%' OR m.student_id ILIKE '%' || $1 || '%')
		  AND ($2 = '' OR m.registration_status = $2)
		ORDER BY m.created_at DESC
	`

	rows, err := r.conn.Query(ctx, query, searchTerm, status)
	if err != nil {
		return nil, fmt.Errorf("failed to list members: %w", err)
	}
	defer rows.Close()

	var members []config.MemberWithPayment
	for rows.Next() {
		var m config.MemberWithPayment
		var contactNumber, course, approverName *string
		err := rows.Scan(
			&m.ID, &m.Name, &m.Email, &m.StudentID, &course, &contactNumber, &m.RegistrationStatus, &m.CreatedAt, &m.LastUpdated,
			&m.LatestPaymentID, &m.LatestPaymentStatus, &m.LatestSubmission, &m.LatestApprovalDate, &approverName,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan member: %w", err)
		}
		if course != nil {
			m.Course = *course
		}
		if contactNumber != nil {
			m.ContactNumber = *contactNumber
		}
		m.ApproverName = approverName // m.ApproverName is already *string in config.MemberWithPayment
		members = append(members, m)
	}

	return members, nil
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
