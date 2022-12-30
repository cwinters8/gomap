package results_test

import (
	"os"
	"testing"

	"github.com/cwinters8/gomap/client"
	"github.com/cwinters8/gomap/methods"
	"github.com/cwinters8/gomap/requests"
	"github.com/cwinters8/gomap/utils"
)

func TestResultSetJSON(t *testing.T) {
	// validate json unmarshal using a Request
	if err := utils.Env("../.env"); err != nil {
		t.Fatalf("failed to load env: %s", err.Error())
	}
	c, err := client.NewClient(os.Getenv("FASTMAIL_SESSION_URL"), os.Getenv("FASTMAIL_TOKEN"))
	if err != nil {
		t.Fatalf("failed to instantiate a new client: %s", err.Error())
	}
	msg, err := methods.NewMessage(
		methods.Mailboxes{"Drafts"},
		&methods.Address{
			Name:  "Gopher Clark",
			Email: "dev@clarkwinters.com",
		}, &methods.Address{
			Name:  "Setter Tester",
			Email: "tester@clarkwinters.com",
		},
		"hope this works",
		"trying to parse result of set request to json",
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
	resp, err := req.Send(c)
	if err != nil {
		t.Fatalf("request failed: %s", err.Error())
	}
	if resp == nil {
		t.Fatalf("response is nil")
	}
}
