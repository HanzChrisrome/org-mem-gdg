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
	CreateSession(ownerID, ownerType, refreshToken, userAgent, ip string) (*config.Session, error)
	ValidateSession(session *config.Session, refreshToken string) error
	RevokeSession(session *config.Session)
}

type DefaultSessionManager struct {
	refreshTTL time.Duration
}

func NewDefaultSessionManager(ttlHours int) *DefaultSessionManager {
	if ttlHours <= 0 {
		ttlHours = 168
	}
	return &DefaultSessionManager{refreshTTL: time.Duration(ttlHours) * time.Hour}
}

func (m *DefaultSessionManager) CreateSession(ownerID, ownerType, refreshToken, userAgent, ip string) (*config.Session, error) {
	if ownerID == "" || refreshToken == "" {
		return nil, config.ErrInvalidInput
	}

	sessionID, err := randomHex(32)
	if err != nil {
		return nil, config.ErrInternal
	}

	now := time.Now().UTC()
	return &config.Session{
		ID:               sessionID,
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

func hashToken(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}

func randomHex(size int) (string, error) {
	buf := make([]byte, size)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	return hex.EncodeToString(buf), nil
}
