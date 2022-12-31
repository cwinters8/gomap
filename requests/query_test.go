package requests_test

import (
	"encoding/json"
	"testing"

	"github.com/cwinters8/gomap/utils"
	"github.com/google/uuid"

	"github.com/cwinters8/gomap/requests"
)

func TestQueryJSON(t *testing.T) {
	want := requests.Query{
		Prefix: "Mailbox",
		Body: &requests.QueryBody{
			AccountID: "xyz",
			Filter: &requests.Filter{
				Name: "tester",
			},
		},
	}
	q, err := requests.NewQuery(want.Body.AccountID, want.Prefix, want.Body.Filter.Name)
	if err != nil {
		t.Fatalf("failed to instantiate new query: %s", err.Error())
	}
	b, err := json.Marshal(q)
	if err != nil {
		t.Fatalf("failed to marshal query to json: %s", err.Error())
	}
	var raw [3]any
	if err := json.Unmarshal(b, &raw); err != nil {
		t.Fatalf("failed to unmarshal raw query to array: %s", err.Error())
	}
	method, ok := raw[0].(string)
	if !ok {
		t.Fatalf("failed to cast method to string. %s", utils.Describe(raw[0]))
	}
	wantMethod, err := want.Method()
	if err != nil {
		t.Fatalf("failed to get wanted method: %s", err.Error())
	}
	body, ok := raw[1].(map[string]any)
	if !ok {
		t.Fatalf("failed to cast body to map. %s", utils.Describe(raw[1]))
	}
	filter, ok := body["filter"].(map[string]any)
	if !ok {
		t.Fatalf("failed to cast filter to map. %s", utils.Describe(body["filter"]))
	}
	id, ok := raw[2].(string)
	if !ok {
		t.Fatalf("failed to cast id to string. %s", utils.Describe(raw[2]))
	}
	if _, err := uuid.Parse(id); err != nil {
		t.Errorf("failed to parse id: %s", err.Error())
	}
	cases := utils.Cases{
		utils.NewCase(
			method != wantMethod,
			"wanted method %s; got %s",
			wantMethod, method,
		),
		utils.NewCase(
			want.Body.AccountID != body["accountId"],
			"wanted account id %s; got %s",
			want.Body.AccountID, body["accountId"],
		),
		utils.NewCase(
			want.Body.Filter.Name != filter["name"],
			"wanted name filter %s; got %s",
			want.Body.Filter.Name, filter["name"],
		),
	}
	cases.Iterator(func(c *utils.Case) {
		t.Error(c.Message)
	})
}
