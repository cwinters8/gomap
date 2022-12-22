package gomap

import (
	"encoding/json"
	"fmt"
	"net/http"
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
		return fmt.Errorf("m.Name must be non-empty")
	}
	var args MailboxQuery
	args.AccountID = m.client.Session.PrimaryAccounts.Mail
	args.Filter.Name = m.Name
	b, err := json.Marshal(args)
	if err != nil {
		return fmt.Errorf("failed to marshal args to json: %w", err)
	}
	status, _, err := m.client.makeRequest(http.MethodPost, m.client.Session.APIURL, b)
	if err != nil {
		return fmt.Errorf("request status %d\nfailed to make query request: %w", status, err)
	}
	// TODO: unmarshal body to mailbox
	return fmt.Errorf("not implemented")
}

func (m Mailbox) MarshalJSON() ([]byte, error) {

	return nil, fmt.Errorf("not implemented")
}
