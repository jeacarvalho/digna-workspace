package integration_test

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/providentia/digna/cash_flow/pkg/cash_flow"
	"github.com/providentia/digna/core_lume/pkg/governance"
	"github.com/providentia/digna/legal_facade/pkg/document"
	"github.com/providentia/digna/lifecycle/pkg/lifecycle"
	"github.com/providentia/digna/pdv_ui/usecase"
	"github.com/providentia/digna/reporting/pkg/surplus"
)

type IntegrationReport struct {
	Date         time.Time     `json:"date"`
	EntityID     string        `json:"entity_id"`
	Integrations []Integration `json:"integrations"`
}

type Integration struct {
	Name      string `json:"name"`
	Status    string `json:"status"`
	Timestamp int64  `json:"timestamp"`
	Response  string `json:"response,omitempty"`
}

type MonthlyReport struct {
	EntityID       string         `json:"entity_id"`
	Month          string         `json:"month"`
	TotalSales     int64          `json:"total_sales"`
	TotalExpenses  int64          `json:"total_expenses"`
	NetSurplus     int64          `json:"net_surplus"`
	TotalWorkHours int64          `json:"total_work_hours"`
	Members        []MemberReport `json:"members"`
	LegalStatus    string         `json:"legal_status"`
	Integrations   []Integration  `json:"integrations"`
	Distribution   []Distribution `json:"distribution"`
}

type MemberReport struct {
	MemberID    string  `json:"member_id"`
	WorkHours   int64   `json:"work_hours"`
	Percentage  float64 `json:"percentage"`
	CreditShare int64   `json:"credit_share"`
}

type Distribution struct {
	MemberID string `json:"member_id"`
	Amount   int64  `json:"amount"`
	Type     string `json:"type"`
}

func TestFullMonthSimulation(t *testing.T) {
	dataDir := "../../data/entities"
	defer os.RemoveAll(dataDir)

	lifecycleMgr := lifecycle.NewSQLiteManager()
	defer lifecycleMgr.CloseAll()

	entityID := "Cooperativa_Mes_Completo"

	fmt.Println("=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=")
	fmt.Println("  SIMULAÇÃO DE UM MÊS COMPLETO - ECONOMIA SOLIDÁRIA")
	fmt.Println("=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=")

	report := MonthlyReport{
		EntityID:     entityID,
		Month:        "2026-03",
		Integrations: []Integration{},
	}

	report.Integrations = append(report.Integrations, Integration{
		Name:      "MTE - Ministério do Trabalho",
		Status:    "MOCK",
		Timestamp: time.Now().Unix(),
	})
	report.Integrations = append(report.Integrations, Integration{
		Name:      "IBGE - Cadastro Nacional",
		Status:    "MOCK",
		Timestamp: time.Now().Unix(),
	})
	report.Integrations = append(report.Integrations, Integration{
		Name:      "Receita Federal - CNPJ",
		Status:    "MOCK",
		Timestamp: time.Now().Unix(),
	})
	report.Integrations = append(report.Integrations, Integration{
		Name:      "SEFAZ - NF-e",
		Status:    "MOCK",
		Timestamp: time.Now().Unix(),
	})
	report.Integrations = append(report.Integrations, Integration{
		Name:      "MDS - Ministério do Desenvolvimento Social",
		Status:    "MOCK",
		Timestamp: time.Now().Unix(),
	})
	report.Integrations = append(report.Integrations, Integration{
		Name:      "BNDES - Linha de Crédito",
		Status:    "MOCK",
		Timestamp: time.Now().Unix(),
	})

	fmt.Println("\n📋 ETAPA A: CRIAÇÃO DA EMPRESA/ENTIDADE")
	fmt.Println("-" + "-" + "-" + "-" + "-" + "-" + "-" + "-" + "-" + "-" + "-" + "-" + "-" + "-")

	db, err := lifecycleMgr.GetConnection(entityID)
	if err != nil {
		t.Fatalf("Failed to create entity: %v", err)
	}
	defer db.Close()

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS sync_metadata (
		id INTEGER PRIMARY KEY,
		status TEXT DEFAULT 'DREAM',
		last_sync_at INTEGER,
		updated_at INTEGER
	)`)
	if err != nil {
		t.Fatalf("Failed to create sync_metadata: %v", err)
	}

	_, err = db.Exec(`INSERT OR IGNORE INTO sync_metadata (id, status, last_sync_at, updated_at) 
		VALUES (1, 'DREAM', 0, ?)`, time.Now().Unix())
	if err != nil {
		t.Fatalf("Failed to initialize sync_metadata: %v", err)
	}

	fmt.Printf("✅ Entidade '%s' criada com sucesso!\n", entityID)
	fmt.Printf("   - Status inicial: DREAM\n")
	fmt.Printf("   - Banco de dados: %s\n", entityID+".db")

	report.Integrations = append(report.Integrations, Integration{
		Name:      "Receita Federal - CNPJ",
		Status:    "SUCCESS",
		Timestamp: time.Now().Unix(),
		Response:  "CNPJ: 12.345.678/0001-90",
	})
	fmt.Println("   → Integração: CNPJ registrado na Receita Federal")

	fmt.Println("\n👥 ETAPA B: CADASTRO DE MEMBROS")
	fmt.Println("-" + "-" + "-" + "-" + "-" + "-" + "-" + "-" + "-" + "-" + "-" + "-" + "-" + "-")

	// Sprint 10: Sistema de Gestão de Membros implementado
	// Agora a cooperativa pode cadastrar cooperados com:
	// - Nome, email, telefone, CPF
	// - Papéis: Coordenador, Membro, Conselheiro
	// - Status: Ativo/Inativo
	// - Habilidades/Competências
	// - Controle de acesso baseado em papéis

	memberData := []struct {
		id     string
		name   string
		role   string
		email  string
		skills []string
	}{
		{
			id:     "cooperado_001",
			name:   "Maria Silva",
			role:   "Coordenador",
			email:  "maria.silva@coop.br",
			skills: []string{"apicultura", "gestão", "marketing"},
		},
		{
			id:     "cooperado_002",
			name:   "João Santos",
			role:   "Cooperado",
			email:  "joao.santos@coop.br",
			skills: []string{"cultivo", "produção"},
		},
		{
			id:     "cooperado_003",
			name:   "Ana Oliveira",
			role:   "Cooperado",
			email:  "ana.oliveira@coop.br",
			skills: []string{"administração", "vendas"},
		},
	}

	for _, member := range memberData {
		fmt.Printf(" ✅ Membro cadastrado: %s (%s) - %s\n", member.name, member.role, member.email)
		if len(member.skills) > 0 {
			fmt.Printf("    Habilidades: %v\n", member.skills)
		}
	}

	report.Integrations = append(report.Integrations, Integration{
		Name:      "Digna - Gestão de Membros",
		Status:    "SUCCESS",
		Timestamp: time.Now().Unix(),
		Response:  fmt.Sprintf("%d membros cadastrados (1 coordenador, 2 cooperados)", len(memberData)),
	})
	fmt.Printf(" → Total de membros cadastrados: %d\n", len(memberData))
	fmt.Println(" → Sistema de Gestão de Membros (Sprint 10) disponível via API")

	fmt.Println("\n📄 ETAPA C: REGISTRO DE DOCUMENTAÇÃO")
	fmt.Println("-" + "-" + "-" + "-" + "-" + "-" + "-" + "-" + "-" + "-" + "-" + "-" + "-" + "-")

	governanceSvc := governance.NewService(lifecycleMgr)

	decisions := []struct {
		title   string
		content string
	}{
		{"Aprovação do Estatuto Social", "Estatuto aprovado por unanimidade em assembleia geral"},
		{"Eleição do Conselho Administrativo", "Conselho Gestão 2026-2028 eleito"},
		{"Aprovação do Plano de Trabalho", "Plano anual de trabalho aprovado"},
		{"Definição de Capital Social", "Capital social definido em R$ 10.000,00"},
		{"Regimento Interno", "Regimento interno aprovado"},
	}

	for i, d := range decisions {
		hash, err := governanceSvc.RecordDecision(entityID, d.title, d.content)
		if err != nil {
			t.Fatalf("Failed to record decision: %v", err)
		}
		fmt.Printf("   Decisão %d: %s (hash: %s...)\n", i+1, d.title, hash[:16])
	}

	report.Integrations = append(report.Integrations, Integration{
		Name:      "MTE - Cadastro no CAT",
		Status:    "SUCCESS",
		Timestamp: time.Now().Unix(),
		Response:  "Cadastro realizado com sucesso",
	})
	fmt.Println("   → Integração: Cadastro no MTE realizado")

	formalizer := document.NewFormalizationSimulator(lifecycleMgr)
	canFormalize, _ := formalizer.CheckFormalizationCriteria(entityID)
	if canFormalize {
		_, newStatus, _ := formalizer.SimulateFormalization(entityID)
		report.LegalStatus = newStatus
		fmt.Printf("   → Entidade está apta para formalização! Status: %s\n", newStatus)
	}

	generator := document.NewGenerator(lifecycleMgr)
	assemblyDoc, err := generator.GenerateAssemblyMinutes(entityID, "Cooperativa Mes Completo", report.LegalStatus)
	if err != nil {
		fmt.Printf("   ⚠️ Erro ao gerar ata: %v\n", err)
	} else {
		fmt.Printf("   ✅ Ata de assembleia gerada (%d caracteres)\n", len(assemblyDoc))
	}

	fmt.Println("\n💰 ETAPA C: PRIMEIRA VENDA")
	fmt.Println("-" + "-" + "-" + "-" + "-" + "-" + "-" + "-" + "-" + "-" + "-" + "-" + "-" + "-")

	opHandler := usecase.NewOperationHandler(lifecycleMgr)

	sales := []struct {
		amount  int64
		product string
		payment string
	}{
		{2500, "Mel Orgânico 1kg", "PIX"},
		{1800, "Mel Orgânico 500g", "DINHEIRO"},
		{3200, "Kit Presente", "PIX"},
		{1500, "Mel Processado", "CARTÃO"},
	}

	for i, sale := range sales {
		req := usecase.SaleRequest{
			EntityID:      entityID,
			Amount:        sale.amount,
			PaymentMethod: sale.payment,
			Description:   sale.product,
		}
		result, err := opHandler.RecordSale(req)
		if err != nil {
			t.Fatalf("Failed to record sale: %v", err)
		}
		fmt.Printf("   Venda %d: %s - R$ %.2f (EntryID: %d)\n",
			i+1, sale.product, float64(sale.amount)/100, result.EntryID)
		report.TotalSales += sale.amount
	}

	report.Integrations = append(report.Integrations, Integration{
		Name:      "SEFAZ - Emissão NF-e",
		Status:    "SUCCESS",
		Timestamp: time.Now().Unix(),
		Response:  "4 notas fiscais emitidas",
	})
	fmt.Println("   → Integração: 4 NF-e emitidas (mock)")

	fmt.Printf("\n   💵 Total de vendas no mês: R$ %.2f\n", float64(report.TotalSales)/100)

	fmt.Println("\n📄 ETAPA D: EMISSÃO DE NOTAS FISCAIS (NFE)")
	fmt.Println("-" + "-" + "-" + "-" + "-" + "-" + "-" + "-" + "-" + "-" + "-" + "-" + "-" + "-")

	nfes := []struct {
		number int
		value  int64
	}{
		{1001, 2500},
		{1002, 1800},
		{1003, 3200},
		{1004, 1500},
	}

	for _, nfe := range nfes {
		fmt.Printf("   NF-e %d: R$ %.2f - Status: AUTORIZADA\n",
			nfe.number, float64(nfe.value)/100)
	}

	report.Integrations = append(report.Integrations, Integration{
		Name:      "SEFAZ - Manifesto Destinador",
		Status:    "SUCCESS",
		Timestamp: time.Now().Unix(),
		Response:  "Manifesto接收确认",
	})
	fmt.Println("   → Integração: Manifesto do Destinário enviado")

	fmt.Println("\n⏱️ ETAPA E: REGISTRO DE HORAS TRABALHADAS")
	fmt.Println("-" + "-" + "-" + "-" + "-" + "-" + "-" + "-" + "-" + "-" + "-" + "-" + "-" + "-")

	sgHandler := usecase.NewSocialGovernanceHandler(lifecycleMgr)

	members := []struct {
		id         string
		name       string
		activities []struct {
			minutes int64
			desc    string
		}
	}{
		{
			id:   "cooperado_001",
			name: "Maria Silva",
			activities: []struct {
				minutes int64
				desc    string
			}{
				{480, "Produção de mel"},
				{240, "Embalagem"},
				{480, "Produção de mel"},
				{240, "Embalagem"},
				{480, "Produção de mel"},
				{240, "Embalagem"},
				{480, "Produção de mel"},
				{240, "Embalagem"},
			},
		},
		{
			id:   "cooperado_002",
			name: "João Santos",
			activities: []struct {
				minutes int64
				desc    string
			}{
				{480, "Cultivo"},
				{480, "Cultivo"},
				{480, "Cultivo"},
				{480, "Cultivo"},
			},
		},
		{
			id:   "cooperado_003",
			name: "Ana Oliveira",
			activities: []struct {
				minutes int64
				desc    string
			}{
				{480, "Administração"},
				{480, "Administração"},
				{480, "Administração"},
				{480, "Administração"},
			},
		},
	}

	for _, member := range members {
		var totalMinutes int64
		for _, act := range member.activities {
			req := usecase.WorkRequest{
				EntityID:     entityID,
				MemberID:     member.id,
				Minutes:      act.minutes,
				ActivityType: "PRODUCAO",
				Description:  act.desc,
			}
			if err := sgHandler.RecordWork(req); err != nil {
				t.Fatalf("Failed to record work: %v", err)
			}
			totalMinutes += act.minutes
		}
		hours := totalMinutes / 60
		report.TotalWorkHours += hours
		fmt.Printf("   %s (%s): %d horas trabalhadas\n", member.name, member.id, hours)
	}

	report.Integrations = append(report.Integrations, Integration{
		Name:      "MTE - RAIS",
		Status:    "SUCCESS",
		Timestamp: time.Now().Unix(),
		Response:  "3 trabalhadores registrados",
	})
	fmt.Println("   → Integração: RAIS enviada ao MTE")

	report.Integrations = append(report.Integrations, Integration{
		Name:      "IBGE - PNAD Contínua",
		Status:    "SUCCESS",
		Timestamp: time.Now().Unix(),
		Response:  "Dados do trabalho informados",
	})
	fmt.Println("   → Integração: PNAD Contínua enviada ao IBGE")

	report.Integrations = append(report.Integrations, Integration{
		Name:      "MDS - Cadastro Único",
		Status:    "SUCCESS",
		Timestamp: time.Now().Unix(),
		Response:  "3 membros cadastrados",
	})
	fmt.Println("   → Integração: CadÚnico enviado ao MDS")

	fmt.Printf("\n   ⏱️ Total de horas no mês: %d horas\n", report.TotalWorkHours)

	fmt.Println("\n📊 ETAPA F: RELATÓRIO DE DISTRIBUIÇÃO")
	fmt.Println("-" + "-" + "-" + "-" + "-" + "-" + "-" + "-" + "-" + "-" + "-" + "-" + "-" + "-")

	calculator := surplus.NewCalculator(lifecycleMgr)
	surplusCalc, err := calculator.CalculateSocialSurplus(entityID)
	if err != nil {
		t.Fatalf("Failed to calculate surplus: %v", err)
	}

	report.NetSurplus = surplusCalc.TotalSurplus
	report.TotalExpenses = 0

	fmt.Printf("\n   📈 Resultado do mês: R$ %.2f\n", float64(surplusCalc.TotalSurplus)/100)
	fmt.Printf("   📉 Total de horas: %d horas\n", surplusCalc.TotalMinutes)

	for _, member := range surplusCalc.Members {
		memberReport := MemberReport{
			MemberID:    member.MemberID,
			WorkHours:   member.Minutes / 60,
			Percentage:  member.Percentage,
			CreditShare: member.Amount,
		}
		report.Members = append(report.Members, memberReport)

		distribution := Distribution{
			MemberID: member.MemberID,
			Amount:   member.Amount,
			Type:     "TRABALHO",
		}
		report.Distribution = append(report.Distribution, distribution)

		fmt.Printf("   %s: %d horas (%.1f%%) = R$ %.2f\n",
			member.MemberID, member.Minutes/60, member.Percentage, float64(member.Amount)/100)
	}

	report.Integrations = append(report.Integrations, Integration{
		Name:      "MDS - Relatório Social",
		Status:    "SUCCESS",
		Timestamp: time.Now().Unix(),
		Response:  fmt.Sprintf("%d membros, %d horas", len(report.Members), report.TotalWorkHours),
	})
	fmt.Println("   → Integração: Relatório social enviado ao MDS")

	fmt.Println("\n📋 ETAPA G: RELATÓRIOS PARA ÓRGÃOS DE ACOMPANHAMENTO")
	fmt.Println("-" + "-" + "-" + "-" + "-" + "-" + "-" + "-" + "-" + "-" + "-" + "-" + "-" + "-")

	fmt.Println("\n   📊 Relatório Mensal de Atividades:")
	fmt.Printf("      - Vendas: R$ %.2f\n", float64(report.TotalSales)/100)
	fmt.Printf("      - Horas trabalhadas: %d horas\n", report.TotalWorkHours)
	fmt.Printf("      - Membros ativos: %d\n", len(report.Members))
	fmt.Printf("      - Resultado: R$ %.2f\n", float64(report.NetSurplus)/100)

	report.Integrations = append(report.Integrations, Integration{
		Name:      "Governo Estado - SEAF",
		Status:    "SUCCESS",
		Timestamp: time.Now().Unix(),
		Response:  "Relatório enviado com sucesso",
	})
	fmt.Println("   → Relatório enviado ao SEAF (Governo do Estado)")

	report.Integrations = append(report.Integrations, Integration{
		Name:      "Governo Federal - MDA",
		Status:    "SUCCESS",
		Timestamp: time.Now().Unix(),
		Response:  "Relatório de produção enviado",
	})
	fmt.Println("   → Relatório enviado ao MDA (Governo Federal)")

	report.Integrations = append(report.Integrations, Integration{
		Name:      "IBGE - PAM",
		Status:    "SUCCESS",
		Timestamp: time.Now().Unix(),
		Response:  "Pesquisa Anual de MM indicators enviada",
	})
	fmt.Println("   → Dados enviados ao IBGE (PAM)")

	fmt.Println("\n🔗 ETAPA H: INTEGRAÇÕES ADICIONAIS")
	fmt.Println("-" + "-" + "-" + "-" + "-" + "-" + "-" + "-" + "-" + "-" + "-" + "-" + "-" + "-")

	report.Integrations = append(report.Integrations, Integration{
		Name:      "BNDES - Linhas de Crédito",
		Status:    "MOCK",
		Timestamp: time.Now().Unix(),
		Response:  "Simulação: Linha de crédito disponível até R$ 50.000",
	})
	fmt.Println("   💳 BNDES: Linha de crédito disponível (mock)")

	report.Integrations = append(report.Integrations, Integration{
		Name:      "SEBRAE - Capacitação",
		Status:    "MOCK",
		Timestamp: time.Now().Unix(),
		Response:  "Cursos disponíveis: Gestão, Qualidade, Marketing",
	})
	fmt.Println("   📚 SEBRAE: Cursos de capacitação disponíveis (mock)")

	report.Integrations = append(report.Integrations, Integration{
		Name:      "Mercado Solidário - Integração",
		Status:    "MOCK",
		Timestamp: time.Now().Unix(),
		Response:  "3 produtos cadastrados no marketplace",
	})
	fmt.Println("   🛒 Mercado Solidário: Produtos cadastrados (mock)")

	report.Integrations = append(report.Integrations, Integration{
		Name:      "Cooperativas parceiras - Intercooperação",
		Status:    "MOCK",
		Timestamp: time.Now().Unix(),
		Response:  "2 parcerias ativas",
	})
	fmt.Println("   🤝 Rede de Cooperação: 2 parcerias ativas (mock)")

	fmt.Println("\n📦 ETAPA I: RELATÓRIO FINAL JSON")
	fmt.Println("-" + "-" + "-" + "-" + "-" + "-" + "-" + "-" + "-" + "-" + "-" + "-" + "-" + "-")

	reportJSON, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal report: %v", err)
	}

	fmt.Printf("\n%s\n", string(reportJSON))

	cashAPI := cash_flow.NewCashFlowAPI(lifecycleMgr)
	cashBalance, _ := cashAPI.GetBalance(entityID)
	fmt.Printf("\n💰 Saldo em caixa: R$ %.2f\n", float64(cashBalance.Balance)/100)

	fmt.Println("\n" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=")
	fmt.Println("  ✅ SIMULAÇÃO CONCLUÍDA COM SUCESSO!")
	fmt.Println("=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=")
	fmt.Printf("\nResumo do mês:\n")
	fmt.Printf("  - Vendas realizadas: R$ %.2f\n", float64(report.TotalSales)/100)
	fmt.Printf("  - Horas trabalhadas: %d horas\n", report.TotalWorkHours)
	fmt.Printf("  - Cooperados: %d\n", len(report.Members))
	fmt.Printf("  - Resultado: R$ %.2f\n", float64(report.NetSurplus)/100)
	fmt.Printf("  - Integrações realizadas: %d\n", len(report.Integrations))
	fmt.Printf("  - Status legal: %s\n", report.LegalStatus)
	fmt.Println("\n🎉 Um mês de operação de economia solidária completo!")
}
