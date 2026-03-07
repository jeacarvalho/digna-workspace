package governance

import (
	"time"

	"github.com/providentia/digna/core_lume/internal/domain"
	"github.com/providentia/digna/core_lume/internal/repository"
	"github.com/providentia/digna/core_lume/internal/service"
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
	decisionService *service.DecisionService
}

func NewService(lm lifecycle.LifecycleManager) *Service {
	decisionRepo := repository.NewSQLiteDecisionRepository(lm)
	return &Service{
		decisionService: service.NewDecisionService(decisionRepo),
	}
}

func (s *Service) RecordDecision(entityID string, title string, content string) (string, error) {
	if title == "" {
		return "", service.ErrTitleEmpty
	}
	if content == "" {
		return "", service.ErrContentEmpty
	}

	return s.decisionService.RecordDecision(entityID, title, content)
}

func (s *Service) GetDecisionByHash(entityID string, hash string) (*DecisionRecord, error) {
	decision, err := s.decisionService.GetDecisionByHash(entityID, hash)
	if err != nil {
		return nil, err
	}

	return &DecisionRecord{
		ID:           decision.ID,
		Title:        decision.Title,
		ContentHash:  decision.ContentHash,
		Status:       DecisionStatus(decision.Status),
		DecisionDate: time.Time{},
	}, nil
}

func (s *Service) UpdateDecisionStatus(entityID string, decisionID int64, status DecisionStatus) error {
	return s.decisionService.UpdateDecisionStatus(entityID, decisionID, domain.DecisionStatus(status))
}
