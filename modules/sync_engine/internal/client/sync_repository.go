package client

import (
	"database/sql"

	"github.com/providentia/digna/lifecycle/pkg/lifecycle"
)

type SyncRepository interface {
	GetConnection(entityID string) (*sql.DB, error)
	GetSyncState(entityID string) (int64, error)
	GetTotalSales(entityID string) (int64, error)
	GetTotalWorkAndMembers(entityID string) (int64, int64, error)
	GetDecisionCount(entityID string) (int64, error)
	GetLastEntryRef(entityID string) (string, error)
	GetLastDecisionHash(entityID string) (string, error)
	GetEntriesCountSince(entityID string, since int64) (int64, error)
	GetWorkLogsCountSince(entityID string, since int64) (int64, error)
	GetDecisionsCountSince(entityID string, since int64) (int64, error)
	UpdateLastSync(entityID string, timestamp int64) error
}

type SQLiteSyncRepository struct {
	lifecycleManager lifecycle.LifecycleManager
}

func NewSQLiteSyncRepository(lm lifecycle.LifecycleManager) *SQLiteSyncRepository {
	return &SQLiteSyncRepository{
		lifecycleManager: lm,
	}
}

func (r *SQLiteSyncRepository) GetConnection(entityID string) (*sql.DB, error) {
	return r.lifecycleManager.GetConnection(entityID)
}

func (r *SQLiteSyncRepository) GetSyncState(entityID string) (int64, error) {
	db, err := r.GetConnection(entityID)
	if err != nil {
		return 0, err
	}

	var lastSyncAt int64
	err = db.QueryRow(
		"SELECT COALESCE(last_sync_at, 0) FROM sync_metadata WHERE id = 1",
	).Scan(&lastSyncAt)
	if err == sql.ErrNoRows {
		return 0, nil
	}
	return lastSyncAt, err
}

func (r *SQLiteSyncRepository) GetTotalSales(entityID string) (int64, error) {
	db, err := r.GetConnection(entityID)
	if err != nil {
		return 0, err
	}

	var total int64
	err = db.QueryRow(
		"SELECT COALESCE(SUM(amount), 0) FROM postings WHERE direction = 'CREDIT' AND account_id = 2",
	).Scan(&total)
	return total, err
}

func (r *SQLiteSyncRepository) GetTotalWorkAndMembers(entityID string) (int64, int64, error) {
	db, err := r.GetConnection(entityID)
	if err != nil {
		return 0, 0, err
	}

	var totalMinutes, memberCount int64
	err = db.QueryRow(
		"SELECT COALESCE(SUM(minutes), 0), COUNT(DISTINCT member_id) FROM work_logs",
	).Scan(&totalMinutes, &memberCount)
	return totalMinutes / 60, memberCount, err
}

func (r *SQLiteSyncRepository) GetDecisionCount(entityID string) (int64, error) {
	db, err := r.GetConnection(entityID)
	if err != nil {
		return 0, err
	}

	var count int64
	err = db.QueryRow("SELECT COUNT(*) FROM decisions_log").Scan(&count)
	return count, err
}

func (r *SQLiteSyncRepository) GetLastEntryRef(entityID string) (string, error) {
	db, err := r.GetConnection(entityID)
	if err != nil {
		return "", err
	}

	var ref string
	err = db.QueryRow(
		"SELECT COALESCE(MAX(reference), '') FROM entries",
	).Scan(&ref)
	return ref, err
}

func (r *SQLiteSyncRepository) GetLastDecisionHash(entityID string) (string, error) {
	db, err := r.GetConnection(entityID)
	if err != nil {
		return "", err
	}

	var hash string
	err = db.QueryRow(
		"SELECT COALESCE(MAX(content_hash), '') FROM decisions_log",
	).Scan(&hash)
	return hash, err
}

func (r *SQLiteSyncRepository) GetEntriesCountSince(entityID string, since int64) (int64, error) {
	db, err := r.GetConnection(entityID)
	if err != nil {
		return 0, err
	}

	var count int64
	err = db.QueryRow(
		"SELECT COUNT(*) FROM entries WHERE created_at > ?",
		since,
	).Scan(&count)
	return count, err
}

func (r *SQLiteSyncRepository) GetWorkLogsCountSince(entityID string, since int64) (int64, error) {
	db, err := r.GetConnection(entityID)
	if err != nil {
		return 0, err
	}

	var count int64
	err = db.QueryRow(
		"SELECT COUNT(*) FROM work_logs WHERE created_at > ?",
		since,
	).Scan(&count)
	return count, err
}

func (r *SQLiteSyncRepository) GetDecisionsCountSince(entityID string, since int64) (int64, error) {
	db, err := r.GetConnection(entityID)
	if err != nil {
		return 0, err
	}

	var count int64
	err = db.QueryRow(
		"SELECT COUNT(*) FROM decisions_log WHERE created_at > ?",
		since,
	).Scan(&count)
	return count, err
}

func (r *SQLiteSyncRepository) UpdateLastSync(entityID string, timestamp int64) error {
	db, err := r.GetConnection(entityID)
	if err != nil {
		return err
	}

	_, err = db.Exec(
		"UPDATE sync_metadata SET last_sync_at = ? WHERE id = 1",
		timestamp,
	)
	return err
}
