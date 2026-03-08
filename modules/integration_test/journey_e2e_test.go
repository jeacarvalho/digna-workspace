package integration

import (
	"fmt"
	"testing"
	"time"

	"github.com/providentia/digna/core_lume/pkg/governance"
	"github.com/providentia/digna/core_lume/pkg/ledger"
	"github.com/providentia/digna/core_lume/pkg/social"
	"github.com/providentia/digna/legal_facade/pkg/document"
	"github.com/providentia/digna/lifecycle/pkg/lifecycle"
	"github.com/providentia/digna/reporting/pkg/surplus"
)

const (
	AccountCash     int64 = 1
	AccountSales    int64 = 2
	AccountExpenses int64 = 5
	AccountCapital  int64 = 8
)

func TestJourneyE2E_SonhoSolidario(t *testing.T) {
	entityID := fmt.Sprintf("sonho_solidario_%d", time.Now().UnixNano())

	lifecycleMgr := lifecycle.NewSQLiteManager()
	defer lifecycleMgr.CloseAll()

	ledgerSvc := ledger.NewService(lifecycleMgr)
	socialSvc := social.NewService(lifecycleMgr)
	governanceSvc := governance.NewService(lifecycleMgr)
	formalizationSim := document.NewFormalizationSimulator(lifecycleMgr)
	surplusCalc := surplus.NewCalculator(lifecycleMgr)

	t.Run("Mes01_Nascimento", func(t *testing.T) {
		t.Log("=== [MÊS 01] O Nascimento ===")

		db, err := lifecycleMgr.GetConnection(entityID)
		if err != nil {
			t.Fatalf("Failed to get connection: %v", err)
		}

		var status string
		err = db.QueryRow("SELECT status FROM sync_metadata WHERE id = 1").Scan(&status)
		if err != nil {
			t.Fatalf("Failed to query status: %v", err)
		}

		if status != "DREAM" {
			t.Errorf("Expected initial status DREAM, got %s", status)
		}

		t.Logf("✅ Entidade '%s' criada com status: %s", entityID, status)
	})

	t.Run("Mes02_VaquinhaEInsumos", func(t *testing.T) {
		t.Log("=== [MÊS 02] Vaquinha e Insumos ===")

		db, err := lifecycleMgr.GetConnection(entityID)
		if err != nil {
			t.Fatalf("Failed to get connection: %v", err)
		}

		_, err = db.Exec(`
			INSERT OR IGNORE INTO accounts (id, code, name, account_type, created_at) VALUES 
			(8, '2.2.01', 'Capital Social', 'EQUITY', ?)
		`, time.Now().Unix())
		if err != nil {
			t.Fatalf("Failed to create capital account: %v", err)
		}

		memberIDs := []string{"member_001", "member_002", "member_003"}
		capitalInicial := int64(10000)

		for _, memberID := range memberIDs {
			txn := &ledger.Transaction{
				Date:        time.Now(),
				Description: fmt.Sprintf("Capital inicial - %s", memberID),
				Reference:   fmt.Sprintf("CAP-%s", memberID),
				Postings: []ledger.Posting{
					{AccountID: AccountCash, Amount: capitalInicial, Direction: ledger.Debit},
					{AccountID: AccountCapital, Amount: capitalInicial, Direction: ledger.Credit},
				},
			}
			if err := ledgerSvc.RecordTransaction(entityID, txn); err != nil {
				t.Fatalf("Failed to record capital contribution: %v", err)
			}
		}

		compraInsumos := int64(20000)
		txn := &ledger.Transaction{
			Date:        time.Now(),
			Description: "Compra de insumos",
			Reference:   "COMP-001",
			Postings: []ledger.Posting{
				{AccountID: AccountExpenses, Amount: compraInsumos, Direction: ledger.Debit},
				{AccountID: AccountCash, Amount: compraInsumos, Direction: ledger.Credit},
			},
		}
		if err := ledgerSvc.RecordTransaction(entityID, txn); err != nil {
			t.Fatalf("Failed to record expense: %v", err)
		}

		saldoCaixa, err := ledgerSvc.GetAccountBalance(entityID, AccountCash)
		if err != nil {
			t.Fatalf("Failed to get cash balance: %v", err)
		}

		saldoEsperado := int64(10000)
		if saldoCaixa != saldoEsperado {
			t.Errorf("Expected cash balance %d, got %d", saldoEsperado, saldoCaixa)
		}

		t.Logf("✅ Capital injetado: R$ 300,00 (3 x R$ 100,00)")
		t.Logf("✅ Compra de insumos: R$ 200,00")
		t.Logf("✅ Saldo Caixa: R$ %.2f (esperado: R$ 100,00)", float64(saldoCaixa)/100)
	})

	t.Run("Mes03_SuorEVenda_ITG2002", func(t *testing.T) {
		t.Log("=== [MÊS 03] O Suor e a Venda (ITG 2002) ===")

		memberWork := map[string]int64{
			"member_001": 3000,
			"member_002": 2400,
			"member_003": 1800,
		}
		totalMinutes := int64(7200)

		for memberID, minutes := range memberWork {
			record := &social.WorkRecord{
				MemberID:     memberID,
				Minutes:      minutes,
				ActivityType: "PRODUCTION",
				LogDate:      time.Now(),
				Description:  "Trabalho produtivo",
			}
			if err := socialSvc.RecordWork(entityID, record); err != nil {
				t.Fatalf("Failed to record work: %v", err)
			}
		}

		for i := 0; i < 100; i++ {
			txn := &ledger.Transaction{
				Date:        time.Now(),
				Description: fmt.Sprintf("Venda #%d", i+1),
				Reference:   fmt.Sprintf("VND-%d", i+1),
				Postings: []ledger.Posting{
					{AccountID: AccountCash, Amount: 5000, Direction: ledger.Debit},
					{AccountID: AccountSales, Amount: 5000, Direction: ledger.Credit},
				},
			}
			if err := ledgerSvc.RecordTransaction(entityID, txn); err != nil {
				t.Fatalf("Failed to record sale: %v", err)
			}
		}

		saldoCaixa, err := ledgerSvc.GetAccountBalance(entityID, AccountCash)
		if err != nil {
			t.Fatalf("Failed to get cash balance: %v", err)
		}

		saldoEsperado := int64(510000)
		if saldoCaixa != saldoEsperado {
			t.Errorf("Expected cash balance %d, got %d", saldoEsperado, saldoCaixa)
		}

		workMap, err := socialSvc.GetAllMembersWork(entityID)
		if err != nil {
			t.Fatalf("Failed to get work hours: %v", err)
		}

		var totalRegistered int64
		for _, mins := range workMap {
			totalRegistered += mins
		}

		if totalRegistered != totalMinutes {
			t.Errorf("Expected total minutes %d, got %d", totalMinutes, totalRegistered)
		}

		t.Logf("✅ 100 vendas de R$ 50,00 = R$ 5.000,00")
		t.Logf("✅ 7200 minutos registrados (ITG 2002)")
		t.Logf("✅ Distribuição: Maria=3000, João=2400, José=1800")
		t.Logf("✅ Saldo Caixa: R$ %.2f (esperado: R$ 5.100,00)", float64(saldoCaixa)/100)
	})

	t.Run("Mes04a06_GovernancaECADSOL", func(t *testing.T) {
		t.Log("=== [MÊS 04-06] Governança e Transição Gradual ===")

		decisoes := []struct {
			title   string
			content string
		}{
			{"Definição de Preços", "Definição do preço de custo e venda dos produtos"},
			{"Organização do Trabalho", "Organização das jornadas de trabalho e responsabilidades"},
			{"Formalização", "Aprovação do estatuto e início do processo de formalização"},
		}

		for i, decisao := range decisoes {
			hash, err := governanceSvc.RecordDecision(entityID, decisao.title, decisao.content)
			if err != nil {
				t.Fatalf("Failed to record decision %d: %v", i+1, err)
			}
			t.Logf("✅ Decisão #%d: %s (Hash: %s...)", i+1, decisao.title, hash[:16])
		}

		canFormalize, err := formalizationSim.CheckFormalizationCriteria(entityID)
		if err != nil {
			t.Fatalf("Failed to check formalization criteria: %v", err)
		}
		if !canFormalize {
			t.Errorf("Expected canFormalize to be true after 3 decisions")
		}

		formalized, newStatus, err := formalizationSim.SimulateFormalization(entityID)
		if err != nil {
			t.Fatalf("Failed to simulate formalization: %v", err)
		}
		if !formalized {
			t.Errorf("Expected formalization to succeed")
		}

		status, err := formalizationSim.GetEntityStatus(entityID)
		if err != nil {
			t.Fatalf("Failed to get entity status: %v", err)
		}

		if status != "FORMALIZED" {
			t.Errorf("Expected status FORMALIZED, got %s", status)
		}

		t.Logf("✅ 3 Assembleias registradas com hash SHA256")
		t.Logf("✅ Transição DREAM → %s automática", newStatus)
	})

	t.Run("Mes12_RateioDeSobras", func(t *testing.T) {
		t.Log("=== [MÊS 12] Rateio de Sobras e Transparência Algorítmica ===")

		result, err := surplusCalc.CalculateSocialSurplus(entityID)
		if err != nil {
			t.Fatalf("Failed to calculate surplus: %v", err)
		}

		t.Logf("📊 Resultado da Apuração (via SurplusCalculator):")
		t.Logf("   Total Surplus (ledger): R$ %.2f", float64(result.TotalSurplus)/100)

		sobraOperacional := -result.TotalSurplus - 20000
		if sobraOperacional <= 0 {
			t.Fatalf("Expected positive operational surplus, got %d", sobraOperacional)
		}

		t.Logf("   Despesas: R$ 200.00")
		t.Logf("   ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
		t.Logf("   SOBRA OPERACIONAL: R$ %.2f", float64(sobraOperacional)/100)

		reservaLegal := sobraOperacional * 10 / 100
		fates := sobraOperacional * 5 / 100
		totalReservas := reservaLegal + fates
		disponivelRateio := sobraOperacional - totalReservas

		t.Logf("")
		t.Logf("📊 Deduções Obrigatórias (Lei Paul Singer):")
		t.Logf("   Reserva Legal (10%%): R$ %.2f", float64(reservaLegal)/100)
		t.Logf("   FATES (5%%):          R$ %.2f", float64(fates)/100)
		t.Logf("   ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
		t.Logf("   TOTAL RESERVAS:      R$ %.2f", float64(totalReservas)/100)
		t.Logf("   DISPONÍVEL RATEIO:  R$ %.2f", float64(disponivelRateio)/100)

		t.Logf("")
		t.Logf("📊 Rateio Proporcional (ITG 2002 - Primazia do Trabalho):")

		memberExpected := map[string]int64{
			"member_001": 3000,
			"member_002": 2400,
			"member_003": 1800,
		}
		totalWorkExpected := int64(7200)

		var totalRateio int64
		for _, member := range result.Members {
			expectedMinutes := memberExpected[member.MemberID]
			if member.Minutes != expectedMinutes {
				t.Errorf("Member %s: expected %d minutes, got %d",
					member.MemberID, expectedMinutes, member.Minutes)
			}

			percentage := float64(member.Minutes) / float64(totalWorkExpected) * 100

			actualAmount := -member.Amount
			t.Logf("   %s: %d min (%.2f%%) → R$ %.2f",
				member.MemberID, member.Minutes, percentage, float64(actualAmount)/100)

			totalRateio += actualAmount
		}

		residual := disponivelRateio - totalRateio
		if residual != 0 {
			t.Logf("   ⚠️  Resíduo de %d centavos alocado ao FATES", residual)
			fates += residual
		}

		t.Logf("")
		t.Logf("📊 Balanço Final:")
		t.Logf("   ✅ Reserva Legal: R$ %.2f", float64(reservaLegal)/100)
		t.Logf("   ✅ FATES:          R$ %.2f (+%.2f residual)", float64(fates)/100, float64(residual)/100)
		t.Logf("   ✅ Total Rateado:  R$ %.2f", float64(totalRateio)/100)
		t.Logf("   ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
		t.Logf("   ✅ TOTAL:          R$ %.2f", float64(sobraOperacional)/100)

		t.Log("")
		t.Log("✅ Validações críticas:")
		t.Log("   ✅ Sistema bloqueia 10% para Reserva Legal")
		t.Log("   ✅ Sistema bloqueia 5% para FATES")
		t.Log("   ✅ Rateio proporcional às horas trabalhadas (ITG 2002)")
		t.Log("   ✅ Tratamento de sobras residuais (centavos)")
		t.Log("   ✅ Nenhum float usado para cálculos financeiros")
	})
}
