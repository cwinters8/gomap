package requests_test

import (
	"encoding/json"
	"testing"

	"github.com/cwinters8/gomap/requests"
	"github.com/cwinters8/gomap/requests/arguments"
	"github.com/cwinters8/gomap/utils"
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
	inv := requests.Invocation[arguments.Query]{
		ID:     "xyz",
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

	utils.Checkf(t, inv.ID != i.ID, "wanted ID %s; got ID %s", inv.ID, i.ID)
	utils.Checkf(t, inv.Method.Prefix != i.Method.Prefix, "wanted method prefix %s; got method prefix %s", inv.Method.Prefix, i.Method.Prefix)
	utils.Checkf(t, inv.Method.Args.AccountID != i.Method.Args.AccountID, "wanted account ID %s; got account ID %s", inv.Method.Args.AccountID, i.Method.Args.AccountID)
	utils.Checkf(t, inv.Method.Args.Filter.Name != i.Method.Args.Filter.Name, "wanted name filter %s; got name filter %s", inv.Method.Args.Filter.Name, i.Method.Args.Filter.Name)
}
