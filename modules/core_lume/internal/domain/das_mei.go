package domain

import (
	"errors"
	"fmt"
	"time"
)

// DASMEIStatus representa o status do DAS MEI
type DASMEIStatus string

const (
	DASMEIStatusPending  DASMEIStatus = "PENDENTE"
	DASMEIStatusPaid     DASMEIStatus = "PAGO"
	DASMEIStatusOverdue  DASMEIStatus = "VENCIDO"
	DASMEIStatusCanceled DASMEIStatus = "CANCELADO"
)

// ActivityType representa o tipo de atividade do MEI
type ActivityType string

const (
	ActivityTypeCommerce ActivityType = "COMERCIO" // ICMS
	ActivityTypeService  ActivityType = "SERVICOS" // ISS
	ActivityTypeMixed    ActivityType = "MISTO"    // ICMS + ISS
)

// DASMEI representa o Documento de Arrecadação do Simples Nacional para MEI
type DASMEI struct {
	ID             string
	EntityID       string
	Competencia    string // YYYY-MM (mês de referência)
	ValorDevido    int64  // centavos - Anti-Float
	ValorPago      int64  // centavos - Anti-Float
	DataVencimento int64  // Unix timestamp
	DataPagamento  int64  // Unix timestamp (0 se não pago)
	Status         DASMEIStatus
	SalarioMinimo  int64        // Salário mínimo de referência (centavos)
	ActivityType   ActivityType // Tipo de atividade
	CreatedAt      int64        // Unix timestamp
	UpdatedAt      int64        // Unix timestamp
}

// Tabela de salário mínimo versionada por ano (valores em centavos)
var MinimumWageTable = map[int]int64{
	2024: 141200, // R$ 1.412,00
	2025: 151800, // R$ 1.518,00
	2026: 151800, // R$ 1.518,00 (ajustar quando houver decreto)
}

// Constantes de valores fixos para DAS MEI (em centavos)
const (
	ICMSFixedAmount   int64 = 100 // R$ 1,00 para comércio
	ISSFixedAmount    int64 = 500 // R$ 5,00 para serviços
	DASPercentage     int64 = 5   // 5%
	PercentageDivisor int64 = 100 // Para cálculo de porcentagem
)

var (
	ErrDASMEIInvalidEntityID    = errors.New("entity ID is required")
	ErrDASMEIInvalidCompetencia = errors.New("competencia is required and must be in format YYYY-MM")
	ErrDASMEIInvalidStatus      = errors.New("invalid DAS MEI status")
	ErrDASMEIInvalidActivity    = errors.New("invalid activity type")
	ErrDASMEIAlreadyPaid        = errors.New("DAS MEI is already paid")
	ErrDASMEINotPending         = errors.New("DAS MEI is not pending")
)

// Validate verifica se o DAS MEI está válido
func (d *DASMEI) Validate() error {
	if d.EntityID == "" {
		return ErrDASMEIInvalidEntityID
	}

	if d.Competencia == "" || len(d.Competencia) != 7 {
		return ErrDASMEIInvalidCompetencia
	}

	if !isValidDASMEIStatus(d.Status) {
		return ErrDASMEIInvalidStatus
	}

	if !isValidActivityType(d.ActivityType) {
		return ErrDASMEIInvalidActivity
	}

	if d.ValorDevido <= 0 {
		return errors.New("valor devido must be positive")
	}

	if d.SalarioMinimo <= 0 {
		return errors.New("salario minimo must be positive")
	}

	return nil
}

// IsOverdue verifica se o DAS MEI está vencido
func (d *DASMEI) IsOverdue() bool {
	if d.Status == DASMEIStatusPaid || d.Status == DASMEIStatusCanceled {
		return false
	}
	return time.Now().Unix() > d.DataVencimento
}

// CalculateAmount calcula o valor do DAS MEI baseado no salário mínimo e atividade
// Usa int64 (centavos) - Anti-Float
func CalculateDASMEIAmount(salarioMinimo int64, activity ActivityType) int64 {
	// 5% do salário mínimo
	baseAmount := (salarioMinimo * DASPercentage) / PercentageDivisor

	switch activity {
	case ActivityTypeCommerce:
		// Comércio: 5% do SM + ICMS fixo
		return baseAmount + ICMSFixedAmount
	case ActivityTypeService:
		// Serviços: 5% do SM + ISS fixo
		return baseAmount + ISSFixedAmount
	case ActivityTypeMixed:
		// Misto: 5% do SM + ICMS + ISS
		return baseAmount + ICMSFixedAmount + ISSFixedAmount
	default:
		// Padrão: apenas 5% do SM
		return baseAmount
	}
}

// GetMinimumWageForYear retorna o salário mínimo para um ano específico
func GetMinimumWageForYear(year int) int64 {
	if wage, ok := MinimumWageTable[year]; ok {
		return wage
	}
	// Retorna o valor mais recente se o ano não estiver na tabela
	return MinimumWageTable[2026]
}

// MarkAsPaid marca o DAS MEI como pago
func (d *DASMEI) MarkAsPaid() error {
	if d.Status == DASMEIStatusPaid {
		return ErrDASMEIAlreadyPaid
	}

	if d.Status == DASMEIStatusCanceled {
		return errors.New("cannot pay a canceled DAS MEI")
	}

	now := time.Now().Unix()
	d.Status = DASMEIStatusPaid
	d.ValorPago = d.ValorDevido
	d.DataPagamento = now
	d.UpdatedAt = now

	return nil
}

// Cancel cancela o DAS MEI
func (d *DASMEI) Cancel() error {
	if d.Status == DASMEIStatusPaid {
		return errors.New("cannot cancel a paid DAS MEI")
	}

	d.Status = DASMEIStatusCanceled
	d.UpdatedAt = time.Now().Unix()

	return nil
}

// UpdateStatus atualiza o status baseado na data atual
func (d *DASMEI) UpdateStatus() {
	if d.Status == DASMEIStatusPaid || d.Status == DASMEIStatusCanceled {
		return
	}

	if d.IsOverdue() {
		d.Status = DASMEIStatusOverdue
		d.UpdatedAt = time.Now().Unix()
	}
}

// String retorna representação textual do DAS MEI
func (d *DASMEI) String() string {
	return fmt.Sprintf("DASMEI{ID: %s, Competencia: %s, Valor: %d, Status: %s}",
		d.ID, d.Competencia, d.ValorDevido, d.Status)
}

// GetValorDevidoReal retorna o valor devido em reais (para exibição)
func (d *DASMEI) GetValorDevidoReal() float64 {
	return float64(d.ValorDevido) / 100.0
}

// GetValorPagoReal retorna o valor pago em reais (para exibição)
func (d *DASMEI) GetValorPagoReal() float64 {
	return float64(d.ValorPago) / 100.0
}

// GetSalarioMinimoReal retorna o salário mínimo em reais (para exibição)
func (d *DASMEI) GetSalarioMinimoReal() float64 {
	return float64(d.SalarioMinimo) / 100.0
}

// GetDueDate retorna a data de vencimento como time.Time
func (d *DASMEI) GetDueDate() time.Time {
	return time.Unix(d.DataVencimento, 0)
}

// GetPaymentDate retorna a data de pagamento como time.Time
func (d *DASMEI) GetPaymentDate() *time.Time {
	if d.DataPagamento == 0 {
		return nil
	}
	t := time.Unix(d.DataPagamento, 0)
	return &t
}

// isValidDASMEIStatus verifica se o status é válido
func isValidDASMEIStatus(status DASMEIStatus) bool {
	switch status {
	case DASMEIStatusPending, DASMEIStatusPaid, DASMEIStatusOverdue, DASMEIStatusCanceled:
		return true
	}
	return false
}

// isValidActivityType verifica se o tipo de atividade é válido
func isValidActivityType(activity ActivityType) bool {
	switch activity {
	case ActivityTypeCommerce, ActivityTypeService, ActivityTypeMixed:
		return true
	}
	return false
}

// ParseCompetencia extrai ano e mês da competência
func ParseCompetencia(competencia string) (year, month int, err error) {
	if len(competencia) != 7 || competencia[4] != '-' {
		return 0, 0, ErrDASMEIInvalidCompetencia
	}

	_, err = fmt.Sscanf(competencia, "%d-%d", &year, &month)
	if err != nil {
		return 0, 0, ErrDASMEIInvalidCompetencia
	}

	if month < 1 || month > 12 {
		return 0, 0, ErrDASMEIInvalidCompetencia
	}

	return year, month, nil
}

// CalculateDueDate calcula a data de vencimento (dia 20 ou dia útil anterior)
func CalculateDueDate(year, month int) int64 {
	// Dia 20 do mês
	dueDate := time.Date(year, time.Month(month), 20, 0, 0, 0, 0, time.Local)

	// Se for fim de semana, volta para sexta-feira
	for dueDate.Weekday() == time.Saturday || dueDate.Weekday() == time.Sunday {
		dueDate = dueDate.AddDate(0, 0, -1)
	}

	return dueDate.Unix()
}

// GenerateCompetencia gera a string de competência para o mês atual
func GenerateCompetencia(t time.Time) string {
	return fmt.Sprintf("%04d-%02d", t.Year(), t.Month())
}
