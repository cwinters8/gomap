package emails

import (
	"fmt"
	"time"

	"github.com/cwinters8/gomap/client"
	"github.com/cwinters8/gomap/parse"
	"github.com/cwinters8/gomap/requests"

	"github.com/google/uuid"
)

type Filter struct {
	InMailboxID string     `json:"inMailbox,omitempty"`
	Text        string     `json:"text,omitempty"`   // searches From, To, Cc, Bcc, and Subject header fields and any text/* body parts
	Before      *time.Time `json:"before,omitempty"` // UTC timestamp the email's receivedAt must be before
	After       *time.Time `json:"after,omitempty"`  // UTC timestamp the email's receivedAt must match or be after
}

func Query(c *client.Client, filter *Filter) (emailIDs []string, err error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return nil, fmt.Errorf("failed to generate new uuid: %w", err)
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
