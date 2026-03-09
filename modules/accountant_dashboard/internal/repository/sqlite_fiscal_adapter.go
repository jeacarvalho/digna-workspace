package repository

import (
	"context"
	"database/sql"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"digna/accountant_dashboard/internal/domain"
	_ "github.com/mattn/go-sqlite3"
)

const DataDir = "../../data/entities"

type SQLiteFiscalAdapter struct {
	basePath string
}

func NewSQLiteFiscalAdapter() *SQLiteFiscalAdapter {
	return &SQLiteFiscalAdapter{
		basePath: DataDir,
	}
}

func (r *SQLiteFiscalAdapter) BasePath() string {
	return r.basePath
}

func (r *SQLiteFiscalAdapter) SetBasePath(path string) {
	r.basePath = path
}

func (r *SQLiteFiscalAdapter) openReadOnly(entityID string) (*sql.DB, error) {
	dbPath := filepath.Join(r.basePath, fmt.Sprintf("%s.db", entityID))

	dsn := fmt.Sprintf("file:%s?mode=ro", dbPath)
	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database in read-only mode: %w", err)
	}

	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return db, nil
}

func (r *SQLiteFiscalAdapter) LoadEntries(ctx context.Context, entityID string, period string) ([]domain.EntryDTO, error) {
	db, err := r.openReadOnly(entityID)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	yearMonth := strings.Split(period, "-")
	if len(yearMonth) != 2 {
		return nil, fmt.Errorf("invalid period format: expected YYYY-MM, got %s", period)
	}

	year := yearMonth[0]
	month := yearMonth[1]
	prefix := fmt.Sprintf("%s-%s", year, month)

	query := `
		SELECT 
			e.id, e.entry_date, e.description, e.reference,
			p.id as posting_id, p.account_id, a.code as account_code, a.name as account_name,
			p.amount, p.direction
		FROM entries e
		LEFT JOIN postings p ON e.id = p.entry_id
		LEFT JOIN accounts a ON p.account_id = a.id
		WHERE strftime('%Y-%m', e.entry_date, 'unixepoch') = ?
		ORDER BY e.id, p.id
	`

	rows, err := db.QueryContext(ctx, query, prefix)
	if err != nil {
		return nil, fmt.Errorf("failed to query entries: %w", err)
	}
	defer rows.Close()

	entriesMap := make(map[int64]*domain.EntryDTO)

	for rows.Next() {
		var entryID int64
		var entryDate int64
		var description, reference string
		var postingID, accountID sql.NullInt64
		var accountCode, accountName sql.NullString
		var amount int64
		var direction string

		err := rows.Scan(
			&entryID, &entryDate, &description, &reference,
			&postingID, &accountID, &accountCode, &accountName,
			&amount, &direction,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		entry, exists := entriesMap[entryID]
		if !exists {
			entry = &domain.EntryDTO{
				ID:          entryID,
				EntityID:    entityID,
				Date:        time.Unix(entryDate, 0),
				Description: description,
				Postings:    []domain.PostingDTO{},
			}
			entriesMap[entryID] = entry
		}

		if postingID.Valid {
			posting := domain.PostingDTO{
				ID:          postingID.Int64,
				EntryID:     entryID,
				AccountID:   accountID.Int64,
				AccountCode: accountCode.String,
				AccountName: accountName.String,
				Debit:       0,
				Credit:      0,
			}

			if direction == "DEBIT" {
				posting.Debit = amount
				entry.TotalDebit += amount
			} else {
				posting.Credit = amount
				entry.TotalCredit += amount
			}

			entry.Postings = append(entry.Postings, posting)
		}
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	result := make([]domain.EntryDTO, 0, len(entriesMap))
	for _, entry := range entriesMap {
		result = append(result, *entry)
	}

	return result, nil
}

func (r *SQLiteFiscalAdapter) RegisterExport(ctx context.Context, entityID string, batch *domain.FiscalBatch) error {
	dbPath := filepath.Join(r.basePath, fmt.Sprintf("%s.db", entityID))

	dsn := fmt.Sprintf("%s?_journal_mode=WAL&_foreign_keys=ON", dbPath)
	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}
	defer db.Close()

	migrationSQL := `
		CREATE TABLE IF NOT EXISTS fiscal_exports (
			id TEXT PRIMARY KEY,
			entity_id TEXT NOT NULL,
			period TEXT NOT NULL,
			batch_id TEXT NOT NULL,
			export_hash TEXT NOT NULL,
			total_entries INTEGER NOT NULL,
			created_at INTEGER NOT NULL,
			UNIQUE(entity_id, period)
		)
	`
	if _, err := db.ExecContext(ctx, migrationSQL); err != nil {
		return fmt.Errorf("failed to create fiscal_exports table: %w", err)
	}

	query := `
		INSERT OR REPLACE INTO fiscal_exports 
		(id, entity_id, period, batch_id, export_hash, total_entries, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`

	_, err = db.ExecContext(ctx, query,
		batch.ID,
		batch.EntityID,
		batch.Period,
		batch.ID,
		batch.ExportHash,
		batch.TotalEntries,
		batch.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to register export: %w", err)
	}

	return nil
}

func (r *SQLiteFiscalAdapter) ListPendingEntities(ctx context.Context, period string) ([]string, error) {
	entries, err := filepath.Glob(filepath.Join(r.basePath, "*.db"))
	if err != nil {
		return nil, fmt.Errorf("failed to glob database files: %w", err)
	}

	var pending []string

	for _, entry := range entries {
		entityID := strings.TrimSuffix(filepath.Base(entry), ".db")

		db, err := r.openReadOnly(entityID)
		if err != nil {
			continue
		}

		yearMonth := strings.Split(period, "-")
		if len(yearMonth) != 2 {
			continue
		}
		prefix := fmt.Sprintf("%s-%s", yearMonth[0], yearMonth[1])

		var count int
		err = db.QueryRowContext(ctx,
			"SELECT COUNT(*) FROM entries WHERE strftime('%Y-%m', entry_date, 'unixepoch') = ?", prefix).Scan(&count)
		db.Close()

		if err != nil {
			continue
		}

		if count > 0 {
			var exportedCount int
			db2, err := r.openReadOnly(entityID)
			if err == nil {
				db2.QueryRowContext(ctx,
					"SELECT COUNT(*) FROM fiscal_exports WHERE period = ?", period).Scan(&exportedCount)
				db2.Close()
			}

			if exportedCount == 0 {
				pending = append(pending, entityID)
			}
		}
	}

	return pending, nil
}

func (r *SQLiteFiscalAdapter) GetExportHistory(ctx context.Context, entityID string, period string) ([]domain.FiscalExportLog, error) {
	db, err := r.openReadOnly(entityID)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	query := `
		SELECT id, entity_id, period, batch_id, export_hash, created_at
		FROM fiscal_exports
		WHERE entity_id = ? AND period = ?
		ORDER BY created_at DESC
	`

	rows, err := db.QueryContext(ctx, query, entityID, period)
	if err != nil {
		return nil, fmt.Errorf("failed to query export history: %w", err)
	}
	defer rows.Close()

	var logs []domain.FiscalExportLog

	for rows.Next() {
		var log domain.FiscalExportLog
		err := rows.Scan(
			&log.ID,
			&log.EntityID,
			&log.Period,
			&log.BatchID,
			&log.ExportHash,
			&log.ExportedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan export log: %w", err)
		}
		logs = append(logs, log)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return logs, nil
}
