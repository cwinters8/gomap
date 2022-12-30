package client_test

import (
	"os"
	"testing"

	"github.com/cwinters8/gomap/client"
	"github.com/cwinters8/gomap/utils"
)

func TestNewClient(t *testing.T) {
	if err := utils.Env("../.env"); err != nil {
		t.Fatalf("failed to load env: %s", err.Error())
	}
	client, err := client.NewClient(os.Getenv("FASTMAIL_SESSION_URL"), os.Getenv("FASTMAIL_TOKEN"))
	if err != nil {
		t.Fatalf("failed to instantiate new client: %s", err.Error())
	}
	wantURL := "https://api.fastmail.com/jmap/api/"

	// test cases that should cause the test to immediately fail
	fatals := []*utils.Case{{
		Check:   client == nil,
		Message: "returned client should not be nil",
	}, {
		Check:   client.Session == nil,
		Message: "session should not be nil",
	}}
	for _, c := range fatals {
		if c.Check {
			t.Fatalf(c.Message, c.Args...)
		}
	}

	errors := []*utils.Case{{
		Check:   client.HTTPClient == nil,
		Message: "http client should not be nil",
	}, {
		Check:   client.Session.APIURL != wantURL,
		Message: "wanted URL %s; got URL %s",
		Args:    []any{wantURL, client.Session.APIURL},
	}}
	for _, c := range errors {
		if c.Check {
			t.Errorf(c.Message, c.Args...)
		}
	}
}
