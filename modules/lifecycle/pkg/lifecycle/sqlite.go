package lifecycle

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/providentia/digna/lifecycle/internal/domain"
	"github.com/providentia/digna/lifecycle/internal/repository"
	"github.com/providentia/digna/lifecycle/internal/service"
)

type SQLiteManager struct {
	connections           map[string]*sql.DB
	centralDB             *sql.DB
	mu                    sync.RWMutex
	migrator              domain.Migrator
	centralMigrator       domain.Migrator
	dataDir               string
	accountantLinkService *service.AccountantLinkService
}

var _ LifecycleManager = (*SQLiteManager)(nil)
var _ AccountantLinkService = (*SQLiteManager)(nil)

func NewSQLiteManager() *SQLiteManager {
	return NewSQLiteManagerWithDataDir("../../data/entities")
}

func NewSQLiteManagerWithDataDir(dataDir string) *SQLiteManager {
	return &SQLiteManager{
		connections:     make(map[string]*sql.DB),
		migrator:        repository.NewMigrator(),
		centralMigrator: repository.NewCentralMigrator(),
		dataDir:         dataDir,
		// accountantLinkService will be lazily initialized
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

func (m *SQLiteManager) GetCentralConnection() (*sql.DB, error) {
	m.mu.RLock()
	if m.centralDB != nil {
		m.mu.RUnlock()
		return m.centralDB, nil
	}
	m.mu.RUnlock()

	m.mu.Lock()
	defer m.mu.Unlock()

	if m.centralDB != nil {
		return m.centralDB, nil
	}

	dbPath := filepath.Join(m.dataDir, "central.db")

	dir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create data directory for central db: %w", err)
	}

	dsn := fmt.Sprintf("%s?_journal_mode=WAL&_foreign_keys=ON&_synchronous=NORMAL", dbPath)
	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open central database: %w", err)
	}

	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to ping central database: %w", err)
	}

	if err := m.applyPragmas(db); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to apply pragmas to central database: %w", err)
	}

	if err := m.centralMigrator.RunMigrations(db); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to run central migrations: %w", err)
	}

	m.centralDB = db
	return db, nil
}

func (m *SQLiteManager) CloseAll() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	var errs []error

	// Fechar conexão central
	if m.centralDB != nil {
		if err := m.centralDB.Close(); err != nil {
			errs = append(errs, fmt.Errorf("failed to close central database: %w", err))
		}
		m.centralDB = nil
	}

	// Fechar conexões das entidades
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

// getAccountantLinkService returns the accountant link service, initializing it if necessary
func (m *SQLiteManager) getAccountantLinkService() (*service.AccountantLinkService, error) {
	m.mu.RLock()
	svc := m.accountantLinkService
	m.mu.RUnlock()

	if svc != nil {
		return svc, nil
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	// Double-check after acquiring lock
	if m.accountantLinkService != nil {
		return m.accountantLinkService, nil
	}

	// Get central database connection
	db, err := m.GetCentralConnection()
	if err != nil {
		return nil, fmt.Errorf("failed to get central database connection: %w", err)
	}

	// Create repository and service
	repo := repository.NewSQLiteEnterpriseAccountantRepository(db)
	svc = service.NewAccountantLinkService(repo)
	m.accountantLinkService = svc

	return svc, nil
}

// GetValidEnterprisesForAccountant returns list of enterprises an accountant can access during a period
func (m *SQLiteManager) GetValidEnterprisesForAccountant(ctx context.Context, accountantID string, startTime, endTime time.Time) ([]string, error) {
	svc, err := m.getAccountantLinkService()
	if err != nil {
		return nil, err
	}
	return svc.GetValidEnterprisesForAccountant(ctx, accountantID, startTime, endTime)
}

// ValidateAccountantAccess checks if an accountant has access to an enterprise during a period
func (m *SQLiteManager) ValidateAccountantAccess(ctx context.Context, accountantID, enterpriseID string, startTime, endTime time.Time) (bool, error) {
	svc, err := m.getAccountantLinkService()
	if err != nil {
		return false, err
	}
	return svc.ValidateAccountantAccess(ctx, accountantID, enterpriseID, startTime, endTime)
}

// CreateLink creates a new accountant-enterprise link
func (m *SQLiteManager) CreateLink(enterpriseID, accountantID, delegatedBy string) (*domain.EnterpriseAccountant, error) {
	svc, err := m.getAccountantLinkService()
	if err != nil {
		return nil, err
	}
	return svc.CreateLink(enterpriseID, accountantID, delegatedBy)
}

// DeactivateLink deactivates a link (Exit Power)
func (m *SQLiteManager) DeactivateLink(linkID, enterpriseID, requestedBy string) error {
	svc, err := m.getAccountantLinkService()
	if err != nil {
		return err
	}
	return svc.DeactivateLink(linkID, enterpriseID, requestedBy)
}

// ReactivateLink reactivates an inactive link
func (m *SQLiteManager) ReactivateLink(linkID string) error {
	svc, err := m.getAccountantLinkService()
	if err != nil {
		return err
	}
	return svc.ReactivateLink(linkID)
}

// GetValidDateRange returns the valid date range for an accountant-enterprise relationship
func (m *SQLiteManager) GetValidDateRange(enterpriseID, accountantID string) (time.Time, time.Time, error) {
	svc, err := m.getAccountantLinkService()
	if err != nil {
		return time.Time{}, time.Time{}, err
	}
	return svc.GetValidDateRange(enterpriseID, accountantID)
}

// GetActiveAccountant returns the active accountant for an enterprise
func (m *SQLiteManager) GetActiveAccountant(enterpriseID string) (*domain.EnterpriseAccountant, error) {
	svc, err := m.getAccountantLinkService()
	if err != nil {
		return nil, err
	}
	return svc.GetActiveAccountant(enterpriseID)
}
