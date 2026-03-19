package utils

import (
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"time"

	"github.com/HanzChrisrome/org-man-app/internal/config"
)

type SessionManager interface {
	CreateSession(ownerID, ownerType, refreshTokenID, refreshToken, userAgent, ip string) (*config.Session, error)
	ValidateSession(session *config.Session, refreshToken string) error
	GetRefreshTTL() time.Duration
	GetMaxSessionTTL() time.Duration
	RevokeSession(session *config.Session)
}

type DefaultSessionManager struct {
	refreshTTL    time.Duration
	maxSessionTTL time.Duration
}

func NewDefaultSessionManager(ttlHours, maxTTLHours int) *DefaultSessionManager {
	if ttlHours <= 0 {
		ttlHours = 1
	}
	if maxTTLHours <= 0 {
		maxTTLHours = 24
	}
	return &DefaultSessionManager{
		refreshTTL:    time.Duration(ttlHours) * time.Hour,
		maxSessionTTL: time.Duration(maxTTLHours) * time.Hour,
	}
}

func (m *DefaultSessionManager) GetRefreshTTL() time.Duration {
	return m.refreshTTL
}

func (m *DefaultSessionManager) GetMaxSessionTTL() time.Duration {
	return m.maxSessionTTL
}

func (m *DefaultSessionManager) CreateSession(ownerID, ownerType, refreshTokenID, refreshToken, userAgent, ip string) (*config.Session, error) {
	if ownerID == "" || refreshTokenID == "" || refreshToken == "" {
		return nil, config.ErrInvalidInput
	}

	now := time.Now().UTC()
	return &config.Session{
		RefreshTokenID:   refreshTokenID,
		OwnerID:          ownerID,
		OwnerType:        ownerType,
		RefreshTokenHash: hashToken(refreshToken),
		UserAgent:        userAgent,
		IPAddress:        ip,
		ExpiresAt:        now.Add(m.refreshTTL),
		CreatedAt:        now,
	}, nil
}

func (m *DefaultSessionManager) ValidateSession(session *config.Session, refreshToken string) error {
	if session == nil {
		return config.ErrSessionNotFound
	}

	if session.RevokedAt != nil {
		return config.ErrSessionExpired
	}

	if time.Now().UTC().After(session.ExpiresAt) {
		return config.ErrSessionExpired
	}

	expected := hashToken(refreshToken)
	if subtle.ConstantTimeCompare([]byte(expected), []byte(session.RefreshTokenHash)) != 1 {
		return config.ErrInvalidToken
	}

	return nil
}

func (m *DefaultSessionManager) RevokeSession(session *config.Session) {
	if session == nil {
		return
	}
	now := time.Now().UTC()
	session.RevokedAt = &now
}

func HashToken(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}

func hashToken(token string) string {
	return HashToken(token)
}

func randomHex(size int) (string, error) {
	buf := make([]byte, size)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	return hex.EncodeToString(buf), nil
}
