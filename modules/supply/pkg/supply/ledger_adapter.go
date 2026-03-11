package supply

import (
	"fmt"
	"time"

	"github.com/providentia/digna/core_lume/pkg/ledger"
	"github.com/providentia/digna/lifecycle/pkg/lifecycle"
)

// CoreLumeLedgerAdapter adapts supply.LedgerPort to core_lume.ledger.Service
type CoreLumeLedgerAdapter struct {
	ledgerService *ledger.Service
}

// NewCoreLumeLedgerAdapter creates a new adapter that connects supply to core_lume ledger
func NewCoreLumeLedgerAdapter(lm lifecycle.LifecycleManager) *CoreLumeLedgerAdapter {
	ledgerService := ledger.NewService(lm)
	return &CoreLumeLedgerAdapter{
		ledgerService: ledgerService,
	}
}

// RecordTransaction implements supply.LedgerPort.RecordTransaction
func (a *CoreLumeLedgerAdapter) RecordTransaction(entityID string, description string, postings []LedgerPosting) error {
	if a.ledgerService == nil {
		return fmt.Errorf("ledger service not initialized")
	}

	// Convert supply.LedgerPosting to ledger.Posting
	ledgerPostings := make([]ledger.Posting, len(postings))
	for i, p := range postings {
		ledgerPostings[i] = ledger.Posting{
			AccountID: p.AccountID,
			Amount:    p.Amount,
			Direction: ledger.Direction(p.Direction),
		}
	}

	// Create ledger transaction
	txn := &ledger.Transaction{
		EntityID:    entityID,
		Date:        time.Now(),
		Description: description,
		Reference:   "",
		Postings:    ledgerPostings,
	}

	// Validate transaction
	if err := txn.Validate(); err != nil {
		return fmt.Errorf("transaction validation failed: %w", err)
	}

	// Record transaction in core_lume ledger
	return a.ledgerService.RecordTransaction(entityID, txn)
}
