package gomap

import (
	"encoding/json"
	"fmt"

	"github.com/cwinters8/gomap/arguments"
	"github.com/cwinters8/gomap/utils"
)

type Mailbox struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	client *Client
}

func (c *Client) NewMailbox() *Mailbox {
	return &Mailbox{
		client: c,
	}
}

// queries for the mailbox and populates m.ID if found
//
// requires m.Name to be populated
func (m *Mailbox) Query() error {
	if len(m.Name) < 1 {
		return fmt.Errorf("mailbox Name must be non-empty")
	}

	i := Invocation[arguments.Query]{
		Method: &Method[arguments.Query]{
			Args: arguments.Query{
				AccountID: m.client.Session.PrimaryAccounts.Mail,
				Filter: arguments.Filter{
					Name: m.Name,
				},
			},
		},
	}

	_, err := json.Marshal(i)
	if err != nil {
		return fmt.Errorf("failed to marshal args to json: %w", err)
	}
	// status, body, err := m.client.httpRequest(http.MethodPost, m.client.Session.APIURL, b)
	// if err != nil {
	// 	return fmt.Errorf("request status %d\nfailed to make query request: %w", status, err)
	// }

	// TODO: unmarshal body to mailbox
	return fmt.Errorf("not implemented")
}

func (m Mailbox) MarshalJSON() ([]byte, error) {

	return nil, utils.ErrNotImplemented
}

func (m *Mailbox) UnmarshalJSON(b []byte) error {

	return utils.ErrNotImplemented
}
