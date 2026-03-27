package domain

import (
	"errors"
	"fmt"
	"time"
)

// FinalidadeCredito representa a finalidade do crédito
type FinalidadeCredito string

const (
	FinalidadeCapitalGiro FinalidadeCredito = "CAPITAL_GIRO"
	FinalidadeEquipamento FinalidadeCredito = "EQUIPAMENTO"
	FinalidadeReforma     FinalidadeCredito = "REFORMA"
	FinalidadeOutro       FinalidadeCredito = "OUTRO"
)

// TipoEntidade representa o tipo da entidade
type TipoEntidade string

const (
	TipoEntidadeMEI         TipoEntidade = "MEI"
	TipoEntidadeME          TipoEntidade = "ME"
	TipoEntidadeEPP         TipoEntidade = "EPP"
	TipoEntidadeCooperativa TipoEntidade = "Cooperativa"
	TipoEntidadeOSC         TipoEntidade = "OSC"
	TipoEntidadeOSCIP       TipoEntidade = "OSCIP"
	TipoEntidadePF          TipoEntidade = "PF"
)

// EligibilityProfile representa o perfil de elegibilidade da entidade
type EligibilityProfile struct {
	ID       string // UUID
	EntityID string // Vínculo com entidade (único por entidade)

	// Dados já capturados pelo ERP (referência, não duplicação)
	CNPJ             string // Copiado de Enterprise
	CNAE             string // Copiado de Enterprise
	Municipio        string // Copiado de Enterprise
	UF               string // Copiado de Enterprise
	FaturamentoAnual int64  // int64 - Anti-Float (centavos)
	RegimeTributario string
	DataAbertura     int64 // Unix timestamp
	SituacaoFiscal   string

	// CAMPOS COMPLEMENTARES (preenchimento único)
	InscritoCadUnico    bool              // Habilita programas sociais
	SocioMulher         bool              // Prioridade em linhas com foco de gênero
	InadimplenciaAtiva  bool              // Direciona ao Desenrola antes de crédito novo
	FinalidadeCredito   FinalidadeCredito // Enum: CAPITAL_GIRO, EQUIPAMENTO, REFORMA, OUTRO
	ValorNecessario     int64             // int64 - Anti-Float (centavos)
	TipoEntidade        TipoEntidade      // Enum: MEI, ME, EPP, Cooperativa, OSC, OSCIP, PF
	ContabilidadeFormal bool              // Requisito de alguns programas

	// Metadados
	PreenchidoEm  int64  // Unix timestamp - primeiro preenchimento
	AtualizadoEm  int64  // Unix timestamp - última atualização
	PreenchidoPor string // ID do usuário que preencheu

	CreatedAt int64
	UpdatedAt int64
}

var (
	ErrInvalidEntityID          = errors.New("entity ID is required")
	ErrInvalidFinalidadeCredito = errors.New("finalidade do crédito inválida")
	ErrInvalidTipoEntidade      = errors.New("tipo de entidade inválido")
	ErrValorNecessarioInvalido  = errors.New("valor necessário deve ser maior que zero quando finalidade é específica")
	ErrProfileNotFound          = errors.New("perfil de elegibilidade não encontrado")
)

// Validate valida o perfil de elegibilidade
func (e *EligibilityProfile) Validate() error {
	if e.EntityID == "" {
		return ErrInvalidEntityID
	}

	if e.FinalidadeCredito != "" && !isValidFinalidadeCredito(e.FinalidadeCredito) {
		return ErrInvalidFinalidadeCredito
	}

	if e.TipoEntidade != "" && !isValidTipoEntidade(e.TipoEntidade) {
		return ErrInvalidTipoEntidade
	}

	// ValorNecessario é obrigatório se FinalidadeCredito for específica (não OUTRO)
	if e.FinalidadeCredito != "" && e.FinalidadeCredito != FinalidadeOutro {
		if e.ValorNecessario <= 0 {
			return ErrValorNecessarioInvalido
		}
	}

	return nil
}

// IsComplete verifica se todos os campos obrigatórios estão preenchidos
func (e *EligibilityProfile) IsComplete() bool {
	if e.EntityID == "" {
		return false
	}

	if e.FinalidadeCredito == "" {
		return false
	}

	if e.TipoEntidade == "" {
		return false
	}

	// Se finalidade for específica, valor necessário deve ser > 0
	if e.FinalidadeCredito != FinalidadeOutro && e.ValorNecessario <= 0 {
		return false
	}

	return true
}

// GetCompletionPercent retorna a porcentagem de campos preenchidos
func (e *EligibilityProfile) GetCompletionPercent() float64 {
	totalFields := 7 // Campos complementares obrigatórios
	filledFields := 0

	if e.InscritoCadUnico { // bool tem valor default
		filledFields++
	}
	if e.SocioMulher {
		filledFields++
	}
	if e.InadimplenciaAtiva {
		filledFields++
	}
	if e.FinalidadeCredito != "" {
		filledFields++
	}
	if e.ValorNecessario > 0 {
		filledFields++
	}
	if e.TipoEntidade != "" {
		filledFields++
	}
	if e.ContabilidadeFormal {
		filledFields++
	}

	return float64(filledFields) / float64(totalFields) * 100.0
}

// CanEdit verifica se o usuário pode editar o perfil (apenas coordenadores)
func (e *EligibilityProfile) CanEdit(userRole string) bool {
	return userRole == "COORDINATOR"
}

// String retorna uma representação textual do perfil
func (e *EligibilityProfile) String() string {
	return fmt.Sprintf("EligibilityProfile{EntityID: %s, Complete: %v, %d%%}",
		e.EntityID, e.IsComplete(), int(e.GetCompletionPercent()))
}

// GetValorNecessarioReal retorna o valor necessário em reais
func (e *EligibilityProfile) GetValorNecessarioReal() float64 {
	return float64(e.ValorNecessario) / 100.0
}

// GetFaturamentoAnualReal retorna o faturamento anual em reais
func (e *EligibilityProfile) GetFaturamentoAnualReal() float64 {
	return float64(e.FaturamentoAnual) / 100.0
}

// Update atualiza campos editáveis do perfil
func (e *EligibilityProfile) Update(input EligibilityInput, userID string) error {
	now := time.Now().Unix()

	if input.InscritoCadUnico != nil {
		e.InscritoCadUnico = *input.InscritoCadUnico
	}
	if input.SocioMulher != nil {
		e.SocioMulher = *input.SocioMulher
	}
	if input.InadimplenciaAtiva != nil {
		e.InadimplenciaAtiva = *input.InadimplenciaAtiva
	}
	if input.FinalidadeCredito != nil && *input.FinalidadeCredito != "" {
		e.FinalidadeCredito = FinalidadeCredito(*input.FinalidadeCredito)
	}
	if input.ValorNecessario != nil {
		e.ValorNecessario = *input.ValorNecessario
	}
	if input.TipoEntidade != nil && *input.TipoEntidade != "" {
		e.TipoEntidade = TipoEntidade(*input.TipoEntidade)
	}
	if input.ContabilidadeFormal != nil {
		e.ContabilidadeFormal = *input.ContabilidadeFormal
	}

	// Atualizar metadados
	if e.PreenchidoEm == 0 {
		e.PreenchidoEm = now
		e.PreenchidoPor = userID
	}
	e.AtualizadoEm = now
	e.UpdatedAt = now

	return e.Validate()
}

// EligibilityInput representa os campos editáveis do perfil
type EligibilityInput struct {
	InscritoCadUnico    *bool   `json:"inscrito_cad_unico,omitempty"`
	SocioMulher         *bool   `json:"socio_mulher,omitempty"`
	InadimplenciaAtiva  *bool   `json:"inadimplencia_ativa,omitempty"`
	FinalidadeCredito   *string `json:"finalidade_credito,omitempty"`
	ValorNecessario     *int64  `json:"valor_necessario,omitempty"`
	TipoEntidade        *string `json:"tipo_entidade,omitempty"`
	ContabilidadeFormal *bool   `json:"contabilidade_formal,omitempty"`
}

// Helper functions
func isValidFinalidadeCredito(f FinalidadeCredito) bool {
	switch f {
	case FinalidadeCapitalGiro, FinalidadeEquipamento, FinalidadeReforma, FinalidadeOutro:
		return true
	}
	return false
}

func isValidTipoEntidade(t TipoEntidade) bool {
	switch t {
	case TipoEntidadeMEI, TipoEntidadeME, TipoEntidadeEPP,
		TipoEntidadeCooperativa, TipoEntidadeOSC, TipoEntidadeOSCIP, TipoEntidadePF:
		return true
	}
	return false
}

// GetFinalidadeCreditoLabel retorna o label da finalidade
func GetFinalidadeCreditoLabel(f FinalidadeCredito) string {
	labels := map[FinalidadeCredito]string{
		FinalidadeCapitalGiro: "Capital de Giro",
		FinalidadeEquipamento: "Equipamento",
		FinalidadeReforma:     "Reforma",
		FinalidadeOutro:       "Outro",
	}
	if label, ok := labels[f]; ok {
		return label
	}
	return string(f)
}

// GetTipoEntidadeLabel retorna o label do tipo de entidade
func GetTipoEntidadeLabel(t TipoEntidade) string {
	labels := map[TipoEntidade]string{
		TipoEntidadeMEI:         "Microempreendedor Individual (MEI)",
		TipoEntidadeME:          "Microempresa (ME)",
		TipoEntidadeEPP:         "Empresa de Pequeno Porte (EPP)",
		TipoEntidadeCooperativa: "Cooperativa",
		TipoEntidadeOSC:         "Organização da Sociedade Civil (OSC)",
		TipoEntidadeOSCIP:       "OSCIP",
		TipoEntidadePF:          "Pessoa Física",
	}
	if label, ok := labels[t]; ok {
		return label
	}
	return string(t)
}
