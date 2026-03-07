package service

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/providentia/digna/cash_flow/internal/domain"
	"github.com/providentia/digna/core_lume/pkg/ledger"
	"github.com/providentia/digna/lifecycle/pkg/lifecycle"
)

const (
	AccountCash      int64 = 1
	AccountBank      int64 = 3
	AccountSuppliers int64 = 4
)

type SQLiteCashManager struct {
	lifecycleManager lifecycle.LifecycleManager
	ledgerService    *ledger.Service
}

func NewSQLiteCashManager(lm lifecycle.LifecycleManager) *SQLiteCashManager {
	return &SQLiteCashManager{
		lifecycleManager: lm,
		ledgerService:    ledger.NewService(lm),
	}
}

func (m *SQLiteCashManager) RecordEntry(entityID string, entry *domain.CashEntry) error {
	if entry.Amount <= 0 {
		return fmt.Errorf("amount must be positive")
	}
	if entityID == "" {
		return fmt.Errorf("entity_id cannot be empty")
	}

	var postings []ledger.Posting

	if entry.Type == domain.EntryTypeCredit {
		postings = []ledger.Posting{
			{
				AccountID: AccountCash,
				Amount:    entry.Amount,
				Direction: ledger.Debit,
			},
			{
				AccountID: ledger.GetAccountByName(entry.Category),
				Amount:    entry.Amount,
				Direction: ledger.Credit,
			},
		}
	} else {
		postings = []ledger.Posting{
			{
				AccountID: AccountCash,
				Amount:    entry.Amount,
				Direction: ledger.Credit,
			},
			{
				AccountID: ledger.GetAccountByName(entry.Category),
				Amount:    entry.Amount,
				Direction: ledger.Debit,
			},
		}
	}

	txn := &ledger.Transaction{
		Date:        entry.Date,
		Description: entry.Description,
		Reference:   fmt.Sprintf("CASH-%s-%d", entry.Type, time.Now().Unix()),
		Postings:    postings,
	}

	return m.ledgerService.RecordTransaction(entityID, txn)
}

func (m *SQLiteCashManager) GetBalance(entityID string) (int64, error) {
	balance, err := m.ledgerService.GetAccountBalance(entityID, AccountCash)
	if err != nil {
		return 0, fmt.Errorf("failed to get cash balance: %w", err)
	}
	return balance, nil
}

func (m *SQLiteCashManager) GetCashFlow(entityID string, startDate, endDate time.Time) (*domain.CashFlow, error) {
	db, err := m.lifecycleManager.GetConnection(entityID)
	if err != nil {
		return nil, fmt.Errorf("failed to get connection: %w", err)
	}

	flow := &domain.CashFlow{
		EntityID:    entityID,
		PeriodStart: startDate,
		PeriodEnd:   endDate,
	}

	err = db.QueryRow(`
		SELECT COALESCE(SUM(CASE WHEN direction = 'DEBIT' THEN amount ELSE 0 END), 0)
		FROM postings 
		WHERE account_id = ? AND created_at >= ? AND created_at <= ?`,
		AccountCash, startDate.Unix(), endDate.Unix(),
	).Scan(&flow.TotalCredit)
	if err != nil {
		return nil, fmt.Errorf("failed to get total credit: %w", err)
	}

	err = db.QueryRow(`
		SELECT COALESCE(SUM(CASE WHEN direction = 'CREDIT' THEN amount ELSE 0 END), 0)
		FROM postings 
		WHERE account_id = ? AND created_at >= ? AND created_at <= ?`,
		AccountCash, startDate.Unix(), endDate.Unix(),
	).Scan(&flow.TotalDebit)
	if err != nil {
		return nil, fmt.Errorf("failed to get total debit: %w", err)
	}

	flow.Balance = flow.TotalCredit - flow.TotalDebit

	entries, err := m.getEntriesByPeriod(db, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("failed to get entries: %w", err)
	}
	flow.Entries = entries

	return flow, nil
}

func (m *SQLiteCashManager) GetEntries(entityID string, limit int) ([]domain.CashEntry, error) {
	db, err := m.lifecycleManager.GetConnection(entityID)
	if err != nil {
		return nil, fmt.Errorf("failed to get connection: %w", err)
	}

	query := `
		SELECT p.id, p.created_at, p.amount, p.direction, e.description
		FROM postings p
		JOIN entries e ON p.entry_id = e.id
		WHERE p.account_id = ?
		ORDER BY p.created_at DESC
		LIMIT ?`

	rows, err := db.Query(query, AccountCash, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to query entries: %w", err)
	}
	defer rows.Close()

	var entries []domain.CashEntry
	for rows.Next() {
		var entry domain.CashEntry
		var direction string
		var createdAt int64

		err := rows.Scan(&entry.ID, &createdAt, &entry.Amount, &direction, &entry.Description)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		entry.Date = time.Unix(createdAt, 0)
		entry.CreatedAt = time.Unix(createdAt, 0)
		if direction == "DEBIT" {
			entry.Type = domain.EntryTypeCredit
		} else {
			entry.Type = domain.EntryTypeDebit
		}

		entries = append(entries, entry)
	}

	return entries, nil
}

func (m *SQLiteCashManager) getEntriesByPeriod(db *sql.DB, startDate, endDate time.Time) ([]domain.CashEntry, error) {
	query := `
		SELECT p.id, p.created_at, p.amount, p.direction, e.description
		FROM postings p
		JOIN entries e ON p.entry_id = e.id
		WHERE p.account_id = ? AND p.created_at >= ? AND p.created_at <= ?
		ORDER BY p.created_at DESC`

	rows, err := db.Query(query, AccountCash, startDate.Unix(), endDate.Unix())
	if err != nil {
		return nil, fmt.Errorf("failed to query entries: %w", err)
	}
	defer rows.Close()

	var entries []domain.CashEntry
	for rows.Next() {
		var entry domain.CashEntry
		var direction string
		var createdAt int64

		err := rows.Scan(&entry.ID, &createdAt, &entry.Amount, &direction, &entry.Description)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		entry.Date = time.Unix(createdAt, 0)
		entry.CreatedAt = time.Unix(createdAt, 0)
		if direction == "DEBIT" {
			entry.Type = domain.EntryTypeCredit
		} else {
			entry.Type = domain.EntryTypeDebit
		}

		entries = append(entries, entry)
	}

	return entries, nil
}

type CashFlowService struct {
	cashManager CashManager
}

func NewCashFlowService(lm lifecycle.LifecycleManager) *CashFlowService {
	return &CashFlowService{
		cashManager: NewSQLiteCashManager(lm),
	}
}

func (s *CashFlowService) RecordCredit(entityID, description, category string, amount int64) error {
	entry := &domain.CashEntry{
		Type:        domain.EntryTypeCredit,
		Amount:      amount,
		Description: description,
		Category:    category,
		Date:        time.Now(),
	}
	return s.cashManager.RecordEntry(entityID, entry)
}

func (s *CashFlowService) RecordDebit(entityID, description, category string, amount int64) error {
	entry := &domain.CashEntry{
		Type:        domain.EntryTypeDebit,
		Amount:      amount,
		Description: description,
		Category:    category,
		Date:        time.Now(),
	}
	return s.cashManager.RecordEntry(entityID, entry)
}

func (s *CashFlowService) GetBalance(entityID string) (int64, error) {
	return s.cashManager.GetBalance(entityID)
}

func (s *CashFlowService) GetCashFlow(entityID string, startDate, endDate time.Time) (*domain.CashFlow, error) {
	return s.cashManager.GetCashFlow(entityID, startDate, endDate)
}

func (s *CashFlowService) GetRecentEntries(entityID string, limit int) ([]domain.CashEntry, error) {
	if limit <= 0 {
		limit = 50
	}
	return s.cashManager.GetEntries(entityID, limit)
}
