package social

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/providentia/digna/lifecycle/pkg/lifecycle"
)

var (
	ErrInvalidMinutes = fmt.Errorf("minutes must be positive")
	ErrEmptyMemberID  = fmt.Errorf("member_id cannot be empty")
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
	lifecycleManager lifecycle.LifecycleManager
}

func NewService(lm lifecycle.LifecycleManager) *Service {
	return &Service{
		lifecycleManager: lm,
	}
}

func (s *Service) RecordWork(entityID string, record *WorkRecord) error {
	if record.MemberID == "" {
		return ErrEmptyMemberID
	}
	if record.Minutes <= 0 {
		return ErrInvalidMinutes
	}

	db, err := s.lifecycleManager.GetConnection(entityID)
	if err != nil {
		return fmt.Errorf("failed to get connection: %w", err)
	}

	logDate := record.LogDate.Unix()
	if record.LogDate.IsZero() {
		logDate = time.Now().Unix()
	}

	_, err = db.Exec(
		"INSERT INTO work_logs (member_id, minutes, activity_type, log_date, description, created_at) VALUES (?, ?, ?, ?, ?, ?)",
		record.MemberID, record.Minutes, record.ActivityType, logDate, record.Description, time.Now().Unix(),
	)
	if err != nil {
		return fmt.Errorf("failed to insert work log: %w", err)
	}

	return nil
}

func (s *Service) GetTotalWorkByMember(entityID string, memberID string) (int64, int64, error) {
	if memberID == "" {
		return 0, 0, ErrEmptyMemberID
	}

	db, err := s.lifecycleManager.GetConnection(entityID)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to get connection: %w", err)
	}

	var totalMinutes sql.NullInt64
	err = db.QueryRow(
		"SELECT COALESCE(SUM(minutes), 0) FROM work_logs WHERE member_id = ?",
		memberID,
	).Scan(&totalMinutes)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to get total work: %w", err)
	}

	var count sql.NullInt64
	err = db.QueryRow(
		"SELECT COUNT(*) FROM work_logs WHERE member_id = ?",
		memberID,
	).Scan(&count)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to get work count: %w", err)
	}

	return totalMinutes.Int64, count.Int64, nil
}

func (s *Service) GetAllMembersWork(entityID string) (map[string]int64, error) {
	db, err := s.lifecycleManager.GetConnection(entityID)
	if err != nil {
		return nil, fmt.Errorf("failed to get connection: %w", err)
	}

	rows, err := db.Query(
		"SELECT member_id, COALESCE(SUM(minutes), 0) FROM work_logs GROUP BY member_id",
	)
	if err != nil {
		return nil, fmt.Errorf("failed to query work logs: %w", err)
	}
	defer rows.Close()

	result := make(map[string]int64)
	for rows.Next() {
		var memberID string
		var minutes int64
		if err := rows.Scan(&memberID, &minutes); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		result[memberID] = minutes
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return result, nil
}
