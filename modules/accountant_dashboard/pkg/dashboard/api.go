package dashboard

import (
	"context"
	"time"
)

// FiscalBatch represents a batch of fiscal data ready for export
type FiscalBatch struct {
	ID           string
	EntityID     string
	Period       string
	TotalEntries int
	ExportHash   string
	CreatedAt    int64
}

// EntryDTO represents a journal entry with its postings
type EntryDTO struct {
	ID          int64
	EntityID    string
	Date        time.Time
	Description string
	Postings    []PostingDTO
	TotalDebit  int64
	TotalCredit int64
}

// PostingDTO represents a single posting (debit or credit) within an entry
type PostingDTO struct {
	ID          int64
	EntryID     int64
	AccountID   int64
	AccountCode string
	AccountName string
	Debit       int64
	Credit      int64
}

// FiscalExportLog represents a record of a completed export
type FiscalExportLog struct {
	ID         string
	EntityID   string
	Period     string
	BatchID    string
	ExportHash string
	FilePath   string
	ExportedAt int64
}

// AccountMapping defines how local account codes map to standard codes
type AccountMapping struct {
	LocalCode    string
	LocalName    string
	StandardCode string
	StandardName string
}

// DashboardService is the main public interface for the accountant dashboard
type DashboardService interface {
	// TranslateAndExport loads entries, validates them, and exports to standard format
	TranslateAndExport(ctx context.Context, entityID string, period string) (*FiscalBatch, []byte, error)

	// ListPendingEntities returns entity IDs that have pending fiscal closures for a period
	ListPendingEntities(ctx context.Context, period string) ([]string, error)

	// GetExportHistory returns export history for an entity and period
	GetExportHistory(ctx context.Context, entityID string, period string) ([]FiscalExportLog, error)
}

// RepositoryFactory creates repository instances for specific entity databases
type RepositoryFactory interface {
	// NewRepository creates a repository for accessing a specific entity's database
	NewRepository(entityID string) (FiscalRepository, error)
}

// FiscalRepository defines data access operations for fiscal data
type FiscalRepository interface {
	LoadEntries(ctx context.Context, entityID string, period string) ([]EntryDTO, error)
	RegisterExport(ctx context.Context, entityID string, batch *FiscalBatch) error
	ListPendingEntities(ctx context.Context, period string) ([]string, error)
	GetExportHistory(ctx context.Context, entityID string, period string) ([]FiscalExportLog, error)
}

// AccountMapper defines account mapping operations
type AccountMapper interface {
	GetMapping(localCode string) (AccountMapping, bool)
	GetAllMappings() []AccountMapping
}
