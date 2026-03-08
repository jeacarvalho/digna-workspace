package repository

import (
	"fmt"
	"testing"
	"time"

	"github.com/providentia/digna/core_lume/internal/domain"
	"github.com/providentia/digna/lifecycle/pkg/lifecycle"
)

func generateTestID(prefix string) string {
	return fmt.Sprintf("%s-%d", prefix, time.Now().UnixNano())
}

func setupTestMemberRepo(t *testing.T) (*SQLiteMemberRepository, lifecycle.LifecycleManager, func()) {
	lm := lifecycle.NewSQLiteManager()
	repo := NewSQLiteMemberRepository(lm)

	cleanup := func() {
		lm.CloseAll()
	}

	return repo, lm, cleanup
}

func TestSQLiteMemberRepository_SaveAndFind(t *testing.T) {
	repo, _, cleanup := setupTestMemberRepo(t)
	defer cleanup()

	entityID := generateTestID("entity")

	member := &domain.Member{
		ID:        generateTestID("member"),
		EntityID:  entityID,
		Name:      "Maria Silva",
		Email:     "maria@coop.br",
		Phone:     "+55 11 98765-4321",
		CPF:       "123.456.789-00",
		Role:      domain.RoleCoordinator,
		Status:    domain.StatusActive,
		JoinedAt:  time.Now(),
		Skills:    []string{"mel", "gestão", "marketing"},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Test Save
	err := repo.Save(member)
	if err != nil {
		t.Fatalf("failed to save member: %v", err)
	}

	// Test FindByID
	found, err := repo.FindByID(entityID, member.ID)
	if err != nil {
		t.Fatalf("failed to find member by id: %v", err)
	}

	if found.ID != member.ID {
		t.Errorf("expected ID %s, got %s", member.ID, found.ID)
	}
	if found.Name != member.Name {
		t.Errorf("expected Name %s, got %s", member.Name, found.Name)
	}
	if found.Email != member.Email {
		t.Errorf("expected Email %s, got %s", member.Email, found.Email)
	}
	if found.Role != member.Role {
		t.Errorf("expected Role %s, got %s", member.Role, found.Role)
	}
	if found.Status != member.Status {
		t.Errorf("expected Status %s, got %s", member.Status, found.Status)
	}

	// Test FindByEmail
	foundByEmail, err := repo.FindByEmail(entityID, member.Email)
	if err != nil {
		t.Fatalf("failed to find member by email: %v", err)
	}

	if foundByEmail.ID != member.ID {
		t.Errorf("expected ID %s, got %s", member.ID, foundByEmail.ID)
	}
}

func TestSQLiteMemberRepository_ListByEntity(t *testing.T) {
	repo, _, cleanup := setupTestMemberRepo(t)
	defer cleanup()

	entityID := generateTestID("entity")

	// Create multiple members
	members := []*domain.Member{
		{
			ID:        generateTestID("member"),
			EntityID:  entityID,
			Name:      "João Silva",
			Email:     "joao@coop.br",
			Role:      domain.RoleCoordinator,
			Status:    domain.StatusActive,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			ID:        generateTestID("member"),
			EntityID:  entityID,
			Name:      "Ana Santos",
			Email:     "ana@coop.br",
			Role:      domain.RoleMember,
			Status:    domain.StatusActive,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			ID:        generateTestID("member"),
			EntityID:  entityID,
			Name:      "Carlos Oliveira",
			Email:     "carlos@coop.br",
			Role:      domain.RoleAdvisor,
			Status:    domain.StatusInactive,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	for _, m := range members {
		if err := repo.Save(m); err != nil {
			t.Fatalf("failed to save member: %v", err)
		}
	}

	// Test ListByEntity
	list, err := repo.ListByEntity(entityID)
	if err != nil {
		t.Fatalf("failed to list members: %v", err)
	}

	if len(list) != 3 {
		t.Errorf("expected 3 members, got %d", len(list))
	}

	// Test ListByRole
	coordinators, err := repo.ListByRole(entityID, domain.RoleCoordinator)
	if err != nil {
		t.Fatalf("failed to list coordinators: %v", err)
	}

	if len(coordinators) != 1 {
		t.Errorf("expected 1 coordinator, got %d", len(coordinators))
	}

	if coordinators[0].Name != "João Silva" {
		t.Errorf("expected João Silva, got %s", coordinators[0].Name)
	}

	// Test CountByEntity
	count, err := repo.CountByEntity(entityID)
	if err != nil {
		t.Fatalf("failed to count members: %v", err)
	}

	if count != 3 {
		t.Errorf("expected count 3, got %d", count)
	}

	// Test CountActiveByEntity
	activeCount, err := repo.CountActiveByEntity(entityID)
	if err != nil {
		t.Fatalf("failed to count active members: %v", err)
	}

	if activeCount != 2 {
		t.Errorf("expected active count 2, got %d", activeCount)
	}
}

func TestSQLiteMemberRepository_UpdateStatus(t *testing.T) {
	repo, _, cleanup := setupTestMemberRepo(t)
	defer cleanup()

	entityID := generateTestID("entity")

	member := &domain.Member{
		ID:        generateTestID("member"),
		EntityID:  entityID,
		Name:      "Pedro Costa",
		Email:     "pedro@coop.br",
		Role:      domain.RoleMember,
		Status:    domain.StatusActive,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := repo.Save(member); err != nil {
		t.Fatalf("failed to save member: %v", err)
	}

	// Update status to inactive
	err := repo.UpdateStatus(entityID, member.ID, domain.StatusInactive)
	if err != nil {
		t.Fatalf("failed to update member status: %v", err)
	}

	// Verify status change
	found, err := repo.FindByID(entityID, member.ID)
	if err != nil {
		t.Fatalf("failed to find member: %v", err)
	}

	if found.Status != domain.StatusInactive {
		t.Errorf("expected status INACTIVE, got %s", found.Status)
	}
}

func TestSQLiteMemberRepository_Update(t *testing.T) {
	repo, _, cleanup := setupTestMemberRepo(t)
	defer cleanup()

	entityID := generateTestID("entity")

	member := &domain.Member{
		ID:        generateTestID("member"),
		EntityID:  entityID,
		Name:      "Fernanda Lima",
		Email:     "fernanda@coop.br",
		Phone:     "+55 11 91234-5678",
		Role:      domain.RoleMember,
		Status:    domain.StatusActive,
		Skills:    []string{"apicultura"},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := repo.Save(member); err != nil {
		t.Fatalf("failed to save member: %v", err)
	}

	// Update member
	member.Name = "Fernanda Lima Santos"
	member.Phone = "+55 11 99876-5432"
	member.AddSkill("gestão")

	err := repo.Update(member)
	if err != nil {
		t.Fatalf("failed to update member: %v", err)
	}

	// Verify changes
	found, err := repo.FindByID(entityID, member.ID)
	if err != nil {
		t.Fatalf("failed to find member: %v", err)
	}

	if found.Name != "Fernanda Lima Santos" {
		t.Errorf("expected name Fernanda Lima Santos, got %s", found.Name)
	}
	if found.Phone != "+55 11 99876-5432" {
		t.Errorf("expected phone +55 11 99876-5432, got %s", found.Phone)
	}
}

func TestSQLiteMemberRepository_InvalidMember(t *testing.T) {
	repo, _, cleanup := setupTestMemberRepo(t)
	defer cleanup()

	entityID := generateTestID("entity")

	// Test with invalid member (empty name)
	invalidMember := &domain.Member{
		ID:        generateTestID("member"),
		EntityID:  entityID,
		Name:      "",
		Email:     "invalid@coop.br",
		Role:      domain.RoleMember,
		Status:    domain.StatusActive,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err := repo.Save(invalidMember)
	if err == nil {
		t.Error("expected error for invalid member, got nil")
	}

	// Test with empty email
	invalidMember2 := &domain.Member{
		ID:        generateTestID("member"),
		EntityID:  entityID,
		Name:      "Nome Válido",
		Email:     "",
		Role:      domain.RoleMember,
		Status:    domain.StatusActive,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err = repo.Save(invalidMember2)
	if err == nil {
		t.Error("expected error for member with empty email, got nil")
	}
}

func TestSQLiteMemberRepository_FindNotFound(t *testing.T) {
	repo, _, cleanup := setupTestMemberRepo(t)
	defer cleanup()

	entityID := generateTestID("entity")

	// Try to find non-existent member
	_, err := repo.FindByID(entityID, "non-existent-id")
	if err == nil {
		t.Error("expected error for non-existent member, got nil")
	}

	// Try to find by email
	_, err = repo.FindByEmail(entityID, "nonexistent@coop.br")
	if err == nil {
		t.Error("expected error for non-existent email, got nil")
	}
}

func TestSQLiteMemberRepository_MemberStats(t *testing.T) {
	repo, _, cleanup := setupTestMemberRepo(t)
	defer cleanup()

	entityID := generateTestID("entity")

	member := &domain.Member{
		ID:        generateTestID("member"),
		EntityID:  entityID,
		Name:      "Lucia Pereira",
		Email:     "lucia@coop.br",
		Role:      domain.RoleMember,
		Status:    domain.StatusActive,
		JoinedAt:  time.Now(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := repo.Save(member); err != nil {
		t.Fatalf("failed to save member: %v", err)
	}

	// Test member voting rights
	if !member.CanVote() {
		t.Error("expected active member to have voting rights")
	}

	if member.IsCoordinator() {
		t.Error("expected member not to be coordinator")
	}

	// Deactivate and test
	member.Deactivate()
	if member.CanVote() {
		t.Error("expected inactive member not to have voting rights")
	}

	// Test coordinator
	coordinator := &domain.Member{
		ID:       generateTestID("member"),
		EntityID: entityID,
		Name:     "Coordenador",
		Email:    "coord@coop.br",
		Role:     domain.RoleCoordinator,
		Status:   domain.StatusActive,
	}

	if !coordinator.IsCoordinator() {
		t.Error("expected coordinator to be identified as coordinator")
	}

	if !coordinator.CanManage() {
		t.Error("expected coordinator to have management rights")
	}
}
