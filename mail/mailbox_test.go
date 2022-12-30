package mail_test

import (
	"os"
	"testing"

	"github.com/cwinters8/gomap/client"
	"github.com/cwinters8/gomap/mail"
	"github.com/cwinters8/gomap/methods"
	"github.com/cwinters8/gomap/utils"
	"github.com/google/uuid"
)

func TestMailbox(t *testing.T) {
	if err := utils.Env("../.env"); err != nil {
		t.Fatalf("failed to load env: %s", err.Error())
	}
	client, err := client.NewClient(os.Getenv("FASTMAIL_SESSION_URL"), os.Getenv("FASTMAIL_TOKEN"))
	if err != nil {
		t.Fatalf("failed to instantiate new client: %s", err.Error())
	}
	t.Run("query", func(t *testing.T) {
		inbox := "Inbox"
		box, err := mail.NewMailbox(client, inbox)
		if err != nil {
			t.Fatalf("failed to instantiate new mailbox %s: %s", inbox, err.Error())
		}
		wantID, err := uuid.Parse(os.Getenv("FASTMAIL_INBOX_ID"))
		if err != nil {
			t.Fatalf("failed to parse uuid: %s", err.Error())
		}
		if wantID != box.ID {
			t.Errorf("wanted ID %s; got ID %s", wantID, box.ID)
		}
	})
	t.Run("draft", func(t *testing.T) {
		drafts := "Drafts"
		box, err := mail.NewMailbox(client, drafts)
		if err != nil {
			t.Fatalf("failed to instantiate new mailbox %s: %s", drafts, err.Error())
		}
		id, err := box.NewEmail(
			&methods.Address{
				Name:  "Clark the Gopher",
				Email: "dev@clarkwinters.com",
			},
			&methods.Address{
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
