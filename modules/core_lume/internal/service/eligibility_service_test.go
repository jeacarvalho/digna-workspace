package service

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/providentia/digna/core_lume/internal/domain"
)

// MockEligibilityRepository implements EligibilityRepository for testing
type MockEligibilityRepository struct {
	profiles map[string]*domain.EligibilityProfile
}

func NewMockEligibilityRepository() *MockEligibilityRepository {
	return &MockEligibilityRepository{
		profiles: make(map[string]*domain.EligibilityProfile),
	}
}

func (m *MockEligibilityRepository) Save(profile *domain.EligibilityProfile) error {
	m.profiles[profile.EntityID] = profile
	return nil
}

func (m *MockEligibilityRepository) FindByEntityID(entityID string) (*domain.EligibilityProfile, error) {
	profile, exists := m.profiles[entityID]
	if !exists {
		return nil, domain.ErrProfileNotFound
	}
	return profile, nil
}

func (m *MockEligibilityRepository) ListIncomplete() ([]*domain.EligibilityProfile, error) {
	var incomplete []*domain.EligibilityProfile
	for _, profile := range m.profiles {
		if !profile.IsComplete() {
			incomplete = append(incomplete, profile)
		}
	}
	return incomplete, nil
}

func (m *MockEligibilityRepository) UpdateFields(entityID string, fields map[string]interface{}) error {
	profile, exists := m.profiles[entityID]
	if !exists {
		return domain.ErrProfileNotFound
	}
	// Simple field update for testing
	if val, ok := fields["valor_necessario"]; ok {
		if v, ok := val.(int64); ok {
			profile.ValorNecessario = v
		}
	}
	if val, ok := fields["cnpj"]; ok {
		if v, ok := val.(string); ok {
			profile.CNPJ = v
		}
	}
	return nil
}

func (m *MockEligibilityRepository) InitTable(entityID string) error {
	return nil
}

// EligibilityMockMemberRepository implements MemberRepository for testing
type EligibilityMockMemberRepository struct {
	members map[string]*domain.Member
}

func NewEligibilityMockMemberRepository() *EligibilityMockMemberRepository {
	return &EligibilityMockMemberRepository{
		members: make(map[string]*domain.Member),
	}
}

func (m *EligibilityMockMemberRepository) Save(member *domain.Member) error {
	m.members[member.ID] = member
	return nil
}

func (m *EligibilityMockMemberRepository) FindByID(entityID, memberID string) (*domain.Member, error) {
	member, exists := m.members[memberID]
	if !exists {
		return nil, fmt.Errorf("member not found")
	}
	return member, nil
}

func (m *EligibilityMockMemberRepository) FindByEmail(entityID, email string) (*domain.Member, error) {
	return nil, fmt.Errorf("not implemented")
}

func (m *EligibilityMockMemberRepository) ListByEntity(entityID string) ([]domain.Member, error) {
	return nil, fmt.Errorf("not implemented")
}

func (m *EligibilityMockMemberRepository) ListByRole(entityID string, role domain.MemberRole) ([]domain.Member, error) {
	return nil, fmt.Errorf("not implemented")
}

func (m *EligibilityMockMemberRepository) Update(member *domain.Member) error {
	return fmt.Errorf("not implemented")
}

func (m *EligibilityMockMemberRepository) UpdateStatus(entityID, memberID string, status domain.MemberStatus) error {
	return fmt.Errorf("not implemented")
}

func (m *EligibilityMockMemberRepository) CountByEntity(entityID string) (int, error) {
	return 0, fmt.Errorf("not implemented")
}

func (m *EligibilityMockMemberRepository) CountActiveByEntity(entityID string) (int, error) {
	return 0, fmt.Errorf("not implemented")
}

func TestEligibilityService_CreateOrUpdate_NewProfile(t *testing.T) {
	ctx := context.Background()
	eligRepo := NewMockEligibilityRepository()
	memberRepo := NewEligibilityMockMemberRepository()

	// Add coordinator member
	memberRepo.members["user-1"] = &domain.Member{
		ID:       "user-1",
		EntityID: "entity-1",
		Role:     domain.RoleCoordinator,
		Status:   domain.StatusActive,
	}

	service := NewEligibilityService(eligRepo, memberRepo)

	trueVal := true
	finalidade := string(domain.FinalidadeCapitalGiro)
	tipo := string(domain.TipoEntidadeMEI)
	valor := int64(50000)

	input := domain.EligibilityInput{
		InscritoCadUnico:  &trueVal,
		FinalidadeCredito: &finalidade,
		TipoEntidade:      &tipo,
		ValorNecessario:   &valor,
	}

	profile, err := service.CreateOrUpdate(ctx, "entity-1", "user-1", input)
	if err != nil {
		t.Fatalf("CreateOrUpdate() unexpected error: %v", err)
	}

	if profile == nil {
		t.Fatal("CreateOrUpdate() returned nil profile")
	}

	if profile.EntityID != "entity-1" {
		t.Errorf("EntityID = %s, expected entity-1", profile.EntityID)
	}

	if !profile.InscritoCadUnico {
		t.Error("InscritoCadUnico should be true")
	}

	if profile.FinalidadeCredito != domain.FinalidadeCapitalGiro {
		t.Errorf("FinalidadeCredito = %v, expected CAPITAL_GIRO", profile.FinalidadeCredito)
	}

	if profile.PreenchidoPor != "user-1" {
		t.Errorf("PreenchidoPor = %s, expected user-1", profile.PreenchidoPor)
	}
}

func TestEligibilityService_CreateOrUpdate_ExistingProfile(t *testing.T) {
	ctx := context.Background()
	eligRepo := NewMockEligibilityRepository()
	memberRepo := NewEligibilityMockMemberRepository()

	// Add coordinator member
	memberRepo.members["user-1"] = &domain.Member{
		ID:       "user-1",
		EntityID: "entity-1",
		Role:     domain.RoleCoordinator,
		Status:   domain.StatusActive,
	}

	// Create existing profile
	now := time.Now().Unix()
	existingProfile := &domain.EligibilityProfile{
		ID:                "profile-1",
		EntityID:          "entity-1",
		FinalidadeCredito: domain.FinalidadeCapitalGiro,
		TipoEntidade:      domain.TipoEntidadeMEI,
		ValorNecessario:   50000,
		PreenchidoEm:      now,
		CreatedAt:         now,
	}
	eligRepo.Save(existingProfile)

	service := NewEligibilityService(eligRepo, memberRepo)

	// Update with new value
	newValor := int64(75000)
	input := domain.EligibilityInput{
		ValorNecessario: &newValor,
	}

	profile, err := service.CreateOrUpdate(ctx, "entity-1", "user-1", input)
	if err != nil {
		t.Fatalf("CreateOrUpdate() unexpected error: %v", err)
	}

	if profile.ValorNecessario != 75000 {
		t.Errorf("ValorNecessario = %d, expected 75000", profile.ValorNecessario)
	}

	// Original values should be preserved
	if profile.FinalidadeCredito != domain.FinalidadeCapitalGiro {
		t.Error("FinalidadeCredito should be preserved")
	}
}

func TestEligibilityService_CreateOrUpdate_NoPermission(t *testing.T) {
	ctx := context.Background()
	eligRepo := NewMockEligibilityRepository()
	memberRepo := NewEligibilityMockMemberRepository()

	// Add coordinator member
	memberRepo.members["user-1"] = &domain.Member{
		ID:       "user-1",
		EntityID: "entity-1",
		Role:     domain.RoleMember,
		Status:   domain.StatusActive,
	}

	service := NewEligibilityService(eligRepo, memberRepo)

	trueVal := true
	input := domain.EligibilityInput{
		InscritoCadUnico: &trueVal,
	}

	_, err := service.CreateOrUpdate(ctx, "entity-1", "user-1", input)
	if err == nil {
		t.Error("CreateOrUpdate() expected error for non-coordinator but got nil")
	}
}

func TestEligibilityService_GetProfile(t *testing.T) {
	ctx := context.Background()
	eligRepo := NewMockEligibilityRepository()
	memberRepo := NewEligibilityMockMemberRepository()

	// Create profile
	now := time.Now().Unix()
	profile := &domain.EligibilityProfile{
		ID:                "profile-1",
		EntityID:          "entity-1",
		FinalidadeCredito: domain.FinalidadeCapitalGiro,
		TipoEntidade:      domain.TipoEntidadeMEI,
		ValorNecessario:   50000,
		CreatedAt:         now,
	}
	eligRepo.Save(profile)

	service := NewEligibilityService(eligRepo, memberRepo)

	found, err := service.GetProfile(ctx, "entity-1")
	if err != nil {
		t.Fatalf("GetProfile() unexpected error: %v", err)
	}

	if found == nil {
		t.Fatal("GetProfile() returned nil")
	}

	if found.EntityID != "entity-1" {
		t.Errorf("EntityID = %s, expected entity-1", found.EntityID)
	}
}

func TestEligibilityService_GetProfile_NotFound(t *testing.T) {
	ctx := context.Background()
	eligRepo := NewMockEligibilityRepository()
	memberRepo := NewEligibilityMockMemberRepository()

	service := NewEligibilityService(eligRepo, memberRepo)

	_, err := service.GetProfile(ctx, "non-existent")
	if err != domain.ErrProfileNotFound {
		t.Errorf("GetProfile() error = %v, expected ErrProfileNotFound", err)
	}
}

func TestEligibilityService_GetCompletionStatus(t *testing.T) {
	ctx := context.Background()
	eligRepo := NewMockEligibilityRepository()
	memberRepo := NewEligibilityMockMemberRepository()

	service := NewEligibilityService(eligRepo, memberRepo)

	// Test with no profile
	percent, err := service.GetCompletionStatus(ctx, "entity-1")
	if err != nil {
		t.Fatalf("GetCompletionStatus() unexpected error: %v", err)
	}
	if percent != 0.0 {
		t.Errorf("GetCompletionStatus() = %f, expected 0.0", percent)
	}

	// Create partial profile
	profile := &domain.EligibilityProfile{
		ID:       "profile-1",
		EntityID: "entity-1",
		// Only some fields filled
		InscritoCadUnico: true,
		SocioMulher:      true,
	}
	eligRepo.Save(profile)

	percent, err = service.GetCompletionStatus(ctx, "entity-1")
	if err != nil {
		t.Fatalf("GetCompletionStatus() unexpected error: %v", err)
	}

	// 2 out of 7 fields filled = ~28.57%
	if percent < 25.0 || percent > 30.0 {
		t.Errorf("GetCompletionStatus() = %f, expected ~28.57", percent)
	}
}

func TestEligibilityService_EnsureProfileExists(t *testing.T) {
	ctx := context.Background()
	eligRepo := NewMockEligibilityRepository()
	memberRepo := NewEligibilityMockMemberRepository()

	service := NewEligibilityService(eligRepo, memberRepo)

	// Ensure profile for non-existing entity
	profile, err := service.EnsureProfileExists(ctx, "entity-1", "user-1")
	if err != nil {
		t.Fatalf("EnsureProfileExists() unexpected error: %v", err)
	}

	if profile == nil {
		t.Fatal("EnsureProfileExists() returned nil")
	}

	if profile.EntityID != "entity-1" {
		t.Errorf("EntityID = %s, expected entity-1", profile.EntityID)
	}

	// Call again - should return existing profile
	profile2, err := service.EnsureProfileExists(ctx, "entity-1", "user-1")
	if err != nil {
		t.Fatalf("EnsureProfileExists() second call error: %v", err)
	}

	if profile2.ID != profile.ID {
		t.Error("EnsureProfileExists() created new profile instead of returning existing")
	}
}

func TestEligibilityService_GetOrCreateProfile(t *testing.T) {
	ctx := context.Background()
	eligRepo := NewMockEligibilityRepository()
	memberRepo := NewEligibilityMockMemberRepository()

	service := NewEligibilityService(eligRepo, memberRepo)

	// GetOrCreate for new entity
	profile, err := service.GetOrCreateProfile(ctx, "entity-1", "user-1")
	if err != nil {
		t.Fatalf("GetOrCreateProfile() unexpected error: %v", err)
	}

	if profile == nil {
		t.Fatal("GetOrCreateProfile() returned nil")
	}

	// Should create empty profile
	if profile.EntityID != "entity-1" {
		t.Errorf("EntityID = %s, expected entity-1", profile.EntityID)
	}
}

func TestEligibilityService_CanUserEditProfile(t *testing.T) {
	ctx := context.Background()
	eligRepo := NewMockEligibilityRepository()
	memberRepo := NewEligibilityMockMemberRepository()

	// Add coordinator
	memberRepo.members["coord-1"] = &domain.Member{
		ID:       "coord-1",
		EntityID: "entity-1",
		Role:     domain.RoleCoordinator,
		Status:   domain.StatusActive,
	}

	// Add regular member
	memberRepo.members["member-1"] = &domain.Member{
		ID:       "member-1",
		EntityID: "entity-1",
		Role:     domain.RoleMember,
		Status:   domain.StatusActive,
	}

	service := NewEligibilityService(eligRepo, memberRepo)

	// Test coordinator can edit
	canEdit, err := service.CanUserEditProfile(ctx, "entity-1", "coord-1")
	if err != nil {
		t.Fatalf("CanUserEditProfile() error: %v", err)
	}
	if !canEdit {
		t.Error("Coordinator should be able to edit")
	}

	// Test member cannot edit
	canEdit, err = service.CanUserEditProfile(ctx, "entity-1", "member-1")
	if err != nil {
		t.Fatalf("CanUserEditProfile() error: %v", err)
	}
	if canEdit {
		t.Error("Member should not be able to edit")
	}
}

func TestEligibilityService_CanUserEditProfile_InactiveCoordinator(t *testing.T) {
	ctx := context.Background()
	eligRepo := NewMockEligibilityRepository()
	memberRepo := NewEligibilityMockMemberRepository()

	// Add inactive coordinator
	memberRepo.members["coord-1"] = &domain.Member{
		ID:       "coord-1",
		EntityID: "entity-1",
		Role:     domain.RoleCoordinator,
		Status:   domain.StatusInactive, // Inactive!
	}

	service := NewEligibilityService(eligRepo, memberRepo)

	canEdit, err := service.CanUserEditProfile(ctx, "entity-1", "coord-1")
	if err != nil {
		t.Fatalf("CanUserEditProfile() error: %v", err)
	}
	if canEdit {
		t.Error("Inactive coordinator should not be able to edit")
	}
}
