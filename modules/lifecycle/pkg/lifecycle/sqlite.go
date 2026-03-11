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

type SQLiteManager struct {
	connections map[string]*sql.DB
	mu          sync.RWMutex
	migrator    domain.Migrator
	dataDir     string
}

var _ LifecycleManager = (*SQLiteManager)(nil)

func NewSQLiteManager() *SQLiteManager {
	return NewSQLiteManagerWithDataDir("../../data/entities")
}

func NewSQLiteManagerWithDataDir(dataDir string) *SQLiteManager {
	return &SQLiteManager{
		connections: make(map[string]*sql.DB),
		migrator:    repository.NewMigrator(),
		dataDir:     dataDir,
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

	dbPath := filepath.Join(m.dataDir, fmt.Sprintf("%s.db", entityID))

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

// EntityExists verifica se o banco de dados da entidade já existe
func (m *SQLiteManager) EntityExists(entityID string) (bool, error) {
	dbPath := filepath.Join(m.dataDir, fmt.Sprintf("%s.db", entityID))

	// Verificar se o arquivo existe
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		return false, nil
	} else if err != nil {
		return false, fmt.Errorf("erro ao verificar arquivo do banco: %w", err)
	}

	return true, nil
}

// CreateEntity cria uma nova entidade com banco de dados inicializado
func (m *SQLiteManager) CreateEntity(entityID, entityName string) error {
	// Verificar se já existe
	exists, err := m.EntityExists(entityID)
	if err != nil {
		return fmt.Errorf("erro ao verificar existência da entidade: %w", err)
	}

	if exists {
		return nil // Já existe, não precisa criar
	}

	// Obter conexão (isso criará o banco se não existir)
	db, err := m.GetConnection(entityID)
	if err != nil {
		return fmt.Errorf("erro ao criar conexão para nova entidade: %w", err)
	}

	// Aqui poderíamos adicionar dados iniciais específicos da entidade
	// Por exemplo, criar contas padrão, configurar nome da entidade, etc.
	// Por enquanto, apenas ping no banco para garantir que está funcionando
	if err := db.Ping(); err != nil {
		return fmt.Errorf("erro ao conectar com banco da nova entidade: %w", err)
	}

	// Para agora, apenas garantir que as migrações foram executadas
	// (já feito pelo GetConnection através do migrator.RunMigrations)

	fmt.Printf("✅ Entidade criada: %s (%s)\n", entityName, entityID)
	return nil
}
