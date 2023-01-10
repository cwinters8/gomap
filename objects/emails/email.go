package emails

import (
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
)

type Email struct {
	ID         string     `json:"id"`
	RequestID  uuid.UUID  `json:"-"`
	MailboxIDs []string   `json:"mailboxIds"`
	Keywords   *Keywords  `json:"keywords"`
	From       []*Address `json:"from"`
	To         []*Address `json:"to"`
	Subject    string     `json:"subject"`
	Body       *Body      `json:"-"`
}
type Address struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}
type Body struct {
	ID    uuid.UUID
	Type  BodyType
	Value string
}
type Keywords struct {
	Seen  bool `json:"$seen"`
	Draft bool `json:"$draft"`
}
type BodyType string

const (
	TextPlain BodyType = "text/plain"
	TextHTML  BodyType = "text/html"
)

func NewEmail(boxIDs []string, from []*Address, to []*Address, subject, body string, bodyType BodyType) (*Email, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return nil, fmt.Errorf("failed to generate new uuid: %w", err)
	}
	bodyID, err := uuid.NewRandom()
	if err != nil {
		return nil, fmt.Errorf("failed to generate new uuid for body: %w", err)
	}
	return &Email{
		RequestID:  id,
		MailboxIDs: boxIDs,
		Keywords: &Keywords{
			Seen:  true,
			Draft: true,
		},
		From:    from,
		To:      to,
		Subject: subject,
		Body: &Body{
			ID:    bodyID,
			Type:  bodyType,
			Value: body,
		},
	}, nil
}

func (e Email) MarshalJSON() ([]byte, error) {
	mailboxes := map[string]bool{}
	for _, box := range e.MailboxIDs {
		mailboxes[box] = true
	}
	raw := map[string]any{
		"mailboxIds": mailboxes,
		"from":       e.From,
		"to":         e.To,
		"subject":    e.Subject,
		"bodyStructure": map[string]string{
			"partId": e.Body.ID.String(),
			"type":   string(e.Body.Type),
		},
		"bodyValues": map[string]map[string]string{
			e.Body.ID.String(): {
				"value": e.Body.Value,
			},
		},
	}
	return json.Marshal(raw)
}
