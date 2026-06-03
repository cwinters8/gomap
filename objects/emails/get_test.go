package emails_test

import (
	"os"
	"testing"

	"github.com/cwinters8/gomap/client"
	"github.com/cwinters8/gomap/objects/emails"
	"github.com/cwinters8/gomap/objects/mailboxes"
	"github.com/cwinters8/gomap/requests"
	"github.com/cwinters8/gomap/utils"
)

func TestEmailGet(t *testing.T) {
	if err := utils.Env(envPath); err != nil {
		t.Fatalf("failed to source env variables from path `%s`: %s", envPath, err.Error())
	}
	c, err := client.NewClient(os.Getenv("FASTMAIL_SESSION_URL"), os.Getenv("FASTMAIL_TOKEN"))
	if err != nil {
		t.Fatalf("failed to construct new client: %s", err.Error())
	}
	drafts := mailboxes.Mailbox{
		Name: "Drafts",
	}
	draftsCall, err := drafts.Query(c.Session.PrimaryAccounts.Mail)
	if err != nil {
		t.Fatalf("failed to construct drafts mailbox query call: %s", err.Error())
	}
	if _, err := requests.Request(c, []*requests.Call{draftsCall}, false); err != nil {
		t.Fatalf("drafts mailbox query request failure: %s", err.Error())
	}
	if len(drafts.ID) == 0 {
		t.Fatal("wanted drafts mailbox id to be populated")
	}

	t.Run("single email", func(t *testing.T) {
		want := newTestDraftEmail(t, drafts.ID, "Setter Tester", "hope this works", "trying to parse result of set request to json")
		if err := emails.Set(c, []*emails.Email{want}); err != nil {
			t.Fatalf("failed to create test email: %s", err.Error())
		}
		got := emails.Email{
			ID: want.ID,
		}
		call, err := got.Get(c.Session.PrimaryAccounts.Mail)
		if err != nil {
			t.Fatalf("failed to construct new call for Email/get: %s", err.Error())
		}
		if _, err := requests.Request(c, []*requests.Call{call}, false); err != nil {
			t.Fatalf("Email/get request failure: %s", err.Error())
		}
		assertEmailMatches(t, want, &got)
	})

	t.Run("multiple emails", func(t *testing.T) {
		first := newTestDraftEmail(t, drafts.ID, "Setter Tester", "hope this works", "trying to parse result of set request to json")
		second := newTestDraftEmail(t, drafts.ID, "Tester McSeterson", "testing Email/set", "hello from TestEmailSet!")
		if err := emails.Set(c, []*emails.Email{first, second}); err != nil {
			t.Fatalf("failed to create test emails: %s", err.Error())
		}
		emailIDs := []string{first.ID, second.ID}
		emailList, notFound, err := emails.GetEmails(c, emailIDs)
		if err != nil {
			t.Fatalf("failed to get emails: %s", err.Error())
		}
		if len(notFound) > 0 {
			t.Errorf("some email IDs were not found: %v", notFound)
		}
		if len(emailList) != len(emailIDs) {
			t.Fatalf("wanted %d emails; got %d", len(emailIDs), len(emailList))
		}
		wantByID := map[string]*emails.Email{
			first.ID:  first,
			second.ID: second,
		}
		for _, got := range emailList {
			want, ok := wantByID[got.ID]
			if !ok {
				t.Errorf("unexpected email id found: %s", got.ID)
				continue
			}
			assertEmailMatches(t, want, got)
		}
	})
}

func newTestDraftEmail(t *testing.T, mailboxID, toName, subject, body string) *emails.Email {
	t.Helper()
	email, err := emails.NewEmail(
		[]string{mailboxID},
		[]*emails.Address{{
			Name:  "Gopher Clark",
			Email: "dev@clarkwinters.com",
		}},
		[]*emails.Address{{
			Name:  toName,
			Email: "tester@clarkwinters.com",
		}},
		subject,
		body,
		emails.TextPlain,
	)
	if err != nil {
		t.Fatalf("failed to construct test email: %s", err.Error())
	}
	return email
}

func assertEmailMatches(t *testing.T, want, got *emails.Email) {
	t.Helper()
	if len(got.MailboxIDs) < 1 {
		t.Fatal("MailboxIDs should contain at least 1 id")
	}
	if len(got.From) < 1 {
		t.Fatal("From should contain at least 1 address")
	}
	if len(got.To) < 1 {
		t.Fatal("To should contain at least 1 address")
	}
	if got.Body == nil {
		t.Fatal("Body should not be nil")
	}

	cases := utils.Cases{
		utils.NewCase(
			want.MailboxIDs[0] != got.MailboxIDs[0],
			"wanted mailbox id %s; got %s",
			want.MailboxIDs[0], got.MailboxIDs[0],
		),
		utils.NewCase(
			want.From[0].Name != got.From[0].Name,
			"wanted from name %s; got %s",
			want.From[0].Name, got.From[0].Name,
		),
		utils.NewCase(
			want.From[0].Email != got.From[0].Email,
			"wanted from email %s; got %s",
			want.From[0].Email, got.From[0].Email,
		),
		utils.NewCase(
			want.To[0].Name != got.To[0].Name,
			"wanted to name %s; got %s",
			want.To[0].Name, got.To[0].Name,
		),
		utils.NewCase(
			want.To[0].Email != got.To[0].Email,
			"wanted to email %s; got %s",
			want.To[0].Email, got.To[0].Email,
		),
		utils.NewCase(
			want.Subject != got.Subject,
			"wanted subject %s; got %s",
			want.Subject, got.Subject,
		),
		utils.NewCase(
			want.Body.Type != got.Body.Type,
			"wanted body type %s; got %s",
			want.Body.Type, got.Body.Type,
		),
		utils.NewCase(
			want.Body.Value != got.Body.Value,
			"wanted body value %s; got %s",
			want.Body.Value, got.Body.Value,
		),
	}
	cases.Iterator(func(c *utils.Case) {
		t.Error(c.Message)
	})
}
