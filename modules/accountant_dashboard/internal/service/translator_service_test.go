package service

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"digna/accountant_dashboard/internal/domain"
)

type mockFiscalRepository struct {
	entries           []domain.EntryDTO
	exportLogs        []domain.FiscalExportLog
	pendingEntities   []string
	loadEntriesErr    error
	registerExportErr error
	listPendingErr    error
	getHistoryErr     error
}

func (m *mockFiscalRepository) LoadEntries(ctx context.Context, entityID string, period string) ([]domain.EntryDTO, error) {
	if m.loadEntriesErr != nil {
		return nil, m.loadEntriesErr
	}
	return m.entries, nil
}

func (m *mockFiscalRepository) RegisterExport(ctx context.Context, entityID string, batch *domain.FiscalBatch) error {
	return m.registerExportErr
}

func (m *mockFiscalRepository) ListPendingEntities(ctx context.Context, period string) ([]string, error) {
	if m.listPendingErr != nil {
		return nil, m.listPendingErr
	}
	return m.pendingEntities, nil
}

func (m *mockFiscalRepository) GetExportHistory(ctx context.Context, entityID string, period string) ([]domain.FiscalExportLog, error) {
	if m.getHistoryErr != nil {
		return nil, m.getHistoryErr
	}
	return m.exportLogs, nil
}

func TestTranslatorService_GenerateHash(t *testing.T) {
	repo := &mockFiscalRepository{}
	mapper := domain.NewDefaultAccountMapper()
	svc := NewTranslatorService(repo, mapper)

	data := []byte("test data content")
	hash := svc.GenerateHash(data)

	if len(hash) != 64 {
		t.Errorf("GenerateHash should return 64 char SHA256 hash, got %d", len(hash))
	}

	hash2 := svc.GenerateHash(data)
	if hash != hash2 {
		t.Error("GenerateHash should return consistent results")
	}
}

func TestTranslatorService_ValidateSomaZero(t *testing.T) {
	repo := &mockFiscalRepository{}
	mapper := domain.NewDefaultAccountMapper()
	svc := NewTranslatorService(repo, mapper)

	validEntries := []domain.EntryDTO{
		{
			ID:          1,
			TotalDebit:  1000,
			TotalCredit: 1000,
		},
	}

	err := svc.validateSomaZero(validEntries)
	if err != nil {
		t.Errorf("validateSomaZero should not fail for valid entries: %v", err)
	}

	invalidEntries := []domain.EntryDTO{
		{
			ID:          1,
			TotalDebit:  1000,
			TotalCredit: 999,
		},
	}

	err = svc.validateSomaZero(invalidEntries)
	if err == nil {
		t.Error("validateSomaZero should fail for invalid entries")
	}
}

func TestTranslatorService_TranslateToStandardFormat(t *testing.T) {
	repo := &mockFiscalRepository{}
	mapper := domain.NewDefaultAccountMapper()
	svc := NewTranslatorService(repo, mapper)

	entries := []domain.EntryDTO{
		{
			ID:          1,
			EntityID:    "test_entity",
			Date:        time.Now(),
			Description: "Venda de produtos",
			Postings: []domain.PostingDTO{
				{AccountCode: "1.1.01", Debit: 1000},
				{AccountCode: "3.1.01", Credit: 1000},
			},
			TotalDebit:  1000,
			TotalCredit: 1000,
		},
	}

	data, err := svc.TranslateToStandardFormat(entries)
	if err != nil {
		t.Errorf("TranslateToStandardFormat failed: %v", err)
	}

	if len(data) == 0 {
		t.Error("TranslateToStandardFormat should return data")
	}
}

func TestGenerateBatchID(t *testing.T) {
	batchID := generateBatchID("entity1", "2026-03")

	if batchID == "" {
		t.Error("generateBatchID should not return empty string")
	}

	if batchID[:8] != "entity1_" {
		t.Errorf("generateBatchID should start with entityID_, got %s", batchID[:8])
	}
}

func TestGenerateEntryHash(t *testing.T) {
	entry := domain.EntryDTO{
		ID:          1,
		Date:        time.Date(2026, 3, 8, 0, 0, 0, 0, time.UTC),
		TotalDebit:  1000,
		TotalCredit: 1000,
	}

	hash := generateEntryHash(entry)

	if len(hash) != 64 {
		t.Errorf("generateEntryHash should return 64 char hash, got %d", len(hash))
	}

	hash2 := generateEntryHash(entry)
	if hash != hash2 {
		t.Error("generateEntryHash should return consistent results")
	}
}

func TestTranslatorService_TranslateAndExport_Success(t *testing.T) {
	repo := &mockFiscalRepository{
		entries: []domain.EntryDTO{
			{
				ID:          1,
				EntityID:    "test_entity",
				Date:        time.Now(),
				Description: "Venda de produtos",
				Postings: []domain.PostingDTO{
					{AccountCode: "1.1.01", Debit: 1000},
					{AccountCode: "3.1.01", Credit: 1000},
				},
				TotalDebit:  1000,
				TotalCredit: 1000,
			},
		},
	}
	mapper := domain.NewDefaultAccountMapper()
	svc := NewTranslatorService(repo, mapper)

	batch, data, err := svc.TranslateAndExport(context.Background(), "test_entity", "2026-03")
	if err != nil {
		t.Errorf("TranslateAndExport should succeed: %v", err)
	}

	if batch == nil {
		t.Error("TranslateAndExport should return batch")
	}

	if len(data) == 0 {
		t.Error("TranslateAndExport should return data")
	}

	if batch.EntityID != "test_entity" {
		t.Errorf("Batch should have correct entity ID, got %s", batch.EntityID)
	}

	if batch.Period != "2026-03" {
		t.Errorf("Batch should have correct period, got %s", batch.Period)
	}

	if batch.TotalEntries != 1 {
		t.Errorf("Batch should have correct entry count, got %d", batch.TotalEntries)
	}

	if batch.ExportHash == "" {
		t.Error("Batch should have export hash")
	}
}

func TestTranslatorService_TranslateAndExport_NoEntries(t *testing.T) {
	repo := &mockFiscalRepository{
		entries: []domain.EntryDTO{},
	}
	mapper := domain.NewDefaultAccountMapper()
	svc := NewTranslatorService(repo, mapper)

	batch, data, err := svc.TranslateAndExport(context.Background(), "test_entity", "2026-03")
	if err == nil {
		t.Error("TranslateAndExport should fail when no entries")
	}

	if !contains(err.Error(), "no entries found") {
		t.Errorf("Error should mention no entries, got: %v", err)
	}

	if batch != nil {
		t.Error("TranslateAndExport should not return batch on error")
	}

	if data != nil {
		t.Error("TranslateAndExport should not return data on error")
	}
}

func TestTranslatorService_TranslateAndExport_ValidationFailure(t *testing.T) {
	repo := &mockFiscalRepository{
		entries: []domain.EntryDTO{
			{
				ID:          1,
				EntityID:    "test_entity",
				Date:        time.Now(),
				Description: "Invalid entry",
				Postings: []domain.PostingDTO{
					{AccountCode: "1.1.01", Debit: 1000},
					{AccountCode: "3.1.01", Credit: 999},
				},
				TotalDebit:  1000,
				TotalCredit: 999,
			},
		},
	}
	mapper := domain.NewDefaultAccountMapper()
	svc := NewTranslatorService(repo, mapper)

	batch, data, err := svc.TranslateAndExport(context.Background(), "test_entity", "2026-03")
	if err == nil {
		t.Error("TranslateAndExport should fail on validation error")
	}

	if !contains(err.Error(), "audit validation failed") {
		t.Errorf("Error should mention validation failure, got: %v", err)
	}

	if batch != nil {
		t.Error("TranslateAndExport should not return batch on validation error")
	}

	if data != nil {
		t.Error("TranslateAndExport should not return data on validation error")
	}
}

func TestTranslatorService_TranslateAndExport_LoadEntriesError(t *testing.T) {
	repo := &mockFiscalRepository{
		loadEntriesErr: fmt.Errorf("database connection failed"),
	}
	mapper := domain.NewDefaultAccountMapper()
	svc := NewTranslatorService(repo, mapper)

	batch, data, err := svc.TranslateAndExport(context.Background(), "test_entity", "2026-03")
	if err == nil {
		t.Error("TranslateAndExport should fail on repository error")
	}

	if !contains(err.Error(), "failed to load entries") {
		t.Errorf("Error should mention load entries failure, got: %v", err)
	}

	if batch != nil {
		t.Error("TranslateAndExport should not return batch on repository error")
	}

	if data != nil {
		t.Error("TranslateAndExport should not return data on repository error")
	}
}

func TestTranslatorService_TranslateAndExport_RegisterExportError(t *testing.T) {
	repo := &mockFiscalRepository{
		entries: []domain.EntryDTO{
			{
				ID:          1,
				EntityID:    "test_entity",
				Date:        time.Now(),
				Description: "Venda de produtos",
				Postings: []domain.PostingDTO{
					{AccountCode: "1.1.01", Debit: 1000},
					{AccountCode: "3.1.01", Credit: 1000},
				},
				TotalDebit:  1000,
				TotalCredit: 1000,
			},
		},
		registerExportErr: fmt.Errorf("failed to save export"),
	}
	mapper := domain.NewDefaultAccountMapper()
	svc := NewTranslatorService(repo, mapper)

	batch, data, err := svc.TranslateAndExport(context.Background(), "test_entity", "2026-03")
	if err == nil {
		t.Error("TranslateAndExport should fail on register export error")
	}

	if !contains(err.Error(), "failed to register export") {
		t.Errorf("Error should mention register export failure, got: %v", err)
	}

	if batch != nil {
		t.Error("TranslateAndExport should not return batch on register export error")
	}

	if data != nil {
		t.Error("TranslateAndExport should not return data on register export error")
	}
}

func TestTranslatorService_ListPendingEntities(t *testing.T) {
	repo := &mockFiscalRepository{
		pendingEntities: []string{"entity1", "entity2", "entity3"},
	}
	mapper := domain.NewDefaultAccountMapper()
	svc := NewTranslatorService(repo, mapper)

	entities, err := svc.ListPendingEntities(context.Background(), "2026-03")
	if err != nil {
		t.Errorf("ListPendingEntities should succeed: %v", err)
	}

	if len(entities) != 3 {
		t.Errorf("ListPendingEntities should return 3 entities, got %d", len(entities))
	}

	if entities[0] != "entity1" || entities[1] != "entity2" || entities[2] != "entity3" {
		t.Errorf("ListPendingEntities should return correct entities, got %v", entities)
	}
}

func TestTranslatorService_ListPendingEntities_Error(t *testing.T) {
	repo := &mockFiscalRepository{
		listPendingErr: fmt.Errorf("database error"),
	}
	mapper := domain.NewDefaultAccountMapper()
	svc := NewTranslatorService(repo, mapper)

	entities, err := svc.ListPendingEntities(context.Background(), "2026-03")
	if err == nil {
		t.Error("ListPendingEntities should fail on repository error")
	}

	if entities != nil {
		t.Error("ListPendingEntities should not return entities on error")
	}
}

func TestTranslatorService_GetExportHistory(t *testing.T) {
	now := time.Now()
	repo := &mockFiscalRepository{
		exportLogs: []domain.FiscalExportLog{
			{
				ID:         "log1",
				BatchID:    "entity1_2026-03_1234567890",
				EntityID:   "entity1",
				Period:     "2026-03",
				ExportHash: "abc123",
				FilePath:   "/exports/entity1_2026-03.csv",
				ExportedAt: now.Unix(),
			},
			{
				ID:         "log2",
				BatchID:    "entity1_2026-02_1234567891",
				EntityID:   "entity1",
				Period:     "2026-02",
				ExportHash: "def456",
				FilePath:   "/exports/entity1_2026-02.csv",
				ExportedAt: now.Add(-24 * time.Hour).Unix(),
			},
		},
	}
	mapper := domain.NewDefaultAccountMapper()
	svc := NewTranslatorService(repo, mapper)

	history, err := svc.GetExportHistory(context.Background(), "entity1", "2026-03")
	if err != nil {
		t.Errorf("GetExportHistory should succeed: %v", err)
	}

	if len(history) != 2 {
		t.Errorf("GetExportHistory should return 2 logs, got %d", len(history))
	}

	if history[0].BatchID != "entity1_2026-03_1234567890" {
		t.Errorf("GetExportHistory should return correct batch ID, got %s", history[0].BatchID)
	}

	if history[1].BatchID != "entity1_2026-02_1234567891" {
		t.Errorf("GetExportHistory should return correct batch ID, got %s", history[1].BatchID)
	}
}

func TestTranslatorService_GetExportHistory_Error(t *testing.T) {
	repo := &mockFiscalRepository{
		getHistoryErr: fmt.Errorf("database error"),
	}
	mapper := domain.NewDefaultAccountMapper()
	svc := NewTranslatorService(repo, mapper)

	history, err := svc.GetExportHistory(context.Background(), "entity1", "2026-03")
	if err == nil {
		t.Error("GetExportHistory should fail on repository error")
	}

	if history != nil {
		t.Error("GetExportHistory should not return history on error")
	}
}

func TestTranslatorService_getAccountMapping_Unmapped(t *testing.T) {
	repo := &mockFiscalRepository{}
	mapper := domain.NewDefaultAccountMapper()
	svc := NewTranslatorService(repo, mapper)

	entries := []domain.EntryDTO{
		{
			ID:          1,
			EntityID:    "test_entity",
			Date:        time.Now(),
			Description: "Test with unmapped account",
			Postings: []domain.PostingDTO{
				{AccountCode: "UNKNOWN.ACCOUNT", Debit: 1000},
				{AccountCode: "3.1.01", Credit: 1000},
			},
			TotalDebit:  1000,
			TotalCredit: 1000,
		},
	}

	data, err := svc.TranslateToStandardFormat(entries)
	if err != nil {
		t.Errorf("TranslateToStandardFormat should handle unmapped accounts: %v", err)
	}

	if len(data) == 0 {
		t.Error("TranslateToStandardFormat should return data even with unmapped accounts")
	}

	if !strings.Contains(string(data), "9.9.99.99.99") {
		t.Error("TranslateToStandardFormat should use default code for unmapped accounts")
	}

	if !strings.Contains(string(data), "Conta não mapeada") {
		t.Error("TranslateToStandardFormat should use default name for unmapped accounts")
	}
}

func TestTranslatorService_TranslateToStandardFormat_EmptyEntries(t *testing.T) {
	repo := &mockFiscalRepository{}
	mapper := domain.NewDefaultAccountMapper()
	svc := NewTranslatorService(repo, mapper)

	entries := []domain.EntryDTO{}

	data, err := svc.TranslateToStandardFormat(entries)
	if err != nil {
		t.Errorf("TranslateToStandardFormat should handle empty entries: %v", err)
	}

	if len(data) == 0 {
		t.Error("TranslateToStandardFormat should return CSV with only header for empty entries")
	}

	lines := strings.Split(string(data), "\n")
	if len(lines) < 2 {
		t.Errorf("TranslateToStandardFormat should return at least header line, got %d lines", len(lines))
	}

	if !strings.Contains(lines[0], "Data,ID_Lancamento") {
		t.Errorf("TranslateToStandardFormat should include CSV header, got: %s", lines[0])
	}
}

func TestTranslatorService_TranslateToStandardFormat_ComplexPostings(t *testing.T) {
	repo := &mockFiscalRepository{}
	mapper := domain.NewDefaultAccountMapper()
	svc := NewTranslatorService(repo, mapper)

	entries := []domain.EntryDTO{
		{
			ID:          1,
			EntityID:    "test_entity",
			Date:        time.Date(2026, 3, 8, 0, 0, 0, 0, time.UTC),
			Description: "Complex entry with multiple postings",
			Postings: []domain.PostingDTO{
				{AccountCode: "1.1.01", Debit: 500},
				{AccountCode: "1.1.02", Debit: 300},
				{AccountCode: "1.1.03", Debit: 200},
				{AccountCode: "3.1.01", Credit: 400},
				{AccountCode: "3.1.02", Credit: 600},
			},
			TotalDebit:  1000,
			TotalCredit: 1000,
		},
	}

	data, err := svc.TranslateToStandardFormat(entries)
	if err != nil {
		t.Errorf("TranslateToStandardFormat should handle complex postings: %v", err)
	}

	if len(data) == 0 {
		t.Error("TranslateToStandardFormat should return data for complex postings")
	}

	lines := strings.Split(strings.TrimSpace(string(data)), "\n")
	if len(lines) != 4 {
		t.Errorf("TranslateToStandardFormat should create 3 data rows + 1 header for complex postings, got %d lines", len(lines))
	}
}

func contains(s, substr string) bool {
	return strings.Contains(s, substr)
}
