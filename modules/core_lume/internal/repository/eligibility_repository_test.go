package repository

import (
	"database/sql"
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/providentia/digna/core_lume/internal/domain"
)

// mockLifecycleManager é um mock simples do LifecycleManager
type mockLifecycleManager struct {
	db *sql.DB
}

func (m *mockLifecycleManager) GetConnection(entityID string) (*sql.DB, error) {
	return m.db, nil
}

func (m *mockLifecycleManager) GetCentralConnection() (*sql.DB, error) {
	return m.db, nil
}

func (m *mockLifecycleManager) CloseConnection(entityID string) error {
	return nil
}

func (m *mockLifecycleManager) CloseAll() error {
	return m.db.Close()
}

func (m *mockLifecycleManager) EntityExists(entityID string) (bool, error) {
	return true, nil
}

func (m *mockLifecycleManager) CreateEntity(entityID, entityName string) error {
	return nil
}

func setupTestDB(t *testing.T) (*sql.DB, func()) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Failed to open test database: %v", err)
	}

	cleanup := func() {
		db.Close()
	}

	return db, cleanup
}

func TestSQLiteEligibilityRepository_InitTable(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	mockLM := &mockLifecycleManager{db: db}
	repo := NewSQLiteEligibilityRepository(mockLM)

	err := repo.InitTable("test-entity")
	if err != nil {
		t.Errorf("InitTable() unexpected error: %v", err)
	}

	// Verify table was created by trying to query it
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name='eligibility_profiles'").Scan(&count)
	if err != nil {
		t.Errorf("Failed to verify table creation: %v", err)
	}
	if count != 1 {
		t.Errorf("Table eligibility_profiles was not created")
	}

	// Verify index was created
	err = db.QueryRow("SELECT COUNT(*) FROM sqlite_master WHERE type='index' AND name='idx_eligibility_entity'").Scan(&count)
	if err != nil {
		t.Errorf("Failed to verify index creation: %v", err)
	}
	if count != 1 {
		t.Errorf("Index idx_eligibility_entity was not created")
	}
}

func TestSQLiteEligibilityRepository_Save(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	mockLM := &mockLifecycleManager{db: db}
	repo := NewSQLiteEligibilityRepository(mockLM)

	// Initialize table
	err := repo.InitTable("test-entity")
	if err != nil {
		t.Fatalf("Failed to init table: %v", err)
	}

	now := time.Now().Unix()
	profile := &domain.EligibilityProfile{
		ID:                  "profile-1",
		EntityID:            "test-entity",
		CNPJ:                "12345678000195",
		CNAE:                "1234567",
		Municipio:           "São Paulo",
		UF:                  "SP",
		FaturamentoAnual:    1000000, // R$ 10.000,00
		RegimeTributario:    "Simples Nacional",
		DataAbertura:        now,
		SituacaoFiscal:      "Ativa",
		InscritoCadUnico:    true,
		SocioMulher:         true,
		InadimplenciaAtiva:  false,
		FinalidadeCredito:   domain.FinalidadeCapitalGiro,
		ValorNecessario:     50000, // R$ 500,00
		TipoEntidade:        domain.TipoEntidadeMEI,
		ContabilidadeFormal: true,
		PreenchidoEm:        now,
		AtualizadoEm:        now,
		PreenchidoPor:       "user-1",
		CreatedAt:           now,
		UpdatedAt:           now,
	}

	err = repo.Save(profile)
	if err != nil {
		t.Errorf("Save() unexpected error: %v", err)
	}

	// Verify saved data
	saved, err := repo.FindByEntityID("test-entity")
	if err != nil {
		t.Errorf("FindByEntityID() unexpected error: %v", err)
	}
	if saved == nil {
		t.Fatal("Saved profile is nil")
	}

	if saved.EntityID != profile.EntityID {
		t.Errorf("EntityID = %s, expected %s", saved.EntityID, profile.EntityID)
	}
	if saved.CNPJ != profile.CNPJ {
		t.Errorf("CNPJ = %s, expected %s", saved.CNPJ, profile.CNPJ)
	}
	if saved.InscritoCadUnico != profile.InscritoCadUnico {
		t.Errorf("InscritoCadUnico = %v, expected %v", saved.InscritoCadUnico, profile.InscritoCadUnico)
	}
	if saved.ValorNecessario != profile.ValorNecessario {
		t.Errorf("ValorNecessario = %d, expected %d", saved.ValorNecessario, profile.ValorNecessario)
	}
}

func TestSQLiteEligibilityRepository_Save_InvalidProfile(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	mockLM := &mockLifecycleManager{db: db}
	repo := NewSQLiteEligibilityRepository(mockLM)

	// Initialize table
	err := repo.InitTable("test-entity")
	if err != nil {
		t.Fatalf("Failed to init table: %v", err)
	}

	// Invalid profile - missing EntityID
	profile := &domain.EligibilityProfile{
		ID: "profile-1",
		// EntityID is empty
	}

	err = repo.Save(profile)
	if err == nil {
		t.Error("Save() expected error for invalid profile but got nil")
	}
}

func TestSQLiteEligibilityRepository_FindByEntityID_NotFound(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	mockLM := &mockLifecycleManager{db: db}
	repo := NewSQLiteEligibilityRepository(mockLM)

	// Initialize table
	err := repo.InitTable("test-entity")
	if err != nil {
		t.Fatalf("Failed to init table: %v", err)
	}

	// Try to find non-existent profile
	_, err = repo.FindByEntityID("non-existent")
	if err != domain.ErrProfileNotFound {
		t.Errorf("FindByEntityID() error = %v, expected ErrProfileNotFound", err)
	}
}

func TestSQLiteEligibilityRepository_Save_Update(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	mockLM := &mockLifecycleManager{db: db}
	repo := NewSQLiteEligibilityRepository(mockLM)

	// Initialize table
	err := repo.InitTable("test-entity")
	if err != nil {
		t.Fatalf("Failed to init table: %v", err)
	}

	now := time.Now().Unix()

	// Create initial profile
	profile := &domain.EligibilityProfile{
		ID:                "profile-1",
		EntityID:          "test-entity",
		FinalidadeCredito: domain.FinalidadeCapitalGiro,
		TipoEntidade:      domain.TipoEntidadeMEI,
		ValorNecessario:   50000,
		PreenchidoEm:      now,
		AtualizadoEm:      now,
		CreatedAt:         now,
		UpdatedAt:         now,
	}

	err = repo.Save(profile)
	if err != nil {
		t.Fatalf("Failed to save initial profile: %v", err)
	}

	// Update the profile
	profile.ValorNecessario = 100000
	profile.FinalidadeCredito = domain.FinalidadeEquipamento

	err = repo.Save(profile)
	if err != nil {
		t.Errorf("Save() update unexpected error: %v", err)
	}

	// Verify update
	saved, err := repo.FindByEntityID("test-entity")
	if err != nil {
		t.Errorf("FindByEntityID() unexpected error: %v", err)
	}
	if saved.ValorNecessario != 100000 {
		t.Errorf("ValorNecessario = %d, expected 100000", saved.ValorNecessario)
	}
	if saved.FinalidadeCredito != domain.FinalidadeEquipamento {
		t.Errorf("FinalidadeCredito = %v, expected EQUIPAMENTO", saved.FinalidadeCredito)
	}
}

func TestSQLiteEligibilityRepository_UpdateFields(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	mockLM := &mockLifecycleManager{db: db}
	repo := NewSQLiteEligibilityRepository(mockLM)

	// Initialize table
	err := repo.InitTable("test-entity")
	if err != nil {
		t.Fatalf("Failed to init table: %v", err)
	}

	now := time.Now().Unix()

	// Create profile
	profile := &domain.EligibilityProfile{
		ID:                "profile-1",
		EntityID:          "test-entity",
		FinalidadeCredito: domain.FinalidadeCapitalGiro,
		TipoEntidade:      domain.TipoEntidadeMEI,
		ValorNecessario:   50000,
		PreenchidoEm:      now,
		AtualizadoEm:      now,
		CreatedAt:         now,
		UpdatedAt:         now,
	}

	err = repo.Save(profile)
	if err != nil {
		t.Fatalf("Failed to save profile: %v", err)
	}

	// Update specific fields
	fields := map[string]interface{}{
		"valor_necessario": 75000,
		"cnpj":             "98765432000195",
	}

	err = repo.UpdateFields("test-entity", fields)
	if err != nil {
		t.Errorf("UpdateFields() unexpected error: %v", err)
	}

	// Verify update
	saved, err := repo.FindByEntityID("test-entity")
	if err != nil {
		t.Errorf("FindByEntityID() unexpected error: %v", err)
	}
	if saved.ValorNecessario != 75000 {
		t.Errorf("ValorNecessario = %d, expected 75000", saved.ValorNecessario)
	}
	if saved.CNPJ != "98765432000195" {
		t.Errorf("CNPJ = %s, expected 98765432000195", saved.CNPJ)
	}
	// Other fields should remain unchanged
	if saved.FinalidadeCredito != domain.FinalidadeCapitalGiro {
		t.Errorf("FinalidadeCredito was changed unexpectedly")
	}
}

func TestSQLiteEligibilityRepository_ListIncomplete(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	mockLM := &mockLifecycleManager{db: db}
	repo := NewSQLiteEligibilityRepository(mockLM)

	// Initialize table
	err := repo.InitTable("test-entity")
	if err != nil {
		t.Fatalf("Failed to init table: %v", err)
	}

	// ListIncomplete currently returns empty list (as per implementation)
	profiles, err := repo.ListIncomplete()
	if err != nil {
		t.Errorf("ListIncomplete() unexpected error: %v", err)
	}
	if len(profiles) != 0 {
		t.Errorf("ListIncomplete() returned %d profiles, expected 0", len(profiles))
	}
}

func TestSQLiteEligibilityRepository_BoolConversion(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	mockLM := &mockLifecycleManager{db: db}
	repo := NewSQLiteEligibilityRepository(mockLM)

	// Initialize table
	err := repo.InitTable("test-entity")
	if err != nil {
		t.Fatalf("Failed to init table: %v", err)
	}

	now := time.Now().Unix()

	// Test with all bools true
	profile := &domain.EligibilityProfile{
		ID:                  "profile-1",
		EntityID:            "test-entity",
		FinalidadeCredito:   domain.FinalidadeCapitalGiro,
		TipoEntidade:        domain.TipoEntidadeMEI,
		ValorNecessario:     50000,
		InscritoCadUnico:    true,
		SocioMulher:         true,
		InadimplenciaAtiva:  true,
		ContabilidadeFormal: true,
		PreenchidoEm:        now,
		AtualizadoEm:        now,
		CreatedAt:           now,
		UpdatedAt:           now,
	}

	err = repo.Save(profile)
	if err != nil {
		t.Fatalf("Failed to save profile: %v", err)
	}

	saved, err := repo.FindByEntityID("test-entity")
	if err != nil {
		t.Fatalf("Failed to find profile: %v", err)
	}

	if !saved.InscritoCadUnico {
		t.Error("InscritoCadUnico should be true")
	}
	if !saved.SocioMulher {
		t.Error("SocioMulher should be true")
	}
	if !saved.InadimplenciaAtiva {
		t.Error("InadimplenciaAtiva should be true")
	}
	if !saved.ContabilidadeFormal {
		t.Error("ContabilidadeFormal should be true")
	}

	// Test with all bools false
	profile2 := &domain.EligibilityProfile{
		ID:                  "profile-2",
		EntityID:            "test-entity-2",
		FinalidadeCredito:   domain.FinalidadeCapitalGiro,
		TipoEntidade:        domain.TipoEntidadeMEI,
		ValorNecessario:     50000,
		InscritoCadUnico:    false,
		SocioMulher:         false,
		InadimplenciaAtiva:  false,
		ContabilidadeFormal: false,
		PreenchidoEm:        now,
		AtualizadoEm:        now,
		CreatedAt:           now,
		UpdatedAt:           now,
	}

	err = repo.Save(profile2)
	if err != nil {
		t.Fatalf("Failed to save profile2: %v", err)
	}

	saved2, err := repo.FindByEntityID("test-entity-2")
	if err != nil {
		t.Fatalf("Failed to find profile2: %v", err)
	}

	if saved2.InscritoCadUnico {
		t.Error("InscritoCadUnico should be false")
	}
	if saved2.SocioMulher {
		t.Error("SocioMulher should be false")
	}
	if saved2.InadimplenciaAtiva {
		t.Error("InadimplenciaAtiva should be false")
	}
	if saved2.ContabilidadeFormal {
		t.Error("ContabilidadeFormal should be false")
	}
}
