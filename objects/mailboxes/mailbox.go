package mailboxes

import (
	"encoding/json"
	"fmt"

	"github.com/cwinters8/gomap/requests"
	"github.com/cwinters8/gomap/utils"

	"github.com/google/uuid"
)

type Mailbox struct {
	ID   string `json:"id"`
	Name string `json:"name"`
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
		OnSuccess: func(b []byte) error {
			var gotMap map[string]any
			if err := json.Unmarshal(b, &gotMap); err != nil {
				return fmt.Errorf("failed to unmarshal response body to map: %w", err)
			}
			ids, ok := gotMap["ids"].([]string)
			if !ok {
				return fmt.Errorf("failed to cast ids to slice of string. %s", utils.Describe(gotMap["ids"]))
			}
			m.ID = ids[0]
			return nil
		},
	}, nil
}
