package governance

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/providentia/digna/lifecycle/pkg/lifecycle"
)

type DecisionStatus string

const (
	StatusDraft    DecisionStatus = "DRAFT"
	StatusApproved DecisionStatus = "APPROVED"
	StatusRejected DecisionStatus = "REJECTED"
	StatusArchived DecisionStatus = "ARCHIVED"
)

type DecisionRecord struct {
	ID           int64
	Title        string
	Content      string
	ContentHash  string
	Status       DecisionStatus
	DecisionDate time.Time
}

type CADSOLService struct {
	lifecycleManager lifecycle.LifecycleManager
}

func NewCADSOLService(lm lifecycle.LifecycleManager) *CADSOLService {
	return &CADSOLService{
		lifecycleManager: lm,
	}
}

func (cs *CADSOLService) RecordDecision(entityID string, record *DecisionRecord) error {
	if record.Title == "" {
		return fmt.Errorf("title cannot be empty")
	}
	if record.Content == "" {
		return fmt.Errorf("content cannot be empty")
	}

	record.ContentHash = generateHash(record.Content, entityID)

	if record.Status == "" {
		record.Status = StatusDraft
	}

	db, err := cs.lifecycleManager.GetConnection(entityID)
	if err != nil {
		return fmt.Errorf("failed to get connection: %w", err)
	}

	decisionDate := sql.NullInt64{}
	if !record.DecisionDate.IsZero() {
		decisionDate = sql.NullInt64{Int64: record.DecisionDate.Unix(), Valid: true}
	}

	result, err := db.Exec(
		"INSERT INTO decisions_log (title, content_hash, status, decision_date, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?)",
		record.Title, record.ContentHash, string(record.Status), decisionDate, time.Now().Unix(), time.Now().Unix(),
	)
	if err != nil {
		return fmt.Errorf("failed to insert decision: %w", err)
	}

	record.ID, _ = result.LastInsertId()
	return nil
}

func (cs *CADSOLService) GetDecisionByHash(entityID string, hash string) (*DecisionRecord, error) {
	db, err := cs.lifecycleManager.GetConnection(entityID)
	if err != nil {
		return nil, fmt.Errorf("failed to get connection: %w", err)
	}

	var record DecisionRecord
	var decisionDate sql.NullInt64

	err = db.QueryRow(
		"SELECT id, title, content_hash, status, decision_date FROM decisions_log WHERE content_hash = ?",
		hash,
	).Scan(&record.ID, &record.Title, &record.ContentHash, &record.Status, &decisionDate)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("decision not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to query decision: %w", err)
	}

	if decisionDate.Valid {
		record.DecisionDate = time.Unix(decisionDate.Int64, 0)
	}

	return &record, nil
}

func (cs *CADSOLService) UpdateDecisionStatus(entityID string, decisionID int64, status DecisionStatus) error {
	db, err := cs.lifecycleManager.GetConnection(entityID)
	if err != nil {
		return fmt.Errorf("failed to get connection: %w", err)
	}

	_, err = db.Exec(
		"UPDATE decisions_log SET status = ?, updated_at = ?, decision_date = ? WHERE id = ?",
		string(status), time.Now().Unix(), time.Now().Unix(), decisionID,
	)
	if err != nil {
		return fmt.Errorf("failed to update decision: %w", err)
	}

	return nil
}

func generateHash(content string, entityID string) string {
	// Usar entityID como salt para evitar colisões cross-tenant
	// Entidades diferentes com mesmo conteúdo terão hashes diferentes
	salted := fmt.Sprintf("%s:%s:DIGNA_SALT_v1", content, entityID)
	hash := sha256.Sum256([]byte(salted))
	return hex.EncodeToString(hash[:])
}
