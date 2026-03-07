package tracker

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"fmt"

	"github.com/providentia/digna/lifecycle/pkg/lifecycle"
)

type Delta struct {
	TableName     string
	OperationType string
	RecordID      int64
	Timestamp     int64
	DataHash      string
}

type SyncState struct {
	EntityID       string
	LastSyncAt     int64
	LastEntryID    int64
	LastWorkLogID  int64
	LastDecisionID int64
	ChainDigest    string
	PendingChanges int64
}

type DeltaTracker struct {
	lifecycleManager lifecycle.LifecycleManager
}

func NewDeltaTracker(lm lifecycle.LifecycleManager) *DeltaTracker {
	return &DeltaTracker{
		lifecycleManager: lm,
	}
}

func (dt *DeltaTracker) GetCurrentState(entityID string) (*SyncState, error) {
	db, err := dt.lifecycleManager.GetConnection(entityID)
	if err != nil {
		return nil, fmt.Errorf("failed to get connection: %w", err)
	}

	state := &SyncState{EntityID: entityID}

	var lastSyncAt sql.NullInt64
	err = db.QueryRow(
		"SELECT last_sync_at FROM sync_metadata WHERE id = 1",
	).Scan(&lastSyncAt)
	if err == nil && lastSyncAt.Valid {
		state.LastSyncAt = lastSyncAt.Int64
	}

	var lastEntryID sql.NullInt64
	err = db.QueryRow(
		"SELECT MAX(id) FROM entries",
	).Scan(&lastEntryID)
	if err == nil && lastEntryID.Valid {
		state.LastEntryID = lastEntryID.Int64
	}

	var lastWorkLogID sql.NullInt64
	err = db.QueryRow(
		"SELECT MAX(id) FROM work_logs",
	).Scan(&lastWorkLogID)
	if err == nil && lastWorkLogID.Valid {
		state.LastWorkLogID = lastWorkLogID.Int64
	}

	var lastDecisionID sql.NullInt64
	err = db.QueryRow(
		"SELECT MAX(id) FROM decisions_log",
	).Scan(&lastDecisionID)
	if err == nil && lastDecisionID.Valid {
		state.LastDecisionID = lastDecisionID.Int64
	}

	chainDigest, err := dt.calculateChainDigest(db)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate chain digest: %w", err)
	}
	state.ChainDigest = chainDigest

	pendingChanges, err := dt.countPendingChanges(db, state.LastSyncAt)
	if err != nil {
		return nil, fmt.Errorf("failed to count pending changes: %w", err)
	}
	state.PendingChanges = pendingChanges

	return state, nil
}

func (dt *DeltaTracker) DetectDeltas(entityID string, since int64) ([]Delta, error) {
	db, err := dt.lifecycleManager.GetConnection(entityID)
	if err != nil {
		return nil, fmt.Errorf("failed to get connection: %w", err)
	}

	var deltas []Delta

	entryRows, err := db.Query(
		"SELECT id, created_at FROM entries WHERE created_at > ?",
		since,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to query entries: %w", err)
	}
	defer entryRows.Close()

	for entryRows.Next() {
		var id, createdAt int64
		if err := entryRows.Scan(&id, &createdAt); err != nil {
			return nil, fmt.Errorf("failed to scan entry: %w", err)
		}
		deltas = append(deltas, Delta{
			TableName:     "entries",
			OperationType: "INSERT",
			RecordID:      id,
			Timestamp:     createdAt,
			DataHash:      dt.hashRecord(fmt.Sprintf("entry:%d", id)),
		})
	}

	workRows, err := db.Query(
		"SELECT id, created_at FROM work_logs WHERE created_at > ?",
		since,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to query work_logs: %w", err)
	}
	defer workRows.Close()

	for workRows.Next() {
		var id, createdAt int64
		if err := workRows.Scan(&id, &createdAt); err != nil {
			return nil, fmt.Errorf("failed to scan work_log: %w", err)
		}
		deltas = append(deltas, Delta{
			TableName:     "work_logs",
			OperationType: "INSERT",
			RecordID:      id,
			Timestamp:     createdAt,
			DataHash:      dt.hashRecord(fmt.Sprintf("work:%d", id)),
		})
	}

	decisionRows, err := db.Query(
		"SELECT id, created_at FROM decisions_log WHERE created_at > ?",
		since,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to query decisions: %w", err)
	}
	defer decisionRows.Close()

	for decisionRows.Next() {
		var id, createdAt int64
		if err := decisionRows.Scan(&id, &createdAt); err != nil {
			return nil, fmt.Errorf("failed to scan decision: %w", err)
		}
		deltas = append(deltas, Delta{
			TableName:     "decisions_log",
			OperationType: "INSERT",
			RecordID:      id,
			Timestamp:     createdAt,
			DataHash:      dt.hashRecord(fmt.Sprintf("decision:%d", id)),
		})
	}

	return deltas, nil
}

func (dt *DeltaTracker) MarkSynced(entityID string, syncTime int64) error {
	db, err := dt.lifecycleManager.GetConnection(entityID)
	if err != nil {
		return fmt.Errorf("failed to get connection: %w", err)
	}

	_, err = db.Exec(
		"UPDATE sync_metadata SET last_sync_at = ? WHERE id = 1",
		syncTime,
	)
	if err != nil {
		return fmt.Errorf("failed to update sync metadata: %w", err)
	}

	return nil
}

func (dt *DeltaTracker) calculateChainDigest(db *sql.DB) (string, error) {
	var lastEntryHash, lastPostingHash string

	row := db.QueryRow(
		"SELECT COALESCE(MAX(reference), '') FROM entries",
	)
	if err := row.Scan(&lastEntryHash); err != nil {
		return "", err
	}

	row = db.QueryRow(
		"SELECT COALESCE(MAX(id), 0) FROM postings",
	)
	var lastPostingID int64
	if err := row.Scan(&lastPostingID); err != nil {
		return "", err
	}
	lastPostingHash = fmt.Sprintf("posting:%d", lastPostingID)

	hash := sha256.Sum256([]byte(lastEntryHash + lastPostingHash))
	return hex.EncodeToString(hash[:])[:16], nil
}

func (dt *DeltaTracker) countPendingChanges(db *sql.DB, since int64) (int64, error) {
	var count int64

	err := db.QueryRow(
		"SELECT COUNT(*) FROM entries WHERE created_at > ?",
		since,
	).Scan(&count)
	if err != nil {
		return 0, err
	}

	var workCount int64
	err = db.QueryRow(
		"SELECT COUNT(*) FROM work_logs WHERE created_at > ?",
		since,
	).Scan(&workCount)
	if err != nil {
		return 0, err
	}

	var decisionCount int64
	err = db.QueryRow(
		"SELECT COUNT(*) FROM decisions_log WHERE created_at > ?",
		since,
	).Scan(&decisionCount)
	if err != nil {
		return 0, err
	}

	return count + workCount + decisionCount, nil
}

func (dt *DeltaTracker) hashRecord(data string) string {
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:])[:8]
}

func (dt *DeltaTracker) HasChanges(entityID string) (bool, error) {
	state, err := dt.GetCurrentState(entityID)
	if err != nil {
		return false, err
	}

	return state.PendingChanges > 0, nil
}
