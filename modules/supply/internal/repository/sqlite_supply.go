package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/providentia/digna/lifecycle/pkg/lifecycle"
	"github.com/providentia/digna/supply/internal/domain"
)

type SQLiteSupplyRepository struct {
	lifecycleManager lifecycle.LifecycleManager
	dbCache          map[string]*sql.DB // Cache de conexões por entityID
}

func NewSQLiteSupplyRepository(lm lifecycle.LifecycleManager) *SQLiteSupplyRepository {
	return &SQLiteSupplyRepository{
		lifecycleManager: lm,
		dbCache:          make(map[string]*sql.DB),
	}
}

func (r *SQLiteSupplyRepository) initDB(ctx context.Context, entityID string) (*sql.DB, error) {
	// Verificar cache primeiro
	if db, ok := r.dbCache[entityID]; ok {
		// Verificar se a conexão ainda está válida
		if err := db.PingContext(ctx); err == nil {
			return db, nil
		}
		// Conexão inválida, remover do cache
		delete(r.dbCache, entityID)
	}

	// Obter nova conexão
	db, err := r.lifecycleManager.GetConnection(entityID)
	if err != nil {
		return nil, fmt.Errorf("failed to get database connection: %w", err)
	}

	// Verificar se a conexão está válida
	if err := db.PingContext(ctx); err != nil {
		// Não fechar aqui - deixar o lifecycle manager gerenciar
		return nil, fmt.Errorf("database connection is closed: %w", err)
	}

	// Executar migrações apenas uma vez por conexão
	if err := r.runMigrations(db); err != nil {
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	// Armazenar no cache
	r.dbCache[entityID] = db

	return db, nil
}

func (r *SQLiteSupplyRepository) runMigrations(db *sql.DB) error {
	// Criar tabela de controle de migrações
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS supply_migrations (
		name TEXT PRIMARY KEY,
		applied_at INTEGER NOT NULL
	)`)
	if err != nil {
		return fmt.Errorf("failed to create migrations table: %w", err)
	}

	// Verificar quais migrações já foram aplicadas
	rows, err := db.Query("SELECT name FROM supply_migrations")
	if err != nil {
		return fmt.Errorf("failed to query applied migrations: %w", err)
	}
	defer rows.Close()

	applied := make(map[string]bool)
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return fmt.Errorf("failed to scan migration name: %w", err)
		}
		applied[name] = true
	}

	migrations := []struct {
		name string
		sql  string
	}{
		{
			name: "create_suppliers",
			sql: `CREATE TABLE IF NOT EXISTS suppliers (
				id TEXT PRIMARY KEY,
				name TEXT NOT NULL,
				contact_info TEXT,
				created_at INTEGER NOT NULL
			)`,
		},
		{
			name: "create_stock_items",
			sql: `CREATE TABLE IF NOT EXISTS stock_items (
				id TEXT PRIMARY KEY,
				name TEXT NOT NULL,
				description TEXT,
				type TEXT NOT NULL CHECK(type IN ('INSUMO', 'PRODUTO', 'MERCADORIA')),
				unit TEXT NOT NULL DEFAULT 'UNIDADE' CHECK(unit IN ('UNIDADE', 'KG', 'G', 'L', 'M', 'CM', 'PACOTE', 'CAIXA', 'SACO')),
				quantity INTEGER NOT NULL DEFAULT 0,
				min_quantity INTEGER NOT NULL DEFAULT 0,
				unit_cost INTEGER NOT NULL,
				created_at INTEGER NOT NULL
			)`,
		},
		{
			name: "create_purchases",
			sql: `CREATE TABLE IF NOT EXISTS purchases (
				id TEXT PRIMARY KEY,
				supplier_id TEXT NOT NULL,
				total_value INTEGER NOT NULL,
				date INTEGER NOT NULL,
				created_at INTEGER NOT NULL,
				FOREIGN KEY (supplier_id) REFERENCES suppliers(id)
			)`,
		},
		{
			name: "create_purchase_items",
			sql: `CREATE TABLE IF NOT EXISTS purchase_items (
				id TEXT PRIMARY KEY,
				purchase_id TEXT NOT NULL,
				stock_item_id TEXT NOT NULL,
				quantity INTEGER NOT NULL,
				unit_cost INTEGER NOT NULL,
				total_cost INTEGER NOT NULL,
				FOREIGN KEY (purchase_id) REFERENCES purchases(id) ON DELETE CASCADE,
				FOREIGN KEY (stock_item_id) REFERENCES stock_items(id)
			)`,
		},
	}

	for _, migration := range migrations {
		// Aplicar apenas migrações não aplicadas
		if !applied[migration.name] {
			if _, err := db.Exec(migration.sql); err != nil {
				return fmt.Errorf("failed to run migration %s: %w", migration.name, err)
			}

			// Registrar migração aplicada
			_, err = db.Exec("INSERT OR REPLACE INTO supply_migrations (name, applied_at) VALUES (?, ?)",
				migration.name, time.Now().Unix())
			if err != nil {
				return fmt.Errorf("failed to record migration %s: %w", migration.name, err)
			}
		}
	}

	return nil
}

// Suppliers
func (r *SQLiteSupplyRepository) SaveSupplier(ctx context.Context, entityID string, supplier *domain.Supplier) error {
	db, err := r.initDB(ctx, entityID)
	if err != nil {
		return err
	}

	if supplier.ID == "" {
		supplier.ID = fmt.Sprintf("supp_%d", time.Now().UnixNano())
	}
	if supplier.CreatedAt.IsZero() {
		supplier.CreatedAt = time.Now()
	}

	query := `INSERT OR REPLACE INTO suppliers (id, name, contact_info, created_at) VALUES (?, ?, ?, ?)`
	_, err = db.ExecContext(ctx, query,
		supplier.ID,
		supplier.Name,
		supplier.ContactInfo,
		supplier.CreatedAt.Unix(),
	)
	return err
}

func (r *SQLiteSupplyRepository) GetSupplier(ctx context.Context, entityID, supplierID string) (*domain.Supplier, error) {
	db, err := r.initDB(ctx, entityID)
	if err != nil {
		return nil, err
	}

	query := `SELECT id, name, contact_info, created_at FROM suppliers WHERE id = ?`
	row := db.QueryRowContext(ctx, query, supplierID)

	var supplier domain.Supplier
	var createdAt int64
	err = row.Scan(&supplier.ID, &supplier.Name, &supplier.ContactInfo, &createdAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	supplier.CreatedAt = time.Unix(createdAt, 0)
	return &supplier, nil
}

func (r *SQLiteSupplyRepository) ListSuppliers(ctx context.Context, entityID string) ([]*domain.Supplier, error) {
	db, err := r.initDB(ctx, entityID)
	if err != nil {
		return nil, err
	}

	query := `SELECT id, name, contact_info, created_at FROM suppliers ORDER BY name`
	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var suppliers []*domain.Supplier
	for rows.Next() {
		var supplier domain.Supplier
		var createdAt int64
		if err := rows.Scan(&supplier.ID, &supplier.Name, &supplier.ContactInfo, &createdAt); err != nil {
			return nil, err
		}
		supplier.CreatedAt = time.Unix(createdAt, 0)
		suppliers = append(suppliers, &supplier)
	}

	return suppliers, nil
}

// Stock Items
func (r *SQLiteSupplyRepository) SaveStockItem(ctx context.Context, entityID string, item *domain.StockItem) error {
	db, err := r.initDB(ctx, entityID)
	if err != nil {
		return err
	}

	if item.ID == "" {
		item.ID = fmt.Sprintf("item_%d", time.Now().UnixNano())
	}
	if item.CreatedAt.IsZero() {
		item.CreatedAt = time.Now()
	}

	query := `INSERT OR REPLACE INTO stock_items (id, name, description, type, unit, quantity, min_quantity, unit_cost, created_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`
	_, err = db.ExecContext(ctx, query,
		item.ID,
		item.Name,
		item.Description,
		string(item.Type),
		string(item.Unit),
		item.Quantity,
		item.MinQuantity,
		item.UnitCost,
		item.CreatedAt.Unix(),
	)
	return err
}

func (r *SQLiteSupplyRepository) GetStockItem(ctx context.Context, entityID, itemID string) (*domain.StockItem, error) {
	db, err := r.initDB(ctx, entityID)
	if err != nil {
		return nil, err
	}

	query := `SELECT id, name, description, type, unit, quantity, min_quantity, unit_cost, created_at FROM stock_items WHERE id = ?`
	row := db.QueryRowContext(ctx, query, itemID)

	var item domain.StockItem
	var itemType string
	var itemUnit string
	var createdAt int64
	err = row.Scan(&item.ID, &item.Name, &item.Description, &itemType, &itemUnit, &item.Quantity, &item.MinQuantity, &item.UnitCost, &createdAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	item.Type = domain.StockItemType(itemType)
	item.Unit = domain.StockItemUnit(itemUnit)
	item.CreatedAt = time.Unix(createdAt, 0)
	return &item, nil
}

func (r *SQLiteSupplyRepository) ListStockItems(ctx context.Context, entityID string) ([]*domain.StockItem, error) {
	db, err := r.initDB(ctx, entityID)
	if err != nil {
		return nil, err
	}

	query := `SELECT id, name, description, type, unit, quantity, min_quantity, unit_cost, created_at FROM stock_items ORDER BY name`
	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []*domain.StockItem
	for rows.Next() {
		var item domain.StockItem
		var itemType string
		var itemUnit string
		var createdAt int64
		if err := rows.Scan(&item.ID, &item.Name, &item.Description, &itemType, &itemUnit, &item.Quantity, &item.MinQuantity, &item.UnitCost, &createdAt); err != nil {
			return nil, err
		}
		item.Type = domain.StockItemType(itemType)
		item.Unit = domain.StockItemUnit(itemUnit)
		item.CreatedAt = time.Unix(createdAt, 0)
		items = append(items, &item)
	}

	return items, nil
}

func (r *SQLiteSupplyRepository) ListStockItemsByType(ctx context.Context, entityID string, itemType domain.StockItemType) ([]*domain.StockItem, error) {
	db, err := r.initDB(ctx, entityID)
	if err != nil {
		return nil, err
	}

	query := `SELECT id, name, description, type, unit, quantity, min_quantity, unit_cost, created_at FROM stock_items WHERE type = ? ORDER BY name`
	rows, err := db.QueryContext(ctx, query, string(itemType))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []*domain.StockItem
	for rows.Next() {
		var item domain.StockItem
		var itemTypeStr string
		var itemUnit string
		var createdAt int64
		if err := rows.Scan(&item.ID, &item.Name, &item.Description, &itemTypeStr, &itemUnit, &item.Quantity, &item.MinQuantity, &item.UnitCost, &createdAt); err != nil {
			return nil, err
		}
		item.Type = domain.StockItemType(itemTypeStr)
		item.Unit = domain.StockItemUnit(itemUnit)
		item.CreatedAt = time.Unix(createdAt, 0)
		items = append(items, &item)
	}

	return items, nil
}

func (r *SQLiteSupplyRepository) UpdateStockQuantity(ctx context.Context, entityID, itemID string, delta int) error {
	db, err := r.initDB(ctx, entityID)
	if err != nil {
		return err
	}

	query := `UPDATE stock_items SET quantity = quantity + ? WHERE id = ?`
	result, err := db.ExecContext(ctx, query, delta, itemID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("stock item not found: %s", itemID)
	}

	return nil
}

// Purchases
func (r *SQLiteSupplyRepository) SavePurchase(ctx context.Context, entityID string, purchase *domain.Purchase) error {
	db, err := r.initDB(ctx, entityID)
	if err != nil {
		return err
	}

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if purchase.ID == "" {
		purchase.ID = fmt.Sprintf("pur_%d", time.Now().UnixNano())
	}
	if purchase.CreatedAt.IsZero() {
		purchase.CreatedAt = time.Now()
	}
	if purchase.Date.IsZero() {
		purchase.Date = time.Now()
	}

	// Salvar compra
	query := `INSERT INTO purchases (id, supplier_id, total_value, date, created_at) VALUES (?, ?, ?, ?, ?)`
	_, err = tx.ExecContext(ctx, query,
		purchase.ID,
		purchase.SupplierID,
		purchase.TotalValue,
		purchase.Date.Unix(),
		purchase.CreatedAt.Unix(),
	)
	if err != nil {
		return err
	}

	// Salvar itens da compra
	for _, item := range purchase.Items {
		if item.ID == "" {
			item.ID = fmt.Sprintf("pi_%d", time.Now().UnixNano())
		}
		item.PurchaseID = purchase.ID

		query := `INSERT INTO purchase_items (id, purchase_id, stock_item_id, quantity, unit_cost, total_cost) VALUES (?, ?, ?, ?, ?, ?)`
		_, err = tx.ExecContext(ctx, query,
			item.ID,
			item.PurchaseID,
			item.StockItemID,
			item.Quantity,
			item.UnitCost,
			item.TotalCost,
		)
		if err != nil {
			return err
		}

		// Atualizar estoque
		updateQuery := `UPDATE stock_items SET quantity = quantity + ? WHERE id = ?`
		_, err = tx.ExecContext(ctx, updateQuery, item.Quantity, item.StockItemID)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (r *SQLiteSupplyRepository) GetPurchase(ctx context.Context, entityID, purchaseID string) (*domain.Purchase, error) {
	db, err := r.initDB(ctx, entityID)
	if err != nil {
		return nil, err
	}

	// Buscar compra
	query := `SELECT id, supplier_id, total_value, date, created_at FROM purchases WHERE id = ?`
	row := db.QueryRowContext(ctx, query, purchaseID)

	var purchase domain.Purchase
	var dateUnix, createdAtUnix int64
	err = row.Scan(&purchase.ID, &purchase.SupplierID, &purchase.TotalValue, &dateUnix, &createdAtUnix)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	purchase.Date = time.Unix(dateUnix, 0)
	purchase.CreatedAt = time.Unix(createdAtUnix, 0)

	// Buscar itens da compra
	itemsQuery := `SELECT id, stock_item_id, quantity, unit_cost, total_cost FROM purchase_items WHERE purchase_id = ?`
	rows, err := db.QueryContext(ctx, itemsQuery, purchaseID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var item domain.PurchaseItem
		if err := rows.Scan(&item.ID, &item.StockItemID, &item.Quantity, &item.UnitCost, &item.TotalCost); err != nil {
			return nil, err
		}
		item.PurchaseID = purchaseID
		purchase.Items = append(purchase.Items, item)
	}

	return &purchase, nil
}

func (r *SQLiteSupplyRepository) ListPurchases(ctx context.Context, entityID string) ([]*domain.Purchase, error) {
	db, err := r.initDB(ctx, entityID)
	if err != nil {
		return nil, err
	}

	query := `SELECT id, supplier_id, total_value, date, created_at FROM purchases ORDER BY date DESC`
	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var purchases []*domain.Purchase
	for rows.Next() {
		var purchase domain.Purchase
		var dateUnix, createdAtUnix int64
		if err := rows.Scan(&purchase.ID, &purchase.SupplierID, &purchase.TotalValue, &dateUnix, &createdAtUnix); err != nil {
			return nil, err
		}
		purchase.Date = time.Unix(dateUnix, 0)
		purchase.CreatedAt = time.Unix(createdAtUnix, 0)
		purchases = append(purchases, &purchase)
	}

	return purchases, nil
}

func (r *SQLiteSupplyRepository) ListPurchasesBySupplier(ctx context.Context, entityID, supplierID string) ([]*domain.Purchase, error) {
	db, err := r.initDB(ctx, entityID)
	if err != nil {
		return nil, err
	}

	query := `SELECT id, supplier_id, total_value, date, created_at FROM purchases WHERE supplier_id = ? ORDER BY date DESC`
	rows, err := db.QueryContext(ctx, query, supplierID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var purchases []*domain.Purchase
	for rows.Next() {
		var purchase domain.Purchase
		var dateUnix, createdAtUnix int64
		if err := rows.Scan(&purchase.ID, &purchase.SupplierID, &purchase.TotalValue, &dateUnix, &createdAtUnix); err != nil {
			return nil, err
		}
		purchase.Date = time.Unix(dateUnix, 0)
		purchase.CreatedAt = time.Unix(createdAtUnix, 0)
		purchases = append(purchases, &purchase)
	}

	return purchases, nil
}

// Transaction management
func (r *SQLiteSupplyRepository) BeginTx(ctx context.Context, entityID string) (interface{}, error) {
	db, err := r.initDB(ctx, entityID)
	if err != nil {
		return nil, err
	}
	return db.BeginTx(ctx, nil)
}

func (r *SQLiteSupplyRepository) CommitTx(tx interface{}) error {
	sqlTx, ok := tx.(*sql.Tx)
	if !ok {
		return fmt.Errorf("invalid transaction type")
	}
	return sqlTx.Commit()
}

func (r *SQLiteSupplyRepository) RollbackTx(tx interface{}) error {
	sqlTx, ok := tx.(*sql.Tx)
	if !ok {
		return fmt.Errorf("invalid transaction type")
	}
	return sqlTx.Rollback()
}

// Helper function para serializar itens (usado em testes)
func serializeItems(items []domain.PurchaseItem) (string, error) {
	data, err := json.Marshal(items)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func deserializeItems(data string) ([]domain.PurchaseItem, error) {
	var items []domain.PurchaseItem
	if err := json.Unmarshal([]byte(data), &items); err != nil {
		return nil, err
	}
	return items, nil
}
