package arguments

import (
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
)

type Set struct {
	AccountID string   `json:"accountId"`
	Create    *Message `json:"create"`
}

type Message struct {
	ID         uuid.UUID  `json:"-"` // auto generated value that can be overridden if desired
	MailboxIDs Mailboxes  `json:"mailboxIds"`
	Keywords   *Keywords  `json:"keywords"`
	From       []*Address `json:"from"`
	To         []*Address `json:"to"`
	Subject    string     `json:"subject"`
	Body       *Body      `json:"body"`
}

func NewMessage(mailboxes Mailboxes, from, to *Address, subject, msg string) (*Message, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return nil, fmt.Errorf("failed to generate new uuid: %w", err)
	}
	body, err := NewBody(TextPlain, msg)
	if err != nil {
		return nil, fmt.Errorf("failed to instantiate new body: %w", err)
	}
	return &Message{
		ID:         id,
		MailboxIDs: mailboxes,
		From:       []*Address{from},
		To:         []*Address{to},
		Subject:    subject,
		Body:       body,
	}, nil
}

type Mailboxes []string

type Body struct {
	ID        uuid.UUID
	Type      BodyType
	Value     string
	Structure *BodyStructure // do not try to set by hand - only exported for json munging
	BodyValue *BodyValue     // do not try to set by hand - only exported for json munging
}

func NewBody(t BodyType, value string) (*Body, error) {
	b := Body{
		Value: value,
		Structure: &BodyStructure{
			Type: t,
		},
		BodyValue: &BodyValue{
			Value: value,
		},
	}
	id, err := uuid.NewRandom()
	if err != nil {
		return nil, fmt.Errorf("failed to generate new uuid: %w", err)
	}
	if err := b.SetID(id); err != nil {
		return nil, fmt.Errorf("failed to update ID: %w", err)
	}
	return &b, nil
}

func (b *Body) SetID(id uuid.UUID) error {
	if id == uuid.Nil {
		return fmt.Errorf("id must not be nil")
	}
	b.ID = id
	b.Structure.ID = id
	b.BodyValue.ID = id
	return nil
}

type BodyStructure struct {
	ID   uuid.UUID `json:"partId"` // mustMatch body.BodyValue.ID
	Type BodyType  `json:"type"`
}

type BodyValue struct {
	ID    uuid.UUID // must match body.Structure.ID
	Value string
}

type BodyType string

const (
	TextPlain BodyType = "text/plain"
	TextHTML  BodyType = "text/html"
)

type Keywords struct {
	Seen  bool `json:"$seen"`
	Draft bool `json:"$draft"`
}

type Address struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

func (s Set) MarshalJSON() ([]byte, error) {
	set := map[string]any{
		"accountId": s.AccountID,
		"create": map[string]Message{
			s.Create.ID.String(): *s.Create,
		},
	}
	return json.Marshal(set)
}

func (m Mailboxes) MarshalJSON() ([]byte, error) {
	boxes := map[string]bool{}
	for _, v := range m {
		boxes[v] = true
	}
	return json.Marshal(boxes)
}

func (b Body) MarshalJSON() ([]byte, error) {
	body := struct {
		BodyStructure `json:"bodyStructure"`
		BodyValue     `json:"bodyValue"`
	}{
		BodyStructure: *b.Structure,
		BodyValue:     *b.BodyValue,
	}
	return json.Marshal(body)
}

func (b BodyValue) MarshalJSON() ([]byte, error) {
	body := map[string]map[string]string{
		b.ID.String(): {
			"value": b.Value,
		},
	}
	return json.Marshal(body)
}
