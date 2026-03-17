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
	Create(ctx context.Context, exec *config.Executive) error
	List(ctx context.Context) ([]config.Executive, error)
	Update(ctx context.Context, exec *config.Executive) error
	Delete(ctx context.Context, id string) error
}
