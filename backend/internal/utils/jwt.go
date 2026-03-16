package utils

import (
	"fmt"
	"time"

	"github.com/HanzChrisrome/org-man-app/internal/config"
	"github.com/golang-jwt/jwt/v5"
)

type AccessClaims struct {
	UserID string `json:"uid"`
	jwt.RegisteredClaims
}

type JWTManager interface {
	GenerateAccessToken(user *config.Member) (string, int64, error)
	ValidateAccessToken(token string) (*AccessClaims, error)
}

type HMACJWTManager struct {
	secret    []byte
	issuer    string
	accessTTL time.Duration
}

func NewHMACJWTManager(secret, issuer string, ttlMinutes int) *HMACJWTManager {
	if ttlMinutes <= 0 {
		ttlMinutes = 15
	}
	if issuer == "" {
		issuer = "org-man-app"
	}
	return &HMACJWTManager{
		secret:    []byte(secret),
		issuer:    issuer,
		accessTTL: time.Duration(ttlMinutes) * time.Minute,
	}
}

func (m *HMACJWTManager) GenerateAccessToken(user *config.Member) (string, int64, error) {
	if user == nil || user.ID == "" {
		return "", 0, config.ErrInvalidInput
	}

	now := time.Now().UTC()
	expiresAt := now.Add(m.accessTTL)
	claims := AccessClaims{
		UserID: user.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   user.ID,
			Issuer:    m.issuer,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(expiresAt),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(m.secret)
	if err != nil {
		return "", 0, fmt.Errorf("failed to sign jwt: %w", err)
	}

	return tokenString, int64(m.accessTTL.Seconds()), nil
}

func (m *HMACJWTManager) ValidateAccessToken(token string) (*AccessClaims, error) {
	claims := &AccessClaims{}
	parsed, err := jwt.ParseWithClaims(token, claims, func(t *jwt.Token) (interface{}, error) {
		if t.Method != jwt.SigningMethodHS256 {
			return nil, config.ErrInvalidToken
		}
		return m.secret, nil
	}, jwt.WithIssuer(m.issuer), jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
	if err != nil {
		return nil, config.ErrInvalidToken
	}

	if !parsed.Valid {
		return nil, config.ErrInvalidToken
	}

	return claims, nil
}
