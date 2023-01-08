package emails

import (
	"fmt"

	"github.com/cwinters8/gomap/requests"
	"github.com/cwinters8/gomap/utils"

	"github.com/google/uuid"
)

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
