package arguments_test

import (
	"encoding/json"
	"testing"

	"github.com/cwinters8/gomap/arguments"
	"github.com/cwinters8/gomap/utils"
)

func TestQueryJSON(t *testing.T) {
	q := arguments.Query{
		AccountID: "xyz",
		Filter: arguments.Filter{
			Name: "felix the cat",
		},
	}
	b, err := json.Marshal(q)
	if err != nil {
		utils.Failf(t, "failed to marshal query to json: %s", err.Error())
	}
	var got arguments.Query
	if err := json.Unmarshal(b, &got); err != nil {
		utils.Failf(t, "failed to unmarshal json to query: %s", err.Error())
	}
	utils.Checkf(t, q.AccountID != got.AccountID, "wanted account id %s; got %s", q.AccountID, got.AccountID)
	utils.Checkf(t, q.Filter.Name != got.Filter.Name, "wanted name %s; got %s", q.Filter.Name, got.Filter.Name)
}
