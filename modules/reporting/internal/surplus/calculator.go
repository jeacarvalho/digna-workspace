package surplus

import (
	"fmt"

	"github.com/providentia/digna/core_lume/pkg/ledger"
)

const (
	AccountSales    int64 = 2
	AccountExpenses int64 = 5
)

type MemberShare struct {
	MemberID   string
	Minutes    int64
	Percentage float64
	Amount     int64
}

type SurplusCalculation struct {
	EntityID     string
	TotalSurplus int64
	TotalMinutes int64
	Members      []MemberShare
}

type SurplusRepository interface {
	GetAccountBalance(entityID string, accountID int64) (int64, error)
	GetAllMembersWork(entityID string) (map[string]int64, error)
}

type Calculator struct {
	surplusRepo SurplusRepository
}

func NewCalculator(ledgerRepo ledger.LedgerRepository, workRepo ledger.WorkRepository) *Calculator {
	return &Calculator{
		surplusRepo: &SurplusAdapter{
			ledgerRepo: ledgerRepo,
			workRepo:   workRepo,
		},
	}
}

type SurplusAdapter struct {
	ledgerRepo ledger.LedgerRepository
	workRepo   ledger.WorkRepository
}

func (a *SurplusAdapter) GetAccountBalance(entityID string, accountID int64) (int64, error) {
	return a.ledgerRepo.GetAccountBalance(entityID, accountID)
}

func (a *SurplusAdapter) GetAllMembersWork(entityID string) (map[string]int64, error) {
	return a.workRepo.GetAllMembersWork(entityID)
}

func (c *Calculator) CalculateSocialSurplus(entityID string) (*SurplusCalculation, error) {
	revenue, err := c.surplusRepo.GetAccountBalance(entityID, AccountSales)
	if err != nil {
		return nil, fmt.Errorf("failed to get revenue: %w", err)
	}

	expenses, err := c.surplusRepo.GetAccountBalance(entityID, AccountExpenses)
	if err != nil {
		return nil, fmt.Errorf("failed to get expenses: %w", err)
	}

	surplus := revenue - expenses

	memberMinutes, err := c.surplusRepo.GetAllMembersWork(entityID)
	if err != nil {
		return nil, fmt.Errorf("failed to get member work minutes: %w", err)
	}

	var totalMinutes int64
	for _, minutes := range memberMinutes {
		totalMinutes += minutes
	}

	members := make([]MemberShare, 0, len(memberMinutes))
	for memberID, minutes := range memberMinutes {
		percentage := 0.0
		if totalMinutes > 0 {
			percentage = float64(minutes) / float64(totalMinutes) * 100
		}

		amount := int64(0)
		if totalMinutes > 0 {
			amount = (surplus * minutes) / totalMinutes
		}

		members = append(members, MemberShare{
			MemberID:   memberID,
			Minutes:    minutes,
			Percentage: percentage,
			Amount:     amount,
		})
	}

	return &SurplusCalculation{
		EntityID:     entityID,
		TotalSurplus: surplus,
		TotalMinutes: totalMinutes,
		Members:      members,
	}, nil
}
