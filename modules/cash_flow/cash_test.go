package cash_flow

import (
	"testing"

	"github.com/providentia/digna/cash_flow/pkg/cash_flow"
	"github.com/providentia/digna/lifecycle/pkg/lifecycle"
)

func TestCashFlow_RecordAndGetBalance(t *testing.T) {
	lifecycleMgr := lifecycle.NewSQLiteManager()
	defer lifecycleMgr.CloseAll()

	entityID := "test_cash_entity"

	api := cash_flow.NewCashFlowAPI(lifecycleMgr)

	t.Run("Step1_RecordCredit", func(t *testing.T) {
		req := cash_flow.EntryRequest{
			EntityID:    entityID,
			Type:        "CREDIT",
			Amount:      10000,
			Category:    "SALES",
			Description: "Test credit entry",
		}

		resp, err := api.RecordEntry(req)
		if err != nil {
			t.Fatalf("Failed to record credit: %v", err)
		}
		if !resp.Success {
			t.Fatalf("Credit entry failed: %s", resp.Error)
		}
		t.Log("Step1_RecordCredit - PASS")
	})

	t.Run("Step2_GetBalance", func(t *testing.T) {
		resp, err := api.GetBalance(entityID)
		if err != nil {
			t.Fatalf("Failed to get balance: %v", err)
		}
		if resp.Balance != 10000 {
			t.Fatalf("Expected balance 10000, got %d", resp.Balance)
		}
		t.Logf("Step2_GetBalance - PASS (Balance: %d)", resp.Balance)
	})

	t.Run("Step3_RecordDebit", func(t *testing.T) {
		req := cash_flow.EntryRequest{
			EntityID:    entityID,
			Type:        "DEBIT",
			Amount:      3000,
			Category:    "EXPENSES",
			Description: "Test debit entry",
		}

		resp, err := api.RecordEntry(req)
		if err != nil {
			t.Fatalf("Failed to record debit: %v", err)
		}
		if !resp.Success {
			t.Fatalf("Debit entry failed: %s", resp.Error)
		}
		t.Log("Step3_RecordDebit - PASS")
	})

	t.Run("Step4_VerifyFinalBalance", func(t *testing.T) {
		resp, err := api.GetBalance(entityID)
		if err != nil {
			t.Fatalf("Failed to get balance: %v", err)
		}
		expected := int64(7000)
		if resp.Balance != expected {
			t.Fatalf("Expected balance %d, got %d", expected, resp.Balance)
		}
		t.Logf("Step4_VerifyFinalBalance - PASS (Balance: %d)", resp.Balance)
	})
}

func TestCashFlow_InvalidAmount(t *testing.T) {
	lifecycleMgr := lifecycle.NewSQLiteManager()
	defer lifecycleMgr.CloseAll()

	entityID := "test_invalid_entity"

	api := cash_flow.NewCashFlowAPI(lifecycleMgr)

	req := cash_flow.EntryRequest{
		EntityID:    entityID,
		Type:        "CREDIT",
		Amount:      0,
		Category:    "SALES",
		Description: "Invalid entry",
	}

	resp, err := api.RecordEntry(req)
	if err != nil {
		t.Logf("Error returned as expected: %v", err)
	}
	if resp != nil && resp.Success {
		t.Fatal("Expected failure for zero amount")
	}

	t.Log("TestCashFlow_InvalidAmount - PASS")
}

func TestCashFlow_MultipleEntries(t *testing.T) {
	lifecycleMgr := lifecycle.NewSQLiteManager()
	defer lifecycleMgr.CloseAll()

	entityID := "test_multiple_entity"

	api := cash_flow.NewCashFlowAPI(lifecycleMgr)

	entries := []cash_flow.EntryRequest{
		{EntityID: entityID, Type: "CREDIT", Amount: 5000, Category: "SALES", Description: "Sale 1"},
		{EntityID: entityID, Type: "CREDIT", Amount: 3000, Category: "OTHER_INCOME", Description: "Income 1"},
		{EntityID: entityID, Type: "DEBIT", Amount: 2000, Category: "EXPENSES", Description: "Expense 1"},
		{EntityID: entityID, Type: "DEBIT", Amount: 1000, Category: "SUPPLIERS", Description: "Supplier payment"},
	}

	for i, req := range entries {
		resp, err := api.RecordEntry(req)
		if err != nil {
			t.Fatalf("Failed to record entry %d: %v", i+1, err)
		}
		if !resp.Success {
			t.Fatalf("Entry %d failed: %s", i+1, resp.Error)
		}
	}

	resp, err := api.GetBalance(entityID)
	if err != nil {
		t.Fatalf("Failed to get balance: %v", err)
	}

	expected := int64(5000)
	if resp.Balance != expected {
		t.Fatalf("Expected balance %d, got %d", expected, resp.Balance)
	}

	t.Logf("TestCashFlow_MultipleEntries - PASS (Final Balance: %d)", resp.Balance)
}
