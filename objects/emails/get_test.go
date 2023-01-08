package emails_test

import (
	"os"
	"testing"

	"github.com/cwinters8/gomap/client"
	"github.com/cwinters8/gomap/objects/emails"
	"github.com/cwinters8/gomap/requests"
	"github.com/cwinters8/gomap/utils"
)

func TestEmailGet(t *testing.T) {
	if err := utils.Env(envPath); err != nil {
		t.Fatalf("failed to source env variables from path `%s`: %s", envPath, err.Error())
	}
	c, err := client.NewClient(os.Getenv("FASTMAIL_SESSION_URL"), os.Getenv("FASTMAIL_TOKEN"))
	if err != nil {
		t.Fatalf("failed to construct new client: %s", err.Error())
	}
	t.Run("single email", func(t *testing.T) {
		e := emails.Email{
			ID: "M1cb24edb211ae50b4ed508ad",
		}
		call, err := e.Get(c.Session.PrimaryAccounts.Mail)
		if err != nil {
			t.Fatalf("failed to construct new call for Email/get: %s", err.Error())
		}
		if _, err := requests.Request(c, []*requests.Call{call}, false); err != nil {
			t.Fatalf("Email/get request failure: %s", err.Error())
		}
		want := emails.Email{
			MailboxIDs: []string{"60b77041-ee8f-4429-aaf7-39b94d40c9eb"},
			From: []*emails.Address{{
				Name:  "Gopher Clark",
				Email: "dev@clarkwinters.com",
			}},
			To: []*emails.Address{{
				Name:  "Setter Tester",
				Email: "tester@clarkwinters.com",
			}},
			Subject: "hope this works",
			Body: &emails.Body{
				Value: "trying to parse result of set request to json",
				Type:  emails.TextPlain,
			},
		}

		// if 1 or more of these cases are true, test should immediately fail
		fatals := utils.Cases{
			utils.NewCase(
				len(e.MailboxIDs) < 1,
				"MailboxIDs should contain at least 1 id",
			),
			utils.NewCase(
				len(e.From) < 1,
				"From should contain at least 1 address",
			),
			utils.NewCase(
				len(e.To) < 1,
				"To should contain at least 1 address",
			),
			utils.NewCase(
				e.Body == nil,
				"Body should not be nil",
			),
		}
		failed := 0
		fatals.Iterator(func(c *utils.Case) {
			t.Error(c.Message)
			failed++
		})
		if failed > 0 {
			t.FailNow()
		}

		cases := utils.Cases{
			utils.NewCase(
				want.MailboxIDs[0] != e.MailboxIDs[0],
				"wanted mailbox id %s; got %s",
				want.MailboxIDs[0], e.MailboxIDs[0],
			),
			utils.NewCase(
				want.From[0].Name != e.From[0].Name,
				"wanted from address name %s; got %s",
				want.From[0].Name, e.From[0].Name,
			),
			utils.NewCase(
				want.From[0].Email != e.From[0].Email,
				"wanted from address email %s; got %s",
				want.From[0].Email, e.From[0].Email,
			),
			utils.NewCase(
				want.To[0].Name != e.To[0].Name,
				"wanted to address name %s; got %s",
				want.To[0].Name, e.To[0].Name,
			),
			utils.NewCase(
				want.To[0].Email != e.To[0].Email,
				"wanted to address email %s; got %s",
				want.To[0].Email, e.To[0].Email,
			),
			utils.NewCase(
				want.Body.Value != e.Body.Value,
				"wanted body value %s; got %s",
				want.Body.Value, e.Body.Value,
			),
			utils.NewCase(
				want.Body.Type != e.Body.Type,
				"wanted body type %s; got %s",
				want.Body.Type, e.Body.Type,
			),
		}
		cases.Iterator(func(c *utils.Case) {
			t.Error(c.Message)
		})
	})

	t.Run("multiple emails", func(t *testing.T) {
		hopeThisWorks := "M1cb24edb211ae50b4ed508ad"
		testingSet := "Mf59d21cfc25fbde3d9676e7d"
		emailIDs := []string{hopeThisWorks, testingSet}
		emailList, notFound, err := emails.GetEmails(c, emailIDs)
		if err != nil {
			t.Fatalf("failed to get emails: %s", err.Error())
		}
		if len(notFound) > 0 {
			t.Errorf("some email IDs were not found: %v", notFound)
		}
		wantFrom := []*emails.Address{{
			Name:  "Gopher Clark",
			Email: "dev@clarkwinters.com",
		}}
		wantToAddr := "tester@clarkwinters.com"
		wantBoxIDs := []string{"60b77041-ee8f-4429-aaf7-39b94d40c9eb"}
		cases := utils.Cases{}
		for _, email := range emailList {
			var want emails.Email
			switch email.ID {
			case hopeThisWorks:
				want = emails.Email{
					ID:         hopeThisWorks,
					MailboxIDs: wantBoxIDs,
					From:       wantFrom,
					To: []*emails.Address{{
						Name:  "Setter Tester",
						Email: wantToAddr,
					}},
					Subject: "hope this works",
					Body: &emails.Body{
						Type:  emails.TextPlain,
						Value: "trying to parse result of set request to json",
					},
				}
			case testingSet:
				want = emails.Email{
					ID:         testingSet,
					MailboxIDs: wantBoxIDs,
					From:       wantFrom,
					To: []*emails.Address{{
						Name:  "Tester McSeterson",
						Email: wantToAddr,
					}},
					Subject: "testing Email/set",
					Body: &emails.Body{
						Type:  emails.TextPlain,
						Value: "hello from TestEmailSet!",
					},
				}
			}
			cases.Append(
				utils.NewCase(
					want.MailboxIDs[0] != email.MailboxIDs[0],
					"wanted mailbox id %s; got %s",
					want.MailboxIDs[0], email.MailboxIDs[0],
				),
				utils.NewCase(
					want.From[0].Name != email.From[0].Name,
					"wanted from name %s; got %s",
					want.From[0].Name, email.From[0].Name,
				),
				utils.NewCase(
					want.From[0].Email != email.From[0].Email,
					"wanted from email %s; got %s",
					want.From[0].Email, email.From[0].Email,
				),
				utils.NewCase(
					want.To[0].Name != email.To[0].Name,
					"wanted to name %s; got %s",
					want.To[0].Name, email.To[0].Name,
				),
				utils.NewCase(
					want.To[0].Email != email.To[0].Email,
					"wanted to email %s; got %s",
					want.To[0].Email, email.To[0].Email,
				),
				utils.NewCase(
					want.Subject != email.Subject,
					"wanted subject %s; got %s",
					want.Subject, email.Subject,
				),
				utils.NewCase(
					want.Body.Type != email.Body.Type,
					"wanted body type %s; got %s",
					want.Body.Type, email.Body.Type,
				),
				utils.NewCase(
					want.Body.Value != email.Body.Value,
					"wanted body value %s; got %s",
					want.Body.Value, email.Body.Value,
				),
			)
		}
		cases.Iterator(func(c *utils.Case) {
			t.Error(c.Message)
		})
	})
}
