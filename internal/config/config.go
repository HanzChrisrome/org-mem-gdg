package config

import (
	"errors"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserNotFound       = errors.New("user not found")
	ErrUserAlreadyExists  = errors.New("user already exists")
	ErrInternal           = errors.New("internal server error")
	ErrWeakPassword       = errors.New("password does not meet security requirements")
	ErrInvalidInput       = errors.New("invalid input data")
	ErrInvalidToken       = errors.New("invalid token")
	ErrSessionNotFound    = errors.New("session not found")
	ErrSessionExpired     = errors.New("session expired")
)

const (
	StatusPending       = "pending"
	StatusApproved      = "approved"
	StatusRejected      = "rejected"
	StatusInactive      = "inactive"
	StatusResubmit      = "resubmitted"
	DashboardMembers    = "members"
	DashboardExecutives = "executives"
)

type Member struct {
	ID                 string    `json:"member_id"`
	Name               string    `json:"name"`
	Email              string    `json:"email"`
	StudentID          string    `json:"student_id"`
	Course             string    `json:"course,omitempty"`
	ContactNumber      string    `json:"contact_number,omitempty"`
	RegistrationStatus string    `json:"registration_status"`
	PasswordHash       string    `json:"-"`
	CreatedAt          time.Time `json:"created_at"`
	LastUpdated        time.Time `json:"last_updated"`
}

type MemberWithPayment struct {
	Member
	LatestPaymentID     *int       `json:"latest_payment_id"`
	LatestPaymentStatus *string    `json:"latest_payment_status"`
	LatestSubmission    *time.Time `json:"latest_submission_date"`
	LatestApprovalDate  *time.Time `json:"latest_approval_date"`
	ApproverName        *string    `json:"approver_name"`
}

type UpdateMemberRequest struct {
	Name               *string `json:"name"`
	Email              *string `json:"email"`
	StudentID          *string `json:"student_id"`
	Course             *string `json:"course"`
	ContactNumber      *string `json:"contact_number"`
	RegistrationStatus *string `json:"registration_status"`
}

type Executive struct {
	ID           string    `json:"executive_id"`
	Name         string    `json:"name"`
	Email        string    `json:"email"`
	StudentID    string    `json:"student_id"`
	RoleID       int       `json:"role_id"`
	PasswordHash string    `json:"-"`
	CreatedAt    time.Time `json:"created_at"`
	LastUpdated  time.Time `json:"last_updated"`
}

type RegisterRequest struct {
	Name            string `json:"name"`
	Email           string `json:"email"`
	StudentID       string `json:"student_id"`
	Password        string `json:"password"`
	SourceDashboard string `json:"source_dashboard" binding:"required"` // "members" or "executives"
}

type CreateExecutiveRequest struct {
	Name      string `json:"name" binding:"required"`
	Email     string `json:"email" binding:"required,email"`
	StudentID string `json:"student_id" binding:"required"`
	Password  string `json:"password" binding:"required"`
	RoleID    *int   `json:"role_id"`
}

type UpdateExecutiveRequest struct {
	Name      *string `json:"name"`
	Email     *string `json:"email" binding:"omitempty,email"`
	StudentID *string `json:"student_id"`
	Password  *string `json:"password"`
	RoleID    *int    `json:"role_id"`
}

type LoginRequest struct {
	Identifier string `json:"identifier"` // Email or StudentID
	Password   string `json:"password"`
}

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int64  `json:"expires_in"`
}

type Session struct {
	ID               string     `json:"id"`
	OwnerID          string     `json:"owner_id"`
	OwnerType        string     `json:"owner_type"` // "member" or "executive"
	RefreshTokenHash string     `json:"-"`
	UserAgent        string     `json:"user_agent"`
	IPAddress        string     `json:"ip_address"`
	ExpiresAt        time.Time  `json:"expires_at"`
	CreatedAt        time.Time  `json:"created_at"`
	RevokedAt        *time.Time `json:"revoked_at,omitempty"`
}

type Config struct {
	DatabaseURL           string
	Port                  string
	BcryptCost            int
	MinPassLen            int
	JWTSecret             string
	JWTIssuer             string
	AccessTokenTTLMinutes int
	RefreshTokenTTLHours  int
}

func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println(".env file not found, using OS environment variables")
	}

	bcryptCostStr := os.Getenv("BCRYPT_COST")
	bcryptCost := 12
	if bcryptCostStr != "" {
		if cost, err := strconv.Atoi(bcryptCostStr); err == nil {
			bcryptCost = cost
		}
	}

	minPassLenStr := os.Getenv("MIN_PASSWORD_LENGTH")
	minPassLen := 12
	if minPassLenStr != "" {
		if length, err := strconv.Atoi(minPassLenStr); err == nil {
			minPassLen = length
		}
	}

	accessTTLStr := os.Getenv("ACCESS_TOKEN_TTL_MINUTES")
	accessTTL := 15
	if accessTTLStr != "" {
		if value, err := strconv.Atoi(accessTTLStr); err == nil {
			accessTTL = value
		}
	}

	refreshTTLStr := os.Getenv("REFRESH_TOKEN_TTL_HOURS")
	refreshTTL := 168
	if refreshTTLStr != "" {
		if value, err := strconv.Atoi(refreshTTLStr); err == nil {
			refreshTTL = value
		}
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "dev-secret-change-me"
	}

	jwtIssuer := os.Getenv("JWT_ISSUER")
	if jwtIssuer == "" {
		jwtIssuer = "org-man-app"
	}

	cfg := &Config{
		DatabaseURL:           os.Getenv("DATABASE_URL"),
		Port:                  os.Getenv("PORT"),
		BcryptCost:            bcryptCost,
		MinPassLen:            minPassLen,
		JWTSecret:             jwtSecret,
		JWTIssuer:             jwtIssuer,
		AccessTokenTTLMinutes: accessTTL,
		RefreshTokenTTLHours:  refreshTTL,
	}

	if cfg.DatabaseURL == "" {
		log.Fatal("DATABASE_URL is not set")
	}

	return cfg
}
