package domain

import "time"

type CashEntryType string

const (
	EntryTypeCredit CashEntryType = "CREDIT"
	EntryTypeDebit  CashEntryType = "DEBIT"
)

type CashEntry struct {
	ID          int64
	EntityID    string
	Type        CashEntryType
	Amount      int64
	Description string
	Category    string
	Date        time.Time
	CreatedAt   time.Time
}

type CashFlow struct {
	EntityID    string
	TotalCredit int64
	TotalDebit  int64
	Balance     int64
	Entries     []CashEntry
	PeriodStart time.Time
	PeriodEnd   time.Time
}
