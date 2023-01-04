package emails_test

import (
	"testing"

	"github.com/cwinters8/gomap/objects/emails"
	"github.com/cwinters8/gomap/utils"
	"github.com/google/uuid"
)

func TestParse(t *testing.T) {
	rawBody := `
	{
		"latestClientVersion": "",
		"sessionState": "cyrus-0;p-8bad9e83f1",
		"methodResponses": [
			[
				"Email/get",
				{
					"state": "4154",
					"list": [
						{
							"mailboxIds": {
								"60b77041-ee8f-4429-aaf7-39b94d40c9eb": true
							},
							"bodyValues": {
								"1": {
									"isEncodingProblem": false,
									"value": "trying to parse result of set request to json",
									"isTruncated": false
								}
							},
							"to": [
								{
									"email": "tester@clarkwinters.com",
									"name": "Setter Tester"
								}
							],
							"from": [
								{
									"name": "Gopher Clark",
									"email": "dev@clarkwinters.com"
								}
							],
							"subject": "hope this works",
							"id": "M1cb24edb211ae50b4ed508ad"
						}
					],
					"accountId": "u69394015",
					"notFound": []
				},
				"6d420433-8733-4add-8cc6-0743d6a27a72"
			]
		]
	}
	`
	id, err := uuid.Parse("6d420433-8733-4add-8cc6-0743d6a27a72")
	if err != nil {
		t.Fatalf("failed to parse uuid: %s", err.Error())
	}
	want := emails.Email{
		ID:         "M1cb24edb211ae50b4ed508ad",
		RequestID:  id,
		MailboxIDs: []string{"60b77041-ee8f-4429-aaf7-39b94d40c9eb"},
		From: &emails.Address{
			Name:  "Gopher Clark",
			Email: "dev@clarkwinters.com",
		},
		To: []*emails.Address{{
			Name:  "Setter Tester",
			Email: "tester@clarkwinters.com",
		}},
		Subject: "hope this works",
		Body: &emails.Body{
			Value: "trying to parse result of set request to json",
		},
	}
	var result emails.Result
	if err := result.Parse(rawBody); err != nil {
		t.Fatalf("failed to parse result body: %s", err.Error())
	}
	if len(result.Body.List) < 1 {
		t.Fatalf("Body.List should not be empty. %s", utils.Describe(result))
	}
	got := result.Body.List[0]
	cases := utils.Cases{
		utils.NewCase(
			want.ID != got.ID,
			"wanted id %s; got %s",
			want.ID, got.ID,
		),
		utils.NewCase(
			want.MailboxIDs[0] != got.MailboxIDs[0],
			"wanted mailbox id %s; got %s",
			want.MailboxIDs[0], got.MailboxIDs[0],
		),
		utils.NewCase(
			want.From.Name != got.From.Name,
			"wanted from name %s; got %s",
			want.From.Name, got.From.Name,
		),
		utils.NewCase(
			want.From.Email != got.From.Email,
			"wanted from email %s; got %s",
			want.From.Email, got.From.Email,
		),
		utils.NewCase(
			want.To[0].Name != got.To[0].Name,
			"wanted to name %s; got %s",
			want.To[0].Name, got.To[0].Name,
		),
		utils.NewCase(
			want.To[0].Email != got.To[0].Email,
			"wanted to email %s; got %s",
			want.To[0].Email, got.To[0].Email,
		),
		utils.NewCase(
			want.Subject != got.Subject,
			"wanted subject %s; got %s",
			want.Subject, got.Subject,
		),
		utils.NewCase(
			want.Body.Value != got.Body.Value,
			"wanted body value %s; got %s",
			want.Body.Value, got.Body.Value,
		),
	}
	cases.Iterator(func(c *utils.Case) {
		t.Error(c.Message)
	})
}
