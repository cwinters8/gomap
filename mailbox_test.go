package gomap_test

import (
	"os"
	"testing"

	"github.com/cwinters8/gomap"
)

func TestMailbox(t *testing.T) {
	env(t)
	client, err := gomap.NewClient(os.Getenv("FASTMAIL_SESSION_URL"), os.Getenv("FASTMAIL_TOKEN"))
	if err != nil {
		failf(t, "failed to instantiate new client: %s", err.Error())
	}
	t.Run("query", func(t *testing.T) {
		box := client.NewMailbox()
		box.Name = "Inbox"
		if err := box.Query(); err != nil {
			failf(t, "failed to query for mailbox %s: %s", box.Name, err.Error())
		}
		wantID := os.Getenv("FASTMAIL_INBOX_ID")
		if wantID != box.ID {
			t.Errorf("wanted ID %s; got ID %s", wantID, box.ID)
		}
	})
}
