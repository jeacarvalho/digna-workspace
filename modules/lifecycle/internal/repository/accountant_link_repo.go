package repository

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/providentia/digna/lifecycle/internal/domain"
)

type EnterpriseAccountantRepository interface {
	Create(link *domain.EnterpriseAccountant) error
	Update(link *domain.EnterpriseAccountant) error
	FindByID(id string) (*domain.EnterpriseAccountant, error)
	FindByEnterpriseID(enterpriseID string) ([]*domain.EnterpriseAccountant, error)
	FindByAccountantID(accountantID string) ([]*domain.EnterpriseAccountant, error)
	FindActiveByEnterpriseID(enterpriseID string) (*domain.EnterpriseAccountant, error)
	FindActiveByAccountantID(accountantID string) ([]*domain.EnterpriseAccountant, error)
	FindByDateRange(enterpriseID, accountantID string, startDate, endDate time.Time) ([]*domain.EnterpriseAccountant, error)
	// New methods for temporal filtering
	FindByAccountantIDAndDateRange(accountantID string, startTime, endTime int64) ([]*domain.EnterpriseAccountant, error)
	FindByAccountantAndEnterpriseInDateRange(accountantID, enterpriseID string, startTime, endTime int64) ([]*domain.EnterpriseAccountant, error)
	FindByAccountantAndEnterprise(accountantID, enterpriseID string) ([]*domain.EnterpriseAccountant, error)
}

type SQLiteEnterpriseAccountantRepository struct {
	db *sql.DB
}

func NewSQLiteEnterpriseAccountantRepository(db *sql.DB) *SQLiteEnterpriseAccountantRepository {
	return &SQLiteEnterpriseAccountantRepository{db: db}
}

func (r *SQLiteEnterpriseAccountantRepository) Create(link *domain.EnterpriseAccountant) error {
	query := `INSERT INTO enterprise_accountants 
		(enterprise_id, accountant_id, status, start_date, end_date, delegated_by, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)`

	startDateUnix := link.StartDate.Unix()
	var endDateUnix interface{}
	if link.EndDate != nil {
		endDateUnix = link.EndDate.Unix()
	} else {
		endDateUnix = nil
	}
	createdAtUnix := link.CreatedAt.Unix()
	updatedAtUnix := link.UpdatedAt.Unix()

	result, err := r.db.Exec(query,
		link.EnterpriseID,
		link.AccountantID,
		string(link.Status),
		startDateUnix,
		endDateUnix,
		link.DelegatedBy,
		createdAtUnix,
		updatedAtUnix,
	)
	if err != nil {
		return fmt.Errorf("failed to create enterprise accountant link: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert id: %w", err)
	}

	link.ID = fmt.Sprintf("%d", id)
	return nil
}

func (r *SQLiteEnterpriseAccountantRepository) Update(link *domain.EnterpriseAccountant) error {
	query := `UPDATE enterprise_accountants 
		SET status = ?, start_date = ?, end_date = ?, updated_at = ?
		WHERE id = ?`

	startDateUnix := link.StartDate.Unix()
	var endDateUnix interface{}
	if link.EndDate != nil {
		endDateUnix = link.EndDate.Unix()
	} else {
		endDateUnix = nil
	}
	updatedAtUnix := link.UpdatedAt.Unix()

	_, err := r.db.Exec(query,
		string(link.Status),
		startDateUnix,
		endDateUnix,
		updatedAtUnix,
		link.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update enterprise accountant link: %w", err)
	}

	return nil
}

func (r *SQLiteEnterpriseAccountantRepository) FindByID(id string) (*domain.EnterpriseAccountant, error) {
	query := `SELECT id, enterprise_id, accountant_id, status, start_date, end_date, delegated_by, created_at, updated_at
		FROM enterprise_accountants WHERE id = ?`

	row := r.db.QueryRow(query, id)
	return r.scanRow(row)
}

func (r *SQLiteEnterpriseAccountantRepository) FindByEnterpriseID(enterpriseID string) ([]*domain.EnterpriseAccountant, error) {
	query := `SELECT id, enterprise_id, accountant_id, status, start_date, end_date, delegated_by, created_at, updated_at
		FROM enterprise_accountants WHERE enterprise_id = ? ORDER BY start_date DESC`

	rows, err := r.db.Query(query, enterpriseID)
	if err != nil {
		return nil, fmt.Errorf("failed to query by enterprise id: %w", err)
	}
	defer rows.Close()

	return r.scanRows(rows)
}

func (r *SQLiteEnterpriseAccountantRepository) FindByAccountantID(accountantID string) ([]*domain.EnterpriseAccountant, error) {
	query := `SELECT id, enterprise_id, accountant_id, status, start_date, end_date, delegated_by, created_at, updated_at
		FROM enterprise_accountants WHERE accountant_id = ? ORDER BY start_date DESC`

	rows, err := r.db.Query(query, accountantID)
	if err != nil {
		return nil, fmt.Errorf("failed to query by accountant id: %w", err)
	}
	defer rows.Close()

	return r.scanRows(rows)
}

func (r *SQLiteEnterpriseAccountantRepository) FindActiveByEnterpriseID(enterpriseID string) (*domain.EnterpriseAccountant, error) {
	query := `SELECT id, enterprise_id, accountant_id, status, start_date, end_date, delegated_by, created_at, updated_at
		FROM enterprise_accountants WHERE enterprise_id = ? AND status = 'ACTIVE' ORDER BY start_date DESC LIMIT 1`

	row := r.db.QueryRow(query, enterpriseID)
	return r.scanRow(row)
}

func (r *SQLiteEnterpriseAccountantRepository) FindActiveByAccountantID(accountantID string) ([]*domain.EnterpriseAccountant, error) {
	query := `SELECT id, enterprise_id, accountant_id, status, start_date, end_date, delegated_by, created_at, updated_at
		FROM enterprise_accountants WHERE accountant_id = ? AND status = 'ACTIVE' ORDER BY start_date DESC`

	rows, err := r.db.Query(query, accountantID)
	if err != nil {
		return nil, fmt.Errorf("failed to query active by accountant id: %w", err)
	}
	defer rows.Close()

	return r.scanRows(rows)
}

func (r *SQLiteEnterpriseAccountantRepository) FindByDateRange(enterpriseID, accountantID string, startDate, endDate time.Time) ([]*domain.EnterpriseAccountant, error) {
	query := `SELECT id, enterprise_id, accountant_id, status, start_date, end_date, delegated_by, created_at, updated_at
		FROM enterprise_accountants 
		WHERE enterprise_id = ? AND accountant_id = ? 
		AND (
			(end_date IS NULL AND start_date <= ?) OR
			(end_date IS NOT NULL AND start_date <= ? AND end_date >= ?)
		)
		ORDER BY start_date DESC`

	startDateUnix := endDate.Unix()
	endDateUnix := startDate.Unix()

	rows, err := r.db.Query(query, enterpriseID, accountantID, startDateUnix, endDateUnix, startDateUnix)
	if err != nil {
		return nil, fmt.Errorf("failed to query by date range: %w", err)
	}
	defer rows.Close()

	return r.scanRows(rows)
}

func (r *SQLiteEnterpriseAccountantRepository) scanRow(row *sql.Row) (*domain.EnterpriseAccountant, error) {
	var link domain.EnterpriseAccountant
	var statusStr string
	var startDateUnix, createdAtUnix, updatedAtUnix int64
	var endDateUnix sql.NullInt64

	err := row.Scan(
		&link.ID,
		&link.EnterpriseID,
		&link.AccountantID,
		&statusStr,
		&startDateUnix,
		&endDateUnix,
		&link.DelegatedBy,
		&createdAtUnix,
		&updatedAtUnix,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to scan row: %w", err)
	}

	link.Status = domain.AccountantStatus(statusStr)
	link.StartDate = time.Unix(startDateUnix, 0).UTC()
	if endDateUnix.Valid {
		endDate := time.Unix(endDateUnix.Int64, 0).UTC()
		link.EndDate = &endDate
	}
	link.CreatedAt = time.Unix(createdAtUnix, 0).UTC()
	link.UpdatedAt = time.Unix(updatedAtUnix, 0).UTC()

	return &link, nil
}

func (r *SQLiteEnterpriseAccountantRepository) scanRows(rows *sql.Rows) ([]*domain.EnterpriseAccountant, error) {
	var links []*domain.EnterpriseAccountant

	for rows.Next() {
		var link domain.EnterpriseAccountant
		var statusStr string
		var startDateUnix, createdAtUnix, updatedAtUnix int64
		var endDateUnix sql.NullInt64

		err := rows.Scan(
			&link.ID,
			&link.EnterpriseID,
			&link.AccountantID,
			&statusStr,
			&startDateUnix,
			&endDateUnix,
			&link.DelegatedBy,
			&createdAtUnix,
			&updatedAtUnix,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		link.Status = domain.AccountantStatus(statusStr)
		link.StartDate = time.Unix(startDateUnix, 0).UTC()
		if endDateUnix.Valid {
			endDate := time.Unix(endDateUnix.Int64, 0).UTC()
			link.EndDate = &endDate
		}
		link.CreatedAt = time.Unix(createdAtUnix, 0).UTC()
		link.UpdatedAt = time.Unix(updatedAtUnix, 0).UTC()

		links = append(links, &link)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return links, nil
}

// FindByAccountantIDAndDateRange finds all links for an accountant within a date range
func (r *SQLiteEnterpriseAccountantRepository) FindByAccountantIDAndDateRange(accountantID string, startTime, endTime int64) ([]*domain.EnterpriseAccountant, error) {
	query := `
		SELECT id, enterprise_id, accountant_id, status, start_date, end_date, delegated_by, created_at, updated_at
		FROM enterprise_accountants
		WHERE accountant_id = ?
		AND status = 'ACTIVE'
		AND (
			(start_date <= ? AND (end_date IS NULL OR end_date >= ?))
			OR (start_date <= ? AND end_date IS NULL)
		)
		ORDER BY start_date DESC
	`

	rows, err := r.db.Query(query, accountantID, endTime, startTime, endTime)
	if err != nil {
		return nil, fmt.Errorf("failed to query links: %w", err)
	}
	defer rows.Close()

	return r.scanRows(rows)
}

// FindByAccountantAndEnterpriseInDateRange finds links for a specific accountant-enterprise pair within a date range
func (r *SQLiteEnterpriseAccountantRepository) FindByAccountantAndEnterpriseInDateRange(accountantID, enterpriseID string, startTime, endTime int64) ([]*domain.EnterpriseAccountant, error) {
	query := `
		SELECT id, enterprise_id, accountant_id, status, start_date, end_date, delegated_by, created_at, updated_at
		FROM enterprise_accountants
		WHERE accountant_id = ?
		AND enterprise_id = ?
		AND status = 'ACTIVE'
		AND (
			(start_date <= ? AND (end_date IS NULL OR end_date >= ?))
			OR (start_date <= ? AND end_date IS NULL)
		)
		ORDER BY start_date DESC
	`

	rows, err := r.db.Query(query, accountantID, enterpriseID, endTime, startTime, endTime)
	if err != nil {
		return nil, fmt.Errorf("failed to query links: %w", err)
	}
	defer rows.Close()

	return r.scanRows(rows)
}

// FindByAccountantAndEnterprise finds all links between an accountant and enterprise
func (r *SQLiteEnterpriseAccountantRepository) FindByAccountantAndEnterprise(accountantID, enterpriseID string) ([]*domain.EnterpriseAccountant, error) {
	query := `
		SELECT id, enterprise_id, accountant_id, status, start_date, end_date, delegated_by, created_at, updated_at
		FROM enterprise_accountants
		WHERE accountant_id = ?
		AND enterprise_id = ?
		ORDER BY start_date DESC
	`

	rows, err := r.db.Query(query, accountantID, enterpriseID)
	if err != nil {
		return nil, fmt.Errorf("failed to query links: %w", err)
	}
	defer rows.Close()

	return r.scanRows(rows)
}
