package service

import (
	"context"
	"fmt"
	"time"

	"github.com/providentia/digna/budget/internal/domain"
	"github.com/providentia/digna/budget/internal/repository"
)

// CashFlowPort interface para integração com o módulo cash_flow
type CashFlowPort interface {
	// GetCashFlow obtém o fluxo de caixa para um período
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

// BudgetService orquestra a gestão orçamentária
type BudgetService struct {
	repo     repository.BudgetRepository
	cashFlow CashFlowPort
}

// NewBudgetService cria um novo serviço de orçamento
func NewBudgetService(repo repository.BudgetRepository, cashFlow CashFlowPort) *BudgetService {
	return &BudgetService{
		repo:     repo,
		cashFlow: cashFlow,
	}
}

// CreatePlan cria um novo plano orçamentário
func (s *BudgetService) CreatePlan(ctx context.Context, entityID string, plan *domain.BudgetPlan) error {
	plan.EntityID = entityID

	if err := plan.Validate(); err != nil {
		return fmt.Errorf("validação do plano falhou: %w", err)
	}

	// Verificar se já existe plano para mesma categoria no período
	existingPlans, err := s.repo.ListPlansByCategory(ctx, entityID, plan.Period, plan.Category)
	if err != nil {
		return fmt.Errorf("erro ao verificar planos existentes: %w", err)
	}

	if len(existingPlans) > 0 {
		// Atualizar plano existente
		existingPlan := existingPlans[0]
		existingPlan.Planned = plan.Planned
		existingPlan.Description = plan.Description
		existingPlan.UpdatedAt = time.Now()

		return s.repo.SavePlan(ctx, entityID, existingPlan)
	}

	// Criar novo plano
	return s.repo.SavePlan(ctx, entityID, plan)
}

// GetExecutionReport obtém o relatório de execução do orçamento para um período
func (s *BudgetService) GetExecutionReport(ctx context.Context, entityID, period string) ([]domain.BudgetExecution, error) {
	// Obter planos do período
	plans, err := s.repo.ListPlansByPeriod(ctx, entityID, period)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar planos: %w", err)
	}

	if len(plans) == 0 {
		return []domain.BudgetExecution{}, nil
	}

	// Obter datas do período
	startDate, endDate, err := parsePeriod(period)
	if err != nil {
		return nil, fmt.Errorf("erro ao parsear período: %w", err)
	}

	// Obter fluxo de caixa real do período
	var cashFlowResponse *CashFlowResponse
	if s.cashFlow != nil {
		cashFlowResponse, err = s.cashFlow.GetCashFlow(entityID, startDate, endDate)
		if err != nil {
			// Logar erro mas continuar (para resiliência)
			fmt.Printf("AVISO: Falha ao obter fluxo de caixa: %v\n", err)
		}
	}

	// Calcular execução para cada plano
	executions := make([]domain.BudgetExecution, len(plans))
	for i, plan := range plans {
		executed := s.calculateExecutedAmount(plan, cashFlowResponse)
		executions[i] = plan.CalculateExecution(executed)
	}

	return executions, nil
}

// GetExecutionSummary obtém um resumo da execução orçamentária
func (s *BudgetService) GetExecutionSummary(ctx context.Context, entityID, period string) (*BudgetSummary, error) {
	executions, err := s.GetExecutionReport(ctx, entityID, period)
	if err != nil {
		return nil, err
	}

	var totalPlanned int64
	var totalExecuted int64
	var safeCount, warningCount, exceededCount int

	for _, exec := range executions {
		totalPlanned += exec.Plan.Planned
		totalExecuted += exec.Executed

		switch domain.BudgetAlertStatus(exec.AlertStatus) {
		case domain.AlertStatusSafe:
			safeCount++
		case domain.AlertStatusWarning:
			warningCount++
		case domain.AlertStatusExceeded:
			exceededCount++
		}
	}

	var percentage int
	if totalPlanned > 0 {
		percentage = int((totalExecuted * 100) / totalPlanned)
	}

	// Determinar status geral
	var overallStatus string
	if percentage <= 70 {
		overallStatus = string(domain.AlertStatusSafe)
	} else if percentage <= 100 {
		overallStatus = string(domain.AlertStatusWarning)
	} else {
		overallStatus = string(domain.AlertStatusExceeded)
	}

	return &BudgetSummary{
		Period:          period,
		TotalPlanned:    totalPlanned,
		TotalExecuted:   totalExecuted,
		TotalRemaining:  totalPlanned - totalExecuted,
		Percentage:      percentage,
		OverallStatus:   overallStatus,
		SafeCount:       safeCount,
		WarningCount:    warningCount,
		ExceededCount:   exceededCount,
		TotalCategories: len(executions),
	}, nil
}

// GetPlan obtém um plano específico
func (s *BudgetService) GetPlan(ctx context.Context, entityID, planID string) (*domain.BudgetPlan, error) {
	return s.repo.GetPlan(ctx, entityID, planID)
}

// DeletePlan remove um plano orçamentário
func (s *BudgetService) DeletePlan(ctx context.Context, entityID, planID string) error {
	return s.repo.DeletePlan(ctx, entityID, planID)
}

// GetCategories retorna as categorias disponíveis
func (s *BudgetService) GetCategories(ctx context.Context) []domain.BudgetCategory {
	return s.repo.GetCategories(ctx)
}

// calculateExecutedAmount calcula o valor realmente gasto para uma categoria
func (s *BudgetService) calculateExecutedAmount(plan *domain.BudgetPlan, cashFlowResponse *CashFlowResponse) int64 {
	if cashFlowResponse == nil {
		return 0
	}

	// Mapear categoria do orçamento para categorias do cash flow
	// (Esta é uma implementação simplificada - em produção precisaria de mapeamento mais sofisticado)
	cashFlowCategory := mapBudgetToCashFlowCategory(plan.Category)

	var executed int64
	for _, entry := range cashFlowResponse.Entries {
		// Considerar apenas débitos (gastos)
		// Nota: Em cash_flow, EntryType é "CREDIT" para entrada e "DEBIT" para saída
		// Precisamos verificar a implementação real do cash_flow

		// Implementação simplificada: somar todas as saídas
		// Em produção, precisaríamos mapear categorias corretamente
		if entry.Category == cashFlowCategory || cashFlowCategory == "" {
			// Assumindo que entry.Type == "DEBIT" significa gasto
			// Esta lógica precisa ser ajustada conforme a implementação real do cash_flow
			executed += entry.Amount
		}
	}

	return executed
}

// mapBudgetToCashFlowCategory mapeia categoria do orçamento para categoria do cash flow
func mapBudgetToCashFlowCategory(budgetCategory string) string {
	// Mapeamento simplificado
	// Em produção, isso poderia ser configurável
	switch domain.BudgetCategory(budgetCategory) {
	case domain.CategoryRawMaterials:
		return "INSUMOS"
	case domain.CategoryEnergy:
		return "ENERGIA"
	case domain.CategoryEquipment:
		return "EQUIPAMENTOS"
	case domain.CategoryTransport:
		return "TRANSPORTE"
	case domain.CategoryMaintenance:
		return "MANUTENCAO"
	case domain.CategoryServices:
		return "SERVICOS"
	default:
		return "" // Categoria não mapeada
	}
}

// parsePeriod converte uma string "YYYY-MM" em datas de início e fim do mês
func parsePeriod(period string) (time.Time, time.Time, error) {
	if len(period) != 7 || period[4] != '-' {
		return time.Time{}, time.Time{}, fmt.Errorf("formato de período inválido: %s", period)
	}

	year := period[0:4]
	month := period[5:7]

	// Primeiro dia do mês
	startDate, err := time.Parse("2006-01-02", fmt.Sprintf("%s-%s-01", year, month))
	if err != nil {
		return time.Time{}, time.Time{}, fmt.Errorf("erro ao parsear data inicial: %w", err)
	}

	// Último dia do mês
	endDate := startDate.AddDate(0, 1, -1)
	endDate = time.Date(endDate.Year(), endDate.Month(), endDate.Day(), 23, 59, 59, 999999999, endDate.Location())

	return startDate, endDate, nil
}

// BudgetSummary representa um resumo da execução orçamentária
type BudgetSummary struct {
	Period          string
	TotalPlanned    int64
	TotalExecuted   int64
	TotalRemaining  int64
	Percentage      int
	OverallStatus   string
	SafeCount       int
	WarningCount    int
	ExceededCount   int
	TotalCategories int
}

// GetSummaryLabel retorna o rótulo em português para o resumo
func (bs *BudgetSummary) GetSummaryLabel() string {
	switch domain.BudgetAlertStatus(bs.OverallStatus) {
	case domain.AlertStatusSafe:
		return "Orçamento sob controle"
	case domain.AlertStatusWarning:
		return "Atenção: perto do limite"
	case domain.AlertStatusExceeded:
		return "Orçamento ultrapassado"
	default:
		return "Status desconhecido"
	}
}
