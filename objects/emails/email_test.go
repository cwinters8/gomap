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

func TestEmailSet(t *testing.T) {
	envPath := "../../.env"
	if err := utils.Env(envPath); err != nil {
		t.Fatalf("failed to load env variables from `%s`: %s", envPath, err.Error())
	}
	c, err := client.NewClient(os.Getenv("FASTMAIL_SESSION_URL"), os.Getenv("FASTMAIL_TOKEN"))
	if err != nil {
		t.Fatalf("failed to construct new client: %s", err.Error())
	}
	m := mailboxes.Mailbox{
		Name: "Drafts",
	}
	boxCall, err := m.Query(c.Session.PrimaryAccounts.Mail)
	if err != nil {
		t.Fatalf("failed to construct new query call: %s", err.Error())
	}
	if err := requests.Request(c, []*requests.Call{boxCall}, false); err != nil {
		t.Fatalf("mailbox query request failure: %s", err.Error())
	}
	if len(m.ID) < 1 {
		t.Fatalf("mailbox id was not populated")
	}
	e, err := emails.NewEmail(
		[]string{m.ID},
		[]*emails.Address{{
			Name:  "Gopher Clark",
			Email: "dev@clarkwinters.com",
		}},
		[]*emails.Address{{
			Name:  "Tester McSeterson",
			Email: "tester@clarkwinters.com",
		}},
		"testing Email/set",
		"hello from TestEmailSet!",
		emails.TextPlain,
	)
	if err != nil {
		t.Fatalf("failed to construct new email: %s", err.Error())
	}
	call, err := e.Set(c.Session.PrimaryAccounts.Mail)
	if err != nil {
		t.Fatalf("failed to construct set call: %s", err.Error())
	}
	if err := requests.Request(c, []*requests.Call{call}, false); err != nil {
		t.Fatalf("email set request failure: %s", err.Error())
	}
	if len(e.ID) < 1 {
		t.Error("wanted non-empty email id")
	}
	// TODO: make Email/get request with returned id to validate email was created with the correct properties
}
