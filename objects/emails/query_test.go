package emails_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/cwinters8/gomap/objects/emails"
	"github.com/cwinters8/gomap/utils"
)

func TestFilterJSON(t *testing.T) {
	timestamp := time.Now().UTC()
	filter := emails.Filter{
		InMailboxID: "xyz",
		Text:        "hello world",
		After:       &timestamp,
	}
	b, err := json.Marshal(filter)
	if err != nil {
		t.Fatalf("failed to marshal filter to json: %s", err.Error())
	}
	var m map[string]string
	if err := json.Unmarshal(b, &m); err != nil {
		t.Fatalf("failed to unmarshal filter to map: %s", err.Error())
	}
	if _, ok := m["before"]; ok {
		t.Error("wanted `before` field to be omitted")
	}
	after, err := time.Parse("2006-01-02T15:04:05.999999Z", m["after"])
	if err != nil {
		t.Fatalf("failed to parse `after` filter to time.Time: %s", err.Error())
	}
	cases := utils.Cases{utils.NewCase(
		m["inMailbox"] != filter.InMailboxID,
		"wanted mailbox id %s; got %s",
		filter.InMailboxID, m["inMailbox"],
	), utils.NewCase(
		m["text"] != filter.Text,
		"wanted text %s; got %s",
		filter.Text, m["text"],
	), utils.NewCase(
		after != *filter.After,
		"wanted after timestamp %v; got %v",
		*filter.After, after,
	)}
	cases.Iterator(func(c *utils.Case) {
		t.Error(c.Message)
	})
}
