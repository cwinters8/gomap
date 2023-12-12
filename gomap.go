package gomap

import (
	"fmt"
	"time"

	"github.com/cwinters8/gomap/client"
	"github.com/cwinters8/gomap/objects/emails"
	"github.com/cwinters8/gomap/objects/mailboxes"
)

type Client struct {
	*client.Client
	Drafts *mailboxes.Mailbox
	Sent   *mailboxes.Mailbox
}

// NewClient creates a new JMAP mail client that can be used for interacting
// with the JMAP mail server specified with jmapSessionURL.
//
// Commonly, draftsMailbox should be "Drafts" and sentMailbox should be "Sent".
// These arguments are available in case customization is necessary.
// You can use the DefaultDrafts and DefaultSent constants for convenience.
func NewClient(jmapSessionURL, bearerToken, draftsMailbox, sentMailbox string) (*Client, error) {
	client, err := client.NewClient(jmapSessionURL, bearerToken)
	if err != nil {
		return nil, fmt.Errorf("failed to instantiate jmap client: %w", err)
	}
	drafts, err := mailboxes.GetMailboxByName(client, draftsMailbox)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve drafts mailbox: %w", err)
	}
	sent, err := mailboxes.GetMailboxByName(client, sentMailbox)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve sent mailbox: %w", err)
	}
	c := Client{client, drafts, sent}
	return &c, nil
}

// SendEmail sends an email using the provided arguments.
//
// Setting isHTML to true will set the body type attribute to HTML instead of plaintext.
// This works best if body is a string that has been output from executing an html/template.
func (c *Client) SendEmail(from, to []*emails.Address, subject, body string, isHTML bool) error {
	bodyType := emails.TextPlain
	if isHTML {
		bodyType = emails.TextHTML
	}
	email, err := emails.NewEmail([]string{c.Drafts.ID}, from, to, subject, body, bodyType)
	if err != nil {
		return fmt.Errorf("failed to instantiate new email: %w", err)
	}
	if err := emails.Set(c.Client, []*emails.Email{email}); err != nil {
		return fmt.Errorf("email set request failure: %w", err)
	}
	if _, err := email.Submit(c.Client, c.Drafts.ID, c.Sent.ID); err != nil {
		return fmt.Errorf("failed to submit email: %w", err)
	}
	return nil
}

// GetEmails retrieves emails based on the provided filter.
// It will continue to query until the first of maxCount or timeout has been reached.
//
// maxCount is used as a metric for breaking out of the query loop,
// not a hard limit on the number of emails returned.
func (c *Client) GetEmails(filter *emails.Filter, maxCount int, timeout time.Duration) ([]*emails.Email, error) {
	var emailIDs []string
	end := time.Now().Add(timeout)
	for time.Now().Compare(end) < 1 {
		newIDs, err := emails.Query(c.Client, filter)
		if err != nil {
			return nil, fmt.Errorf("failed to query for emails: %w", err)
		}
		if len(newIDs) > 0 {
			emailIDs = append(emailIDs, newIDs...)
		}
		if len(emailIDs) >= maxCount {
			break
		}
	}
	if len(emailIDs) == 0 {
		return nil, fmt.Errorf("email IDs matching provided filter not found")
	}
	found, notFound, err := emails.GetEmails(c.Client, emailIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve emails: %w", err)
	}
	if len(notFound) > 0 {
		fmt.Printf("Warning: Unable to retrieve emails with IDs %v\n", notFound)
	}
	return found, nil
}

// GetMailbox retrieves a mailbox with the matching name.
func (c *Client) GetMailbox(name string) (*mailboxes.Mailbox, error) {
	return mailboxes.GetMailboxByName(c.Client, name)
}

const (
	DefaultDrafts = "Drafts"
	DefaultSent   = "Sent"
)
