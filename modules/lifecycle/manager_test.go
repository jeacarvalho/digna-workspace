package lifecycle_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/providentia/digna/lifecycle/pkg/lifecycle"
)

func TestSQLiteManager_CreatesDatabaseFile(t *testing.T) {
	dataDir := "../../data/entities"
	defer os.RemoveAll(dataDir)

	sqliteMgr := lifecycle.NewSQLiteManager()
	defer sqliteMgr.CloseAll()

	entityID := "cooperativa_mel"
	db, err := sqliteMgr.GetConnection(entityID)
	if err != nil {
		t.Fatalf("failed to get connection: %v", err)
	}
	defer db.Close()

	expectedPath := filepath.Join(dataDir, "cooperativa_mel.db")
	if _, err := os.Stat(expectedPath); os.IsNotExist(err) {
		t.Errorf("database file was not created at expected path: %s", expectedPath)
	}
}

func TestSQLiteManager_WorkLogsTableExists(t *testing.T) {
	dataDir := "../../data/entities"
	defer os.RemoveAll(dataDir)

	sqliteMgr := lifecycle.NewSQLiteManager()
	defer sqliteMgr.CloseAll()

	entityID := "test_entity_worklogs"
	db, err := sqliteMgr.GetConnection(entityID)
	if err != nil {
		t.Fatalf("failed to get connection: %v", err)
	}
	defer db.Close()

	var tableName string
	query := `SELECT name FROM sqlite_master WHERE type='table' AND name='work_logs'`
	err = db.QueryRow(query).Scan(&tableName)
	if err != nil {
		t.Fatalf("work_logs table does not exist: %v", err)
	}

	if tableName != "work_logs" {
		t.Errorf("expected table name 'work_logs', got '%s'", tableName)
	}
}

func TestSQLiteManager_AllTablesExist(t *testing.T) {
	dataDir := "../../data/entities"
	defer os.RemoveAll(dataDir)

	sqliteMgr := lifecycle.NewSQLiteManager()
	defer sqliteMgr.CloseAll()

	entityID := "test_entity_all_tables"
	db, err := sqliteMgr.GetConnection(entityID)
	if err != nil {
		t.Fatalf("failed to get connection: %v", err)
	}
	defer db.Close()

	expectedTables := []string{
		"accounts",
		"entries",
		"postings",
		"work_logs",
		"decisions_log",
		"sync_metadata",
	}

	for _, table := range expectedTables {
		var count int
		query := `SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name=?`
		err := db.QueryRow(query, table).Scan(&count)
		if err != nil {
			t.Errorf("error checking table %s: %v", table, err)
			continue
		}
		if count == 0 {
			t.Errorf("table %s does not exist", table)
		}
	}
}

func TestSQLiteManager_WALModeEnabled(t *testing.T) {
	dataDir := "../../data/entities"
	defer os.RemoveAll(dataDir)

	sqliteMgr := lifecycle.NewSQLiteManager()
	defer sqliteMgr.CloseAll()

	entityID := "test_entity_wal"
	db, err := sqliteMgr.GetConnection(entityID)
	if err != nil {
		t.Fatalf("failed to get connection: %v", err)
	}
	defer db.Close()

	var journalMode string
	err = db.QueryRow("PRAGMA journal_mode").Scan(&journalMode)
	if err != nil {
		t.Fatalf("failed to query journal_mode: %v", err)
	}

	if journalMode != "wal" {
		t.Errorf("expected WAL mode, got '%s'", journalMode)
	}
}

func TestSQLiteManager_ForeignKeysEnabled(t *testing.T) {
	dataDir := "../../data/entities"
	defer os.RemoveAll(dataDir)

	sqliteMgr := lifecycle.NewSQLiteManager()
	defer sqliteMgr.CloseAll()

	entityID := "test_entity_fk"
	db, err := sqliteMgr.GetConnection(entityID)
	if err != nil {
		t.Fatalf("failed to get connection: %v", err)
	}
	defer db.Close()

	var foreignKeys int
	err = db.QueryRow("PRAGMA foreign_keys").Scan(&foreignKeys)
	if err != nil {
		t.Fatalf("failed to query foreign_keys: %v", err)
	}

	if foreignKeys != 1 {
		t.Errorf("expected foreign_keys=ON (1), got %d", foreignKeys)
	}
}

func TestSQLiteManager_MultipleConnections(t *testing.T) {
	dataDir := "../../data/entities"
	defer os.RemoveAll(dataDir)

	sqliteMgr := lifecycle.NewSQLiteManager()
	defer sqliteMgr.CloseAll()

	entityID1 := "entity_one"
	entityID2 := "entity_two"

	db1, err := sqliteMgr.GetConnection(entityID1)
	if err != nil {
		t.Fatalf("failed to get connection for entity_one: %v", err)
	}

	db2, err := sqliteMgr.GetConnection(entityID2)
	if err != nil {
		t.Fatalf("failed to get connection for entity_two: %v", err)
	}

	if db1 == db2 {
		t.Error("expected different database connections for different entities")
	}

	db1Again, err := sqliteMgr.GetConnection(entityID1)
	if err != nil {
		t.Fatalf("failed to get connection for entity_one again: %v", err)
	}

	if db1 != db1Again {
		t.Error("expected same connection when requesting same entity twice")
	}
}
