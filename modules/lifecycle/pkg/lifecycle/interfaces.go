package lifecycle

import (
	"context"
	"database/sql"
	"time"

	"github.com/providentia/digna/lifecycle/internal/domain"
)

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
	CreateLink(enterpriseID, accountantID, delegatedBy string) (*domain.EnterpriseAccountant, error)
	// DeactivateLink deactivates a link (Exit Power)
	DeactivateLink(linkID, enterpriseID, requestedBy string) error
}
