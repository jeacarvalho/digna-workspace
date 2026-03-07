package surplus

import (
	"github.com/providentia/digna/lifecycle/pkg/lifecycle"
	"github.com/providentia/digna/reporting/internal/surplus"
)

type MemberShare = surplus.MemberShare
type SurplusCalculation = surplus.SurplusCalculation

type Calculator struct {
	*surplus.Calculator
}

func NewCalculator(lm lifecycle.LifecycleManager) *Calculator {
	return &Calculator{
		Calculator: surplus.NewCalculator(lm),
	}
}

func (c *Calculator) CalculateSocialSurplus(entityID string) (*SurplusCalculation, error) {
	return c.Calculator.CalculateSocialSurplus(entityID)
}
