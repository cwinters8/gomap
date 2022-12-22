package gomap_test

import (
	"encoding/json"
	"testing"

	"github.com/cwinters8/gomap"
)

func TestInvocationJSON(t *testing.T) {
	query := gomap.MailboxQuery{
		AccountID: "xyz",
		Filter: gomap.Filter{
			Name: "felix the cat",
		},
	}
	inv := gomap.Invocation[*gomap.MailboxQuery]{
		ID:     "query",
		Method: query.GetMethod(),
		Args:   &query,
	}
	b, err := json.Marshal(inv)
	if err != nil {
		failf(t, "failed to marshal invocation to json: %s", err.Error())
	}
	var i gomap.Invocation[*gomap.MailboxQuery]
	if err := json.Unmarshal(b, &i); err != nil {
		failf(t, "failed to unmarshal invocation: %s", err.Error())
	}
	if i.Args == nil {
		failf(t, "args must not be nil\nvalue of args: %v", i.Args)
	}
	switch {
	case inv.ID != i.ID:
		t.Errorf("wanted ID %s; got ID %s", inv.ID, i.ID)
	case inv.Method != i.Method:
		t.Errorf("wanted method %s; got method %s", inv.Method, i.Method)
	case inv.Args.AccountID != i.Args.AccountID:
		t.Errorf("wanted account ID %s; got account ID %s", inv.Args.AccountID, i.Args.AccountID)
	case inv.Args.Filter.Name != i.Args.Filter.Name:
		t.Errorf("wanted name filter %s; got name filter %s", inv.Args.Filter.Name, i.Args.Filter.Name)
	}
}
