package mailboxes

import (
	"encoding/json"
	"fmt"

	"github.com/cwinters8/gomap/client"
	"github.com/cwinters8/gomap/requests"

	"github.com/google/uuid"
)

type Mailbox struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func GetMailboxByName(c *client.Client, name string) (*Mailbox, error) {
	m := Mailbox{
		Name: name,
	}
	call, err := m.Query(c.Session.PrimaryAccounts.Mail)
	if err != nil {
		return nil, fmt.Errorf("failed to construct mailbox query call")
	}
	responses, err := requests.Request(c, []*requests.Call{call}, false)
	if err != nil {
		return nil, fmt.Errorf("query request failure: %w", err)
	}
	ids, err := ParseQueryResponseBody(responses[0].Body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse query response: %w", err)
	}
	m.ID = ids[0]
	return &m, nil
}

func ParseQueryResponseBody(body map[string]any) (ids []string, err error) {
	b, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal body to json: %w", err)
	}
	var resp queryResponse
	if err := json.Unmarshal(b, &resp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal query response: %w", err)
	}
	return resp.IDs, nil
}

type queryResponse struct {
	IDs []string `json:"ids"`
}

func (m *Mailbox) Query(acctID string) (*requests.Call, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return nil, fmt.Errorf("failed to generate new uuid: %w", err)
	}
	return &requests.Call{
		ID:        id,
		AccountID: acctID,
		Method:    "Mailbox/query",
		Arguments: map[string]any{
			"filter": map[string]string{
				"name": m.Name,
			},
		},
		OnSuccess: func(gotMap map[string]any) error {
			ids, err := ParseQueryResponseBody(gotMap)
			if err != nil {
				return fmt.Errorf("failed to parse query response: %w", err)
			}
			m.ID = ids[0]
			return nil
		},
	}, nil
}
