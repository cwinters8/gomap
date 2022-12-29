package requests_test

import (
	"os"
	"testing"

	"github.com/cwinters8/gomap/client"
	"github.com/cwinters8/gomap/requests"
	"github.com/cwinters8/gomap/requests/arguments"
	"github.com/cwinters8/gomap/utils"
)

func TestSendRequest(t *testing.T) {
	if err := utils.Env("../.env"); err != nil {
		t.Fatalf("failed to load env: %s", err.Error())
	}
	client, err := client.NewClient(os.Getenv("FASTMAIL_SESSION_URL"), os.Getenv("FASTMAIL_TOKEN"))
	if err != nil {
		utils.Failf(t, "failed to instantiate a new client: %s", err.Error())
	}
	i := requests.Invocation[arguments.Query]{
		ID: "xyz",
		Method: &requests.Method[arguments.Query]{
			Prefix: "Mailbox",
			Args: arguments.Query{
				AccountID: client.Session.PrimaryAccounts.Mail,
				Filter: arguments.Filter{
					Name: "Inbox",
				},
			},
		},
	}
	req := requests.Request[arguments.Query]{
		Using: []requests.Capability{
			requests.UsingCore, // this maybe should just be used by default - I can't think of a case when Core would not be needed
			requests.UsingMail,
		},
		Calls: []*requests.Invocation[arguments.Query]{&i},
	}
	resp, err := req.Send(client)
	if err != nil {
		utils.Failf(t, "failed to send request: %s", err.Error())
	}

	gotInv := resp.Results[0]
	utils.Checkf(t, i.ID != gotInv.ID, "wanted invocation id %s; got %s", i.ID, gotInv.ID)
	utils.Checkf(t, i.Method.Prefix != gotInv.Method.Prefix, "wanted prefix %s; got %s", i.Method.Prefix, gotInv.Method.Prefix)
	gotInboxID := gotInv.Method.Args.IDs[0]
	wantInboxID := os.Getenv("FASTMAIL_INBOX_ID")
	utils.Checkf(t, wantInboxID != gotInboxID, "wanted inbox id %s; got %s", wantInboxID, gotInboxID)
}
