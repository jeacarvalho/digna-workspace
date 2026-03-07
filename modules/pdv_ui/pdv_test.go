package pdv_ui_test

import (
	"os"
	"testing"

	"github.com/providentia/digna/lifecycle/pkg/lifecycle"
	"github.com/providentia/digna/pdv_ui/usecase"
)

func TestSprint02_DoD(t *testing.T) {
	dataDir := "../../data/entities"
	defer os.RemoveAll(dataDir)

	lifecycleMgr := lifecycle.NewSQLiteManager()
	defer lifecycleMgr.CloseAll()

	entityID := "Entidade_Teste"

	opHandler := usecase.NewOperationHandler(lifecycleMgr)
	sgHandler := usecase.NewSocialGovernanceHandler(lifecycleMgr)

	t.Run("Step1_Venda_5000", func(t *testing.T) {
		saleReq := usecase.SaleRequest{
			EntityID:      entityID,
			Amount:        5000,
			PaymentMethod: "PIX",
			Description:   "Venda Teste Sprint 02",
		}

		result, err := opHandler.RecordSale(saleReq)
		if err != nil {
			t.Fatalf("failed to record sale: %v", err)
		}
		if result.EntryID == 0 {
			t.Error("expected entry ID to be set")
		}
		t.Logf("Sale recorded with EntryID: %d", result.EntryID)
	})

	t.Run("Step2_Verificar_Saldo_Caixa", func(t *testing.T) {
		balance, err := opHandler.GetCashBalance(entityID)
		if err != nil {
			t.Fatalf("failed to get cash balance: %v", err)
		}
		if balance != 5000 {
			t.Errorf("expected cash balance 5000, got %d", balance)
		}
		t.Logf("Cash balance verified: %d (expected 5000)", balance)
	})

	t.Run("Step3_Registrar_Trabalho_ITG2002", func(t *testing.T) {
		workReq := usecase.WorkRequest{
			EntityID:     entityID,
			MemberID:     "socio_001",
			Minutes:      480,
			ActivityType: "PRODUCAO",
			Description:  "Trabalho cooperativo",
		}

		err := sgHandler.RecordWork(workReq)
		if err != nil {
			t.Fatalf("failed to record work: %v", err)
		}

		totalMinutes, count, err := sgHandler.GetMemberWorkCapital(entityID, "socio_001")
		if err != nil {
			t.Fatalf("failed to get member work capital: %v", err)
		}

		if totalMinutes != 480 {
			t.Errorf("expected 480 minutes, got %d", totalMinutes)
		}
		if count != 1 {
			t.Errorf("expected 1 work record, got %d", count)
		}

		t.Logf("Work capital verified: %d minutes, %d records", totalMinutes, count)
	})

	t.Run("Step4_Registrar_Decisao_CADSOL", func(t *testing.T) {
		decisionReq := usecase.DecisionRequest{
			EntityID: entityID,
			Title:    "Aprovação Estatuto",
			Content:  "Assembleia geral aprova novo estatuto da cooperativa",
		}

		hash, err := sgHandler.RecordDecision(decisionReq)
		if err != nil {
			t.Fatalf("failed to record decision: %v", err)
		}
		if hash == "" {
			t.Error("expected hash to be generated")
		}

		record, err := sgHandler.GetDecisionByHash(entityID, hash)
		if err != nil {
			t.Fatalf("failed to get decision by hash: %v", err)
		}

		if record.Title != decisionReq.Title {
			t.Errorf("expected title '%s', got '%s'", decisionReq.Title, record.Title)
		}
		if record.ContentHash != hash {
			t.Errorf("expected hash '%s', got '%s'", hash, record.ContentHash)
		}

		t.Logf("Decision recorded and verified with hash: %s", hash)
	})

	t.Run("Step5_Validar_Partidas_Dobradas", func(t *testing.T) {
		saleReq := usecase.SaleRequest{
			EntityID:      entityID,
			Amount:        10000,
			PaymentMethod: "CARTAO",
			Description:   "Teste validação partidas dobradas",
		}

		_, err := opHandler.RecordSale(saleReq)
		if err != nil {
			t.Fatalf("failed to record second sale: %v", err)
		}

		balance, err := opHandler.GetCashBalance(entityID)
		if err != nil {
			t.Fatalf("failed to get updated cash balance: %v", err)
		}

		expectedBalance := int64(15000)
		if balance != expectedBalance {
			t.Errorf("expected cash balance %d (5000+10000), got %d", expectedBalance, balance)
		}

		t.Logf("Double-entry validation passed. Total cash balance: %d", balance)
	})
}

func TestLedger_InvalidTransaction(t *testing.T) {
	dataDir := "../../data/entities"
	defer os.RemoveAll(dataDir)

	lifecycleMgr := lifecycle.NewSQLiteManager()
	defer lifecycleMgr.CloseAll()

	entityID := "Test_Invalid"
	opHandler := usecase.NewOperationHandler(lifecycleMgr)

	saleReq := usecase.SaleRequest{
		EntityID:      entityID,
		Amount:        0,
		PaymentMethod: "PIX",
	}

	_, err := opHandler.RecordSale(saleReq)
	if err == nil {
		t.Error("expected error for zero amount sale")
	}
}

func TestLedger_MultipleEntities_Isolation(t *testing.T) {
	dataDir := "../../data/entities"
	defer os.RemoveAll(dataDir)

	lifecycleMgr := lifecycle.NewSQLiteManager()
	defer lifecycleMgr.CloseAll()

	opHandler := usecase.NewOperationHandler(lifecycleMgr)

	entity1 := "Cooperativa_A"
	entity2 := "Cooperativa_B"

	_, err := opHandler.RecordSale(usecase.SaleRequest{
		EntityID:      entity1,
		Amount:        5000,
		PaymentMethod: "PIX",
	})
	if err != nil {
		t.Fatalf("failed to record sale for entity 1: %v", err)
	}

	_, err = opHandler.RecordSale(usecase.SaleRequest{
		EntityID:      entity2,
		Amount:        3000,
		PaymentMethod: "PIX",
	})
	if err != nil {
		t.Fatalf("failed to record sale for entity 2: %v", err)
	}

	balance1, err := opHandler.GetCashBalance(entity1)
	if err != nil {
		t.Fatalf("failed to get balance for entity 1: %v", err)
	}
	if balance1 != 5000 {
		t.Errorf("expected entity 1 balance 5000, got %d", balance1)
	}

	balance2, err := opHandler.GetCashBalance(entity2)
	if err != nil {
		t.Fatalf("failed to get balance for entity 2: %v", err)
	}
	if balance2 != 3000 {
		t.Errorf("expected entity 2 balance 3000, got %d", balance2)
	}

	t.Logf("Entity isolation verified: A=%d, B=%d", balance1, balance2)
}
