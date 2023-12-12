package emails_test

import (
	"bytes"
	"fmt"
	"html/template"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/cwinters8/gomap"
	"github.com/cwinters8/gomap/client"
	"github.com/cwinters8/gomap/objects/emails"
	"github.com/cwinters8/gomap/objects/mailboxes"
	"github.com/cwinters8/gomap/utils"
	"github.com/google/uuid"
)

func TestSendEmail(t *testing.T) {
	if err := utils.Env(envPath); err != nil {
		t.Fatalf("failed to source env variables from path `%s`: %s", envPath, err.Error())
	}
	c, err := gomap.NewClient(
		os.Getenv("FASTMAIL_SESSION_URL"),
		os.Getenv("FASTMAIL_TOKEN"),
		gomap.DefaultDrafts,
		gomap.DefaultSent,
	)
	if err != nil {
		t.Fatalf("failed to construct new client: %v", err)
	}
	testEmailID, err := uuid.NewRandom()
	if err != nil {
		t.Fatalf("failed to generate new uuid: %v", err)
	}
	strEmailID := testEmailID.String()
	msg := `<!DOCTYPE html>
	<html>
		<head>
			<meta charset="UTF-8">
			<title>{{.Title}}</title>
		</head>
		<body>
			<h1>Hello from {{.TestFunc}}!</h1>
			<p>Test ID: {{.TestID}}</p>
		</body>
	</html>`
	tpl, err := template.New("email").Parse(msg)
	if err != nil {
		t.Fatalf("failed to parse HTML template: %v", err)
	}
	var buf bytes.Buffer
	if err := tpl.Execute(&buf, map[string]string{
		"Title":    "Test JMAP Email",
		"TestFunc": "TestSendEmail",
		"TestID":   strEmailID,
	}); err != nil {
		t.Fatalf("failed to execute HTML template: %v", err)
	}
	from := emails.Address{"Gopher Clark", "dev@clarkwinters.com"}
	to := emails.Address{"Tester McSubmit", "tester@clarkwinters.com"}
	subject := "testing gomap.Client.SendEmail"
	if err := c.SendEmail(
		[]*emails.Address{&from},
		[]*emails.Address{&to},
		subject,
		buf.String(),
		true,
	); err != nil {
		t.Fatalf("failed to send email: %v", err)
	}
	// check for matching emails
	time.Sleep(5 * time.Second)
	box, err := c.GetMailbox("🧪-tester")
	if err != nil {
		t.Fatalf("failed to get tester mailbox: %v", err)
	}
	filter := emails.Filter{
		InMailboxID: box.ID,
		Text:        strEmailID,
	}
	msgs, err := c.GetEmails(&filter, 1, 30*time.Second)
	if err != nil {
		t.Fatalf("failed to retrieve emails: %v", err)
	}
	if len(msgs) == 0 {
		t.Fatalf("email containing ID `%s` not found", strEmailID)
	}
	got := msgs[0]
	gotTo := got.To[0]
	gotFrom := got.From[0]
	cases := utils.Cases{
		utils.NewCase(gotTo.Email != to.Email, "wanted to email `%s`; got `%s`", to.Email, gotTo.Email),
		utils.NewCase(gotFrom.Email != from.Email, "wanted from email `%s`; got `%s`", from.Email, gotFrom.Email),
		utils.NewCase(got.Subject != subject, "wanted subject `%s`; got `%s`", subject, got.Subject),
	}
	cases.Iterator(func(c *utils.Case) {
		t.Error(c.Message)
	})
}

func TestSubmit(t *testing.T) {
	if err := utils.Env(envPath); err != nil {
		t.Fatalf("failed to source env variables from path `%s`: %s", envPath, err.Error())
	}
	c, err := client.NewClient(os.Getenv("FASTMAIL_SESSION_URL"), os.Getenv("FASTMAIL_TOKEN"))
	if err != nil {
		t.Fatalf("failed to construct new client: %s", err.Error())
	}
	draftBox, err := mailboxes.GetMailboxByName(c, "Drafts")
	if err != nil {
		t.Fatalf("failed to retrieve draft mailbox: %s", err.Error())
	}
	testEmailID, err := uuid.NewRandom()
	if err != nil {
		t.Fatalf("failed to generate new uuid: %s", err.Error())
	}
	strID := testEmailID.String()
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
		fmt.Sprintf("hello from TestSubmit!\ntest id: %s", strID),
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
	// retrieve the email from the mailbox it should have landed in
	boxName := "🧪-tester"
	box, err := mailboxes.GetMailboxByName(c, boxName)
	if err != nil {
		t.Fatalf("failed to retrieve mailbox `%s`: %s", boxName, err.Error())
	}
	attempts := 60
	gotID := ""
	filter := emails.Filter{
		InMailboxID: box.ID,
		Text:        strID,
	}
	time.Sleep(5 * time.Second)
	for attempts > 0 {
		emailIDs, err := emails.Query(c, &filter)
		if err != nil {
			t.Fatalf("email query failed: %s", err.Error())
		}
		if len(emailIDs) < 1 {
			attempts--
			continue
		}
		gotID = emailIDs[0]
		break
	}
	if len(gotID) < 1 {
		t.Fatalf("email with test id `%s` not found", testEmailID.String())
	}
	// validate email content
	found, _, err := emails.GetEmails(c, []string{gotID})
	if err != nil {
		t.Fatalf("failed to retrieve email: %s", err.Error())
	}
	if len(found) < 1 {
		t.Fatalf("email id `%s` not found", gotID)
	}
	got := found[0]
	cases := utils.Cases{utils.NewCase(
		e.From[0].Email != got.From[0].Email,
		"wanted from email `%s`; got `%s`",
		e.From[0].Email, got.From[0].Email,
	), utils.NewCase(
		e.To[0].Email != got.To[0].Email,
		"wanted to email `%s`; got `%s`",
		e.To[0].Email, got.To[0].Email,
	), utils.NewCase(
		e.Subject != got.Subject,
		"wanted subject `%s`; got `%s`",
		e.Subject, got.Subject,
	), utils.NewCase(
		!strings.Contains(got.Body.Value, e.Body.Value),
		"wanted body value `%s`; got `%s`",
		e.Body.Value, got.Body.Value,
	)}
	cases.Iterator(func(c *utils.Case) {
		t.Error(c.Message)
	})
}
