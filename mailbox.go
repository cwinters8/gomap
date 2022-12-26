package gomap

import (
	"fmt"

	"github.com/cwinters8/gomap/arguments"
	"github.com/cwinters8/gomap/utils"
)

type Mailbox struct {
	ID     string   `json:"id"`
	Name   string   `json:"name"`
	Emails []string // IDs of emails associated with this mailbox
	client *Client
}

func (c *Client) NewMailbox(name string) *Mailbox {
	return &Mailbox{
		client: c,
		Name:   name,
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

// creates a new email assigned to this mailbox
// and returns the email's ID
//
// note that this does not send the email to the recipient.
// it will simply create a "draft" email in the mailbox
// that can later be sent using SubmitEmail
func (m *Mailbox) NewEmail(to, subject, msg string) (string, error) {

	return "", utils.ErrNotImplemented
}

func (m Mailbox) MarshalJSON() ([]byte, error) {

	return nil, utils.ErrNotImplemented
}

func (m *Mailbox) UnmarshalJSON(b []byte) error {

	return utils.ErrNotImplemented
}
