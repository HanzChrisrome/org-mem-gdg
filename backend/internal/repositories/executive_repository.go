package repositories

import (
	"context"

	"github.com/HanzChrisrome/org-man-app/internal/config"
)

type ExecutiveRepository interface {
	GetByID(ctx context.Context, id string) (*config.Executive, error)
	GetByEmail(ctx context.Context, email string) (*config.Executive, error)
	GetByStudentID(ctx context.Context, studentID string) (*config.Executive, error)
	Exists(ctx context.Context, email, studentID string) (bool, error)
}
