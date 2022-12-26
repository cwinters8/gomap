package arguments_test

import (
	"encoding/json"
	"testing"

	"github.com/cwinters8/gomap/arguments"
	"github.com/cwinters8/gomap/utils"
)

func TestSetJSON(t *testing.T) {
	from := arguments.Address{
		Name:  "Clark",
		Email: "dev@clarkwinters.com",
	}
	to := arguments.Address{
		Name:  "gopher",
		Email: "tester@clarkwinters.com",
	}
	bodyID := "xyz-body"
	s := arguments.Set{
		AccountID: "xyz",
		Create: &arguments.Message{
			ID:         "xyz-message",
			MailboxIDs: []string{"xyz-box"},
			From:       []*arguments.Address{&from},
			To:         []*arguments.Address{&to},
			Subject:    "the meaning of life, the universe, and everything",
			BodyStructure: &arguments.BodyStructure{
				ID:   bodyID,
				Type: arguments.TextPlain,
			},
			BodyValue: &arguments.BodyValue{
				ID:    bodyID,
				Value: "42",
			},
		},
	}
	b, err := json.Marshal(s)
	if err != nil {
		t.Fatalf("failed to marshal set args to json: %s", err.Error())
	}
	var got map[string]any
	if err := json.Unmarshal(b, &got); err != nil {
		t.Fatalf("failed to unmarshal json to set args: %s", err.Error())
	}
	utils.Checkf(t, s.AccountID != got["accountId"], "wanted account id %s; got %s", s.AccountID, got["accountId"])
	wantMailbox := s.Create.MailboxIDs[0]
	gotCreate, ok := got["create"].(map[string]any)
	if !ok {
		t.Fatalf("failed to coerce create map. %s", utils.Describe(got["create"]))
	}
	var msgID string
	for k := range gotCreate {
		if k == s.Create.ID {
			msgID = k
			break
		}
	}
	if len(msgID) < 1 {
		t.Fatalf("message ID %s not found", s.Create.ID)
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
	utils.Checkf(t, from.Name != gotFrom["name"], "wanted from name %s; got %s", from.Name, gotFrom["name"])
	utils.Checkf(t, from.Email != gotFrom["email"], "wanted from name %s; got %s", from.Email, gotFrom["email"])
	gotToAddr, ok := gotEmail["to"].([]any)
	if !ok {
		t.Fatalf("failed to coerce to addresses. %s", utils.Describe(gotEmail["to"]))
	}
	gotTo := gotToAddr[0].(map[string]any)
	utils.Checkf(t, to.Name != gotTo["name"], "wanted to name %s; got %s", to.Name, gotTo["name"])
	utils.Checkf(t, to.Email != gotTo["email"], "wanted to email %s; got %s", to.Email, gotTo["email"])
	utils.Checkf(t, s.Create.Subject != gotEmail["subject"], "wanted subject %s; got %s", s.Create.Subject, gotEmail["subject"])
	gotBodyStructure, ok := gotEmail["bodyStructure"].(map[string]any)
	if !ok {
		t.Fatalf("failed to coerce body structure. %s", utils.Describe(gotEmail["bodyStructure"]))
	}
	utils.Checkf(t, s.Create.BodyStructure.ID != gotBodyStructure["partId"], "wanted body ID %s; got %s", s.Create.BodyStructure.ID, gotBodyStructure["partId"])
	bType, ok := gotBodyStructure["type"].(string)
	if !ok {
		t.Fatalf("failed to coerce body type. %s", utils.Describe(gotBodyStructure["type"]))
	}
	gotBodyType := arguments.BodyType(bType)
	utils.Checkf(t, s.Create.BodyStructure.Type != gotBodyType, "wanted gotBodyStructure type %s; got %s", s.Create.BodyStructure.Type, gotBodyType)
	gotBodyValues, ok := gotEmail["bodyValues"].(map[string]any)
	if !ok {
		t.Fatalf("failed to coerce body values. %s", utils.Describe(gotEmail["bodyValues"]))
	}
	var gotBodyID string
	for k := range gotBodyValues {
		if k == bodyID {
			gotBodyID = k
			break
		}
	}
	if len(gotBodyID) < 1 {
		t.Fatalf("body ID %s not found", bodyID)
	}
	values, ok := gotBodyValues[gotBodyID].(map[string]any)
	if !ok {
		t.Fatalf("failed to coerce body value. %s", utils.Describe(gotBodyValues[gotBodyID]))
	}
	utils.Checkf(t, s.Create.BodyValue.Value != values["value"], "wanted body value %s; got %s", s.Create.BodyValue.Value, values["value"])
}
