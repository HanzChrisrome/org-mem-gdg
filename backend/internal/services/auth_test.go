package services_test

import (
	"context"
	"testing"
	"time"

	"github.com/HanzChrisrome/org-man-app/internal/config"
	"github.com/HanzChrisrome/org-man-app/internal/services"
	"github.com/HanzChrisrome/org-man-app/internal/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock repos and managers for testing
type MockUserRepo struct{ mock.Mock }

func (m *MockUserRepo) Create(ctx context.Context, u *config.Member) error {
	return m.Called(ctx, u).Error(0)
}
func (m *MockUserRepo) GetByEmail(ctx context.Context, email string) (*config.Member, error) {
	args := m.Called(ctx, email)
	return args.Get(0).(*config.Member), args.Error(1)
}
func (m *MockUserRepo) GetByStudentID(ctx context.Context, id string) (*config.Member, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*config.Member), args.Error(1)
}
func (m *MockUserRepo) GetByID(ctx context.Context, id string) (*config.Member, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*config.Member), args.Error(1)
}
func (m *MockUserRepo) Exists(ctx context.Context, email, studentID string) (bool, error) {
	args := m.Called(ctx, email, studentID)
	return args.Bool(0), args.Error(1)
}
func (m *MockUserRepo) Update(ctx context.Context, u *config.Member) error {
	return m.Called(ctx, u).Error(0)
}
func (m *MockUserRepo) Delete(ctx context.Context, id string) error {
	return m.Called(ctx, id).Error(0)
}
func (m *MockUserRepo) List(ctx context.Context, query string, status string) ([]config.MemberWithPayment, error) {
	args := m.Called(ctx, query, status)
	return args.Get(0).([]config.MemberWithPayment), args.Error(1)
}

type MockExecRepo struct{ mock.Mock }

func (m *MockExecRepo) GetByEmail(ctx context.Context, email string) (*config.Executive, error) {
	args := m.Called(ctx, email)
	return args.Get(0).(*config.Executive), args.Error(1)
}
func (m *MockExecRepo) GetByStudentID(ctx context.Context, id string) (*config.Executive, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*config.Executive), args.Error(1)
}
func (m *MockExecRepo) GetByID(ctx context.Context, id string) (*config.Executive, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*config.Executive), args.Error(1)
}
func (m *MockExecRepo) Exists(ctx context.Context, email, studentID string) (bool, error) {
	args := m.Called(ctx, email, studentID)
	return args.Bool(0), args.Error(1)
}

type MockSessionRepo struct{ mock.Mock }

func (m *MockSessionRepo) Create(ctx context.Context, s *config.Session) error {
	return m.Called(ctx, s).Error(0)
}
func (m *MockSessionRepo) Update(ctx context.Context, s *config.Session) error {
	return m.Called(ctx, s).Error(0)
}
func (m *MockSessionRepo) GetByID(ctx context.Context, id string) (*config.Session, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*config.Session), args.Error(1)
}
func (m *MockSessionRepo) Revoke(ctx context.Context, id string, t time.Time) error {
	return m.Called(ctx, id, t).Error(0)
}

func TestAuthService_FullFlow(t *testing.T) {
	// Setup
	userRepo := new(MockUserRepo)
	execRepo := new(MockExecRepo)
	sessionRepo := new(MockSessionRepo)
	hasher := utils.NewBcryptHasher(10)
	validator := utils.NewPasswordValidator(8)
	jwtManager := utils.NewHMACJWTManager("test-secret", "test-issuer", 15)
	sessionManager := utils.NewDefaultSessionManager(168)

	service := services.NewAuthService(userRepo, execRepo, sessionRepo, hasher, validator, jwtManager, sessionManager)

	ctx := context.Background()
	password := "SecurePassword123!"
	hPassword, _ := hasher.HashPassword(password)

	member := &config.Member{
		ID:           "user-1",
		Email:        "test@example.com",
		PasswordHash: hPassword,
	}

	t.Run("Login and Refresh Rotation Test", func(t *testing.T) {
		// 1. Mock Login
		userRepo.On("GetByEmail", ctx, member.Email).Return(member, nil)
		sessionRepo.On("Create", ctx, mock.AnythingOfType("*config.Session")).Return(nil)

		// Login
		uid, oType, tokens, err := service.LoginWithToken(ctx, config.LoginRequest{
			Identifier: member.Email,
			Password:   password,
		}, "agent", "1.1.1.1")

		assert.NoError(t, err)
		assert.Equal(t, member.ID, uid)
		assert.Equal(t, "member", oType)
		assert.NotEmpty(t, tokens.AccessToken)
		assert.Contains(t, tokens.RefreshToken, ".")

		// 2. Mock Refresh with Rotation
		sessionID := tokens.RefreshToken[:64] // sessionID is 32 bytes hex = 64 chars
		originalSession := &config.Session{
			ID:               sessionID,
			OwnerID:          member.ID,
			OwnerType:        "member",
			RefreshTokenHash: utils.HashToken(tokens.RefreshToken[65:]),
			ExpiresAt:        time.Now().Add(1 * time.Hour),
		}

		sessionRepo.On("GetByID", ctx, sessionID).Return(originalSession, nil)
		sessionRepo.On("Update", ctx, originalSession).Return(nil)
		userRepo.On("GetByID", ctx, member.ID).Return(member, nil)

		// Refresh
		newTokens, err := service.RefreshAccessToken(ctx, sessionID, tokens.RefreshToken[65:])

		assert.NoError(t, err)
		assert.NotEqual(t, tokens.RefreshToken, newTokens.RefreshToken, "Refresh token must rotate")
		assert.NotEmpty(t, newTokens.AccessToken)
	})

	t.Run("Logout Revocation Test", func(t *testing.T) {
		sessionID := "session-to-revoke"
		sessionRepo.On("Revoke", ctx, sessionID, mock.AnythingOfType("time.Time")).Return(nil)

		err := service.Logout(ctx, sessionID)
		assert.NoError(t, err)
	})
}
