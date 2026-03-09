package handler

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"digna/accountant_dashboard/internal/domain"
)

type mockTranslator struct {
	listPendingEntitiesFunc func(ctx context.Context, period string) ([]string, error)
	getExportHistoryFunc    func(ctx context.Context, entityID, period string) ([]domain.FiscalExportLog, error)
	translateAndExportFunc  func(ctx context.Context, entityID, period string) (*domain.FiscalBatch, []byte, error)
}

func (m *mockTranslator) ListPendingEntities(ctx context.Context, period string) ([]string, error) {
	if m.listPendingEntitiesFunc != nil {
		return m.listPendingEntitiesFunc(ctx, period)
	}
	return []string{}, nil
}

func (m *mockTranslator) GetExportHistory(ctx context.Context, entityID, period string) ([]domain.FiscalExportLog, error) {
	if m.getExportHistoryFunc != nil {
		return m.getExportHistoryFunc(ctx, entityID, period)
	}
	return []domain.FiscalExportLog{}, nil
}

func (m *mockTranslator) TranslateAndExport(ctx context.Context, entityID, period string) (*domain.FiscalBatch, []byte, error) {
	if m.translateAndExportFunc != nil {
		return m.translateAndExportFunc(ctx, entityID, period)
	}
	return &domain.FiscalBatch{}, []byte{}, nil
}

type mockMapper struct {
	getMappingFunc     func(localCode string) (domain.AccountMapping, bool)
	getAllMappingsFunc func() []domain.AccountMapping
}

func (m *mockMapper) GetMapping(localCode string) (domain.AccountMapping, bool) {
	if m.getMappingFunc != nil {
		return m.getMappingFunc(localCode)
	}
	return domain.AccountMapping{}, false
}

func (m *mockMapper) GetAllMappings() []domain.AccountMapping {
	if m.getAllMappingsFunc != nil {
		return m.getAllMappingsFunc()
	}
	return []domain.AccountMapping{}
}

func TestNewDashboardHandler(t *testing.T) {
	translator := &mockTranslator{}
	mapper := &mockMapper{}

	handler := NewDashboardHandler(translator, mapper)
	if handler == nil {
		t.Error("NewDashboardHandler should return non-nil handler")
	}

	handler = NewDashboardHandler(translator, mapper)
	if handler == nil {
		t.Error("NewDashboardHandler should return non-nil handler with mock translator")
	}
}

func TestDashboardHandler_Dashboard(t *testing.T) {
	translator := &mockTranslator{
		listPendingEntitiesFunc: func(ctx context.Context, period string) ([]string, error) {
			return []string{"entity_001", "entity_002"}, nil
		},
		getExportHistoryFunc: func(ctx context.Context, entityID, period string) ([]domain.FiscalExportLog, error) {
			if entityID == "entity_001" {
				return []domain.FiscalExportLog{
					{
						ID:         "export_001",
						EntityID:   "entity_001",
						Period:     "2026-03",
						BatchID:    "batch_001",
						ExportHash: "hash_001",
						ExportedAt: time.Now().Unix(),
					},
				}, nil
			}
			return []domain.FiscalExportLog{}, nil
		},
	}

	mapper := &mockMapper{
		getAllMappingsFunc: func() []domain.AccountMapping {
			return []domain.AccountMapping{
				{
					LocalCode:    "1.1.01",
					LocalName:    "Caixa",
					StandardCode: "1.1.01.00.00",
					StandardName: "Disponibilidades - Caixa",
				},
			}
		},
	}

	handler := NewDashboardHandler(translator, mapper)

	req := httptest.NewRequest("GET", "/dashboard", nil)
	w := httptest.NewRecorder()

	handler.Dashboard(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Dashboard returned status %d, want %d", resp.StatusCode, http.StatusOK)
	}

	body, _ := io.ReadAll(resp.Body)
	bodyStr := string(body)

	if !strings.Contains(bodyStr, "Painel do Contador Social") {
		t.Error("Dashboard response should contain page title")
	}

	if !strings.Contains(bodyStr, "entity_001") {
		t.Error("Dashboard response should contain entity IDs")
	}

	if !strings.Contains(bodyStr, "1.1.01") {
		t.Error("Dashboard response should contain account mappings")
	}
}

func TestDashboardHandler_Dashboard_WithPeriod(t *testing.T) {
	translator := &mockTranslator{
		listPendingEntitiesFunc: func(ctx context.Context, period string) ([]string, error) {
			if period == "2026-03" {
				return []string{"entity_001"}, nil
			}
			return []string{}, nil
		},
	}

	mapper := &mockMapper{
		getAllMappingsFunc: func() []domain.AccountMapping {
			return []domain.AccountMapping{}
		},
	}

	handler := NewDashboardHandler(translator, mapper)

	req := httptest.NewRequest("GET", "/dashboard?period=2026-03", nil)
	w := httptest.NewRecorder()

	handler.Dashboard(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Dashboard with period returned status %d, want %d", resp.StatusCode, http.StatusOK)
	}
}

func TestDashboardHandler_Dashboard_TranslatorError(t *testing.T) {
	translator := &mockTranslator{
		listPendingEntitiesFunc: func(ctx context.Context, period string) ([]string, error) {
			return nil, context.Canceled
		},
	}

	mapper := &mockMapper{}

	handler := NewDashboardHandler(translator, mapper)

	req := httptest.NewRequest("GET", "/dashboard", nil)
	w := httptest.NewRecorder()

	handler.Dashboard(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusInternalServerError {
		t.Errorf("Dashboard with translator error returned status %d, want %d", resp.StatusCode, http.StatusInternalServerError)
	}
}

func TestDashboardHandler_Dashboard_TemplateError(t *testing.T) {
	translator := &mockTranslator{
		listPendingEntitiesFunc: func(ctx context.Context, period string) ([]string, error) {
			return []string{"entity_001"}, nil
		},
	}

	mapper := &mockMapper{
		getAllMappingsFunc: func() []domain.AccountMapping {
			return []domain.AccountMapping{}
		},
	}

	handler := NewDashboardHandler(translator, mapper)

	req := httptest.NewRequest("GET", "/dashboard", nil)
	w := httptest.NewRecorder()

	handler.Dashboard(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Dashboard with valid data returned status %d, want %d", resp.StatusCode, http.StatusOK)
	}
}

func TestDashboardHandler_ExportFiscal(t *testing.T) {
	translator := &mockTranslator{
		translateAndExportFunc: func(ctx context.Context, entityID, period string) (*domain.FiscalBatch, []byte, error) {
			batch := &domain.FiscalBatch{
				ID:           "batch_2026-03_001",
				EntityID:     entityID,
				Period:       period,
				ExportHash:   "test_hash_123",
				TotalEntries: 10,
				CreatedAt:    time.Now().Unix(),
			}
			data := []byte("entry_date,description,amount\n2026-03-15,Test,1000000\n")
			return batch, data, nil
		},
	}

	mapper := &mockMapper{}

	handler := NewDashboardHandler(translator, mapper)

	req := httptest.NewRequest("GET", "/export?entity_id=entity_001&period=2026-03", nil)
	w := httptest.NewRecorder()

	handler.ExportFiscal(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("ExportFiscal returned status %d, want %d", resp.StatusCode, http.StatusOK)
	}

	contentType := resp.Header.Get("Content-Type")
	if contentType != "text/csv" {
		t.Errorf("ExportFiscal Content-Type = %s, want text/csv", contentType)
	}

	contentDisposition := resp.Header.Get("Content-Disposition")
	expectedDisposition := "attachment; filename=fiscal_entity_001_2026-03.csv"
	if contentDisposition != expectedDisposition {
		t.Errorf("ExportFiscal Content-Disposition = %s, want %s", contentDisposition, expectedDisposition)
	}

	exportHash := resp.Header.Get("X-Export-Hash")
	if exportHash != "test_hash_123" {
		t.Errorf("ExportFiscal X-Export-Hash = %s, want test_hash_123", exportHash)
	}

	body, _ := io.ReadAll(resp.Body)
	if !bytes.Contains(body, []byte("entry_date,description,amount")) {
		t.Error("ExportFiscal response should contain CSV data")
	}
}

func TestDashboardHandler_ExportFiscal_MissingParams(t *testing.T) {
	translator := &mockTranslator{}
	mapper := &mockMapper{}

	handler := NewDashboardHandler(translator, mapper)

	testCases := []struct {
		name     string
		url      string
		wantCode int
	}{
		{
			name:     "missing both params",
			url:      "/export",
			wantCode: http.StatusBadRequest,
		},
		{
			name:     "missing entity_id",
			url:      "/export?period=2026-03",
			wantCode: http.StatusBadRequest,
		},
		{
			name:     "missing period",
			url:      "/export?entity_id=entity_001",
			wantCode: http.StatusBadRequest,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", tc.url, nil)
			w := httptest.NewRecorder()

			handler.ExportFiscal(w, req)

			resp := w.Result()
			defer resp.Body.Close()

			if resp.StatusCode != tc.wantCode {
				t.Errorf("ExportFiscal(%s) returned status %d, want %d", tc.name, resp.StatusCode, tc.wantCode)
			}
		})
	}
}

func TestDashboardHandler_ExportFiscal_TranslatorError(t *testing.T) {
	translator := &mockTranslator{
		translateAndExportFunc: func(ctx context.Context, entityID, period string) (*domain.FiscalBatch, []byte, error) {
			return nil, nil, context.Canceled
		},
	}

	mapper := &mockMapper{}

	handler := NewDashboardHandler(translator, mapper)

	req := httptest.NewRequest("GET", "/export?entity_id=entity_001&period=2026-03", nil)
	w := httptest.NewRecorder()

	handler.ExportFiscal(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusInternalServerError {
		t.Errorf("ExportFiscal with translator error returned status %d, want %d", resp.StatusCode, http.StatusInternalServerError)
	}
}

func TestDashboardPageData(t *testing.T) {
	data := DashboardPageData{
		Title:  "Test Dashboard",
		Period: "2026-03",
		Entities: []EntityInfo{
			{
				ID:           "entity_001",
				Name:         "Test Entity",
				Status:       "PENDING",
				PendingMonth: "2026-03",
				HasExports:   true,
			},
		},
		Mappings: []domain.AccountMapping{
			{
				LocalCode:    "1.1.01",
				LocalName:    "Caixa",
				StandardCode: "1.1.01.00.00",
				StandardName: "Disponibilidades - Caixa",
			},
		},
		ExportHistory: []ExportHistoryItem{
			{
				EntityID:   "entity_001",
				Period:     "2026-03",
				ExportedAt: "2026-03-15",
				EntryCount: 10,
				Hash:       "test_hash",
			},
		},
	}

	if data.Title != "Test Dashboard" {
		t.Errorf("DashboardPageData Title = %s, want Test Dashboard", data.Title)
	}

	if len(data.Entities) != 1 {
		t.Errorf("DashboardPageData has %d entities, want 1", len(data.Entities))
	}

	if len(data.Mappings) != 1 {
		t.Errorf("DashboardPageData has %d mappings, want 1", len(data.Mappings))
	}

	if len(data.ExportHistory) != 1 {
		t.Errorf("DashboardPageData has %d export history items, want 1", len(data.ExportHistory))
	}
}

func TestEntityInfo(t *testing.T) {
	entity := EntityInfo{
		ID:           "entity_001",
		Name:         "Test Entity",
		Status:       "PENDING",
		PendingMonth: "2026-03",
		HasExports:   true,
	}

	if entity.ID != "entity_001" {
		t.Errorf("EntityInfo ID = %s, want entity_001", entity.ID)
	}

	if entity.Status != "PENDING" {
		t.Errorf("EntityInfo Status = %s, want PENDING", entity.Status)
	}

	if !entity.HasExports {
		t.Error("EntityInfo HasExports = false, want true")
	}
}

func TestExportHistoryItem(t *testing.T) {
	item := ExportHistoryItem{
		EntityID:   "entity_001",
		Period:     "2026-03",
		ExportedAt: "2026-03-15",
		EntryCount: 10,
		Hash:       "test_hash",
	}

	if item.EntityID != "entity_001" {
		t.Errorf("ExportHistoryItem EntityID = %s, want entity_001", item.EntityID)
	}

	if item.EntryCount != 10 {
		t.Errorf("ExportHistoryItem EntryCount = %d, want 10", item.EntryCount)
	}
}
