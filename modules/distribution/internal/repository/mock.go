// Package repository implementa repositórios mock para distribuição
// Facilita testes e pode ser substituído por implementação real em SQLite
package repository

import (
	"context"
	"sync"
	"time"

	"github.com/providentia/digna/distribution/internal/domain"
)

// MockDistributionRepository implementação em memória para testes
type MockDistributionRepository struct {
	distributions map[int64]*domain.Distribution
	members       map[int64][]domain.DistributionMember
	ledgerEntries map[int64]*domain.LedgerEntry
	nextID        int64
	mu            sync.RWMutex
}

// NewMockDistributionRepository cria repositório mock
func NewMockDistributionRepository() *MockDistributionRepository {
	return &MockDistributionRepository{
		distributions: make(map[int64]*domain.Distribution),
		members:       make(map[int64][]domain.DistributionMember),
		ledgerEntries: make(map[int64]*domain.LedgerEntry),
		nextID:        1,
	}
}

// Save salva ou atualiza uma distribuição
func (r *MockDistributionRepository) Save(ctx context.Context, dist *domain.Distribution) (int64, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if dist.ID == 0 {
		dist.ID = r.nextID
		r.nextID++
	}
	r.distributions[dist.ID] = dist
	return dist.ID, nil
}

// FindByID busca distribuição por ID
func (r *MockDistributionRepository) FindByID(ctx context.Context, id int64) (*domain.Distribution, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	dist, ok := r.distributions[id]
	if !ok {
		return nil, nil
	}

	// Cria cópia para evitar modificação externa
	result := *dist

	// Carrega membros também
	if members, ok := r.members[id]; ok {
		result.Members = make([]domain.DistributionMember, len(members))
		copy(result.Members, members)
	}

	return &result, nil
}

// FindByEntityAndPeriod busca distribuição por entidade e período
func (r *MockDistributionRepository) FindByEntityAndPeriod(ctx context.Context, entityID, period string) (*domain.Distribution, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, dist := range r.distributions {
		if dist.EntityID == entityID && dist.Period == period {
			result := *dist
			return &result, nil
		}
	}
	return nil, nil
}

// FindByStatus busca distribuições por status
func (r *MockDistributionRepository) FindByStatus(ctx context.Context, entityID string, status domain.DistributionStatus) ([]*domain.Distribution, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []*domain.Distribution
	for _, dist := range r.distributions {
		if dist.EntityID == entityID && dist.Status == status {
			d := *dist
			result = append(result, &d)
		}
	}
	return result, nil
}

// UpdateStatus atualiza o status
func (r *MockDistributionRepository) UpdateStatus(ctx context.Context, id int64, status domain.DistributionStatus) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	dist, ok := r.distributions[id]
	if !ok {
		return nil
	}
	dist.Status = status
	return nil
}

// UpdateAssemblyDecisionID atualiza a referência à decisão
func (r *MockDistributionRepository) UpdateAssemblyDecisionID(ctx context.Context, id int64, decisionID int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	dist, ok := r.distributions[id]
	if !ok {
		return nil
	}
	dist.AssemblyDecisionID = decisionID
	return nil
}

// MarkAsExecuted marca como executada
func (r *MockDistributionRepository) MarkAsExecuted(ctx context.Context, id int64, executedAt time.Time) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	dist, ok := r.distributions[id]
	if !ok {
		return nil
	}
	dist.Status = domain.StatusExecuted
	dist.ExecutedAt = &executedAt
	return nil
}

// SaveMember salva um membro da distribuição (atualiza se já existe)
func (r *MockDistributionRepository) SaveMember(ctx context.Context, member *domain.DistributionMember) (int64, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if member.ID == 0 {
		member.ID = r.nextID
		r.nextID++
	}

	// Busca lista de membros da distribuição
	members := r.members[member.DistributionID]

	// Verifica se o membro já existe (atualiza)
	found := false
	for i := range members {
		if members[i].ID == member.ID {
			members[i] = *member
			found = true
			break
		}
	}

	// Se não existe, adiciona novo
	if !found {
		members = append(members, *member)
	}

	r.members[member.DistributionID] = members
	return member.ID, nil
}

// FindMembersByDistribution busca membros de uma distribuição
func (r *MockDistributionRepository) FindMembersByDistribution(ctx context.Context, distributionID int64) ([]domain.DistributionMember, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	members, ok := r.members[distributionID]
	if !ok {
		return []domain.DistributionMember{}, nil
	}

	// Retorna cópia
	result := make([]domain.DistributionMember, len(members))
	copy(result, members)
	return result, nil
}

// SaveLedgerEntry salva um lançamento contábil
func (r *MockDistributionRepository) SaveLedgerEntry(ctx context.Context, entityID string, entry *domain.LedgerEntry) (int64, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if entry.ID == 0 {
		entry.ID = r.nextID
		r.nextID++
	}
	r.ledgerEntries[entry.ID] = entry
	return entry.ID, nil
}

// MockAssemblyRepository implementação mock para assembly
type MockAssemblyRepository struct {
	decisions map[int64]*domain.AssemblyDecision
}

// NewMockAssemblyRepository cria repositório mock
func NewMockAssemblyRepository() *MockAssemblyRepository {
	return &MockAssemblyRepository{
		decisions: make(map[int64]*domain.AssemblyDecision),
	}
}

// AddDecision adiciona uma decisão (helper para testes)
func (r *MockAssemblyRepository) AddDecision(decision *domain.AssemblyDecision) {
	r.decisions[decision.ID] = decision
}

// FindByID busca decisão por ID
func (r *MockAssemblyRepository) FindByID(ctx context.Context, id int64) (*domain.AssemblyDecision, error) {
	decision, ok := r.decisions[id]
	if !ok {
		return nil, nil
	}
	return decision, nil
}

// IsApproved verifica se a decisão está aprovada
func (r *MockAssemblyRepository) IsApproved(ctx context.Context, id int64) bool {
	decision, ok := r.decisions[id]
	if !ok {
		return false
	}
	return decision.Status == "APPROVED"
}
