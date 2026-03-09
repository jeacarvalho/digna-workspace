package accountant_dashboard_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"digna/accountant_dashboard/internal/domain"
	"digna/accountant_dashboard/internal/service"
	"digna/accountant_dashboard/pkg/dashboard"
)

func TestIntegration_AccountantDashboard(t *testing.T) {
	// Create a mock repository for testing
	mockRepo := &mockFiscalRepository{
		entries: []dashboard.EntryDTO{
			{
				ID:          1,
				EntityID:    "test_entity_001",
				Date:        time.Now(),
				Description: "Test entry",
				Postings: []dashboard.PostingDTO{
					{
						ID:          1,
						EntryID:     1,
						AccountID:   1,
						AccountCode: "1.1.01",
						AccountName: "Caixa",
						Debit:       100000,
						Credit:      0,
					},
					{
						ID:          2,
						EntryID:     1,
						AccountID:   2,
						AccountCode: "3.1.01",
						AccountName: "Receita de Vendas",
						Debit:       0,
						Credit:      100000,
					},
				},
				TotalDebit:  100000,
				TotalCredit: 100000,
			},
		},
	}

	// Create service with default mapper
	svc := dashboard.NewDashboardService(mockRepo)

	ctx := context.Background()
	period := time.Now().Format("2006-01")
	entityID := "test_entity_001"

	// Test 1: List pending entities
	t.Run("ListPendingEntities", func(t *testing.T) {
		pending, err := svc.ListPendingEntities(ctx, period)
		if err != nil {
			t.Errorf("ListPendingEntities failed: %v", err)
		}

		// Should find our test entity
		found := false
		for _, id := range pending {
			if id == entityID {
				found = true
				break
			}
		}

		if !found {
			t.Errorf("Expected to find entity %s in pending list, got: %v", entityID, pending)
		}
	})

	// Test 2: Get export history (should be empty initially)
	t.Run("GetExportHistory_Empty", func(t *testing.T) {
		history, err := svc.GetExportHistory(ctx, entityID, period)
		if err != nil {
			t.Errorf("GetExportHistory failed: %v", err)
		}

		if len(history) != 0 {
			t.Errorf("Expected empty export history, got %d entries", len(history))
		}
	})

	// Test 3: Translate and export
	t.Run("TranslateAndExport", func(t *testing.T) {
		batch, data, err := svc.TranslateAndExport(ctx, entityID, period)
		if err != nil {
			t.Errorf("TranslateAndExport failed: %v", err)
		}

		if batch == nil {
			t.Error("Expected non-nil batch")
		}

		if batch.EntityID != entityID {
			t.Errorf("Expected entity ID %s, got %s", entityID, batch.EntityID)
		}

		if batch.Period != period {
			t.Errorf("Expected period %s, got %s", period, batch.Period)
		}

		if batch.TotalEntries != 1 {
			t.Errorf("Expected 1 entry, got %d", batch.TotalEntries)
		}

		if len(data) == 0 {
			t.Error("Expected non-empty export data")
		}

		// Verify CSV format
		csvStr := string(data)
		if !contains(csvStr, "Data") || !contains(csvStr, "Conta") || !contains(csvStr, "Valor") {
			t.Error("CSV missing expected headers")
		}

		// Check for data rows
		if !contains(csvStr, "1.1.01") || !contains(csvStr, "3.1.01") {
			t.Error("CSV missing expected account codes")
		}
	})

	// Test 4: Get export history after export
	t.Run("GetExportHistory_AfterExport", func(t *testing.T) {
		history, err := svc.GetExportHistory(ctx, entityID, period)
		if err != nil {
			t.Errorf("GetExportHistory failed: %v", err)
		}

		if len(history) != 1 {
			t.Errorf("Expected 1 export history entry, got %d", len(history))
		}

		entry := history[0]
		if entry.EntityID != entityID {
			t.Errorf("Expected entity ID %s, got %s", entityID, entry.EntityID)
		}

		if entry.Period != period {
			t.Errorf("Expected period %s, got %s", period, entry.Period)
		}
	})
}

func TestIntegration_ServiceLayer(t *testing.T) {
	// Test service layer directly with mock repository
	mockRepo := &mockRepository{
		entries: []domain.EntryDTO{
			{
				ID:          1,
				EntityID:    "test_entity",
				Date:        time.Now(),
				Description: "Test entry",
				Postings: []domain.PostingDTO{
					{
						ID:          1,
						EntryID:     1,
						AccountID:   1,
						AccountCode: "1.1.01",
						AccountName: "Caixa",
						Debit:       100000,
						Credit:      0,
					},
					{
						ID:          2,
						EntryID:     1,
						AccountID:   2,
						AccountCode: "3.1.01",
						AccountName: "Receita de Vendas",
						Debit:       0,
						Credit:      100000,
					},
				},
				TotalDebit:  100000,
				TotalCredit: 100000,
			},
		},
	}

	mockMapper := &mockMapper{
		mappings: []domain.AccountMapping{
			{
				LocalCode:    "1.1.01",
				LocalName:    "Caixa",
				StandardCode: "1.1.01.00.00",
				StandardName: "Disponibilidades - Caixa",
			},
			{
				LocalCode:    "3.1.01",
				LocalName:    "Receita de Vendas",
				StandardCode: "3.1.01.00.00",
				StandardName: "Receita Bruta de Vendas",
			},
		},
	}

	translator := service.NewTranslatorService(mockRepo, mockMapper)

	ctx := context.Background()
	entityID := "test_entity"
	period := "2024-01"

	// Test translation with entries
	t.Run("TranslateAndExport_WithEntries", func(t *testing.T) {
		batch, data, err := translator.TranslateAndExport(ctx, entityID, period)
		if err != nil {
			t.Fatalf("TranslateAndExport failed: %v", err)
		}

		if batch == nil {
			t.Error("Expected non-nil batch")
		}

		if batch.EntityID != entityID {
			t.Errorf("Expected entity ID %s, got %s", entityID, batch.EntityID)
		}

		if batch.Period != period {
			t.Errorf("Expected period %s, got %s", period, batch.Period)
		}

		if batch.TotalEntries != 1 {
			t.Errorf("Expected 1 entry, got %d", batch.TotalEntries)
		}

		if len(data) == 0 {
			t.Error("Expected non-empty export data")
		}

		// Verify CSV format
		csvStr := string(data)
		if !contains(csvStr, "Data") || !contains(csvStr, "Conta") || !contains(csvStr, "Valor") {
			t.Error("CSV missing expected headers")
		}

		// Check for data rows
		if !contains(csvStr, "1.1.01") || !contains(csvStr, "3.1.01") {
			t.Error("CSV missing expected account codes")
		}
	})

	// Test translation to standard format
	t.Run("TranslateToStandardFormat", func(t *testing.T) {
		data, err := translator.TranslateToStandardFormat(mockRepo.entries)
		if err != nil {
			t.Fatalf("TranslateToStandardFormat failed: %v", err)
		}

		if len(data) == 0 {
			t.Error("Expected non-empty data")
		}

		csvStr := string(data)
		if !contains(csvStr, "Data") || !contains(csvStr, "Conta") || !contains(csvStr, "Valor") {
			t.Error("CSV missing expected headers")
		}
	})
}

// Helper function
func contains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// Mock repository for testing (implements dashboard.FiscalRepository)
type mockFiscalRepository struct {
	entries []dashboard.EntryDTO
	logs    []dashboard.FiscalExportLog
	err     error
}

func (m *mockFiscalRepository) LoadEntries(ctx context.Context, entityID, period string) ([]dashboard.EntryDTO, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.entries, nil
}

func (m *mockFiscalRepository) RegisterExport(ctx context.Context, entityID string, batch *dashboard.FiscalBatch) error {
	if m.err != nil {
		return m.err
	}
	m.logs = append(m.logs, dashboard.FiscalExportLog{
		ID:         fmt.Sprintf("log_%d", len(m.logs)+1),
		EntityID:   entityID,
		Period:     batch.Period,
		BatchID:    batch.ID,
		ExportHash: batch.ExportHash,
		FilePath:   "",
		ExportedAt: time.Now().Unix(),
	})
	return nil
}

func (m *mockFiscalRepository) ListPendingEntities(ctx context.Context, period string) ([]string, error) {
	if m.err != nil {
		return nil, m.err
	}
	if len(m.entries) > 0 {
		return []string{"test_entity_001"}, nil
	}
	return []string{}, nil
}

func (m *mockFiscalRepository) GetExportHistory(ctx context.Context, entityID, period string) ([]dashboard.FiscalExportLog, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.logs, nil
}

// Mock repository for testing (implements internal domain repository)
type mockRepository struct {
	entries []domain.EntryDTO
	logs    []domain.FiscalExportLog
	err     error
}

func (m *mockRepository) LoadEntries(ctx context.Context, entityID, period string) ([]domain.EntryDTO, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.entries, nil
}

func (m *mockRepository) RegisterExport(ctx context.Context, entityID string, batch *domain.FiscalBatch) error {
	if m.err != nil {
		return m.err
	}
	m.logs = append(m.logs, domain.FiscalExportLog{
		ID:         fmt.Sprintf("log_%d", len(m.logs)+1),
		EntityID:   entityID,
		Period:     batch.Period,
		BatchID:    batch.ID,
		ExportHash: batch.ExportHash,
		ExportedAt: time.Now().Unix(),
	})
	return nil
}

func (m *mockRepository) ListPendingEntities(ctx context.Context, period string) ([]string, error) {
	if m.err != nil {
		return nil, m.err
	}
	if len(m.entries) > 0 {
		return []string{"test_entity"}, nil
	}
	return []string{}, nil
}

func (m *mockRepository) GetExportHistory(ctx context.Context, entityID, period string) ([]domain.FiscalExportLog, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.logs, nil
}

// Mock mapper for testing
type mockMapper struct {
	mappings []domain.AccountMapping
}

func (m *mockMapper) GetMapping(localCode string) (domain.AccountMapping, bool) {
	for _, mapping := range m.mappings {
		if mapping.LocalCode == localCode {
			return mapping, true
		}
	}
	return domain.AccountMapping{}, false
}

func (m *mockMapper) GetAllMappings() []domain.AccountMapping {
	return m.mappings
}
