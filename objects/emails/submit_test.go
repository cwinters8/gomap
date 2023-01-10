package emails_test

import (
	"os"
	"testing"

	"github.com/cwinters8/gomap/client"
	"github.com/cwinters8/gomap/objects/emails"
	"github.com/cwinters8/gomap/objects/mailboxes"
	"github.com/cwinters8/gomap/utils"
)

func TestSubmit(t *testing.T) {
	if err := utils.Env(envPath); err != nil {
		t.Fatalf("failed to source env variables from path `%s`: %s", envPath, err.Error())
	}
	c, err := client.NewClient(os.Getenv("FASTMAIL_SESSION_URL"), os.Getenv("FASTMAIL_TOKEN"))
	if err != nil {
		t.Fatalf("failed to construct new client: %s", err.Error())
	}
	// TODO: get Drafts mailbox ID and create email in it
	draftBox, err := mailboxes.GetMailboxByName(c, "Drafts")
	if err != nil {
		t.Fatalf("failed to retrieve draft mailbox: %s", err.Error())
	}
	e, err := emails.NewEmail(
		[]string{draftBox.ID},
		[]*emails.Address{{
			Name:  "Gopher Clark",
			Email: "dev@clarkwinters.com",
		}},
		[]*emails.Address{{
			Name:  "Tester McSubmit",
			Email: "tester@clarkwinters.com",
		}},
		"testing email submission",
		"hello from TestSubmit!",
		emails.TextPlain,
	)
	if err != nil {
		t.Fatalf("failed to construct new email: %s", err.Error())
	}
	if err := emails.Set(c, []*emails.Email{e}); err != nil {
		t.Fatalf("email set request failure: %s", err.Error())
	}
	if len(e.ID) < 1 {
		t.Fatalf("email id not populated from set request")
	}
	sentBox, err := mailboxes.GetMailboxByName(c, "Sent")
	if err != nil {
		t.Fatalf("failed to retrieve sent mailbox: %s", err.Error())
	}
	id, err := e.Submit(c, draftBox.ID, sentBox.ID)
	if err != nil {
		t.Fatalf("email submission failure: %s", err.Error())
	}
	if len(id) < 1 {
		t.Fatalf("got zero-length submission id")
	}
	// at this time, I'm unable to validate an email was delivered,
	// because Fastmail seems to delete the EmailSubmission object as soon as the email is successfully sent
	// from the jmap mail spec:
	// > "For efficiency, a server MAY destroy EmailSubmission objects at any time after the message is successfully sent or after it has finished retrying to send the message."
	// source: https://jmap.io/spec-mail.html#email-submission
}
