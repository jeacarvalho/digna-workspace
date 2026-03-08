package repository

import (
	"time"

	"github.com/providentia/digna/core_lume/internal/domain"
)

type LedgerRepository interface {
	SaveEntry(entry *domain.Entry) (int64, error)
	SavePosting(posting *domain.Posting) error
	CreateEntryWithPostingsTx(entityID string, entry *domain.Entry, postings []*domain.Posting) (int64, error)
	GetBalance(entityID string, accountID int64) (int64, error)
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

type MemberRepository interface {
	Save(member *domain.Member) error
	FindByID(entityID, memberID string) (*domain.Member, error)
	FindByEmail(entityID, email string) (*domain.Member, error)
	ListByEntity(entityID string) ([]domain.Member, error)
	ListByRole(entityID string, role domain.MemberRole) ([]domain.Member, error)
	Update(member *domain.Member) error
	UpdateStatus(entityID, memberID string, status domain.MemberStatus) error
	CountByEntity(entityID string) (int, error)
	CountActiveByEntity(entityID string) (int, error)
}
