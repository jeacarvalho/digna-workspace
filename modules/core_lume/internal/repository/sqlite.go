package repository

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/providentia/digna/core_lume/internal/domain"
	"github.com/providentia/digna/lifecycle/pkg/lifecycle"
)

type SQLiteLedgerRepository struct {
	lifecycleManager lifecycle.LifecycleManager
}

func NewSQLiteLedgerRepository(lm lifecycle.LifecycleManager) *SQLiteLedgerRepository {
	return &SQLiteLedgerRepository{
		lifecycleManager: lm,
	}
}

func (r *SQLiteLedgerRepository) GetDB(entityID string) (*sql.DB, error) {
	return r.lifecycleManager.GetConnection(entityID)
}

func (r *SQLiteLedgerRepository) SaveEntry(entry *domain.Entry) (int64, error) {
	db, err := r.GetDB(entry.EntityID)
	if err != nil {
		return 0, fmt.Errorf("failed to get connection: %w", err)
	}

	entryDate := entry.Date.Unix()
	if entry.Date.IsZero() {
		entryDate = 0
	}

	result, err := db.Exec(
		"INSERT INTO entries (entry_date, description, reference, created_at) VALUES (?, ?, ?, ?)",
		entryDate, entry.Description, entry.Reference, entry.CreatedAt.Unix(),
	)
	if err != nil {
		return 0, fmt.Errorf("failed to insert entry: %w", err)
	}

	return result.LastInsertId()
}

func (r *SQLiteLedgerRepository) SavePosting(posting *domain.Posting) error {
	db, err := r.GetDB(posting.EntityID)
	if err != nil {
		return fmt.Errorf("failed to get connection: %w", err)
	}

	_, err = db.Exec(
		"INSERT INTO postings (entry_id, account_id, amount, direction, created_at) VALUES (?, ?, ?, ?, ?)",
		posting.EntryID, posting.AccountID, posting.Amount, string(posting.Direction), posting.CreatedAt.Unix(),
	)
	if err != nil {
		return fmt.Errorf("failed to insert posting: %w", err)
	}

	return nil
}

func (r *SQLiteLedgerRepository) GetBalance(entityID string, accountID int64) (int64, error) {
	db, err := r.GetDB(entityID)
	if err != nil {
		return 0, fmt.Errorf("failed to get connection: %w", err)
	}

	// Calculate balance: debits - credits
	// In accounting: positive balance = debit balance, negative = credit balance
	var balance int64
	err = db.QueryRow(`
		SELECT COALESCE(SUM(CASE 
			WHEN direction = 'DEBIT' THEN amount 
			WHEN direction = 'CREDIT' THEN -amount 
			ELSE 0 
		END), 0)
		FROM postings 
		WHERE account_id = ?
	`, accountID).Scan(&balance)

	if err != nil {
		return 0, fmt.Errorf("failed to calculate balance: %w", err)
	}

	return balance, nil
}

func (r *SQLiteLedgerRepository) GetAccountBalance(entityID string, accountID int64) (int64, error) {
	db, err := r.GetDB(entityID)
	if err != nil {
		return 0, fmt.Errorf("failed to get connection: %w", err)
	}

	var debitSum, creditSum sql.NullInt64

	err = db.QueryRow(
		"SELECT COALESCE(SUM(amount), 0) FROM postings WHERE account_id = ? AND direction = 'DEBIT'",
		accountID,
	).Scan(&debitSum)
	if err != nil {
		return 0, fmt.Errorf("failed to get debit sum: %w", err)
	}

	err = db.QueryRow(
		"SELECT COALESCE(SUM(amount), 0) FROM postings WHERE account_id = ? AND direction = 'CREDIT'",
		accountID,
	).Scan(&creditSum)
	if err != nil {
		return 0, fmt.Errorf("failed to get credit sum: %w", err)
	}

	return debitSum.Int64 - creditSum.Int64, nil
}

// CreateEntryWithPostingsTx cria entry e postings de forma atômica (transacional)
// Se qualquer operação falhar, todo o transaction é revertido (rollback)
func (r *SQLiteLedgerRepository) CreateEntryWithPostingsTx(entityID string, entry *domain.Entry, postings []*domain.Posting) (int64, error) {
	db, err := r.GetDB(entityID)
	if err != nil {
		return 0, fmt.Errorf("failed to get connection: %w", err)
	}

	// Inicia transação
	tx, err := db.Begin()
	if err != nil {
		return 0, fmt.Errorf("failed to begin transaction: %w", err)
	}

	// Garante rollback em caso de erro
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// 1. Insere o entry
	entryDate := entry.Date.Unix()
	if entry.Date.IsZero() {
		entryDate = 0
	}

	result, err := tx.Exec(
		"INSERT INTO entries (entry_date, description, reference, created_at) VALUES (?, ?, ?, ?)",
		entryDate, entry.Description, entry.Reference, entry.CreatedAt.Unix(),
	)
	if err != nil {
		return 0, fmt.Errorf("failed to insert entry: %w", err)
	}

	entryID, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to get entry id: %w", err)
	}

	// 2. Insere os postings
	for i, posting := range postings {
		_, err = tx.Exec(
			"INSERT INTO postings (entry_id, account_id, amount, direction, created_at) VALUES (?, ?, ?, ?, ?)",
			entryID, posting.AccountID, posting.Amount, string(posting.Direction), posting.CreatedAt.Unix(),
		)
		if err != nil {
			return 0, fmt.Errorf("failed to insert posting at index %d: %w", i, err)
		}
	}

	// 3. Commit da transação
	if err = tx.Commit(); err != nil {
		return 0, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return entryID, nil
}

type SQLiteDecisionRepository struct {
	lifecycleManager lifecycle.LifecycleManager
}

func NewSQLiteDecisionRepository(lm lifecycle.LifecycleManager) *SQLiteDecisionRepository {
	return &SQLiteDecisionRepository{
		lifecycleManager: lm,
	}
}

func (r *SQLiteDecisionRepository) GetDB(entityID string) (*sql.DB, error) {
	return r.lifecycleManager.GetConnection(entityID)
}

func (r *SQLiteDecisionRepository) Save(decision *domain.Decision) (int64, error) {
	db, err := r.GetDB(decision.EntityID)
	if err != nil {
		return 0, fmt.Errorf("failed to get connection: %w", err)
	}

	var decisionDate *int64
	if decision.DecisionDate != nil {
		ts := decision.DecisionDate.Unix()
		decisionDate = &ts
	}

	result, err := db.Exec(
		"INSERT INTO decisions_log (title, content_hash, status, decision_date, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?)",
		decision.Title, decision.ContentHash, string(decision.Status), decisionDate, decision.CreatedAt.Unix(), decision.UpdatedAt.Unix(),
	)
	if err != nil {
		return 0, fmt.Errorf("failed to insert decision: %w", err)
	}

	return result.LastInsertId()
}

func (r *SQLiteDecisionRepository) FindByHash(entityID string, hash string) (*domain.Decision, error) {
	db, err := r.GetDB(entityID)
	if err != nil {
		return nil, fmt.Errorf("failed to get connection: %w", err)
	}

	var decision domain.Decision
	var decisionDate sql.NullInt64

	err = db.QueryRow(
		"SELECT id, title, content_hash, status, decision_date FROM decisions_log WHERE content_hash = ?",
		hash,
	).Scan(&decision.ID, &decision.Title, &decision.ContentHash, &decision.Status, &decisionDate)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("decision not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to query decision: %w", err)
	}

	if decisionDate.Valid {
		ts := time.Unix(decisionDate.Int64, 0)
		decision.DecisionDate = &ts
	}

	return &decision, nil
}

func (r *SQLiteDecisionRepository) UpdateStatus(entityID string, decisionID int64, status domain.DecisionStatus) error {
	db, err := r.GetDB(entityID)
	if err != nil {
		return fmt.Errorf("failed to get connection: %w", err)
	}

	_, err = db.Exec(
		"UPDATE decisions_log SET status = ?, updated_at = ?, decision_date = ? WHERE id = ?",
		string(status), time.Now().Unix(), time.Now().Unix(), decisionID,
	)
	if err != nil {
		return fmt.Errorf("failed to update decision: %w", err)
	}

	return nil
}

func (r *SQLiteDecisionRepository) FindByEntity(entityID string) ([]domain.Decision, error) {
	db, err := r.GetDB(entityID)
	if err != nil {
		return nil, fmt.Errorf("failed to get connection: %w", err)
	}

	rows, err := db.Query(
		"SELECT id, title, content_hash, status, decision_date, created_at FROM decisions_log ORDER BY created_at DESC",
	)
	if err != nil {
		return nil, fmt.Errorf("failed to query decisions: %w", err)
	}
	defer rows.Close()

	var decisions []domain.Decision
	for rows.Next() {
		var d domain.Decision
		var decisionDate sql.NullInt64
		if err := rows.Scan(&d.ID, &d.Title, &d.ContentHash, &d.Status, &decisionDate, &d.CreatedAt); err != nil {
			return nil, err
		}
		if decisionDate.Valid {
			ts := time.Unix(decisionDate.Int64, 0)
			d.DecisionDate = &ts
		}
		decisions = append(decisions, d)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return decisions, nil
}

type SQLiteWorkRepository struct {
	lifecycleManager lifecycle.LifecycleManager
}

func NewSQLiteWorkRepository(lm lifecycle.LifecycleManager) *SQLiteWorkRepository {
	return &SQLiteWorkRepository{
		lifecycleManager: lm,
	}
}

func (r *SQLiteWorkRepository) GetDB(entityID string) (*sql.DB, error) {
	return r.lifecycleManager.GetConnection(entityID)
}

func (r *SQLiteWorkRepository) Save(work *domain.WorkLog) (int64, error) {
	db, err := r.GetDB(work.EntityID)
	if err != nil {
		return 0, fmt.Errorf("failed to get connection: %w", err)
	}

	logDate := work.LogDate.Unix()
	if work.LogDate.IsZero() {
		logDate = 0
	}

	result, err := db.Exec(
		"INSERT INTO work_logs (member_id, minutes, activity_type, log_date, description, created_at) VALUES (?, ?, ?, ?, ?, ?)",
		work.MemberID, work.Minutes, work.ActivityType, logDate, work.Description, work.CreatedAt.Unix(),
	)
	if err != nil {
		return 0, fmt.Errorf("failed to insert work log: %w", err)
	}

	return result.LastInsertId()
}

func (r *SQLiteWorkRepository) GetTotalByMember(entityID string, memberID string) (int64, int64, error) {
	db, err := r.GetDB(entityID)
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

func (r *SQLiteWorkRepository) GetAllMembersWork(entityID string) (map[string]int64, error) {
	db, err := r.GetDB(entityID)
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

func (r *SQLiteWorkRepository) GetWorkByPeriod(entityID string, startDate, endDate time.Time) ([]domain.WorkLog, error) {
	db, err := r.GetDB(entityID)
	if err != nil {
		return nil, fmt.Errorf("failed to get connection: %w", err)
	}

	rows, err := db.Query(
		"SELECT id, member_id, minutes, activity_type, log_date, description, created_at FROM work_logs WHERE log_date >= ? AND log_date <= ? ORDER BY log_date DESC",
		startDate.Unix(), endDate.Unix(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to query work logs: %w", err)
	}
	defer rows.Close()

	var logs []domain.WorkLog
	for rows.Next() {
		var w domain.WorkLog
		var logDate int64
		if err := rows.Scan(&w.ID, &w.MemberID, &w.Minutes, &w.ActivityType, &logDate, &w.Description, &w.CreatedAt); err != nil {
			return nil, err
		}
		w.LogDate = time.Unix(logDate, 0)
		logs = append(logs, w)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return logs, nil
}

func parseSkills(skillsStr string) []string {
	if skillsStr == "" || skillsStr == "[]" {
		return []string{}
	}
	// Simple parsing: remove brackets and split by comma
	skillsStr = strings.TrimPrefix(skillsStr, "[")
	skillsStr = strings.TrimSuffix(skillsStr, "]")
	if skillsStr == "" {
		return []string{}
	}
	return strings.Split(skillsStr, ",")
}

func joinStrings(strs []string, sep string) string {
	if len(strs) == 0 {
		return ""
	}
	result := strs[0]
	for i := 1; i < len(strs); i++ {
		result += sep + strs[i]
	}
	return result
}

type SQLiteMemberRepository struct {
	lifecycleManager lifecycle.LifecycleManager
}

func NewSQLiteMemberRepository(lm lifecycle.LifecycleManager) *SQLiteMemberRepository {
	return &SQLiteMemberRepository{lifecycleManager: lm}
}

func (r *SQLiteMemberRepository) GetDB(entityID string) (*sql.DB, error) {
	return r.lifecycleManager.GetConnection(entityID)
}

func (r *SQLiteMemberRepository) Save(member *domain.Member) error {
	if err := member.Validate(); err != nil {
		return fmt.Errorf("invalid member: %w", err)
	}

	db, err := r.GetDB(member.EntityID)
	if err != nil {
		return fmt.Errorf("failed to get connection: %w", err)
	}

	skillsJSON := "[]"
	if len(member.Skills) > 0 {
		skillsJSON = fmt.Sprintf("[%s]", joinStrings(member.Skills, ","))
	}

	_, err = db.Exec(
		`INSERT INTO members (id, entity_id, name, email, phone, cpf, role, status, joined_at, skills, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(id) DO UPDATE SET
			name = excluded.name,
			email = excluded.email,
			phone = excluded.phone,
			cpf = excluded.cpf,
			role = excluded.role,
			status = excluded.status,
			skills = excluded.skills,
			updated_at = excluded.updated_at`,
		member.ID, member.EntityID, member.Name, member.Email, member.Phone, member.CPF,
		string(member.Role), string(member.Status), member.JoinedAt.Unix(), skillsJSON,
		member.CreatedAt.Unix(), member.UpdatedAt.Unix(),
	)
	if err != nil {
		return fmt.Errorf("failed to save member: %w", err)
	}

	return nil
}

func (r *SQLiteMemberRepository) FindByID(entityID, memberID string) (*domain.Member, error) {
	db, err := r.GetDB(entityID)
	if err != nil {
		return nil, fmt.Errorf("failed to get connection: %w", err)
	}

	var member domain.Member
	var joinedAt int64
	var createdAt int64
	var updatedAt int64
	var skillsStr string

	err = db.QueryRow(
		`SELECT id, entity_id, name, email, phone, cpf, role, status, joined_at, skills, created_at, updated_at
		FROM members WHERE id = ? AND entity_id = ?`,
		memberID, entityID,
	).Scan(&member.ID, &member.EntityID, &member.Name, &member.Email, &member.Phone, &member.CPF,
		&member.Role, &member.Status, &joinedAt, &skillsStr, &createdAt, &updatedAt)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("member not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to query member: %w", err)
	}

	member.JoinedAt = time.Unix(joinedAt, 0)
	member.CreatedAt = time.Unix(createdAt, 0)
	member.UpdatedAt = time.Unix(updatedAt, 0)
	member.Skills = parseSkills(skillsStr)

	return &member, nil
}

func (r *SQLiteMemberRepository) FindByEmail(entityID, email string) (*domain.Member, error) {
	db, err := r.GetDB(entityID)
	if err != nil {
		return nil, fmt.Errorf("failed to get connection: %w", err)
	}

	var member domain.Member
	var joinedAt int64
	var createdAt int64
	var updatedAt int64
	var skillsStr string

	err = db.QueryRow(
		`SELECT id, entity_id, name, email, phone, cpf, role, status, joined_at, skills, created_at, updated_at
		FROM members WHERE email = ? AND entity_id = ?`,
		email, entityID,
	).Scan(&member.ID, &member.EntityID, &member.Name, &member.Email, &member.Phone, &member.CPF,
		&member.Role, &member.Status, &joinedAt, &skillsStr, &createdAt, &updatedAt)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("member not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to query member: %w", err)
	}

	member.JoinedAt = time.Unix(joinedAt, 0)
	member.CreatedAt = time.Unix(createdAt, 0)
	member.UpdatedAt = time.Unix(updatedAt, 0)
	member.Skills = parseSkills(skillsStr)

	return &member, nil
}

func (r *SQLiteMemberRepository) ListByEntity(entityID string) ([]domain.Member, error) {
	db, err := r.GetDB(entityID)
	if err != nil {
		return nil, fmt.Errorf("failed to get connection: %w", err)
	}

	rows, err := db.Query(
		`SELECT id, entity_id, name, email, phone, cpf, role, status, joined_at, skills, created_at, updated_at
		FROM members WHERE entity_id = ? ORDER BY name`,
		entityID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to query members: %w", err)
	}
	defer rows.Close()

	var members []domain.Member
	for rows.Next() {
		var member domain.Member
		var joinedAt int64
		var createdAt int64
		var updatedAt int64
		var skillsStr string

		if err := rows.Scan(&member.ID, &member.EntityID, &member.Name, &member.Email, &member.Phone, &member.CPF,
			&member.Role, &member.Status, &joinedAt, &skillsStr, &createdAt, &updatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan member: %w", err)
		}

		member.JoinedAt = time.Unix(joinedAt, 0)
		member.CreatedAt = time.Unix(createdAt, 0)
		member.UpdatedAt = time.Unix(updatedAt, 0)
		member.Skills = parseSkills(skillsStr)
		members = append(members, member)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return members, nil
}

func (r *SQLiteMemberRepository) ListByRole(entityID string, role domain.MemberRole) ([]domain.Member, error) {
	db, err := r.GetDB(entityID)
	if err != nil {
		return nil, fmt.Errorf("failed to get connection: %w", err)
	}

	rows, err := db.Query(
		`SELECT id, entity_id, name, email, phone, cpf, role, status, joined_at, skills, created_at, updated_at
		FROM members WHERE entity_id = ? AND role = ? ORDER BY name`,
		entityID, string(role),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to query members: %w", err)
	}
	defer rows.Close()

	var members []domain.Member
	for rows.Next() {
		var member domain.Member
		var joinedAt int64
		var createdAt int64
		var updatedAt int64
		var skillsStr string

		if err := rows.Scan(&member.ID, &member.EntityID, &member.Name, &member.Email, &member.Phone, &member.CPF,
			&member.Role, &member.Status, &joinedAt, &skillsStr, &createdAt, &updatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan member: %w", err)
		}

		member.JoinedAt = time.Unix(joinedAt, 0)
		member.CreatedAt = time.Unix(createdAt, 0)
		member.UpdatedAt = time.Unix(updatedAt, 0)
		member.Skills = parseSkills(skillsStr)
		members = append(members, member)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return members, nil
}

func (r *SQLiteMemberRepository) Update(member *domain.Member) error {
	if err := member.Validate(); err != nil {
		return fmt.Errorf("invalid member: %w", err)
	}

	db, err := r.GetDB(member.EntityID)
	if err != nil {
		return fmt.Errorf("failed to get connection: %w", err)
	}

	skillsJSON := "[]"
	if len(member.Skills) > 0 {
		skillsJSON = fmt.Sprintf("[%s]", joinStrings(member.Skills, ","))
	}

	member.UpdatedAt = time.Now()

	_, err = db.Exec(
		`UPDATE members SET name = ?, email = ?, phone = ?, cpf = ?, role = ?, status = ?, skills = ?, updated_at = ?
		WHERE id = ? AND entity_id = ?`,
		member.Name, member.Email, member.Phone, member.CPF, string(member.Role), string(member.Status),
		skillsJSON, member.UpdatedAt.Unix(), member.ID, member.EntityID,
	)
	if err != nil {
		return fmt.Errorf("failed to update member: %w", err)
	}

	return nil
}

func (r *SQLiteMemberRepository) UpdateStatus(entityID, memberID string, status domain.MemberStatus) error {
	db, err := r.GetDB(entityID)
	if err != nil {
		return fmt.Errorf("failed to get connection: %w", err)
	}

	_, err = db.Exec(
		`UPDATE members SET status = ?, updated_at = ? WHERE id = ? AND entity_id = ?`,
		string(status), time.Now().Unix(), memberID, entityID,
	)
	if err != nil {
		return fmt.Errorf("failed to update member status: %w", err)
	}

	return nil
}

func (r *SQLiteMemberRepository) CountByEntity(entityID string) (int, error) {
	db, err := r.GetDB(entityID)
	if err != nil {
		return 0, fmt.Errorf("failed to get connection: %w", err)
	}

	var count int
	err = db.QueryRow(
		`SELECT COUNT(*) FROM members WHERE entity_id = ?`,
		entityID,
	).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count members: %w", err)
	}

	return count, nil
}

func (r *SQLiteMemberRepository) CountActiveByEntity(entityID string) (int, error) {
	db, err := r.GetDB(entityID)
	if err != nil {
		return 0, fmt.Errorf("failed to get connection: %w", err)
	}

	var count int
	err = db.QueryRow(
		`SELECT COUNT(*) FROM members WHERE entity_id = ? AND status = 'ACTIVE'`,
		entityID,
	).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count active members: %w", err)
	}

	return count, nil
}
