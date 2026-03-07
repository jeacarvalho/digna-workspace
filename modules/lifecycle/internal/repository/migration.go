package repository

import (
	"database/sql"
	"fmt"
)

type Migrator struct{}

func NewMigrator() *Migrator {
	return &Migrator{}
}

func (m *Migrator) RunMigrations(db *sql.DB) error {
	migrations := []struct {
		name string
		sql  string
	}{
		{
			name: "create_accounts",
			sql: `CREATE TABLE IF NOT EXISTS accounts (
				id INTEGER PRIMARY KEY AUTOINCREMENT,
				code TEXT NOT NULL UNIQUE,
				name TEXT NOT NULL,
				parent_id INTEGER,
				account_type TEXT NOT NULL,
				created_at INTEGER NOT NULL,
				FOREIGN KEY (parent_id) REFERENCES accounts(id)
			)`,
		},
		{
			name: "create_entries",
			sql: `CREATE TABLE IF NOT EXISTS entries (
				id INTEGER PRIMARY KEY AUTOINCREMENT,
				entry_date INTEGER NOT NULL,
				description TEXT,
				reference TEXT,
				created_at INTEGER NOT NULL
			)`,
		},
		{
			name: "create_postings",
			sql: `CREATE TABLE IF NOT EXISTS postings (
				id INTEGER PRIMARY KEY AUTOINCREMENT,
				entry_id INTEGER NOT NULL,
				account_id INTEGER NOT NULL,
				amount INTEGER NOT NULL,
				direction TEXT NOT NULL CHECK(direction IN ('DEBIT', 'CREDIT')),
				created_at INTEGER NOT NULL,
				FOREIGN KEY (entry_id) REFERENCES entries(id),
				FOREIGN KEY (account_id) REFERENCES accounts(id)
			)`,
		},
		{
			name: "create_work_logs",
			sql: `CREATE TABLE IF NOT EXISTS work_logs (
				id INTEGER PRIMARY KEY AUTOINCREMENT,
				member_id TEXT NOT NULL,
				minutes INTEGER NOT NULL,
				activity_type TEXT NOT NULL,
				log_date INTEGER NOT NULL,
				description TEXT,
				created_at INTEGER NOT NULL
			)`,
		},
		{
			name: "create_decisions_log",
			sql: `CREATE TABLE IF NOT EXISTS decisions_log (
				id INTEGER PRIMARY KEY AUTOINCREMENT,
				title TEXT NOT NULL,
				content_hash TEXT NOT NULL,
				status TEXT NOT NULL CHECK(status IN ('DRAFT', 'APPROVED', 'REJECTED', 'ARCHIVED')),
				decision_date INTEGER,
				created_at INTEGER NOT NULL,
				updated_at INTEGER NOT NULL
			)`,
		},
		{
			name: "create_sync_metadata",
			sql: `CREATE TABLE IF NOT EXISTS sync_metadata (
				id INTEGER PRIMARY KEY CHECK(id = 1),
				last_sync_at INTEGER,
				version INTEGER NOT NULL DEFAULT 1,
				updated_at INTEGER NOT NULL
			)`,
		},
		{
			name: "init_sync_metadata",
			sql: `INSERT OR IGNORE INTO sync_metadata (id, last_sync_at, version, updated_at) 
				VALUES (1, NULL, 1, strftime('%s', 'now'))`,
		},
		{
			name: "create_indexes",
			sql: `CREATE INDEX IF NOT EXISTS idx_accounts_code ON accounts(code);
				CREATE INDEX IF NOT EXISTS idx_accounts_parent ON accounts(parent_id);
				CREATE INDEX IF NOT EXISTS idx_entries_date ON entries(entry_date);
				CREATE INDEX IF NOT EXISTS idx_postings_entry ON postings(entry_id);
				CREATE INDEX IF NOT EXISTS idx_postings_account ON postings(account_id);
				CREATE INDEX IF NOT EXISTS idx_work_logs_member ON work_logs(member_id);
				CREATE INDEX IF NOT EXISTS idx_work_logs_date ON work_logs(log_date);
				CREATE INDEX IF NOT EXISTS idx_decisions_status ON decisions_log(status);
				CREATE INDEX IF NOT EXISTS idx_decisions_date ON decisions_log(decision_date);`,
		},
		{
			name: "seed_default_accounts",
			sql: `INSERT OR IGNORE INTO accounts (id, code, name, account_type, created_at) VALUES 
				(1, '1.1.01', 'Caixa e Equivalentes', 'ASSET', strftime('%s', 'now')),
				(2, '3.1.01', 'Receita de Vendas', 'REVENUE', strftime('%s', 'now')),
				(3, '1.1.02', 'Bancos', 'ASSET', strftime('%s', 'now')),
				(4, '2.1.01', 'Fornecedores', 'LIABILITY', strftime('%s', 'now'));`,
		},
	}

	for _, migration := range migrations {
		if _, err := db.Exec(migration.sql); err != nil {
			return fmt.Errorf("migration %s failed: %w", migration.name, err)
		}
	}

	return nil
}
