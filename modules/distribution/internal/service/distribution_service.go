// Package service implementa a camada de aplicação para distribuição de sobras
// Coordena o cálculo, aprovação e execução da distribuição
package service

import (
	"fmt"
	"time"

	"github.com/providentia/digna/distribution/internal/domain"
	"github.com/providentia/digna/reporting/pkg/surplus"
)

// CalculatorInterface define a interface para cálculo de sobras
type CalculatorInterface interface {
	CalculateSocialSurplus(entityID string) (*surplus.SurplusCalculation, error)
}

// ErrInvalidSurplus indica que não há sobra para distribuir
var ErrInvalidSurplus = fmt.Errorf("não há sobra disponível para distribuição")

// ErrPendingDecision indica que a distribuição aguarda aprovação da assembleia
var ErrPendingDecision = fmt.Errorf("distribuição aguarda aprovação em assembleia")

// ErrAlreadyExecuted indica que a distribuição já foi executada
var ErrAlreadyExecuted = fmt.Errorf("distribuição já foi executada")

// DistributionService orquestra o fluxo completo de distribuição
type DistributionService struct {
	distributionRepo domain.DistributionRepository
	assemblyRepo     domain.AssemblyRepository
	calculator       CalculatorInterface
}

// NewDistributionService cria um novo serviço de distribuição
func NewDistributionService(
	distRepo domain.DistributionRepository,
	assemblyRepo domain.AssemblyRepository,
	calculator *surplus.Calculator,
) *DistributionService {
	return &DistributionService{
		distributionRepo: distRepo,
		assemblyRepo:     assemblyRepo,
		calculator:       calculator,
	}
}

// NewDistributionServiceWithCalculator cria serviço com interface genérica (para testes)
func NewDistributionServiceWithCalculator(
	distRepo domain.DistributionRepository,
	assemblyRepo domain.AssemblyRepository,
	calculator CalculatorInterface,
) *DistributionService {
	return &DistributionService{
		distributionRepo: distRepo,
		assemblyRepo:     assemblyRepo,
		calculator:       calculator,
	}
}

// CalculateSurplus calcula a proposta de distribuição seguindo a Lei Paul Singer:
// - 10% Reserva Legal
// - 5% FATES
// - 85% Distribuível proporcional às horas trabalhadas
func (s *DistributionService) CalculateSurplus(entityID, period string) (*domain.SurplusCalculation, error) {
	// Calcula sobras usando o reporting
	surplusCalc, err := s.calculator.CalculateSocialSurplus(entityID)
	if err != nil {
		return nil, fmt.Errorf("erro ao calcular sobras: %w", err)
	}

	// Verifica se há sobra positiva
	if surplusCalc.TotalSurplus <= 0 {
		return nil, ErrInvalidSurplus
	}

	// Calcula reservas obrigatórias
	reservaLegal := (surplusCalc.TotalSurplus * 10) / 100           // 10%
	fates := (surplusCalc.TotalSurplus * 5) / 100                   // 5%
	distribuivel := surplusCalc.TotalSurplus - reservaLegal - fates // 85%

	// Calcula valores para cada membro (proporcional às horas)
	members := make([]domain.MemberCalculation, 0, len(surplusCalc.Members))
	for _, m := range surplusCalc.Members {
		// Cálculo proporcional: (distribuível × horas_membro) / horas_totais
		amount := int64(0)
		if surplusCalc.TotalMinutes > 0 {
			amount = (distribuivel * m.Minutes) / surplusCalc.TotalMinutes
		}

		members = append(members, domain.MemberCalculation{
			MemberID:   m.MemberID,
			Minutes:    m.Minutes,
			Percentage: m.Percentage,
			Amount:     amount,
		})
	}

	return &domain.SurplusCalculation{
		EntityID:     entityID,
		Period:       period,
		TotalSurplus: surplusCalc.TotalSurplus,
		ReservaLegal: reservaLegal,
		FATES:        fates,
		Distribuivel: distribuivel,
		TotalMinutes: surplusCalc.TotalMinutes,
		Members:      members,
	}, nil
}

// CreateDistributionProposal cria uma proposta de distribuição
// NÃO gera lançamentos - apenas salva a proposta para aprovação
func (s *DistributionService) CreateDistributionProposal(
	entityID string,
	period string,
	calculation *domain.SurplusCalculation,
) (*domain.Distribution, error) {

	// Cria a distribuição
	dist := &domain.Distribution{
		EntityID:     entityID,
		Period:       period,
		Status:       domain.StatusCalculated,
		TotalSurplus: calculation.TotalSurplus,
		ReservaLegal: calculation.ReservaLegal,
		FATES:        calculation.FATES,
		Distribuivel: calculation.Distribuivel,
		CreatedAt:    time.Now(),
		Members:      make([]domain.DistributionMember, 0, len(calculation.Members)),
	}

	// Converte membros calculados para entidades
	for _, m := range calculation.Members {
		dist.Members = append(dist.Members, domain.DistributionMember{
			MemberID:   m.MemberID,
			Minutes:    m.Minutes,
			Percentage: m.Percentage,
			Amount:     m.Amount,
		})
	}

	// Salva no repositório
	id, err := s.distributionRepo.Save(nil, dist)
	if err != nil {
		return nil, fmt.Errorf("erro ao salvar distribuição: %w", err)
	}
	dist.ID = id

	// Salva membros individuais
	for i := range dist.Members {
		member := &dist.Members[i]
		member.DistributionID = dist.ID
		memberID, err := s.distributionRepo.SaveMember(nil, member)
		if err != nil {
			return nil, fmt.Errorf("erro ao salvar membro: %w", err)
		}
		member.ID = memberID
	}

	return dist, nil
}

// SubmitForApproval submete a distribuição para aprovação em assembleia
// A distribuição fica no status PENDING_APPROVAL
func (s *DistributionService) SubmitForApproval(distributionID int64) error {
	dist, err := s.distributionRepo.FindByID(nil, distributionID)
	if err != nil {
		return fmt.Errorf("distribuição não encontrada: %w", err)
	}

	if dist.Status != domain.StatusCalculated {
		return fmt.Errorf("distribuição deve estar no status CALCULATED")
	}

	// Atualiza para aguardar aprovação
	return s.distributionRepo.UpdateStatus(nil, distributionID, domain.StatusPendingApproval)
}

// ApproveWithDecision aprova a distribuição com uma decisão de assembleia
// Requer que a decisão esteja registrada e aprovada
func (s *DistributionService) ApproveWithDecision(distributionID int64, decisionID int64) error {
	// Verifica se a decisão existe e foi aprovada
	decision, err := s.assemblyRepo.FindByID(nil, decisionID)
	if err != nil {
		return fmt.Errorf("decisão não encontrada: %w", err)
	}

	if decision.Status != "APPROVED" {
		return fmt.Errorf("decisão deve estar aprovada")
	}

	dist, err := s.distributionRepo.FindByID(nil, distributionID)
	if err != nil {
		return fmt.Errorf("distribuição não encontrada: %w", err)
	}

	if dist.Status != domain.StatusPendingApproval {
		return fmt.Errorf("distribuição deve estar aguardando aprovação")
	}

	// Atualiza com referência à decisão
	dist.AssemblyDecisionID = decisionID
	if err := s.distributionRepo.UpdateAssemblyDecisionID(nil, distributionID, decisionID); err != nil {
		return err
	}
	if err := s.distributionRepo.UpdateStatus(nil, distributionID, domain.StatusApproved); err != nil {
		return err
	}

	return nil
}

// ExecuteDistribution executa a distribuição aprovada
// GERA OS LANÇAMENTOS CONTÁBEIS AUTOMÁTICOS
func (s *DistributionService) ExecuteDistribution(distributionID int64) error {
	dist, err := s.distributionRepo.FindByID(nil, distributionID)
	if err != nil {
		return fmt.Errorf("distribuição não encontrada: %w", err)
	}

	if dist.Status == domain.StatusExecuted {
		return ErrAlreadyExecuted
	}

	if dist.Status != domain.StatusApproved {
		return ErrPendingDecision
	}

	// Executa transação: gera todos os lançamentos contábeis
	now := time.Now()

	// Lançamento 1: Sobras do Exercício → Reserva Legal (10%)
	if _, err := s.createLedgerEntry(dist.EntityID, "Reserva Legal",
		[]domain.Posting{
			{AccountID: 9, Amount: -dist.ReservaLegal, Description: "Reserva Legal (10%)"}, // Débito
			{AccountID: 10, Amount: dist.ReservaLegal, Description: "Reserva Legal (10%)"}, // Crédito
		}); err != nil {
		return fmt.Errorf("erro ao lançar reserva legal: %w", err)
	}

	// Lançamento 2: Sobras → FATES (5%)
	if _, err := s.createLedgerEntry(dist.EntityID, "FATES",
		[]domain.Posting{
			{AccountID: 9, Amount: -dist.FATES, Description: "FATES (5%)"},
			{AccountID: 11, Amount: dist.FATES, Description: "FATES (5%)"},
		}); err != nil {
		return fmt.Errorf("erro ao lançar FATES: %w", err)
	}

	// Lançamento 3: Distribuição aos membros (proporcional)
	for i := range dist.Members {
		member := &dist.Members[i]

		// Cria conta de Capital Social para o membro (AccountID fictício, em produção seria lookup)
		memberAccountID := int64(100 + i) // Simulação: contas 100, 101, 102...

		entryID, err := s.createLedgerEntry(dist.EntityID,
			fmt.Sprintf("Distribuição a %s", member.MemberID),
			[]domain.Posting{
				{AccountID: 9, Amount: -member.Amount, Description: "Capital Social de Trabalho"},
				{AccountID: memberAccountID, Amount: member.Amount, Description: "Capital Social de Trabalho"},
			})
		if err != nil {
			return fmt.Errorf("erro ao lançar para membro %s: %w", member.MemberID, err)
		}

		member.LedgerEntryID = entryID

		// Atualiza membro no repositório
		if _, err := s.distributionRepo.SaveMember(nil, member); err != nil {
			return fmt.Errorf("erro ao salvar membro: %w", err)
		}
	}

	// Marca como executada
	if err := s.distributionRepo.MarkAsExecuted(nil, distributionID, now); err != nil {
		return fmt.Errorf("erro ao marcar como executada: %w", err)
	}

	return nil
}

// createLedgerEntry cria um lançamento contábil com partidas dobradas
func (s *DistributionService) createLedgerEntry(entityID, description string, postings []domain.Posting) (int64, error) {
	// Valida partidas dobradas (soma zero)
	var total int64
	for _, p := range postings {
		total += p.Amount
	}
	if total != 0 {
		return 0, fmt.Errorf("partidas não balanceiam: soma = %d", total)
	}

	entry := &domain.LedgerEntry{
		EntityID:    entityID,
		Reference:   fmt.Sprintf("DIST-%d", time.Now().Unix()),
		Description: description,
		Date:        time.Now(),
		Postings:    postings,
	}

	return s.distributionRepo.SaveLedgerEntry(nil, entityID, entry)
}

// GetDistributionReport retorna relatório completo de uma distribuição
func (s *DistributionService) GetDistributionReport(distributionID int64) (*domain.Distribution, error) {
	dist, err := s.distributionRepo.FindByID(nil, distributionID)
	if err != nil {
		return nil, err
	}

	// Carrega membros
	members, err := s.distributionRepo.FindMembersByDistribution(nil, distributionID)
	if err != nil {
		return nil, err
	}
	dist.Members = members

	return dist, nil
}
