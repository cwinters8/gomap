package methods

import (
	"encoding/json"
	"fmt"

	"github.com/cwinters8/gomap/utils"

	"github.com/google/uuid"
)

type Set struct {
	AccountID string   `json:"accountId"`
	Create    *Message `json:"create"`
}

type Message struct {
	ID         uuid.UUID  `json:"-"`
	MailboxIDs *Mailboxes `json:"mailboxIds"`
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
		MailboxIDs: &mailboxes,
		Keywords: &Keywords{
			Seen:  true,
			Draft: true,
		},
		From:    []*Address{from},
		To:      []*Address{to},
		Subject: subject,
		Body:    body,
	}, nil
}

type Mailboxes []string

type Body struct {
	ID        uuid.UUID
	Type      BodyType
	Value     string
	structure *bodyStructure
	value     *bodyValue
}

func NewBody(t BodyType, value string) (*Body, error) {
	b := Body{
		Value: value,
		Type:  t,
		structure: &bodyStructure{
			Type: t,
		},
		value: &bodyValue{
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
	b.structure.ID = id
	b.value.ID = id
	return nil
}

type bodyStructure struct {
	ID   uuid.UUID `json:"partId"` // must match body.value.ID
	Type BodyType  `json:"type"`
}

type bodyValue struct {
	ID    uuid.UUID // must match body.structure.ID
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

func (s *Set) UnmarshalJSON(b []byte) error {
	var set map[string]any
	if err := json.Unmarshal(b, &set); err != nil {
		return fmt.Errorf("failed to unmarshal set to map: %w", err)
	}
	acctID, ok := set["accountId"].(string)
	if !ok {
		return fmt.Errorf("failed to coerce account id to string. %s", utils.Describe(set["accountId"]))
	}
	s.AccountID = acctID
	create, ok := set["create"].(map[string]any)
	if !ok {
		return fmt.Errorf("failed to coerce create value to map. %s", utils.Describe(set["create"]))
	}
	var msgID uuid.UUID
	for k := range create {
		id, err := uuid.Parse(k)
		if err != nil {
			return fmt.Errorf("failed to parse uuid from key: %w", err)
		}
		msgID = id
		break // no need to try to continue iterating as only 1 message is currently supported
	}
	if msgID == uuid.Nil {
		return fmt.Errorf("message ID not found in map keys")
	}
	msg := create[msgID.String()]
	b, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal raw message to json: %w", err)
	}
	var m Message
	if err := json.Unmarshal(b, &m); err != nil {
		return fmt.Errorf("failed to unmarshal message: %w", err)
	}
	m.ID = msgID
	s.Create = &m
	if s.Create == nil {
		return fmt.Errorf("nil Create message found")
	}
	return nil
}

func (m Message) MarshalJSON() ([]byte, error) {
	msg := map[string]any{
		"mailboxIds":    m.MailboxIDs,
		"keywords":      m.Keywords,
		"from":          m.From,
		"to":            m.To,
		"subject":       m.Subject,
		"bodyStructure": m.Body.structure,
		"bodyValues":    m.Body.value,
	}
	return json.Marshal(msg)
}

func (m *Message) UnmarshalJSON(b []byte) error {
	var msg struct {
		Mailboxes     *Mailboxes     `json:"mailboxIds"`
		Keywords      *Keywords      `json:"keywords"`
		From          []*Address     `json:"from"`
		To            []*Address     `json:"to"`
		Subject       string         `json:"subject"`
		BodyStructure *bodyStructure `json:"bodyStructure"`
		BodyValue     *bodyValue     `json:"bodyValues"`
	}
	if err := json.Unmarshal(b, &msg); err != nil {
		return fmt.Errorf("failed to unmarshal message from json: %w", err)
	}
	m.MailboxIDs = msg.Mailboxes
	m.Keywords = msg.Keywords
	m.From = msg.From
	m.To = msg.To
	m.Subject = msg.Subject
	m.Body = &Body{
		ID:        msg.BodyStructure.ID,
		Type:      msg.BodyStructure.Type,
		Value:     msg.BodyValue.Value,
		structure: msg.BodyStructure,
		value:     msg.BodyValue,
	}
	return nil
}

func (m Mailboxes) MarshalJSON() ([]byte, error) {
	boxes := map[string]bool{}
	for _, v := range m {
		boxes[v] = true
	}
	return json.Marshal(boxes)
}

func (m *Mailboxes) UnmarshalJSON(b []byte) error {
	var boxes map[string]bool
	if err := json.Unmarshal(b, &boxes); err != nil {
		return fmt.Errorf("failed to unmarshal mailboxes: %w", err)
	}
	for k, v := range boxes {
		if v {
			*m = append(*m, k)
		}
	}
	return nil
}

func (b bodyValue) MarshalJSON() ([]byte, error) {
	body := map[string]map[string]string{
		b.ID.String(): {
			"value": b.Value,
		},
	}
	return json.Marshal(body)
}

func (bv *bodyValue) UnmarshalJSON(b []byte) error {
	var body map[string]map[string]string
	if err := json.Unmarshal(b, &body); err != nil {
		return fmt.Errorf("failed to unmarshal body value: %w", err)
	}
	for k, v := range body {
		id, err := uuid.Parse(k)
		if err != nil {
			return fmt.Errorf("failed to parse uuid from key: %w", err)
		}
		bv.ID = id
		bv.Value = v["value"]
	}
	return nil
}
