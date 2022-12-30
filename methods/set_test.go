package methods_test

import (
	"encoding/json"
	"testing"

	"github.com/cwinters8/gomap/methods"
	"github.com/cwinters8/gomap/utils"

	"github.com/google/uuid"
)

func TestNewMessage(t *testing.T) {
	msg, err := methods.NewMessage(
		[]string{"xyz-box"},
		&methods.Address{
			Name:  "Clark the Gopher",
			Email: "dev@clarkwinters.com",
		},
		&methods.Address{
			Name:  "Tester McSet",
			Email: "tester@clarkwinters.com",
		},
		"hello world!",
		":)",
	)
	if err != nil {
		t.Fatalf("failed to instantiate new message: %s", err.Error())
	}
	if msg.Keywords == nil {
		t.Fatalf("Keywords must not be nil")
	}
	cases := []*utils.Case{{
		Check:   !msg.Keywords.Seen,
		Message: "$seen keyword value should be true",
	}, {
		Check:   !msg.Keywords.Draft,
		Message: "$draft keyword value should be true",
	}, {
		Check:   msg.ID == uuid.Nil,
		Message: "message ID must not be nil",
	}, {
		Check:   msg.Body.ID == uuid.Nil,
		Message: "message body ID must not be nil",
	}}
	for _, c := range cases {
		if c.Check {
			t.Errorf(c.Message, c.Args...)
		}
	}
}

func TestSetJSON(t *testing.T) {
	// create Set data
	from := methods.Address{
		Name:  "Clark",
		Email: "dev@clarkwinters.com",
	}
	to := methods.Address{
		Name:  "gopher",
		Email: "tester@clarkwinters.com",
	}

	boxName := "xyz-box"
	msg, err := methods.NewMessage(
		[]string{boxName},
		&from,
		&to,
		"the meaning of life, the universe, and everything",
		"42",
	)
	if err != nil {
		t.Fatalf("failed to instantiate new message: %s", err.Error())
	}
	id, err := uuid.NewRandom()
	if err != nil {
		t.Fatalf("failed to generate new uuid: %s", err.Error())
	}
	msg.ID = id
	bodyID, err := uuid.NewRandom()
	if err != nil {
		t.Fatalf("failed to generate new uuid: %s", err.Error())
	}
	if err := msg.Body.SetID(bodyID); err != nil {
		t.Fatalf("failed to set body ID: %s", err.Error())
	}
	s := methods.Set{
		AccountID: "xyz",
		Create:    msg,
	}
	b, err := json.Marshal(s)
	if err != nil {
		t.Fatalf("failed to marshal set args to json: %s", err.Error())
	}

	t.Run("marshal", func(t *testing.T) {
		var got map[string]any
		if err := json.Unmarshal(b, &got); err != nil {
			t.Fatalf("failed to unmarshal json to set args: %s", err.Error())
		}
		cases := []*utils.Case{{
			Check:   s.AccountID != got["accountId"],
			Message: "wanted account id %s; got %s",
			Args:    []any{s.AccountID, got["accountId"]},
		}}
		boxIDs := *s.Create.MailboxIDs
		wantMailbox := boxIDs[0]
		gotCreate, ok := got["create"].(map[string]any)
		if !ok {
			t.Fatalf("failed to coerce create map. %s", utils.Describe(got["create"]))
		}
		var msgID string
		for k := range gotCreate {
			if k == s.Create.ID.String() {
				msgID = k
				break
			}
		}
		if len(msgID) < 1 {
			t.Fatalf("message ID %s not found", s.Create.ID.String())
		}
		gotEmail, ok := gotCreate[msgID].(map[string]any)
		if !ok {
			t.Fatalf("failed to coerce email map. %s", utils.Describe(gotCreate[msgID]))
		}
		gotMailboxes, ok := gotEmail["mailboxIds"].(map[string]any)
		if !ok {
			t.Fatalf("failed to coerce mailbox ids map. %s", utils.Describe(gotEmail["mailboxIds"]))
		}
		var boxID string
		for k, v := range gotMailboxes {
			value, ok := v.(bool)
			if !ok {
				t.Fatalf("failed to coerce mailbox value to bool. %s", utils.Describe(v))
			}
			if k == wantMailbox && value {
				boxID = k
				break
			}
		}
		if len(boxID) < 1 {
			t.Fatalf("mailbox ID %s not found", wantMailbox)
		}
		// validate keywords
		gotKeywords, ok := gotEmail["keywords"].(map[string]any)
		if !ok {
			t.Fatalf("failed to coerce keywords to map. %s", utils.Describe(gotEmail["keywords"]))
		}
		gotSeen, ok := gotKeywords["$seen"].(bool)
		if !ok {
			t.Fatalf("failed to coerce seen keyword to bool. %s", utils.Describe(gotKeywords["$seen"]))
		}
		gotDraft, ok := gotKeywords["$draft"].(bool)
		if !ok {
			t.Fatalf("failed to coerce draft keyword to bool. %s", utils.Describe(gotKeywords["$draft"]))
		}
		cases = append(cases, []*utils.Case{{
			Check:   !gotSeen,
			Message: "wanted seen keyword value to be true",
		}, {
			Check:   !gotDraft,
			Message: "wanted draft keyword value to be true",
		}}...)

		gotFromAddr, ok := gotEmail["from"].([]any)
		if !ok {
			t.Fatalf("failed to coerce from addresses. %s", utils.Describe(gotEmail["from"]))
		}
		gotFrom, ok := gotFromAddr[0].(map[string]any)
		if !ok {
			t.Fatalf("failed to coerce from address. %s", utils.Describe(gotFromAddr[0]))
		}
		cases = append(cases, []*utils.Case{{
			Check:   from.Name != gotFrom["name"],
			Message: "wanted from name %s; got %s",
			Args:    []any{from.Name, gotFrom["name"]},
		}, {
			Check:   from.Email != gotFrom["email"],
			Message: "wanted from email %s; got %s",
			Args:    []any{from.Email, gotFrom["email"]},
		}}...)
		gotToAddr, ok := gotEmail["to"].([]any)
		if !ok {
			t.Fatalf("failed to coerce to addresses. %s", utils.Describe(gotEmail["to"]))
		}
		gotTo := gotToAddr[0].(map[string]any)
		cases = append(cases, []*utils.Case{{
			Check:   to.Name != gotTo["name"],
			Message: "wanted to name %s; got %s",
			Args:    []any{to.Name, gotTo["name"]},
		}, {
			Check:   to.Email != gotTo["email"],
			Message: "wanted to email %s; got %s",
			Args:    []any{to.Email, gotTo["email"]},
		}, {
			Check:   s.Create.Subject != gotEmail["subject"],
			Message: "wanted subject %s; got %s",
			Args:    []any{s.Create.Subject, gotEmail["subject"]},
		}}...)
		gotBodyStructure, ok := gotEmail["bodyStructure"].(map[string]any)
		if !ok {
			t.Fatalf("failed to coerce body structure. %s", utils.Describe(gotEmail))
		}
		wantBodyIDStr := s.Create.Body.ID.String()
		cases = append(cases, &utils.Case{
			Check:   wantBodyIDStr != gotBodyStructure["partId"],
			Message: "wanted body ID %s; got %s",
			Args:    []any{wantBodyIDStr, gotBodyStructure["partId"]},
		})
		bType, ok := gotBodyStructure["type"].(string)
		if !ok {
			t.Fatalf("failed to coerce body type. %s", utils.Describe(gotBodyStructure["type"]))
		}
		gotBodyType := methods.BodyType(bType)
		cases = append(cases, &utils.Case{
			Check:   s.Create.Body.Type != gotBodyType,
			Message: "wanted body type %s; got %s",
			Args:    []any{s.Create.Body.Type, gotBodyType},
		})
		gotBodyValues, ok := gotEmail["bodyValues"].(map[string]any)
		if !ok {
			t.Fatalf("failed to coerce body values. %s", utils.Describe(gotEmail["bodyValues"]))
		}
		var gotBodyID string
		for k := range gotBodyValues {
			if k == bodyID.String() {
				gotBodyID = k
				break
			}
		}
		if len(gotBodyID) < 1 {
			t.Fatalf("body ID %s not found", bodyID.String())
		}
		values, ok := gotBodyValues[gotBodyID].(map[string]any)
		if !ok {
			t.Fatalf("failed to coerce body value. %s", utils.Describe(gotBodyValues[gotBodyID]))
		}
		cases = append(cases, &utils.Case{
			Check:   s.Create.Body.Value != values["value"],
			Message: "wanted body value %s; got %s",
			Args:    []any{s.Create.Body.Value, values["value"]},
		})

		// evaluate cases
		for _, c := range cases {
			if c.Check {
				t.Errorf(c.Message, c.Args...)
			}
		}
	})
	t.Run("unmarshal", func(t *testing.T) {
		if len(b) < 1 {
			t.Fatalf("empty byte slice")
		}
		var gotSet methods.Set
		if err := json.Unmarshal(b, &gotSet); err != nil {
			t.Fatalf("failed to unmarshal set: %s", err.Error())
		}
		if gotSet.Create == nil {
			t.Fatalf("Set's Create message must not be nil")
		}
		gotBoxes := *gotSet.Create.MailboxIDs
		if len(gotBoxes) < 1 {
			t.Fatalf("no mailboxes 😢")
		}
		cases := []*utils.Case{{
			Check:   s.AccountID != gotSet.AccountID,
			Message: "wanted account id %s; got %s",
			Args:    []any{s.AccountID, gotSet.AccountID},
		}, {
			Check:   id != gotSet.Create.ID,
			Message: "wanted message id %v; got %v",
			Args:    []any{id, gotSet.Create.ID},
		}, {
			Check:   boxName != gotBoxes[0],
			Message: "wanted mailbox id %s; got %s",
			Args:    []any{boxName, gotBoxes[0]},
		}, {
			Check:   from.Name != gotSet.Create.From[0].Name,
			Message: "wanted from name %s; got %s",
			Args:    []any{from.Name, gotSet.Create.From[0].Name},
		}, {
			Check:   from.Email != gotSet.Create.From[0].Email,
			Message: "wanted from email %s; got %s",
			Args:    []any{from.Email, gotSet.Create.From[0].Email},
		}, {
			Check:   to.Name != gotSet.Create.To[0].Name,
			Message: "wanted to name %s; got %s",
			Args:    []any{to.Name, gotSet.Create.To[0].Name},
		}, {
			Check:   to.Email != gotSet.Create.To[0].Email,
			Message: "wanted to email %s; got %s",
			Args:    []any{to.Email, gotSet.Create.To[0].Email},
		}, {
			Check:   msg.Subject != gotSet.Create.Subject,
			Message: "wanted subject %s; got %s",
			Args:    []any{msg.Subject, gotSet.Create.Subject},
		}, {
			Check:   msg.Body.Value != gotSet.Create.Body.Value,
			Message: "wanted body value %s; got %s",
			Args:    []any{msg.Body.Value, gotSet.Create.Body.Value},
		}, {
			Check:   bodyID != gotSet.Create.Body.ID,
			Message: "wanted body id %v; got %v",
			Args:    []any{bodyID, gotSet.Create.Body.ID},
		}, {
			Check:   s.Create.Keywords.Seen != gotSet.Create.Keywords.Seen,
			Message: "wanted $seen boolean %t; got %t",
			Args:    []any{s.Create.Keywords.Seen, gotSet.Create.Keywords.Seen},
		}, {
			Check:   s.Create.Keywords.Draft != gotSet.Create.Keywords.Draft,
			Message: "wanted $draft boolean %t; got %t",
			Args:    []any{s.Create.Keywords.Draft, gotSet.Create.Keywords.Draft},
		}}

		for _, c := range cases {
			if c.Check {
				t.Errorf(c.Message, c.Args...)
			}
		}
	})
}

func TestNewBody(t *testing.T) {
	value := "hello"
	body, err := methods.NewBody(methods.TextPlain, value)
	if err != nil {
		t.Fatalf("failed to instantiate new body: %s", err.Error())
	}
	cases := []*utils.Case{{
		Check:   body.Type != methods.TextPlain,
		Message: "wanted body type %s; got %s",
		Args:    []any{methods.TextPlain, body.Type},
	}, {
		Check:   body.ID == uuid.Nil,
		Message: "body id must not be nil",
	}, {
		Check:   body.Value != value,
		Message: "wanted body value %s; got %s",
		Args:    []any{value, body.Value},
	}}
	for _, c := range cases {
		if c.Check {
			t.Errorf(c.Message, c.Args...)
		}
	}
}
