package budget

import (
	"time"

	"github.com/providentia/digna/budget/internal/service"
	"github.com/providentia/digna/cash_flow/pkg/cash_flow"
)

// CashFlowAdapter adapta a interface do cash_flow para o budget service
type CashFlowAdapter struct {
	cashFlowAPI *cash_flow.CashFlowAPI
}

// NewCashFlowAdapter cria um novo adaptador para cash_flow
func NewCashFlowAdapter(cashFlowAPI *cash_flow.CashFlowAPI) *CashFlowAdapter {
	return &CashFlowAdapter{
		cashFlowAPI: cashFlowAPI,
	}
}

// GetCashFlow implementa service.CashFlowPort.GetCashFlow
func (a *CashFlowAdapter) GetCashFlow(entityID string, startDate, endDate time.Time) (*service.CashFlowResponse, error) {
	if a.cashFlowAPI == nil {
		return &service.CashFlowResponse{
			EntityID:    entityID,
			TotalCredit: 0,
			TotalDebit:  0,
			Balance:     0,
			Entries:     []service.CashEntry{},
			PeriodStart: startDate,
			PeriodEnd:   endDate,
		}, nil
	}

	cfResponse, err := a.cashFlowAPI.GetCashFlow(entityID, startDate, endDate)
	if err != nil {
		return nil, err
	}

	// Converter cash_flow.CashFlowResponse para service.CashFlowResponse
	entries := make([]service.CashEntry, len(cfResponse.Entries))
	for i, entry := range cfResponse.Entries {
		entries[i] = service.CashEntry{
			ID:          entry.ID,
			EntityID:    entry.EntityID,
			Type:        string(entry.Type),
			Amount:      entry.Amount,
			Description: entry.Description,
			Category:    entry.Category,
			Date:        entry.Date,
			CreatedAt:   entry.CreatedAt,
		}
	}

	return &service.CashFlowResponse{
		EntityID:    cfResponse.EntityID,
		TotalCredit: cfResponse.TotalCredit,
		TotalDebit:  cfResponse.TotalDebit,
		Balance:     cfResponse.Balance,
		Entries:     entries,
		PeriodStart: cfResponse.PeriodStart,
		PeriodEnd:   cfResponse.PeriodEnd,
	}, nil
}
