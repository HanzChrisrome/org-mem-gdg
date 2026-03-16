package repositories

import (
	"context"
	"time"

	"github.com/HanzChrisrome/org-man-app/internal/config"
)

type SessionRepository interface {
	Create(ctx context.Context, session *config.Session) error
	GetByID(ctx context.Context, id string) (*config.Session, error)
	Revoke(ctx context.Context, id string, revokedAt time.Time) error
}
