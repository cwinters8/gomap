package requests_test

import (
	"encoding/json"
	"testing"

	"github.com/cwinters8/gomap/requests"
	"github.com/cwinters8/gomap/utils"

	"github.com/google/uuid"
)

func TestCallMarshal(t *testing.T) {
	id, err := uuid.NewRandom()
	if err != nil {
		t.Fatalf("failed to generate new uuid: %s", err.Error())
	}
	c := requests.Call{
		AccountID: "xyz",
		ID:        id,
		Method:    "Mailbox/query",
		Arguments: map[string]any{
			"hello": "world",
		},
		OnSuccess: func(b map[string]any) error {
			return nil
		},
	}
	b, err := json.Marshal(c)
	if err != nil {
		t.Fatalf("failed to marshal call to json: %s", err.Error())
	}
	var slice []any
	if err := json.Unmarshal(b, &slice); err != nil {
		t.Fatalf("failed to unmarshal call to map: %s", err.Error())
	}
	m, ok := slice[0].(string)
	if !ok {
		t.Fatalf("failed to cast method to string. %s", utils.Describe(slice[0]))
	}
	cases := utils.Cases{utils.NewCase(
		m != c.Method,
		"wanted method %s; got %s",
		c.Method, m,
	)}
	args, ok := slice[1].(map[string]any)
	if !ok {
		t.Fatalf("failed to cast args to map. %s", utils.Describe(slice[1]))
	}
	acctID, ok := args["accountId"].(string)
	if !ok {
		t.Fatalf("failed to cast account id to string. %s", utils.Describe(args["accountId"]))
	}
	hello, ok := args["hello"].(string)
	if !ok {
		t.Fatalf("failed to cast hello arg to string. %s", utils.Describe(args["hello"]))
	}
	cases.Append(utils.NewCase(
		acctID != c.AccountID,
		"wanted account id %s; got %s",
		c.AccountID, acctID,
	), utils.NewCase(
		hello != c.Arguments["hello"],
		"wanted hello message %s; got %s",
		c.Arguments["hello"], hello,
	))
	strID, ok := slice[2].(string)
	if !ok {
		t.Fatalf("failed to cast id to string. %s", utils.Describe(slice[2]))
	}
	gotID, err := uuid.Parse(strID)
	if err != nil {
		t.Fatalf("failed to parse `%s` as uuid: %s", strID, err.Error())
	}
	cases.Append(utils.NewCase(
		gotID != c.ID,
		"wanted id %s; got %s",
		c.ID, gotID,
	))
	cases.Iterator(func(c *utils.Case) {
		t.Error(c.Message)
	})
}
