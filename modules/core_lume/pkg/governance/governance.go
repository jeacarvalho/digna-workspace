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
	ContentHash  string
	Status       DecisionStatus
	DecisionDate time.Time
}

type Service struct {
	lifecycleManager lifecycle.LifecycleManager
}

func NewService(lm lifecycle.LifecycleManager) *Service {
	return &Service{
		lifecycleManager: lm,
	}
}

func (s *Service) RecordDecision(entityID string, title string, content string) (string, error) {
	if title == "" {
		return "", fmt.Errorf("title cannot be empty")
	}
	if content == "" {
		return "", fmt.Errorf("content cannot be empty")
	}

	hash := generateHashWithSalt(content, entityID)

	db, err := s.lifecycleManager.GetConnection(entityID)
	if err != nil {
		return "", fmt.Errorf("failed to get connection: %w", err)
	}

	_, err = db.Exec(
		"INSERT INTO decisions_log (title, content_hash, status, decision_date, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?)",
		title, hash, string(StatusDraft), sql.NullInt64{}, time.Now().Unix(), time.Now().Unix(),
	)
	if err != nil {
		return "", fmt.Errorf("failed to insert decision: %w", err)
	}

	return hash, nil
}

func generateHashWithSalt(content string, entityID string) string {
	// Usar entityID como salt para evitar colisões cross-tenant
	salted := fmt.Sprintf("%s:%s:DIGNA_SALT_v1", content, entityID)
	hash := sha256.Sum256([]byte(salted))
	return hex.EncodeToString(hash[:])
}

func (s *Service) GetDecisionByHash(entityID string, hash string) (*DecisionRecord, error) {
	db, err := s.lifecycleManager.GetConnection(entityID)
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

func (s *Service) UpdateDecisionStatus(entityID string, decisionID int64, status DecisionStatus) error {
	db, err := s.lifecycleManager.GetConnection(entityID)
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

func generateHash(content string) string {
	hash := sha256.Sum256([]byte(content))
	return hex.EncodeToString(hash[:])
}
