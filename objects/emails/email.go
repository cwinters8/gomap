package emails

import (
	"fmt"

	"github.com/google/uuid"
)

type Email struct {
	ID         string
	RequestID  uuid.UUID
	MailboxIDs []string
	Keywords   *Keywords
	From       *Address
	To         []*Address
	Subject    string
	Body       *Body
}

func NewEmail(boxIDs []string, from *Address, to []*Address, subject, body string, bodyType BodyType) (Email, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return NilEmail, fmt.Errorf("failed to generate new uuid: %w", err)
	}
	bodyID, err := uuid.NewRandom()
	if err != nil {
		return NilEmail, fmt.Errorf("failed to generate new uuid for body: %w", err)
	}
	return Email{
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

func (e Email) GetReqID() uuid.UUID {
	return e.RequestID
}

func (e Email) Name() string {
	return "Email"
}

func (e Email) Map() map[uuid.UUID]map[string]any {
	boxes := map[string]bool{}
	for _, m := range e.MailboxIDs {
		boxes[m] = true
	}
	return map[uuid.UUID]map[string]any{
		e.RequestID: {
			"mailboxIds": boxes,
			"from":       []*Address{e.From},
			"to":         e.To,
			"subject":    e.Subject,
			"bodyStructure": map[string]any{
				"partId": e.Body.ID,
				"type":   e.Body.Type,
			},
			"bodyValues": map[uuid.UUID]map[string]string{
				e.Body.ID: {
					"value": e.Body.Value,
				},
			},
		},
	}
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

// empty instance of Email struct
var NilEmail = Email{}
