package service

import (
	"context"
	"testing"
	"time"

	"github.com/providentia/digna/core_lume/internal/domain"
)

// MockMemberRepository implements MemberRepository for testing
type MockMemberRepository struct {
	members map[string]*domain.Member
}

func NewMockMemberRepository() *MockMemberRepository {
	return &MockMemberRepository{
		members: make(map[string]*domain.Member),
	}
}

func (m *MockMemberRepository) Save(member *domain.Member) error {
	m.members[member.ID] = member
	return nil
}

func (m *MockMemberRepository) FindByID(entityID, memberID string) (*domain.Member, error) {
	member, exists := m.members[memberID]
	if !exists {
		return nil, domain.ErrMemberInvalidName // Using existing error
	}
	return member, nil
}

func (m *MockMemberRepository) FindByEmail(entityID, email string) (*domain.Member, error) {
	for _, member := range m.members {
		if member.Email == email {
			return member, nil
		}
	}
	return nil, domain.ErrMemberInvalidEmail
}

func (m *MockMemberRepository) ListByEntity(entityID string) ([]domain.Member, error) {
	var result []domain.Member
	for _, member := range m.members {
		if member.EntityID == entityID {
			result = append(result, *member)
		}
	}
	return result, nil
}

func (m *MockMemberRepository) ListByRole(entityID string, role domain.MemberRole) ([]domain.Member, error) {
	var result []domain.Member
	for _, member := range m.members {
		if member.EntityID == entityID && member.Role == role {
			result = append(result, *member)
		}
	}
	return result, nil
}

func (m *MockMemberRepository) Update(member *domain.Member) error {
	m.members[member.ID] = member
	return nil
}

func (m *MockMemberRepository) UpdateStatus(entityID, memberID string, status domain.MemberStatus) error {
	if member, exists := m.members[memberID]; exists {
		member.Status = status
	}
	return nil
}

func (m *MockMemberRepository) CountByEntity(entityID string) (int, error) {
	count := 0
	for _, member := range m.members {
		if member.EntityID == entityID {
			count++
		}
	}
	return count, nil
}

func (m *MockMemberRepository) CountActiveByEntity(entityID string) (int, error) {
	count := 0
	for _, member := range m.members {
		if member.EntityID == entityID && member.Status == domain.StatusActive {
			count++
		}
	}
	return count, nil
}

// MockWorkRepository implements WorkRepository for testing
type MockWorkRepository struct {
	workLogs map[string][]domain.WorkLog
}

func NewMockWorkRepository() *MockWorkRepository {
	return &MockWorkRepository{
		workLogs: make(map[string][]domain.WorkLog),
	}
}

func (m *MockWorkRepository) Save(work *domain.WorkLog) (int64, error) {
	return 1, nil
}

func (m *MockWorkRepository) GetTotalByMember(entityID, memberID string) (int64, int64, error) {
	return 480, 10, nil // 480 minutes, 10 logs
}

func (m *MockWorkRepository) GetAllMembersWork(entityID string) (map[string]int64, error) {
	return make(map[string]int64), nil
}

func (m *MockWorkRepository) GetWorkByPeriod(entityID string, startDate, endDate time.Time) ([]domain.WorkLog, error) {
	return []domain.WorkLog{}, nil
}

func TestMemberService_RegisterMember(t *testing.T) {
	memberRepo := NewMockMemberRepository()
	workRepo := NewMockWorkRepository()
	service := NewMemberService(memberRepo, workRepo)

	ctx := context.Background()
	entityID := "test-entity-001"

	req := &RegisterMemberRequest{
		Name:   "Maria Silva",
		Email:  "maria@coop.br",
		Phone:  "+55 11 98765-4321",
		Role:   domain.RoleCoordinator,
		Skills: []string{"mel", "gestão"},
	}

	member, err := service.RegisterMember(ctx, entityID, req)
	if err != nil {
		t.Fatalf("failed to register member: %v", err)
	}

	if member.Name != req.Name {
		t.Errorf("expected name %s, got %s", req.Name, member.Name)
	}
	if member.Email != req.Email {
		t.Errorf("expected email %s, got %s", req.Email, member.Email)
	}
	if member.Role != req.Role {
		t.Errorf("expected role %s, got %s", req.Role, member.Role)
	}
	if member.Status != domain.StatusActive {
		t.Errorf("expected status ACTIVE, got %s", member.Status)
	}
}

func TestMemberService_RegisterMember_DuplicateEmail(t *testing.T) {
	memberRepo := NewMockMemberRepository()
	workRepo := NewMockWorkRepository()
	service := NewMemberService(memberRepo, workRepo)

	ctx := context.Background()
	entityID := "test-entity-001"

	// Register first member
	req1 := &RegisterMemberRequest{
		Name:  "Maria Silva",
		Email: "maria@coop.br",
		Role:  domain.RoleMember,
	}
	_, err := service.RegisterMember(ctx, entityID, req1)
	if err != nil {
		t.Fatalf("failed to register first member: %v", err)
	}

	// Try to register with same email
	req2 := &RegisterMemberRequest{
		Name:  "Maria Silva 2",
		Email: "maria@coop.br",
		Role:  domain.RoleMember,
	}
	_, err = service.RegisterMember(ctx, entityID, req2)
	if err == nil {
		t.Error("expected error for duplicate email, got nil")
	}
}

func TestMemberService_RegisterMember_InvalidData(t *testing.T) {
	memberRepo := NewMockMemberRepository()
	workRepo := NewMockWorkRepository()
	service := NewMemberService(memberRepo, workRepo)

	ctx := context.Background()
	entityID := "test-entity-001"

	// Test empty name
	req := &RegisterMemberRequest{
		Name:  "",
		Email: "test@coop.br",
		Role:  domain.RoleMember,
	}
	_, err := service.RegisterMember(ctx, entityID, req)
	if err == nil {
		t.Error("expected error for empty name, got nil")
	}

	// Test empty email
	req2 := &RegisterMemberRequest{
		Name:  "Test",
		Email: "",
		Role:  domain.RoleMember,
	}
	_, err = service.RegisterMember(ctx, entityID, req2)
	if err == nil {
		t.Error("expected error for empty email, got nil")
	}
}

func TestMemberService_DeactivateMember(t *testing.T) {
	memberRepo := NewMockMemberRepository()
	workRepo := NewMockWorkRepository()
	service := NewMemberService(memberRepo, workRepo)

	ctx := context.Background()
	entityID := "test-entity-001"

	// Register coordinator
	req := &RegisterMemberRequest{
		Name:  "Coordenador",
		Email: "coord@coop.br",
		Role:  domain.RoleCoordinator,
	}
	member, _ := service.RegisterMember(ctx, entityID, req)

	// Register another coordinator
	req2 := &RegisterMemberRequest{
		Name:  "Coordenador 2",
		Email: "coord2@coop.br",
		Role:  domain.RoleCoordinator,
	}
	_, _ = service.RegisterMember(ctx, entityID, req2)

	// Deactivate first coordinator
	err := service.DeactivateMember(ctx, entityID, member.ID)
	if err != nil {
		t.Fatalf("failed to deactivate member: %v", err)
	}

	// Verify status
	deactivated, _ := service.GetMember(ctx, entityID, member.ID)
	if deactivated.Status != domain.StatusInactive {
		t.Errorf("expected status INACTIVE, got %s", deactivated.Status)
	}
}

func TestMemberService_DeactivateLastCoordinator(t *testing.T) {
	memberRepo := NewMockMemberRepository()
	workRepo := NewMockWorkRepository()
	service := NewMemberService(memberRepo, workRepo)

	ctx := context.Background()
	entityID := "test-entity-001"

	// Register only one coordinator
	req := &RegisterMemberRequest{
		Name:  "Coordenador",
		Email: "coord@coop.br",
		Role:  domain.RoleCoordinator,
	}
	member, _ := service.RegisterMember(ctx, entityID, req)

	// Try to deactivate the only coordinator
	err := service.DeactivateMember(ctx, entityID, member.ID)
	if err == nil {
		t.Error("expected error when deactivating last coordinator, got nil")
	}
}

func TestMemberService_UpdateMember(t *testing.T) {
	memberRepo := NewMockMemberRepository()
	workRepo := NewMockWorkRepository()
	service := NewMemberService(memberRepo, workRepo)

	ctx := context.Background()
	entityID := "test-entity-001"

	// Register member
	req := &RegisterMemberRequest{
		Name:  "Original Name",
		Email: "original@coop.br",
		Role:  domain.RoleMember,
	}
	member, _ := service.RegisterMember(ctx, entityID, req)

	// Update member
	updateReq := &UpdateMemberRequest{
		ID:    member.ID,
		Name:  "Updated Name",
		Email: "updated@coop.br",
	}
	updated, err := service.UpdateMember(ctx, entityID, updateReq)
	if err != nil {
		t.Fatalf("failed to update member: %v", err)
	}

	if updated.Name != "Updated Name" {
		t.Errorf("expected name Updated Name, got %s", updated.Name)
	}
	if updated.Email != "updated@coop.br" {
		t.Errorf("expected email updated@coop.br, got %s", updated.Email)
	}
}

func TestMemberService_ListMembers(t *testing.T) {
	memberRepo := NewMockMemberRepository()
	workRepo := NewMockWorkRepository()
	service := NewMemberService(memberRepo, workRepo)

	ctx := context.Background()
	entityID := "test-entity-001"

	// Register multiple members
	members := []RegisterMemberRequest{
		{Name: "João", Email: "joao@coop.br", Role: domain.RoleCoordinator},
		{Name: "Maria", Email: "maria@coop.br", Role: domain.RoleMember},
		{Name: "Pedro", Email: "pedro@coop.br", Role: domain.RoleMember},
	}

	for _, req := range members {
		_, _ = service.RegisterMember(ctx, entityID, &req)
	}

	// List all members
	list, err := service.ListMembers(ctx, entityID)
	if err != nil {
		t.Fatalf("failed to list members: %v", err)
	}

	if len(list) != 3 {
		t.Errorf("expected 3 members, got %d", len(list))
	}

	// List by role
	coordinators, err := service.ListMembersByRole(ctx, entityID, domain.RoleCoordinator)
	if err != nil {
		t.Fatalf("failed to list coordinators: %v", err)
	}

	if len(coordinators) != 1 {
		t.Errorf("expected 1 coordinator, got %d", len(coordinators))
	}
}

func TestMemberService_GetMemberStats(t *testing.T) {
	memberRepo := NewMockMemberRepository()
	workRepo := NewMockWorkRepository()
	service := NewMemberService(memberRepo, workRepo)

	ctx := context.Background()
	entityID := "test-entity-001"

	req := &RegisterMemberRequest{
		Name:  "Maria Silva",
		Email: "maria@coop.br",
		Role:  domain.RoleMember,
	}
	member, _ := service.RegisterMember(ctx, entityID, req)

	stats, err := service.GetMemberStats(ctx, entityID, member.ID)
	if err != nil {
		t.Fatalf("failed to get member stats: %v", err)
	}

	if stats.MemberID != member.ID {
		t.Errorf("expected member ID %s, got %s", member.ID, stats.MemberID)
	}
	if stats.TotalWorkMinutes != 480 {
		t.Errorf("expected 480 minutes, got %d", stats.TotalWorkMinutes)
	}
	if stats.TotalWorkHours != 8.0 {
		t.Errorf("expected 8 hours, got %f", stats.TotalWorkHours)
	}
}

func TestMemberService_ValidateCoordinatorExists(t *testing.T) {
	memberRepo := NewMockMemberRepository()
	workRepo := NewMockWorkRepository()
	service := NewMemberService(memberRepo, workRepo)

	ctx := context.Background()
	entityID := "test-entity-001"

	// No coordinator - should fail
	err := service.ValidateCoordinatorExists(ctx, entityID)
	if err == nil {
		t.Error("expected error when no coordinator exists, got nil")
	}

	// Add coordinator
	req := &RegisterMemberRequest{
		Name:  "Coordenador",
		Email: "coord@coop.br",
		Role:  domain.RoleCoordinator,
	}
	_, _ = service.RegisterMember(ctx, entityID, req)

	// Now should succeed
	err = service.ValidateCoordinatorExists(ctx, entityID)
	if err != nil {
		t.Errorf("expected no error with coordinator, got %v", err)
	}
}

func TestMemberService_GetEntityStats(t *testing.T) {
	memberRepo := NewMockMemberRepository()
	workRepo := NewMockWorkRepository()
	service := NewMemberService(memberRepo, workRepo)

	ctx := context.Background()
	entityID := "test-entity-001"

	// Register members
	_, _ = service.RegisterMember(ctx, entityID, &RegisterMemberRequest{
		Name:  "Coordenador",
		Email: "coord@coop.br",
		Role:  domain.RoleCoordinator,
	})
	_, _ = service.RegisterMember(ctx, entityID, &RegisterMemberRequest{
		Name:  "Membro 1",
		Email: "membro1@coop.br",
		Role:  domain.RoleMember,
	})
	_, _ = service.RegisterMember(ctx, entityID, &RegisterMemberRequest{
		Name:  "Membro 2",
		Email: "membro2@coop.br",
		Role:  domain.RoleMember,
	})

	stats, err := service.GetEntityStats(ctx, entityID)
	if err != nil {
		t.Fatalf("failed to get entity stats: %v", err)
	}

	if stats.TotalMembers != 3 {
		t.Errorf("expected 3 total members, got %d", stats.TotalMembers)
	}
	if stats.ActiveMembers != 3 {
		t.Errorf("expected 3 active members, got %d", stats.ActiveMembers)
	}
	if stats.Coordinators != 1 {
		t.Errorf("expected 1 coordinator, got %d", stats.Coordinators)
	}
}
