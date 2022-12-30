package results_test

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/cwinters8/gomap/results"
	"github.com/cwinters8/gomap/utils"
	"github.com/google/uuid"
)

func TestResultsJSON(t *testing.T) {
	b, err := os.ReadFile("testdata/set_error_response.json")
	if err != nil {
		t.Fatalf("failed to read set_error_response.json: %s", err.Error())
	}
	var r results.Results
	if err := json.Unmarshal(b, &r); err != nil {
		t.Fatalf("failed to unmarshal results: %s", err.Error())
	}
	s, ok := r.Results[0].(*results.Set)
	if !ok {
		t.Fatalf("failed to cast result to Set. %s", utils.Describe(r.Results[0]))
	}
	wantNotCreatedID := "4f17fc1f-f68a-4d5a-84a8-13e20bdaef1a"
	wantMethod := "Email/set"
	gotMethod, err := s.Method()
	if err != nil {
		t.Fatalf("failed to get set method: %s", err.Error())
	}
	wantID := "c5de7ad2-889e-44c3-a573-917d551dc856"
	errID, err := uuid.Parse("5112cf7a-d596-4dbd-9f34-456873aea3ef")
	if err != nil {
		t.Fatalf("failed to parse uuid: %s", err.Error())
	}
	wantError := results.Error{
		ID:         errID,
		Type:       "invalidProperties",
		Properties: []string{"onSuccessUpdateEmail/#sub"},
	}
	cases := []*utils.Case{{
		Check:   len(r.Results) != 1,
		Message: "wanted 1 result; got %d",
		Args:    []any{len(r.Results)},
	}, {
		Check:   len(r.Errors) != 1,
		Message: "wanted 1 error; got %d",
		Args:    []any{len(r.Errors)},
	}, {
		Check:   s.Body.NotCreated.ID.String() != wantNotCreatedID,
		Message: "wanted not created ID %s; got %s",
		Args:    []any{wantNotCreatedID, s.Body.NotCreated.ID.String()},
	}, {
		Check:   gotMethod != wantMethod,
		Message: "wanted result method %s; got %s",
		Args:    []any{wantMethod, gotMethod},
	}, {
		Check:   s.ID.String() != wantID,
		Message: "wanted result ID %s; got %s",
		Args:    []any{wantID, s.ID.String()},
	}, {
		Check:   r.Errors[0].ID != wantError.ID,
		Message: "wanted error ID %s; got %s",
		Args:    []any{wantError.ID, r.Errors[0].ID},
	}, {
		Check:   r.Errors[0].Type != wantError.Type,
		Message: "wanted error type %s; got %s",
		Args:    []any{wantError.Type, r.Errors[0].Type},
	}, utils.NewCase(
		r.Errors[0].Properties[0] != wantError.Properties[0],
		"wanted error property %s; got %s",
		wantError.Properties[0], r.Errors[0].Properties[0],
	)}
	for _, c := range cases {
		if c.Check {
			t.Errorf(c.Message, c.Args...)
		}
	}
}
