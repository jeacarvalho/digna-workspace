package service

import (
	"time"

	"github.com/providentia/digna/cash_flow/internal/domain"
)

type CashManager interface {
	RecordEntry(entityID string, entry *domain.CashEntry) error
	GetBalance(entityID string) (int64, error)
	GetCashFlow(entityID string, startDate, endDate time.Time) (*domain.CashFlow, error)
	GetEntries(entityID string, limit int) ([]domain.CashEntry, error)
}
