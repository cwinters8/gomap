package arguments

import (
	"encoding/json"
)

type Set struct {
	AccountID string   `json:"accountId"`
	Create    *Message `json:"create"`
}

type Message struct {
	ID            string         `json:"-"` // auto generated value that can be overridden if desired
	MailboxIDs    Mailboxes      `json:"mailboxIds"`
	Keywords      *Keywords      `json:"keywords"`
	From          []*Address     `json:"from"`
	To            []*Address     `json:"to"`
	Subject       string         `json:"subject"`
	BodyStructure *BodyStructure `json:"bodyStructure"`
	BodyValue     *BodyValue     `json:"bodyValues"`
}

type Mailboxes []string

type BodyStructure struct {
	ID   string   `json:"partId"` // auto generated value that can be overridden if desired. when used in a Message, must match m.BodyValue.ID
	Type BodyType `json:"type"`
}

type BodyValue struct {
	ID    string // auto generated value that can be overridden if desired. when used in a Message, must match m.BodyStructure.ID
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
			s.Create.ID: *s.Create,
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

func (b BodyValue) MarshalJSON() ([]byte, error) {
	body := map[string]map[string]string{
		b.ID: {
			"value": b.Value,
		},
	}
	return json.Marshal(body)
}
