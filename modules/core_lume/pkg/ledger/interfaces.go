package ledger

type LedgerRepository interface {
	SaveEntry(entryID int64, accountID int64, amount int64, direction Direction) error
	GetAccountBalance(entityID string, accountID int64) (int64, error)
}

type WorkRepository interface {
	GetAllMembersWork(entityID string) (map[string]int64, error)
}
