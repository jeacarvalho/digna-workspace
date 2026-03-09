package budget

import (
	"context"
	"time"

	"github.com/providentia/digna/budget/internal/domain"
	"github.com/providentia/digna/budget/internal/repository"
	"github.com/providentia/digna/budget/internal/service"
	"github.com/providentia/digna/lifecycle/pkg/lifecycle"
)

// BudgetAPIImpl implementa BudgetAPI
type BudgetAPIImpl struct {
	service *service.BudgetService
	repo    repository.BudgetRepository
}

// NewBudgetAPI cria uma nova instância da API de budget
func NewBudgetAPI(lm lifecycle.LifecycleManager, cashFlowPort service.CashFlowPort) BudgetAPI {
	repo := repository.NewSQLiteBudgetRepository(lm)
	budgetService := service.NewBudgetService(repo, cashFlowPort)

	return &BudgetAPIImpl{
		service: budgetService,
		repo:    repo,
	}
}

// CreatePlan implementa BudgetAPI.CreatePlan
func (api *BudgetAPIImpl) CreatePlan(ctx context.Context, req BudgetPlanRequest) (*BudgetPlanResponse, error) {
	plan := &domain.BudgetPlan{
		Period:      req.Period,
		Category:    req.Category,
		Planned:     req.Planned,
		Description: req.Description,
	}

	if err := api.service.CreatePlan(ctx, req.EntityID, plan); err != nil {
		return &BudgetPlanResponse{
			Success: false,
			Error:   err.Error(),
		}, err
	}

	return &BudgetPlanResponse{
		PlanID:  plan.ID,
		Success: true,
	}, nil
}

// GetPlan implementa BudgetAPI.GetPlan
func (api *BudgetAPIImpl) GetPlan(ctx context.Context, entityID, planID string) (*BudgetExecution, error) {
	plan, err := api.service.GetPlan(ctx, entityID, planID)
	if err != nil {
		return nil, err
	}

	// Obter execução do plano
	executions, err := api.service.GetExecutionReport(ctx, entityID, plan.Period)
	if err != nil {
		return nil, err
	}

	// Encontrar execução para este plano
	for _, exec := range executions {
		if exec.Plan.ID == planID {
			return &BudgetExecution{
				PlanID:      exec.Plan.ID,
				Period:      exec.Plan.Period,
				Category:    exec.Plan.Category,
				Planned:     exec.Plan.Planned,
				Executed:    exec.Executed,
				Remaining:   exec.Remaining,
				Percentage:  exec.Percentage,
				AlertStatus: exec.AlertStatus,
				Description: exec.Plan.Description,
			}, nil
		}
	}

	// Se não encontrou execução, retornar apenas o plano
	return &BudgetExecution{
		PlanID:      plan.ID,
		Period:      plan.Period,
		Category:    plan.Category,
		Planned:     plan.Planned,
		Executed:    0,
		Remaining:   plan.Planned,
		Percentage:  0,
		AlertStatus: string(domain.AlertStatusSafe),
		Description: plan.Description,
	}, nil
}

// DeletePlan implementa BudgetAPI.DeletePlan
func (api *BudgetAPIImpl) DeletePlan(ctx context.Context, entityID, planID string) (*BudgetPlanResponse, error) {
	if err := api.service.DeletePlan(ctx, entityID, planID); err != nil {
		return &BudgetPlanResponse{
			Success: false,
			Error:   err.Error(),
		}, err
	}

	return &BudgetPlanResponse{
		PlanID:  planID,
		Success: true,
	}, nil
}

// GetExecutionReport implementa BudgetAPI.GetExecutionReport
func (api *BudgetAPIImpl) GetExecutionReport(ctx context.Context, entityID, period string) ([]*BudgetExecution, error) {
	executions, err := api.service.GetExecutionReport(ctx, entityID, period)
	if err != nil {
		return nil, err
	}

	result := make([]*BudgetExecution, len(executions))
	for i, exec := range executions {
		result[i] = &BudgetExecution{
			PlanID:      exec.Plan.ID,
			Period:      exec.Plan.Period,
			Category:    exec.Plan.Category,
			Planned:     exec.Plan.Planned,
			Executed:    exec.Executed,
			Remaining:   exec.Remaining,
			Percentage:  exec.Percentage,
			AlertStatus: exec.AlertStatus,
			Description: exec.Plan.Description,
		}
	}

	return result, nil
}

// GetExecutionSummary implementa BudgetAPI.GetExecutionSummary
func (api *BudgetAPIImpl) GetExecutionSummary(ctx context.Context, entityID, period string) (*BudgetSummary, error) {
	summary, err := api.service.GetExecutionSummary(ctx, entityID, period)
	if err != nil {
		return nil, err
	}

	return &BudgetSummary{
		Period:          summary.Period,
		TotalPlanned:    summary.TotalPlanned,
		TotalExecuted:   summary.TotalExecuted,
		TotalRemaining:  summary.TotalRemaining,
		Percentage:      summary.Percentage,
		OverallStatus:   summary.OverallStatus,
		SafeCount:       summary.SafeCount,
		WarningCount:    summary.WarningCount,
		ExceededCount:   summary.ExceededCount,
		TotalCategories: summary.TotalCategories,
	}, nil
}

// GetCategories implementa BudgetAPI.GetCategories
func (api *BudgetAPIImpl) GetCategories(ctx context.Context) ([]*Category, error) {
	domainCategories := api.service.GetCategories(ctx)
	categories := make([]*Category, len(domainCategories))

	for i, cat := range domainCategories {
		categories[i] = &Category{
			ID:    string(cat),
			Label: domain.GetCategoryLabel(string(cat)),
		}
	}

	return categories, nil
}

// GetAvailablePeriods implementa BudgetAPI.GetAvailablePeriods
func (api *BudgetAPIImpl) GetAvailablePeriods(ctx context.Context, entityID string) ([]string, error) {
	// Implementação simplificada: retorna períodos dos últimos 12 meses
	var periods []string
	now := time.Now()

	for i := 0; i < 12; i++ {
		date := now.AddDate(0, -i, 0)
		period := date.Format("2006-01")
		periods = append(periods, period)
	}

	return periods, nil
}

// Helper function para converter domain.BudgetExecution para BudgetExecution
func domainToBudgetExecution(exec domain.BudgetExecution) *BudgetExecution {
	return &BudgetExecution{
		PlanID:      exec.Plan.ID,
		Period:      exec.Plan.Period,
		Category:    exec.Plan.Category,
		Planned:     exec.Plan.Planned,
		Executed:    exec.Executed,
		Remaining:   exec.Remaining,
		Percentage:  exec.Percentage,
		AlertStatus: exec.AlertStatus,
		Description: exec.Plan.Description,
	}
}

// Helper function para converter service.BudgetSummary para BudgetSummary
func serviceToBudgetSummary(summary *service.BudgetSummary) *BudgetSummary {
	return &BudgetSummary{
		Period:          summary.Period,
		TotalPlanned:    summary.TotalPlanned,
		TotalExecuted:   summary.TotalExecuted,
		TotalRemaining:  summary.TotalRemaining,
		Percentage:      summary.Percentage,
		OverallStatus:   summary.OverallStatus,
		SafeCount:       summary.SafeCount,
		WarningCount:    summary.WarningCount,
		ExceededCount:   summary.ExceededCount,
		TotalCategories: summary.TotalCategories,
	}
}

// MockCashFlowPort implementa CashFlowPort para testes
type MockCashFlowPort struct{}

func (m *MockCashFlowPort) GetCashFlow(entityID string, startDate, endDate time.Time) (*service.CashFlowResponse, error) {
	// Implementação mock para testes
	return &service.CashFlowResponse{
		EntityID:    entityID,
		TotalCredit: 0,
		TotalDebit:  0,
		Balance:     0,
		Entries:     []service.CashEntry{},
		PeriodStart: startDate,
		PeriodEnd:   endDate,
	}, nil
}

// NewMockBudgetAPI cria uma API de budget com mock para testes
func NewMockBudgetAPI(lm lifecycle.LifecycleManager) BudgetAPI {
	repo := repository.NewSQLiteBudgetRepository(lm)
	mockCashFlow := &MockCashFlowPort{}
	budgetService := service.NewBudgetService(repo, mockCashFlow)

	return &BudgetAPIImpl{
		service: budgetService,
		repo:    repo,
	}
}
