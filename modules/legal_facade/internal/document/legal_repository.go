package document

import (
	"database/sql"
	"time"

	"github.com/providentia/digna/lifecycle/pkg/lifecycle"
)

type DecisionInfo struct {
	ID        int64
	Title     string
	Hash      string
	Status    string
	CreatedAt int64
}

type LegalRepository interface {
	GetConnection(entityID string) (*sql.DB, error)
	GetDecisionCount(entityID string) (int, error)
	GetEntityStatus(entityID string) (string, error)
	UpdateEntityStatus(entityID string, status string) error
	GetAllDecisions(entityID string) ([]DecisionInfo, error)
}

type SQLiteLegalRepository struct {
	lifecycleManager lifecycle.LifecycleManager
}

func NewSQLiteLegalRepository(lm lifecycle.LifecycleManager) *SQLiteLegalRepository {
	return &SQLiteLegalRepository{
		lifecycleManager: lm,
	}
}

func (r *SQLiteLegalRepository) GetConnection(entityID string) (*sql.DB, error) {
	return r.lifecycleManager.GetConnection(entityID)
}

func (r *SQLiteLegalRepository) GetDecisionCount(entityID string) (int, error) {
	db, err := r.GetConnection(entityID)
	if err != nil {
		return 0, err
	}

	var count int
	err = db.QueryRow(
		"SELECT COUNT(*) FROM decisions_log",
	).Scan(&count)
	return count, err
}

func (r *SQLiteLegalRepository) GetEntityStatus(entityID string) (string, error) {
	db, err := r.GetConnection(entityID)
	if err != nil {
		return "", err
	}

	var status string
	err = db.QueryRow(
		"SELECT status FROM sync_metadata WHERE id = 1",
	).Scan(&status)
	if err == sql.ErrNoRows {
		return "DREAM", nil
	}
	if err != nil {
		return "", err
	}

	if status == "" {
		return "DREAM", nil
	}

	return status, nil
}

func (r *SQLiteLegalRepository) UpdateEntityStatus(entityID string, status string) error {
	db, err := r.GetConnection(entityID)
	if err != nil {
		return err
	}

	_, err = db.Exec(
		"UPDATE sync_metadata SET status = ?, updated_at = ? WHERE id = 1",
		status, time.Now().Unix(),
	)
	return err
}

func (r *SQLiteLegalRepository) GetAllDecisions(entityID string) ([]DecisionInfo, error) {
	db, err := r.GetConnection(entityID)
	if err != nil {
		return nil, err
	}

	rows, err := db.Query(
		"SELECT id, title, content_hash, status, created_at FROM decisions_log ORDER BY created_at DESC",
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var decisions []DecisionInfo
	for rows.Next() {
		var d DecisionInfo
		if err := rows.Scan(&d.ID, &d.Title, &d.Hash, &d.Status, &d.CreatedAt); err != nil {
			return nil, err
		}
		decisions = append(decisions, d)
	}

	return decisions, nil
}
