package requests_test

import (
	"encoding/json"
	"testing"

	"github.com/cwinters8/gomap/methods"
	"github.com/cwinters8/gomap/requests"
	"github.com/cwinters8/gomap/utils"

	"github.com/google/uuid"
)

func TestInvocationJSON(t *testing.T) {
	query := methods.Method[methods.Query]{
		Prefix: "Mailbox",
		Args: methods.Query{
			AccountID: "xyz",
			Filter: methods.Filter{
				Name: "felix the cat",
			},
		},
	}
	id, err := uuid.NewRandom()
	if err != nil {
		t.Fatalf("failed to generate new uuid: %s", err.Error())
	}
	inv := requests.Invocation[methods.Query]{
		ID:     id,
		Method: &query,
	}
	b, err := json.Marshal(inv)
	if err != nil {
		t.Fatalf("failed to marshal invocation to json: %s", err.Error())
	}
	var i requests.Invocation[methods.Query]
	if err := json.Unmarshal(b, &i); err != nil {
		t.Fatalf("failed to unmarshal invocation: %s", err.Error())
	}

	cases := []*utils.Case{{
		Check:   inv.ID != i.ID,
		Message: "wanted ID %s; got ID %s",
		Args:    []any{inv.ID, i.ID},
	}, {
		Check:   inv.Method.Prefix != i.Method.Prefix,
		Message: "wanted method prefix %s; got method prefix %s",
		Args:    []any{inv.Method.Prefix, i.Method.Prefix},
	}, {
		Check:   inv.Method.Args.AccountID != i.Method.Args.AccountID,
		Message: "wanted account ID %s; got account ID %s",
		Args:    []any{inv.Method.Args.AccountID, i.Method.Args.AccountID},
	}, {
		Check:   inv.Method.Args.Filter.Name != i.Method.Args.Filter.Name,
		Message: "wanted name filter %s; got name filter %s",
		Args:    []any{inv.Method.Args.Filter.Name, i.Method.Args.Filter.Name},
	}}
	for _, c := range cases {
		if c.Check {
			t.Errorf(c.Message, c.Args...)
		}
	}
}

func TestNewInvocation(t *testing.T) {
	prefix := "Email"
	i, err := requests.NewInvocation(methods.Set{}, prefix, methods.SetMethod)
	if err != nil {
		t.Fatalf("failed to instantiate new invocation: %s", err.Error())
	}
	fatals := []*utils.Case{{
		Check:   i == nil,
		Message: "invocation is nil",
	}, {
		Check:   i.Method == nil,
		Message: "method is nil",
	}}
	for _, c := range fatals {
		if c.Check {
			t.Fatalf(c.Message, c.Args...)
		}
	}

	cases := []*utils.Case{{
		Check:   len(i.ID) < 16,
		Message: "uuid slice must be 16 bytes long",
	}, {
		Check:   prefix != i.Method.Prefix,
		Message: "wanted method prefix %s; got %s",
		Args:    []any{prefix, i.Method.Prefix},
	}, {
		Check:   methods.SetMethod != i.Method.Type,
		Message: "wanted method type %s; got %s",
		Args:    []any{methods.SetMethod, i.Method.Type},
	}}
	for _, c := range cases {
		if c.Check {
			t.Errorf(c.Message, c.Args...)
		}
	}
}
