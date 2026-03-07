package service

import (
	"fmt"
	"time"

	"github.com/providentia/digna/core_lume/internal/domain"
	"github.com/providentia/digna/core_lume/internal/repository"
)

var (
	ErrInvalidMinutes = fmt.Errorf("minutes must be positive")
	ErrEmptyMemberID  = fmt.Errorf("member_id cannot be empty")
)

type WorkService struct {
	workRepo repository.WorkRepository
}

func NewWorkService(workRepo repository.WorkRepository) *WorkService {
	return &WorkService{
		workRepo: workRepo,
	}
}

func (s *WorkService) RecordWork(entityID, memberID string, minutes int64, activityType, description string) error {
	if memberID == "" {
		return ErrEmptyMemberID
	}
	if minutes <= 0 {
		return ErrInvalidMinutes
	}

	work := &domain.WorkLog{
		EntityID:     entityID,
		MemberID:     memberID,
		Minutes:      minutes,
		ActivityType: activityType,
		LogDate:      time.Now(),
		Description:  description,
		CreatedAt:    time.Now(),
	}

	_, err := s.workRepo.Save(work)
	if err != nil {
		return fmt.Errorf("failed to save work log: %w", err)
	}

	return nil
}

func (s *WorkService) GetTotalWorkByMember(entityID, memberID string) (int64, int64, error) {
	if memberID == "" {
		return 0, 0, ErrEmptyMemberID
	}

	return s.workRepo.GetTotalByMember(entityID, memberID)
}

func (s *WorkService) GetAllMembersWork(entityID string) (map[string]int64, error) {
	return s.workRepo.GetAllMembersWork(entityID)
}
