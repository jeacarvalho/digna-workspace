package social

import (
	"time"

	"github.com/providentia/digna/core_lume/internal/repository"
	"github.com/providentia/digna/core_lume/internal/service"
	"github.com/providentia/digna/lifecycle/pkg/lifecycle"
)

var (
	ErrInvalidMinutes = service.ErrInvalidMinutes
	ErrEmptyMemberID  = service.ErrEmptyMemberID
)

type WorkRecord struct {
	ID           int64
	MemberID     string
	Minutes      int64
	ActivityType string
	LogDate      time.Time
	Description  string
}

type Service struct {
	workService *service.WorkService
}

func NewService(lm lifecycle.LifecycleManager) *Service {
	workRepo := repository.NewSQLiteWorkRepository(lm)
	return &Service{
		workService: service.NewWorkService(workRepo),
	}
}

func (s *Service) RecordWork(entityID string, record *WorkRecord) error {
	if record.MemberID == "" {
		return ErrEmptyMemberID
	}
	if record.Minutes <= 0 {
		return ErrInvalidMinutes
	}

	return s.workService.RecordWork(entityID, record.MemberID, record.Minutes, record.ActivityType, record.Description)
}

func (s *Service) GetTotalWorkByMember(entityID string, memberID string) (int64, int64, error) {
	if memberID == "" {
		return 0, 0, ErrEmptyMemberID
	}

	return s.workService.GetTotalWorkByMember(entityID, memberID)
}

func (s *Service) GetAllMembersWork(entityID string) (map[string]int64, error) {
	return s.workService.GetAllMembersWork(entityID)
}
