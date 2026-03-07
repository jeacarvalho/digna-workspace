package sprint03_test

import (
	"os"
	"strings"
	"testing"

	"github.com/providentia/digna/legal_facade/pkg/document"
	"github.com/providentia/digna/lifecycle/pkg/lifecycle"
	"github.com/providentia/digna/pdv_ui/usecase"
	"github.com/providentia/digna/reporting/pkg/surplus"
)

func TestSprint03_DoD(t *testing.T) {
	dataDir := "../../data/entities"
	defer os.RemoveAll(dataDir)

	lifecycleMgr := lifecycle.NewSQLiteManager()
	defer lifecycleMgr.CloseAll()

	entityID := "Entidade_Rateio_Teste"

	opHandler := usecase.NewOperationHandler(lifecycleMgr)
	sgHandler := usecase.NewSocialGovernanceHandler(lifecycleMgr)
	calculator := surplus.NewCalculator(lifecycleMgr)
	generator := document.NewGenerator(lifecycleMgr)
	formalizer := document.NewFormalizationSimulator(lifecycleMgr)

	t.Run("Step1_Criar_Socios_com_Horas_Diferentes", func(t *testing.T) {
		socio1Work := usecase.WorkRequest{
			EntityID:     entityID,
			MemberID:     "socio_001",
			Minutes:      600,
			ActivityType: "PRODUCAO",
			Description:  "Sócio 1 - 10 horas de trabalho",
		}
		socio2Work := usecase.WorkRequest{
			EntityID:     entityID,
			MemberID:     "socio_002",
			Minutes:      300,
			ActivityType: "PRODUCAO",
			Description:  "Sócio 2 - 5 horas de trabalho",
		}

		if err := sgHandler.RecordWork(socio1Work); err != nil {
			t.Fatalf("failed to record work for socio_001: %v", err)
		}
		if err := sgHandler.RecordWork(socio2Work); err != nil {
			t.Fatalf("failed to record work for socio_002: %v", err)
		}

		t.Log("2 sócios criados: socio_001 (600 min), socio_002 (300 min)")
	})

	t.Run("Step2_Realizar_Venda_10000", func(t *testing.T) {
		saleReq := usecase.SaleRequest{
			EntityID:      entityID,
			Amount:        10000,
			PaymentMethod: "PIX",
			Description:   "Venda para teste de rateio - R$ 100,00",
		}

		result, err := opHandler.RecordSale(saleReq)
		if err != nil {
			t.Fatalf("failed to record sale: %v", err)
		}

		t.Logf("Venda registrada: EntryID=%d, Valor=10000 (R$ 100,00)", result.EntryID)
	})

	t.Run("Step3_Calcular_Rateio_Social", func(t *testing.T) {
		calc, err := calculator.CalculateSocialSurplus(entityID)
		if err != nil {
			t.Fatalf("failed to calculate surplus: %v", err)
		}

		t.Logf("Excedente total: %d (R$ %.2f)", calc.TotalSurplus, float64(calc.TotalSurplus)/100)
		t.Logf("Total de minutos trabalhados: %d", calc.TotalMinutes)

		if len(calc.Members) != 2 {
			t.Fatalf("expected 2 members, got %d", len(calc.Members))
		}

		var socio1, socio2 *surplus.MemberShare
		for i := range calc.Members {
			if calc.Members[i].MemberID == "socio_001" {
				socio1 = &calc.Members[i]
			} else if calc.Members[i].MemberID == "socio_002" {
				socio2 = &calc.Members[i]
			}
		}

		if socio1 == nil || socio2 == nil {
			t.Fatal("could not find both members in calculation")
		}

		t.Logf("Socio_001: %d min (%.1f%%) = R$ %.2f",
			socio1.Minutes, socio1.Percentage, float64(socio1.Amount)/100)
		t.Logf("Socio_002: %d min (%.1f%%) = R$ %.2f",
			socio2.Minutes, socio2.Percentage, float64(socio2.Amount)/100)

		if socio1.Minutes <= socio2.Minutes {
			t.Error("socio_001 should have more minutes than socio_002")
		}

		if calc.TotalSurplus >= 0 {
			if socio1.Amount <= socio2.Amount {
				t.Errorf("socio_001 should receive more credit (%.2f vs %.2f)",
					float64(socio1.Amount)/100, float64(socio2.Amount)/100)
			}
		} else {
			if socio1.Amount >= socio2.Amount {
				t.Errorf("socio_001 should receive less debt (more negative): (%.2f vs %.2f)",
					float64(socio1.Amount)/100, float64(socio2.Amount)/100)
			}
		}

		t.Log("✅ Rateio validado: quem trabalhou mais recebeu mais crédito")
	})

	t.Run("Step4_Gerar_3_Decisoes", func(t *testing.T) {
		decisions := []struct {
			title   string
			content string
		}{
			{
				title:   "Aprovação do Estatuto Social",
				content: "Assembleia aprova o estatuto da cooperativa por unanimidade",
			},
			{
				title:   "Eleição do Conselho Administrativo",
				content: "Eleitos membros do conselho para gestão 2026-2028",
			},
			{
				title:   "Aprovação do Plano de Negócios",
				content: "Aprovado plano de negócios para o exercício fiscal",
			},
		}

		for _, d := range decisions {
			req := usecase.DecisionRequest{
				EntityID: entityID,
				Title:    d.title,
				Content:  d.content,
			}
			hash, err := sgHandler.RecordDecision(req)
			if err != nil {
				t.Fatalf("failed to record decision '%s': %v", d.title, err)
			}
			t.Logf("Decisão registrada: '%s' (hash: %s...)", d.title, hash[:16])
		}
	})

	t.Run("Step5_Verificar_Formalizacao", func(t *testing.T) {
		canFormalize, err := formalizer.CheckFormalizationCriteria(entityID)
		if err != nil {
			t.Fatalf("failed to check formalization criteria: %v", err)
		}

		if !canFormalize {
			t.Fatal("entity should be eligible for formalization with 3 decisions")
		}

		success, newStatus, err := formalizer.SimulateFormalization(entityID)
		if err != nil {
			t.Fatalf("failed to formalize: %v", err)
		}

		if !success {
			t.Fatal("formalization should succeed")
		}

		if newStatus != "FORMALIZED" {
			t.Errorf("expected status FORMALIZED, got %s", newStatus)
		}

		t.Logf("✅ Entidade formalizada com sucesso! Status: %s", newStatus)
	})

	t.Run("Step6_Gerar_Ata_Assembleia", func(t *testing.T) {
		assemblyDoc, err := generator.GenerateAssemblyMinutes(entityID, "Cooperativa Teste Rateio", "FORMALIZED")
		if err != nil {
			t.Fatalf("failed to generate assembly minutes: %v", err)
		}

		if assemblyDoc == "" {
			t.Fatal("generated assembly document is empty")
		}

		requiredSections := []string{
			"ATA DE ASSEMBLEIA GERAL EXTRAORDINÁRIA",
			"Cooperativa Teste Rateio",
			"DECISÕES REGISTRADAS",
			"Aprovação do Estatuto Social",
			"Eleição do Conselho Administrativo",
			"Aprovação do Plano de Negócios",
			"Hash de Auditoria",
			"Assinaturas digitais requeridas",
			"CADSOL",
		}

		for _, section := range requiredSections {
			if !strings.Contains(assemblyDoc, section) {
				t.Errorf("assembly document missing section: '%s'", section)
			}
		}

		t.Log("✅ Ata de Assembleia gerada com sucesso!")
		t.Logf("Document preview:\n---\n%s\n---", assemblyDoc[:500])
	})
}

func TestRateio_Proporcionalidade(t *testing.T) {
	dataDir := "../../data/entities"
	defer os.RemoveAll(dataDir)

	lifecycleMgr := lifecycle.NewSQLiteManager()
	defer lifecycleMgr.CloseAll()

	entityID := "Rateio_Proporcional"
	sgHandler := usecase.NewSocialGovernanceHandler(lifecycleMgr)
	calculator := surplus.NewCalculator(lifecycleMgr)
	opHandler := usecase.NewOperationHandler(lifecycleMgr)

	sgHandler.RecordWork(usecase.WorkRequest{
		EntityID: entityID, MemberID: "A", Minutes: 1000, ActivityType: "PROD",
	})
	sgHandler.RecordWork(usecase.WorkRequest{
		EntityID: entityID, MemberID: "B", Minutes: 500, ActivityType: "PROD",
	})
	sgHandler.RecordWork(usecase.WorkRequest{
		EntityID: entityID, MemberID: "C", Minutes: 500, ActivityType: "PROD",
	})

	opHandler.RecordSale(usecase.SaleRequest{
		EntityID: entityID, Amount: 20000, PaymentMethod: "PIX",
	})

	calc, err := calculator.CalculateSocialSurplus(entityID)
	if err != nil {
		t.Fatalf("failed to calculate: %v", err)
	}

	if calc.TotalMinutes != 2000 {
		t.Errorf("expected 2000 total minutes, got %d", calc.TotalMinutes)
	}

	var memberA, memberB, memberC *surplus.MemberShare
	for i := range calc.Members {
		switch calc.Members[i].MemberID {
		case "A":
			memberA = &calc.Members[i]
		case "B":
			memberB = &calc.Members[i]
		case "C":
			memberC = &calc.Members[i]
		}
	}

	if memberA == nil || memberB == nil || memberC == nil {
		t.Fatal("all members should be present")
	}

	if memberA.Minutes != 1000 || memberB.Minutes != 500 || memberC.Minutes != 500 {
		t.Error("minutes mismatch")
	}

	if memberA.Percentage != 50.0 {
		t.Errorf("member A should have 50%%, got %.1f%%", memberA.Percentage)
	}

	expectedA := calc.TotalSurplus / 2
	if memberA.Amount != expectedA {
		t.Errorf("member A should receive %d, got %d", expectedA, memberA.Amount)
	}

	t.Logf("Rateio proporcional validado: A=%d%%, B=%d%%, C=%d%%",
		int(memberA.Percentage), int(memberB.Percentage), int(memberC.Percentage))
}
