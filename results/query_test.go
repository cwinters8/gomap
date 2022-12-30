package results_test

import (
	"testing"

	"github.com/cwinters8/gomap/results"
	"github.com/cwinters8/gomap/utils"
	"github.com/google/uuid"
)

func TestParseQuery(t *testing.T) {
	rawBody := `
	{
		"accountId": "u69394015",
		"canCalculateChanges": true,
		"filter": {
			"name": "Drafts"
		},
		"ids": [
			"60b77041-ee8f-4429-aaf7-39b94d40c9eb"
		],
		"position": 0,
		"queryState": "17",
		"total": 1
	}
	`
	var q results.Query
	if err := q.Parse(rawBody); err != nil {
		t.Fatalf("failed to parse query body: %s", err.Error())
	}
	id, err := uuid.Parse("60b77041-ee8f-4429-aaf7-39b94d40c9eb")
	if err != nil {
		t.Fatalf("failed to parse uuid: %s", err.Error())
	}
	want := results.QueryBody{
		AccountID: "u69394015",
		Filter: &results.Filter{
			Name: "Drafts",
		},
		IDs:   []uuid.UUID{id},
		Total: 1,
	}
	got := q.Body
	cases := []*utils.Case{
		utils.NewCase(
			got.AccountID != want.AccountID,
			"wanted account id %s; got %s",
			want.AccountID, got.AccountID,
		),
		utils.NewCase(
			got.Filter.Name != want.Filter.Name,
			"wanted name filter %s; got %s",
			want.Filter.Name, got.Filter.Name,
		),
		utils.NewCase(
			got.IDs[0] != want.IDs[0],
			"wanted id %s; got %s",
			want.IDs[0], got.IDs[0],
		),
		utils.NewCase(
			got.Total != want.Total,
			"wanted total %d; got %d",
			want.Total, got.Total,
		),
	}
	for _, c := range cases {
		if c.Check {
			t.Errorf(c.Message)
		}
	}
}
