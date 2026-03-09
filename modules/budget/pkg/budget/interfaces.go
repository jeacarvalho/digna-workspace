package budget

import (
	"context"
	"time"
)

// BudgetPlanRequest representa uma requisição para criar/atualizar plano orçamentário
type BudgetPlanRequest struct {
	EntityID    string `json:"entity_id"`
	Period      string `json:"period"`      // Formato: "YYYY-MM"
	Category    string `json:"category"`    // Categoria: "INSUMOS", "ENERGIA", etc.
	Planned     int64  `json:"planned"`     // Valor planejado em centavos
	Description string `json:"description"` // Descrição opcional
}

// BudgetPlanResponse representa a resposta de uma operação de plano
type BudgetPlanResponse struct {
	PlanID  string `json:"plan_id"`
	Success bool   `json:"success"`
	Error   string `json:"error,omitempty"`
}

// BudgetExecution representa a execução de um plano orçamentário
type BudgetExecution struct {
	PlanID      string `json:"plan_id"`
	Period      string `json:"period"`
	Category    string `json:"category"`
	Planned     int64  `json:"planned"`      // Em centavos
	Executed    int64  `json:"executed"`     // Em centavos
	Remaining   int64  `json:"remaining"`    // Em centavos
	Percentage  int    `json:"percentage"`   // 0-100
	AlertStatus string `json:"alert_status"` // "SAFE", "WARNING", "EXCEEDED"
	Description string `json:"description,omitempty"`
}

// BudgetSummary representa um resumo da execução orçamentária
type BudgetSummary struct {
	Period          string `json:"period"`
	TotalPlanned    int64  `json:"total_planned"`   // Em centavos
	TotalExecuted   int64  `json:"total_executed"`  // Em centavos
	TotalRemaining  int64  `json:"total_remaining"` // Em centavos
	Percentage      int    `json:"percentage"`      // 0-100
	OverallStatus   string `json:"overall_status"`  // "SAFE", "WARNING", "EXCEEDED"
	SafeCount       int    `json:"safe_count"`
	WarningCount    int    `json:"warning_count"`
	ExceededCount   int    `json:"exceeded_count"`
	TotalCategories int    `json:"total_categories"`
}

// Category representa uma categoria de orçamento
type Category struct {
	ID    string `json:"id"`    // Código: "INSUMOS", "ENERGIA", etc.
	Label string `json:"label"` // Nome em português: "Insumos", "Energia", etc.
}

// BudgetAPI interface pública do módulo budget
type BudgetAPI interface {
	// Plan Management
	CreatePlan(ctx context.Context, req BudgetPlanRequest) (*BudgetPlanResponse, error)
	GetPlan(ctx context.Context, entityID, planID string) (*BudgetExecution, error)
	DeletePlan(ctx context.Context, entityID, planID string) (*BudgetPlanResponse, error)

	// Reports
	GetExecutionReport(ctx context.Context, entityID, period string) ([]*BudgetExecution, error)
	GetExecutionSummary(ctx context.Context, entityID, period string) (*BudgetSummary, error)

	// Categories
	GetCategories(ctx context.Context) ([]*Category, error)

	// Periods
	GetAvailablePeriods(ctx context.Context, entityID string) ([]string, error)
}

// CashFlowPort interface para integração com cash_flow
type CashFlowPort interface {
	GetCashFlow(entityID string, startDate, endDate time.Time) (*CashFlowResponse, error)
}

// CashFlowResponse representa a resposta do cash_flow
type CashFlowResponse struct {
	EntityID    string      `json:"entity_id"`
	TotalCredit int64       `json:"total_credit"`
	TotalDebit  int64       `json:"total_debit"`
	Balance     int64       `json:"balance"`
	Entries     []CashEntry `json:"entries"`
	PeriodStart time.Time   `json:"period_start"`
	PeriodEnd   time.Time   `json:"period_end"`
}

// CashEntry representa uma entrada de caixa
type CashEntry struct {
	ID          int64     `json:"id"`
	EntityID    string    `json:"entity_id"`
	Type        string    `json:"type"`
	Amount      int64     `json:"amount"`
	Description string    `json:"description"`
	Category    string    `json:"category"`
	Date        time.Time `json:"date"`
	CreatedAt   time.Time `json:"created_at"`
}
