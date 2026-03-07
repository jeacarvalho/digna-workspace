package domain

import "time"

type Direction string

const (
	Debit  Direction = "DEBIT"
	Credit Direction = "CREDIT"
)

type AccountType string

const (
	AccountTypeAsset     AccountType = "ASSET"
	AccountTypeLiability AccountType = "LIABILITY"
	AccountTypeRevenue   AccountType = "REVENUE"
	AccountTypeExpense   AccountType = "EXPENSE"
	AccountTypeEquity    AccountType = "EQUITY"
)

type Account struct {
	ID          int64
	Code        string
	Name        string
	ParentID    *int64
	AccountType AccountType
	CreatedAt   time.Time
}

type Entry struct {
	ID          int64
	EntityID    string
	Date        time.Time
	Description string
	Reference   string
	CreatedAt   time.Time
}

type Posting struct {
	ID        int64
	EntityID  string
	EntryID   int64
	AccountID int64
	Amount    int64
	Direction Direction
	CreatedAt time.Time
}

type DecisionStatus string

const (
	StatusDraft    DecisionStatus = "DRAFT"
	StatusApproved DecisionStatus = "APPROVED"
	StatusRejected DecisionStatus = "REJECTED"
	StatusArchived DecisionStatus = "ARCHIVED"
)

type Decision struct {
	ID           int64
	EntityID     string
	Title        string
	Content      string
	ContentHash  string
	Status       DecisionStatus
	DecisionDate *time.Time
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type WorkLog struct {
	ID           int64
	EntityID     string
	MemberID     string
	Minutes      int64
	ActivityType string
	LogDate      time.Time
	Description  string
	CreatedAt    time.Time
}
