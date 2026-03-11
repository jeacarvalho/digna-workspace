package service

import (
	"context"
	"fmt"
	"time"

	"github.com/providentia/digna/lifecycle/internal/domain"
	"github.com/providentia/digna/lifecycle/internal/repository"
)

type AccountantLinkService struct {
	repo repository.EnterpriseAccountantRepository
}

func NewAccountantLinkService(repo repository.EnterpriseAccountantRepository) *AccountantLinkService {
	return &AccountantLinkService{repo: repo}
}

// CreateLink cria um novo vínculo contábil com validação de cardinalidade
// Regra: Uma cooperativa só pode ter 1 contador ATIVO por vez
func (s *AccountantLinkService) CreateLink(enterpriseID, accountantID, delegatedBy string) (*domain.EnterpriseAccountant, error) {
	// Verificar se já existe um contador ativo para esta cooperativa
	activeLink, err := s.repo.FindActiveByEnterpriseID(enterpriseID)
	if err != nil {
		return nil, fmt.Errorf("failed to check active links: %w", err)
	}

	// Se já existe um ativo, inativá-lo (regra de cardinalidade)
	if activeLink != nil {
		if err := s.deactivateLink(activeLink, time.Now().UTC(), "system"); err != nil {
			return nil, fmt.Errorf("failed to deactivate previous active link: %w", err)
		}
	}

	// Criar novo vínculo
	link := domain.NewEnterpriseAccountant(enterpriseID, accountantID, delegatedBy)
	if err := s.repo.Create(link); err != nil {
		return nil, fmt.Errorf("failed to create link: %w", err)
	}

	return link, nil
}

// DeactivateLink inativa um vínculo (Exit Power da cooperativa)
// Apenas a cooperativa (delegatedBy) pode inativar o vínculo
func (s *AccountantLinkService) DeactivateLink(linkID, enterpriseID, requestedBy string) error {
	link, err := s.repo.FindByID(linkID)
	if err != nil {
		return fmt.Errorf("failed to find link: %w", err)
	}
	if link == nil {
		return fmt.Errorf("link not found: %s", linkID)
	}

	// Validar que o link pertence à cooperativa
	if link.EnterpriseID != enterpriseID {
		return fmt.Errorf("link does not belong to enterprise: %s", enterpriseID)
	}

	// Validar Exit Power: apenas quem delegou pode encerrar
	if link.DelegatedBy != requestedBy {
		return fmt.Errorf("only the delegator can deactivate the link: %s", link.DelegatedBy)
	}

	// Inativar o vínculo
	return s.deactivateLink(link, time.Now().UTC(), requestedBy)
}

// GetValidDateRange retorna o período de vigência de um contador para uma cooperativa
// Útil para filtrar acesso no accountant_dashboard
func (s *AccountantLinkService) GetValidDateRange(enterpriseID, accountantID string) (startDate, endDate time.Time, err error) {
	links, err := s.repo.FindByDateRange(enterpriseID, accountantID, time.Time{}, time.Now().UTC())
	if err != nil {
		return time.Time{}, time.Time{}, fmt.Errorf("failed to find links by date range: %w", err)
	}

	if len(links) == 0 {
		return time.Time{}, time.Time{}, fmt.Errorf("no valid link found for accountant %s in enterprise %s", accountantID, enterpriseID)
	}

	// Encontrar o link mais recente (primeiro da lista ordenada por start_date DESC)
	latestLink := links[0]
	startDate, endDate = latestLink.GetDateRange()
	return startDate, endDate, nil
}

// GetActiveAccountant retorna o contador ativo atual para uma cooperativa
func (s *AccountantLinkService) GetActiveAccountant(enterpriseID string) (*domain.EnterpriseAccountant, error) {
	return s.repo.FindActiveByEnterpriseID(enterpriseID)
}

// GetAccountantLinks retorna todos os vínculos de um contador
func (s *AccountantLinkService) GetAccountantLinks(accountantID string) ([]*domain.EnterpriseAccountant, error) {
	return s.repo.FindByAccountantID(accountantID)
}

// GetEnterpriseLinks retorna todos os vínculos de uma cooperativa
func (s *AccountantLinkService) GetEnterpriseLinks(enterpriseID string) ([]*domain.EnterpriseAccountant, error) {
	return s.repo.FindByEnterpriseID(enterpriseID)
}

// IsAccountantValidForDate verifica se um contador tinha acesso válido em uma data específica
func (s *AccountantLinkService) IsAccountantValidForDate(enterpriseID, accountantID string, checkDate time.Time) (bool, error) {
	links, err := s.repo.FindByDateRange(enterpriseID, accountantID, checkDate, checkDate)
	if err != nil {
		return false, fmt.Errorf("failed to check accountant validity: %w", err)
	}

	return len(links) > 0, nil
}

// deactivateLink é um método interno para inativar um vínculo
func (s *AccountantLinkService) deactivateLink(link *domain.EnterpriseAccountant, endDate time.Time, requestedBy string) error {
	if err := link.Deactivate(endDate); err != nil {
		return fmt.Errorf("failed to deactivate link: %w", err)
	}

	if err := s.repo.Update(link); err != nil {
		return fmt.Errorf("failed to update deactivated link: %w", err)
	}

	return nil
}

// ReactivateLink reativa um vínculo inativo (apenas para administração)
func (s *AccountantLinkService) ReactivateLink(linkID string) error {
	link, err := s.repo.FindByID(linkID)
	if err != nil {
		return fmt.Errorf("failed to find link: %w", err)
	}
	if link == nil {
		return fmt.Errorf("link not found: %s", linkID)
	}

	// Verificar se já existe um ativo para esta cooperativa
	activeLink, err := s.repo.FindActiveByEnterpriseID(link.EnterpriseID)
	if err != nil {
		return fmt.Errorf("failed to check active links: %w", err)
	}

	// Se já existe um ativo, não pode reativar (regra de cardinalidade)
	if activeLink != nil && activeLink.ID != link.ID {
		return fmt.Errorf("enterprise already has an active accountant: %s", activeLink.AccountantID)
	}

	link.Reactivate()
	if err := s.repo.Update(link); err != nil {
		return fmt.Errorf("failed to update reactivated link: %w", err)
	}

	return nil
}

// GetValidEnterprisesForAccountant returns list of enterprises an accountant can access during a period
func (s *AccountantLinkService) GetValidEnterprisesForAccountant(ctx context.Context, accountantID string, startTime, endTime time.Time) ([]string, error) {
	// Convert to Unix timestamps for database query
	startUnix := startTime.Unix()
	endUnix := endTime.Unix()

	// Find all links for this accountant that are active during the period
	links, err := s.repo.FindByAccountantIDAndDateRange(accountantID, startUnix, endUnix)
	if err != nil {
		return nil, fmt.Errorf("failed to find links for accountant: %w", err)
	}

	// Extract unique enterprise IDs
	enterprises := make([]string, 0, len(links))
	seen := make(map[string]bool)
	for _, link := range links {
		if !seen[link.EnterpriseID] {
			enterprises = append(enterprises, link.EnterpriseID)
			seen[link.EnterpriseID] = true
		}
	}

	return enterprises, nil
}

// ValidateAccountantAccess checks if an accountant has access to an enterprise during a period
func (s *AccountantLinkService) ValidateAccountantAccess(ctx context.Context, accountantID, enterpriseID string, startTime, endTime time.Time) (bool, error) {
	// Convert to Unix timestamps for database query
	startUnix := startTime.Unix()
	endUnix := endTime.Unix()

	// Find links for this accountant-enterprise pair during the period
	links, err := s.repo.FindByAccountantAndEnterpriseInDateRange(accountantID, enterpriseID, startUnix, endUnix)
	if err != nil {
		return false, fmt.Errorf("failed to validate access: %w", err)
	}

	// If we found any links, access is valid
	return len(links) > 0, nil
}
