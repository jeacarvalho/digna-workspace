package lifecycle

import (
	"context"
	"database/sql"
	"time"
)

// EnterpriseAccountantPublic é uma versão pública da estrutura EnterpriseAccountant
// para uso em handlers e outros módulos
type EnterpriseAccountantPublic struct {
	ID           string
	EnterpriseID string
	AccountantID string
	Status       string
	StartDate    time.Time
	EndDate      *time.Time
	DelegatedBy  string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type LifecycleManager interface {
	GetConnection(entityID string) (*sql.DB, error)
	GetCentralConnection() (*sql.DB, error)
	CloseConnection(entityID string) error
	CloseAll() error
	EntityExists(entityID string) (bool, error)
	CreateEntity(entityID, entityName string) error
}

// AccountantLinkService provides temporal filtering for accountant-enterprise relationships
type AccountantLinkService interface {
	// GetValidEnterprisesForAccountant returns list of enterprises an accountant can access during a period
	GetValidEnterprisesForAccountant(ctx context.Context, accountantID string, startTime, endTime time.Time) ([]string, error)
	// ValidateAccountantAccess checks if an accountant has access to an enterprise during a period
	ValidateAccountantAccess(ctx context.Context, accountantID, enterpriseID string, startTime, endTime time.Time) (bool, error)
	// CreateLink creates a new accountant-enterprise link
	CreateLink(enterpriseID, accountantID, delegatedBy string) (*EnterpriseAccountantPublic, error)
	// DeactivateLink deactivates a link (Exit Power)
	DeactivateLink(linkID, enterpriseID, requestedBy string) error
	// GetEnterpriseLinks returns all links for an enterprise
	GetEnterpriseLinks(enterpriseID string) ([]*EnterpriseAccountantPublic, error)
	// GetAccountantLinks returns all links for an accountant
	GetAccountantLinks(accountantID string) ([]*EnterpriseAccountantPublic, error)
}
