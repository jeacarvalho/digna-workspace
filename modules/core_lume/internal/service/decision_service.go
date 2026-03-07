package service

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/providentia/digna/core_lume/internal/domain"
	"github.com/providentia/digna/core_lume/internal/repository"
)

type DecisionService struct {
	decisionRepo repository.DecisionRepository
}

func NewDecisionService(decisionRepo repository.DecisionRepository) *DecisionService {
	return &DecisionService{
		decisionRepo: decisionRepo,
	}
}

func (s *DecisionService) RecordDecision(entityID, title, content string) (string, error) {
	if title == "" {
		return "", fmt.Errorf("title cannot be empty")
	}
	if content == "" {
		return "", fmt.Errorf("content cannot be empty")
	}

	contentHash := generateHash(content, entityID)

	decision := &domain.Decision{
		EntityID:    entityID,
		Title:       title,
		Content:     content,
		ContentHash: contentHash,
		Status:      domain.StatusDraft,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	_, err := s.decisionRepo.Save(decision)
	if err != nil {
		return "", fmt.Errorf("failed to save decision: %w", err)
	}

	return contentHash, nil
}

func (s *DecisionService) GetDecisionByHash(entityID, hash string) (*domain.Decision, error) {
	return s.decisionRepo.FindByHash(entityID, hash)
}

func (s *DecisionService) UpdateDecisionStatus(entityID string, decisionID int64, status domain.DecisionStatus) error {
	return s.decisionRepo.UpdateStatus(entityID, decisionID, status)
}

func generateHash(content string, entityID string) string {
	salted := fmt.Sprintf("%s:%s:DIGNA_SALT_v1", content, entityID)
	hash := sha256.Sum256([]byte(salted))
	return hex.EncodeToString(hash[:])
}
