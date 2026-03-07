package surplus

import (
	"github.com/providentia/digna/core_lume/pkg/ledger"
	"github.com/providentia/digna/lifecycle/pkg/lifecycle"
	"github.com/providentia/digna/reporting/internal/surplus"
)

type MemberShare = surplus.MemberShare
type SurplusCalculation = surplus.SurplusCalculation

type Calculator struct {
	calc *surplus.Calculator
}

func NewCalculator(lm lifecycle.LifecycleManager) *Calculator {
	ledgerRepo := ledger.NewSQLiteLedgerRepository(lm)
	workRepo := ledger.NewSQLiteWorkRepository(lm)
	return &Calculator{
		calc: surplus.NewCalculator(ledgerRepo, workRepo),
	}
}

func (c *Calculator) CalculateSocialSurplus(entityID string) (*SurplusCalculation, error) {
	return c.calc.CalculateSocialSurplus(entityID)
}
