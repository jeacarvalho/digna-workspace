package dashboard

import (
	"context"
	"path/filepath"
	"reflect"
	"testing"
	"time"
)

// mockFiscalRepository is a mock implementation of FiscalRepository
type mockFiscalRepository struct {
	loadEntriesFunc         func(ctx context.Context, entityID string, period string) ([]EntryDTO, error)
	registerExportFunc      func(ctx context.Context, entityID string, batch *FiscalBatch) error
	listPendingEntitiesFunc func(ctx context.Context, period string) ([]string, error)
	getExportHistoryFunc    func(ctx context.Context, entityID string, period string) ([]FiscalExportLog, error)
}

func (m *mockFiscalRepository) LoadEntries(ctx context.Context, entityID string, period string) ([]EntryDTO, error) {
	if m.loadEntriesFunc != nil {
		return m.loadEntriesFunc(ctx, entityID, period)
	}
	return []EntryDTO{}, nil
}

func (m *mockFiscalRepository) RegisterExport(ctx context.Context, entityID string, batch *FiscalBatch) error {
	if m.registerExportFunc != nil {
		return m.registerExportFunc(ctx, entityID, batch)
	}
	return nil
}

func (m *mockFiscalRepository) ListPendingEntities(ctx context.Context, period string) ([]string, error) {
	if m.listPendingEntitiesFunc != nil {
		return m.listPendingEntitiesFunc(ctx, period)
	}
	return []string{}, nil
}

func (m *mockFiscalRepository) GetExportHistory(ctx context.Context, entityID string, period string) ([]FiscalExportLog, error) {
	if m.getExportHistoryFunc != nil {
		return m.getExportHistoryFunc(ctx, entityID, period)
	}
	return []FiscalExportLog{}, nil
}

func TestNewSQLiteRepositoryFactory(t *testing.T) {
	tests := []struct {
		name     string
		dataDir  string
		expected string
	}{
		{
			name:     "dataDir without entities suffix",
			dataDir:  "/test/data",
			expected: filepath.Join("/test/data", "entities"),
		},
		{
			name:     "dataDir with forward slash entities suffix",
			dataDir:  "/test/data/entities",
			expected: "/test/data/entities",
		},
		{
			name:     "dataDir with backslash entities suffix",
			dataDir:  "C:\\test\\data\\entities",
			expected: "C:\\test\\data\\entities",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			factory := NewSQLiteRepositoryFactory(tt.dataDir)
			if factory == nil {
				t.Fatal("factory should not be nil")
			}

			repo, err := factory.NewRepository("test-entity")
			if err != nil {
				t.Fatalf("NewRepository should not return error: %v", err)
			}
			if repo == nil {
				t.Fatal("repository should not be nil")
			}
		})
	}
}

func TestNewDashboardService(t *testing.T) {
	mockRepo := &mockFiscalRepository{}

	service := NewDashboardService(mockRepo)
	if service == nil {
		t.Fatal("service should not be nil")
	}

	// Test that service implements DashboardService interface
	var _ DashboardService = service
}

func TestCustomPathRepositoryAdapter(t *testing.T) {
	t.Run("adapter creation", func(t *testing.T) {
		factory := NewSQLiteRepositoryFactory("/test/data")
		repo, err := factory.NewRepository("test-entity")
		if err != nil {
			t.Fatalf("NewRepository should not return error: %v", err)
		}
		if repo == nil {
			t.Fatal("repository should not be nil")
		}
	})
}

func TestRepositoryAdapter(t *testing.T) {
	t.Run("convertEntryToPublic", func(t *testing.T) {
		testTime := time.Now()
		publicEntry := EntryDTO{
			ID:          1,
			EntityID:    "test-entity",
			Date:        testTime,
			Description: "Test entry",
			Postings: []PostingDTO{
				{
					ID:          1,
					EntryID:     1,
					AccountID:   100,
					AccountCode: "1.01.01",
					AccountName: "Cash",
					Debit:       10000,
					Credit:      0,
				},
			},
			TotalDebit:  10000,
			TotalCredit: 0,
		}

		internalEntry := convertEntryToInternal(publicEntry)
		convertedBack := convertEntryToPublic(internalEntry)

		if publicEntry.ID != convertedBack.ID {
			t.Errorf("ID mismatch: got %v, want %v", convertedBack.ID, publicEntry.ID)
		}
		if publicEntry.EntityID != convertedBack.EntityID {
			t.Errorf("EntityID mismatch: got %v, want %v", convertedBack.EntityID, publicEntry.EntityID)
		}
		if publicEntry.Description != convertedBack.Description {
			t.Errorf("Description mismatch: got %v, want %v", convertedBack.Description, publicEntry.Description)
		}
		if publicEntry.TotalDebit != convertedBack.TotalDebit {
			t.Errorf("TotalDebit mismatch: got %v, want %v", convertedBack.TotalDebit, publicEntry.TotalDebit)
		}
		if publicEntry.TotalCredit != convertedBack.TotalCredit {
			t.Errorf("TotalCredit mismatch: got %v, want %v", convertedBack.TotalCredit, publicEntry.TotalCredit)
		}
		if len(convertedBack.Postings) != 1 {
			t.Errorf("Postings length mismatch: got %v, want 1", len(convertedBack.Postings))
		}
		if publicEntry.Postings[0].AccountCode != convertedBack.Postings[0].AccountCode {
			t.Errorf("AccountCode mismatch: got %v, want %v", convertedBack.Postings[0].AccountCode, publicEntry.Postings[0].AccountCode)
		}
		if publicEntry.Postings[0].Debit != convertedBack.Postings[0].Debit {
			t.Errorf("Debit mismatch: got %v, want %v", convertedBack.Postings[0].Debit, publicEntry.Postings[0].Debit)
		}
	})

	t.Run("convertBatchToPublic", func(t *testing.T) {
		publicBatch := &FiscalBatch{
			ID:           "batch-123",
			EntityID:     "test-entity",
			Period:       "2024-01",
			TotalEntries: 10,
			ExportHash:   "abc123",
			CreatedAt:    time.Now().Unix(),
		}

		internalBatch := convertBatchToInternal(publicBatch)
		convertedBack := convertBatchToPublic(internalBatch)

		if publicBatch.ID != convertedBack.ID {
			t.Errorf("ID mismatch: got %v, want %v", convertedBack.ID, publicBatch.ID)
		}
		if publicBatch.EntityID != convertedBack.EntityID {
			t.Errorf("EntityID mismatch: got %v, want %v", convertedBack.EntityID, publicBatch.EntityID)
		}
		if publicBatch.Period != convertedBack.Period {
			t.Errorf("Period mismatch: got %v, want %v", convertedBack.Period, publicBatch.Period)
		}
		if publicBatch.TotalEntries != convertedBack.TotalEntries {
			t.Errorf("TotalEntries mismatch: got %v, want %v", convertedBack.TotalEntries, publicBatch.TotalEntries)
		}
		if publicBatch.ExportHash != convertedBack.ExportHash {
			t.Errorf("ExportHash mismatch: got %v, want %v", convertedBack.ExportHash, publicBatch.ExportHash)
		}
		if publicBatch.CreatedAt != convertedBack.CreatedAt {
			t.Errorf("CreatedAt mismatch: got %v, want %v", convertedBack.CreatedAt, publicBatch.CreatedAt)
		}
	})

	t.Run("convertBatchToPublic with nil", func(t *testing.T) {
		result := convertBatchToPublic(nil)
		if result != nil {
			t.Errorf("convertBatchToPublic(nil) should return nil, got %v", result)
		}
	})

	t.Run("convertBatchToInternal with nil", func(t *testing.T) {
		result := convertBatchToInternal(nil)
		if result != nil {
			t.Errorf("convertBatchToInternal(nil) should return nil, got %v", result)
		}
	})

	t.Run("convertLogToPublic", func(t *testing.T) {
		publicLog := FiscalExportLog{
			ID:         "log-123",
			EntityID:   "test-entity",
			Period:     "2024-01",
			BatchID:    "batch-123",
			ExportHash: "abc123",
			FilePath:   "/exports/test.csv",
			ExportedAt: time.Now().Unix(),
		}

		internalLog := convertLogToInternal(publicLog)
		convertedBack := convertLogToPublic(internalLog)

		if publicLog.ID != convertedBack.ID {
			t.Errorf("ID mismatch: got %v, want %v", convertedBack.ID, publicLog.ID)
		}
		if publicLog.EntityID != convertedBack.EntityID {
			t.Errorf("EntityID mismatch: got %v, want %v", convertedBack.EntityID, publicLog.EntityID)
		}
		if publicLog.Period != convertedBack.Period {
			t.Errorf("Period mismatch: got %v, want %v", convertedBack.Period, publicLog.Period)
		}
		if publicLog.BatchID != convertedBack.BatchID {
			t.Errorf("BatchID mismatch: got %v, want %v", convertedBack.BatchID, publicLog.BatchID)
		}
		if publicLog.ExportHash != convertedBack.ExportHash {
			t.Errorf("ExportHash mismatch: got %v, want %v", convertedBack.ExportHash, publicLog.ExportHash)
		}
		if publicLog.FilePath != convertedBack.FilePath {
			t.Errorf("FilePath mismatch: got %v, want %v", convertedBack.FilePath, publicLog.FilePath)
		}
		if publicLog.ExportedAt != convertedBack.ExportedAt {
			t.Errorf("ExportedAt mismatch: got %v, want %v", convertedBack.ExportedAt, publicLog.ExportedAt)
		}
	})
}

func TestServiceAdapterWrapper(t *testing.T) {
	ctx := context.Background()

	called := false
	mockRepo := &mockFiscalRepository{
		listPendingEntitiesFunc: func(ctx context.Context, period string) ([]string, error) {
			called = true
			if period != "2024-01" {
				t.Errorf("period mismatch: got %v, want 2024-01", period)
			}
			return []string{"entity-1", "entity-2"}, nil
		},
	}

	service := NewDashboardService(mockRepo)

	t.Run("ListPendingEntities", func(t *testing.T) {
		entities, err := service.ListPendingEntities(ctx, "2024-01")
		if err != nil {
			t.Fatalf("ListPendingEntities should not return error: %v", err)
		}
		if !called {
			t.Error("mock repository should have been called")
		}
		expected := []string{"entity-1", "entity-2"}
		if !reflect.DeepEqual(entities, expected) {
			t.Errorf("entities mismatch: got %v, want %v", entities, expected)
		}
	})
}

func TestTypesAndInterfaces(t *testing.T) {
	t.Run("FiscalBatch fields", func(t *testing.T) {
		batch := FiscalBatch{
			ID:           "test",
			EntityID:     "entity",
			Period:       "2024-01",
			TotalEntries: 5,
			ExportHash:   "hash",
			CreatedAt:    1234567890,
		}

		if batch.ID != "test" {
			t.Errorf("ID mismatch: got %v, want test", batch.ID)
		}
		if batch.EntityID != "entity" {
			t.Errorf("EntityID mismatch: got %v, want entity", batch.EntityID)
		}
		if batch.Period != "2024-01" {
			t.Errorf("Period mismatch: got %v, want 2024-01", batch.Period)
		}
		if batch.TotalEntries != 5 {
			t.Errorf("TotalEntries mismatch: got %v, want 5", batch.TotalEntries)
		}
		if batch.ExportHash != "hash" {
			t.Errorf("ExportHash mismatch: got %v, want hash", batch.ExportHash)
		}
		if batch.CreatedAt != 1234567890 {
			t.Errorf("CreatedAt mismatch: got %v, want 1234567890", batch.CreatedAt)
		}
	})

	t.Run("EntryDTO fields", func(t *testing.T) {
		testTime := time.Now()
		entry := EntryDTO{
			ID:          1,
			EntityID:    "entity",
			Date:        testTime,
			Description: "Test",
			Postings:    []PostingDTO{},
			TotalDebit:  1000,
			TotalCredit: 1000,
		}

		if entry.ID != 1 {
			t.Errorf("ID mismatch: got %v, want 1", entry.ID)
		}
		if entry.EntityID != "entity" {
			t.Errorf("EntityID mismatch: got %v, want entity", entry.EntityID)
		}
		if !entry.Date.Equal(testTime) {
			t.Errorf("Date mismatch: got %v, want %v", entry.Date, testTime)
		}
		if entry.Description != "Test" {
			t.Errorf("Description mismatch: got %v, want Test", entry.Description)
		}
		if entry.TotalDebit != 1000 {
			t.Errorf("TotalDebit mismatch: got %v, want 1000", entry.TotalDebit)
		}
		if entry.TotalCredit != 1000 {
			t.Errorf("TotalCredit mismatch: got %v, want 1000", entry.TotalCredit)
		}
	})

	t.Run("PostingDTO fields", func(t *testing.T) {
		posting := PostingDTO{
			ID:          1,
			EntryID:     1,
			AccountID:   100,
			AccountCode: "1.01.01",
			AccountName: "Cash",
			Debit:       5000,
			Credit:      0,
		}

		if posting.ID != 1 {
			t.Errorf("ID mismatch: got %v, want 1", posting.ID)
		}
		if posting.EntryID != 1 {
			t.Errorf("EntryID mismatch: got %v, want 1", posting.EntryID)
		}
		if posting.AccountID != 100 {
			t.Errorf("AccountID mismatch: got %v, want 100", posting.AccountID)
		}
		if posting.AccountCode != "1.01.01" {
			t.Errorf("AccountCode mismatch: got %v, want 1.01.01", posting.AccountCode)
		}
		if posting.AccountName != "Cash" {
			t.Errorf("AccountName mismatch: got %v, want Cash", posting.AccountName)
		}
		if posting.Debit != 5000 {
			t.Errorf("Debit mismatch: got %v, want 5000", posting.Debit)
		}
		if posting.Credit != 0 {
			t.Errorf("Credit mismatch: got %v, want 0", posting.Credit)
		}
	})

	t.Run("AccountMapping fields", func(t *testing.T) {
		mapping := AccountMapping{
			LocalCode:    "101",
			LocalName:    "Cash",
			StandardCode: "1.01.01",
			StandardName: "Cash and Cash Equivalents",
		}

		if mapping.LocalCode != "101" {
			t.Errorf("LocalCode mismatch: got %v, want 101", mapping.LocalCode)
		}
		if mapping.LocalName != "Cash" {
			t.Errorf("LocalName mismatch: got %v, want Cash", mapping.LocalName)
		}
		if mapping.StandardCode != "1.01.01" {
			t.Errorf("StandardCode mismatch: got %v, want 1.01.01", mapping.StandardCode)
		}
		if mapping.StandardName != "Cash and Cash Equivalents" {
			t.Errorf("StandardName mismatch: got %v, want Cash and Cash Equivalents", mapping.StandardName)
		}
	})
}
