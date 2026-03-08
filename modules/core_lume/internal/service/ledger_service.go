package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/providentia/digna/core_lume/internal/domain"
	"github.com/providentia/digna/core_lume/internal/repository"
)

type Direction string

var (
	ErrInvalidTransaction = errors.New("transaction sum must be zero (debits = credits)")
	ErrEmptyTransaction   = errors.New("transaction must have at least two postings")
	ErrDuplicatePostings  = errors.New("cannot have duplicate account postings in same entry")
)

type LedgerService struct {
	ledgerRepo repository.LedgerRepository
}

func NewLedgerService(ledgerRepo repository.LedgerRepository) *LedgerService {
	return &LedgerService{
		ledgerRepo: ledgerRepo,
	}
}

type Transaction struct {
	ID          int64
	EntityID    string
	Date        time.Time
	Description string
	Reference   string
	Postings    []Posting
}

type Posting struct {
	AccountID int64
	Amount    int64
	Direction domain.Direction
}

func ValidateTransaction(t *Transaction) error {
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
		case domain.Debit:
			sum += p.Amount
		case domain.Credit:
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

func TranslateDirection(d Direction) domain.Direction {
	return domain.Direction(d)
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
		case domain.Debit:
			sum += p.Amount
		case domain.Credit:
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

func (s *LedgerService) RecordTransaction(txn *Transaction) error {
	if err := txn.Validate(); err != nil {
		return fmt.Errorf("transaction validation failed: %w", err)
	}

	entry := &domain.Entry{
		EntityID:    txn.EntityID,
		Date:        txn.Date,
		Description: txn.Description,
		Reference:   txn.Reference,
		CreatedAt:   time.Now(),
	}

	postings := make([]*domain.Posting, len(txn.Postings))
	for i, p := range txn.Postings {
		postings[i] = &domain.Posting{
			EntityID:  txn.EntityID,
			AccountID: p.AccountID,
			Amount:    p.Amount,
			Direction: domain.Direction(p.Direction),
			CreatedAt: time.Now(),
		}
	}

	entryID, err := s.ledgerRepo.CreateEntryWithPostingsTx(txn.EntityID, entry, postings)
	if err != nil {
		return fmt.Errorf("failed to record transaction: %w", err)
	}

	txn.ID = entryID
	return nil
}

func (s *LedgerService) GetAccountBalance(entityID string, accountID int64) (int64, error) {
	return s.ledgerRepo.GetAccountBalance(entityID, accountID)
}
