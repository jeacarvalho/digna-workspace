package document

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/providentia/digna/lifecycle/pkg/lifecycle"
)

const (
	MinDecisionsForFormalization = 3
)

type FormalizationSimulator struct {
	lifecycleManager lifecycle.LifecycleManager
}

func NewFormalizationSimulator(lm lifecycle.LifecycleManager) *FormalizationSimulator {
	return &FormalizationSimulator{
		lifecycleManager: lm,
	}
}

func (fs *FormalizationSimulator) CheckFormalizationCriteria(entityID string) (bool, error) {
	db, err := fs.lifecycleManager.GetConnection(entityID)
	if err != nil {
		return false, fmt.Errorf("failed to get connection: %w", err)
	}

	var decisionCount int
	err = db.QueryRow(
		"SELECT COUNT(*) FROM decisions_log",
	).Scan(&decisionCount)
	if err != nil {
		return false, fmt.Errorf("failed to count decisions: %w", err)
	}

	return decisionCount >= MinDecisionsForFormalization, nil
}

func (fs *FormalizationSimulator) GetEntityStatus(entityID string) (string, error) {
	db, err := fs.lifecycleManager.GetConnection(entityID)
	if err != nil {
		return "", fmt.Errorf("failed to get connection: %w", err)
	}

	var status string
	err = db.QueryRow(
		"SELECT status FROM sync_metadata WHERE id = 1",
	).Scan(&status)
	if err == sql.ErrNoRows {
		return "DREAM", nil
	}
	if err != nil {
		return "", fmt.Errorf("failed to get status: %w", err)
	}

	if status == "" {
		return "DREAM", nil
	}

	return status, nil
}

func (fs *FormalizationSimulator) UpdateEntityStatus(entityID string, newStatus string) error {
	db, err := fs.lifecycleManager.GetConnection(entityID)
	if err != nil {
		return fmt.Errorf("failed to get connection: %w", err)
	}

	_, err = db.Exec(
		"UPDATE sync_metadata SET status = ?, updated_at = ? WHERE id = 1",
		newStatus, time.Now().Unix(),
	)
	if err != nil {
		return fmt.Errorf("failed to update status: %w", err)
	}

	return nil
}

func (fs *FormalizationSimulator) SimulateFormalization(entityID string) (bool, string, error) {
	canFormalize, err := fs.CheckFormalizationCriteria(entityID)
	if err != nil {
		return false, "", err
	}

	if !canFormalize {
		currentStatus, _ := fs.GetEntityStatus(entityID)
		return false, currentStatus, fmt.Errorf("insufficient decisions: need at least %d", MinDecisionsForFormalization)
	}

	err = fs.UpdateEntityStatus(entityID, "FORMALIZED")
	if err != nil {
		return false, "", err
	}

	return true, "FORMALIZED", nil
}
