package cash_flow

import (
	"time"

	"github.com/providentia/digna/cash_flow/internal/domain"
	"github.com/providentia/digna/cash_flow/internal/service"
	"github.com/providentia/digna/lifecycle/pkg/lifecycle"
)

type CashFlowAPI struct {
	service *service.CashFlowService
}

func NewCashFlowAPI(lm lifecycle.LifecycleManager) *CashFlowAPI {
	return &CashFlowAPI{
		service: service.NewCashFlowService(lm),
	}
}

type EntryRequest struct {
	EntityID    string `json:"entity_id"`
	Type        string `json:"type"`
	Amount      int64  `json:"amount"`
	Category    string `json:"category"`
	Description string `json:"description"`
}

type EntryResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error,omitempty"`
}

type BalanceResponse struct {
	EntityID string `json:"entity_id"`
	Balance  int64  `json:"balance"`
}

type CashFlowResponse struct {
	EntityID    string             `json:"entity_id"`
	TotalCredit int64              `json:"total_credit"`
	TotalDebit  int64              `json:"total_debit"`
	Balance     int64              `json:"balance"`
	Entries     []domain.CashEntry `json:"entries"`
	PeriodStart time.Time          `json:"period_start"`
	PeriodEnd   time.Time          `json:"period_end"`
}

func (api *CashFlowAPI) RecordEntry(req EntryRequest) (*EntryResponse, error) {
	var err error

	switch req.Type {
	case "CREDIT":
		err = api.service.RecordCredit(req.EntityID, req.Description, req.Category, req.Amount)
	case "DEBIT":
		err = api.service.RecordDebit(req.EntityID, req.Description, req.Category, req.Amount)
	default:
		return &EntryResponse{Success: false, Error: "invalid type"}, nil
	}

	if err != nil {
		return &EntryResponse{Success: false, Error: err.Error()}, nil
	}

	return &EntryResponse{Success: true}, nil
}

func (api *CashFlowAPI) GetBalance(entityID string) (*BalanceResponse, error) {
	balance, err := api.service.GetBalance(entityID)
	if err != nil {
		return nil, err
	}

	return &BalanceResponse{
		EntityID: entityID,
		Balance:  balance,
	}, nil
}

func (api *CashFlowAPI) GetCashFlow(entityID string, startDate, endDate time.Time) (*CashFlowResponse, error) {
	flow, err := api.service.GetCashFlow(entityID, startDate, endDate)
	if err != nil {
		return nil, err
	}

	return &CashFlowResponse{
		EntityID:    flow.EntityID,
		TotalCredit: flow.TotalCredit,
		TotalDebit:  flow.TotalDebit,
		Balance:     flow.Balance,
		Entries:     flow.Entries,
		PeriodStart: flow.PeriodStart,
		PeriodEnd:   flow.PeriodEnd,
	}, nil
}

func (api *CashFlowAPI) GetRecentEntries(entityID string, limit int) ([]domain.CashEntry, error) {
	return api.service.GetRecentEntries(entityID, limit)
}
