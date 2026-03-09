package repository

import (
	"context"

	"github.com/providentia/digna/budget/internal/domain"
)

// BudgetRepository interface para persistência de orçamentos
type BudgetRepository interface {
	// SavePlan salva ou atualiza um plano orçamentário
	SavePlan(ctx context.Context, entityID string, plan *domain.BudgetPlan) error

	// GetPlan obtém um plano orçamentário por ID
	GetPlan(ctx context.Context, entityID, planID string) (*domain.BudgetPlan, error)

	// ListPlansByPeriod lista todos os planos de um período
	ListPlansByPeriod(ctx context.Context, entityID, period string) ([]*domain.BudgetPlan, error)

	// ListPlansByCategory lista planos por categoria em um período
	ListPlansByCategory(ctx context.Context, entityID, period, category string) ([]*domain.BudgetPlan, error)

	// DeletePlan remove um plano orçamentário
	DeletePlan(ctx context.Context, entityID, planID string) error

	// GetCategories retorna as categorias disponíveis
	GetCategories(ctx context.Context) []domain.BudgetCategory

	// Transaction management
	BeginTx(ctx context.Context, entityID string) (interface{}, error)
	CommitTx(tx interface{}) error
	RollbackTx(tx interface{}) error
}
