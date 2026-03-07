package document

import (
	"fmt"

	"github.com/providentia/digna/lifecycle/pkg/lifecycle"
)

const (
	MinDecisionsForFormalization = 3
)

type FormalizationSimulator struct {
	legalRepo LegalRepository
}

func NewFormalizationSimulator(lm lifecycle.LifecycleManager) *FormalizationSimulator {
	return &FormalizationSimulator{
		legalRepo: NewSQLiteLegalRepository(lm),
	}
}

func (fs *FormalizationSimulator) CheckFormalizationCriteria(entityID string) (bool, error) {
	decisionCount, err := fs.legalRepo.GetDecisionCount(entityID)
	if err != nil {
		return false, fmt.Errorf("failed to count decisions: %w", err)
	}

	return decisionCount >= MinDecisionsForFormalization, nil
}

func (fs *FormalizationSimulator) GetEntityStatus(entityID string) (string, error) {
	return fs.legalRepo.GetEntityStatus(entityID)
}

func (fs *FormalizationSimulator) UpdateEntityStatus(entityID string, newStatus string) error {
	return fs.legalRepo.UpdateEntityStatus(entityID, newStatus)
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
