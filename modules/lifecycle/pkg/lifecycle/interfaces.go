package lifecycle

import (
	"database/sql"
)

type LifecycleManager interface {
	GetConnection(entityID string) (*sql.DB, error)
	CloseConnection(entityID string) error
	CloseAll() error
	EntityExists(entityID string) (bool, error)
	CreateEntity(entityID, entityName string) error
}
