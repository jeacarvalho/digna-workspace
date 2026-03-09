package repository

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"digna/accountant_dashboard/internal/domain"
)

func setupTestDatabase(t *testing.T, entityID string) (string, func()) {
	t.Helper()

	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, fmt.Sprintf("%s.db", entityID))

	dsn := fmt.Sprintf("%s?_journal_mode=WAL&_foreign_keys=ON", dbPath)
	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		t.Fatalf("failed to open test database: %v", err)
	}
	defer db.Close()

	schemaSQL := `
		CREATE TABLE entries (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			entry_date INTEGER NOT NULL,
			description TEXT NOT NULL,
			reference TEXT,
			created_at INTEGER NOT NULL DEFAULT (unixepoch())
		);

		CREATE TABLE accounts (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			code TEXT NOT NULL,
			name TEXT NOT NULL,
			parent_id INTEGER,
			account_type TEXT NOT NULL,
			created_at INTEGER NOT NULL DEFAULT (unixepoch()),
			FOREIGN KEY (parent_id) REFERENCES accounts(id)
		);

		CREATE TABLE postings (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			entry_id INTEGER NOT NULL,
			account_id INTEGER NOT NULL,
			amount INTEGER NOT NULL,
			direction TEXT NOT NULL CHECK (direction IN ('DEBIT', 'CREDIT')),
			created_at INTEGER NOT NULL DEFAULT (unixepoch()),
			FOREIGN KEY (entry_id) REFERENCES entries(id),
			FOREIGN KEY (account_id) REFERENCES accounts(id)
		);

		CREATE TABLE fiscal_exports (
			id TEXT PRIMARY KEY,
			entity_id TEXT NOT NULL,
			period TEXT NOT NULL,
			batch_id TEXT NOT NULL,
			export_hash TEXT NOT NULL,
			total_entries INTEGER NOT NULL,
			created_at INTEGER NOT NULL,
			UNIQUE(entity_id, period)
		);
	`

	if _, err := db.Exec(schemaSQL); err != nil {
		t.Fatalf("failed to create test schema: %v", err)
	}

	now := time.Now()
	entryDate := time.Date(2026, 3, 15, 0, 0, 0, 0, time.UTC).Unix()

	insertDataSQL := `
		INSERT INTO accounts (code, name, account_type) VALUES
		('1.1.01', 'Caixa', 'ASSET'),
		('3.1.01', 'Receita de Vendas', 'REVENUE'),
		('4.1.01', 'Despesas Administrativas', 'EXPENSE');

		INSERT INTO entries (entry_date, description, reference) VALUES
		(?, 'Venda à vista', 'REF001'),
		(?, 'Compra de material', 'REF002');

		INSERT INTO postings (entry_id, account_id, amount, direction) VALUES
		(1, 1, 1000000, 'DEBIT'),
		(1, 2, 1000000, 'CREDIT'),
		(2, 3, 500000, 'DEBIT'),
		(2, 1, 500000, 'CREDIT');

		INSERT INTO fiscal_exports (id, entity_id, period, batch_id, export_hash, total_entries, created_at) VALUES
		('export_001', ?, '2026-02', 'batch_001', 'hash_001', 50, ?);
	`

	if _, err := db.Exec(insertDataSQL, entryDate, entryDate, entityID, now.Unix()); err != nil {
		t.Fatalf("failed to insert test data: %v", err)
	}

	cleanup := func() {
		os.RemoveAll(tempDir)
	}

	return tempDir, cleanup
}

func TestNewSQLiteFiscalAdapter(t *testing.T) {
	adapter := NewSQLiteFiscalAdapter()
	if adapter == nil {
		t.Error("NewSQLiteFiscalAdapter should return non-nil adapter")
	}

	if adapter.BasePath() != DataDir {
		t.Errorf("BasePath = %s, want %s", adapter.BasePath(), DataDir)
	}
}

func TestSQLiteFiscalAdapter_SetBasePath(t *testing.T) {
	adapter := NewSQLiteFiscalAdapter()
	newPath := "/tmp/test_path"
	adapter.SetBasePath(newPath)

	if adapter.BasePath() != newPath {
		t.Errorf("BasePath = %s, want %s", adapter.BasePath(), newPath)
	}
}

func TestSQLiteFiscalAdapter_LoadEntries(t *testing.T) {
	entityID := "test_entity_001"
	tempDir, cleanup := setupTestDatabase(t, entityID)
	defer cleanup()

	adapter := NewSQLiteFiscalAdapter()
	adapter.SetBasePath(tempDir)

	ctx := context.Background()
	period := "2026-03"

	entries, err := adapter.LoadEntries(ctx, entityID, period)
	if err != nil {
		t.Fatalf("LoadEntries failed: %v", err)
	}

	if len(entries) != 2 {
		t.Errorf("LoadEntries returned %d entries, want 2", len(entries))
	}

	for _, entry := range entries {
		if entry.EntityID != entityID {
			t.Errorf("Entry EntityID = %s, want %s", entry.EntityID, entityID)
		}

		if len(entry.Postings) == 0 {
			t.Error("Entry should have postings")
		}

		if entry.TotalDebit != entry.TotalCredit {
			t.Errorf("Entry %d: TotalDebit (%d) != TotalCredit (%d)", entry.ID, entry.TotalDebit, entry.TotalCredit)
		}
	}

	entry1 := entries[0]
	if entry1.Description != "Venda à vista" {
		t.Errorf("Entry 1 Description = %s, want 'Venda à vista'", entry1.Description)
	}

	if entry1.TotalDebit != 1000000 || entry1.TotalCredit != 1000000 {
		t.Errorf("Entry 1 amounts: Debit=%d, Credit=%d, want 1000000 each", entry1.TotalDebit, entry1.TotalCredit)
	}
}

func TestSQLiteFiscalAdapter_LoadEntries_InvalidPeriod(t *testing.T) {
	entityID := "test_entity_002"
	tempDir, cleanup := setupTestDatabase(t, entityID)
	defer cleanup()

	adapter := NewSQLiteFiscalAdapter()
	adapter.SetBasePath(tempDir)

	ctx := context.Background()
	invalidPeriod := "2026-03-15"

	_, err := adapter.LoadEntries(ctx, entityID, invalidPeriod)
	if err == nil {
		t.Error("LoadEntries should fail with invalid period format")
	}
}

func TestSQLiteFiscalAdapter_LoadEntries_NoDatabase(t *testing.T) {
	adapter := NewSQLiteFiscalAdapter()
	adapter.SetBasePath("/tmp/nonexistent")

	ctx := context.Background()
	period := "2026-03"

	_, err := adapter.LoadEntries(ctx, "nonexistent_entity", period)
	if err == nil {
		t.Error("LoadEntries should fail when database doesn't exist")
	}
}

func TestSQLiteFiscalAdapter_RegisterExport(t *testing.T) {
	entityID := "test_entity_003"
	tempDir, cleanup := setupTestDatabase(t, entityID)
	defer cleanup()

	adapter := NewSQLiteFiscalAdapter()
	adapter.SetBasePath(tempDir)

	ctx := context.Background()
	batch := &domain.FiscalBatch{
		ID:           "batch_2026-03_001",
		EntityID:     entityID,
		Period:       "2026-03",
		ExportHash:   "test_hash_123",
		TotalEntries: 2,
		CreatedAt:    time.Now().Unix(),
	}

	err := adapter.RegisterExport(ctx, entityID, batch)
	if err != nil {
		t.Fatalf("RegisterExport failed: %v", err)
	}

	db, err := adapter.openReadOnly(entityID)
	if err != nil {
		t.Fatalf("Failed to open database for verification: %v", err)
	}
	defer db.Close()

	var count int
	err = db.QueryRowContext(ctx, "SELECT COUNT(*) FROM fiscal_exports WHERE period = ?", batch.Period).Scan(&count)
	if err != nil {
		t.Fatalf("Failed to query exports: %v", err)
	}

	if count != 1 {
		t.Errorf("Expected 1 export record, got %d", count)
	}
}

func TestSQLiteFiscalAdapter_ListPendingEntities(t *testing.T) {
	entity1 := "entity_001"
	entity2 := "entity_002"
	entity3 := "entity_003"

	tempDir := t.TempDir()

	createTestDB := func(entityID string, hasEntries bool, hasExport bool) {
		dbPath := filepath.Join(tempDir, fmt.Sprintf("%s.db", entityID))
		dsn := fmt.Sprintf("%s?_journal_mode=WAL&_foreign_keys=ON", dbPath)
		db, err := sql.Open("sqlite3", dsn)
		if err != nil {
			t.Fatalf("failed to create test database for %s: %v", entityID, err)
		}
		defer db.Close()

		schemaSQL := `
			CREATE TABLE IF NOT EXISTS entries (
				id INTEGER PRIMARY KEY AUTOINCREMENT,
				entry_date INTEGER NOT NULL,
				description TEXT NOT NULL,
				reference TEXT,
				created_at INTEGER NOT NULL DEFAULT (unixepoch())
			);

			CREATE TABLE IF NOT EXISTS fiscal_exports (
				id TEXT PRIMARY KEY,
				entity_id TEXT NOT NULL,
				period TEXT NOT NULL,
				batch_id TEXT NOT NULL,
				export_hash TEXT NOT NULL,
				total_entries INTEGER NOT NULL,
				created_at INTEGER NOT NULL,
				UNIQUE(entity_id, period)
			);
		`

		if _, err := db.Exec(schemaSQL); err != nil {
			t.Fatalf("failed to create schema for %s: %v", entityID, err)
		}

		if hasEntries {
			entryDate := time.Date(2026, 3, 15, 0, 0, 0, 0, time.UTC).Unix()
			if _, err := db.Exec("INSERT INTO entries (entry_date, description) VALUES (?, 'Test entry')", entryDate); err != nil {
				t.Fatalf("failed to insert entry for %s: %v", entityID, err)
			}
		}

		if hasExport {
			now := time.Now().Unix()
			if _, err := db.Exec(
				"INSERT INTO fiscal_exports (id, entity_id, period, batch_id, export_hash, total_entries, created_at) VALUES (?, ?, ?, ?, ?, ?, ?)",
				fmt.Sprintf("export_%s", entityID), entityID, "2026-03", "batch_001", "hash_001", 1, now,
			); err != nil {
				t.Fatalf("failed to insert export for %s: %v", entityID, err)
			}
		}
	}

	createTestDB(entity1, true, false)
	createTestDB(entity2, true, true)
	createTestDB(entity3, false, false)

	adapter := NewSQLiteFiscalAdapter()
	adapter.SetBasePath(tempDir)

	ctx := context.Background()
	testPeriod := "2026-03"

	pending, err := adapter.ListPendingEntities(ctx, testPeriod)
	if err != nil {
		t.Fatalf("ListPendingEntities failed: %v", err)
	}

	if len(pending) != 1 {
		t.Errorf("ListPendingEntities returned %d entities, want 1", len(pending))
	}

	if len(pending) > 0 && pending[0] != entity1 {
		t.Errorf("ListPendingEntities returned %s, want %s", pending[0], entity1)
	}
}

func TestSQLiteFiscalAdapter_GetExportHistory(t *testing.T) {
	entityID := "test_entity_004"
	tempDir, cleanup := setupTestDatabase(t, entityID)
	defer cleanup()

	adapter := NewSQLiteFiscalAdapter()
	adapter.SetBasePath(tempDir)

	ctx := context.Background()
	period := "2026-02"

	history, err := adapter.GetExportHistory(ctx, entityID, period)
	if err != nil {
		t.Fatalf("GetExportHistory failed: %v", err)
	}

	if len(history) != 1 {
		t.Errorf("GetExportHistory returned %d logs, want 1", len(history))
	}

	if len(history) > 0 {
		log := history[0]
		if log.EntityID != entityID {
			t.Errorf("Export log EntityID = %s, want %s", log.EntityID, entityID)
		}
		if log.Period != period {
			t.Errorf("Export log Period = %s, want %s", log.Period, period)
		}
		if log.BatchID != "batch_001" {
			t.Errorf("Export log BatchID = %s, want batch_001", log.BatchID)
		}
	}
}

func TestSQLiteFiscalAdapter_GetExportHistory_NoExports(t *testing.T) {
	entityID := "test_entity_005"
	tempDir, cleanup := setupTestDatabase(t, entityID)
	defer cleanup()

	adapter := NewSQLiteFiscalAdapter()
	adapter.SetBasePath(tempDir)

	ctx := context.Background()
	period := "2026-04"

	history, err := adapter.GetExportHistory(ctx, entityID, period)
	if err != nil {
		t.Fatalf("GetExportHistory failed: %v", err)
	}

	if len(history) != 0 {
		t.Errorf("GetExportHistory returned %d logs, want 0", len(history))
	}
}

func TestSQLiteFiscalAdapter_openReadOnly(t *testing.T) {
	entityID := "test_entity_006"
	tempDir, cleanup := setupTestDatabase(t, entityID)
	defer cleanup()

	adapter := NewSQLiteFiscalAdapter()
	adapter.SetBasePath(tempDir)

	db, err := adapter.openReadOnly(entityID)
	if err != nil {
		t.Fatalf("openReadOnly failed: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		t.Errorf("Database ping failed: %v", err)
	}
}

func TestSQLiteFiscalAdapter_openReadOnly_NoDatabase(t *testing.T) {
	adapter := NewSQLiteFiscalAdapter()
	adapter.SetBasePath("/tmp/nonexistent")

	_, err := adapter.openReadOnly("nonexistent_entity")
	if err == nil {
		t.Error("openReadOnly should fail when database doesn't exist")
	}
}
