package surplus

import (
	"database/sql"
	"fmt"

	"github.com/providentia/digna/lifecycle/pkg/lifecycle"
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

type Calculator struct {
	lifecycleManager lifecycle.LifecycleManager
}

func NewCalculator(lm lifecycle.LifecycleManager) *Calculator {
	return &Calculator{
		lifecycleManager: lm,
	}
}

func (c *Calculator) CalculateSocialSurplus(entityID string) (*SurplusCalculation, error) {
	db, err := c.lifecycleManager.GetConnection(entityID)
	if err != nil {
		return nil, fmt.Errorf("failed to get connection: %w", err)
	}

	surplus, err := c.calculateTotalSurplus(db)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate surplus: %w", err)
	}

	memberMinutes, totalMinutes, err := c.getMemberWorkMinutes(db)
	if err != nil {
		return nil, fmt.Errorf("failed to get member work minutes: %w", err)
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

func (c *Calculator) calculateTotalSurplus(db *sql.DB) (int64, error) {
	var revenue, expenses sql.NullInt64

	err := db.QueryRow(
		"SELECT COALESCE(SUM(CASE WHEN direction = 'CREDIT' THEN amount ELSE -amount END), 0) FROM postings WHERE account_id = ?",
		AccountSales,
	).Scan(&revenue)
	if err != nil {
		return 0, fmt.Errorf("failed to get revenue: %w", err)
	}

	err = db.QueryRow(
		"SELECT COALESCE(SUM(CASE WHEN direction = 'DEBIT' THEN amount ELSE -amount END), 0) FROM postings WHERE account_id = ?",
		AccountExpenses,
	).Scan(&expenses)
	if err != nil {
		return 0, fmt.Errorf("failed to get expenses: %w", err)
	}

	surplus := revenue.Int64 - expenses.Int64
	return surplus, nil
}

func (c *Calculator) getMemberWorkMinutes(db *sql.DB) (map[string]int64, int64, error) {
	rows, err := db.Query(
		"SELECT member_id, COALESCE(SUM(minutes), 0) FROM work_logs GROUP BY member_id",
	)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to query work logs: %w", err)
	}
	defer rows.Close()

	memberMinutes := make(map[string]int64)
	var totalMinutes int64

	for rows.Next() {
		var memberID string
		var minutes int64
		if err := rows.Scan(&memberID, &minutes); err != nil {
			return nil, 0, fmt.Errorf("failed to scan row: %w", err)
		}
		memberMinutes[memberID] = minutes
		totalMinutes += minutes
	}

	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("error iterating rows: %w", err)
	}

	return memberMinutes, totalMinutes, nil
}
