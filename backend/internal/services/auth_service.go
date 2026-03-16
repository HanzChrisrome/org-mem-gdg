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
	sessionRepo    repositories.SessionRepository
	hasher         utils.PasswordHasher
	validator      *utils.PasswordValidator
	jwtManager     utils.JWTManager
	sessionManager utils.SessionManager
}

func NewAuthService(
	repo repositories.UserRepository,
	sessionRepo repositories.SessionRepository,
	hasher utils.PasswordHasher,
	validator *utils.PasswordValidator,
	jwtManager utils.JWTManager,
	sessionManager utils.SessionManager,
) *AuthService {
	return &AuthService{
		repo:           repo,
		sessionRepo:    sessionRepo,
		hasher:         hasher,
		validator:      validator,
		jwtManager:     jwtManager,
		sessionManager: sessionManager,
	}
}

func (s *AuthService) RegisterMember(ctx context.Context, req config.RegisterRequest) (*config.Member, error) {
	// 1. Validate password policy
	if err := s.validator.Validate(req.Password); err != nil {
		return nil, fmt.Errorf("%w: %v", config.ErrWeakPassword, err)
	}

	// 2. Check for existence (avoid duplicates)
	exists, err := s.repo.Exists(ctx, req.Email, req.StudentID)
	if err != nil {
		return nil, fmt.Errorf("existence check failed: %w", err)
	}
	if exists {
		return nil, config.ErrUserAlreadyExists
	}

	// 3. Hash password
	hPassword, err := s.hasher.HashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("password hashing failed: %w", err)
	}

	// 4. Create user domain model
	user := &config.Member{
		Name:         req.Name,
		Email:        req.Email,
		StudentID:    req.StudentID,
		PasswordHash: hPassword,
	}

	// 5. Persist
	if err := s.repo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("user creation failed: %w", err)
	}

	return user, nil
}

func (s *AuthService) Login(ctx context.Context, req config.LoginRequest) (*config.Member, error) {
	// Normalize identifier (email-like check or student ID format)
	var user *config.Member
	var err error

	if strings.Contains(req.Identifier, "@") {
		user, err = s.repo.GetByEmail(ctx, req.Identifier)
	} else {
		user, err = s.repo.GetByStudentID(ctx, req.Identifier)
	}

	if err != nil {
		// Generic credential failure for both not-found and mismatch
		return nil, config.ErrInvalidCredentials
	}

	// Verify login hash
	if err := s.hasher.VerifyPassword(req.Password, user.PasswordHash); err != nil {
		return nil, config.ErrInvalidCredentials
	}

	return user, nil
}

func (s *AuthService) LoginWithToken(ctx context.Context, req config.LoginRequest, userAgent, ip string) (*config.Member, *config.TokenPair, error) {
	user, err := s.Login(ctx, req)
	if err != nil {
		return nil, nil, err
	}

	accessToken, expiresIn, err := s.jwtManager.GenerateAccessToken(user)
	if err != nil {
		return nil, nil, fmt.Errorf("%w: access token generation failed: %v", config.ErrInternal, err)
	}

	refreshToken, err := generateRandomToken(32)
	if err != nil {
		return nil, nil, fmt.Errorf("%w: refresh token generation failed: %v", config.ErrInternal, err)
	}

	session, err := s.sessionManager.CreateSession(user.ID, "member", refreshToken, userAgent, ip)
	if err != nil {
		return nil, nil, fmt.Errorf("%w: session creation failed: %v", config.ErrInternal, err)
	}

	if err := s.sessionRepo.Create(ctx, session); err != nil {
		return nil, nil, fmt.Errorf("%w: session persistence failed: %v", config.ErrInternal, err)
	}

	tokens := &config.TokenPair{
		AccessToken:  accessToken,
		RefreshToken: fmt.Sprintf("%s.%s", session.ID, refreshToken),
		TokenType:    "Bearer",
		ExpiresIn:    expiresIn,
	}

	return user, tokens, nil
}

func (s *AuthService) RefreshAccessToken(ctx context.Context, sessionID, refreshToken string) (*config.TokenPair, error) {
	session, err := s.sessionRepo.GetByID(ctx, sessionID)
	if err != nil {
		return nil, err
	}

	if err := s.sessionManager.ValidateSession(session, refreshToken); err != nil {
		return nil, err
	}

	var user *config.Member
	if session.OwnerType == "member" {
		user, err = s.repo.GetByID(ctx, session.OwnerID)
	} else {
		// Placeholder for executive lookup if needed in the future
		// user, err = s.execRepo.GetByID(ctx, session.OwnerID)
		return nil, fmt.Errorf("executive refresh not yet implemented in service")
	}

	if err != nil {
		if err == config.ErrUserNotFound {
			return nil, config.ErrInvalidCredentials
		}
		return nil, config.ErrInternal
	}

	accessToken, expiresIn, err := s.jwtManager.GenerateAccessToken(user)
	if err != nil {
		return nil, config.ErrInternal
	}

	return &config.TokenPair{
		AccessToken: accessToken,
		TokenType:   "Bearer",
		ExpiresIn:   expiresIn,
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
