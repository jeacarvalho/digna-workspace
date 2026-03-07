package service

import (
	"fmt"
	"time"

	"github.com/providentia/digna/cash_flow/internal/domain"
	"github.com/providentia/digna/core_lume/pkg/ledger"
)

const (
	AccountCash      int64 = 1
	AccountBank      int64 = 3
	AccountSuppliers int64 = 4
)

type LedgerPort interface {
	RecordTransaction(entityID string, txn *ledger.Transaction) error
	GetAccountBalance(entityID string, accountID int64) (int64, error)
}

type CashManager struct {
	ledger LedgerPort
}

func NewCashManager(ledgerPort LedgerPort) *CashManager {
	return &CashManager{
		ledger: ledgerPort,
	}
}

func (m *CashManager) RecordEntry(entityID string, entry *domain.CashEntry) error {
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

	return m.ledger.RecordTransaction(entityID, txn)
}

func (m *CashManager) GetBalance(entityID string) (int64, error) {
	return m.ledger.GetAccountBalance(entityID, AccountCash)
}

type CashFlowService struct {
	cashManager *CashManager
}

func NewCashFlowService(ledgerPort LedgerPort) *CashFlowService {
	return &CashFlowService{
		cashManager: NewCashManager(ledgerPort),
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
	balance, err := s.cashManager.GetBalance(entityID)
	if err != nil {
		return nil, err
	}

	return &domain.CashFlow{
		EntityID:    entityID,
		Balance:     balance,
		PeriodStart: startDate,
		PeriodEnd:   endDate,
	}, nil
}

func (s *CashFlowService) GetRecentEntries(entityID string, limit int) ([]domain.CashEntry, error) {
	return []domain.CashEntry{}, nil
}
