package domain

import (
	"database/sql"
)

type LifecycleManager interface {
	GetConnection(entityID string) (*sql.DB, error)
	CloseConnection(entityID string) error
	CloseAll() error
}

type Migrator interface {
	RunMigrations(db *sql.DB) error
}
