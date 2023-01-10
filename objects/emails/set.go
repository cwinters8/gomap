package emails

import (
	"encoding/json"
	"fmt"

	"github.com/cwinters8/gomap/client"
	"github.com/cwinters8/gomap/requests"
	"github.com/cwinters8/gomap/utils"

	"github.com/google/uuid"
)

func Set(c *client.Client, emails []*Email) (err error) {
	call, err := SetCall(c.Session.PrimaryAccounts.Mail, emails)
	if err != nil {
		return fmt.Errorf("failed to construct email set call: %w", err)
	}
	responses, err := requests.Request(c, []*requests.Call{call}, false)
	if err != nil {
		return fmt.Errorf("set request failure: %w", err)
	}
	if len(responses) < 1 {
		return fmt.Errorf("no responses returned")
	}
	requestsNotFound, err := ParseSetResponseBody(responses[0].Body, emails)
	if err != nil {
		return fmt.Errorf("failed to parse set response: %w", err)
	}
	if len(requestsNotFound) > 0 {
		return fmt.Errorf("some set requests were not found. request IDs not found: %v", requestsNotFound)
	}
	return nil
}

func SetCall(acctID string, emails []*Email) (*requests.Call, error) {
	if len(emails) < 1 {
		return nil, fmt.Errorf("no emails provided")
	}
	id, err := uuid.NewRandom()
	if err != nil {
		return nil, fmt.Errorf("failed to generate new uuid: %w", err)
	}
	create := map[string]*Email{}
	for _, e := range emails {
		create[e.RequestID.String()] = e
	}
	return &requests.Call{
		ID:        id,
		AccountID: acctID,
		Method:    "Email/set",
		Arguments: map[string]any{
			"create": create,
		},
	}, nil
}

func ParseSetResponseBody(body map[string]any, emails []*Email) (requestsNotFound []uuid.UUID, err error) {
	b, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal response body to json: %w", err)
	}
	var resp setResponse
	if err := json.Unmarshal(b, &resp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal set response: %w", err)
	}
	for _, e := range emails {
		if id, ok := resp.Created[e.RequestID.String()]; ok {
			e.ID = id.ID
		} else {
			requestsNotFound = append(requestsNotFound, e.RequestID)
		}
	}
	return requestsNotFound, nil
}

type setResponse struct {
	Created    map[string]id            `json:"created"`
	NotCreated map[string]notCreatedErr `json:"notCreated"`
}

type notCreatedErr struct {
	Description string `json:"description"`
	Type        string `json:"type"`
}

type id struct {
	ID string `json:"id"`
}

func (e *Email) Set(acctID string) (*requests.Call, error) {
	id, err := uuid.NewRandom()
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
				e.RequestID.String(): {
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
			value, ok := created[e.RequestID.String()].(map[string]any)
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
