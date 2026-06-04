package requests_test

import (
	"os"
	"testing"

	"github.com/cwinters8/gomap/client"
	"github.com/cwinters8/gomap/objects/mailboxes"
	"github.com/cwinters8/gomap/parse"
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
		responses, err := requests.Request(c, []*requests.Call{call}, false)
		if err != nil {
			t.Fatalf("request failure: %s", err.Error())
		}
		if len(responses) != 1 {
			t.Fatalf("wanted 1 response; got %d", len(responses))
		}
		ids, err := parse.QueryResponseBody(responses[0].Body)
		if err != nil {
			t.Fatalf("failed to parse query response body: %s", err.Error())
		}
		if len(ids) == 0 {
			t.Fatal("wanted at least 1 mailbox id")
		}
		if box.ID != ids[0] {
			t.Errorf("wanted mailbox id %s; got %s", ids[0], box.ID)
		}
	})

	t.Run("multiple calls", func(t *testing.T) {
		firstBox := mailboxes.Mailbox{
			Name: "Inbox",
		}
		firstCall, err := firstBox.Query(c.Session.PrimaryAccounts.Mail)
		if err != nil {
			t.Fatalf("failed to construct first mailbox query call: %s", err.Error())
		}
		secondBox := mailboxes.Mailbox{
			Name: "Inbox",
		}
		secondCall, err := secondBox.Query(c.Session.PrimaryAccounts.Mail)
		if err != nil {
			t.Fatalf("failed to construct second mailbox query call: %s", err.Error())
		}
		responses, err := requests.Request(c, []*requests.Call{firstCall, secondCall}, false)
		if err != nil {
			t.Fatalf("request failure: %s", err.Error())
		}
		if len(responses) != 2 {
			t.Fatalf("wanted 2 responses; got %d", len(responses))
		}
		firstSeen := false
		secondSeen := false
		for _, resp := range responses {
			if resp == nil {
				t.Fatalf("found nil response")
			}
			ids, err := parse.QueryResponseBody(resp.Body)
			if err != nil {
				t.Fatalf("failed to parse query response body: %s", err.Error())
			}
			if len(ids) == 0 {
				t.Fatal("wanted at least 1 mailbox id")
			}
			switch resp.ID {
			case firstCall.ID:
				firstSeen = true
				if firstBox.ID != ids[0] {
					t.Errorf("wanted first mailbox id %s; got %s", ids[0], firstBox.ID)
				}
			case secondCall.ID:
				secondSeen = true
				if secondBox.ID != ids[0] {
					t.Errorf("wanted second mailbox id %s; got %s", ids[0], secondBox.ID)
				}
			default:
				t.Errorf("unknown request id found in response: %s", resp.ID)
			}
		}
		if !firstSeen {
			t.Error("missing response for first call")
		}
		if !secondSeen {
			t.Error("missing response for second call")
		}
	})
}
