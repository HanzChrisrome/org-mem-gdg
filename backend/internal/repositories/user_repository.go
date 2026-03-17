package repositories

import (
	"context"

	"github.com/HanzChrisrome/org-man-app/internal/config"
)

type UserRepository interface {
	GetByID(ctx context.Context, member_id string) (*config.Member, error)
	GetByEmail(ctx context.Context, email string) (*config.Member, error)
	GetByStudentID(ctx context.Context, studentID string) (*config.Member, error)
	Create(ctx context.Context, user *config.Member) error
	Update(ctx context.Context, user *config.Member) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, query string, status string) ([]config.MemberWithPayment, error)
	Exists(ctx context.Context, email, studentID string) (bool, error)
}
