package gomap_test

import (
	"os"
	"testing"

	"github.com/cwinters8/gomap"
	"github.com/cwinters8/gomap/utils"
)

func TestMailbox(t *testing.T) {
	utils.Env(t)
	client, err := gomap.NewClient(os.Getenv("FASTMAIL_SESSION_URL"), os.Getenv("FASTMAIL_TOKEN"))
	if err != nil {
		failf(t, "failed to instantiate new client: %s", err.Error())
	}
	t.Run("query", func(t *testing.T) {
		box := client.NewMailbox("Inbox")
		if err := box.Query(); err != nil {
			failf(t, "failed to query for mailbox %s: %s", box.Name, err.Error())
		}
		wantID := os.Getenv("FASTMAIL_INBOX_ID")
		if wantID != box.ID {
			t.Errorf("wanted ID %s; got ID %s", wantID, box.ID)
		}
	})
	t.Run("draft", func(t *testing.T) {
		box := client.NewMailbox("Draft")
		id, err := box.NewEmail(
			"tester@clarkwinters.com",
			"hello world!",
			"this should land in the Drafts folder",
		)
		if err != nil {
			failf(t, "failed to create draft email: %s", err.Error())
		}
		if len(box.Emails) < 1 {
			utils.Failf(t, "new email not assigned to mailbox")
		}
		if id != box.Emails[0] {
			t.Errorf("returned ID and stored ID do not match. %s returned; %s stored", id, box.Emails[0])
		}
	})
}
