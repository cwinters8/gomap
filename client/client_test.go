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
	switch {
	case client == nil:
		t.Fatalf("returned client should not be nil")
	case client.Session == nil:
		t.Error("session should not be nil")
	case client.HTTPClient == nil:
		t.Error("http client should not be nil")
	case client.Session.APIURL != wantURL:
		t.Errorf("wanted URL %s; got URL %s", wantURL, client.Session.APIURL)
	}
}
