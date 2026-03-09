package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/providentia/digna/budget/internal/domain"
	"github.com/providentia/digna/lifecycle/pkg/lifecycle"
)

// SQLiteBudgetRepository implementa BudgetRepository usando SQLite
type SQLiteBudgetRepository struct {
	lifecycleManager lifecycle.LifecycleManager
}

// NewSQLiteBudgetRepository cria um novo repositório SQLite para orçamentos
func NewSQLiteBudgetRepository(lm lifecycle.LifecycleManager) *SQLiteBudgetRepository {
	return &SQLiteBudgetRepository{
		lifecycleManager: lm,
	}
}

// SavePlan salva ou atualiza um plano orçamentário
func (r *SQLiteBudgetRepository) SavePlan(ctx context.Context, entityID string, plan *domain.BudgetPlan) error {
	db, err := r.lifecycleManager.GetConnection(entityID)
	if err != nil {
		return fmt.Errorf("erro ao obter conexão: %w", err)
	}

	// Verificar se a tabela existe, criar se necessário
	if err := r.ensureTables(ctx, entityID); err != nil {
		return fmt.Errorf("erro ao garantir tabelas: %w", err)
	}

	// Gerar ID se não existir
	if plan.ID == "" {
		plan.ID = fmt.Sprintf("budget_%d", time.Now().UnixNano())
	}

	now := time.Now()
	if plan.CreatedAt.IsZero() {
		plan.CreatedAt = now
	}
	plan.UpdatedAt = now

	// Validar plano
	if err := plan.Validate(); err != nil {
		return err
	}

	// Inserir ou atualizar
	query := `
		INSERT OR REPLACE INTO budget_plans (
			id, entity_id, period, category, planned, description, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err = db.ExecContext(ctx, query,
		plan.ID,
		entityID,
		plan.Period,
		plan.Category,
		plan.Planned,
		plan.Description,
		plan.CreatedAt,
		plan.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("erro ao salvar plano orçamentário: %w", err)
	}

	return nil
}

// GetPlan obtém um plano orçamentário por ID
func (r *SQLiteBudgetRepository) GetPlan(ctx context.Context, entityID, planID string) (*domain.BudgetPlan, error) {
	db, err := r.lifecycleManager.GetConnection(entityID)
	if err != nil {
		return nil, fmt.Errorf("erro ao obter conexão: %w", err)
	}

	// Verificar se a tabela existe, criar se necessário
	if err := r.ensureTables(ctx, entityID); err != nil {
		return nil, fmt.Errorf("erro ao garantir tabelas: %w", err)
	}

	query := `
		SELECT id, entity_id, period, category, planned, description, created_at, updated_at
		FROM budget_plans
		WHERE id = ? AND entity_id = ?
	`

	row := db.QueryRowContext(ctx, query, planID, entityID)

	var plan domain.BudgetPlan
	var createdAt, updatedAt string

	err = row.Scan(
		&plan.ID,
		&plan.EntityID,
		&plan.Period,
		&plan.Category,
		&plan.Planned,
		&plan.Description,
		&createdAt,
		&updatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, domain.ErrBudgetPlanNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar plano orçamentário: %w", err)
	}

	// Converter timestamps
	if plan.CreatedAt, err = time.Parse(time.RFC3339, createdAt); err != nil {
		return nil, fmt.Errorf("erro ao converter created_at: %w", err)
	}
	if plan.UpdatedAt, err = time.Parse(time.RFC3339, updatedAt); err != nil {
		return nil, fmt.Errorf("erro ao converter updated_at: %w", err)
	}

	return &plan, nil
}

// ListPlansByPeriod lista todos os planos de um período
func (r *SQLiteBudgetRepository) ListPlansByPeriod(ctx context.Context, entityID, period string) ([]*domain.BudgetPlan, error) {
	db, err := r.lifecycleManager.GetConnection(entityID)
	if err != nil {
		return nil, fmt.Errorf("erro ao obter conexão: %w", err)
	}

	// Verificar se a tabela existe, criar se necessário
	if err := r.ensureTables(ctx, entityID); err != nil {
		return nil, fmt.Errorf("erro ao garantir tabelas: %w", err)
	}

	query := `
		SELECT id, entity_id, period, category, planned, description, created_at, updated_at
		FROM budget_plans
		WHERE entity_id = ? AND period = ?
		ORDER BY category, created_at
	`

	rows, err := db.QueryContext(ctx, query, entityID, period)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar planos por período: %w", err)
	}
	defer rows.Close()

	var plans []*domain.BudgetPlan
	for rows.Next() {
		var plan domain.BudgetPlan
		var createdAt, updatedAt string

		err := rows.Scan(
			&plan.ID,
			&plan.EntityID,
			&plan.Period,
			&plan.Category,
			&plan.Planned,
			&plan.Description,
			&createdAt,
			&updatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("erro ao scanear plano: %w", err)
		}

		// Converter timestamps
		if plan.CreatedAt, err = time.Parse(time.RFC3339, createdAt); err != nil {
			return nil, fmt.Errorf("erro ao converter created_at: %w", err)
		}
		if plan.UpdatedAt, err = time.Parse(time.RFC3339, updatedAt); err != nil {
			return nil, fmt.Errorf("erro ao converter updated_at: %w", err)
		}

		plans = append(plans, &plan)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("erro ao iterar planos: %w", err)
	}

	return plans, nil
}

// ListPlansByCategory lista planos por categoria em um período
func (r *SQLiteBudgetRepository) ListPlansByCategory(ctx context.Context, entityID, period, category string) ([]*domain.BudgetPlan, error) {
	db, err := r.lifecycleManager.GetConnection(entityID)
	if err != nil {
		return nil, fmt.Errorf("erro ao obter conexão: %w", err)
	}

	// Verificar se a tabela existe, criar se necessário
	if err := r.ensureTables(ctx, entityID); err != nil {
		return nil, fmt.Errorf("erro ao garantir tabelas: %w", err)
	}

	query := `
		SELECT id, entity_id, period, category, planned, description, created_at, updated_at
		FROM budget_plans
		WHERE entity_id = ? AND period = ? AND category = ?
		ORDER BY created_at
	`

	rows, err := db.QueryContext(ctx, query, entityID, period, category)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar planos por categoria: %w", err)
	}
	defer rows.Close()

	var plans []*domain.BudgetPlan
	for rows.Next() {
		var plan domain.BudgetPlan
		var createdAt, updatedAt string

		err := rows.Scan(
			&plan.ID,
			&plan.EntityID,
			&plan.Period,
			&plan.Category,
			&plan.Planned,
			&plan.Description,
			&createdAt,
			&updatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("erro ao scanear plano: %w", err)
		}

		// Converter timestamps
		if plan.CreatedAt, err = time.Parse(time.RFC3339, createdAt); err != nil {
			return nil, fmt.Errorf("erro ao converter created_at: %w", err)
		}
		if plan.UpdatedAt, err = time.Parse(time.RFC3339, updatedAt); err != nil {
			return nil, fmt.Errorf("erro ao converter updated_at: %w", err)
		}

		plans = append(plans, &plan)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("erro ao iterar planos: %w", err)
	}

	return plans, nil
}

// DeletePlan remove um plano orçamentário
func (r *SQLiteBudgetRepository) DeletePlan(ctx context.Context, entityID, planID string) error {
	db, err := r.lifecycleManager.GetConnection(entityID)
	if err != nil {
		return fmt.Errorf("erro ao obter conexão: %w", err)
	}

	// Verificar se a tabela existe, criar se necessário
	if err := r.ensureTables(ctx, entityID); err != nil {
		return fmt.Errorf("erro ao garantir tabelas: %w", err)
	}

	query := "DELETE FROM budget_plans WHERE id = ? AND entity_id = ?"
	result, err := db.ExecContext(ctx, query, planID, entityID)
	if err != nil {
		return fmt.Errorf("erro ao deletar plano: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("erro ao verificar linhas afetadas: %w", err)
	}

	if rowsAffected == 0 {
		return domain.ErrBudgetPlanNotFound
	}

	return nil
}

// GetCategories retorna as categorias disponíveis
func (r *SQLiteBudgetRepository) GetCategories(ctx context.Context) []domain.BudgetCategory {
	return []domain.BudgetCategory{
		domain.CategoryRawMaterials,
		domain.CategoryEnergy,
		domain.CategoryEquipment,
		domain.CategoryTransport,
		domain.CategoryMaintenance,
		domain.CategoryServices,
		domain.CategoryOther,
	}
}

// ensureTables garante que as tabelas necessárias existam
func (r *SQLiteBudgetRepository) ensureTables(ctx context.Context, entityID string) error {
	db, err := r.lifecycleManager.GetConnection(entityID)
	if err != nil {
		return fmt.Errorf("erro ao obter conexão: %w", err)
	}

	// Verificar se a conexão está válida
	if err := db.PingContext(ctx); err != nil {
		return fmt.Errorf("conexão com banco inválida: %w", err)
	}

	// Criar tabela de planos orçamentários
	createTableSQL := `
		CREATE TABLE IF NOT EXISTS budget_plans (
			id TEXT PRIMARY KEY,
			entity_id TEXT NOT NULL,
			period TEXT NOT NULL, -- Formato: YYYY-MM
			category TEXT NOT NULL,
			planned INTEGER NOT NULL, -- Em centavos (int64)
			description TEXT,
			created_at TIMESTAMP NOT NULL,
			updated_at TIMESTAMP NOT NULL,
			UNIQUE(entity_id, period, category)
		);

		CREATE INDEX IF NOT EXISTS idx_budget_plans_entity_period 
		ON budget_plans(entity_id, period);

		CREATE INDEX IF NOT EXISTS idx_budget_plans_entity_category 
		ON budget_plans(entity_id, category);
	`

	_, err = db.ExecContext(ctx, createTableSQL)
	if err != nil {
		return fmt.Errorf("erro ao criar tabelas: %w", err)
	}

	return nil
}

// BeginTx inicia uma transação
func (r *SQLiteBudgetRepository) BeginTx(ctx context.Context, entityID string) (interface{}, error) {
	db, err := r.lifecycleManager.GetConnection(entityID)
	if err != nil {
		return nil, fmt.Errorf("erro ao obter conexão: %w", err)
	}

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("erro ao iniciar transação: %w", err)
	}

	return tx, nil
}

// CommitTx commita uma transação
func (r *SQLiteBudgetRepository) CommitTx(tx interface{}) error {
	sqlTx, ok := tx.(*sql.Tx)
	if !ok {
		return fmt.Errorf("tipo de transação inválido")
	}

	return sqlTx.Commit()
}

// RollbackTx faz rollback de uma transação
func (r *SQLiteBudgetRepository) RollbackTx(tx interface{}) error {
	sqlTx, ok := tx.(*sql.Tx)
	if !ok {
		return fmt.Errorf("tipo de transação inválido")
	}

	return sqlTx.Rollback()
}
