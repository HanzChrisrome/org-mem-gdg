package handlers

import (
	"github.com/HanzChrisrome/org-man-app/internal/config"
)

// ErrorResponse represents a generic error message
// @name ErrorResponse
type ErrorResponse struct {
	Error string `json:"error"`
}

// MessageResponse represents a generic success message
// @name MessageResponse
type MessageResponse struct {
	Message string `json:"message"`
}

// TokenPairResponse represents a pair of access and refresh tokens
// @name TokenPair
type TokenPairResponse config.TokenPair

// LoginResponse represents the successful login response
// @name LoginResponse
type LoginResponse struct {
	UserID    string             `json:"user_id"`
	OwnerType string             `json:"owner_type"`
	Token     *TokenPairResponse `json:"token"`
}

// RegisterResponse represents the successful registration response
// @name RegisterResponse
type RegisterResponse struct {
	User      interface{} `json:"user"`
	OwnerType string      `json:"owner_type"`
}

// RefreshResponse represents the successful token refresh response
// @name RefreshResponse
type RefreshResponse struct {
	Token *TokenPairResponse `json:"token"`
}

// Named types for existing config models to ensure clean names in Swagger
// These are just for documentation purposes to remove the package prefix

// MemberResponse represents a member record
// @name MemberResponse
type MemberResponse config.Member

// MemberWithPaymentResponse represents a member with payment details
// @name MemberWithPaymentResponse
type MemberWithPaymentResponse config.MemberWithPayment

// ExecutiveResponse represents an executive record
// @name ExecutiveResponse
type ExecutiveResponse config.Executive

// CreateExecutiveRequestDoc represents the request to create an executive
// @name CreateExecutiveRequest
type CreateExecutiveRequestDoc config.CreateExecutiveRequest

// UpdateExecutiveRequestDoc represents the request to update an executive
// @name UpdateExecutiveRequest
type UpdateExecutiveRequestDoc config.UpdateExecutiveRequest

// RegisterRequestDoc represents the request to register
// @name RegisterRequest
type RegisterRequestDoc config.RegisterRequest

// LoginRequestDoc represents the request to login
// @name LoginRequest
type LoginRequestDoc config.LoginRequest

// UpdateMemberRequestDoc represents the request to update a member
// @name UpdateMemberRequest
type UpdateMemberRequestDoc config.UpdateMemberRequest
