package repository

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/providentia/digna/core_lume/internal/domain"
	"github.com/providentia/digna/lifecycle/pkg/lifecycle"
)

type SQLiteDASMEIRepository struct {
	lifecycleManager lifecycle.LifecycleManager
}

func NewSQLiteDASMEIRepository(lm lifecycle.LifecycleManager) *SQLiteDASMEIRepository {
	return &SQLiteDASMEIRepository{lifecycleManager: lm}
}

func (r *SQLiteDASMEIRepository) GetDB(entityID string) (*sql.DB, error) {
	return r.lifecycleManager.GetConnection(entityID)
}

func (r *SQLiteDASMEIRepository) Save(das *domain.DASMEI) error {
	if err := das.Validate(); err != nil {
		return fmt.Errorf("invalid DAS MEI: %w", err)
	}

	db, err := r.GetDB(das.EntityID)
	if err != nil {
		return fmt.Errorf("failed to get connection: %w", err)
	}

	_, err = db.Exec(
		`INSERT INTO das_mei (id, entity_id, competencia, valor_devido, valor_pago, data_vencimento, 
		data_pagamento, status, salario_minimo, activity_type, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(id) DO UPDATE SET
			valor_devido = excluded.valor_devido,
			valor_pago = excluded.valor_pago,
			data_vencimento = excluded.data_vencimento,
			data_pagamento = excluded.data_pagamento,
			status = excluded.status,
			salario_minimo = excluded.salario_minimo,
			activity_type = excluded.activity_type,
			updated_at = excluded.updated_at`,
		das.ID, das.EntityID, das.Competencia, das.ValorDevido, das.ValorPago,
		das.DataVencimento, das.DataPagamento, string(das.Status), das.SalarioMinimo,
		string(das.ActivityType), das.CreatedAt, das.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to save DAS MEI: %w", err)
	}

	return nil
}

func (r *SQLiteDASMEIRepository) FindByID(entityID, dasID string) (*domain.DASMEI, error) {
	db, err := r.GetDB(entityID)
	if err != nil {
		return nil, fmt.Errorf("failed to get connection: %w", err)
	}

	var das domain.DASMEI
	var statusStr, activityStr string

	err = db.QueryRow(
		`SELECT id, entity_id, competencia, valor_devido, valor_pago, data_vencimento,
		data_pagamento, status, salario_minimo, activity_type, created_at, updated_at
		FROM das_mei WHERE id = ? AND entity_id = ?`,
		dasID, entityID,
	).Scan(&das.ID, &das.EntityID, &das.Competencia, &das.ValorDevido, &das.ValorPago,
		&das.DataVencimento, &das.DataPagamento, &statusStr, &das.SalarioMinimo,
		&activityStr, &das.CreatedAt, &das.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("DAS MEI not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to query DAS MEI: %w", err)
	}

	das.Status = domain.DASMEIStatus(statusStr)
	das.ActivityType = domain.ActivityType(activityStr)

	return &das, nil
}

func (r *SQLiteDASMEIRepository) FindByCompetencia(entityID, competencia string) (*domain.DASMEI, error) {
	db, err := r.GetDB(entityID)
	if err != nil {
		return nil, fmt.Errorf("failed to get connection: %w", err)
	}

	var das domain.DASMEI
	var statusStr, activityStr string

	err = db.QueryRow(
		`SELECT id, entity_id, competencia, valor_devido, valor_pago, data_vencimento,
		data_pagamento, status, salario_minimo, activity_type, created_at, updated_at
		FROM das_mei WHERE entity_id = ? AND competencia = ?`,
		entityID, competencia,
	).Scan(&das.ID, &das.EntityID, &das.Competencia, &das.ValorDevido, &das.ValorPago,
		&das.DataVencimento, &das.DataPagamento, &statusStr, &das.SalarioMinimo,
		&activityStr, &das.CreatedAt, &das.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("DAS MEI not found for competencia %s", competencia)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to query DAS MEI: %w", err)
	}

	das.Status = domain.DASMEIStatus(statusStr)
	das.ActivityType = domain.ActivityType(activityStr)

	return &das, nil
}

func (r *SQLiteDASMEIRepository) ListByEntity(entityID string) ([]*domain.DASMEI, error) {
	db, err := r.GetDB(entityID)
	if err != nil {
		return nil, fmt.Errorf("failed to get connection: %w", err)
	}

	rows, err := db.Query(
		`SELECT id, entity_id, competencia, valor_devido, valor_pago, data_vencimento,
		data_pagamento, status, salario_minimo, activity_type, created_at, updated_at
		FROM das_mei WHERE entity_id = ? ORDER BY competencia DESC`,
		entityID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to query DAS MEI list: %w", err)
	}
	defer rows.Close()

	var dasList []*domain.DASMEI
	for rows.Next() {
		var das domain.DASMEI
		var statusStr, activityStr string

		if err := rows.Scan(&das.ID, &das.EntityID, &das.Competencia, &das.ValorDevido, &das.ValorPago,
			&das.DataVencimento, &das.DataPagamento, &statusStr, &das.SalarioMinimo,
			&activityStr, &das.CreatedAt, &das.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan DAS MEI: %w", err)
		}

		das.Status = domain.DASMEIStatus(statusStr)
		das.ActivityType = domain.ActivityType(activityStr)
		dasList = append(dasList, &das)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return dasList, nil
}

func (r *SQLiteDASMEIRepository) ListPending(entityID string) ([]*domain.DASMEI, error) {
	db, err := r.GetDB(entityID)
	if err != nil {
		return nil, fmt.Errorf("failed to get connection: %w", err)
	}

	rows, err := db.Query(
		`SELECT id, entity_id, competencia, valor_devido, valor_pago, data_vencimento,
		data_pagamento, status, salario_minimo, activity_type, created_at, updated_at
		FROM das_mei WHERE entity_id = ? AND status = 'PENDENTE' ORDER BY data_vencimento ASC`,
		entityID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to query pending DAS MEI: %w", err)
	}
	defer rows.Close()

	var dasList []*domain.DASMEI
	for rows.Next() {
		var das domain.DASMEI
		var statusStr, activityStr string

		if err := rows.Scan(&das.ID, &das.EntityID, &das.Competencia, &das.ValorDevido, &das.ValorPago,
			&das.DataVencimento, &das.DataPagamento, &statusStr, &das.SalarioMinimo,
			&activityStr, &das.CreatedAt, &das.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan DAS MEI: %w", err)
		}

		das.Status = domain.DASMEIStatus(statusStr)
		das.ActivityType = domain.ActivityType(activityStr)
		dasList = append(dasList, &das)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return dasList, nil
}

func (r *SQLiteDASMEIRepository) ListOverdue(entityID string) ([]*domain.DASMEI, error) {
	db, err := r.GetDB(entityID)
	if err != nil {
		return nil, fmt.Errorf("failed to get connection: %w", err)
	}

	now := time.Now().Unix()

	rows, err := db.Query(
		`SELECT id, entity_id, competencia, valor_devido, valor_pago, data_vencimento,
		data_pagamento, status, salario_minimo, activity_type, created_at, updated_at
		FROM das_mei WHERE entity_id = ? AND status = 'PENDENTE' AND data_vencimento < ? ORDER BY data_vencimento ASC`,
		entityID, now,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to query overdue DAS MEI: %w", err)
	}
	defer rows.Close()

	var dasList []*domain.DASMEI
	for rows.Next() {
		var das domain.DASMEI
		var statusStr, activityStr string

		if err := rows.Scan(&das.ID, &das.EntityID, &das.Competencia, &das.ValorDevido, &das.ValorPago,
			&das.DataVencimento, &das.DataPagamento, &statusStr, &das.SalarioMinimo,
			&activityStr, &das.CreatedAt, &das.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan DAS MEI: %w", err)
		}

		das.Status = domain.DASMEIStatusPending
		das.ActivityType = domain.ActivityType(activityStr)
		dasList = append(dasList, &das)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return dasList, nil
}

func (r *SQLiteDASMEIRepository) Update(das *domain.DASMEI) error {
	if err := das.Validate(); err != nil {
		return fmt.Errorf("invalid DAS MEI: %w", err)
	}

	db, err := r.GetDB(das.EntityID)
	if err != nil {
		return fmt.Errorf("failed to get connection: %w", err)
	}

	das.UpdatedAt = time.Now().Unix()

	_, err = db.Exec(
		`UPDATE das_mei SET valor_devido = ?, valor_pago = ?, data_vencimento = ?,
		data_pagamento = ?, status = ?, salario_minimo = ?, activity_type = ?, updated_at = ?
		WHERE id = ? AND entity_id = ?`,
		das.ValorDevido, das.ValorPago, das.DataVencimento, das.DataPagamento,
		string(das.Status), das.SalarioMinimo, string(das.ActivityType), das.UpdatedAt,
		das.ID, das.EntityID,
	)
	if err != nil {
		return fmt.Errorf("failed to update DAS MEI: %w", err)
	}

	return nil
}

func (r *SQLiteDASMEIRepository) MarkAsPaid(entityID, dasID string, valorPago int64) error {
	db, err := r.GetDB(entityID)
	if err != nil {
		return fmt.Errorf("failed to get connection: %w", err)
	}

	now := time.Now().Unix()

	_, err = db.Exec(
		`UPDATE das_mei SET status = 'PAGO', valor_pago = ?, data_pagamento = ?, updated_at = ?
		WHERE id = ? AND entity_id = ?`,
		valorPago, now, now, dasID, entityID,
	)
	if err != nil {
		return fmt.Errorf("failed to mark DAS MEI as paid: %w", err)
	}

	return nil
}

// InitTable initializes the das_mei table for a specific entity
func (r *SQLiteDASMEIRepository) InitTable(entityID string) error {
	db, err := r.GetDB(entityID)
	if err != nil {
		return fmt.Errorf("failed to get connection: %w", err)
	}
	return InitDASMEITable(db)
}

// InitDASMEITable cria a tabela das_mei se não existir
func InitDASMEITable(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS das_mei (
			id TEXT PRIMARY KEY,
			entity_id TEXT NOT NULL,
			competencia TEXT NOT NULL,
			valor_devido INTEGER NOT NULL,
			valor_pago INTEGER DEFAULT 0,
			data_vencimento INTEGER NOT NULL,
			data_pagamento INTEGER DEFAULT 0,
			status TEXT NOT NULL,
			salario_minimo INTEGER NOT NULL,
			activity_type TEXT NOT NULL,
			created_at INTEGER,
			updated_at INTEGER,
			UNIQUE(entity_id, competencia)
		)
	`)
	if err != nil {
		return fmt.Errorf("failed to create das_mei table: %w", err)
	}

	// Cria índice para melhorar performance
	_, err = db.Exec(`
		CREATE INDEX IF NOT EXISTS idx_das_mei_entity_competencia 
		ON das_mei(entity_id, competencia)
	`)
	if err != nil {
		return fmt.Errorf("failed to create index: %w", err)
	}

	_, err = db.Exec(`
		CREATE INDEX IF NOT EXISTS idx_das_mei_entity_status 
		ON das_mei(entity_id, status)
	`)
	if err != nil {
		return fmt.Errorf("failed to create status index: %w", err)
	}

	return nil
}
