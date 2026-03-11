package repository

import (
	"database/sql"
	"fmt"
)

type CentralMigrator struct{}

func NewCentralMigrator() *CentralMigrator {
	return &CentralMigrator{}
}

func (m *CentralMigrator) RunMigrations(db *sql.DB) error {
	migrations := []struct {
		name string
		sql  string
	}{
		{
			name: "create_enterprise_accountants",
			sql: `CREATE TABLE IF NOT EXISTS enterprise_accountants (
				id INTEGER PRIMARY KEY AUTOINCREMENT,
				enterprise_id TEXT NOT NULL,
				accountant_id TEXT NOT NULL,
				status TEXT NOT NULL CHECK(status IN ('ACTIVE', 'INACTIVE')),
				start_date INTEGER NOT NULL,
				end_date INTEGER,
				delegated_by TEXT NOT NULL,
				created_at INTEGER NOT NULL,
				updated_at INTEGER NOT NULL,
				UNIQUE(enterprise_id, accountant_id)
			)`,
		},
		{
			name: "create_enterprise_accountants_indexes",
			sql: `CREATE INDEX IF NOT EXISTS idx_enterprise_accountants_enterprise ON enterprise_accountants(enterprise_id, status);
				CREATE INDEX IF NOT EXISTS idx_enterprise_accountants_accountant ON enterprise_accountants(accountant_id, status);
				CREATE INDEX IF NOT EXISTS idx_enterprise_accountants_dates ON enterprise_accountants(start_date, end_date);`,
		},
	}

	for _, migration := range migrations {
		if _, err := db.Exec(migration.sql); err != nil {
			return fmt.Errorf("central migration %s failed: %w", migration.name, err)
		}
	}

	return nil
}
