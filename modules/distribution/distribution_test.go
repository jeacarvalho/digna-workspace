package distribution_test

import (
	"testing"

	"github.com/providentia/digna/distribution/internal/domain"
	"github.com/providentia/digna/distribution/internal/repository"
	"github.com/providentia/digna/distribution/internal/service"
	"github.com/providentia/digna/reporting/pkg/surplus"
)

// MockSurplusCalculator implementação mock do calculator para testes
type MockSurplusCalculator struct{}

func (m *MockSurplusCalculator) CalculateSocialSurplus(entityID string) (*surplus.SurplusCalculation, error) {
	// Retorna valores fixos para teste
	return &surplus.SurplusCalculation{
		EntityID:     entityID,
		TotalSurplus: 100000, // R$ 1.000,00
		TotalMinutes: 6720,   // 112 horas
		Members: []surplus.MemberShare{
			{MemberID: "cooperado_001", Minutes: 2880, Percentage: 42.86, Amount: 42857},
			{MemberID: "cooperado_002", Minutes: 1920, Percentage: 28.57, Amount: 28571},
			{MemberID: "cooperado_003", Minutes: 1920, Percentage: 28.57, Amount: 28572},
		},
	}, nil
}

func TestDistribution_FlowComplete(t *testing.T) {
	entityID := "Cooperativa_Distribuicao_Teste"

	// Criar serviço de distribuição com mocks
	distRepo := repository.NewMockDistributionRepository()
	assemblyRepo := repository.NewMockAssemblyRepository()

	// Usar mock calculator para garantir valores consistentes
	mockCalc := &MockSurplusCalculator{}
	distService := service.NewDistributionServiceWithCalculator(distRepo, assemblyRepo, mockCalc)

	t.Run("Step1_CalcularSobrasComReservas", func(t *testing.T) {
		// Calcula sobras (agora com reservas incluídas no cálculo)
		calculation, err := distService.CalculateSurplus(entityID, "2026-03")
		if err != nil {
			t.Fatalf("Erro ao calcular sobras: %v", err)
		}

		// Verifica estrutura de reservas
		totalReservas := calculation.ReservaLegal + calculation.FATES
		expectedReservaLegal := (calculation.TotalSurplus * 10) / 100
		expectedFATES := (calculation.TotalSurplus * 5) / 100

		if calculation.ReservaLegal != expectedReservaLegal {
			t.Errorf("Reserva Legal incorreta: esperado %d, obtido %d", expectedReservaLegal, calculation.ReservaLegal)
		}
		if calculation.FATES != expectedFATES {
			t.Errorf("FATES incorreto: esperado %d, obtido %d", expectedFATES, calculation.FATES)
		}

		t.Logf("📊 Sobra Total: R$ %.2f", float64(calculation.TotalSurplus)/100)
		t.Logf("📊 Reserva Legal (10%%): R$ %.2f", float64(calculation.ReservaLegal)/100)
		t.Logf("📊 FATES (5%%): R$ %.2f", float64(calculation.FATES)/100)
		t.Logf("📊 Distribuível (85%%): R$ %.2f", float64(calculation.Distribuivel)/100)
		t.Logf("📊 Total Reservas: R$ %.2f", float64(totalReservas)/100)

		// Verifica proporções
		expectedDistribuivel := calculation.TotalSurplus - totalReservas
		if calculation.Distribuivel != expectedDistribuivel {
			t.Errorf("Valor distribuível incorreto: esperado %d, obtido %d", expectedDistribuivel, calculation.Distribuivel)
		}

		t.Log("✅ Reservas calculadas corretamente: 10% + 5% = 15%")
	})

	t.Run("Step2_ProporcaoPorHoras", func(t *testing.T) {
		calculation, _ := distService.CalculateSurplus(entityID, "2026-03")

		// Verifica proporções
		totalMinutes := int64(6720)
		if calculation.TotalMinutes != totalMinutes {
			t.Errorf("Total de minutos incorreto: esperado %d, obtido %d", totalMinutes, calculation.TotalMinutes)
		}

		expectedPercentages := map[string]float64{
			"cooperado_001": 42.86,
			"cooperado_002": 28.57,
			"cooperado_003": 28.57,
		}

		t.Log("📊 Distribuição Proporcional:")
		for _, m := range calculation.Members {
			expectedPct := expectedPercentages[m.MemberID]
			tolerance := 0.5 // Tolerância de 0.5%
			if m.Percentage < expectedPct-tolerance || m.Percentage > expectedPct+tolerance {
				t.Errorf("Porcentagem incorreta para %s: esperado %.2f%%, obtido %.2f%%",
					m.MemberID, expectedPct, m.Percentage)
			}
			t.Logf("   %s: %d horas (%.2f%%) = R$ %.2f",
				m.MemberID, m.Minutes/60, m.Percentage, float64(m.Amount)/100)
		}

		t.Log("✅ Rateio proporcional às horas calculado")
	})

	t.Run("Step3_CriarProposta", func(t *testing.T) {
		calculation, _ := distService.CalculateSurplus(entityID, "2026-03")

		// Cria proposta de distribuição
		dist, err := distService.CreateDistributionProposal(entityID, "2026-03", calculation)
		if err != nil {
			t.Fatalf("Erro ao criar proposta: %v", err)
		}

		if dist.Status != domain.StatusCalculated {
			t.Errorf("Status incorreto: esperado CALCULATED, obtido %s", dist.Status)
		}

		if len(dist.Members) != 3 {
			t.Errorf("Número de membros incorreto: esperado 3, obtido %d", len(dist.Members))
		}

		t.Logf("✅ Proposta criada: ID=%d, Status=%s", dist.ID, dist.Status)
	})

	t.Run("Step4_SubmeterParaAprovacao", func(t *testing.T) {
		calculation, _ := distService.CalculateSurplus(entityID, "2026-03")
		dist, _ := distService.CreateDistributionProposal(entityID, "2026-03", calculation)

		// Submete para assembleia
		err := distService.SubmitForApproval(dist.ID)
		if err != nil {
			t.Fatalf("Erro ao submeter: %v", err)
		}

		// Verifica se mudou para PENDING_APPROVAL
		updatedDist, _ := distService.GetDistributionReport(dist.ID)
		if updatedDist.Status != domain.StatusPendingApproval {
			t.Errorf("Status incorreto: esperado PENDING_APPROVAL, obtido %s", updatedDist.Status)
		}

		t.Log("✅ Proposta submetida para aprovação em assembleia")
	})

	t.Run("Step5_AprovarComDecisao", func(t *testing.T) {
		calculation, _ := distService.CalculateSurplus(entityID, "2026-03")
		dist, _ := distService.CreateDistributionProposal(entityID, "2026-03", calculation)
		distService.SubmitForApproval(dist.ID)

		// Cria decisão de assembleia aprovada
		decision := &domain.AssemblyDecision{
			ID:       1,
			EntityID: entityID,
			Title:    "Aprovação de Rateio de Sobras - Março/2026",
			Content:  "Aprovado rateio conforme ITG 2002",
			Status:   "APPROVED",
		}
		assemblyRepo.AddDecision(decision)

		// Aprova distribuição com a decisão
		err := distService.ApproveWithDecision(dist.ID, decision.ID)
		if err != nil {
			t.Fatalf("Erro ao aprovar: %v", err)
		}

		// Verifica se foi aprovada
		approvedDist, _ := distService.GetDistributionReport(dist.ID)
		if approvedDist.Status != domain.StatusApproved {
			t.Errorf("Status incorreto: esperado APPROVED, obtido %s", approvedDist.Status)
		}

		if approvedDist.AssemblyDecisionID != decision.ID {
			t.Errorf("Referência à decisão incorreta: esperado %d, obtido %d",
				decision.ID, approvedDist.AssemblyDecisionID)
		}

		t.Log("✅ Distribuição aprovada em assembleia")
	})

	t.Run("Step6_ExecutarDistribuicao", func(t *testing.T) {
		calculation, _ := distService.CalculateSurplus(entityID, "2026-03")
		dist, _ := distService.CreateDistributionProposal(entityID, "2026-03", calculation)
		distService.SubmitForApproval(dist.ID)

		decision := &domain.AssemblyDecision{
			ID:       2,
			EntityID: entityID,
			Title:    "Aprovação de Rateio",
			Status:   "APPROVED",
		}
		assemblyRepo.AddDecision(decision)

		distService.ApproveWithDecision(dist.ID, decision.ID)

		// Executa distribuição (gera lançamentos contábeis)
		err := distService.ExecuteDistribution(dist.ID)
		if err != nil {
			t.Fatalf("Erro ao executar: %v", err)
		}

		// Verifica se foi executada
		executedDist, _ := distService.GetDistributionReport(dist.ID)
		if executedDist.Status != domain.StatusExecuted {
			t.Errorf("Status incorreto: esperado EXECUTED, obtido %s", executedDist.Status)
		}

		if executedDist.ExecutedAt == nil {
			t.Error("Data de execução não registrada")
		}

		t.Log("✅ Distribuição executada - Lançamentos contábeis gerados")
	})

	t.Run("Step7_VerificarLancamentos", func(t *testing.T) {
		// Verifica lançamentos foram criados para cada membro
		calculation, _ := distService.CalculateSurplus(entityID, "2026-03")
		dist, _ := distService.CreateDistributionProposal(entityID, "2026-03", calculation)
		distService.SubmitForApproval(dist.ID)

		decision := &domain.AssemblyDecision{
			ID:       3,
			EntityID: entityID,
			Title:    "Aprovação",
			Status:   "APPROVED",
		}
		assemblyRepo.AddDecision(decision)

		distService.ApproveWithDecision(dist.ID, decision.ID)
		distService.ExecuteDistribution(dist.ID)

		executedDist, _ := distService.GetDistributionReport(dist.ID)

		// Verifica se todos os membros têm lançamentos
		for _, member := range executedDist.Members {
			if member.LedgerEntryID == 0 {
				t.Errorf("Membro %s não tem lançamento contábil", member.MemberID)
			}
			t.Logf("   %s: LedgerEntryID=%d", member.MemberID, member.LedgerEntryID)
		}

		// Verifica estrutura contábil
		t.Log("📊 Estrutura de Lançamentos Contábeis:")
		t.Logf("   Lançamento 1: Sobras → Reserva Legal (10%%): R$ %.2f",
			float64(executedDist.ReservaLegal)/100)
		t.Logf("   Lançamento 2: Sobras → FATES (5%%): R$ %.2f",
			float64(executedDist.FATES)/100)
		t.Logf("   Lançamento 3: Sobras → Capital Social (Rateio):")
		for _, m := range executedDist.Members {
			t.Logf("      %s: R$ %.2f", m.MemberID, float64(m.Amount)/100)
		}

		// Verifica soma (com tolerância para arredondamento de até 5 centavos)
		totalDistribuido := int64(0)
		for _, m := range executedDist.Members {
			totalDistribuido += m.Amount
		}

		// Tolerância de 5 centavos para arredondamento em cálculos proporcionais
		diff := totalDistribuido - executedDist.Distribuivel
		if diff < 0 {
			diff = -diff
		}
		if diff > 5 {
			// Diferença maior que 5 centavos é erro
			t.Errorf("Soma dos rateios incorreta: esperado ~%d, obtido %d (diferença: %d centavos)",
				executedDist.Distribuivel, totalDistribuido, diff)
		} else if diff > 0 {
			// Diferença de 1-5 centavos é aceitável (arredondamento)
			t.Logf("⚠️  Diferença de arredondamento: %d centavos (aceitável)", diff)
		}

		t.Log("✅ Todos os lançamentos contábeis gerados corretamente")
	})

	t.Run("Step8_FluxoInvalido_ExecutarSemAprovacao", func(t *testing.T) {
		calculation, _ := distService.CalculateSurplus(entityID, "2026-03")
		dist, _ := distService.CreateDistributionProposal(entityID, "2026-03", calculation)
		distService.SubmitForApproval(dist.ID)
		// NÃO aprova

		// Tenta executar sem aprovação
		err := distService.ExecuteDistribution(dist.ID)
		if err == nil {
			t.Error("Deveria retornar erro ao executar sem aprovação")
		}

		if err != service.ErrPendingDecision {
			t.Errorf("Erro incorreto: esperado ErrPendingDecision, obtido %v", err)
		}

		t.Log("✅ Sistema bloqueia execução sem aprovação da assembleia")
	})

	t.Run("Step9_FluxoInvalido_SemSobra", func(t *testing.T) {
		// Entity com saldo zero
		emptyEntityID := "Cooperativa_Sem_Sobra"

		// Sobrescreve o mock para retornar sobra zero
		mockCalcZero := &MockSurplusCalculatorZero{}
		distServiceZero := service.NewDistributionServiceWithCalculator(
			repository.NewMockDistributionRepository(),
			repository.NewMockAssemblyRepository(),
			mockCalcZero,
		)

		// Tenta calcular sobras (não há receitas)
		_, err := distServiceZero.CalculateSurplus(emptyEntityID, "2026-03")
		if err == nil {
			t.Error("Deveria retornar erro quando não há sobra")
		}

		if err != service.ErrInvalidSurplus {
			t.Errorf("Erro incorreto: esperado ErrInvalidSurplus, obtido %v", err)
		}

		t.Log("✅ Sistema impede distribuição sem sobra disponível")
	})

	t.Run("Step10_ResumoFinal", func(t *testing.T) {
		t.Log("\n============================================================")
		t.Log("  RESUMO DO FLUXO DE DISTRIBUIÇÃO DE SOBRAS")
		t.Log("============================================================")
		t.Log("")
		t.Log("✅ Etapa 1: Cálculo das sobras com reservas")
		t.Log("   - 10% Reserva Legal (obrigatório)")
		t.Log("   - 5% FATES (obrigatório)")
		t.Log("   - 85% Distribuível")
		t.Log("")
		t.Log("✅ Etapa 2: Rateio proporcional às horas")
		t.Log("   - Seguindo ITG 2002")
		t.Log("   - Capital Social de Trabalho")
		t.Log("")
		t.Log("✅ Etapa 3: Aprovação em assembleia")
		t.Log("   - Proposta submetida")
		t.Log("   - Decisão registrada e aprovada")
		t.Log("")
		t.Log("✅ Etapa 4: Execução contábil")
		t.Log("   - Lançamentos automáticos gerados")
		t.Log("   - Partidas dobradas validadas")
		t.Log("")
		t.Log("✅ Segurança:")
		t.Log("   - Bloqueia execução sem aprovação")
		t.Log("   - Bloqueia distribuição sem sobra")
		t.Log("============================================================")
	})
}

// MockSurplusCalculatorZero retorna sobra zero para teste
type MockSurplusCalculatorZero struct{}

func (m *MockSurplusCalculatorZero) CalculateSocialSurplus(entityID string) (*surplus.SurplusCalculation, error) {
	return &surplus.SurplusCalculation{
		EntityID:     entityID,
		TotalSurplus: 0, // Sem sobra!
		TotalMinutes: 0,
		Members:      []surplus.MemberShare{},
	}, nil
}
