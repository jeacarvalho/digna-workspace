package domain

import (
	"context"
	"time"
)

// BudgetPlan representa o planejamento orçamentário de uma categoria para um período
type BudgetPlan struct {
	ID          string
	EntityID    string
	Period      string // Formato: "YYYY-MM" (ex: "2026-04")
	Category    string // Categoria: "Insumos", "Energia", "Equipamentos", "Transporte", "Outros"
	Planned     int64  // Valor planejado em centavos (int64 - anti-float)
	Description string // Descrição opcional
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// BudgetExecution representa a execução do orçamento (planejado vs realizado)
type BudgetExecution struct {
	Plan        BudgetPlan
	Executed    int64  // Valor realmente gasto (buscado do Ledger/Caixa) em centavos
	Remaining   int64  // Planejado - Executado (em centavos)
	Percentage  int    // Percentual executado (0-100)
	AlertStatus string // Status de alerta: "SAFE", "WARNING", "EXCEEDED"
}

// BudgetAlertStatus define os possíveis status de alerta
type BudgetAlertStatus string

const (
	AlertStatusSafe     BudgetAlertStatus = "SAFE"     // Até 70% do planejado
	AlertStatusWarning  BudgetAlertStatus = "WARNING"  // 71% a 100% do planejado
	AlertStatusExceeded BudgetAlertStatus = "EXCEEDED" // Acima de 100% do planejado
)

// BudgetCategory define categorias pré-definidas para orçamento
type BudgetCategory string

const (
	CategoryRawMaterials BudgetCategory = "INSUMOS"      // Matérias-primas
	CategoryEnergy       BudgetCategory = "ENERGIA"      // Eletricidade, gás
	CategoryEquipment    BudgetCategory = "EQUIPAMENTOS" // Ferramentas, máquinas
	CategoryTransport    BudgetCategory = "TRANSPORTE"   // Combustível, frete
	CategoryMaintenance  BudgetCategory = "MANUTENCAO"   // Reparos, conservação
	CategoryServices     BudgetCategory = "SERVICOS"     // Terceirizados
	CategoryOther        BudgetCategory = "OUTROS"       // Outras despesas
)

// Validate valida se o BudgetPlan é válido
func (bp *BudgetPlan) Validate() error {
	if bp.EntityID == "" {
		return ErrInvalidBudgetEntity
	}
	if bp.Period == "" {
		return ErrInvalidBudgetPeriod
	}
	if bp.Category == "" {
		return ErrInvalidBudgetCategory
	}

	// Validar categoria
	validCategories := map[string]bool{
		"INSUMOS":      true,
		"ENERGIA":      true,
		"EQUIPAMENTOS": true,
		"TRANSPORTE":   true,
		"MANUTENCAO":   true,
		"SERVICOS":     true,
		"OUTROS":       true,
	}
	if !validCategories[bp.Category] {
		return ErrInvalidBudgetCategory
	}

	if bp.Planned <= 0 {
		return ErrInvalidBudgetAmount
	}

	// Validar formato do período (YYYY-MM)
	if len(bp.Period) != 7 || bp.Period[4] != '-' {
		return ErrInvalidBudgetPeriodFormat
	}

	return nil
}

// CalculateExecution calcula a execução do orçamento com base no valor executado
func (bp *BudgetPlan) CalculateExecution(executed int64) BudgetExecution {
	if executed < 0 {
		executed = 0
	}

	remaining := bp.Planned - executed
	if remaining < 0 {
		remaining = 0
	}

	// Calcular percentual (evitar divisão por zero)
	var percentage int
	var alertStatus BudgetAlertStatus

	if bp.Planned > 0 {
		percentage = int((executed * 100) / bp.Planned)

		// Determinar status de alerta
		if percentage <= 70 {
			alertStatus = AlertStatusSafe
		} else if percentage <= 100 {
			alertStatus = AlertStatusWarning
		} else {
			alertStatus = AlertStatusExceeded
			percentage = 100 // Manter 100% para exibição
		}
	} else {
		// Se não há valor planejado, considera-se safe
		alertStatus = AlertStatusSafe
	}

	return BudgetExecution{
		Plan:        *bp,
		Executed:    executed,
		Remaining:   remaining,
		Percentage:  percentage,
		AlertStatus: string(alertStatus),
	}
}

// GetCategoryLabel retorna o rótulo em português para a categoria
func GetCategoryLabel(category string) string {
	switch BudgetCategory(category) {
	case CategoryRawMaterials:
		return "Insumos"
	case CategoryEnergy:
		return "Energia"
	case CategoryEquipment:
		return "Equipamentos"
	case CategoryTransport:
		return "Transporte"
	case CategoryMaintenance:
		return "Manutenção"
	case CategoryServices:
		return "Serviços"
	case CategoryOther:
		return "Outros"
	default:
		return category
	}
}

// GetAlertStatusLabel retorna o rótulo em português para o status de alerta
func GetAlertStatusLabel(status string) string {
	switch BudgetAlertStatus(status) {
	case AlertStatusSafe:
		return "Dentro do planejado"
	case AlertStatusWarning:
		return "Atenção: perto do limite"
	case AlertStatusExceeded:
		return "Ultrapassou o planejado"
	default:
		return status
	}
}

// BudgetRepository interface para persistência de orçamentos
type BudgetRepository interface {
	// SavePlan salva ou atualiza um plano orçamentário
	SavePlan(ctx context.Context, entityID string, plan *BudgetPlan) error

	// GetPlan obtém um plano orçamentário por ID
	GetPlan(ctx context.Context, entityID, planID string) (*BudgetPlan, error)

	// ListPlansByPeriod lista todos os planos de um período
	ListPlansByPeriod(ctx context.Context, entityID, period string) ([]*BudgetPlan, error)

	// ListPlansByCategory lista planos por categoria em um período
	ListPlansByCategory(ctx context.Context, entityID, period, category string) ([]*BudgetPlan, error)

	// DeletePlan remove um plano orçamentário
	DeletePlan(ctx context.Context, entityID, planID string) error

	// GetCategories retorna as categorias disponíveis
	GetCategories(ctx context.Context) []BudgetCategory
}

// Erros de domínio para orçamento
var (
	ErrInvalidBudgetEntity       = newBudgetError("entidade do orçamento inválida")
	ErrInvalidBudgetPeriod       = newBudgetError("período do orçamento inválido")
	ErrInvalidBudgetPeriodFormat = newBudgetError("formato do período inválido (use YYYY-MM)")
	ErrInvalidBudgetCategory     = newBudgetError("categoria do orçamento inválida")
	ErrInvalidBudgetAmount       = newBudgetError("valor do orçamento inválido")
	ErrBudgetPlanNotFound        = newBudgetError("plano orçamentário não encontrado")
	ErrBudgetExecutionFailed     = newBudgetError("falha ao calcular execução do orçamento")
)

// BudgetError representa um erro de domínio do módulo budget
type BudgetError struct {
	message string
}

func newBudgetError(message string) *BudgetError {
	return &BudgetError{message: message}
}

func (e *BudgetError) Error() string {
	return e.message
}
