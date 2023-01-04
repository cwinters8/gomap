package requests_test

import (
	"os"
	"testing"

	"github.com/cwinters8/gomap/client"
	"github.com/cwinters8/gomap/objects/emails"
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
		query, err := requests.NewQuery(c.Session.PrimaryAccounts.Mail, "Mailbox", "Inbox")
		if err != nil {
			t.Fatalf("failed to instantiate new Query: %s", err.Error())
		}
		req := requests.NewRequest([]requests.Call{query})
		resp, err := req.SendAndParse(c)
		if err != nil {
			t.Fatalf("failed to send request: %s", err.Error())
		}

		q, ok := resp.Results[0].(*results.Query)
		if !ok {
			t.Fatalf("failed to cast result to Set. %s", utils.Describe(resp.Results[0]))
		}

		gotInboxID := q.Body.IDs[0]
		wantInboxID := os.Getenv("FASTMAIL_INBOX_ID")
		cases := []*utils.Case{{
			Check:   query.ID != q.ID,
			Message: "wanted invocation id %s; got %s",
			Args:    []any{query.ID, q.ID},
		}, {
			Check:   query.Prefix != q.Prefix,
			Message: "wanted prefix %s; got %s",
			Args:    []any{query.Prefix, q.Prefix},
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
	t.Run("set and get", func(t *testing.T) {
		query, err := requests.NewQuery(c.Session.PrimaryAccounts.Mail, "Mailbox", "Drafts")
		if err != nil {
			t.Fatalf("failed to instantiate new Query: %s", err.Error())
		}
		resp, err := requests.NewRequest([]requests.Call{query}).SendAndParse(c)
		if err != nil {
			t.Fatalf("request failure: %s", err.Error())
		}
		q, ok := resp.Results[0].(*results.Query)
		if !ok {
			t.Fatalf("failed to cast result to Query. %s", utils.Describe(resp.Results[0]))
		}
		wantEmail := emails.Email{
			MailboxIDs: q.Body.IDs,
			From: &emails.Address{
				Name:  "Gopher Clark",
				Email: "dev@clarkwinters.com",
			},
			To: []*emails.Address{{
				Name:  "Request the Setter",
				Email: "tester@clarkwinters.com",
			}},
			Subject: "requesting email/set",
			Body: &emails.Body{
				Type:  emails.TextPlain,
				Value: "attempt to create a new draft email",
			},
		}
		email, err := emails.NewEmail(
			wantEmail.MailboxIDs,
			wantEmail.From,
			wantEmail.To,
			wantEmail.Subject,
			wantEmail.Body.Value,
			wantEmail.Body.Type,
		)
		if err != nil {
			t.Fatalf("failed to create new email: %s", err.Error())
		}
		set, err := requests.NewSet(c.Session.PrimaryAccounts.Mail, email)
		if err != nil {
			t.Fatalf("failed to instantiate new set: %s", err.Error())
		}
		result, err := requests.NewRequest([]requests.Call{set}).SendAndParse(c)
		if err != nil {
			t.Fatalf("request failure: %s", err.Error())
		}
		s, ok := result.Results[0].(*results.Set)
		if !ok {
			t.Fatalf("failed to cast result to set. %s", utils.Describe(result.Results[0]))
		}
		if s.Body.Created == nil {
			t.Fatalf("Body.Created is nil. remaining Body attributes: %v", *s.Body)
		}
		fatals := utils.Cases{
			utils.NewCase(
				s.Body.Created.ID == uuid.Nil,
				"wanted Body.Created.ID to be a non-nil uuid; got nil value",
			),
			utils.NewCase(
				s.Body.NotCreated != nil,
				"wanted Body.NotCreated to be nil; got %v",
				s.Body.NotCreated,
			),
		}
		failed := 0
		fatals.Iterator(func(c *utils.Case) {
			t.Error(c.Message)
			failed++
		})
		if failed > 0 {
			t.FailNow()
		}

		get, err := requests.NewGet(
			c.Session.PrimaryAccounts.Mail,
			"Email",
			[]string{s.Body.Created.ServerID},
			[]string{"from", "to", "subject"},
		)
		if err != nil {
			t.Fatalf("failed to instantiate new Get")
		}
		getResult, err := requests.NewRequest([]requests.Call{get}).SendAndParse(c)
		if err != nil {
			t.Fatalf("get request failure: %s", err.Error())
		}
		g, ok := getResult.Results[0].(*emails.Result)
		if !ok {
			t.Fatalf("failed to cast result to emails.Result. %s", utils.Describe(getResult.Results[0]))
		}
		gotMethod, err := g.Method()
		if err != nil {
			t.Fatalf("failed to get method: %s", err.Error())
		}
		wantMethod := "Email/get"
		gotEmail := g.Body.List[0]
		cases := utils.Cases{
			utils.NewCase(
				gotMethod != wantMethod,
				"wanted method %s; got %s",
				wantMethod, gotMethod,
			),
			utils.NewCase(
				g.Body.AccountID != c.Session.PrimaryAccounts.Mail,
				"wanted account id %s; got %s",
				c.Session.PrimaryAccounts.Mail, g.Body.AccountID,
			),
			utils.NewCase(
				gotEmail.ID != s.Body.Created.ServerID,
				"wanted email id %s; got %s",
				s.Body.Created.ServerID, gotEmail.ID,
			),
			utils.NewCase(
				gotEmail.Subject != wantEmail.Subject,
				"wanted subject %s; got %s",
				wantEmail.Subject, gotEmail.Subject,
			),
			utils.NewCase(
				gotEmail.Body.Value != wantEmail.Body.Value,
				"wanted body value %s; got %s",
				wantEmail.Body.Value, gotEmail.Body.Value,
			),
			utils.NewCase(
				gotEmail.Body.Type != wantEmail.Body.Type,
				"wanted body type %s; got %s",
				wantEmail.Body.Type, gotEmail.Body.Type,
			),
		}
		cases.Iterator(func(c *utils.Case) {
			t.Error(c.Message)
		})
	})
}
