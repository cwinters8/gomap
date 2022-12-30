package results_test

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/cwinters8/gomap/client"
	"github.com/cwinters8/gomap/mail"
	"github.com/cwinters8/gomap/methods"
	"github.com/cwinters8/gomap/requests"
	"github.com/cwinters8/gomap/results"
	"github.com/cwinters8/gomap/utils"
	"github.com/google/uuid"
)

func TestResultsJSON(t *testing.T) {
	b, err := os.ReadFile("testdata/set_error_response.json")
	if err != nil {
		t.Fatalf("failed to read set_error_response.json: %s", err.Error())
	}
	var r results.Results
	if err := json.Unmarshal(b, &r); err != nil {
		t.Fatalf("failed to unmarshal results: %s", err.Error())
	}
	s, ok := r.Results[0].(*results.Set)
	if !ok {
		t.Fatalf("failed to cast result to Set. %s", utils.Describe(r.Results[0]))
	}
	wantNotCreatedID := "4f17fc1f-f68a-4d5a-84a8-13e20bdaef1a"
	wantMethod := "Email/set"
	gotMethod, err := s.Method()
	if err != nil {
		t.Fatalf("failed to get set method: %s", err.Error())
	}
	wantID := "c5de7ad2-889e-44c3-a573-917d551dc856"
	errID, err := uuid.Parse("5112cf7a-d596-4dbd-9f34-456873aea3ef")
	if err != nil {
		t.Fatalf("failed to parse uuid: %s", err.Error())
	}
	wantError := results.Error{
		ID:         errID,
		Type:       "invalidProperties",
		Properties: []string{"onSuccessUpdateEmail/#sub"},
	}
	cases := []*utils.Case{{
		Check:   len(r.Results) != 1,
		Message: "wanted 1 result; got %d",
		Args:    []any{len(r.Results)},
	}, {
		Check:   len(r.Errors) != 1,
		Message: "wanted 1 error; got %d",
		Args:    []any{len(r.Errors)},
	}, {
		Check:   s.Body.NotCreated.ID.String() != wantNotCreatedID,
		Message: "wanted not created ID %s; got %s",
		Args:    []any{wantNotCreatedID, s.Body.NotCreated.ID.String()},
	}, {
		Check:   gotMethod != wantMethod,
		Message: "wanted result method %s; got %s",
		Args:    []any{wantMethod, gotMethod},
	}, {
		Check:   s.ID.String() != wantID,
		Message: "wanted result ID %s; got %s",
		Args:    []any{wantID, s.ID.String()},
	}, {
		Check:   r.Errors[0].ID != wantError.ID,
		Message: "wanted error ID %s; got %s",
		Args:    []any{wantError.ID, r.Errors[0].ID},
	}, {
		Check:   r.Errors[0].Type != wantError.Type,
		Message: "wanted error type %s; got %s",
		Args:    []any{wantError.Type, r.Errors[0].Type},
	}, utils.NewCase(
		r.Errors[0].Properties[0] != wantError.Properties[0],
		"wanted error property %s; got %s",
		wantError.Properties[0], r.Errors[0].Properties[0],
	)}
	for _, c := range cases {
		if c.Check {
			t.Errorf(c.Message, c.Args...)
		}
	}
}

func TestRequestResults(t *testing.T) {
	if err := utils.Env("../.env"); err != nil {
		t.Fatalf("failed to load env: %s", err.Error())
	}
	c, err := client.NewClient(os.Getenv("FASTMAIL_SESSION_URL"), os.Getenv("FASTMAIL_TOKEN"))
	if err != nil {
		t.Fatalf("failed to instantiate a new client: %s", err.Error())
	}

	t.Run("query", func(t *testing.T) {
		q := methods.Query{
			AccountID: c.Session.PrimaryAccounts.Mail,
			Filter: methods.Filter{
				Name: "Inbox",
			},
		}
		i, err := requests.NewInvocation(q, "Mailbox", methods.QueryMethod)
		if err != nil {
			t.Fatalf("failed to instantiate new invocation: %s", err.Error())
		}
		resp, err := requests.NewRequest([]*requests.Invocation[methods.Query]{i}).Send(c)
		if err != nil {
			t.Fatalf("request failure: %s", err.Error())
		}
		got, ok := resp.Results[0].(*results.Query)
		if !ok {
			t.Fatalf("failed to case result to Query. %s", utils.Describe(resp.Results[0]))
		}
		wantBoxID, err := uuid.Parse(os.Getenv("FASTMAIL_INBOX_ID"))
		if err != nil {
			t.Fatalf("failed to parse inbox uuid: %s", err.Error())
		}
		cases := []*utils.Case{
			utils.NewCase(
				got.ID != i.ID,
				"wanted id %s; got %s",
				i.ID, got.ID,
			),
			utils.NewCase(
				got.Prefix != i.Method.Prefix,
				"wanted method prefix %s; got %s",
				i.Method.Prefix, got.Prefix,
			),
			utils.NewCase(
				got.Body.AccountID != q.AccountID,
				"wanted account id %s; got %s",
				q.AccountID, got.Body.AccountID,
			),
			utils.NewCase(
				got.Body.IDs[0] != wantBoxID,
				"wanted mailbox id %s; got %s",
				wantBoxID, got.Body.IDs[0],
			),
			utils.NewCase(
				got.Body.Filter.Name != q.Filter.Name,
				"wanted name filter %s; got %s",
				q.Filter.Name, got.Body.Filter.Name,
			),
			utils.NewCase(
				got.Body.Total != 1,
				"wanted total 1; got %d",
				got.Body.Total,
			),
		}
		for _, c := range cases {
			if c.Check {
				t.Error(c.Message)
			}
		}
	})

	t.Run("set", func(t *testing.T) {
		box, err := mail.NewMailbox(c, "Drafts")
		if err != nil {
			t.Fatalf("failed to instantiate new mailbox: %s", err.Error())
		}
		msg, err := methods.NewMessage(
			methods.Mailboxes{box.ID},
			&methods.Address{
				Name:  "Gopher Clark",
				Email: "dev@clarkwinters.com",
			}, &methods.Address{
				Name:  "Setter Tester",
				Email: "tester@clarkwinters.com",
			},
			"hope this works",
			"trying to parse result of set request to json",
		)
		if err != nil {
			t.Fatalf("failed to instantiate new message: %s", err.Error())
		}
		set := methods.Set{
			AccountID: c.Session.PrimaryAccounts.Mail,
			Create:    msg,
		}
		i, err := requests.NewInvocation(set, "Email", methods.SetMethod)
		if err != nil {
			t.Fatalf("failed to instantiate new invocation: %s", err.Error())
		}
		req := requests.NewRequest([]*requests.Invocation[methods.Set]{i})
		got, err := req.Send(c)
		if err != nil {
			t.Fatalf("request failed: %s", err.Error())
		}
		if got == nil {
			t.Fatalf("response is nil")
		}
		if len(got.Errors) > 0 {
			errs := []results.Error{}
			for _, e := range got.Errors {
				errs = append(errs, *e)
			}
			t.Fatalf("found errors: %v", errs)
		}
		s, ok := got.Results[0].(*results.Set)
		if !ok {
			t.Fatalf("failed to cast result to Set. %s", utils.Describe(got.Results[0]))
		}
		if s.Body.Created == nil {
			t.Fatalf("created field in body is nil")
		}
		cases := []*utils.Case{
			utils.NewCase(
				s.Body.AccountID != set.AccountID,
				"wanted account id %s; got %s",
				set.AccountID, s.Body.AccountID,
			),
			utils.NewCase(
				len(s.Body.Created.ID) != 16,
				"wanted created id to be a valid 16 character uuid; got %s",
				s.Body.Created.ID,
			),
			utils.NewCase(
				s.Body.NotCreated != nil,
				"wanted NotCreated to be nil; got %v",
				s.Body.NotCreated,
			),
		}
		for _, c := range cases {
			if c.Check {
				t.Error(c.Message)
			}
		}
	})
}
