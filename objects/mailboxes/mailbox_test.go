package mailboxes_test

import (
	"testing"

	"github.com/cwinters8/gomap/objects/mailboxes"
	"github.com/cwinters8/gomap/utils"
	"github.com/google/uuid"
)

func TestQueryCall(t *testing.T) {
	m := mailboxes.Mailbox{
		Name: "Inbox",
	}
	acctID := "xyz"
	call, err := m.Query(acctID)
	if err != nil {
		t.Fatalf("failed to construct query call: %s", err.Error())
	}
	cases := utils.Cases{utils.NewCase(
		acctID != call.AccountID,
		"wanted account id %s; got %s",
		acctID, call.AccountID,
	), utils.NewCase(
		call.ID == uuid.Nil,
		"wanted call id to be non-nil",
	)}
	filter, ok := call.Arguments["filter"].(map[string]string)
	if !ok {
		t.Fatalf("failed to cast filter args to map. %s", utils.Describe(call.Arguments["filter"]))
	}
	cases.Append(utils.NewCase(
		m.Name != filter["name"],
		"wanted mailbox name %s; got %s",
		m.Name, filter["name"],
	))
	cases.Iterator(func(c *utils.Case) {
		t.Error(c.Message)
	})
}
