package mail

import (
	"fmt"

	"github.com/cwinters8/gomap/client"
	"github.com/cwinters8/gomap/methods"
	"github.com/cwinters8/gomap/requests"
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

	i := requests.Invocation[methods.Query]{
		Method: &methods.Method[methods.Query]{
			Prefix: "Mailbox",
			Type:   methods.QueryMethod,
			Args: methods.Query{
				AccountID: m.Client.Session.PrimaryAccounts.Mail,
				Filter: methods.Filter{
					Name: m.Name,
				},
			},
		},
	}
	resp, err := requests.NewRequest([]*requests.Invocation[methods.Query]{&i}).Send(m.Client)
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
func (m *Mailbox) NewEmail(from, to *methods.Address, subject, msg string) (uuid.UUID, error) {
	if len(m.ID) < 1 {
		return uuid.Nil, fmt.Errorf("m.ID must have a valid ID in order to create a new email. try running m.Query first")
	}
	message, err := methods.NewMessage(
		methods.Mailboxes{m.ID},
		from,
		to,
		subject,
		msg,
	)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to instantiate new message: %w", err)
	}
	i := requests.Invocation[methods.Set]{
		Method: &methods.Method[methods.Set]{
			Prefix: "Email",
			Type:   methods.SetMethod,
			Args: methods.Set{
				AccountID: m.Client.Session.PrimaryAccounts.Mail,
				Create:    message,
			},
		},
	}
	// create and send request for i
	resp, err := requests.NewRequest([]*requests.Invocation[methods.Set]{&i}).Send(m.Client)
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
