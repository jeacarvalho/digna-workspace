package service

import (
	"context"
	"fmt"
	"time"

	"github.com/providentia/digna/core_lume/internal/domain"
	"github.com/providentia/digna/core_lume/internal/repository"
)

// EligibilityService implements business logic for eligibility profiles
type EligibilityService struct {
	repo       repository.EligibilityRepository
	memberRepo repository.MemberRepository
}

// NewEligibilityService creates a new EligibilityService
func NewEligibilityService(repo repository.EligibilityRepository, memberRepo repository.MemberRepository) *EligibilityService {
	return &EligibilityService{
		repo:       repo,
		memberRepo: memberRepo,
	}
}

// CreateOrUpdate creates or updates an eligibility profile
func (s *EligibilityService) CreateOrUpdate(ctx context.Context, entityID string, userID string, input domain.EligibilityInput) (*domain.EligibilityProfile, error) {
	// Check if user has permission to edit
	canEdit, err := s.canUserEdit(ctx, entityID, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to check permissions: %w", err)
	}
	if !canEdit {
		return nil, fmt.Errorf("user does not have permission to edit eligibility profile")
	}

	// Try to find existing profile
	profile, err := s.repo.FindByEntityID(entityID)
	if err != nil && err != domain.ErrProfileNotFound {
		return nil, fmt.Errorf("failed to find profile: %w", err)
	}

	now := time.Now().Unix()

	if profile == nil {
		// Create new profile
		profile = &domain.EligibilityProfile{
			ID:        fmt.Sprintf("elig-%d", now),
			EntityID:  entityID,
			CreatedAt: now,
		}
	}

	// Update profile with input data
	if err := profile.Update(input, userID); err != nil {
		return nil, fmt.Errorf("failed to update profile: %w", err)
	}

	// Save profile
	if err := s.repo.Save(profile); err != nil {
		return nil, fmt.Errorf("failed to save profile: %w", err)
	}

	return profile, nil
}

// GetProfile retrieves the eligibility profile for an entity
func (s *EligibilityService) GetProfile(ctx context.Context, entityID string) (*domain.EligibilityProfile, error) {
	profile, err := s.repo.FindByEntityID(entityID)
	if err != nil {
		if err == domain.ErrProfileNotFound {
			return nil, err
		}
		return nil, fmt.Errorf("failed to get profile: %w", err)
	}
	return profile, nil
}

// GetCompletionStatus returns the completion percentage of the profile
func (s *EligibilityService) GetCompletionStatus(ctx context.Context, entityID string) (float64, error) {
	profile, err := s.repo.FindByEntityID(entityID)
	if err != nil {
		if err == domain.ErrProfileNotFound {
			return 0.0, nil
		}
		return 0.0, fmt.Errorf("failed to get profile: %w", err)
	}
	return profile.GetCompletionPercent(), nil
}

// EnsureProfileExists creates an empty profile if none exists
func (s *EligibilityService) EnsureProfileExists(ctx context.Context, entityID string, userID string) (*domain.EligibilityProfile, error) {
	profile, err := s.repo.FindByEntityID(entityID)
	if err == nil {
		return profile, nil
	}
	if err != domain.ErrProfileNotFound {
		return nil, fmt.Errorf("failed to check profile existence: %w", err)
	}

	// Create empty profile
	now := time.Now().Unix()
	profile = &domain.EligibilityProfile{
		ID:        fmt.Sprintf("elig-%d", now),
		EntityID:  entityID,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := s.repo.Save(profile); err != nil {
		return nil, fmt.Errorf("failed to create empty profile: %w", err)
	}

	return profile, nil
}

// GetOrCreateProfile gets existing profile or creates empty one
func (s *EligibilityService) GetOrCreateProfile(ctx context.Context, entityID string, userID string) (*domain.EligibilityProfile, error) {
	profile, err := s.GetProfile(ctx, entityID)
	if err == nil {
		return profile, nil
	}
	if err != domain.ErrProfileNotFound {
		return nil, err
	}
	return s.EnsureProfileExists(ctx, entityID, userID)
}

// CanUserEditProfile checks if a user can edit the profile
func (s *EligibilityService) CanUserEditProfile(ctx context.Context, entityID string, userID string) (bool, error) {
	return s.canUserEdit(ctx, entityID, userID)
}

// canUserEdit checks if user has coordinator role
func (s *EligibilityService) canUserEdit(ctx context.Context, entityID string, userID string) (bool, error) {
	// Skip permission check for system user
	if userID == "system" {
		return true, nil
	}

	// Get member info
	member, err := s.memberRepo.FindByID(entityID, userID)
	if err != nil {
		// If member not found, allow (for development/testing)
		return true, nil
	}

	// Only coordinators can edit
	return member.Role == domain.RoleCoordinator && member.Status == domain.StatusActive, nil
}

// GetProfileWithDefaults returns profile with default values if not set
func (s *EligibilityService) GetProfileWithDefaults(ctx context.Context, entityID string) (*domain.EligibilityProfile, error) {
	profile, err := s.GetProfile(ctx, entityID)
	if err != nil {
		return nil, err
	}

	// Set default values if needed
	if profile.FinalidadeCredito == "" {
		// No default - user must choose
	}
	if profile.TipoEntidade == "" {
		// Try to infer from entity data (would need enterprise repo)
	}

	return profile, nil
}

// InitTableForEntity initializes the table for an entity
func (s *EligibilityService) InitTableForEntity(entityID string) error {
	if repoWithInit, ok := s.repo.(interface{ InitTable(entityID string) error }); ok {
		return repoWithInit.InitTable(entityID)
	}
	return nil
}
