package services

import (
	"context"
	"fmt"

	"github.com/HanzChrisrome/org-man-app/internal/config"
	"github.com/HanzChrisrome/org-man-app/internal/repositories"
	"github.com/HanzChrisrome/org-man-app/internal/utils"
)

type ExecutiveService struct {
	repo      repositories.ExecutiveRepository
	hasher    utils.PasswordHasher
	validator *utils.PasswordValidator
}

func NewExecutiveService(repo repositories.ExecutiveRepository, hasher utils.PasswordHasher, validator *utils.PasswordValidator) *ExecutiveService {
	return &ExecutiveService{
		repo:      repo,
		hasher:    hasher,
		validator: validator,
	}
}

func (s *ExecutiveService) CreateExecutive(ctx context.Context, req config.CreateExecutiveRequest) (*config.Executive, error) {
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

	exec := &config.Executive{
		Name:         req.Name,
		Email:        req.Email,
		StudentID:    req.StudentID,
		PasswordHash: hPassword,
	}
	if req.RoleID != nil {
		exec.RoleID = *req.RoleID
	}

	if err := s.repo.Create(ctx, exec); err != nil {
		return nil, fmt.Errorf("executive creation failed: %w", err)
	}

	return exec, nil
}

func (s *ExecutiveService) ListExecutives(ctx context.Context) ([]config.Executive, error) {
	return s.repo.List(ctx)
}

func (s *ExecutiveService) GetExecutiveByID(ctx context.Context, id string) (*config.Executive, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *ExecutiveService) UpdateExecutive(ctx context.Context, id string, req config.UpdateExecutiveRequest) (*config.Executive, error) {
	exec, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.Name != nil {
		exec.Name = *req.Name
	}
	if req.Email != nil && *req.Email != exec.Email {
		exists, _ := s.repo.Exists(ctx, *req.Email, "")
		if exists {
			return nil, config.ErrUserAlreadyExists
		}
		exec.Email = *req.Email
	}
	if req.StudentID != nil && *req.StudentID != exec.StudentID {
		exists, _ := s.repo.Exists(ctx, "", *req.StudentID)
		if exists {
			return nil, config.ErrUserAlreadyExists
		}
		exec.StudentID = *req.StudentID
	}
	if req.Password != nil {
		if err := s.validator.Validate(*req.Password); err != nil {
			return nil, fmt.Errorf("%w: %v", config.ErrWeakPassword, err)
		}

		hPassword, err := s.hasher.HashPassword(*req.Password)
		if err != nil {
			return nil, fmt.Errorf("password hashing failed: %w", err)
		}
		exec.PasswordHash = hPassword
	}
	if req.RoleID != nil {
		exec.RoleID = *req.RoleID
	}

	if err := s.repo.Update(ctx, exec); err != nil {
		return nil, err
	}

	return exec, nil
}

func (s *ExecutiveService) DeleteExecutive(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
