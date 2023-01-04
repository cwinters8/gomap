package mailboxes

import (
	"github.com/google/uuid"
)

type Mailbox struct {
	ID          string    `json:"id"`
	RequestID   uuid.UUID `json:"-"`
	BoxName     string    `json:"name"`
	TotalEmails int       `json:"totalEmails"`
}

func (m Mailbox) GetReqID() uuid.UUID {
	return m.RequestID
}

func (m Mailbox) Name() string {
	return "Mailbox"
}

func (m Mailbox) Map() map[uuid.UUID]map[string]any {
	return nil
}
