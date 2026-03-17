package services

import (
	"context"
	"fmt"

	"github.com/HanzChrisrome/org-man-app/internal/config"
	"github.com/HanzChrisrome/org-man-app/internal/repositories"
	"github.com/HanzChrisrome/org-man-app/internal/utils"
)

type MemberService struct {
	repo      repositories.UserRepository
	hasher    utils.PasswordHasher
	validator *utils.PasswordValidator
}

func NewMemberService(repo repositories.UserRepository, hasher utils.PasswordHasher, validator *utils.PasswordValidator) *MemberService {
	return &MemberService{
		repo:      repo,
		hasher:    hasher,
		validator: validator,
	}
}

func (s *MemberService) CreateMember(ctx context.Context, req config.RegisterRequest) (*config.Member, error) {
	if err := s.validator.Validate(req.Password); err != nil {
		return nil, fmt.Errorf("%w: %v", config.ErrWeakPassword, err)
	}

	exists, err := s.repo.Exists(ctx, req.Email, req.StudentID)
	if err != nil {
		return nil, fmt.Errorf("existence check failed: %w", err)
	}
	if exists {
		return nil, config.ErrUserAlreadyExists
	}

	hPassword, err := s.hasher.HashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("password hashing failed: %w", err)
	}

	user := &config.Member{
		Name:               req.Name,
		Email:              req.Email,
		StudentID:          req.StudentID,
		PasswordHash:       hPassword,
		RegistrationStatus: config.StatusPending,
	}

	if err := s.repo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("member creation failed: %w", err)
	}

	return user, nil
}

func (s *MemberService) ListMembers(ctx context.Context, query string, status string) ([]config.MemberWithPayment, error) {
	members, err := s.repo.List(ctx, query, status)
	if err != nil {
		return nil, fmt.Errorf("service failed to list members: %w", err)
	}
	return members, nil
}

func (s *MemberService) GetMemberByID(ctx context.Context, id string) (*config.Member, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *MemberService) UpdateMember(ctx context.Context, id string, req config.UpdateMemberRequest) (*config.Member, error) {
	member, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.Name != nil {
		member.Name = *req.Name
	}
	if req.Email != nil && *req.Email != member.Email {
		// Check uniqueness if email changes
		exists, _ := s.repo.Exists(ctx, *req.Email, "")
		if exists {
			return nil, config.ErrUserAlreadyExists
		}
		member.Email = *req.Email
	}
	if req.StudentID != nil && *req.StudentID != member.StudentID {
		// Check uniqueness if student ID changes
		exists, _ := s.repo.Exists(ctx, "", *req.StudentID)
		if exists {
			return nil, config.ErrUserAlreadyExists
		}
		member.StudentID = *req.StudentID
	}
	if req.Course != nil {
		member.Course = *req.Course
	}
	if req.ContactNumber != nil {
		member.ContactNumber = *req.ContactNumber
	}
	if req.RegistrationStatus != nil {
		member.RegistrationStatus = *req.RegistrationStatus
	}

	if err := s.repo.Update(ctx, member); err != nil {
		return nil, err
	}

	return member, nil
}

func (s *MemberService) DeleteMember(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
