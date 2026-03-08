package domain

import (
	"errors"
	"fmt"
)

// ValidationError represents a validation error
type ValidationError struct {
	Field   string
	Message string
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("validation error [%s]: %s", e.Field, e.Message)
}

// EntryValidator validates ledger entries before persistence
type EntryValidator struct {
	accountRepo AccountRepository
}

// NewEntryValidator creates a new entry validator
func NewEntryValidator(accountRepo AccountRepository) *EntryValidator {
	return &EntryValidator{
		accountRepo: accountRepo,
	}
}

// ValidateEntry validates an entry with its postings
// Returns nil if valid, or error with details if invalid
func (v *EntryValidator) ValidateEntry(entry *Entry, postings []*Posting) error {
	if entry == nil {
		return ValidationError{Field: "entry", Message: "entry cannot be nil"}
	}

	if len(postings) == 0 {
		return ValidationError{Field: "postings", Message: "at least one posting required"}
	}

	// Rule 1: Total debits must equal total credits (double-entry)
	var totalDebits, totalCredits int64
	for _, posting := range postings {
		if posting == nil {
			return ValidationError{Field: "postings", Message: "posting cannot be nil"}
		}

		// Rule 2: Amount must be positive
		if posting.Amount <= 0 {
			return ValidationError{
				Field:   fmt.Sprintf("posting[%d].amount", posting.AccountID),
				Message: "amount must be positive",
			}
		}

		// Accumulate totals
		if posting.Direction == "DEBIT" {
			totalDebits += posting.Amount
		} else if posting.Direction == "CREDIT" {
			totalCredits += posting.Amount
		} else {
			return ValidationError{
				Field:   fmt.Sprintf("posting[%d].direction", posting.AccountID),
				Message: "direction must be DEBIT or CREDIT",
			}
		}

		// Rule 3: Account must exist
		if v.accountRepo != nil {
			exists, err := v.accountRepo.Exists(entry.EntityID, posting.AccountID)
			if err != nil {
				return fmt.Errorf("failed to check account existence: %w", err)
			}
			if !exists {
				return ValidationError{
					Field:   fmt.Sprintf("posting[%d].account_id", posting.AccountID),
					Message: fmt.Sprintf("account %d does not exist", posting.AccountID),
				}
			}
		}
	}

	// Verify double-entry balance
	if totalDebits != totalCredits {
		return ValidationError{
			Field:   "balance",
			Message: fmt.Sprintf("debits (%d) must equal credits (%d)", totalDebits, totalCredits),
		}
	}

	return nil
}

// ValidateEntrySimple validates without account existence check
// Use this when account repository is not available
func (v *EntryValidator) ValidateEntrySimple(entry *Entry, postings []*Posting) error {
	if entry == nil {
		return ValidationError{Field: "entry", Message: "entry cannot be nil"}
	}

	if len(postings) == 0 {
		return ValidationError{Field: "postings", Message: "at least one posting required"}
	}

	// Rule 1: Total debits must equal total credits (double-entry)
	var totalDebits, totalCredits int64
	for i, posting := range postings {
		if posting == nil {
			return ValidationError{Field: "postings", Message: fmt.Sprintf("posting at index %d is nil", i)}
		}

		// Rule 2: Amount must be positive
		if posting.Amount <= 0 {
			return ValidationError{
				Field:   fmt.Sprintf("postings[%d].amount", i),
				Message: "amount must be positive",
			}
		}

		// Accumulate totals
		if posting.Direction == "DEBIT" {
			totalDebits += posting.Amount
		} else if posting.Direction == "CREDIT" {
			totalCredits += posting.Amount
		} else {
			return ValidationError{
				Field:   fmt.Sprintf("postings[%d].direction", i),
				Message: "direction must be DEBIT or CREDIT",
			}
		}
	}

	// Verify double-entry balance
	if totalDebits != totalCredits {
		return ValidationError{
			Field:   "balance",
			Message: fmt.Sprintf("debits (%d) must equal credits (%d)", totalDebits, totalCredits),
		}
	}

	return nil
}

// AccountRepository interface for account validation
type AccountRepository interface {
	Exists(entityID string, accountID int64) (bool, error)
}

// Common validation errors
var (
	ErrNilEntry         = errors.New("entry cannot be nil")
	ErrNoPostings       = errors.New("at least one posting required")
	ErrInvalidAmount    = errors.New("amount must be positive")
	ErrInvalidDirection = errors.New("direction must be DEBIT or CREDIT")
	ErrUnbalanced       = errors.New("debits must equal credits")
)
