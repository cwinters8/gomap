package objects

import "github.com/google/uuid"

type Object interface {
	GetID() uuid.UUID
	Name() string
	Map() map[string]any
}
