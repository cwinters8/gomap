package mail

import (
	"fmt"

	"github.com/cwinters8/gomap/client"
	"github.com/cwinters8/gomap/requests"
	"github.com/cwinters8/gomap/requests/arguments"
	"github.com/cwinters8/gomap/utils"

	"github.com/google/uuid"
)

type Mailbox struct {
	ID     string      `json:"id"`
	Name   string      `json:"name"`
	Emails []uuid.UUID // IDs of emails associated with this mailbox
	Client *client.Client
}

func NewMailbox(client *client.Client, name string) (*Mailbox, error) {
	m := &Mailbox{
		Client: client,
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

	i := requests.Invocation[arguments.Query]{
		Method: &requests.Method[arguments.Query]{
			Prefix: "Mailbox",
			Type:   requests.QueryMethod,
			Args: arguments.Query{
				AccountID: m.Client.Session.PrimaryAccounts.Mail,
				Filter: arguments.Filter{
					Name: m.Name,
				},
			},
		},
	}
	resp, err := requests.NewRequest([]*requests.Invocation[arguments.Query]{&i}).Send(m.Client)
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
	i := requests.Invocation[arguments.Set]{
		Method: &requests.Method[arguments.Set]{
			Prefix: "Email",
			Type:   requests.SetMethod,
			Args: arguments.Set{
				AccountID: m.Client.Session.PrimaryAccounts.Mail,
				Create:    message,
			},
		},
	}
	// create and send request for i
	resp, err := requests.NewRequest([]*requests.Invocation[arguments.Set]{&i}).Send(m.Client)
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
