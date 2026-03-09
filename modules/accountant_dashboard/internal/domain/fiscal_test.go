package domain

import (
	"testing"
)

func TestDefaultAccountMapper_GetMapping(t *testing.T) {
	mapper := NewDefaultAccountMapper()

	tests := []struct {
		localCode string
		wantFound bool
		wantCode  string
		wantName  string
	}{
		{"1.1.01", true, "1.1.01.00.00", "Disponibilidades - Caixa"},
		{"1.1.02", true, "1.1.02.00.00", "Disponibilidades - Bancos Conta Movimento"},
		{"3.1.01", true, "3.1.01.00.00", "Receita Bruta de Vendas"},
		{"9.9.99", false, "", ""},
	}

	for _, tt := range tests {
		t.Run(tt.localCode, func(t *testing.T) {
			mapping, found := mapper.GetMapping(tt.localCode)
			if found != tt.wantFound {
				t.Errorf("GetMapping(%s) found = %v, want %v", tt.localCode, found, tt.wantFound)
			}
			if found {
				if mapping.StandardCode != tt.wantCode {
					t.Errorf("GetMapping(%s) StandardCode = %v, want %v", tt.localCode, mapping.StandardCode, tt.wantCode)
				}
				if mapping.StandardName != tt.wantName {
					t.Errorf("GetMapping(%s) StandardName = %v, want %v", tt.localCode, mapping.StandardName, tt.wantName)
				}
			}
		})
	}
}

func TestDefaultAccountMapper_GetAllMappings(t *testing.T) {
	mapper := NewDefaultAccountMapper()
	mappings := mapper.GetAllMappings()

	if len(mappings) == 0 {
		t.Error("GetAllMappings should return at least one mapping")
	}

	if mappings[0].LocalCode != "1.1.01" {
		t.Errorf("First mapping should be 1.1.01, got %s", mappings[0].LocalCode)
	}
}

func TestFiscalBatch_TotalEntries(t *testing.T) {
	batch := FiscalBatch{
		ID:           "test_2026-03_123",
		EntityID:     "test_entity",
		Period:       "2026-03",
		TotalEntries: 100,
		ExportHash:   "abc123",
		CreatedAt:    1234567890,
	}

	if batch.TotalEntries != 100 {
		t.Errorf("TotalEntries = %d, want 100", batch.TotalEntries)
	}

	if batch.EntityID != "test_entity" {
		t.Errorf("EntityID = %s, want test_entity", batch.EntityID)
	}
}
