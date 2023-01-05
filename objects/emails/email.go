package emails

import (
	"fmt"

	"github.com/cwinters8/gomap/requests"
	"github.com/cwinters8/gomap/utils"

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
	Body       *Body
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

func (e *Email) Set(acctID string) (*requests.Call, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return nil, fmt.Errorf("failed to generate new uuid: %w", err)
	}
	reqID, err := uuid.NewRandom()
	if err != nil {
		return nil, fmt.Errorf("failed to generate new uuid: %w", err)
	}
	mailboxes := map[string]bool{}
	for _, box := range e.MailboxIDs {
		mailboxes[box] = true
	}
	return &requests.Call{
		ID:        id,
		AccountID: acctID,
		Method:    "Email/set",
		Arguments: map[string]any{
			"create": map[string]map[string]any{
				reqID.String(): {
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
				},
			},
		},
		OnSuccess: func(m map[string]any) error {
			created, ok := m["created"].(map[string]any)
			if !ok {
				return fmt.Errorf("failed to cast created value to map. %s", utils.Describe(m["created"]))
			}
			value, ok := created[reqID.String()].(map[string]any)
			if !ok {
				return fmt.Errorf("failed to cast result value to map. %s", err.Error())
			}
			resultID, ok := value["id"].(string)
			if !ok {
				return fmt.Errorf("failed to cast id to string. %s", err.Error())
			}
			e.ID = resultID
			return nil
		},
	}, nil
}
