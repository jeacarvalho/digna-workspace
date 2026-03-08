package repository

import (
	"fmt"
	"testing"
	"time"

	"github.com/providentia/digna/core_lume/internal/domain"
	"github.com/providentia/digna/lifecycle/pkg/lifecycle"
)

func TestSQLiteLedgerRepository_CreateEntryWithPostingsTx(t *testing.T) {
	// Setup
	lifecycleMgr := lifecycle.NewSQLiteManager()
	defer lifecycleMgr.CloseAll()

	entityID := "test-entity-tx"
	repo := NewSQLiteLedgerRepository(lifecycleMgr)

	db, err := repo.GetDB(entityID)
	if err != nil {
		t.Fatalf("failed to get db: %v", err)
	}

	// Create test account
	accountCode := fmt.Sprintf("TX-ACCT-%d", time.Now().UnixNano())
	_, err = db.Exec(`
		INSERT INTO accounts (code, name, account_type, created_at)
		VALUES (?, ?, ?, ?)
	`, accountCode, "Test Account TX", "ASSET", time.Now().Unix())
	if err != nil {
		t.Fatalf("failed to create account: %v", err)
	}

	var accountID int64
	err = db.QueryRow("SELECT last_insert_rowid()").Scan(&accountID)
	if err != nil {
		t.Fatalf("failed to get account id: %v", err)
	}

	// Teste 1: Transação bem-sucedida
	t.Run("SuccessfulTransaction", func(t *testing.T) {
		entry := &domain.Entry{
			EntityID:    entityID,
			Date:        time.Now(),
			Description: "Test Transaction Entry",
			Reference:   "TX-001",
			CreatedAt:   time.Now(),
		}

		postings := []*domain.Posting{
			{
				EntityID:  entityID,
				AccountID: accountID,
				Amount:    1000,
				Direction: "DEBIT",
				CreatedAt: time.Now(),
			},
			{
				EntityID:  entityID,
				AccountID: accountID,
				Amount:    1000,
				Direction: "CREDIT",
				CreatedAt: time.Now(),
			},
		}

		entryID, err := repo.CreateEntryWithPostingsTx(entityID, entry, postings)
		if err != nil {
			t.Fatalf("transaction failed: %v", err)
		}

		if entryID == 0 {
			t.Error("expected non-zero entry id")
		}

		// Verifica se entry foi criado
		var count int
		err = db.QueryRow("SELECT COUNT(*) FROM entries WHERE id = ?", entryID).Scan(&count)
		if err != nil {
			t.Fatalf("failed to query entry: %v", err)
		}
		if count != 1 {
			t.Errorf("expected entry to exist, got count %d", count)
		}

		// Verifica se postings foram criados
		err = db.QueryRow("SELECT COUNT(*) FROM postings WHERE entry_id = ?", entryID).Scan(&count)
		if err != nil {
			t.Fatalf("failed to query postings: %v", err)
		}
		if count != 2 {
			t.Errorf("expected 2 postings, got %d", count)
		}

		// Verifica balanço
		balance, err := repo.GetBalance(entityID, accountID)
		if err != nil {
			t.Fatalf("failed to get balance: %v", err)
		}
		// 1000 debit - 1000 credit = 0
		if balance != 0 {
			t.Errorf("expected balance 0, got %d", balance)
		}
	})

	// Teste 2: Falha em posting deve fazer rollback
	t.Run("RollbackOnPostingFailure", func(t *testing.T) {
		// Contar entries antes
		var entriesBefore int
		err := db.QueryRow("SELECT COUNT(*) FROM entries").Scan(&entriesBefore)
		if err != nil {
			t.Fatalf("failed to count entries: %v", err)
		}

		entry := &domain.Entry{
			EntityID:    entityID,
			Date:        time.Now(),
			Description: "Should Fail Entry",
			Reference:   "TX-FAIL-001",
			CreatedAt:   time.Now(),
		}

		// Criar posting com account_id inexistente (vai falhar na FK)
		postings := []*domain.Posting{
			{
				EntityID:  entityID,
				AccountID: 999999, // account inexistente - vai falhar
				Amount:    500,
				Direction: "DEBIT",
				CreatedAt: time.Now(),
			},
		}

		_, err = repo.CreateEntryWithPostingsTx(entityID, entry, postings)
		if err == nil {
			t.Fatal("expected transaction to fail")
		}

		t.Logf("Expected error: %v", err)

		// Verifica se NADA foi persistido (rollback)
		var entriesAfter int
		err = db.QueryRow("SELECT COUNT(*) FROM entries").Scan(&entriesAfter)
		if err != nil {
			t.Fatalf("failed to count entries: %v", err)
		}

		if entriesAfter != entriesBefore {
			t.Errorf("entries should not increase after rollback: before=%d, after=%d", entriesBefore, entriesAfter)
		}
	})

	// Teste 3: Postings vazios devem funcionar (apenas entry)
	t.Run("EntryOnlyNoPostings", func(t *testing.T) {
		entry := &domain.Entry{
			EntityID:    entityID,
			Date:        time.Now(),
			Description: "Entry Only",
			Reference:   "TX-ENTRY-ONLY",
			CreatedAt:   time.Now(),
		}

		entryID, err := repo.CreateEntryWithPostingsTx(entityID, entry, []*domain.Posting{})
		if err != nil {
			t.Fatalf("transaction failed: %v", err)
		}

		// Verifica se entry foi criado
		var count int
		err = db.QueryRow("SELECT COUNT(*) FROM entries WHERE id = ?", entryID).Scan(&count)
		if err != nil {
			t.Fatalf("failed to query entry: %v", err)
		}
		if count != 1 {
			t.Errorf("expected entry to exist, got count %d", count)
		}

		// Verifica que não há postings
		err = db.QueryRow("SELECT COUNT(*) FROM postings WHERE entry_id = ?", entryID).Scan(&count)
		if err != nil {
			t.Fatalf("failed to query postings: %v", err)
		}
		if count != 0 {
			t.Errorf("expected 0 postings, got %d", count)
		}
	})
}
