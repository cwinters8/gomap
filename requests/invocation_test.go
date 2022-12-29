package requests_test

import (
	"encoding/json"
	"testing"

	"github.com/cwinters8/gomap/requests"
	"github.com/cwinters8/gomap/requests/arguments"
	"github.com/cwinters8/gomap/utils"

	"github.com/google/uuid"
)

func TestInvocationJSON(t *testing.T) {
	query := requests.Method[arguments.Query]{
		Prefix: "Mailbox",
		Args: arguments.Query{
			AccountID: "xyz",
			Filter: arguments.Filter{
				Name: "felix the cat",
			},
		},
	}
	id, err := uuid.NewRandom()
	if err != nil {
		t.Fatalf("failed to generate new uuid: %s", err.Error())
	}
	inv := requests.Invocation[arguments.Query]{
		ID:     id,
		Method: &query,
	}
	b, err := json.Marshal(inv)
	if err != nil {
		t.Fatalf("failed to marshal invocation to json: %s", err.Error())
	}
	var i requests.Invocation[arguments.Query]
	if err := json.Unmarshal(b, &i); err != nil {
		t.Fatalf("failed to unmarshal invocation: %s", err.Error())
	}

	cases := []*utils.Case{{
		Check:  inv.ID != i.ID,
		Format: "wanted ID %s; got ID %s",
		Args:   []any{inv.ID, i.ID},
	}, {
		Check:  inv.Method.Prefix != i.Method.Prefix,
		Format: "wanted method prefix %s; got method prefix %s",
		Args:   []any{inv.Method.Prefix, i.Method.Prefix},
	}, {
		Check:  inv.Method.Args.AccountID != i.Method.Args.AccountID,
		Format: "wanted account ID %s; got account ID %s",
		Args:   []any{inv.Method.Args.AccountID, i.Method.Args.AccountID},
	}, {
		Check:  inv.Method.Args.Filter.Name != i.Method.Args.Filter.Name,
		Format: "wanted name filter %s; got name filter %s",
		Args:   []any{inv.Method.Args.Filter.Name, i.Method.Args.Filter.Name},
	}}
	for _, c := range cases {
		if c.Check {
			t.Errorf(c.Format, c.Args...)
		}
	}
}
