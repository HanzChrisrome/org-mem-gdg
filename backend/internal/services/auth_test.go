package services_test

import (
	"context"
	"strings"
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

func (m *MockExecRepo) Create(ctx context.Context, e *config.Executive) error {
	return m.Called(ctx, e).Error(0)
}
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
func (m *MockExecRepo) List(ctx context.Context) ([]config.Executive, error) {
	args := m.Called(ctx)
	return args.Get(0).([]config.Executive), args.Error(1)
}
func (m *MockExecRepo) Update(ctx context.Context, exec *config.Executive) error {
	return m.Called(ctx, exec).Error(0)
}
func (m *MockExecRepo) Delete(ctx context.Context, id string) error {
	return m.Called(ctx, id).Error(0)
}

type MockSessionRepo struct{ mock.Mock }

func (m *MockSessionRepo) Create(ctx context.Context, s *config.Session) error {
	return m.Called(ctx, s).Error(0)
}
func (m *MockSessionRepo) Update(ctx context.Context, s *config.Session) error {
	return m.Called(ctx, s).Error(0)
}
func (m *MockSessionRepo) GetByRefreshTokenID(ctx context.Context, id string) (*config.Session, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*config.Session), args.Error(1)
}
func (m *MockSessionRepo) RevokeByRefreshTokenID(ctx context.Context, id string, t time.Time) error {
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
	sessionManager := utils.NewDefaultSessionManager(1, 24)

	service := services.NewAuthService(userRepo, execRepo, sessionRepo, hasher, validator, jwtManager, sessionManager)

	ctx := context.Background()
	password := "SecurePassword123!"
	hPassword, _ := hasher.HashPassword(password)

	t.Run("Login and Refresh Rotation Test", func(t *testing.T) {
		// Create an executive for login (members no longer supported)
		exec := &config.Executive{
			ID:           "exec-1",
			Email:        "exec@example.com",
			PasswordHash: hPassword,
		}

		// 1. Mock Login
		execRepo.On("GetByEmail", ctx, exec.Email).Return(exec, nil)
		sessionRepo.On("Create", ctx, mock.AnythingOfType("*config.Session")).Return(nil)

		// Login
		uid, oType, tokens, err := service.LoginWithToken(ctx, config.LoginRequest{
			Identifier: exec.Email,
			Password:   password,
		}, "agent", "1.1.1.1")

		assert.NoError(t, err)
		assert.Equal(t, exec.ID, uid)
		assert.Equal(t, "executive", oType)
		assert.NotEmpty(t, tokens.AccessToken)
		assert.Contains(t, tokens.RefreshToken, ".")

		// 2. Mock Refresh with Rotation
		parts := strings.Split(tokens.RefreshToken, ".")
		refreshTokenID := parts[0]
		refreshTokenSecret := parts[1]

		originalSession := &config.Session{
			ID:               "550e8400-e29b-41d4-a716-446655440000",
			RefreshTokenID:   refreshTokenID,
			OwnerID:          exec.ID,
			OwnerType:        "executive",
			RefreshTokenHash: utils.HashToken(refreshTokenSecret),
			ExpiresAt:        time.Now().Add(1 * time.Hour),
		}

		sessionRepo.On("GetByRefreshTokenID", ctx, refreshTokenID).Return(originalSession, nil)
		sessionRepo.On("Update", ctx, originalSession).Return(nil)
		execRepo.On("GetByID", ctx, exec.ID).Return(exec, nil)

		// Refresh
		newTokens, err := service.RefreshAccessToken(ctx, refreshTokenID, refreshTokenSecret)

		assert.NoError(t, err)
		assert.NotEqual(t, tokens.RefreshToken, newTokens.RefreshToken, "Refresh token must rotate")
		assert.NotEmpty(t, newTokens.AccessToken)
	})

	t.Run("Member Login Rejection Test", func(t *testing.T) {
		execRepo.On("GetByEmail", ctx, "member@example.com").Return((*config.Executive)(nil), config.ErrUserNotFound)

		_, _, _, err := service.Login(ctx, config.LoginRequest{
			Identifier: "member@example.com",
			Password:   password,
		})

		assert.Error(t, err)
		assert.Equal(t, config.ErrInvalidCredentials, err)
	})

	t.Run("Dashboard-Routed Registration Test", func(t *testing.T) {
		// Member registration
		reqMember := config.RegisterRequest{
			Name:            "Member User",
			Email:           "new-member@example.com",
			StudentID:       "MEM-001",
			Password:        password,
			SourceDashboard: config.DashboardMembers,
		}

		userRepo.On("Exists", ctx, reqMember.Email, reqMember.StudentID).Return(false, nil)
		userRepo.On("Create", ctx, mock.AnythingOfType("*config.Member")).Return(nil)

		resMember, ownerType, err := service.Register(ctx, reqMember)
		assert.NoError(t, err)
		assert.Equal(t, "member", ownerType)
		assert.IsType(t, &config.Member{}, resMember)

		// Executive registration
		reqExec := config.RegisterRequest{
			Name:            "Exec User",
			Email:           "new-exec@example.com",
			StudentID:       "EXEC-001",
			Password:        password,
			SourceDashboard: config.DashboardExecutives,
		}

		execRepo.On("Exists", ctx, reqExec.Email, reqExec.StudentID).Return(false, nil)
		execRepo.On("Create", ctx, mock.AnythingOfType("*config.Executive")).Return(nil)

		resExec, ownerType, err := service.Register(ctx, reqExec)
		assert.NoError(t, err)
		assert.Equal(t, "executive", ownerType)
		assert.IsType(t, &config.Executive{}, resExec)
	})

	t.Run("Logout Test Cases", func(t *testing.T) {
		refreshTokenID := "test-sid"
		refreshTokenSecret := "test-secret"
		revokedAt := mock.AnythingOfType("time.Time")

		// 1. Standard logout
		sessionRepo.On("RevokeByRefreshTokenID", ctx, refreshTokenID, revokedAt).Return(nil).Once()
		err := service.Logout(ctx, refreshTokenID)
		assert.NoError(t, err)

		// 2. Refresh token logout
		session := &config.Session{
			RefreshTokenID:   refreshTokenID,
			RefreshTokenHash: utils.HashToken(refreshTokenSecret),
			ExpiresAt:        time.Now().Add(1 * time.Hour),
		}
		sessionRepo.On("GetByRefreshTokenID", ctx, refreshTokenID).Return(session, nil).Once()
		sessionRepo.On("RevokeByRefreshTokenID", ctx, refreshTokenID, revokedAt).Return(nil).Once()

		err = service.LogoutWithRefreshToken(ctx, refreshTokenID, refreshTokenSecret)
		assert.NoError(t, err)

		// 3. Refresh token logout - Invalid secret
		sessionRepo.On("GetByRefreshTokenID", ctx, refreshTokenID).Return(session, nil).Once()
		err = service.LogoutWithRefreshToken(ctx, refreshTokenID, "wrong-secret")
		assert.Error(t, err)
		assert.Equal(t, config.ErrInvalidToken, err)
	})
}
