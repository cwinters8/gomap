package objects

import (
	"github.com/google/uuid"
)

type Mailbox struct {
	ID      uuid.UUID
	BoxName string
}

func (m Mailbox) GetID() uuid.UUID {
	return m.ID
}

func (m Mailbox) Name() string {
	return "Mailbox"
}

func (m Mailbox) Map() map[string]any {
	return nil
}

func (m Mailbox) Call()
