package service

import (
	"context"
	"fmt"
	"time"

	"github.com/providentia/digna/core_lume/internal/domain"
	"github.com/providentia/digna/core_lume/internal/repository"
)

// MemberService implements application logic for member management
type MemberService struct {
	memberRepo repository.MemberRepository
	workRepo   repository.WorkRepository
}

// NewMemberService creates a new MemberService
func NewMemberService(memberRepo repository.MemberRepository, workRepo repository.WorkRepository) *MemberService {
	return &MemberService{
		memberRepo: memberRepo,
		workRepo:   workRepo,
	}
}

// RegisterMemberRequest represents the request to register a new member
type RegisterMemberRequest struct {
	Name   string
	Email  string
	Phone  string
	CPF    string
	Role   domain.MemberRole
	Skills []string
}

// RegisterMember creates a new member in the entity
func (s *MemberService) RegisterMember(ctx context.Context, entityID string, req *RegisterMemberRequest) (*domain.Member, error) {
	// Validate request
	if req.Name == "" {
		return nil, fmt.Errorf("%w: name is required", domain.ErrMemberInvalidName)
	}
	if req.Email == "" {
		return nil, fmt.Errorf("%w: email is required", domain.ErrMemberInvalidEmail)
	}
	if !isValidRole(req.Role) {
		return nil, fmt.Errorf("%w: role must be COORDINATOR, MEMBER, or ADVISOR", domain.ErrMemberInvalidRole)
	}

	// Check if email already exists
	_, err := s.memberRepo.FindByEmail(entityID, req.Email)
	if err == nil {
		return nil, fmt.Errorf("%w: %s", domain.ErrMemberDuplicateEmail, req.Email)
	}

	now := time.Now()
	member := &domain.Member{
		ID:        generateID(),
		EntityID:  entityID,
		Name:      req.Name,
		Email:     req.Email,
		Phone:     req.Phone,
		CPF:       req.CPF,
		Role:      req.Role,
		Status:    domain.StatusActive,
		JoinedAt:  now,
		Skills:    req.Skills,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := s.memberRepo.Save(member); err != nil {
		return nil, fmt.Errorf("failed to save member: %w", err)
	}

	return member, nil
}

// UpdateMemberRequest represents the request to update a member
type UpdateMemberRequest struct {
	ID     string
	Name   string
	Email  string
	Phone  string
	CPF    string
	Role   domain.MemberRole
	Skills []string
}

// UpdateMember updates member information
func (s *MemberService) UpdateMember(ctx context.Context, entityID string, req *UpdateMemberRequest) (*domain.Member, error) {
	// Find existing member
	member, err := s.memberRepo.FindByID(entityID, req.ID)
	if err != nil {
		return nil, fmt.Errorf("member not found: %w", err)
	}

	// Validate new email if changed
	if req.Email != "" && req.Email != member.Email {
		_, err := s.memberRepo.FindByEmail(entityID, req.Email)
		if err == nil {
			return nil, fmt.Errorf("%w: %s", domain.ErrMemberDuplicateEmail, req.Email)
		}
		member.Email = req.Email
	}

	// Update fields
	if req.Name != "" {
		member.Name = req.Name
	}
	if req.Phone != "" {
		member.Phone = req.Phone
	}
	if req.CPF != "" {
		member.CPF = req.CPF
	}
	if isValidRole(req.Role) {
		member.Role = req.Role
	}
	if len(req.Skills) > 0 {
		member.Skills = req.Skills
	}

	if err := s.memberRepo.Update(member); err != nil {
		return nil, fmt.Errorf("failed to update member: %w", err)
	}

	return member, nil
}

// DeactivateMember deactivates a member
func (s *MemberService) DeactivateMember(ctx context.Context, entityID, memberID string) error {
	// Verify member exists
	_, err := s.memberRepo.FindByID(entityID, memberID)
	if err != nil {
		return fmt.Errorf("member not found: %w", err)
	}

	// Check if this is the last coordinator
	member, _ := s.memberRepo.FindByID(entityID, memberID)
	if member.IsCoordinator() {
		coordinators, err := s.memberRepo.ListByRole(entityID, domain.RoleCoordinator)
		if err != nil {
			return fmt.Errorf("failed to check coordinators: %w", err)
		}

		activeCoordinators := 0
		for _, c := range coordinators {
			if c.Status == domain.StatusActive && c.ID != memberID {
				activeCoordinators++
			}
		}

		if activeCoordinators == 0 {
			return fmt.Errorf("cannot deactivate the last active coordinator")
		}
	}

	if err := s.memberRepo.UpdateStatus(entityID, memberID, domain.StatusInactive); err != nil {
		return fmt.Errorf("failed to deactivate member: %w", err)
	}

	return nil
}

// ActivateMember activates a member
func (s *MemberService) ActivateMember(ctx context.Context, entityID, memberID string) error {
	_, err := s.memberRepo.FindByID(entityID, memberID)
	if err != nil {
		return fmt.Errorf("member not found: %w", err)
	}

	if err := s.memberRepo.UpdateStatus(entityID, memberID, domain.StatusActive); err != nil {
		return fmt.Errorf("failed to activate member: %w", err)
	}

	return nil
}

// GetMember returns a member by ID
func (s *MemberService) GetMember(ctx context.Context, entityID, memberID string) (*domain.Member, error) {
	return s.memberRepo.FindByID(entityID, memberID)
}

// GetMemberByEmail returns a member by email
func (s *MemberService) GetMemberByEmail(ctx context.Context, entityID, email string) (*domain.Member, error) {
	return s.memberRepo.FindByEmail(entityID, email)
}

// ListMembers returns all members of an entity
func (s *MemberService) ListMembers(ctx context.Context, entityID string) ([]domain.Member, error) {
	return s.memberRepo.ListByEntity(entityID)
}

// ListMembersByRole returns members filtered by role
func (s *MemberService) ListMembersByRole(ctx context.Context, entityID string, role domain.MemberRole) ([]domain.Member, error) {
	return s.memberRepo.ListByRole(entityID, role)
}

// GetMemberStats returns statistics for a member
type MemberStats struct {
	MemberID            string
	Name                string
	Email               string
	Role                domain.MemberRole
	Status              domain.MemberStatus
	TotalWorkMinutes    int64
	TotalWorkHours      float64
	NumberOfWorkLogs    int64
	ContributionPercent float64
}

// GetMemberStats returns work statistics for a member
func (s *MemberService) GetMemberStats(ctx context.Context, entityID, memberID string) (*MemberStats, error) {
	member, err := s.memberRepo.FindByID(entityID, memberID)
	if err != nil {
		return nil, fmt.Errorf("member not found: %w", err)
	}

	totalMinutes, logCount, err := s.workRepo.GetTotalByMember(entityID, memberID)
	if err != nil {
		return nil, fmt.Errorf("failed to get work stats: %w", err)
	}

	return &MemberStats{
		MemberID:         memberID,
		Name:             member.Name,
		Email:            member.Email,
		Role:             member.Role,
		Status:           member.Status,
		TotalWorkMinutes: totalMinutes,
		TotalWorkHours:   float64(totalMinutes) / 60.0,
		NumberOfWorkLogs: logCount,
	}, nil
}

// GetEntityStats returns statistics for the entire entity
type EntityStats struct {
	TotalMembers         int
	ActiveMembers        int
	Coordinators         int
	TotalWorkHours       float64
	AverageWorkPerMember float64
}

// GetEntityStats returns statistics for the entity
func (s *MemberService) GetEntityStats(ctx context.Context, entityID string) (*EntityStats, error) {
	totalCount, err := s.memberRepo.CountByEntity(entityID)
	if err != nil {
		return nil, fmt.Errorf("failed to count members: %w", err)
	}

	activeCount, err := s.memberRepo.CountActiveByEntity(entityID)
	if err != nil {
		return nil, fmt.Errorf("failed to count active members: %w", err)
	}

	coordinators, err := s.memberRepo.ListByRole(entityID, domain.RoleCoordinator)
	if err != nil {
		return nil, fmt.Errorf("failed to list coordinators: %w", err)
	}

	activeCoordinators := 0
	for _, c := range coordinators {
		if c.Status == domain.StatusActive {
			activeCoordinators++
		}
	}

	return &EntityStats{
		TotalMembers:   totalCount,
		ActiveMembers:  activeCount,
		Coordinators:   activeCoordinators,
		TotalWorkHours: 0, // Would need to aggregate all members
	}, nil
}

// ValidateCoordinatorExists checks if there's at least one coordinator
func (s *MemberService) ValidateCoordinatorExists(ctx context.Context, entityID string) error {
	coordinators, err := s.memberRepo.ListByRole(entityID, domain.RoleCoordinator)
	if err != nil {
		return fmt.Errorf("failed to list coordinators: %w", err)
	}

	for _, c := range coordinators {
		if c.Status == domain.StatusActive {
			return nil
		}
	}

	return fmt.Errorf("no active coordinator found for entity")
}

// HasCoordinatorRole checks if a member has coordinator role
func (s *MemberService) HasCoordinatorRole(ctx context.Context, entityID, memberID string) (bool, error) {
	member, err := s.memberRepo.FindByID(entityID, memberID)
	if err != nil {
		return false, fmt.Errorf("member not found: %w", err)
	}

	return member.IsCoordinator() && member.Status == domain.StatusActive, nil
}

// Helper functions
func isValidRole(role domain.MemberRole) bool {
	switch role {
	case domain.RoleCoordinator, domain.RoleMember, domain.RoleAdvisor:
		return true
	}
	return false
}

func generateID() string {
	return fmt.Sprintf("mem-%d", time.Now().UnixNano())
}
