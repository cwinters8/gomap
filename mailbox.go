package gomap

import (
	"fmt"

	"github.com/cwinters8/gomap/arguments"
	"github.com/cwinters8/gomap/utils"

	"github.com/google/uuid"
)

type Mailbox struct {
	ID     string      `json:"id"`
	Name   string      `json:"name"`
	Emails []uuid.UUID // IDs of emails associated with this mailbox
	client *Client
}

func (c *Client) NewMailbox(name string) (*Mailbox, error) {
	m := &Mailbox{
		client: c,
		Name:   name,
	}
	if err := m.Query(); err != nil {
		return nil, fmt.Errorf("failed to query for mailbox: %w", err)
	}
	return m, nil
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
			Type:   QueryMethod,
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
//
// note that this does not send the email to the recipient.
// it will simply create a "draft" email in the mailbox
// that can later be sent using SubmitEmail
func (m *Mailbox) NewEmail(from, to *arguments.Address, subject, msg string) (uuid.UUID, error) {
	if len(m.ID) < 1 {
		return uuid.Nil, fmt.Errorf("m.ID must have a valid ID in order to create a new email. try running m.Query first")
	}
	message, err := arguments.NewMessage(
		arguments.Mailboxes{m.ID},
		from,
		to,
		subject,
		msg,
	)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to instantiate new message: %w", err)
	}
	i := Invocation[arguments.Set]{
		Method: &Method[arguments.Set]{
			Prefix: "Email",
			Type:   SetMethod,
			Args: arguments.Set{
				AccountID: m.client.Session.PrimaryAccounts.Mail,
				Create:    message,
			},
		},
	}
	// create and send request for i
	resp, err := NewRequest([]*Invocation[arguments.Set]{&i}).Send(m.client)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to send email set request: %w", err)
	}
	if len(resp.Results) < 1 {
		return uuid.Nil, fmt.Errorf("results slice is empty")
	}
	create := resp.Results[0].Method.Args.Create
	if create == nil {
		return uuid.Nil, fmt.Errorf("result method's Create message is nil")
	}
	id := resp.Results[0].Method.Args.Create.ID
	m.Emails = append(m.Emails, id)
	return id, utils.ErrNotImplemented
}

func (m Mailbox) MarshalJSON() ([]byte, error) {

	return nil, utils.ErrNotImplemented
}

func (m *Mailbox) UnmarshalJSON(b []byte) error {

	return utils.ErrNotImplemented
}
