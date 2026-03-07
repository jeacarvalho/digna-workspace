package lifecycle

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	_ "github.com/mattn/go-sqlite3"
	"github.com/providentia/digna/lifecycle/internal/domain"
	"github.com/providentia/digna/lifecycle/internal/repository"
)

const (
	DataDir = "../../data/entities"
)

type SQLiteManager struct {
	connections map[string]*sql.DB
	mu          sync.RWMutex
	migrator    domain.Migrator
}

var _ LifecycleManager = (*SQLiteManager)(nil)

func NewSQLiteManager() *SQLiteManager {
	return &SQLiteManager{
		connections: make(map[string]*sql.DB),
		migrator:    repository.NewMigrator(),
	}
}

func (m *SQLiteManager) GetConnection(entityID string) (*sql.DB, error) {
	m.mu.RLock()
	if db, exists := m.connections[entityID]; exists {
		m.mu.RUnlock()
		return db, nil
	}
	m.mu.RUnlock()

	m.mu.Lock()
	defer m.mu.Unlock()

	if db, exists := m.connections[entityID]; exists {
		return db, nil
	}

	dbPath := filepath.Join(DataDir, fmt.Sprintf("%s.db", entityID))

	dir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create data directory: %w", err)
	}

	dsn := fmt.Sprintf("%s?_journal_mode=WAL&_foreign_keys=ON&_synchronous=NORMAL", dbPath)
	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	if err := m.applyPragmas(db); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to apply pragmas: %w", err)
	}

	if err := m.migrator.RunMigrations(db); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	m.connections[entityID] = db
	return db, nil
}

func (m *SQLiteManager) applyPragmas(db *sql.DB) error {
	pragmas := []string{
		"PRAGMA journal_mode=WAL",
		"PRAGMA foreign_keys=ON",
		"PRAGMA synchronous=NORMAL",
		"PRAGMA temp_store=MEMORY",
		"PRAGMA mmap_size=268435456",
	}

	for _, pragma := range pragmas {
		if _, err := db.Exec(pragma); err != nil {
			return fmt.Errorf("failed to execute %s: %w", pragma, err)
		}
	}

	return nil
}

func (m *SQLiteManager) CloseConnection(entityID string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if db, exists := m.connections[entityID]; exists {
		delete(m.connections, entityID)
		return db.Close()
	}

	return nil
}

func (m *SQLiteManager) CloseAll() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	var errs []error
	for entityID, db := range m.connections {
		if err := db.Close(); err != nil {
			errs = append(errs, fmt.Errorf("failed to close connection for %s: %w", entityID, err))
		}
		delete(m.connections, entityID)
	}

	if len(errs) > 0 {
		return fmt.Errorf("errors closing connections: %v", errs)
	}

	return nil
}
