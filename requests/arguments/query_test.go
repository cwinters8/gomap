package arguments_test

import (
	"encoding/json"
	"testing"

	"github.com/cwinters8/gomap/requests/arguments"
	"github.com/cwinters8/gomap/utils"
)

// this test may not be necessary, because custom marshal/unmarshal is not needed for arguments.Query
func TestQueryJSON(t *testing.T) {
	q := arguments.Query{
		AccountID: "xyz",
		Filter: arguments.Filter{
			Name: "felix the cat",
		},
	}
	b, err := json.Marshal(q)
	if err != nil {
		t.Fatalf("failed to marshal query to json: %s", err.Error())
	}
	var got arguments.Query
	if err := json.Unmarshal(b, &got); err != nil {
		t.Fatalf("failed to unmarshal json to query: %s", err.Error())
	}
	cases := []*utils.Case{{
		Check:  q.AccountID != got.AccountID,
		Format: "wanted account id %s; got %s",
		Args:   []any{q.AccountID, got.AccountID},
	}, {
		Check:  q.Filter.Name != got.Filter.Name,
		Format: "wanted name %s; got %s",
		Args:   []any{q.Filter.Name, got.Filter.Name},
	}}
	for _, c := range cases {
		if c.Check {
			t.Errorf(c.Format, c.Args...)
		}
	}
}
