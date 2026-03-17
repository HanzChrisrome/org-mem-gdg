package repositories

import (
	"context"
	"time"

	"github.com/HanzChrisrome/org-man-app/internal/config"
)

type SessionRepository interface {
	Create(ctx context.Context, session *config.Session) error
	Update(ctx context.Context, session *config.Session) error
	GetByRefreshTokenID(ctx context.Context, id string) (*config.Session, error)
	RevokeByRefreshTokenID(ctx context.Context, id string, revokedAt time.Time) error
}
