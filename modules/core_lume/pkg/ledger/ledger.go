package ledger

import (
	"time"

	"github.com/providentia/digna/core_lume/internal/domain"
	"github.com/providentia/digna/core_lume/internal/repository"
	"github.com/providentia/digna/core_lume/internal/service"
	"github.com/providentia/digna/lifecycle/pkg/lifecycle"
)

var (
	ErrInvalidTransaction = service.ErrInvalidTransaction
	ErrEmptyTransaction   = service.ErrEmptyTransaction
	ErrDuplicatePostings  = service.ErrDuplicatePostings
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
	EntityID    string
	Date        time.Time
	Description string
	Reference   string
	Postings    []Posting
}

func (t *Transaction) Validate() error {
	svcPostings := make([]service.Posting, len(t.Postings))
	for i, p := range t.Postings {
		svcPostings[i] = service.Posting{
			AccountID: p.AccountID,
			Amount:    p.Amount,
			Direction: domain.Direction(p.Direction),
		}
	}

	txn := &service.Transaction{
		EntityID:    t.EntityID,
		Date:        t.Date,
		Description: t.Description,
		Reference:   t.Reference,
		Postings:    svcPostings,
	}

	return service.ValidateTransaction(txn)
}

type Service struct {
	ledgerService *service.LedgerService
}

func NewService(lm lifecycle.LifecycleManager) *Service {
	ledgerRepo := repository.NewSQLiteLedgerRepository(lm)
	return &Service{
		ledgerService: service.NewLedgerService(ledgerRepo),
	}
}

func (s *Service) RecordTransaction(entityID string, txn *Transaction) error {
	svcPostings := make([]service.Posting, len(txn.Postings))
	for i, p := range txn.Postings {
		svcPostings[i] = service.Posting{
			AccountID: p.AccountID,
			Amount:    p.Amount,
			Direction: domain.Direction(p.Direction),
		}
	}

	serviceTxn := &service.Transaction{
		EntityID:    entityID,
		Date:        txn.Date,
		Description: txn.Description,
		Reference:   txn.Reference,
		Postings:    svcPostings,
	}

	err := s.ledgerService.RecordTransaction(serviceTxn)
	if err != nil {
		return err
	}

	txn.ID = serviceTxn.ID
	return nil
}

func (s *Service) GetAccountBalance(entityID string, accountID int64) (int64, error) {
	return s.ledgerService.GetAccountBalance(entityID, accountID)
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
