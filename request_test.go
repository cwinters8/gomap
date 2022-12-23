package gomap_test

import (
	"os"
	"testing"

	"github.com/cwinters8/gomap"
	"github.com/cwinters8/gomap/arguments"
	"github.com/cwinters8/gomap/utils"
)

func TestSendRequest(t *testing.T) {
	utils.Env(t)
	client, err := gomap.NewClient(os.Getenv("FASTMAIL_SESSION_URL"), os.Getenv("FASTMAIL_TOKEN"))
	if err != nil {
		utils.Failf(t, "failed to instantiate a new client: %s", err.Error())
	}
	i := gomap.Invocation[arguments.Query]{
		ID: "xyz",
		Method: &gomap.Method[arguments.Query]{
			Prefix: "Mailbox",
			Args: arguments.Query{
				AccountID: client.Session.PrimaryAccounts.Mail,
				Filter: arguments.Filter{
					Name: "Inbox",
				},
			},
		},
	}
	req := gomap.Request[arguments.Query]{
		Using: []gomap.Capability{
			gomap.UsingCore, // this maybe should just be used by default - I can't think of a case when Core would not be needed
			gomap.UsingMail,
		},
		Calls: []*gomap.Invocation[arguments.Query]{&i},
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
