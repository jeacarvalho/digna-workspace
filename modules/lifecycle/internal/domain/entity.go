package domain

import (
	"time"
)

type EntityStatus string

const (
	StatusDream      EntityStatus = "DREAM"
	StatusFormalized EntityStatus = "FORMALIZED"
)

type Entity struct {
	ID        string       `json:"id"`
	Name      string       `json:"name"`
	Status    EntityStatus `json:"status"`
	CreatedAt time.Time    `json:"created_at"`
}

func NewEntity(id, name string, status EntityStatus) *Entity {
	return &Entity{
		ID:        id,
		Name:      name,
		Status:    status,
		CreatedAt: time.Now().UTC(),
	}
}

func (e *Entity) IsFormalized() bool {
	return e.Status == StatusFormalized
}
