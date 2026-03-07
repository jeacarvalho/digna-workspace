package ledger

import (
	"database/sql"

	"github.com/providentia/digna/lifecycle/pkg/lifecycle"
)

type sqliteLedgerRepository struct {
	lifecycleManager lifecycle.LifecycleManager
}

func NewSQLiteLedgerRepository(lm lifecycle.LifecycleManager) LedgerRepository {
	return &sqliteLedgerRepository{
		lifecycleManager: lm,
	}
}

func (r *sqliteLedgerRepository) GetDB(entityID string) (*sql.DB, error) {
	return r.lifecycleManager.GetConnection(entityID)
}

func (r *sqliteLedgerRepository) SaveEntry(entryID int64, accountID int64, amount int64, direction Direction) error {
	// This method is not used directly - posting is saved via service
	return nil
}

func (r *sqliteLedgerRepository) GetAccountBalance(entityID string, accountID int64) (int64, error) {
	db, err := r.GetDB(entityID)
	if err != nil {
		return 0, err
	}

	var debitSum, creditSum int64
	err = db.QueryRow(
		"SELECT COALESCE(SUM(amount), 0) FROM postings WHERE account_id = ? AND direction = 'DEBIT'",
		accountID,
	).Scan(&debitSum)
	if err != nil {
		return 0, err
	}

	err = db.QueryRow(
		"SELECT COALESCE(SUM(amount), 0) FROM postings WHERE account_id = ? AND direction = 'CREDIT'",
		accountID,
	).Scan(&creditSum)
	if err != nil {
		return 0, err
	}

	return debitSum - creditSum, nil
}

type sqliteWorkRepository struct {
	lifecycleManager lifecycle.LifecycleManager
}

func NewSQLiteWorkRepository(lm lifecycle.LifecycleManager) WorkRepository {
	return &sqliteWorkRepository{
		lifecycleManager: lm,
	}
}

func (r *sqliteWorkRepository) GetDB(entityID string) (*sql.DB, error) {
	return r.lifecycleManager.GetConnection(entityID)
}

func (r *sqliteWorkRepository) GetAllMembersWork(entityID string) (map[string]int64, error) {
	db, err := r.GetDB(entityID)
	if err != nil {
		return nil, err
	}

	rows, err := db.Query(
		"SELECT member_id, COALESCE(SUM(minutes), 0) FROM work_logs GROUP BY member_id",
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make(map[string]int64)
	for rows.Next() {
		var memberID string
		var minutes int64
		if err := rows.Scan(&memberID, &minutes); err != nil {
			return nil, err
		}
		result[memberID] = minutes
	}

	return result, nil
}
