package requests_test

import (
	"os"
	"testing"

	"github.com/cwinters8/gomap/client"
	"github.com/cwinters8/gomap/objects/emails"
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

	t.Run("single call", func(t *testing.T) {
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
	})

	t.Run("multiple calls", func(t *testing.T) {
		// TODO: figure out which 2 calls to make in a single request
		e := emails.Email{
			ID: "M55f24a60e9751599460d8fa4",
		}
		mailCall, err := e.Get(c.Session.PrimaryAccounts.Mail)
		if err != nil {
			t.Fatalf("failed to construct email get call: %s", err.Error())
		}
		box := mailboxes.Mailbox{
			Name: "🧪-tester",
		}
		boxCall, err := box.Query(c.Session.PrimaryAccounts.Mail)
		if err != nil {
			t.Fatalf("failed to construct mailbox query call: %s", err.Error())
		}
		responses, err := requests.Request(c, []*requests.Call{mailCall, boxCall}, false)
		if err != nil {
			t.Fatalf("request failure: %s", err.Error())
		}
		cases := utils.Cases{}
		for _, resp := range responses {
			if resp == nil {
				t.Fatalf("found nil response")
			}
			switch resp.ID {
			case mailCall.ID:
				found, notFound, err := emails.ParseRawResponseBody(resp.Body)
				if err != nil {
					t.Fatalf("failed to parse email response body: %s", err.Error())
				}
				if len(notFound) > 0 {
					t.Errorf("at least 1 email id not found: %v", notFound)
				}
				got := found[0]
				wantSubject := "attempting to set and get an email in a single request"
				cases.Append(utils.NewCase(
					got.ID != e.ID,
					"wanted email id %s; got %s",
					e.ID, got.ID,
				), utils.NewCase(
					got.Subject != wantSubject,
					"wanted email subject %s; got %s",
					wantSubject, got.Subject,
				))
			case boxCall.ID:
				wantID := "87d37719-fe88-44da-b80a-1ca15aec71c0"
				cases.Append(utils.NewCase(
					wantID != box.ID,
					"wanted mailbox id %s; got %s",
					wantID, box.ID,
				))
			default:
				t.Errorf("unknown request id found in response: %s", resp.ID)
			}
		}
		cases.Iterator(func(c *utils.Case) {
			t.Error(c.Message)
		})
	})
}
