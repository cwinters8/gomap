package requests_test

import (
	"os"
	"testing"

	"github.com/cwinters8/gomap/client"
	"github.com/cwinters8/gomap/requests"
	"github.com/cwinters8/gomap/requests/arguments"
	"github.com/cwinters8/gomap/utils"

	"github.com/google/uuid"
)

func TestSendRequest(t *testing.T) {
	if err := utils.Env("../.env"); err != nil {
		t.Fatalf("failed to load env: %s", err.Error())
	}
	client, err := client.NewClient(os.Getenv("FASTMAIL_SESSION_URL"), os.Getenv("FASTMAIL_TOKEN"))
	if err != nil {
		t.Fatalf("failed to instantiate a new client: %s", err.Error())
	}
	id, err := uuid.NewRandom()
	if err != nil {
		t.Fatalf("failed to generate new uuid: %s", err.Error())
	}
	i := requests.Invocation[arguments.Query]{
		ID: id,
		Method: &requests.Method[arguments.Query]{
			Prefix: "Mailbox",
			Type:   requests.QueryMethod,
			Args: arguments.Query{
				AccountID: client.Session.PrimaryAccounts.Mail,
				Filter: arguments.Filter{
					Name: "Inbox",
				},
			},
		},
	}
	req := requests.NewRequest([]*requests.Invocation[arguments.Query]{&i})
	resp, err := req.Send(client)
	if err != nil {
		t.Fatalf("failed to send request: %s", err.Error())
	}

	gotInv := resp.Results[0]

	gotInboxID := gotInv.Method.Args.IDs[0]
	wantInboxID := os.Getenv("FASTMAIL_INBOX_ID")

	cases := []*utils.Case{{
		Check:  i.ID != gotInv.ID,
		Format: "wanted invocation id %s; got %s",
		Args:   []any{i.ID, gotInv.ID},
	}, {
		Check:  i.Method.Prefix != gotInv.Method.Prefix,
		Format: "wanted prefix %s; got %s",
		Args:   []any{i.Method.Prefix, gotInv.Method.Prefix},
	}, {
		Check:  wantInboxID != gotInboxID,
		Format: "wanted inbox id %s; got %s",
		Args:   []any{wantInboxID, gotInboxID},
	}}
	for _, c := range cases {
		if c.Check {
			t.Errorf(c.Format, c.Args...)
		}
	}
}
