package emails

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/cwinters8/gomap/utils"
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

func (e *Email) UnmarshalJSON(b []byte) error {
	var m map[string]any
	if err := json.Unmarshal(b, &m); err != nil {
		return fmt.Errorf("failed to unmarshal email: %w", err)
	}
	id, ok := m["id"].(string)
	if !ok {
		return fmt.Errorf("failed to cast id to string. %s", utils.Describe(m["id"]))
	}
	e.ID = id
	subject, ok := m["subject"].(string)
	if !ok {
		return fmt.Errorf("failed to cast subject to string. %s", utils.Describe(m["subject"]))
	}
	e.Subject = subject
	boxIDs, ok := m["mailboxIds"].(map[string]bool)
	if !ok {
		return fmt.Errorf("failed to cast mailbox IDs to map. %s", utils.Describe(m["mailboxIds"]))
	}
	boxes := []string{}
	for k, v := range boxIDs {
		if v {
			boxes = append(boxes, k)
		}
	}
	e.MailboxIDs = boxes
	rawTo, ok := m["to"].([]map[string]string)
	if !ok {
		return fmt.Errorf("failed to cast `to` to slice of maps. %s", utils.Describe(m["to"]))
	}
	to := []*Address{}
	for _, addr := range rawTo {
		to = append(to, parseAddress(addr))
	}
	e.To = to
	rawFrom, ok := m["from"].([]map[string]string)
	if !ok {
		return fmt.Errorf("failed to cast `from` to slice of maps. %s", utils.Describe(m["from"]))
	}
	e.From = parseAddress(rawFrom[0])
	rawBody, ok := m["bodyValues"].(map[string]map[string]any)
	if !ok {
		return fmt.Errorf("failed to cast body to map of maps. %s", utils.Describe(m["bodyValues"]))
	}
	values := []string{}
	for k, v := range rawBody {
		key, err := strconv.Atoi(k)
		if err != nil {
			return fmt.Errorf("failed to convert `%s` to int: %w", k, err)
		}
		val, ok := v["value"].(string)
		if !ok {
			return fmt.Errorf("failed to cast value to string. %s", utils.Describe(v["value"]))
		}
		values[key-1] = val
	}
	e.Body.Value = strings.Join(values, " ")
	return nil
}

func parseAddress(m map[string]string) *Address {
	a := Address{}
	for k, v := range m {
		switch k {
		case "name":
			a.Name = v
		case "email":
			a.Email = v
		}
	}
	return &a
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
