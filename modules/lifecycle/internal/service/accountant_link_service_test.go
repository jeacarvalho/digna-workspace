package service

import (
	"testing"
	"time"

	"github.com/providentia/digna/lifecycle/internal/domain"
	"github.com/providentia/digna/lifecycle/internal/repository"
)

// MockRepository implementa EnterpriseAccountantRepository para testes
type MockRepository struct {
	links map[string]*domain.EnterpriseAccountant
}

func NewMockRepository() *MockRepository {
	return &MockRepository{
		links: make(map[string]*domain.EnterpriseAccountant),
	}
}

func (m *MockRepository) Create(link *domain.EnterpriseAccountant) error {
	if link.ID == "" {
		// Gerar ID único baseado no timestamp com nanossegundos
		link.ID = "mock_" + time.Now().Format("20060102150405.000000000")
	}
	// Criar uma cópia para evitar referências compartilhadas
	linkCopy := *link
	m.links[link.ID] = &linkCopy
	return nil
}

func (m *MockRepository) Update(link *domain.EnterpriseAccountant) error {
	if _, exists := m.links[link.ID]; !exists {
		return nil // Simulando "not found"
	}
	// Atualizar o link existente
	existingLink := m.links[link.ID]
	existingLink.Status = link.Status
	existingLink.EndDate = link.EndDate
	existingLink.UpdatedAt = link.UpdatedAt
	return nil
}

func (m *MockRepository) FindByID(id string) (*domain.EnterpriseAccountant, error) {
	link, exists := m.links[id]
	if !exists {
		return nil, nil
	}
	return link, nil
}

func (m *MockRepository) FindByEnterpriseID(enterpriseID string) ([]*domain.EnterpriseAccountant, error) {
	var result []*domain.EnterpriseAccountant
	for _, link := range m.links {
		if link.EnterpriseID == enterpriseID {
			result = append(result, link)
		}
	}
	return result, nil
}

func (m *MockRepository) FindByAccountantID(accountantID string) ([]*domain.EnterpriseAccountant, error) {
	var result []*domain.EnterpriseAccountant
	for _, link := range m.links {
		if link.AccountantID == accountantID {
			result = append(result, link)
		}
	}
	return result, nil
}

func (m *MockRepository) FindActiveByEnterpriseID(enterpriseID string) (*domain.EnterpriseAccountant, error) {
	for _, link := range m.links {
		if link.EnterpriseID == enterpriseID && link.IsActive() {
			return link, nil
		}
	}
	return nil, nil
}

func (m *MockRepository) FindActiveByAccountantID(accountantID string) ([]*domain.EnterpriseAccountant, error) {
	var result []*domain.EnterpriseAccountant
	for _, link := range m.links {
		if link.AccountantID == accountantID && link.IsActive() {
			result = append(result, link)
		}
	}
	return result, nil
}

func (m *MockRepository) FindByDateRange(enterpriseID, accountantID string, startDate, endDate time.Time) ([]*domain.EnterpriseAccountant, error) {
	var result []*domain.EnterpriseAccountant
	for _, link := range m.links {
		if link.EnterpriseID == enterpriseID && link.AccountantID == accountantID {
			if link.IsValidForDate(startDate) || link.IsValidForDate(endDate) {
				result = append(result, link)
			}
		}
	}
	return result, nil
}

func (m *MockRepository) FindByAccountantAndEnterprise(accountantID, enterpriseID string) ([]*domain.EnterpriseAccountant, error) {
	var result []*domain.EnterpriseAccountant
	for _, link := range m.links {
		if link.AccountantID == accountantID && link.EnterpriseID == enterpriseID {
			result = append(result, link)
		}
	}
	return result, nil
}

func (m *MockRepository) FindByAccountantAndEnterpriseInDateRange(accountantID, enterpriseID string, startTime, endTime int64) ([]*domain.EnterpriseAccountant, error) {
	var result []*domain.EnterpriseAccountant
	startDate := time.Unix(startTime, 0)
	endDate := time.Unix(endTime, 0)
	for _, link := range m.links {
		if link.AccountantID == accountantID && link.EnterpriseID == enterpriseID {
			if link.IsValidForDate(startDate) || link.IsValidForDate(endDate) {
				result = append(result, link)
			}
		}
	}
	return result, nil
}

func (m *MockRepository) FindByAccountantIDAndDateRange(accountantID string, startTime, endTime int64) ([]*domain.EnterpriseAccountant, error) {
	var result []*domain.EnterpriseAccountant
	startDate := time.Unix(startTime, 0)
	endDate := time.Unix(endTime, 0)
	for _, link := range m.links {
		if link.AccountantID == accountantID && link.IsActive() {
			if link.IsValidForDate(startDate) || link.IsValidForDate(endDate) {
				result = append(result, link)
			}
		}
	}
	return result, nil
}

var _ repository.EnterpriseAccountantRepository = (*MockRepository)(nil)

func TestAccountantLinkService_CreateLink(t *testing.T) {
	repo := NewMockRepository()
	service := NewAccountantLinkService(repo)

	// Testar criação de primeiro vínculo
	link, err := service.CreateLink("ent_123", "acc_456", "user_789")
	if err != nil {
		t.Fatalf("Unexpected error creating first link: %v", err)
	}
	if link == nil {
		t.Fatal("Expected link to be created")
	}
	if !link.IsActive() {
		t.Error("Expected new link to be active")
	}

	// Testar criação de segundo vínculo (deve inativar o primeiro)
	link2, err := service.CreateLink("ent_123", "acc_789", "user_999")
	if err != nil {
		t.Fatalf("Unexpected error creating second link: %v", err)
	}
	if !link2.IsActive() {
		t.Error("Expected second link to be active")
	}

	// Debug: listar todos os links no repositório
	t.Logf("Total links in repo: %d", len(repo.links))
	for id, l := range repo.links {
		t.Logf("Link %s: Enterprise=%s, Accountant=%s, Status=%s, EndDate=%v",
			id, l.EnterpriseID, l.AccountantID, l.Status, l.EndDate)
	}

	// Verificar que o primeiro link foi inativado
	firstLink, _ := repo.FindByID(link.ID)
	if firstLink == nil {
		t.Fatal("First link not found in repository")
	}
	if firstLink.IsActive() {
		t.Errorf("Expected first link to be inactive after creating second link. Status: %s, EndDate: %v",
			firstLink.Status, firstLink.EndDate)
	}
}

func TestAccountantLinkService_DeactivateLink(t *testing.T) {
	repo := NewMockRepository()
	service := NewAccountantLinkService(repo)

	// Criar um link
	link, _ := service.CreateLink("ent_123", "acc_456", "user_789")

	// Testar desativação válida
	err := service.DeactivateLink(link.ID, "ent_123", "user_789")
	if err != nil {
		t.Fatalf("Unexpected error deactivating link: %v", err)
	}

	// Verificar que o link foi inativado
	updatedLink, _ := repo.FindByID(link.ID)
	if updatedLink.IsActive() {
		t.Error("Expected link to be inactive after deactivation")
	}

	// Testar desativação com delegator errado (Exit Power)
	link2, _ := service.CreateLink("ent_456", "acc_789", "user_111")
	err = service.DeactivateLink(link2.ID, "ent_456", "user_wrong")
	if err == nil {
		t.Error("Expected error when wrong delegator tries to deactivate")
	}

	// Testar desativação de link que não pertence à cooperativa
	err = service.DeactivateLink(link2.ID, "ent_wrong", "user_111")
	if err == nil {
		t.Error("Expected error when enterprise ID doesn't match")
	}

	// Testar desativação de link não encontrado
	err = service.DeactivateLink("non_existent", "ent_123", "user_789")
	if err == nil {
		t.Error("Expected error when link not found")
	}
}

func TestAccountantLinkService_GetValidDateRange(t *testing.T) {
	repo := NewMockRepository()
	service := NewAccountantLinkService(repo)

	// Criar um link
	link, _ := service.CreateLink("ent_123", "acc_456", "user_789")
	link.StartDate = time.Now().UTC().Add(-24 * time.Hour)

	// Testar obtenção do período válido
	start, end, err := service.GetValidDateRange("ent_123", "acc_456")
	if err != nil {
		t.Fatalf("Unexpected error getting date range: %v", err)
	}
	if start.IsZero() || end.IsZero() {
		t.Error("Expected non-zero dates")
	}
	if !start.Before(end) {
		t.Error("Expected start date to be before end date")
	}

	// Testar com contador sem vínculo
	_, _, err = service.GetValidDateRange("ent_123", "acc_nonexistent")
	if err == nil {
		t.Error("Expected error when no valid link found")
	}
}

func TestAccountantLinkService_GetActiveAccountant(t *testing.T) {
	repo := NewMockRepository()
	service := NewAccountantLinkService(repo)

	// Testar quando não há contador ativo
	active, err := service.GetActiveAccountant("ent_123")
	if err != nil {
		t.Fatalf("Unexpected error getting active accountant: %v", err)
	}
	if active != nil {
		t.Error("Expected no active accountant for new enterprise")
	}

	// Criar um link ativo
	_, _ = service.CreateLink("ent_123", "acc_456", "user_789")

	// Testar obtenção do contador ativo
	active, err = service.GetActiveAccountant("ent_123")
	if err != nil {
		t.Fatalf("Unexpected error getting active accountant: %v", err)
	}
	if active == nil {
		t.Fatal("Expected active accountant")
	}
	if active.AccountantID != "acc_456" {
		t.Errorf("Expected accountant ID acc_456, got %s", active.AccountantID)
	}
}

func TestAccountantLinkService_IsAccountantValidForDate(t *testing.T) {
	repo := NewMockRepository()
	service := NewAccountantLinkService(repo)

	// Criar um link com data específica
	link, _ := service.CreateLink("ent_123", "acc_456", "user_789")
	yesterday := time.Now().UTC().Add(-24 * time.Hour)
	link.StartDate = yesterday
	_ = repo.Update(link)

	// Testar validade para data dentro do período
	valid, err := service.IsAccountantValidForDate("ent_123", "acc_456", time.Now().UTC())
	if err != nil {
		t.Fatalf("Unexpected error checking validity: %v", err)
	}
	if !valid {
		t.Error("Expected accountant to be valid for current date")
	}

	// Testar validade para data anterior ao início
	valid, err = service.IsAccountantValidForDate("ent_123", "acc_456", yesterday.Add(-24*time.Hour))
	if err != nil {
		t.Fatalf("Unexpected error checking validity: %v", err)
	}
	if valid {
		t.Error("Expected accountant to be invalid for date before start")
	}

	// Inativar o link
	link.Deactivate(time.Now().UTC())
	_ = repo.Update(link)

	// Testar validade para data após término
	valid, err = service.IsAccountantValidForDate("ent_123", "acc_456", time.Now().UTC().Add(24*time.Hour))
	if err != nil {
		t.Fatalf("Unexpected error checking validity: %v", err)
	}
	if valid {
		t.Error("Expected accountant to be invalid for date after end")
	}
}

func TestAccountantLinkService_ReactivateLink(t *testing.T) {
	repo := NewMockRepository()
	service := NewAccountantLinkService(repo)

	// Criar e inativar um link
	link, _ := service.CreateLink("ent_123", "acc_456", "user_789")
	service.DeactivateLink(link.ID, "ent_123", "user_789")

	// Testar reativação
	err := service.ReactivateLink(link.ID)
	if err != nil {
		t.Fatalf("Unexpected error reactivating link: %v", err)
	}

	reactivatedLink, _ := repo.FindByID(link.ID)
	if !reactivatedLink.IsActive() {
		t.Error("Expected link to be active after reactivation")
	}

	// Criar outro link ativo para mesma cooperativa (isso inativará o link reativado)
	_, _ = service.CreateLink("ent_123", "acc_789", "user_999")

	// Verificar que o primeiro link foi inativado pelo CreateLink
	inactivatedLink, _ := repo.FindByID(link.ID)
	if inactivatedLink.IsActive() {
		t.Error("Expected first link to be inactive after creating new active link")
	}

	// Tentar reativar o primeiro link (deve falhar porque já existe um ativo)
	err = service.ReactivateLink(link.ID)
	if err == nil {
		t.Error("Expected error when trying to reactivate while another link is active")
	}

	// Testar reativação de link não encontrado
	err = service.ReactivateLink("non_existent")
	if err == nil {
		t.Error("Expected error when link not found")
	}
}
