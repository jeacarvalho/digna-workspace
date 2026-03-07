package document

import (
	"github.com/providentia/digna/legal_facade/internal/document"
	"github.com/providentia/digna/lifecycle/pkg/lifecycle"
)

type Generator struct {
	*document.Generator
}

func NewGenerator(lm lifecycle.LifecycleManager) *Generator {
	return &Generator{
		Generator: document.NewGenerator(lm),
	}
}

func (g *Generator) GenerateAssemblyMinutes(entityID string, entityName string, status string) (string, error) {
	return g.Generator.GenerateAssemblyMinutes(entityID, entityName, status)
}

type IdentityGenerator struct {
	*document.IdentityGenerator
}

func NewIdentityGenerator(lm lifecycle.LifecycleManager) *IdentityGenerator {
	return &IdentityGenerator{
		IdentityGenerator: document.NewIdentityGenerator(lm),
	}
}

func (ig *IdentityGenerator) GenerateIdentityCard(entityID string, entityName string, status string) (string, error) {
	return ig.IdentityGenerator.GenerateIdentityCard(entityID, entityName, status)
}

type FormalizationSimulator struct {
	*document.FormalizationSimulator
}

func NewFormalizationSimulator(lm lifecycle.LifecycleManager) *FormalizationSimulator {
	return &FormalizationSimulator{
		FormalizationSimulator: document.NewFormalizationSimulator(lm),
	}
}

func (fs *FormalizationSimulator) CheckFormalizationCriteria(entityID string) (bool, error) {
	return fs.FormalizationSimulator.CheckFormalizationCriteria(entityID)
}

func (fs *FormalizationSimulator) GetEntityStatus(entityID string) (string, error) {
	return fs.FormalizationSimulator.GetEntityStatus(entityID)
}

func (fs *FormalizationSimulator) UpdateEntityStatus(entityID string, newStatus string) error {
	return fs.FormalizationSimulator.UpdateEntityStatus(entityID, newStatus)
}

func (fs *FormalizationSimulator) SimulateFormalization(entityID string) (bool, string, error) {
	return fs.FormalizationSimulator.SimulateFormalization(entityID)
}
