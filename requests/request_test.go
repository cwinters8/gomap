package requests_test

import (
	"os"
	"testing"

	"github.com/cwinters8/gomap/client"
	"github.com/cwinters8/gomap/objects/mailboxes"
	"github.com/cwinters8/gomap/requests"
	"github.com/cwinters8/gomap/utils"
)

const envPath = "../.env"

func TestRequest(t *testing.T) {
	if err := utils.Env(envPath); err != nil {
		t.Fatalf("failed to load `%s`: %s", envPath, err.Error())
	}
	c, err := client.NewClient(os.Getenv("FASTMAIL_SESSION_URL"), os.Getenv("FASTMAIL_TOKEN"))
	if err != nil {
		t.Fatalf("failed to construct new client: %s", err.Error())
	}
	box := mailboxes.Mailbox{
		Name: "Inbox",
	}
	call, err := box.Query(c.Session.PrimaryAccounts.Mail)
	if err != nil {
		t.Fatalf("failed to construct new mailbox query: %s", err.Error())
	}
	if _, err := requests.Request(c, []*requests.Call{call}, false); err != nil {
		t.Fatalf("request failure: %s", err.Error())
	}
	wantID := os.Getenv("FASTMAIL_INBOX_ID")
	if box.ID != wantID {
		t.Errorf("wanted mailbox id %s; got %s", wantID, box.ID)
	}
}
