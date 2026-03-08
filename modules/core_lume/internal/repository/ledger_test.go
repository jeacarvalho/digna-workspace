package repository

import (
	"fmt"
	"testing"
	"time"

	"github.com/providentia/digna/lifecycle/pkg/lifecycle"
)

func TestSQLiteLedgerRepository_GetBalance(t *testing.T) {
	// Setup
	lifecycleMgr := lifecycle.NewSQLiteManager()
	defer lifecycleMgr.CloseAll()

	entityID := "test-entity-balance"
	repo := NewSQLiteLedgerRepository(lifecycleMgr)

	// Get database connection
	db, err := repo.GetDB(entityID)
	if err != nil {
		t.Fatalf("failed to get db: %v", err)
	}

	// Teste 1: Conta vazia deve retornar 0
	t.Run("EmptyAccountReturnsZero", func(t *testing.T) {
		balance, err := repo.GetBalance(entityID, 9999)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if balance != 0 {
			t.Errorf("expected balance 0, got %d", balance)
		}
	})

	// Create test account first
	accountCode := fmt.Sprintf("TEST-%d", time.Now().UnixNano())
	_, err = db.Exec(`
		INSERT INTO accounts (code, name, account_type, created_at)
		VALUES (?, ?, ?, ?)
	`, accountCode, "Test Account", "ASSET", time.Now().Unix())
	if err != nil {
		t.Fatalf("failed to create account: %v", err)
	}

	var accountID int64
	err = db.QueryRow("SELECT last_insert_rowid()").Scan(&accountID)
	if err != nil {
		t.Fatalf("failed to get account id: %v", err)
	}

	// Insert entry
	_, err = db.Exec(`
		INSERT INTO entries (entry_date, description, reference, created_at)
		VALUES (?, ?, ?, ?)
	`, time.Now().Unix(), "Test Entry", "TEST-001", time.Now().Unix())
	if err != nil {
		t.Fatalf("failed to insert entry: %v", err)
	}

	var entryID int64
	err = db.QueryRow("SELECT last_insert_rowid()").Scan(&entryID)
	if err != nil {
		t.Fatalf("failed to get entry id: %v", err)
	}

	// Insert DEBIT posting of 1000
	_, err = db.Exec(`
		INSERT INTO postings (entry_id, account_id, amount, direction, created_at)
		VALUES (?, ?, ?, ?, ?)
	`, entryID, accountID, 1000, "DEBIT", time.Now().Unix())
	if err != nil {
		t.Fatalf("failed to insert debit posting: %v", err)
	}

	// Teste 2: Apenas débitos
	t.Run("OnlyDebits", func(t *testing.T) {
		balance, err := repo.GetBalance(entityID, accountID)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if balance != 1000 {
			t.Errorf("expected balance 1000, got %d", balance)
		}
	})

	// Insert CREDIT posting of 300
	_, err = db.Exec(`
		INSERT INTO postings (entry_id, account_id, amount, direction, created_at)
		VALUES (?, ?, ?, ?, ?)
	`, entryID, accountID, 300, "CREDIT", time.Now().Unix())
	if err != nil {
		t.Fatalf("failed to insert credit posting: %v", err)
	}

	// Teste 3: Débitos e créditos mistos
	t.Run("MixedPostings", func(t *testing.T) {
		balance, err := repo.GetBalance(entityID, accountID)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		// 1000 (debit) - 300 (credit) = 700
		if balance != 700 {
			t.Errorf("expected balance 700, got %d", balance)
		}
	})

	// Insert CREDIT posting of 800
	_, err = db.Exec(`
		INSERT INTO postings (entry_id, account_id, amount, direction, created_at)
		VALUES (?, ?, ?, ?, ?)
	`, entryID, accountID, 800, "CREDIT", time.Now().Unix())
	if err != nil {
		t.Fatalf("failed to insert credit posting: %v", err)
	}

	// Teste 4: Saldo negativo
	t.Run("NegativeBalance", func(t *testing.T) {
		balance, err := repo.GetBalance(entityID, accountID)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		// 1000 - 300 - 800 = -100
		if balance != -100 {
			t.Errorf("expected balance -100, got %d", balance)
		}
	})
}

func TestSQLiteLedgerRepository_GetAccountBalance(t *testing.T) {
	// Setup
	lifecycleMgr := lifecycle.NewSQLiteManager()
	defer lifecycleMgr.CloseAll()

	entityID := "test-entity-account-balance"
	repo := NewSQLiteLedgerRepository(lifecycleMgr)

	db, err := repo.GetDB(entityID)
	if err != nil {
		t.Fatalf("failed to get db: %v", err)
	}

	// Create test account first
	accountCode := fmt.Sprintf("TEST-ACCT-%d", time.Now().UnixNano())
	_, err = db.Exec(`
		INSERT INTO accounts (code, name, account_type, created_at)
		VALUES (?, ?, ?, ?)
	`, accountCode, "Test Account 2", "ASSET", time.Now().Unix())
	if err != nil {
		t.Fatalf("failed to create account: %v", err)
	}

	var accountID int64
	err = db.QueryRow("SELECT last_insert_rowid()").Scan(&accountID)
	if err != nil {
		t.Fatalf("failed to get account id: %v", err)
	}

	// Insert entry
	_, err = db.Exec(`
		INSERT INTO entries (entry_date, description, reference, created_at)
		VALUES (?, ?, ?, ?)
	`, time.Now().Unix(), "Test Entry", "TEST-002", time.Now().Unix())
	if err != nil {
		t.Fatalf("failed to insert entry: %v", err)
	}

	var entryID int64
	err = db.QueryRow("SELECT last_insert_rowid()").Scan(&entryID)
	if err != nil {
		t.Fatalf("failed to get entry id: %v", err)
	}

	// Insert postings: 500 debit, 200 credit, 300 debit
	postings := []struct {
		amount    int64
		direction string
	}{
		{500, "DEBIT"},
		{200, "CREDIT"},
		{300, "DEBIT"},
	}

	for _, p := range postings {
		_, err = db.Exec(`
			INSERT INTO postings (entry_id, account_id, amount, direction, created_at)
			VALUES (?, ?, ?, ?, ?)
		`, entryID, accountID, p.amount, p.direction, time.Now().Unix())
		if err != nil {
			t.Fatalf("failed to insert posting: %v", err)
		}
	}

	// Teste: Verificar cálculo
	balance, err := repo.GetAccountBalance(entityID, accountID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// 500 + 300 - 200 = 600
	if balance != 600 {
		t.Errorf("expected balance 600, got %d", balance)
	}
}
