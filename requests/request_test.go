package requests_test

import (
	"os"
	"testing"

	"github.com/cwinters8/gomap/client"
	"github.com/cwinters8/gomap/mail"
	"github.com/cwinters8/gomap/methods"
	"github.com/cwinters8/gomap/requests"
	"github.com/cwinters8/gomap/results"
	"github.com/cwinters8/gomap/utils"
	"github.com/google/uuid"
)

func TestSendRequest(t *testing.T) {
	if err := utils.Env("../.env"); err != nil {
		t.Fatalf("failed to load env: %s", err.Error())
	}
	c, err := client.NewClient(os.Getenv("FASTMAIL_SESSION_URL"), os.Getenv("FASTMAIL_TOKEN"))
	if err != nil {
		t.Fatalf("failed to instantiate a new client: %s", err.Error())
	}

	t.Run("query", func(t *testing.T) {
		query := methods.Query{
			AccountID: c.Session.PrimaryAccounts.Mail,
			Filter: methods.Filter{
				Name: "Inbox",
			},
		}
		i, err := requests.NewInvocation(query, "Mailbox", methods.QueryMethod)
		if err != nil {
			t.Fatalf("failed to instantiate new invocation: %s", err.Error())
		}
		req := requests.NewRequest([]*requests.Invocation[methods.Query]{i})
		resp, err := req.Send(c)
		if err != nil {
			t.Fatalf("failed to send request: %s", err.Error())
		}

		q, ok := resp.Results[0].(*results.Query)
		if !ok {
			t.Fatalf("failed to cast result to Set. %s", utils.Describe(resp.Results[0]))
		}

		gotInboxID := q.Body.IDs[0]
		wantInboxID, err := uuid.Parse(os.Getenv("FASTMAIL_INBOX_ID"))
		if err != nil {
			t.Fatalf("failed to parse inbox uuid: %s", err.Error())
		}

		cases := []*utils.Case{{
			Check:   i.ID != q.ID,
			Message: "wanted invocation id %s; got %s",
			Args:    []any{i.ID, q.ID},
		}, {
			Check:   i.Method.Prefix != q.Prefix,
			Message: "wanted prefix %s; got %s",
			Args:    []any{i.Method.Prefix, q.Prefix},
		}, {
			Check:   wantInboxID != gotInboxID,
			Message: "wanted inbox id %s; got %s",
			Args:    []any{wantInboxID, gotInboxID},
		}}
		for _, c := range cases {
			if c.Check {
				t.Errorf(c.Message, c.Args...)
			}
		}
	})
	t.Run("set", func(t *testing.T) {
		box, err := mail.NewMailbox(c, "Drafts")
		if err != nil {
			t.Fatalf("failed to get Drafts mailbox: %s", err.Error())
		}
		msg, err := methods.NewMessage(
			methods.Mailboxes{box.ID},
			&methods.Address{
				Name:  "Gopher Clark",
				Email: "dev@clarkwinters.com",
			},
			&methods.Address{
				Name:  "Request the Setter",
				Email: "tester@clarkwinters.com",
			},
			"requesting Email/set",
			"attempting to make a request to create a new email",
		)
		if err != nil {
			t.Fatalf("failed to instantiate new message: %s", err.Error())
		}
		set := methods.Set{
			AccountID: c.Session.PrimaryAccounts.Mail,
			Create:    msg,
		}
		i, err := requests.NewInvocation(set, "Email", methods.SetMethod)
		if err != nil {
			t.Fatalf("failed to instantiate new invocation: %s", err.Error())
		}
		req := requests.NewRequest([]*requests.Invocation[methods.Set]{i})
		if _, err := req.Send(c); err != nil {
			t.Fatalf("request failure: %s", err.Error())
		}
	})
}
