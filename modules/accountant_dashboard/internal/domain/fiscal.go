package domain

import (
	"context"
	"time"
)

type FiscalBatch struct {
	ID           string
	EntityID     string
	Period       string
	TotalEntries int
	ExportHash   string
	CreatedAt    int64
}

type EntryDTO struct {
	ID          int64
	EntityID    string
	Date        time.Time
	Description string
	Postings    []PostingDTO
	TotalDebit  int64
	TotalCredit int64
}

type PostingDTO struct {
	ID          int64
	EntryID     int64
	AccountID   int64
	AccountCode string
	AccountName string
	Debit       int64
	Credit      int64
}

type FiscalExportLog struct {
	ID         string
	EntityID   string
	Period     string
	BatchID    string
	ExportHash string
	FilePath   string
	ExportedAt int64
}

type FiscalRepository interface {
	LoadEntries(ctx context.Context, entityID string, period string) ([]EntryDTO, error)
	RegisterExport(ctx context.Context, entityID string, batch *FiscalBatch) error
	ListPendingEntities(ctx context.Context, period string) ([]string, error)
	GetExportHistory(ctx context.Context, entityID string, period string) ([]FiscalExportLog, error)
}

type FiscalTranslator interface {
	TranslateToStandardFormat(entries []EntryDTO) ([]byte, error)
	GenerateHash(data []byte) string
}

type AccountMapping struct {
	LocalCode    string
	LocalName    string
	StandardCode string
	StandardName string
}

type AccountMapper interface {
	GetMapping(localCode string) (AccountMapping, bool)
	GetAllMappings() []AccountMapping
}
