package dashboard

import (
	"context"
	"path/filepath"
	"strings"

	"digna/accountant_dashboard/internal/domain"
	"digna/accountant_dashboard/internal/repository"
	"digna/accountant_dashboard/internal/service"
)

// NewSQLiteRepositoryFactory creates a repository factory for SQLite databases
func NewSQLiteRepositoryFactory(dataDir string) RepositoryFactory {
	// Ensure dataDir ends with /entities if it doesn't already
	entitiesDir := dataDir
	if !strings.HasSuffix(dataDir, "/entities") && !strings.HasSuffix(dataDir, "\\entities") {
		entitiesDir = filepath.Join(dataDir, "entities")
	}

	return &sqliteRepositoryFactory{
		dataDir: entitiesDir,
	}
}

type sqliteRepositoryFactory struct {
	dataDir string
}

func (f *sqliteRepositoryFactory) NewRepository(entityID string) (FiscalRepository, error) {
	// Create internal repository adapter with custom data directory
	internalRepo := repository.NewSQLiteFiscalAdapter()

	// We need to modify the internal adapter to use our data directory
	// Since the internal adapter has a private basePath field, we'll create a wrapper
	// that overrides the openReadOnly method
	return &customPathRepositoryAdapter{
		repo:     internalRepo,
		dataDir:  f.dataDir,
		entityID: entityID,
	}, nil
}

// NewDashboardService creates a new dashboard service with default account mapper
func NewDashboardService(repo FiscalRepository) DashboardService {
	// Create internal account mapper
	internalMapper := domain.NewDefaultAccountMapper()

	// Create internal service
	internalService := service.NewTranslatorService(
		&serviceAdapter{repo: repo},
		internalMapper,
	)

	// Wrap it in an adapter
	return &serviceAdapterWrapper{
		service: internalService,
		repo:    repo,
	}
}

// customPathRepositoryAdapter adapts internal repository with custom data directory
type customPathRepositoryAdapter struct {
	repo     *repository.SQLiteFiscalAdapter
	dataDir  string
	entityID string
}

func (a *customPathRepositoryAdapter) LoadEntries(ctx context.Context, entityID string, period string) ([]EntryDTO, error) {
	// Temporarily override the base path
	originalBasePath := a.repo.BasePath()
	a.repo.SetBasePath(a.dataDir)
	defer a.repo.SetBasePath(originalBasePath)

	internalEntries, err := a.repo.LoadEntries(ctx, entityID, period)
	if err != nil {
		return nil, err
	}

	// Convert internal entries to public entries
	entries := make([]EntryDTO, len(internalEntries))
	for i, ie := range internalEntries {
		entries[i] = convertEntryToPublic(ie)
	}
	return entries, nil
}

func (a *customPathRepositoryAdapter) RegisterExport(ctx context.Context, entityID string, batch *FiscalBatch) error {
	// Temporarily override the base path
	originalBasePath := a.repo.BasePath()
	a.repo.SetBasePath(a.dataDir)
	defer a.repo.SetBasePath(originalBasePath)

	internalBatch := convertBatchToInternal(batch)
	return a.repo.RegisterExport(ctx, entityID, internalBatch)
}

func (a *customPathRepositoryAdapter) ListPendingEntities(ctx context.Context, period string) ([]string, error) {
	// Temporarily override the base path
	originalBasePath := a.repo.BasePath()
	a.repo.SetBasePath(a.dataDir)
	defer a.repo.SetBasePath(originalBasePath)

	return a.repo.ListPendingEntities(ctx, period)
}

func (a *customPathRepositoryAdapter) GetExportHistory(ctx context.Context, entityID string, period string) ([]FiscalExportLog, error) {
	// Temporarily override the base path
	originalBasePath := a.repo.BasePath()
	a.repo.SetBasePath(a.dataDir)
	defer a.repo.SetBasePath(originalBasePath)

	internalLogs, err := a.repo.GetExportHistory(ctx, entityID, period)
	if err != nil {
		return nil, err
	}

	logs := make([]FiscalExportLog, len(internalLogs))
	for i, il := range internalLogs {
		logs[i] = convertLogToPublic(il)
	}
	return logs, nil
}

// repositoryAdapter adapts internal repository to public interface
type repositoryAdapter struct {
	repo domain.FiscalRepository
}

func (a *repositoryAdapter) LoadEntries(ctx context.Context, entityID string, period string) ([]EntryDTO, error) {
	internalEntries, err := a.repo.LoadEntries(ctx, entityID, period)
	if err != nil {
		return nil, err
	}

	// Convert internal entries to public entries
	entries := make([]EntryDTO, len(internalEntries))
	for i, ie := range internalEntries {
		entries[i] = convertEntryToPublic(ie)
	}
	return entries, nil
}

func (a *repositoryAdapter) RegisterExport(ctx context.Context, entityID string, batch *FiscalBatch) error {
	internalBatch := convertBatchToInternal(batch)
	return a.repo.RegisterExport(ctx, entityID, internalBatch)
}

func (a *repositoryAdapter) ListPendingEntities(ctx context.Context, period string) ([]string, error) {
	return a.repo.ListPendingEntities(ctx, period)
}

func (a *repositoryAdapter) GetExportHistory(ctx context.Context, entityID string, period string) ([]FiscalExportLog, error) {
	internalLogs, err := a.repo.GetExportHistory(ctx, entityID, period)
	if err != nil {
		return nil, err
	}

	logs := make([]FiscalExportLog, len(internalLogs))
	for i, il := range internalLogs {
		logs[i] = convertLogToPublic(il)
	}
	return logs, nil
}

// serviceAdapter adapts public repository to internal interface
type serviceAdapter struct {
	repo FiscalRepository
}

func (a *serviceAdapter) LoadEntries(ctx context.Context, entityID string, period string) ([]domain.EntryDTO, error) {
	publicEntries, err := a.repo.LoadEntries(ctx, entityID, period)
	if err != nil {
		return nil, err
	}

	entries := make([]domain.EntryDTO, len(publicEntries))
	for i, pe := range publicEntries {
		entries[i] = convertEntryToInternal(pe)
	}
	return entries, nil
}

func (a *serviceAdapter) RegisterExport(ctx context.Context, entityID string, batch *domain.FiscalBatch) error {
	publicBatch := convertBatchToPublic(batch)
	return a.repo.RegisterExport(ctx, entityID, publicBatch)
}

func (a *serviceAdapter) ListPendingEntities(ctx context.Context, period string) ([]string, error) {
	return a.repo.ListPendingEntities(ctx, period)
}

func (a *serviceAdapter) GetExportHistory(ctx context.Context, entityID string, period string) ([]domain.FiscalExportLog, error) {
	publicLogs, err := a.repo.GetExportHistory(ctx, entityID, period)
	if err != nil {
		return nil, err
	}

	logs := make([]domain.FiscalExportLog, len(publicLogs))
	for i, pl := range publicLogs {
		logs[i] = convertLogToInternal(pl)
	}
	return logs, nil
}

// serviceAdapterWrapper wraps internal service with public interface
type serviceAdapterWrapper struct {
	service *service.TranslatorService
	repo    FiscalRepository
}

func (w *serviceAdapterWrapper) TranslateAndExport(ctx context.Context, entityID string, period string) (*FiscalBatch, []byte, error) {
	internalBatch, data, err := w.service.TranslateAndExport(ctx, entityID, period)
	if err != nil {
		return nil, nil, err
	}

	publicBatch := convertBatchToPublic(internalBatch)
	return publicBatch, data, nil
}

func (w *serviceAdapterWrapper) ListPendingEntities(ctx context.Context, period string) ([]string, error) {
	return w.service.ListPendingEntities(ctx, period)
}

func (w *serviceAdapterWrapper) GetExportHistory(ctx context.Context, entityID string, period string) ([]FiscalExportLog, error) {
	internalLogs, err := w.service.GetExportHistory(ctx, entityID, period)
	if err != nil {
		return nil, err
	}

	logs := make([]FiscalExportLog, len(internalLogs))
	for i, il := range internalLogs {
		logs[i] = convertLogToPublic(il)
	}
	return logs, nil
}

// Conversion functions
func convertEntryToPublic(internal domain.EntryDTO) EntryDTO {
	postings := make([]PostingDTO, len(internal.Postings))
	for i, p := range internal.Postings {
		postings[i] = PostingDTO{
			ID:          p.ID,
			EntryID:     p.EntryID,
			AccountID:   p.AccountID,
			AccountCode: p.AccountCode,
			AccountName: p.AccountName,
			Debit:       p.Debit,
			Credit:      p.Credit,
		}
	}

	return EntryDTO{
		ID:          internal.ID,
		EntityID:    internal.EntityID,
		Date:        internal.Date,
		Description: internal.Description,
		Postings:    postings,
		TotalDebit:  internal.TotalDebit,
		TotalCredit: internal.TotalCredit,
	}
}

func convertEntryToInternal(public EntryDTO) domain.EntryDTO {
	postings := make([]domain.PostingDTO, len(public.Postings))
	for i, p := range public.Postings {
		postings[i] = domain.PostingDTO{
			ID:          p.ID,
			EntryID:     p.EntryID,
			AccountID:   p.AccountID,
			AccountCode: p.AccountCode,
			AccountName: p.AccountName,
			Debit:       p.Debit,
			Credit:      p.Credit,
		}
	}

	return domain.EntryDTO{
		ID:          public.ID,
		EntityID:    public.EntityID,
		Date:        public.Date,
		Description: public.Description,
		Postings:    postings,
		TotalDebit:  public.TotalDebit,
		TotalCredit: public.TotalCredit,
	}
}

func convertBatchToPublic(internal *domain.FiscalBatch) *FiscalBatch {
	if internal == nil {
		return nil
	}

	return &FiscalBatch{
		ID:           internal.ID,
		EntityID:     internal.EntityID,
		Period:       internal.Period,
		TotalEntries: internal.TotalEntries,
		ExportHash:   internal.ExportHash,
		CreatedAt:    internal.CreatedAt,
	}
}

func convertBatchToInternal(public *FiscalBatch) *domain.FiscalBatch {
	if public == nil {
		return nil
	}

	return &domain.FiscalBatch{
		ID:           public.ID,
		EntityID:     public.EntityID,
		Period:       public.Period,
		TotalEntries: public.TotalEntries,
		ExportHash:   public.ExportHash,
		CreatedAt:    public.CreatedAt,
	}
}

func convertLogToPublic(internal domain.FiscalExportLog) FiscalExportLog {
	return FiscalExportLog{
		ID:         internal.ID,
		EntityID:   internal.EntityID,
		Period:     internal.Period,
		BatchID:    internal.BatchID,
		ExportHash: internal.ExportHash,
		FilePath:   internal.FilePath,
		ExportedAt: internal.ExportedAt,
	}
}

func convertLogToInternal(public FiscalExportLog) domain.FiscalExportLog {
	return domain.FiscalExportLog{
		ID:         public.ID,
		EntityID:   public.EntityID,
		Period:     public.Period,
		BatchID:    public.BatchID,
		ExportHash: public.ExportHash,
		FilePath:   public.FilePath,
		ExportedAt: public.ExportedAt,
	}
}
