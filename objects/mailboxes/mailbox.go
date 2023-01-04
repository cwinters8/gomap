package mailboxes

import (
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
		OnSuccess: func(gotMap map[string]any) error {
			ids, ok := gotMap["ids"].([]any)
			if !ok {
				return fmt.Errorf("failed to cast ids to slice of any. %s", utils.Describe(gotMap["ids"]))
			}
			strID, ok := ids[0].(string)
			if !ok {
				return fmt.Errorf("failed to cast id to string. %s", utils.Describe(ids[0]))
			}
			m.ID = strID
			return nil
		},
	}, nil
}
