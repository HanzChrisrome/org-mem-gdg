package services

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	"github.com/HanzChrisrome/org-man-app/internal/config"
	"github.com/HanzChrisrome/org-man-app/internal/repositories"
	"github.com/HanzChrisrome/org-man-app/internal/utils"
)

type AuthService struct {
	repo           repositories.UserRepository
	execRepo       repositories.ExecutiveRepository
	sessionRepo    repositories.SessionRepository
	hasher         utils.PasswordHasher
	validator      *utils.PasswordValidator
	jwtManager     utils.JWTManager
	sessionManager utils.SessionManager
}

func NewAuthService(
	repo repositories.UserRepository,
	execRepo repositories.ExecutiveRepository,
	sessionRepo repositories.SessionRepository,
	hasher utils.PasswordHasher,
	validator *utils.PasswordValidator,
	jwtManager utils.JWTManager,
	sessionManager utils.SessionManager,
) *AuthService {
	return &AuthService{
		repo:           repo,
		execRepo:       execRepo,
		sessionRepo:    sessionRepo,
		hasher:         hasher,
		validator:      validator,
		jwtManager:     jwtManager,
		sessionManager: sessionManager,
	}
}

func (s *AuthService) Register(ctx context.Context, req config.RegisterRequest) (interface{}, string, error) {
	// 1. Validate password policy
	if err := s.validator.Validate(req.Password); err != nil {
		return nil, "", fmt.Errorf("%w: %v", config.ErrWeakPassword, err)
	}

	// 2. Hash password
	hPassword, err := s.hasher.HashPassword(req.Password)
	if err != nil {
		return nil, "", fmt.Errorf("password hashing failed: %w", err)
	}

	// 3. Route based on source dashboard
	switch req.SourceDashboard {
	case config.DashboardMembers:
		// Check for existence
		exists, err := s.repo.Exists(ctx, req.Email, req.StudentID)
		if err != nil {
			return nil, "", fmt.Errorf("existence check failed: %w", err)
		}
		if exists {
			return nil, "", config.ErrUserAlreadyExists
		}

		user := &config.Member{
			Name:         req.Name,
			Email:        req.Email,
			StudentID:    req.StudentID,
			PasswordHash: hPassword,
		}

		if err := s.repo.Create(ctx, user); err != nil {
			return nil, "", fmt.Errorf("user creation failed: %w", err)
		}
		return user, "member", nil

	case config.DashboardExecutives:
		// Check for existence
		exists, err := s.execRepo.Exists(ctx, req.Email, req.StudentID)
		if err != nil {
			return nil, "", fmt.Errorf("existence check failed: %w", err)
		}
		if exists {
			return nil, "", config.ErrUserAlreadyExists
		}

		exec := &config.Executive{
			Name:         req.Name,
			Email:        req.Email,
			StudentID:    req.StudentID,
			PasswordHash: hPassword,
		}

		if err := s.execRepo.Create(ctx, exec); err != nil {
			return nil, "", fmt.Errorf("executive creation failed: %w", err)
		}
		return exec, "executive", nil

	default:
		return nil, "", fmt.Errorf("%w: invalid source dashboard", config.ErrInvalidInput)
	}
}

func (s *AuthService) Login(ctx context.Context, req config.LoginRequest) (string, string, string, error) {
	// 1. Try executives table only
	var exec *config.Executive
	var err error
	if strings.Contains(req.Identifier, "@") {
		exec, err = s.execRepo.GetByEmail(ctx, req.Identifier)
	} else {
		exec, err = s.execRepo.GetByStudentID(ctx, req.Identifier)
	}

	if err != nil {
		if err == config.ErrUserNotFound {
			return "", "", "", config.ErrInvalidCredentials
		}
		return "", "", "", config.ErrInternal
	}

	// Verify login hash
	if err := s.hasher.VerifyPassword(req.Password, exec.PasswordHash); err != nil {
		return "", "", "", config.ErrInvalidCredentials
	}

	return exec.ID, "executive", "", nil
}

func (s *AuthService) LoginWithToken(ctx context.Context, req config.LoginRequest, userAgent, ip string) (string, string, *config.TokenPair, error) {
	userID, ownerType, _, err := s.Login(ctx, req)
	if err != nil {
		return "", "", nil, err
	}

	refreshToken, err := generateRandomToken(32)
	if err != nil {
		return "", "", nil, fmt.Errorf("%w: refresh token generation failed: %v", config.ErrInternal, err)
	}

	session, err := s.sessionManager.CreateSession(userID, ownerType, refreshToken, userAgent, ip)
	if err != nil {
		return "", "", nil, fmt.Errorf("%w: session creation failed: %v", config.ErrInternal, err)
	}

	if err := s.sessionRepo.Create(ctx, session); err != nil {
		return "", "", nil, fmt.Errorf("%w: session persistence failed: %v", config.ErrInternal, err)
	}

	accessToken, expiresIn, err := s.jwtManager.GenerateAccessTokenWithSession(userID, session.ID, ownerType)
	if err != nil {
		return "", "", nil, fmt.Errorf("%w: access token generation failed: %v", config.ErrInternal, err)
	}

	tokens := &config.TokenPair{
		AccessToken:  accessToken,
		RefreshToken: fmt.Sprintf("%s.%s", session.ID, refreshToken),
		TokenType:    "Bearer",
		ExpiresIn:    expiresIn,
	}

	return userID, ownerType, tokens, nil
}

func (s *AuthService) RefreshAccessToken(ctx context.Context, sessionID, refreshToken string) (*config.TokenPair, error) {
	session, err := s.sessionRepo.GetByID(ctx, sessionID)
	if err != nil {
		return nil, err
	}

	if err := s.sessionManager.ValidateSession(session, refreshToken); err != nil {
		return nil, err
	}

	// 1. Generate new refresh token
	newRefreshToken, err := generateRandomToken(32)
	if err != nil {
		return nil, fmt.Errorf("%w: refresh token generation failed: %v", config.ErrInternal, err)
	}

	// 2. Update session with new hash and extend expiry
	session.RefreshTokenHash = utils.HashToken(newRefreshToken)
	session.ExpiresAt = time.Now().UTC().Add(7 * 24 * time.Hour) // Use a constant or value from config

	// 3. Persist update in repo
	if err := s.sessionRepo.Update(ctx, session); err != nil {
		return nil, fmt.Errorf("%w: session update failed: %v", config.ErrInternal, err)
	}

	var userID string
	if session.OwnerType == "member" {
		user, err := s.repo.GetByID(ctx, session.OwnerID)
		if err != nil {
			if err == config.ErrUserNotFound {
				return nil, config.ErrInvalidCredentials
			}
			return nil, config.ErrInternal
		}
		userID = user.ID
	} else if session.OwnerType == "executive" {
		exec, err := s.execRepo.GetByID(ctx, session.OwnerID)
		if err != nil {
			if err == config.ErrUserNotFound {
				return nil, config.ErrInvalidCredentials
			}
			return nil, config.ErrInternal
		}
		userID = exec.ID
	} else {
		return nil, fmt.Errorf("unknown owner type: %s", session.OwnerType)
	}

	accessToken, expiresIn, err := s.jwtManager.GenerateAccessTokenWithSession(userID, session.ID, session.OwnerType)
	if err != nil {
		return nil, config.ErrInternal
	}

	return &config.TokenPair{
		AccessToken:  accessToken,
		RefreshToken: fmt.Sprintf("%s.%s", session.ID, newRefreshToken),
		TokenType:    "Bearer",
		ExpiresIn:    expiresIn,
	}, nil
}

func (s *AuthService) Logout(ctx context.Context, sessionID string) error {
	if sessionID == "" {
		return config.ErrInvalidInput
	}

	return s.sessionRepo.Revoke(ctx, sessionID, time.Now().UTC())
}

func generateRandomToken(size int) (string, error) {
	if size <= 0 {
		return "", config.ErrInvalidInput
	}

	buf := make([]byte, size)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}

	return hex.EncodeToString(buf), nil
}
