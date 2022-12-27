package arguments_test

import (
	"encoding/json"
	"testing"

	"github.com/cwinters8/gomap/arguments"
	"github.com/cwinters8/gomap/utils"

	"github.com/google/uuid"
)

func TestNewMessage(t *testing.T) {
	msg, err := arguments.NewMessage(
		[]string{"xyz-box"},
		&arguments.Address{
			Name:  "Clark the Gopher",
			Email: "dev@clarkwinters.com",
		},
		&arguments.Address{
			Name:  "Tester McSet",
			Email: "tester@clarkwinters.com",
		},
		"hello world!",
		":)",
	)
	if err != nil {
		t.Fatalf("failed to instantiate new message: %s", err.Error())
	}
	cases := []*utils.Case{{
		Check:  msg.ID == uuid.Nil,
		Format: "message ID must not be nil",
	}, {
		Check:  msg.Body.ID == uuid.Nil,
		Format: "message body ID must not be nil",
	}}
	for _, c := range cases {
		if c.Check {
			t.Errorf(c.Format, c.Args...)
		}
	}
}

func TestSetJSON(t *testing.T) {
	from := arguments.Address{
		Name:  "Clark",
		Email: "dev@clarkwinters.com",
	}
	to := arguments.Address{
		Name:  "gopher",
		Email: "tester@clarkwinters.com",
	}

	msg, err := arguments.NewMessage(
		[]string{"xyz-box"},
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
	s := arguments.Set{
		AccountID: "xyz",
		Create:    msg,
	}
	b, err := json.Marshal(s)
	if err != nil {
		t.Fatalf("failed to marshal set args to json: %s", err.Error())
	}
	var got map[string]any
	if err := json.Unmarshal(b, &got); err != nil {
		t.Fatalf("failed to unmarshal json to set args: %s", err.Error())
	}
	cases := []*utils.Case{{
		Check:  s.AccountID != got["accountId"],
		Format: "wanted account id %s; got %s",
		Args:   []any{s.AccountID, got["accountId"]},
	}}
	wantMailbox := s.Create.MailboxIDs[0]
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
	gotFromAddr, ok := gotEmail["from"].([]any)
	if !ok {
		t.Fatalf("failed to coerce from addresses. %s", utils.Describe(gotEmail["from"]))
	}
	gotFrom, ok := gotFromAddr[0].(map[string]any)
	if !ok {
		t.Fatalf("failed to coerce from address. %s", utils.Describe(gotFromAddr[0]))
	}
	cases = append(cases, []*utils.Case{{
		Check:  from.Name != gotFrom["name"],
		Format: "wanted from name %s; got %s",
		Args:   []any{from.Name, gotFrom["name"]},
	}, {
		Check:  from.Email != gotFrom["email"],
		Format: "wanted from email %s; got %s",
		Args:   []any{from.Email, gotFrom["email"]},
	}}...)
	gotToAddr, ok := gotEmail["to"].([]any)
	if !ok {
		t.Fatalf("failed to coerce to addresses. %s", utils.Describe(gotEmail["to"]))
	}
	gotTo := gotToAddr[0].(map[string]any)
	cases = append(cases, []*utils.Case{{
		Check:  to.Name != gotTo["name"],
		Format: "wanted to name %s; got %s",
		Args:   []any{to.Name, gotTo["name"]},
	}, {
		Check:  to.Email != gotTo["email"],
		Format: "wanted to email %s; got %s",
		Args:   []any{to.Email, gotTo["email"]},
	}, {
		Check:  s.Create.Subject != gotEmail["subject"],
		Format: "wanted subject %s; got %s",
		Args:   []any{s.Create.Subject, gotEmail["subject"]},
	}}...)
	gotBodyStructure, ok := gotEmail["bodyStructure"].(map[string]any)
	if !ok {
		t.Fatalf("failed to coerce body structure. %s", utils.Describe(gotEmail["bodyStructure"]))
	}
	wantBodyIDStr := s.Create.Body.ID.String()
	cases = append(cases, &utils.Case{
		Check:  wantBodyIDStr != gotBodyStructure["partId"],
		Format: "wanted body ID %s; got %s",
		Args:   []any{wantBodyIDStr, gotBodyStructure["partId"]},
	})
	bType, ok := gotBodyStructure["type"].(string)
	if !ok {
		t.Fatalf("failed to coerce body type. %s", utils.Describe(gotBodyStructure["type"]))
	}
	gotBodyType := arguments.BodyType(bType)
	cases = append(cases, &utils.Case{
		Check:  s.Create.Body.Type != gotBodyType,
		Format: "wanted body type %s; got %s",
		Args:   []any{s.Create.Body.Type, gotBodyType},
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
		Check:  s.Create.Body.Value != values["value"],
		Format: "wanted body value %s; got %s",
		Args:   []any{s.Create.Body.Value, values["value"]},
	})

	// evaluate cases
	for _, c := range cases {
		if c.Check {
			t.Errorf(c.Format, c.Args...)
		}
	}
}
