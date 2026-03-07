package ledger

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/providentia/digna/lifecycle/pkg/lifecycle"
)

var (
	ErrInvalidTransaction = errors.New("transaction sum must be zero (debits = credits)")
	ErrEmptyTransaction   = errors.New("transaction must have at least two postings")
	ErrDuplicatePostings  = errors.New("cannot have duplicate account postings in same entry")
)

type Direction string

const (
	Debit  Direction = "DEBIT"
	Credit Direction = "CREDIT"
)

type Posting struct {
	AccountID int64
	Amount    int64
	Direction Direction
}

type Transaction struct {
	ID          int64
	Date        time.Time
	Description string
	Reference   string
	Postings    []Posting
}

func (t *Transaction) Validate() error {
	if len(t.Postings) < 2 {
		return ErrEmptyTransaction
	}

	var sum int64
	accountSet := make(map[int64]bool)

	for _, p := range t.Postings {
		if accountSet[p.AccountID] {
			return ErrDuplicatePostings
		}
		accountSet[p.AccountID] = true

		switch p.Direction {
		case Debit:
			sum += p.Amount
		case Credit:
			sum -= p.Amount
		default:
			return fmt.Errorf("invalid direction: %s", p.Direction)
		}
	}

	if sum != 0 {
		return fmt.Errorf("%w: balance is %d", ErrInvalidTransaction, sum)
	}

	return nil
}

type Service struct {
	lifecycleManager lifecycle.LifecycleManager
}

func NewService(lm lifecycle.LifecycleManager) *Service {
	return &Service{
		lifecycleManager: lm,
	}
}

func (s *Service) RecordTransaction(entityID string, txn *Transaction) error {
	if err := txn.Validate(); err != nil {
		return fmt.Errorf("transaction validation failed: %w", err)
	}

	db, err := s.lifecycleManager.GetConnection(entityID)
	if err != nil {
		return fmt.Errorf("failed to get connection: %w", err)
	}

	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	entryDate := txn.Date.Unix()
	if txn.Date.IsZero() {
		entryDate = time.Now().Unix()
	}

	result, err := tx.Exec(
		"INSERT INTO entries (entry_date, description, reference, created_at) VALUES (?, ?, ?, ?)",
		entryDate, txn.Description, txn.Reference, time.Now().Unix(),
	)
	if err != nil {
		return fmt.Errorf("failed to insert entry: %w", err)
	}

	entryID, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get entry ID: %w", err)
	}

	for _, posting := range txn.Postings {
		_, err := tx.Exec(
			"INSERT INTO postings (entry_id, account_id, amount, direction, created_at) VALUES (?, ?, ?, ?, ?)",
			entryID, posting.AccountID, posting.Amount, string(posting.Direction), time.Now().Unix(),
		)
		if err != nil {
			return fmt.Errorf("failed to insert posting for account %d: %w", posting.AccountID, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	txn.ID = entryID
	return nil
}

func (s *Service) GetAccountBalance(entityID string, accountID int64) (int64, error) {
	db, err := s.lifecycleManager.GetConnection(entityID)
	if err != nil {
		return 0, fmt.Errorf("failed to get connection: %w", err)
	}

	var debitSum, creditSum sql.NullInt64

	err = db.QueryRow(
		"SELECT COALESCE(SUM(amount), 0) FROM postings WHERE account_id = ? AND direction = 'DEBIT'",
		accountID,
	).Scan(&debitSum)
	if err != nil {
		return 0, fmt.Errorf("failed to get debit sum: %w", err)
	}

	err = db.QueryRow(
		"SELECT COALESCE(SUM(amount), 0) FROM postings WHERE account_id = ? AND direction = 'CREDIT'",
		accountID,
	).Scan(&creditSum)
	if err != nil {
		return 0, fmt.Errorf("failed to get credit sum: %w", err)
	}

	balance := debitSum.Int64 - creditSum.Int64
	return balance, nil
}

func GetAccountByName(category string) int64 {
	accountMap := map[string]int64{
		"SALES":         2,
		"EXPENSES":      5,
		"SUPPLIERS":     4,
		"BANK":          3,
		"OTHER_INCOME":  6,
		"OTHER_EXPENSE": 7,
	}
	if id, ok := accountMap[category]; ok {
		return id
	}
	return 5
}
