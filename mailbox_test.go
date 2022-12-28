package gomap_test

import (
	"os"
	"testing"

	"github.com/cwinters8/gomap"
	"github.com/cwinters8/gomap/arguments"
	"github.com/cwinters8/gomap/utils"
)

func TestMailbox(t *testing.T) {
	utils.Env(t)
	client, err := gomap.NewClient(os.Getenv("FASTMAIL_SESSION_URL"), os.Getenv("FASTMAIL_TOKEN"))
	if err != nil {
		failf(t, "failed to instantiate new client: %s", err.Error())
	}
	t.Run("query", func(t *testing.T) {
		inbox := "Inbox"
		box, err := client.NewMailbox(inbox)
		if err != nil {
			t.Fatalf("failed to instantiate new mailbox %s: %s", inbox, err.Error())
		}
		wantID := os.Getenv("FASTMAIL_INBOX_ID")
		if wantID != box.ID {
			t.Errorf("wanted ID %s; got ID %s", wantID, box.ID)
		}
	})
	t.Run("draft", func(t *testing.T) {
		drafts := "Drafts"
		box, err := client.NewMailbox(drafts)
		if err != nil {
			t.Fatalf("failed to instantiate new mailbox %s: %s", drafts, err.Error())
		}
		id, err := box.NewEmail(
			&arguments.Address{
				Name:  "Clark the Gopher",
				Email: "dev@clarkwinters.com",
			},
			&arguments.Address{
				Name:  "Tester McTesterson",
				Email: "tester@clarkwinters.com",
			},
			"hello world!",
			"this should land in the Drafts folder",
		)
		if err != nil {
			t.Fatalf("failed to create draft email: %s", err.Error())
		}
		if len(box.Emails) < 1 {
			t.Fatalf("new email not assigned to mailbox")
		}
		if id != box.Emails[0] {
			t.Errorf("returned ID and stored ID do not match. %s returned; %s stored", id, box.Emails[0])
		}
	})
}
