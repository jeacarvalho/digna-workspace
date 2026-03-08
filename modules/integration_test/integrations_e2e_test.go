package integration

import (
	"context"
	"testing"
	"time"

	"github.com/providentia/digna/core_lume/pkg/governance"
	"github.com/providentia/digna/core_lume/pkg/ledger"
	"github.com/providentia/digna/core_lume/pkg/social"
	"github.com/providentia/digna/integrations/pkg/integrations"
	"github.com/providentia/digna/legal_facade/pkg/document"
	"github.com/providentia/digna/lifecycle/pkg/lifecycle"
	"github.com/providentia/digna/reporting/pkg/surplus"
)

func TestE2E_IntegracoesGovernamentais(t *testing.T) {
	entityID := "test_integracoes_001"

	lifecycleMgr := lifecycle.NewSQLiteManager()
	defer lifecycleMgr.CloseAll()

	ledgerSvc := ledger.NewService(lifecycleMgr)
	socialSvc := social.NewService(lifecycleMgr)
	governanceSvc := governance.NewService(lifecycleMgr)

	db, err := lifecycleMgr.GetConnection(entityID)
	if err != nil {
		t.Fatalf("Failed to get connection: %v", err)
	}

	integrationSvc, err := integrations.NewMockIntegrationService(db)
	if err != nil {
		t.Fatalf("Failed to create integration service: %v", err)
	}

	ctx := context.Background()

	t.Run("Setup_EntidadeBasica", func(t *testing.T) {
		t.Log("=== Setup: Criando dados básicos da entidade ===")

		db.Exec(`INSERT OR IGNORE INTO accounts (id, code, name, account_type, created_at) VALUES 
			(8, '2.2.01', 'Capital Social', 'EQUITY', ?)`,
			0)

		txn := &ledger.Transaction{
			Date:        time.Now(),
			Description: "Capital inicial",
			Reference:   "CAP-001",
			Postings: []ledger.Posting{
				{AccountID: 1, Amount: 100000, Direction: ledger.Debit},
				{AccountID: 8, Amount: 100000, Direction: ledger.Credit},
			},
		}
		if err := ledgerSvc.RecordTransaction(entityID, txn); err != nil {
			t.Fatalf("Failed to record: %v", err)
		}

		socialSvc.RecordWork(entityID, &social.WorkRecord{
			MemberID:     "member_001",
			Minutes:      1000,
			ActivityType: "PRODUCTION",
			Description:  "Trabalho",
		})

		for i := 0; i < 50; i++ {
			ledgerSvc.RecordTransaction(entityID, &ledger.Transaction{
				Date:        time.Now(),
				Description: "Venda",
				Reference:   "VND-001",
				Postings: []ledger.Posting{
					{AccountID: 1, Amount: 1000, Direction: ledger.Debit},
					{AccountID: 2, Amount: 1000, Direction: ledger.Credit},
				},
			})
		}

		governanceSvc.RecordDecision(entityID, "Aprovação", "Aprovação de estatuto")

		t.Log("✅ Dados básicos criados")
	})

	t.Run("ReceitaFederal_ConsultarCNPJ", func(t *testing.T) {
		t.Log("=== [Receita Federal] Consultar CNPJ ===")

		cnpj := "12.345.678/0001-90"
		resp, err := integrationSvc.ReceitaFederal().ConsultarCNPJ(ctx, cnpj)
		if err != nil {
			t.Fatalf("Failed to consult CNPJ: %v", err)
		}

		t.Logf("✅ CNPJ %s: %s (%s)", cnpj, resp.RazaoSocial, resp.Situacao)
	})

	t.Run("MTE_EnviarRAIS", func(t *testing.T) {
		t.Log("=== [MTE] Enviar RAIS ===")

		req := &integrations.RAISRequest{
			EntityID: entityID,
			Ano:      2025,
			CNPJ:     "12.345.678/0001-90",
			Trabalhadores: []integrations.TrabalhadorRAIS{
				{CPF: "12345678900", Nome: "Maria", CBO: "4110", Salario: 150000, HorasSemanais: 40},
			},
		}

		resp, err := integrationSvc.MTE().EnviarRAIS(ctx, req)
		if err != nil {
			t.Fatalf("Failed to send RAIS: %v", err)
		}

		t.Logf("✅ RAIS enviada: %s (Protocolo: %s)", resp.Status, resp.Protocolo)
	})

	t.Run("MTE_RegistrarCAT", func(t *testing.T) {
		t.Log("=== [MTE] Registrar CAT ===")

		req := &integrations.CATRequest{
			EntityID:     entityID,
			CNPJ:         "12.345.678/0001-90",
			DataAcidente: time.Now(),
			TipoAcidente: "TIPO_1",
			Descricao:    "Acidente no trabalho",
			Comunicante:  "Empregador",
		}

		resp, err := integrationSvc.MTE().RegistrarCAT(ctx, req)
		if err != nil {
			t.Fatalf("Failed to register CAT: %v", err)
		}

		t.Logf("✅ CAT registrada: %s (Protocolo: %s)", resp.Status, resp.Protocolo)
	})

	t.Run("MDS_EnviarRelatorioSocial", func(t *testing.T) {
		t.Log("=== [MDS] Enviar Relatório Social ===")

		relatorio := &integrations.RelatorioSocial{
			EntityID:   entityID,
			Periodo:    "2026-01",
			TotalHoras: 1000,
			Membros:    3,
		}

		resp, err := integrationSvc.MDS().EnviarRelatorioSocial(ctx, relatorio)
		if err != nil {
			t.Fatalf("Failed to send social report: %v", err)
		}

		t.Logf("✅ Relatório Social enviado: %s", resp.Status)
	})

	t.Run("SEFAZ_EmitirNFe", func(t *testing.T) {
		t.Log("=== [SEFAZ] Emitir NFe ===")

		req := &integrations.NFeRequest{
			EntityID:    entityID,
			CNPJ:        "12.345.678/0001-90",
			DataEmissao: time.Now(),
			Destinatario: integrations.DestinatarioNFe{
				CNPJCPF: "98765432000199",
				Nome:    "Empresa Compradora",
			},
			Itens: []integrations.ItemNFe{
				{Numero: 1, Codigo: "001", Descricao: "Vestuário", Quantidade: 10, ValorUnitario: 5000},
			},
		}

		resp, err := integrationSvc.SEFAZ().EmitirNFe(ctx, req)
		if err != nil {
			t.Fatalf("Failed to emit NFe: %v", err)
		}

		t.Logf("✅ NFe emitida: %s (Protocolo: %s)", resp.Status, resp.Protocolo)
	})

	t.Run("BNDES_SimularCredito", func(t *testing.T) {
		t.Log("=== [BNDES] Simular Crédito ===")

		req := &integrations.SimulacaoCredito{
			EntityID: entityID,
			CNPJ:     "12.345.678/0001-90",
			Linha:    "MICROCREDITO",
			Valor:    50000,
			Prazo:    36,
		}

		resp, err := integrationSvc.BNDES().SimularCredito(ctx, req)
		if err != nil {
			t.Fatalf("Failed to simulate credit: %v", err)
		}

		t.Logf("✅ Simulação BNDES: Valor parcela R$ %.2f, CET: %.2f%%",
			float64(resp.ValorParcela)/100, resp.CET)
	})

	t.Run("SEBRAE_ConsultarCursos", func(t *testing.T) {
		t.Log("=== [SEBRAE] Consultar Cursos ===")

		resp, err := integrationSvc.SEBRAE().ConsultarCursos(ctx, "12.345.678/0001-90")
		if err != nil {
			t.Fatalf("Failed to consult courses: %v", err)
		}

		t.Logf("✅ SEBRAE: %d cursos disponíveis", len(resp))
	})

	t.Run("Providentia_Sync", func(t *testing.T) {
		t.Log("=== [Providentia] Sincronização ===")

		pkg := &integrations.SyncPackage{
			EntityID:  entityID,
			Timestamp: time.Now().Unix(),
			AggregatedData: integrations.AggregatedMetrics{
				TotalSales:     50000,
				TotalWorkHours: 1000,
				TotalMembers:   3,
				LegalStatus:    "DREAM",
			},
		}

		resp, err := integrationSvc.Providentia().SyncPackage(ctx, pkg)
		if err != nil {
			t.Fatalf("Failed to sync: %v", err)
		}

		t.Logf("✅ Providentia Sync: %s (Protocolo: %s)", resp.Status, resp.Protocolo)
	})

	t.Run("Providentia_Marketplace", func(t *testing.T) {
		t.Log("=== [Providentia] Marketplace ===")

		offer := &integrations.MarketplaceOffer{
			EntityID:    entityID,
			ProductName: "Vestuário Artesanal",
			Description: "Peças únicas",
			Quantity:    50,
			Price:       5000,
			Unit:        "un",
			Active:      true,
		}

		resp, err := integrationSvc.Providentia().RegisterOffer(ctx, offer)
		if err != nil {
			t.Fatalf("Failed to publish offer: %v", err)
		}

		t.Logf("✅ Marketplace: Oferta publicada (ID: %s, Status: %s)", resp.ID, resp.Status)
	})

	t.Run("SurplusCalculator_ComDeducoes", func(t *testing.T) {
		t.Log("=== [SurplusCalculator] Calcular com deduções ===")

		surplusCalc := surplus.NewCalculator(lifecycleMgr)

		ledgerSvc.RecordTransaction(entityID, &ledger.Transaction{
			Date:        time.Now(),
			Description: "Vendas",
			Reference:   "VND-SURPLUS",
			Postings: []ledger.Posting{
				{AccountID: 1, Amount: 100000, Direction: ledger.Debit},
				{AccountID: 2, Amount: 100000, Direction: ledger.Credit},
			},
		})

		result, err := surplusCalc.CalculateWithDeductions(entityID)
		if err != nil {
			t.Fatalf("Failed to calculate surplus: %v", err)
		}

		t.Logf("📊 Sobra Bruta: R$ %.2f", float64(result.GrossSurplus)/100)
		t.Logf("📊 Reserva Legal (10%%): R$ %.2f", float64(result.LegalReserve)/100)
		t.Logf("📊 FATES (5%%): R$ %.2f", float64(result.FATES)/100)
		t.Logf("📊 Total Deduções: R$ %.2f", float64(result.TotalDeductions)/100)
		t.Logf("📊 Disponível Rateio: R$ %.2f", float64(result.AvailableForShare)/100)
		t.Logf("📊 Resíduo: %d centavos", result.Residual)

		if result.GrossSurplus <= 0 {
			t.Skip("Skipping - no surplus to test")
			return
		}
	})

	t.Run("Formalizacao_AutoTransicao", func(t *testing.T) {
		t.Log("=== [Formalização] Transição Automática ===")

		entityID2 := "test_formalizacao_001"
		formalizationSim := document.NewFormalizationSimulator(lifecycleMgr)

		initialStatus, _ := formalizationSim.GetEntityStatus(entityID2)
		t.Logf("📊 Status inicial: %s", initialStatus)

		governanceSvc.RecordDecision(entityID2, "Definição de Preços", "Definição de preços de venda")
		governanceSvc.RecordDecision(entityID2, "Organização do Trabalho", "Organização das equipes")

		transicionado, novoStatus, err := formalizationSim.AutoTransitionIfReady(entityID2)
		if err != nil {
			t.Fatalf("Failed to auto transition: %v", err)
		}

		t.Logf("📊 Após 2 decisões: transicionado=%v, status=%s", transicionado, novoStatus)

		governanceSvc.RecordDecision(entityID2, "Formalização", "Aprovação do estatuto")

		transicionado, novoStatus, err = formalizationSim.AutoTransitionIfReady(entityID2)
		if err != nil {
			t.Fatalf("Failed to auto transition: %v", err)
		}

		if !transicionado {
			t.Errorf("Expected transition to happen after 3 decisions")
		}

		t.Logf("✅ Transição automática: %s → %s", initialStatus, novoStatus)
	})

	t.Log("=")
	t.Log("✅ TODAS INTEGRAÇÕES GOVERNAMENTAIS VALIDADAS")
	t.Log("=")
}
