package gomap

import (
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
			Prefix: "Mailbox",
			Args: arguments.Query{
				AccountID: m.client.Session.PrimaryAccounts.Mail,
				Filter: arguments.Filter{
					Name: m.Name,
				},
			},
		},
	}
	resp, err := NewRequest([]*Invocation[arguments.Query]{&i}).Send(m.client)
	if err != nil {
		return fmt.Errorf("failed to send query request: %w", err)
	}
	m.ID = resp.Results[0].Method.Args.IDs[0]
	return nil
}

func (m Mailbox) MarshalJSON() ([]byte, error) {

	return nil, utils.ErrNotImplemented
}

func (m *Mailbox) UnmarshalJSON(b []byte) error {

	return utils.ErrNotImplemented
}
