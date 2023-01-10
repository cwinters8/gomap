package emails

import (
	"fmt"

	"github.com/cwinters8/gomap/client"
	"github.com/cwinters8/gomap/parse"
	"github.com/cwinters8/gomap/requests"

	"github.com/google/uuid"
)

func Query(c *client.Client, text string, inMailboxID string) (emailIDs []string, err error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return nil, fmt.Errorf("failed to generate new uuid: %w", err)
	}
	filter := map[string]string{"text": text}
	if len(inMailboxID) > 0 {
		filter["inMailbox"] = inMailboxID
	}
	call := requests.Call{
		ID:        id,
		AccountID: c.Session.PrimaryAccounts.Mail,
		Method:    "Email/query",
		Arguments: map[string]any{
			"filter": filter,
			"sort": []map[string]any{{
				"isAscending": false,
				"property":    "receivedAt",
			}},
		},
	}
	responses, err := requests.Request(c, []*requests.Call{&call}, false)
	if err != nil {
		return nil, fmt.Errorf("query request failure: %w", err)
	}
	if len(responses) < 1 {
		return nil, fmt.Errorf("no responses returned from request")
	}
	return parse.QueryResponseBody(responses[0].Body)
}
