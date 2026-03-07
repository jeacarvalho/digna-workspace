package repository

import (
	"time"

	"github.com/providentia/digna/core_lume/internal/domain"
)

type LedgerRepository interface {
	SaveEntry(entry *domain.Entry) (int64, error)
	SavePosting(posting *domain.Posting) error
	GetBalance(accountID int64) (int64, error)
	GetAccountBalance(entityID string, accountID int64) (int64, error)
}

type DecisionRepository interface {
	Save(decision *domain.Decision) (int64, error)
	FindByHash(entityID, hash string) (*domain.Decision, error)
	UpdateStatus(entityID string, decisionID int64, status domain.DecisionStatus) error
	FindByEntity(entityID string) ([]domain.Decision, error)
}

type WorkRepository interface {
	Save(work *domain.WorkLog) (int64, error)
	GetTotalByMember(entityID, memberID string) (int64, int64, error)
	GetAllMembersWork(entityID string) (map[string]int64, error)
	GetWorkByPeriod(entityID string, startDate, endDate time.Time) ([]domain.WorkLog, error)
}

type AccountRepository interface {
	FindByID(id int64) (*domain.Account, error)
	FindByCode(code string) (*domain.Account, error)
	ListAll() ([]domain.Account, error)
}
